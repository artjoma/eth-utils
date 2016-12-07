package main

import (
	"fmt"
	"math/big"
	"time"

	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	ETH_DAEMON_URL       = "http://127.0.0.1:8546"
	PRIVATE_KEY_HEX      = "f293dedfa929af4595591f4f2c46c6d891ab7389d003936d31e2952379023a57"
	CONTRACT_ADDRESS_HEX = "bfa218c460ff5a9bc5ab1261e986eb1fc0b69e32"
)

/*
	Tested on Geth 1.5.3/4
*/
func main() {
	//create http rpc client
	client, err := ethclient.Dial(ETH_DAEMON_URL)
	ctxParent := context.Background()
	// The request has a timeout, so create a context that is
	// canceled automatically when the timeout expires.
	ctx, cancel := context.WithTimeout(ctxParent, time.Second*3)
	defer cancel()

	if err == nil {
		fmt.Println(client)
	} else {
		fmt.Println("Can't ctreate client")
		panic(err)
	}
	//create test request
	block, err := client.BlockByNumber(ctx, big.NewInt(2))
	if err != nil {
		fmt.Println("Can't get block by number")
		panic(err)
	}
	fmt.Println("Block: ", block.Number(), err)

	//step 1
	service, err := InitBankContractService(PRIVATE_KEY_HEX, "", client)
	if err == nil {
		service.deployContract()
	} else {
		fmt.Println("Err create contract binding instance", err)
		panic(err)
	}

	/*
		//step 2
		service, err := InitBankContractService(PRIVATE_KEY_HEX, CONTRACT_ADDRESS_HEX, client)
		if err == nil {
			service.deposit()
		} else {
			fmt.Println("Err create contract binding instance", err)
			panic(err)
		}

		//step 3
		service, err = InitBankContractService(PRIVATE_KEY_HEX, CONTRACT_ADDRESS_HEX, client)
		if err == nil {
			fmt.Println("Count:", service.getClientsCount())
		} else {
			fmt.Println("Err create contract binding instance", err)
			panic(err)
		}

		//step 4
		service, err = InitBankContractService(PRIVATE_KEY_HEX, CONTRACT_ADDRESS_HEX, client)
		if err == nil {
			service.withdraw()
		} else {
			fmt.Println("Err create contract binding instance", err)
			panic(err)
		}
	*/
}
