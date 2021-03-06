package main

import (
	"github.com/mukatee/go-naive/chain"
	"github.com/mukatee/go-naive/wallet"
	"io"
	"log"
	"os"
)

func main() {
	print(wallet.HelpText)
	setupLogging()
	addr, loaded := wallet.InitWallet()
	loaded = chain.InitBlockChain()
	if !loaded {
		chain.CreateTestChain(addr, 2)
	}
	wallet.ReadConsole()
}

func setupLogging() {
	logFile, err := os.OpenFile("tc-log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	//https://stackoverflow.com/questions/36719525/how-to-log-messages-to-the-console-and-a-file-both-in-golang
	mw := io.MultiWriter(os.Stdout, logFile)
	//log.SetPrefix("LOG: ")
	//log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.SetOutput(mw)
	log.Println("Starting system, logging setup done.")
	//	log.SetOutput(logFile)
}
