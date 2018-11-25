package chain

//bitcoin stores blocks on disk using leveldb:
//https://bitcoin.stackexchange.com/questions/69628/what-data-structure-should-i-use-to-model-a-blockchain

//https://stackoverflow.com/questions/16900938/how-to-place-golang-project-a-set-of-packages-to-github

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var chainPath = "node/blocks/"
var chainFileName = "blocks.json"

var Curve = elliptic.P256()

//this is the current chain this node is on
var globalChain []Block

type Block struct {
	Index        int           //the block index in the chain
	Hash         string        //hash for this block
	PreviousHash string        //hash for previous block
	Timestamp    time.Time     //time when this block was created
	Data         string        //the data in this block. could be anything. not really needed since real data is transaction but for fun..
	Transactions []Transaction //the transactions in this block
	Difficulty   int           //block difficulty when created
	Nonce        int           //nonce used to find the hash for this block
}

//list of all transactions in the blockchain
var allTransactions []Transaction

//list of all unspent tx-outs in the blockchain
var unspentTxOuts []UnspentTxOut

func addTransaction(tx Transaction) {
	log.Println("Adding transaction:", tx)
	oldTx := findUnspentTransaction(tx.Sender, tx.Id)
	if oldTx >= 0 {
		log.Println("transaction already exists, not adding: ", tx.Id)
		return
	}
	allTransactions = append(allTransactions, tx)
	for _, txIn := range tx.TxIns {
		log.Println("Deleteing tx output: ", txIn.TxId)
		deleteUnspentTransaction(tx.Sender, txIn.TxId)
	}
	for idx, txOut := range tx.TxOuts {
		utx := UnspentTxOut{tx.Id, idx, txOut.Address, txOut.Amount}
		log.Println("Created unspent tx out:", utx)
		unspentTxOuts = append(unspentTxOuts, utx)
	}
}

//check that the blockchain has a transaction with the given id
//returns the index of matching (block, transaction) in the blockchain or -1, -1 if not found
func findTransaction(txId string) (int, int) {
	print("looking for txid:", txId)
	for bIdx, block := range globalChain {
		for tIdx, tx := range block.Transactions {
			if tx.Id == txId {
				return bIdx, tIdx
			}
		}
	}
	return -1, -1
}

//check that the blockchain has a given unspent transaction for the given public key (user)
//returns the index of matching transaction in the list of that users unspent transactions or -1 if not found
func findUnspentTransaction(pubKey string, txId string) int {
	log.Println("Looking to find unspent tx, pubkey=:", pubKey, ", txid=", txId)
	for idx, utx := range unspentTxOuts {
		if utx.TxId == txId && utx.Address == pubKey {
			log.Println("tx found at index:", idx)
			return idx
		}
	}
	log.Println("no matching tx found, returning -1")
	return -1
}

//remove a transaction from the list of unspent transactions
func deleteUnspentTransaction(pubKey string, txId string) bool {
	log.Println("deleting tx, pubkey=", pubKey, ", txid=", txId)
	idx := findUnspentTransaction(pubKey, txId)
	if idx < 0 {
		log.Println("no matching tx found")
		return false
	}
	log.Println("matching tx found, removing it")

	//val, _ := userTxMap[pubKey] //this would test for existtence of key
	//https://stackoverflow.com/questions/21326109/why-are-lists-used-infrequently-in-go
	//http://yourbasic.org/golang/three-dots-ellipsis/
	//it seems "..." stands for "unpacking the slice", which is needed since append 2nd argument is varargs, not slice
	unspentTxOuts = append(unspentTxOuts[:idx], unspentTxOuts[idx+1:]...)

	return true
}

//https://stackoverflow.com/questions/15323767/does-golang-have-if-x-in-construct-similar-to-python#15323988
func stringInSlice(a string, list []string) int {
	for idx, b := range list {
		if b == a {
			return idx
		}
	}
	return -1
}

//calculate hash string for the given block
func hash(block *Block) string {
	log.Println("hashing block:", block)
	indexStr := strconv.Itoa(block.Index)
	timeStr := strconv.FormatUint(uint64(block.Timestamp.Unix()), 16) //base 16 output
	nonceStr := strconv.Itoa(block.Nonce)
	diffStr := strconv.Itoa(block.Difficulty)
	txBytes, _ := json.Marshal(block.Transactions)
	txStr := string(txBytes)
	//this joins all the block elements to one long string with all elements appended after another, to produce the hash
	blockStr := strings.Join([]string{indexStr, block.PreviousHash, timeStr, diffStr, block.Data, txStr, nonceStr}, " ")
	//	print(data)
	bytes := []byte(blockStr)
	hash := sha256.Sum256(bytes)
	s := hex.EncodeToString(hash[:]) //encode the Hash as a hex-string. the [:] is slicing to match datatypes in args
	log.Println("Created hash:", s)
	return s
}

