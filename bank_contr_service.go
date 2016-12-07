package main

import (
	_ "encoding/hex"
	"eth-utils/contract"
	"fmt"
	"math/big"
	"time"

	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BankContractService struct {
	key             *ecdsa.PrivateKey
	contractAddress common.Address
	client          *ethclient.Client
}

/*
	Constructor
*/
func InitBankContractService(privateKey string, contractAddress string, client *ethclient.Client) (*BankContractService, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err == nil {
		var address common.Address
		if contractAddress != "" {
			address = common.HexToAddress(contractAddress)
		}
		return &BankContractService{key, address, client}, nil
	} else {
		return nil, err
	}

}

func (self *BankContractService) deployContract() {
	txOpt := bind.NewKeyedTransactor(self.key)
	txOpt.GasLimit = big.NewInt(250000)
	fmt.Println("TX: ", txOpt)
	//Deploy a new contract
	triggerAddr, tx, trigger, err := contract.DeployBank(txOpt, self.client)
	if err == nil {
		fmt.Println("Contract address: ", triggerAddr)
		fmt.Println("tx:", tx.Hash())
		fmt.Println("TX fields: ", tx)
		fmt.Println("Addr:", triggerAddr)
		fmt.Println("Trigger obj:", trigger)

	} else {
		fmt.Println("Failed to deploy new trigger contract: %v", err)
	}
}

func (self *BankContractService) deposit() {
	ctxParent := context.Background()
	ctx, cancel := context.WithTimeout(ctxParent, time.Second*3)
	defer cancel()

	addressNonce, err := self.client.PendingNonceAt(ctx, crypto.PubkeyToAddress(self.key.PublicKey))
	if err != nil {
		fmt.Println("Err get nonce", err)
	}
	fmt.Println("Address nonce", addressNonce)

	gasPrice, err := self.client.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Println("Err get gasPrice", err)
		panic(err)
	}
	fmt.Println("GasPrice", gasPrice)

	amount := big.NewInt(200000)
	gasLimit := big.NewInt(85000)

	tx := types.NewTransaction(addressNonce, self.contractAddress, amount, gasLimit, gasPrice, nil)
	txSigned, err := tx.SignECDSA(types.HomesteadSigner{}, self.key)
	if err != nil {
		fmt.Println("Err sign", err)
		panic(err)
	}
	err = self.client.SendTransaction(ctx, txSigned)
	if err == nil {
		fmt.Println(tx)
	} else {
		fmt.Println("Err send tx", err)
		panic(err)
	}
}

func (self *BankContractService) withdraw() {
	txOpt := bind.NewKeyedTransactor(self.key)
	txOpt.GasLimit = big.NewInt(100000)
	fmt.Println("txOpt:", txOpt)

	transactor, err := contract.NewBankTransactor(self.contractAddress, self.client)
	if err != nil {
		fmt.Println("Failed to create transactor: %v", err)
	}
	tx, err := transactor.Withdraw(txOpt)
	if err != nil {
		fmt.Println("tx err", err)
	}
	fmt.Println("TX", tx)
}

func (self *BankContractService) getClientsCount() uint64 {
	caller, err := contract.NewBankCaller(self.contractAddress, self.client)
	if err != nil {
		fmt.Println("Failed to instantiate contract: %v", err)
	}

	//call contract without tx get count
	count, err := caller.GetCount(nil)
	if err != nil {
		fmt.Println("Err 'GetCount' :", err)
		panic(err)
	}

	fmt.Println("Clients count:", count)
	return count
}
