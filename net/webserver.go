package net

import (
	"fmt"
	"github.com/mukatee/go-naive/chain"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

//https://tutorialedge.net/golang/creating-simple-web-server-with-golang/
//https://astaxie.gitbooks.io/build-web-application-with-golang/en/03.2.html
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // send data to client side
}

func rpcBlocks(w http.ResponseWriter, r *http.Request) {
	response := chain.JsonChain(chain.GlobalChain)
	fmt.Fprintf(w, response) // send data to client side
}

func rpcMineBlock(w http.ResponseWriter, r *http.Request) {
	block := chain.CreateBlock(nil, "RPC test block", 0)
	response := chain.JsonBlock(block)
	fmt.Fprintf(w, response) // send data to client side
}

func rpcListPeers(w http.ResponseWriter, r *http.Request) {
	response := jsonPeers(peers)
	fmt.Fprintf(w, response) // send data to client side
}

func rpcAddPeer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // parse arguments, you have to call this by yourself
	for k, v := range r.Form {
		if k == "ip" {
			if len(v) > 1 {
				str := strings.Join(v, " ")
				fmt.Println("Too many values (was len " + string(len(v)) + " expected 1): " + str)
				return
			}
			addPeer(Peer{v[0]})
		}
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
}

func StartServer() {
	http.HandleFunc("/hello", sayhelloName)     // set router
	http.HandleFunc("/blocks", rpcBlocks)       // set router
	http.HandleFunc("/mineblock", rpcMineBlock) // set router
	http.HandleFunc("/peers", rpcListPeers)     // set router
	http.HandleFunc("/addPeer", rpcAddPeer)     // set router
	//https://stackoverflow.com/questions/49067160/what-is-the-difference-in-listening-on-0-0-0-080-and-80
	//https://grokbase.com/t/gg/golang-nuts/141ee4dqyg/go-nuts-how-to-know-when-listenandserve-is-ready-to-handle-connections
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		os.Exit(1)
	}
	go http.Serve(listener, nil)
}
