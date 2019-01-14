package net

import (
	"encoding/json"
	"fmt"
	"github.com/mukatee/go-naive/chain"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestGetBlocks(t *testing.T) {
	chain.CreateTestChain(chain.GenesisAddress, 1)
	StartServer()
	time.Sleep(1)
	resp, err := http.Get("http://127.0.0.1:9090/blocks")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bytes2 := body
	rpcChain := []chain.Block{}
	json.Unmarshal(bytes2, &rpcChain)
	assert.Equal(t, 2, len(rpcChain))
	genesisBlock := rpcChain[0]
	AssertGenesisBlock(t, genesisBlock)
	testBlock := rpcChain[1]
	AssertTestBlock(t, 2, testBlock, genesisBlock)
}

func AssertTestBlock(t *testing.T, idx int, block, prevBlock chain.Block) {
	testBlock := block
	assert.Equal(t, 0, testBlock.Difficulty)
	assert.Equal(t, 0, testBlock.Nonce)
	assert.Equal(t, 1, len(testBlock.Transactions))
	cbTx := testBlock.Transactions[0]
	assertCoinbaseTx(t, cbTx)
	assert.Equal(t, idx, testBlock.Index)
	data := fmt.Sprintf("Test%d", idx)
	assert.Equal(t, data, testBlock.Data)
	assert.Equal(t, prevBlock.Hash, testBlock.PreviousHash)
	//not checking hash since it varies from block
	//	assert.Equal(t, "", testBlock.Hash)
}

func AssertGenesisBlock(t *testing.T, block chain.Block) {
	genesisBlock := block
	assert.Equal(t, 1, genesisBlock.Index)
	assert.Equal(t, "Teemu oli täällä", genesisBlock.Data)
	assert.Equal(t, "2c0bb21918368e747a9f0a1dbd6f059291c2d20c0e3a75954088210b77d76ad5", genesisBlock.Hash)
	assert.Equal(t, "0", genesisBlock.PreviousHash)
	assert.Equal(t, chain.GenesisTime, genesisBlock.Timestamp)
	assert.Equal(t, 1, genesisBlock.Difficulty)
	assert.Equal(t, 1, genesisBlock.Nonce)
	assert.Equal(t, 1, len(genesisBlock.Transactions))
	coinbase := genesisBlock.Transactions[0]
	assertCoinbaseTx(t, coinbase)
}

func assertCoinbaseTx(t *testing.T, cbTx chain.Transaction) {
	assert.Equal(t, 0, len(cbTx.TxIns))
	assert.Equal(t, 1, len(cbTx.TxOuts))
	txOut := cbTx.TxOuts[0]
	assert.Equal(t, chain.GenesisAddress, txOut.Address)
	assert.Equal(t, chain.COINBASE_AMOUNT, txOut.Amount)
	assert.Equal(t, "coinbase", cbTx.Signature)
	assert.Equal(t, "", cbTx.Sender)
}
