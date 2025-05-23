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

### Client Configuration

The default configuration should work for most use cases, but you can customize it:

```go
// Create a custom configuration
config := fiber.NewConfig()
    .SetReadBufferSize(1024)
    .SetIdleTimeout(10 * time.Second) // Restart connections idle for 10 seconds
    .SetHealthCheckInterval(10 * time.Second) // Check connection health every 10 seconds

// Use the configuration with a client
client := fiber.NewClientWithConfig(endpoint, apiKey, config)
```

#### Available Configuration Options

- `EnableCompression()`: Enables gzip compression for all requests and responses
- `SetWriteBufferSize(size int)`: Sets the gRPC write buffer size
- `SetReadBufferSize(size int)`: Sets the gRPC read buffer size
- `SetConnWindowSize(size int32)`: Sets the gRPC connection window size
- `SetWindowSize(size int32)`: Sets the gRPC window size
- `SetIdleTimeout(timeout time.Duration)`: Sets a timeout after which idle connections will be restarted automatically. Set to 0 to disable (default).
- `SetHealthCheckInterval(interval time.Duration)`: Sets the interval for health checks. Set to 0 to disable (default).
- `SetLogLevel(level string)`: Sets the log level. Default is `panic` (which never logs).

### Subscriptions

You can find some examples on how to subscribe below. `fiber-go` uses the familiar `go-ethereum` core types where possible,
making it easy to integrate with existing applications.

#### Transactions

Transactions are returned as `*fiber.TransactionWithSender` which is a wrapper around `go-ethereum` `*types.Transaction` plus the sender's address.
The sender address is included in the message to avoid having to recompute it from ECDSA signature recovery in the client, which can be slow.

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

    ch := make(chan *fiber.TransactionWithSender)
    go func() {
        if err := client.SubscribeNewTxs(nil, ch); err != nil {
            log.Fatal(err)
        }
    }()

    for message := range ch {
        handleTransaction(message.Transaction)
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

    ch := make(chan *fiber.TransactionWithSender)
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

- To
- From
- MethodID
- Value (greater than, less than, equal to)

#### Execution Payloads (new blocks with transactions)

Execution payloads are returned as `*fiber.Block` which is a wrapper around `go-ethereum` native types such as `Header`, `Transaction` and `Withdrawal`.

```go
import (
    ...
    fiber "github.com/chainbound/fiber-go"
)

func main() {
    ...

    ch := make(chan *fiber.Block)

    go func() {
        if err := client.SubscribeNewExecutionPayloads(ch); err != nil {
            log.Fatal(err)
        }
    }()

    for block := range ch {
        handlePayload(block)
    }
}
```

#### Beacon Blocks

Beacon blocks follow the [Consensus specs](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#signedbeaconblock).
The returned items are `*fiber.BeaconBlock` which is a wrapper around `go-eth-2` `SignedBeaconBlock` depending on the hardfork version:

Each `*fiber.BeaconBlock` contains the `DataVersion` field which indicates the hardfork version of the beacon block.
The returned type will contain either a Bellatrix (3), Capella (4) or Deneb (5) hardfork block depending on the specified DataVersion.

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

#### Raw Beacon Blocks

Raw beacon blocks are raw, SSZ-encoded bytes that you can manually decode into [SignedBeaconBlocks](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#signedbeaconblock) in your application.

```go
import (
    ...
    fiber "github.com/chainbound/fiber-go"
)

func main() {
    ...

    ch := make(chan []byte)

    go func() {
        if err := client.SubscribeNewRawBeaconBlocks(ch); err != nil {
            log.Fatal(err)
        }
    }()

    for block := range ch {
        handleRawBeaconBlock(block)
    }
}
```

### Sending Transactions

#### `SendTransaction`

This method supports sending a single `go-ethereum` `*types.Transaction` object to the Fiber Network.

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

This method supports sending a sequence of transactions to the Fiber Network.

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

This method supports sending a single raw, RLP-encoded transaction to the Fiber Network.

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

This method supports sending a sequence of raw, RLP-encoded transactions to the Fiber Network.

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
