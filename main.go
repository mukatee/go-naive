package main

import "github.com/mukatee/go-naive/chain"

func main() {
	addr, loaded := chain.InitWallet()
	loaded = chain.InitBlockChain()
	if !loaded {
		chain.BootstrapTestEnv(addr)
	}
	chain.ReadConsole()
}

