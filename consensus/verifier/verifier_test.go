package verifier

import (
	"fmt"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/iost-official/Go-IOS-Protocol/account"
	"github.com/iost-official/Go-IOS-Protocol/common"
	"github.com/iost-official/Go-IOS-Protocol/core/block"
	"github.com/iost-official/Go-IOS-Protocol/core/tx"
	"github.com/iost-official/Go-IOS-Protocol/crypto"
	"github.com/iost-official/Go-IOS-Protocol/db"
	"github.com/iost-official/Go-IOS-Protocol/vm"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVerifyBlockHead(t *testing.T) {
	Convey("Test of verify block head", t, func() {
		parentBlk := &block.Block{
			Head: &block.BlockHead{
				Number: 3,
				Time:   common.GetCurrentTimestamp().Slot - 1,
			},
		}
		chainTop := &block.Block{
			Head: &block.BlockHead{
				Number: 1,
				Time:   common.GetCurrentTimestamp().Slot - 4,
			},
		}
		hash := parentBlk.HeadHash()
		blk := &block.Block{
			Head: &block.BlockHead{
				ParentHash: hash,
				Number:     4,
				Time:       common.GetCurrentTimestamp().Slot,
				TxsHash:    common.Sha3([]byte{}),
				MerkleHash: []byte{},
			},
		}
		convey.Convey("Pass", func() {
			err := VerifyBlockHead(blk, parentBlk, chainTop)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Wrong time", func() {
			blk.Head.Time = common.GetCurrentTimestamp().Slot - 5
			err := VerifyBlockHead(blk, parentBlk, chainTop)
			convey.So(err, convey.ShouldEqual, errOldBlk)
			blk.Head.Time = common.GetCurrentTimestamp().Slot + 2
			err = VerifyBlockHead(blk, parentBlk, chainTop)
			convey.So(err, convey.ShouldEqual, errFutureBlk)
		})

		convey.Convey("Wrong parent", func() {
			blk.Head.ParentHash = []byte("fake hash")
			err := VerifyBlockHead(blk, parentBlk, chainTop)
			convey.So(err, convey.ShouldEqual, errParentHash)
		})

		convey.Convey("Wrong number", func() {
			blk.Head.Number = 5
			err := VerifyBlockHead(blk, parentBlk, chainTop)
			convey.So(err, convey.ShouldEqual, errNumber)
		})

		convey.Convey("Wrong tx hash", func() {
			tx0 := tx.NewTx(nil, nil, 1000, 1, 300)
			blk.Txs = append(blk.Txs, tx0)
			blk.Head.TxsHash = blk.CalculateTxsHash()
			err := VerifyBlockHead(blk, parentBlk, chainTop)
			convey.So(err, convey.ShouldBeNil)
			blk.Head.TxsHash = []byte("fake hash")
			err = VerifyBlockHead(blk, parentBlk, chainTop)
			convey.So(err, convey.ShouldEqual, errTxHash)
		})
	})
}

func loadBytes(s string) []byte {
	if s[len(s)-1] == 10 {
		s = s[:len(s)-1]
	}
	buf := common.Base58Decode(s)
	return buf
}

func transfer() *tx.Tx {
	action := tx.NewAction("iost.system", "Transfer", `["IOST5FhLBhVXMnwWRwhvz5j9NyWpBSchAMzpSMZT21xZqT8w7icwJ5","IOSTponZK9JJqZAsEWMF1BCZkSKnRP7abGbKjZb49nidfYW8",1]`)
	acc, _ := account.NewAccount(loadBytes("BCV7fV37aSWNx1N1Yjk3TdQXeHMmLhyqsqGms1PkqwPT"), crypto.Secp256k1)
	// fmt.Println(acc.Pubkey, account.GetIDByPubkey(acc.Pubkey))
	trx := tx.NewTx([]*tx.Action{&action}, [][]byte{}, 1000, 1, 0)
	stx, err := tx.SignTx(trx, acc)
	//fmt.Println("verify", stx.VerifySelf())
	if err != nil {
		return nil
	}
	return stx
}

var (
	WitnessList = []string{
		"IOST5FhLBhVXMnwWRwhvz5j9NyWpBSchAMzpSMZT21xZqT8w7icwJ5",
		"IOST6Jymdka3EFLAv8954MJ1nBHytNMwBkZfcXevE2PixZHsSrRkbR",
		"IOST7gKuvHVXtRYupUixCcuhW95izkHymaSsgKTXGDjsyy5oTMvAAm",
	}
)

func BenchmarkVerifier(b *testing.B) {
	var acts []*tx.Action
	for i := 0; i < len(WitnessList); i++ {
		act := tx.NewAction("iost.system", "IssueIOST", fmt.Sprintf(`["%v", %v]`, WitnessList[i], 10000000000))
		acts = append(acts, &act)
	}
	trx := tx.NewTx(acts, nil, 0, 0, 0)
	trx.Time = 0
	acc, _ := account.NewAccount(common.Base58Decode("2vj2Ab8Taz1TT2MSQHxmSffGnvsc9EVrmjx1W7SBQthCpuykhbRn2it8DgNkcm4T9tdBgsue3uBiAzxLpLJoDUbc"), crypto.Ed25519)
	blockHead := block.BlockHead{
		Version:    0,
		ParentHash: nil,
		Number:     0,
		Witness:    acc.ID,
		Time:       0,
	}
	stateDB, _ := db.NewMVCCDB("StateDB")
	defer func() {
		stateDB.Close()
		os.RemoveAll("stateDB")
	}()
	engine := vm.NewEngine(&blockHead, stateDB)
	engine.Exec(trx)
	hash, _ := blockHead.Hash()
	stateDB.Tag(string(hash))
	var blk block.Block
	blk.Head = &block.BlockHead{
		Version:    0,
		ParentHash: nil,
		Number:     1,
		Witness:    acc.ID,
		Time:       1,
	}
	for i := 0; i < 6000; i++ {
		trans := transfer()
		blk.Txs = append(blk.Txs, trans)
		receipt, _ := engine.Exec(trans)
		blk.Receipts = append(blk.Receipts, receipt)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := VerifyBlockWithVM(&blk, stateDB)
		if err != nil {
			fmt.Println(err)
		}
	}
}
