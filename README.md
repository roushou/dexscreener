# Dex Screener

This repository is a Go library to interact with [Dex Screener API](https://docs.dexscreener.com/api/reference).

## Features

- Token Profiles: Retrieve the latest token profiles.
- Boosted Tokens: Get lists of both latest and most active boosted tokens.
- Token Orders: Fetch orders for specific tokens by chain ID and token address.
- Token Pairs: Query for token pair details including liquidity, volume, and price changes.

## Installation

```sh
go get github.com/roushou/dexscreener
```

## Usage

Here's how to use this library:

```go
package main

import "github.com/roushou/dexscreener"

func main() {
    client := dexscreener.NewClient()

    // Fetch latest boosted tokens
    boosted, err := client.GetLatestBoostedTokens()
    if err != nil {
        log.Fatalf("Failed to get boosted tokens: %v", err)
    }
    fmt.Printf("Number of boosted tokens: %d\n", len(boosted))

    // Fetch token pairs
    pairs, err := client.GetTokenPairs("JUPyiwrYJFskUPiHa7hkeR8VUtAeFoSYbKedZNsDvCN")
    if err != nil {
        log.Fatalf("Failed to get token pairs: %v", err)
    }
    fmt.Printf("Number of pairs for token: %d\n", len(pairs.Pairs))
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