//create genesis block, the first one on the chain to bootstrap the chain
func createGenesisBlock(addToChain bool) Block {
	log.Println("Creating genesis block")
	genesisTime, _ := time.Parse("Jan 2 15:04 2006", "Mar 15 19:00 2018")
	genesisAddress := "3uGQcE7wPMYoJDGStgWMkm6qj7u83TDmcvXUYVoVeDq4sYRMZXisB6vxMwQWCMPe1eX5rGPgoJ9oyYoFNGwpqNcPU"
	cbTx := createCoinbaseTx(genesisAddress)
	txs := []Transaction{cbTx}
	block := Block{1, "", "0", genesisTime, "Teemu oli täällä", txs, 1, 1}
	hash := hash(&block)
	block.Hash = hash
	if addToChain {
		log.Println("Adding genesis block to chain")
		globalChain = append(globalChain, block)
	}
	log.Println("Genesis block creation finished:", block)
	return block
}

//check if the given block matches the genesis block.
//since the genesis block has no previous block to compare, need to do separate check
func checkGenesisBlock(block Block) bool {
	genesis := createGenesisBlock(false)
	if block.Hash != genesis.Hash {
		return false
	}
	if block.Index != genesis.Index {
		return false
	}
	if block.PreviousHash != genesis.PreviousHash {
		return false
	}
	if block.Timestamp != genesis.Timestamp {
		return false
	}
	if block.Data != genesis.Data {
		return false
	}

	return true
}

//validate the overall chain, starting from genesis block all the way through the whole chain until the last block
func validateChain(chain []Block) bool {
	checkGenesisBlock(chain[0])
	for i := 1; i < len(chain); i++ {
		//validate index is in sequence and is +1 from previous block
		thisIndex := chain[i].Index
		prevIndex := chain[i-1].Index + 1
		if thisIndex != prevIndex {
			log.Println("Index issue at " + strconv.Itoa(i) + " - " + strconv.Itoa(thisIndex) + " vs " + strconv.Itoa(prevIndex))
			return false
		}
		//validate that previous hash stored in this block matches the hash stored for previous block in chain
		prevHash1 := chain[i].PreviousHash
		prevHash2 := chain[i-1].Hash
		if prevHash1 != prevHash2 {
			log.Println("Hash mismatch at " + strconv.Itoa(i) + " - " + prevHash1 + " vs " + prevHash2)
			return false
		}
		//validate the hash stored in this block is a valid hash for this block
		hash := hash(&chain[i])
		if hash != chain[i].Hash {
			log.Println("Hash mismatch with itself at index", i)
			return false
		}
	}
	log.Println("chain validated")
	return true
}

//create a block from the given parameters, and find a nonce to produce a hash matching the difficulty
//finally, append new block to current chain
func createBlock(txs []Transaction, blockData string, difficulty int) Block {
	log.Println("Creating new block, tx count = ", len(txs), "difficult=", difficulty, "block-data=", blockData)
	chainLength := len(globalChain)
	log.Println("current chain len:", chainLength)
	previous := globalChain[chainLength-1]
	index := previous.Index + 1
	timestamp := time.Now().UTC()
	nonce := 0
	newBlock := Block{index, "", previous.Hash, timestamp, blockData, txs, difficulty, nonce}
	log.Println("starting pow for block template:", newBlock)
	for {
		hash := hash(&newBlock)
		newBlock.Hash = hash
		//TODO: exit/update if peer finds hash
		if verifyHashVsDifficulty(hash, difficulty) {
			addBlock(newBlock)
			log.Println("found pow hash:", hash)
			//			globalChain = append(globalChain, newBlock)
			return newBlock
		}
		nonce++
		newBlock.Nonce = nonce
	}
}

//add a new block to the existing chain
func addBlock(block Block) {
	chainLength := len(globalChain)
	log.Println("adding block to chain. current height=", chainLength, ", block=", block)
	previousBlock := globalChain[chainLength-1]
	block.PreviousHash = previousBlock.Hash
	globalChain = append(globalChain, block)
	log.Println("Adding " + strconv.Itoa(len(block.Transactions)) + " transactions from block.")
	for _, tx := range block.Transactions {
		addTransaction(tx)
	}
	//todo: check block hash matches difficulty
}

func printBlock(block Block) {
	fmt.Printf("block %d:%s %s %d %s\n", block.Index, block.Hash, block.Timestamp.String(), block.Difficulty, block.Data)
	//txStrs := make(map[string]int)
	for txIdx, tx := range block.Transactions {
		fmt.Printf("-tx: %d\n", txIdx)
		if len(tx.TxIns) == 0 {
			fmt.Print("--No txIn found, assuming coinbase tx.\n")
		} else {
			for _, txIn := range tx.TxIns {
				blockIdx, txIdx := findTransaction(txIn.TxId)
				fmt.Printf("--txin: block id = %d, txid = %d\n", blockIdx, txIdx)
				txOut := globalChain[blockIdx].Transactions[txIdx]
				fmt.Printf("---in from %s: (%s, %d)\n", txOut.Sender, txIn.TxId)
			}
		}
		for _, txOut := range tx.TxOuts {
			fmt.Printf("--out to %s: %d\n", txOut.Address, txOut.Amount)
		}
	}
	//	println()
}

