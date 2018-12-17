package chain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/akamensky/base58"
	"github.com/mukatee/go-naive/cryptoff"
	"log"
	"math/big"
	"strings"
)

var COINBASE_AMOUNT = 1000

type TxOut struct {
	Address string //receiving public key
	Amount  int    //amount of coin units to send/receive
}

type TxIn struct {
	TxId  string //id of the transaction inside which this TxIn should be found
	TxIdx int    //index of TxOut this refers to inside the transaction
}

type UnspentTxOut struct {
	TxId    string //transaction id
	TxIdx   int    //index of txout in transaction
	Address string //public key of owner
	Amount  int    //amount coin units that was sent/received
}

//transaction is from a single person/entity
//inputs from that entity use funds, outputs show how much and where
//all the time there should be some list kept and updated based on this
type Transaction struct {
	Id        string
	Sender    string //the address/public key of the sender
	Signature string //signature including all txin and txout. in this case we sign Transaction.Id since that already has all the TxIn and TxOut
	TxIns     []TxIn
	TxOuts    []TxOut
}

type ecdsaKeyElements struct {
	R, S *big.Int
}

//build a string with all transaction inputs and ouputs concatenated, hash it, and encode the hash into a hex-string
func calculateTxId(tx Transaction) string {
	log.Print("Calculating tx id (hash string)")
	var strBuilder strings.Builder
	for _, txIn := range tx.TxIns {
		fmt.Fprintf(&strBuilder, "%d%d", txIn.TxIdx, txIn.TxIdx)
	}
	for _, txOut := range tx.TxOuts {
		fmt.Fprintf(&strBuilder, "%s%d", txOut.Address, txOut.Amount)
	}
	raw := strBuilder.String()
	hash := sha256.Sum256([]byte(raw))
	id := hex.EncodeToString(hash[:])
	log.Print("Calculated tx id (hash string):", id)
	return id
}

//signData creates an ECDSA signature for given message (byte slice).
//the created signature is returned as base 58 encoded
func signData(privKey *ecdsa.PrivateKey, msg []byte) string {
	log.Print("Creating ECDSA signature for data size ", len(msg))
	esig := cryptoff.CreateSignature(msg, privKey)
	rBytes := esig.R.Bytes()
	sBytes := esig.S.Bytes()
	whole := cryptoff.MergeTwoByteSlices(rBytes, sBytes)
	signature := base58.Encode(whole)
	log.Print("Signature created: ", signature)
	return signature
}

//createCoinbaseTx build a new coinbase transaction and assigns it to the given address
func CreateCoinbaseTx(address string) Transaction {
	log.Print("Creating coinbase transaction for ", address)
	var cbTx Transaction

	//no txin for coinbase tx

	var txOut TxOut
	txOut.Amount = COINBASE_AMOUNT
	txOut.Address = address
	cbTx.TxOuts = append(cbTx.TxOuts, txOut)

	cbTx.Id = calculateTxId(cbTx)
	cbTx.Signature = "coinbase"

	log.Print("Coinbase tx created")
	return cbTx
}

//SendCoins sends "count" number of coins to the "to" address, from the owner of given private key
func SendCoins(privKey *ecdsa.PrivateKey, to string, count int) Transaction {
	from := cryptoff.EncodePublicKey(&privKey.PublicKey)
	log.Print("Creating tx to send ", count, " coins from ", from, " to ", to)
	//TODO: error handling (insufficient funds)
	txIns, total := findTxInsFor(from, count)
	txOuts := SplitTxIns(from, to, count, total)
	tx := createTx(privKey, txIns, txOuts)
	log.Print("Send-tx created")
	return tx
}

//createTx builds a new transaction where the sender is identified by the given private key,
//and where the transaction includes the given tXins and txOuts
func createTx(privKey *ecdsa.PrivateKey, txIns []TxIn, txOuts []TxOut) Transaction {
	pubKey := cryptoff.EncodePublicKey(&privKey.PublicKey)
	log.Print("Creating tx from ", pubKey, " with ", len(txIns), " tx-ins, ", len(txOuts), " tx-outs")
	tx := Transaction{"", pubKey, "", txIns, txOuts}

	signTxIns(tx, privKey)

	tx.Id = calculateTxId(tx)
	tx.Signature = signData(privKey, []byte(tx.Id))
	log.Print("Created tx and signed with signature: ", tx.Signature)
	return tx
}

//processTransaction takes txIns from transaction and removes any matching unspent txOuts,
//and creates new txOuts matching the transaction
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

//consumeTxOut scans the list of unspent txouts to find one with matching transaction id and index, removes that when found
func consumeTxOut(txId string, txIdx int) {
	log.Print("Consuming tx-out with tx-id=", txId, ", tx-idx=", txIdx)
	for idx, val := range unspentTxOuts {
		if val.TxId == txId && val.TxIdx == txIdx {
			//remove the matching unspend txout from the list of unspents
			//unpacking with ... https://www.reddit.com/r/golang/comments/1y8ytg/what_mean_in_append_function_when_appending_byte/
			unspentTxOuts = append(unspentTxOuts[:idx], unspentTxOuts[idx+1:]...)
			return
		}
	}
	//TODO: error when not found
}

//createTxOut creates a txout from given parameters (pubKey is recipient) and adds it to list of unspent txouts
func createTxOut(txId string, amount int, pubKey string, txIdx int) {
	log.Print("Creating tx-out from: tx-in=", txId, ", amount=", amount, ", pubkey=", pubKey, ", tx-idx=", txIdx)
	newUTXO := UnspentTxOut{txId, txIdx, pubKey, amount}
	unspentTxOuts = append(unspentTxOuts, newUTXO)
}

//signTxIns verifies that the given transaction is valid, i.e. all txin exist as unspent txout for the spending user
//TODO: check why did i call this sign... when no signing appears to happen -> rename this
func signTxIns(tx Transaction, privKey *ecdsa.PrivateKey) bool {
	//key from string https://stackoverflow.com/questions/48392334/how-to-sign-a-message-with-an-ecdsa-string-privatekey
	myAddress := cryptoff.EncodePublicKey(&privKey.PublicKey)
	//first param from range is index, second is the value
	for _, val := range tx.TxIns {
		errorStatus := false
		internalIdx := findUnspentTransaction(myAddress, val.TxId)
		if internalIdx < 0 {
			//TODO: error logging
			log.Print("Error: trying to spend a transaction that does not exist (as unspent..)")
			errorStatus = true
		}
		if errorStatus {
			return false
		}
	}
	return true
}

//hexToPublicKey converts a hex-encoded string into a goland public-key
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

	pub.Curve = cryptoff.Curve

	return pub
}
