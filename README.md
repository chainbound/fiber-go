# `fiber-go`
Go client package for interacting with Fiber Network.

## Installation
```bash
go get github.com/chainbound/fiber-go
```

## Usage
### Connecting
```go
import (
    fiber "github.com/chainbound/fiber-go"
)

func main() {
    client := fiber.NewClient("YOUR_API_HERE")
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### Subscriptions
You can find some examples on how to subscribe to the various subscriptions below. `fiber-go` works with
`go-ethereum` Transactions internally, but uses a protobuf block representation. This will be updated soon.
#### Transactions
```go
import (
    "context"
    "log"
    "time"
    
    fiber "github.com/chainbound/fiber-go"
    "github.com/chainbound/fiber-go/protobuf/api"

    "github.com/ethereum/go-ethereum/core/types"
)

func main() {
    client := fiber.NewClient("YOUR_API_HERE")
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }

    ch := make(chan *types.Transaction)
    go func() {
        if err := client.SubscribeNewTxs(&api.TxFilter{}, ch); err != nil {
            log.Fatal(err)
        }
    }()

    for tx := range ch {
        handleTransaction(tx)
    }
}
```
#### Blocks
```go
import (
    "context"
    "log"
    "time"
    
    fiber "github.com/chainbound/fiber-go"
    "github.com/chainbound/fiber-go/protobuf/api"
    "github.com/chainbound/fiber-go/protobuf/eth"
)

func main() {
    client := fiber.NewClient("YOUR_API_HERE")
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }

    ch := make(chan *eth.Block)
    go func() {
        if err := client.SubscribeNewBlocks(&api.BlockFilter{}, ch); err != nil {
            log.Fatal(err)
        }
    }()

    for tx := range ch {
        handleBlock(tx)
    }
}
```

### Sending Transactions
#### `SendTransaction`
```go
import (
    "context"
    "log"
    "math/big"
    "time"

    fiber "github.com/chainbound/fiber-go"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"

)

func main() {
    client := fiber.NewClient("YOUR_API_HERE")
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }

    // Example transaction
    tx := types.NewTx(&types.DynamicFeeTx{
        Nonce:     nonce,
        To:        common.HexToAddress("0x...."),
        Value:     big.NewInt(100),
        Gas:       21000,
        GasFeeCap: big.NewInt(x),
        GasTipCap: big.NewInt(y),
        Data:      nil,
    })

    pk, _ := crypto.HexToECDSA("PRIVATE_KEY")
    signer := types.NewLondonSigner(common.Big1)

    signed, err := types.SignTx(tx, signer, pk)
    if err != nil {
        log.Fatal(err)
    }

    hash, timestamp, err := client.SendTransaction(ctx, tx)
    if err != nil {
        log.Fatal(err)
    }

    doSomething(hash, timestamp)
}
```
#### `BackrunTransaction`
```go
import (
    "context"
    "log"
    "math/big"
    "time"

    fiber "github.com/chainbound/fiber-go"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"

)

func main() {
    client := fiber.NewClient("YOUR_API_HERE")
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }

    // Example transaction
    tx := types.NewTx(&types.DynamicFeeTx{
        Nonce:     nonce,
        To:        common.HexToAddress("0x...."),
        Value:     big.NewInt(100),
        Gas:       21000,
        GasFeeCap: big.NewInt(x),
        GasTipCap: big.NewInt(y),
        Data:      nil,
    })

    pk, _ := crypto.HexToECDSA("PRIVATE_KEY")
    signer := types.NewLondonSigner(common.Big1)

    signed, err := types.SignTx(tx, signer, pk)
    if err != nil {
        log.Fatal(err)
    }

    // type should be common.Hash
    target := someTargetTransactionHash

    hash, timestamp, err := client.BackrunTransaction(ctx, target, tx)
    if err != nil {
        log.Fatal(err)
    }

    doSomething(hash, timestamp)
}
```
