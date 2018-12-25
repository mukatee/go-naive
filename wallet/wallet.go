package wallet

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/mukatee/go-naive/chain"
	"github.com/mukatee/go-naive/cryptoff"
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
	fmt.Print("wallet> ")

	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "balance":
			balance := chain.BalanceFor(publicAddr)
			fmt.Println(balance)
		case "exit":
			writeWallet()
			chain.WriteBlockChain()
			os.Exit(1)
			//break readloop
		case "save":
			writeWallet()
			chain.WriteBlockChain()
		case "send":
			walletSend()
		case "show address":
			log.Print("Wallet address: ")
			log.Print("        pubkey: ", cryptoff.EncodePublicKey(&walletKey.PublicKey))
			log.Print("       privkey: ", cryptoff.EncodePrivateKey(walletKey))
		case "create address":
			privKey, pubKey, addressStr := cryptoff.CreateAddress()
			privStr := cryptoff.EncodePrivateKey(privKey)
			pubStr := cryptoff.EncodePublicKey(pubKey)
			log.Print("Created address: ", addressStr)
			log.Print("Created pubkey: ", pubStr)
			log.Print("Created priveky: ", privStr)
		case "address":
			println(publicAddr)
		case "private key":
			print(cryptoff.EncodePrivateKey(walletKey))
		case "blocks":
			chain.PrintChain(chain.GlobalChain)
		case "mine block":
			tx := chain.CreateCoinbaseTx(publicAddr)
			txs := []chain.Transaction{tx}
			chain.CreateBlock(txs, "Hello", 0)
		default:
			println("Unknown command: ", input)
		}

		//		fmt.Println(input)
		fmt.Print("wallet> ")
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
		walletKey, _, publicAddr = cryptoff.CreateAddress()
		log.Println("Created address: ", publicAddr)
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
	walletKey = cryptoff.DecodePrivateKey(data["priv"].(string))
	//println("wallet key:", walletKey.D.String())
	//println("", walletKey.PublicKey)
	publicAddr = cryptoff.EncodePublicKey(&walletKey.PublicKey)
	log.Println("wallet public key: ", publicAddr)
}

func writeWallet() {
	privStr := cryptoff.EncodePrivateKey(walletKey)
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
	scanner := bufio.NewScanner(os.Stdin)
	print("Receiver address:")
	scanner.Scan()
	receiver := scanner.Text()
	print("Amount to send:")
	scanner.Scan()
	amountStr := scanner.Text()
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		println("oh no, error occurred, no coins sent:", err)
		return
	}
	println("sending ", amount, "coins to", receiver)
	chain.SendCoins(walletKey, receiver, amount)
}
