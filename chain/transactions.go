package chain

import (
	"strings"
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"math/big"
	"github.com/akamensky/base58"
)

var COINBASE_AMOUNT = 1000

type TxOut struct {
	Address string	//receiving public key
	Amount int		//amount of coin units to send/receive
}

type TxIn struct {
	TxId      string	//id of the transaction inside which this TxIn should be found (in the list of TxOut)
	TxIdx     int		//index of TxOut this refers to inside the transaction
}

type UnspentTxOut struct {
	TxId    string	//transaction id
	TxIdx   int		//index of txout in transaction
	Address string  //public key of owner
	Amount  int		//amount coin units that was sent/received
}

//transaction is from a single person/entity
//inputs from that entity use funds, outputs show how much and where
//all the time there should be some list kept and updated based on this
type Transaction struct {
	Id string
	Sender string //the address/public key of the sender
	Signature string	//signature including all txin and txout. in this case we sign Transaction.Id since that already has all the TxIn and TxOut
	TxIns []TxIn
	TxOuts []TxOut
}

type ecdsaKeyElements struct {
	R, S *big.Int
}

//build a string with all transaction inputs and ouputs concatenated, hash it, and encode the hash into a hex-string
func calculateTxId(tx Transaction) string {
	var strBuilder strings.Builder
	for _, txIn := range tx.TxIns {
		fmt.Fprintf(&strBuilder, "%d%d", txIn.TxIdx, txIn.TxIdx)
	}
	for _, txOut := range tx.TxOuts {
		fmt.Fprintf(&strBuilder, "%s%d", txOut.Address, txOut.Amount)
	}
	raw := strBuilder.String()
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func signData(privKey *ecdsa.PrivateKey, msg []byte) string {
	esig := createSignature(msg, privKey)
	rBytes := esig.R.Bytes()
	sBytes := esig.S.Bytes()
	whole := mergeTwoByteSlices(rBytes, sBytes)
	return base58.Encode(whole)
}

func createCoinbaseTx(address string) Transaction {
	var cbTx Transaction

	var txIn TxIn
	txIn.TxIdx = len(globalChain)
	cbTx.TxIns = append(cbTx.TxIns, txIn)

	var txOut TxOut
	txOut.Amount = COINBASE_AMOUNT
	txOut.Address = address
	cbTx.TxOuts = append(cbTx.TxOuts, txOut)

	cbTx.Id = calculateTxId(cbTx)
	cbTx.Signature = "coinbase"

	return cbTx
}

func sendCoins(privKey *ecdsa.PrivateKey, to string, count int) Transaction {
	from := encodePublicKey(privKey.PublicKey)
	//TODO: error handling
	txIns, total := findTxInsFor(from, count)
	txOuts := splitTxIns(from, to, count, total)
	tx := createTx(privKey, txIns, txOuts)
	return tx
}

func createTx(privKey *ecdsa.PrivateKey, txIns []TxIn, txOuts []TxOut) Transaction {
	pubKey := encodePublicKey(privKey.PublicKey)
	tx := Transaction{"",  pubKey,"", txIns, txOuts}

	signTxIns(tx, privKey)

	tx.Id = calculateTxId(tx)
	tx.Signature = signData(privKey, []byte(tx.Id))
	return tx
}

func processTransaction(tx Transaction) {
	//TODO: check that tx is valid, and all txin and txout are valid as well
	//first param from range is index, second is the value
	for _, val := range tx.TxIns {
		consumeTxOut(val.TxId, val.TxIdx)
	}
	for idx, val := range tx.TxOuts {
		createTxOut(tx.Id, val.Amount, val.Address, idx)
	}
}

func consumeTxOut(txId string, txIdx int) {
	for idx, val := range unspentTxOuts {
		if val.TxId == txId && val.TxIdx == txIdx {
			//remove the matching unspend txout from the list of unspents
			//unpacking with ... https://www.reddit.com/r/golang/comments/1y8ytg/what_mean_in_append_function_when_appending_byte/
			unspentTxOuts = append(unspentTxOuts[:idx], unspentTxOuts[idx+1:]...)
			return
		}
	}
}

func createTxOut(txId string, amount int, pubKey string, txIdx int) {
	newUTXO := UnspentTxOut{txId, txIdx, pubKey, amount}
	unspentTxOuts = append(unspentTxOuts, newUTXO)
}


func signTxIns(tx Transaction, privKey *ecdsa.PrivateKey) bool {
	//key from string https://stackoverflow.com/questions/48392334/how-to-sign-a-message-with-an-ecdsa-string-privatekey
	myAddress := encodePublicKey(privKey.PublicKey)
	//first param from range is index, second is the value
	for _, val := range tx.TxIns {
		errorStatus := false
		internalIdx := findUnspentTransaction(myAddress, val.TxId)
		if internalIdx < 0 {
			println("trying to spend a transaction that does not exist (as unspent..)")
			errorStatus = true
		}
		if errorStatus {
			return false
		}
	}
	return true
}

func hexToPublicKey(xHex string, yHex string) *ecdsa.PublicKey {
	xBytes, _ := hex.DecodeString(xHex)
	x := new(big.Int)
	x.SetBytes(xBytes)

	yBytes, _ := hex.DecodeString(yHex)
	y := new(big.Int)
	y.SetBytes(yBytes)

	pub := new(ecdsa.PublicKey)
	pub.X = x
	pub.Y = y

	pub.Curve = Curve

	return pub
}
