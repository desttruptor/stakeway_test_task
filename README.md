# Validator Deposit Script for Holesky Testnet

This script allows you to submit a validator deposit to the staking deposit contract on the Holesky testnet.

## Prerequisites

1. Go 1.16 or newer
2. Ethereum private key with Holesky ETH balance
3. `deposit_data.json` file in the same directory as the script

## Installing Dependencies

From the project root directory:

```bash
# Navigate to the project root if you're not already there
cd stakeway_test_task

# Initialize the Go module (if not already done)
go mod init stakeway_test_task

# Get dependencies
go get github.com/ethereum/go-ethereum
```

## Usage

1. Get test ETH from the Holesky faucet: https://holesky-faucet.pk910.de/#/
2. Make sure your `deposit_data.json` file is in the project's root directory
3. Run the script with your private key (with or without the 0x prefix):

```bash
# From the project root
go run cmd/deposit/main.go YOUR_PRIVATE_KEY

# OR navigate to the script directory and run
cd cmd/deposit
go run main.go YOUR_PRIVATE_KEY
```

## How It Works

The script performs the following steps:

1. Loads your private key and the validator deposit data from deposit_data.json
2. Connects to the Holesky testnet through a public RPC endpoint
3. Encodes the deposit function call with the validator data
4. Creates, signs, and sends the transaction to the deposit contract
5. Waits for transaction confirmation and displays the transaction hash

## Example Output

```
Transaction sent: 0x3a9273d7e0e30e63668725c6c9bd25a39985dcd8e38b2195f95a2fdc6e34b03e
Waiting for confirmation...
Transaction confirmed in block: 123456
Transaction hash: 0x3a9273d7e0e30e63668725c6c9bd25a39985dcd8e38b2195f95a2fdc6e34b03e
Gas used: 321000
```

## Notes

- The deposit contract address is set to 0x4242424242424242424242424242424242424242 as specified in the requirements
- Make sure to place the `deposit_data.json` file in the project root directory