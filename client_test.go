package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/chainbound/fiber-go/protobuf/api"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestSub(t *testing.T) {
	key := "xxx"
	fiber := NewClient("localhost:8080", key)
	defer fiber.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := fiber.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	ch := make(chan *types.Transaction)
	go func() {
		if err := fiber.SubscribeNewTxs(&api.TxFilter{
			To: common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D").Bytes(),
		}, ch); err != nil {
			log.Println(err)
		}
	}()

	for tx := range ch {
		fmt.Println(tx.Hash())
		fmt.Println(tx.To())
	}
}
