package server

import (
	"github.com/paradoxxl/gonnect/msg"
	"net"
	"sync"
)

type Client struct {
	cli  msg.Peer
	nw   msg.Network
	conn net.Conn
}

type Network struct {
	Networkname    string
	Networkmembers map[string]Client
	Networkpass    string
	Networkip      net.IP
	sync.RWMutex
}

var State = struct {
	sync.RWMutex
	m map[string]Network
}{m: make(map[string]Network)}

