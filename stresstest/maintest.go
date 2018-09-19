package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/iost-official/Go-IOS-Protocol/account"
	"github.com/iost-official/Go-IOS-Protocol/common"
	"github.com/iost-official/Go-IOS-Protocol/core/tx"
	"github.com/iost-official/Go-IOS-Protocol/crypto"
	"github.com/iost-official/Go-IOS-Protocol/ilog"
	pb "github.com/iost-official/Go-IOS-Protocol/rpc"
	"google.golang.org/grpc"
)

var conns []*grpc.ClientConn

func initConn(num int) {
	conns = make([]*grpc.ClientConn, num)
	allServers := []string{"54.88.65.72:30002", "18.223.226.249:30002", "54.67.42.228:30002", "52.88.239.19:30002", "13.232.151.244:30002"}
	for i := 0; i < num; i++ {
		conn, err := grpc.Dial(allServers[i%len(allServers)], grpc.WithInsecure())
		if err != nil {
			continue
		}
		conns[i] = conn
	}
}

func transParallel(num int) {
	if conns == nil {
		initConn(num)
	}
	wg := new(sync.WaitGroup)
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			transfer(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func sendTx(stx *tx.Tx, i int) ([]byte, error) {
	if conns[i] == nil {
		return nil, errors.New("nil conn")
	}
	client := pb.NewApisClient(conns[i])
	resp, err := client.SendRawTx(context.Background(), &pb.RawTxReq{Data: stx.Encode()})
	if err != nil {
		return nil, err
	}
	return []byte(resp.Hash), nil
}

func loadBytes(s string) []byte {
	if s[len(s)-1] == 10 {
		s = s[:len(s)-1]
	}
	buf := common.Base58Decode(s)
	return buf
}

func transfer(i int) {
	action := tx.NewAction("iost.system", "Transfer", `["IOST2g5LzaXkjAwpxCnCm29HK69wdbyRKbfG4BQQT7Yuqk57bgTFkY","IOST22TgXbjvgrDd3DuMkVufcWbYDMy7vpmQoCgZXmgi8eqM7doxWD",1]`)
	acc, _ := account.NewAccount(loadBytes("319xGCaLZP5D4sAVCEX4LDAMgzaZ3LJiXgCVxB8y1igTmUCkHj6DJRCH4C8myor1P3rZHttFneApzznHqvqqTpiu"), crypto.Ed25519)
	trx := tx.NewTx([]*tx.Action{&action}, [][]byte{}, 1000, 1, time.Now().Add(time.Second*time.Duration(10000)).UnixNano())
	stx, err := tx.SignTx(trx, acc)
	if err != nil {
		fmt.Println("signtx", stx, err)
		return
	}

	var txHash []byte
	txHash, err = sendTx(stx, i)
	if err != nil {
		fmt.Println("sendtx", txHash, err)
		return
	}
}

func main() {
	go transParallel(5000)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	i := <-c
	ilog.Infof("IOST server received interrupt[%v], shutting down...", i)
	for i := 0; i < len(conns); i++ {
		conns[i].Close()
	}
	ilog.Error("close the connections")

}
