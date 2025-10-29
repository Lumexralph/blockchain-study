package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, "https://cloudflare-eth.com")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	fmt.Println("Successfully connected to the Ethereum client")

	address := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")
	fmt.Printf("Address: Hex: %x, bytes: %v\n", address.Hex(), address.Bytes())

	blockNumber := big.NewInt(5532993)
	balance, err := client.BalanceAt(ctx, address, blockNumber)
	if err != nil {
		log.Fatalf("Failed to retrieve balance: %v", err)
	}
	fmt.Printf("Balance: %s wei\n", balance.String())

}
