package chain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/mukatee/go-naive/cryptoff"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoinbaseOnly(t *testing.T) {
	createGenesisBlock(true)

	privKey, _ := ecdsa.GenerateKey(cryptoff.Curve, rand.Reader)
	pubKey := &privKey.PublicKey
	address := cryptoff.EncodePublicKey(pubKey)
	cbTx := CreateCoinbaseTx(address)
	txs := []Transaction{cbTx}
	CreateBlock(txs, "My data", 0)

	assert.Equal(t, len(GlobalChain), 2, "Genesis block + single block with only coinbase transaction expected")

	valid := validateChain(GlobalChain)
	assert.True(t, valid, "Blockchain should be valid")
}

func TestCoinbaseAndUsers(t *testing.T) {
	createGenesisBlock(true)

	privKey1, _, address1 := cryptoff.CreateAddress()
	_, _, address2 := cryptoff.CreateAddress()

	cbTx := CreateCoinbaseTx(address1)

	txs := []Transaction{cbTx}
	CreateBlock(txs, "My data", 0)

	u1Tx := SendCoins(privKey1, address2, 50)

	/*	txIn := TxIn{cbTx.Id, 0}
		txIns := []TxIn{txIn}

		txOut1 := TxOut{address2, 50}
		txOut2 := TxOut{address1, 950}
		txOuts := []TxOut{txOut1, txOut2}


		u1Tx := createTx(privKey1, txIns, txOuts)*/

	txs = []Transaction{u1Tx}
	CreateBlock(txs, "My data", 0)

	assert.Equal(t, len(GlobalChain), 3, "Genesis block + two blocks expected")

	valid := validateChain(GlobalChain)
	assert.True(t, valid, "Blockchain should be valid")
	assert.Equal(t, 50, BalanceFor(address2))
	assert.Equal(t, COINBASE_AMOUNT-50, BalanceFor(address1))
}
