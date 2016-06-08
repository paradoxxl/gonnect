package msg

import "net"


const (
	Version = iota
	CreateNwType
	JoinNwType
	DisconnectNw
	NotifyClientDisconect
	NotifyClientJoin
)


type Peer struct {
	Peername string
	Peeraddress net.UDPAddr
	VirtualAddress net.IP
}

type Network struct {
	Networkname string
	Networkmembers []*Peer
}

type CreateNetworkCommand struct{
	Networkname string
	Networkpass string
	Networkip net.IP
	Peername string
	Peerport string
}

type JoinNetworkCommand struct{
	Networkname string
	Networkpass string
	Peername string
	Peerport string
}


func CheckMsg (msg []*byte) (version,cmdType int, err error){
	version = msg[0]

	length := uint32(msg[2:5])
	if length < 6 {
		err = "Bad Packet length"
		return
	}
	cmdType = msg[5]
	return
}

func EncodeCreateNetworkCommand(data CreateNetworkCommand) []*byte{
	return
}
func EncodeJoinNetworkCommand(data JoinNetworkCommand) []*byte{
	return
}

func DecodeCreateNetworkCommand(msg []*byte) CreateNetworkCommand{
	return
}
func DecodeJoinNetworkCommand(msg []*byte) JoinNetworkCommand{
	return
}