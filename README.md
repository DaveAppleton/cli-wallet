# Spacemesh CLI Wallet

## Overview
A basic reference Spacemesh CLI Wallet designed to work together with a [go-spacemesh full node](https://github.com/spacemeshos/go-spacemesh).

The CLI Wallet is designed for developers who want to hack on the Spacemesh platform. For non-devs we recommend using Smapp - the Spacemesh App. [Smapp](https://github.com/spacemeshos/smapp) is available for all major desktop platforms and includes a wallet.

## Functionality
The CLI wallet implements this [mini spec](https://github.com/spacemeshos/product/blob/master/cli_wallet_spec.md). See below on how to use with any network.

## Building

### Build for your current platform with go:

```bash
go get && go build
```

### Build for all platforms:
```bash
make
```

### Build for a specific platforms:
```bash
make build-win
```

```bash
make build-mac
```

```bash
make build-linux
```

### With `docker`:
```
make dockerbuild-go
```
---

## Using with a local full node and an open testnet
1. Follow the instructions in the readme to build and run [go-spacemesh](https://github.com/spacemeshos/go-spacemesh) from source code.

2. Build the CLI Wallet from this repository and run it:

```bash
./cli_wallet
```



---

## Using with a Local testnet
Please follow the instructions in our [local testnet guide](https://testnet.spacemesh.io/#/local)
