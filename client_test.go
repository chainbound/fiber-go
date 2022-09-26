package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

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
		// if err := fiber.SubscribeNewTxs(&api.TxFilter{
		// 	To:       common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48").Bytes(),
		// 	MethodID: common.Hex2Bytes("7ff36ab5"),
		// }, ch); err != nil {
		// 	log.Println(err)
		// }
		if err := fiber.SubscribeNewTxs(nil, ch); err != nil {
			log.Println(err)
		}
	}()

	for tx := range ch {
		fmt.Println(tx.Hash())
		fmt.Println(tx.To())
	}
}
