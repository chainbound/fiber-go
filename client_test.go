package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSub(t *testing.T) {
	key := "xxx"
	fiber := NewClient("localhost:8080", key)
	// fiber := NewClient("fiber-node.fly.dev:8080", "fiber/v0.0.2-alpha/28820807-4315-491c-bfbe-b38d9513b687")
	defer fiber.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := fiber.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	// ch := make(chan *types.Transaction)
	// go func() {
	// 	// if err := fiber.SubscribeNewTxs(&api.TxFilter{
	// 	// 	To:       common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48").Bytes(),
	// 	// 	MethodID: common.Hex2Bytes("7ff36ab5"),
	// 	// }, ch); err != nil {
	// 	// 	log.Println(err)
	// 	// }
	// 	if err := fiber.SubscribeNewTxs(nil, ch); err != nil {
	// 		log.Println(err)
	// 	}
	// }()
	start := time.Now()

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     1,
		To:        nil,
		Value:     big.NewInt(100),
		Gas:       21000,
		GasFeeCap: big.NewInt(21000),
		GasTipCap: big.NewInt(21000),
		Data:      nil,
	})

	pk, _ := crypto.GenerateKey()
	signer := types.NewLondonSigner(common.Big1)

	signed, err := types.SignTx(tx, signer, pk)
	if err != nil {
		log.Fatal(err)
	}

	start2 := time.Now()
	bytes, _ := signed.MarshalBinary()
	// hash, timestamp, err := fiber.BackrunTransaction(ctx, common.HexToHash("0xf98be8adcb57a4ebd50e21768d45d0fb131ca6a56c589c9721ebb5fc47804b96"), signed)
	// hash, timestamp, err := fiber.SendTransaction(ctx, signed)
	hash, timestamp, err := fiber.SendRawTransaction(ctx, bytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hash)
	fmt.Println(timestamp)
	fmt.Println("Sending took", time.Since(start2))
	fmt.Println("Total took", time.Since(start))

}
