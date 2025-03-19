package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"strings"
)

type DepositData struct {
	Pubkey                string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	Amount                uint64 `json:"amount"`
	Signature             string `json:"signature"`
	DepositMessageRoot    string `json:"deposit_message_root"`
	DepositDataRoot       string `json:"deposit_data_root"`
	ForkVersion           string `json:"fork_version"`
	NetworkName           string `json:"network_name"`
}

const depositContractABI = `[
	{
		"inputs": [
			{"name": "pubkey", "type": "bytes"},
			{"name": "withdrawal_credentials", "type": "bytes"},
			{"name": "signature", "type": "bytes"},
			{"name": "deposit_data_root", "type": "bytes32"}
		],
		"name": "deposit",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	}
]`

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please set private key as an argument")
	}
	privateKeyHex := os.Args[1]

	if !strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = "0x" + privateKeyHex
	}
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	jsonFile, err := os.ReadFile("deposit_data.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var depositDataList []DepositData
	if err := json.Unmarshal(jsonFile, &depositDataList); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	if len(depositDataList) == 0 {
		log.Fatal("File does not contain any deposit data")
	}

	depositData := depositDataList[0]

	client, err := ethclient.Dial("https://ethereum-holesky.publicnode.com")
	if err != nil {
		log.Fatalf("Error connecting to Holesky: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Error getting nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Error getting gas price: %v", err)
	}

	depositContractAddress := common.HexToAddress("0x4242424242424242424242424242424242424242")

	pubkeyBytes, err := hex.DecodeString(strings.TrimPrefix(depositData.Pubkey, "0x"))
	if err != nil {
		log.Fatalf("Error decoding pubkey: %v", err)
	}

	withdrawalCredentialsBytes, err := hex.DecodeString(strings.TrimPrefix(depositData.WithdrawalCredentials, "0x"))
	if err != nil {
		log.Fatalf("Error decoding withdrawal credentials: %v", err)
	}

	signatureBytes, err := hex.DecodeString(strings.TrimPrefix(depositData.Signature, "0x"))
	if err != nil {
		log.Fatalf("Error decoding signature: %v", err)
	}

	depositDataRootBytes, err := hex.DecodeString(strings.TrimPrefix(depositData.DepositDataRoot, "0x"))
	if err != nil {
		log.Fatalf("Error decoding deposit data root: %v", err)
	}

	var depositDataRoot [32]byte
	copy(depositDataRoot[:], depositDataRootBytes)

	parsedABI, err := abi.JSON(strings.NewReader(depositContractABI))
	if err != nil {
		log.Fatalf("Error parsing ABI: %v", err)
	}

	data, err := parsedABI.Pack("deposit", pubkeyBytes, withdrawalCredentialsBytes, signatureBytes, depositDataRoot)
	if err != nil {
		log.Fatalf("Error encoding function call: %v", err)
	}

	amount := big.NewInt(10000000000000000) // 0.01 ETH in wei

	gasLimit := uint64(500000)

	tx := types.NewTransaction(
		nonce,
		depositContractAddress,
		amount,
		gasLimit,
		gasPrice,
		data,
	)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Error getting chain ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Error signing transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Error sending transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	fmt.Println("Waiting for confirmation...")

	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		log.Fatalf("Error waiting for confirmation: %v", err)
	}

	fmt.Printf("Transaction confirmed in block: %d\n", receipt.BlockNumber)
	fmt.Printf("Transaction hash: %s\n", receipt.TxHash.Hex())
	fmt.Printf("Gas used: %d\n", receipt.GasUsed)
}