func printChain(chain []Block) {
	for _, block := range chain {
		printBlock(block)
	}
}

//turn the given chain into a json description, to copy to other location
func JsonChain(chain []Block) string {
	log.Println("Converting chain into JSON")
	bytes, _ := json.Marshal(chain)
	json := string(bytes)
	return json
}

//turn the given block into a json description, to copy to other location
func JsonBlock(block Block) string {
	log.Println("Converting block into JSON")
	bytes, _ := json.Marshal(block)
	json := string(bytes)
	return json
}

//print the given chain in json format
func printChainJSON(chain []Block) {
	bytes, _ := json.Marshal(chain)
	println(string(bytes))
}

//turn the given chain into json and back into a go-struct. for testing
func marshallDemarshallChain(chain []Block) {
	bytes, _ := json.Marshal(chain)
	str := string(bytes)
	bytes2 := []byte(str)
	chain2 := []Block{}
	json.Unmarshal(bytes2, &chain2)
	printChainJSON(chain2)
}

//validate given chain, compare to current, and if new is longer replace current chain with new
//returns true if new chain is selected
//matches the first version of blockchain from naivecoin tutorial:
//https://lhartikk.github.io/jekyll/update/2017/07/14/chapter1.html
func takeLongestChain(newChain []Block) bool {
	newLength := len(newChain)
	oldLength := len(globalChain)
	log.Println("Comparing new chain vs old chain length:", newLength, " vs. ", oldLength)
	fine := validateChain(newChain)
	if !fine {
		log.Println("New chain failed to validate, dropping it.")
		return false
	}
	if newLength > oldLength {
		log.Println("New chain longer, replacing old.")
		globalChain = newChain
	} else {
		log.Println("New chain not longer, keeping old.")
	}
	return true
}

//replaces current global chain with new if new is valid and has higher diff
//return true if new chain has higher difficulty than current global chain
//matches second version of blockchain from naivecoin tutorial:
//https://lhartikk.github.io/jekyll/update/2017/07/13/chapter2.html
func takeMostDifficultChain(newChain []Block) bool {
	fine := validateChain(newChain)
	if !fine {
		log.Println("New chain validation failed, dropping it.")
		return false
	}

	totalDiff1 := calculateChainDifficulty(globalChain)
	totalDiff2 := calculateChainDifficulty(newChain)

	if totalDiff2 < totalDiff1 {
		log.Println("not switching chain")
		return false
	}
	log.Println("switching chain to more difficult")
	globalChain = newChain
	return true
}

func calculateChainDifficulty(chain []Block) float64 {
	totalDiff := 0.0
	for i := 0; i < len(chain); i++ {
		totalDiff += math.Pow((float64(chain[i].Difficulty)), 2)
	}
	log.Println("Calculated chain difficulty:", totalDiff)
	return totalDiff
}

//create a chain with a single coinbase transaction to given address
func BootstrapTestEnv(address string) {
	log.Println("Bootstrapping test env.")
	createGenesisBlock(true)

	cbTx := createCoinbaseTx(address)
	txs := []Transaction{cbTx}
	createBlock(txs, "My data", 0)
	log.Println("Test env bootstrapped")
}

func writeBlockChain() {
	log.Println("Starting to write blockchain to disk")
	chainJson := JsonChain(globalChain)
	err := os.MkdirAll(chainPath, os.ModePerm)
	log.Println("File path created/found, " + chainPath)
	f, err := os.Create(chainPath + chainFileName)
	log.Println("Done OK, now to write.., " + chainPath)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer f.Close()

	f.WriteString(chainJson)
}

func InitBlockChain() bool {
	err := os.MkdirAll(chainPath, os.ModePerm)
	_, err = os.Stat(chainPath + chainFileName)
	loaded := false
	if os.IsNotExist(err) {
		log.Println("No blockchain file found.")
		//TODO: start downloading chain
	} else {
		// file/dir with chain path already exists
		readBlockChain()
		loaded = true
	}
	return loaded
}

//readBlockChain() reads the chain from disk.
func readBlockChain() {
	fullPath := chainPath + chainFileName
	log.Println("Reading blockchain from disk, path = " + fullPath)
	bytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Println(err)
		//TODO: panic?
	}
	//convert bytes in file to golang objecs in the chain
	json.Unmarshal(bytes, &globalChain)
	//a slice is passed as copy of header, so this does not copy the whole array
	//https://stackoverflow.com/questions/39993688/are-golang-slices-pass-by-value#39993797
	valid := validateChain(globalChain)
	if !valid {
		panic("Invalid chain loaded, exit")
	}
}
