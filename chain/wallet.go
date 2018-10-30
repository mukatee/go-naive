package chain

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var walletPath = "node/wallet/"
var walletFileName = "wallet.json"
var walletKey *ecdsa.PrivateKey
var publicAddr string
var walletBalance int64

func ReadConsole() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome, sir!")
	fmt.Print("> ")
readloop:
	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "balance":
			balance := balanceFor(publicAddr)
			fmt.Println(balance)
		case "exit":
			writeWallet()
			writeBlockChain()
			break readloop
		case "save":
			writeWallet()
			writeBlockChain()
		case "send":
			walletSend()
		}
		//		fmt.Println(input)
		fmt.Print("> ")
	}

	if scanner.Err() != nil {
		// handle error.
		fmt.Println("error", scanner.Err())
	}
}

func InitWallet() (string, bool) {
	err := os.MkdirAll(walletPath, os.ModePerm)
	_, err = os.Stat(walletPath + walletFileName)
	loaded := false
	if os.IsNotExist(err) {
		log.Print("No wallet file found. Creating new.")
		walletKey, _, publicAddr = createAddress()
		writeWallet()
	} else {
		// file/dir with wallet path already exists
		log.Print("Existing wallet file found, reading wallet from disk.")
		readWallet()
		loaded = true
	}
	return publicAddr, loaded
}

func readWallet() {
	//https://gobyexample.com/json
	var data map[string]interface{}
	fullPath := walletPath + walletFileName
	log.Print("Reading wallet from path:" + fullPath)
	bytes, err := ioutil.ReadFile(fullPath)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		panic(err)
	}
	//TODO: is this intended to be int or float?
	walletBalance = int64(data["balance"].(float64))
	//TODO: why storing empty private key leads to this weirdness?
	//println("priv:", data["priv"].(string))
	walletKey = decodePrivateKey(data["priv"].(string))
	//println("wallet key:", walletKey.D.String())
	//println("", walletKey.PublicKey)
	pubStr := encodePublicKey(walletKey.PublicKey)
	log.Println("wallet public key: ", pubStr)
}

func writeWallet() {
	privStr := encodePrivateKey(walletKey)
	fullPath := walletPath + walletFileName
	log.Print("Writing wallet to path:" + fullPath)
	err := os.MkdirAll(walletPath, os.ModePerm)
	f, err := os.Create(walletPath + walletFileName)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer f.Close()

	//TODO: wallet struct
	//https://gobyexample.com/json
	//https://stackoverflow.com/questions/18526046/mapping-strings-to-multiple-types-for-json-objects
	content := map[string]interface{}{"priv": privStr, "pubaddr": publicAddr, "balance": walletBalance}
	contentB, _ := json.Marshal(content)
	f.WriteString(string(contentB))
}

func walletSend() {
	//TODO: add sending funds to another user
}

//balanceFor counts the unspent balance for given address (as count of unspent txouts)
//address parameter given is the base58 encoded public key
func balanceFor(address string) int {
	log.Print("Calculating balance for address:" + address)
	balance := 0
	for _, val := range unspentTxOuts {
		if val.Address == address {
			balance += val.Amount
		}
	}
	log.Print("Balance for " + address + " = " + strconv.Itoa(balance))
	return balance
}

//findTxInsFor looks for unspent txouts for the given address to match the amount wanting to spend
func findTxInsFor(address string, amount int) ([]TxIn, int) {
	log.Print("Searching for unspent txOuts for " + address + ", to amount of " + strconv.Itoa(amount))
	balance := 0
	var unspents []TxIn
	for _, val := range unspentTxOuts {
		if val.Address == address {
			balance += val.Amount
			txIn := TxIn{val.TxId, val.TxIdx}
			unspents = append(unspents, txIn)
		}
		if balance >= amount {
			log.Print("Found unspent txOuts: ", unspents, ", total funds = "+strconv.Itoa(balance))
			return unspents, balance
		}
	}
	log.Print("Did not find suffient funds for " + address + ", requested " + strconv.Itoa(amount) + ", found " + strconv.Itoa(balance))
	return nil, -1
}

//splitTxIns produces two txouts, by taking the total sum of txins and the amount to send
//and splitting this to one txout for the coins to send, and another for the remains to send back to self
func splitTxIns(from string, to string, toSend int, total int) []TxOut {
	log.Print("Creating txIn splits for transaction from " + from + " to " + to)
	diff := total - toSend
	txOut := TxOut{to, toSend}
	var txOuts []TxOut
	txOuts = append(txOuts, txOut)
	if diff == 0 {
		log.Print("send and txin amount equal, only creating single txout:", txOuts)
		return txOuts
	}
	log.Print("created sending txout and change txout:", txOuts)
	txOut2 := TxOut{from, diff}
	txOuts = append(txOuts, txOut2)
	return txOuts
}
