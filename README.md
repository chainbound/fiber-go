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
    "context"
    "log"
    "time"

    fiber "github.com/chainbound/fiber-go"
)

func main() {
    endpoint := "fiber.example.io"
    apiKey := "YOUR_API_KEY"
    client := fiber.NewClient(endpoint, apiKey)
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### Subscriptions
You can find some examples on how to subscribe below. `fiber-go` uses it's own
`Transaction` struct, which you can convert to a `go-ethereum` transaction using `tx.ToNative()`.
#### Transactions
```go
import (
    "context"
    "log"
    "time"
    
    fiber "github.com/chainbound/fiber-go"
    "github.com/chainbound/fiber-go/filter"
    "github.com/chainbound/fiber-go/protobuf/api"
)

func main() {
    endpoint := "fiber.example.io"
    apiKey := "YOUR_API_KEY"
    client := fiber.NewClient(endpoint, apiKey)
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }

    ch := make(chan *fiber.Transaction)
    go func() {
        if err := client.SubscribeNewTxs(nil, ch); err != nil {
            log.Fatal(err)
        }
    }()

    for tx := range ch {
        handleTransaction(tx)
    }
}
```

#### Filtering
The first argument to `SubscribeNewTxs` is a filter, which can be `nil` if you want to get all transactions.
A filter can be built with the `filter` package:
```go
import (
    ...
    "github.com/chainbound/fiber-go/filter"
)

func main() {
    ...

    // Construct filter
    // example 1: all transactions with either of these addresses as the receiver
    f := filter.New(filter.Or(
        filter.To("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
        filter.To("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"),
    ))

    // example 2: all transactions with a value greater than or equal to 1 ETH
    f := filter.New(filter.ValueGte(big.NewInt(1) * big.NewInt(1e18)))
    
    // example 3: all transactions with a value equal to 1 ETH
    f := filter.New(filter.ValueEq(big.NewInt(1) * big.NewInt(1e18)))
    
    // example 4: all transactions with a value less than or equal to 1 ETH
    f := filter.New(filter.ValueLte(big.NewInt(1) * big.NewInt(1e18)))

    // example 5: all ERC20 transfers on the 2 tokens below
    f := filter.New(filter.And(
        filter.MethodID("0xa9059cbb"),
        filter.Or(
            filter.To("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
            filter.To("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
        ),
    ))

    ch := make(chan *fiber.Transaction)
    go func() {
        // apply filter
        if err := client.SubscribeNewTxs(f, ch); err != nil {
            log.Fatal(err)
        }
    }()

    ...
}
```
You can currently filter the following properties
* To
* From
* MethodID
* Value (greater than, less than, equal to)
#### Execution Headers (new block headers)
```go
import (
    ...
    fiber "github.com/chainbound/fiber-go"
)

func main() {
    ...

    ch := make(chan *fiber.ExecutionPayloadHeader)

    go func() {
        if err := client.SubscribeNewExecutionPayloadHeaders(ch); err != nil {
            log.Fatal(err)
        }
    }()

    for header := range ch {
        handleHeader(block)        
    }
}
```

#### Execution Payloads (new blocks with transactions)
```go
import (
    ...
    fiber "github.com/chainbound/fiber-go"
)

func main() {
    ...

    ch := make(chan *fiber.ExecutionPayload)

    go func() {
        if err := client.SubscribeNewExecutionPayloads(ch); err != nil {
            log.Fatal(err)
        }
    }()

    for block := range ch {
        handleBlock(block)        
    }
}
```

#### Beacon Blocks
Beacon blocks follow the [Consensus specs](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconblock), with the exception of the `ExecutionPayload`, which is not included to
allow for a smaller payload size. Please use the `SubscribeNewExecutionPayloads` stream if you need it.

```go
import (
    ...
    fiber "github.com/chainbound/fiber-go"
)

func main() {
    ...

    ch := make(chan *fiber.BeaconBlock)

    go func() {
        if err := client.SubscribeNewBeaconBlocks(ch); err != nil {
            log.Fatal(err)
        }
    }()

    for block := range ch {
        handleBeaconBlock(block)
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
    endpoint := "fiber.example.io"
    apiKey := "YOUR_API_KEY"
    client := fiber.NewClient(endpoint, apiKey)
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

    hash, timestamp, err := client.SendTransaction(ctx, signed)
    if err != nil {
        log.Fatal(err)
    }

    doSomething(hash, timestamp)
}
```
#### `SendTransactionSequence`
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
    endpoint := "fiber.example.io"
    apiKey := "YOUR_API_KEY"
    client := fiber.NewClient(endpoint, apiKey)
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

    // type should be *types.Transaction (but signed, e.g. v,r,s fields filled in)
    target := someTargetTransaction

    hashes, timestamp, err := client.SendTransactionSequence(ctx, target, signed)
    if err != nil {
        log.Fatal(err)
    }

    doSomething(hashes, timestamp)
}
```
#### `SendRawTransaction`
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
    endpoint := "fiber.example.io"
    apiKey := "YOUR_API_KEY"
    client := fiber.NewClient(endpoint, apiKey)
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

    bytes, err := signed.MarshalBinary()
    if err != nil {
        log.Fatal(err)
    }

    hash, timestamp, err := client.SendRawTransaction(ctx, bytes)
    if err != nil {
        log.Fatal(err)
    }

    doSomething(hash, timestamp)
}
```

#### `SendRawTransactionSequence`
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
    endpoint := "fiber.example.io"
    apiKey := "YOUR_API_KEY"
    client := fiber.NewClient(endpoint, apiKey)
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

    bytes, err := signed.MarshalBinary()
    if err != nil {
        log.Fatal(err)
    }

    // Type should be []byte
    targetTransaction := someTargetTransaction

    hashes, timestamp, err := client.SendRawTransactionSequence(ctx, targetTransaction, bytes)
    if err != nil {
        log.Fatal(err)
    }

    doSomething(hashes, timestamp)
}
```
