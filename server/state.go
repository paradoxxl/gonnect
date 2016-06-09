package server

import (
	"github.com/paradoxxl/gonnect/msg"
	"net"
	"sync"
)

type PeerState struct {
	peer    msg.Peer
	network msg.Network
	conn    net.Conn
}

type Network struct {
	Networkname    string
	Networkmembers map[string]PeerState
	Networkpass    string
	Networkip      net.IP
	sync.RWMutex
}

var State = struct {
	sync.RWMutex
	m map[string]Network
}{m: make(map[string]Network)}

