package net

import "encoding/json"

type Peer struct {
	Address string
}

var peers = []Peer{}

func addPeer(peer Peer) {
	peers = append(peers, peer)
}

func jsonPeers(peers []Peer) string {
	bytes, _ := json.Marshal(peers)
	json := string(bytes)
	return json
}
