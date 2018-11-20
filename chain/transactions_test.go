package chain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoinbaseOnly(t *testing.T) {
	createGenesisBlock(true)

	privKey, _ := ecdsa.GenerateKey(Curve, rand.Reader)
	pubKey := &privKey.PublicKey
	address := encodePublicKey(pubKey)
	cbTx := createCoinbaseTx(address)
	txs := []Transaction{cbTx}
	createBlock(txs, "My data", 0)

	assert.Equal(t, len(globalChain), 2, "Genesis block + single block with only coinbase transaction expected")

	valid := validateChain(globalChain)
	assert.True(t, valid, "Blockchain should be valid")
}

func TestCoinbaseAndUsers(t *testing.T) {
	createGenesisBlock(true)

	privKey1, _, address1 := createAddress()
	_, _, address2 := createAddress()

	cbTx := createCoinbaseTx(address1)

	txs := []Transaction{cbTx}
	createBlock(txs, "My data", 0)

	u1Tx := sendCoins(privKey1, address2, 50)

	/*	txIn := TxIn{cbTx.Id, 0}
		txIns := []TxIn{txIn}

		txOut1 := TxOut{address2, 50}
		txOut2 := TxOut{address1, 950}
		txOuts := []TxOut{txOut1, txOut2}


		u1Tx := createTx(privKey1, txIns, txOuts)*/

	txs = []Transaction{u1Tx}
	createBlock(txs, "My data", 0)

	assert.Equal(t, len(globalChain), 3, "Genesis block + two blocks expected")

	valid := validateChain(globalChain)
	assert.True(t, valid, "Blockchain should be valid")
	assert.Equal(t, 50, balanceFor(address2))
	assert.Equal(t, COINBASE_AMOUNT-50, balanceFor(address1))
}
