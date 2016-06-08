package server

import (
	"net"
	"github.com/paradoxxl/gonnect/msg"
)

type Network struct{
	Networkname string
	Networkmembers map[string]msg.Peer
	Networkpass string
	Networkip net.IP
}

var State = make(map[string]Network)
