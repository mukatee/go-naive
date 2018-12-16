package chain

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
	response := JsonChain(GlobalChain)
	fmt.Fprintf(w, response) // send data to client side
}

func rpcMineBlock(w http.ResponseWriter, r *http.Request) {
	block := CreateBlock(nil, "RPC test block", 0)
	response := JsonBlock(block)
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

func RunServer() {
	http.HandleFunc("/hello", sayhelloName)     // set router
	http.HandleFunc("/blocks", rpcBlocks)       // set router
	http.HandleFunc("/mineblock", rpcMineBlock) // set router
	http.HandleFunc("/peers", rpcListPeers)     // set router
	http.HandleFunc("/addPeer", rpcAddPeer)     // set router
	err := http.ListenAndServe(":9090", nil)    // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
