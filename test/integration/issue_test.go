package integration

import (
	"encoding/json"
	"testing"

	"github.com/iost-official/go-iost/core/tx"

	"github.com/iost-official/go-iost/ilog"

	"github.com/iost-official/go-iost/common"
	. "github.com/iost-official/go-iost/verifier"
	. "github.com/smartystreets/goconvey/convey"
)

func prepareIssue(s *Simulator, acc *TestAccount) (*tx.TxReceipt, error) {
	s.Head.Number = 0

	// deploy issue.iost
	setNonNativeContract(s, "issue.iost", "issue.js", ContractPath)
	s.Call("issue.iost", "init", `[]`, acc.ID, acc.KeyPair)

	witness := common.Witness{
		ID:      acc0.ID,
		Owner:   acc0.KeyPair.ID,
		Active:  acc0.KeyPair.ID,
		Balance: 123000,
	}
	params := []interface{}{
		acc0.ID,
		common.TokenInfo{
			FoundationAccount: acc1.ID,
			IOSTTotalSupply:   90000000000,
			IOSTDecimal:       8,
			RAMTotalSupply:    9000000000000000000,
			RAMGenesisAmount:  128,
		},
		[]interface{}{witness},
	}
	b, _ := json.Marshal(params)
	r, err := s.Call("issue.iost", "InitGenesis", string(b), acc.ID, acc.KeyPair)
	s.Visitor.Commit()
	return r, err
}

func Test_IOSTIssue(t *testing.T) {
	ilog.Stop()
	Convey("test issue.iost", t, func() {
		s := NewSimulator()
		defer s.Clear()

		createAccountsWithResource(s)
		r, err := prepareIssue(s, acc0)

		Convey("test init", func() {
			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(s.Visitor.TokenBalance("iost", acc0.ID), ShouldEqual, int64(123000*1e8))
		})

		Convey("test IssueIOST", func() {
			s.Head.Time += 4 * 3 * 1e9
			r, err := s.Call("issue.iost", "IssueIOST", `[]`, acc0.ID, acc0.KeyPair)
			s.Visitor.Commit()

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")

			So(s.Visitor.TokenBalance("iost", "bonus.iost"), ShouldEqual, int64(45654))
			So(s.Visitor.TokenBalance("iost", acc1.ID), ShouldEqual, int64(92691))
		})
	})
}
