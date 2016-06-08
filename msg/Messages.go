package msg

import (
	"net"
	"errors"
	"encoding/json"
	"encoding/binary"
	"io"
)

const (
	Version byte = 10
)

const (
	CreateNwType byte = iota
	JoinNwType
	DisconnectNw
	NotifyNetworkJoined
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

type DisconnectNetworkCommand struct {
	VirtualIPAddress net.IP
}

type ClientDisconectNotification struct {
	VirtualIPAddress net.IP
}
type ClientJoinNotification struct {
	VirtualIPAddress net.IP
	RemoteAddress net.UDPAddr
	Peername string
}
type ClientDisconnectNotification struct {
	VirtualIPAddress net.IP
}
type NetworkJoinNotification struct {
	VirtualIPAddress net.IP
	Peers []*Peer
}


type MessageHandler interface {
	OnCommand(commandType byte, message interface{})
}

func ReadMessage(rd io.Reader, handler MessageHandler) (error) {
	var hdr [6]byte
	_, err := io.ReadFull(rd, hdr[:])
	if err != nil {
		return err
	}

	if hdr[0] != Version {
		return errors.New("Version mismatch")
	}

	pktlen := binary.BigEndian.Uint32(hdr[1:5])
	if pktlen < 6 {
		return errors.New("Invalid length")
	}
	pktlen -= 6
	if pktlen == 0 {
		handler.OnCommand(hdr[5], nil)
		return nil
	}

	buf := make([]byte, pktlen)

	_, err = io.ReadFull(rd, buf)
	if err != nil {
		return err
	}

	var msg interface{}
	cmdType := hdr[5]
	switch cmdType {
	case CreateNwType:
		msg, err = DecodeCreateNetworkCommand(buf)
	case JoinNwType:
		msg, err = DecodeCreateNetworkCommand(buf)
	default:
		err = errors.New("Unsupported command type")
	}

	if err != nil {
		return err
	}

	handler.OnCommand(cmdType, msg)
	return nil
}
func EncodeNetworkJoinNotification(data *NetworkJoinNotification) []byte {
	header := []byte{Version, 0, 0, 0, 0, NotifyNetworkJoined}
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return p
}
func EncodeClientJoinNotification(data *ClientJoinNotification) []byte {
	header := []byte{Version, 0, 0, 0, 0, NotifyClientJoin}
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return p
}
func EncodeClientDisconnectNotification(data *ClientDisconnectNotification) []byte {
	header := []byte{Version, 0, 0, 0, 0, NotifyClientDisconect}
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return p
}
func EncodeDisconnectNetworkCommand(data *DisconnectNetworkCommand) []byte {
	header := []byte{Version, 0, 0, 0, 0, DisconnectNw}
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return p
}
func EncodeCreateNetworkCommand(data *CreateNetworkCommand) []byte {
	header := []byte{Version, 0, 0, 0, 0, CreateNwType}
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return p
}
func EncodeJoinNetworkCommand(data *JoinNetworkCommand) []byte{
	header := []byte{Version, 0, 0, 0, 0, JoinNwType}
	p, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return p
}

func DecodeDisconnectNetworkCommand(msg []byte) (*DisconnectNetworkCommand, error) {
	var cmd DisconnectNetworkCommand
	err := json.Unmarshal(msg, &cmd)
	return &cmd, err
}
func DecodeClientDisconnectNotification(msg []byte) (*ClientDisconnectNotification, error) {
	var cmd ClientDisconnectNotification
	err := json.Unmarshal(msg, &cmd)
	return &cmd, err
}
func DecodeClientJoinNotification(msg []byte) (*ClientJoinNotification, error) {
	var cmd ClientJoinNotification
	err := json.Unmarshal(msg, &cmd)
	return &cmd, err
}
func DecodeJoinNetworkCommand(msg []byte) (*JoinNetworkCommand, error) {
	var cmd JoinNetworkCommand
	err := json.Unmarshal(msg, &cmd)
	return &cmd, err
}
func DecodeCreateNetworkCommand(msg []byte) (*CreateNetworkCommand, error) {
	var cmd CreateNetworkCommand
	err := json.Unmarshal(msg, &cmd)
	return &cmd, err
}
func DecodeNetworkJoinNotification(msg []byte) (*NetworkJoinNotification, error) {
	var cmd NetworkJoinNotification
	err := json.Unmarshal(msg, &cmd)
	return &cmd, err
}