package client

import (
    "sync"
    "net"
    "github.com/arcpop/tun"
	"encoding/binary"
	"log"
)


type Network struct {
    Networkname string
    Networkmembers []*Client
}
type Client struct {
    Nickname string
    RemoteAddress net.UDPAddr
    VirtualAddress net.IP
}

type TunNetwork struct {
    NetworkMembers map[uint32] *Client
    NetworkMembersLock sync.RWMutex
    RemoteConn *net.UDPConn
    Stop chan interface{}
    TunAdapter tun.TunInterface
}

func (tn* TunNetwork)listenUDP() {
    for {
        select {
            case _, ok := <- tn.Stop:
                if ok {
                    log.Println("Stopping ListenTUN!")
                }
                return
            default:
                tn.udpRxSinglePacket()
        }
    }
}

func (tn* TunNetwork)udpRxSinglePacket() {
    var buffer[2048]byte
    n, _, _, _, err := tn.RemoteConn.ReadMsgUDP(buffer[:], nil)
    if err != nil {
        return
    }
    pkt := buffer[:n]
    _, err = tn.TunAdapter.Write(pkt)
    if err != nil {
        log.Println("TunAdapter.Write() -> ", err)
    }
}


func (tn* TunNetwork)listenTUN() {
    for {
        select {
            case _, ok := <- tn.Stop:
                if ok {
                    log.Println("Stopping ListenTUN!")
                }
                return
            default:
                tn.tunParseSinglePacket()
        }
    }
}

func (tn* TunNetwork)tunParseSinglePacket() {
    var buffer [2048]byte
    n, err := tn.TunAdapter.Read(buffer[:])
    if err != nil {
        return
    }
    pkt := buffer[:n]
    ip := parseDestinationIPAddress(pkt)
    if ip == nil {
        return
    }
    if ip[3] == 255 {
        tn.NetworkMembersLock.RLock()
        for _, client := range tn.NetworkMembers {
            _, _, err = tn.RemoteConn.WriteMsgUDP(pkt, nil, &client.RemoteAddress)
            if err != nil {
                log.Println("Broadcast: RemoteConn.WriteMsgUDP() to ", client.RemoteAddress, " -> ", err)
            }
        }
        tn.NetworkMembersLock.RUnlock()
    } else {
        tn.NetworkMembersLock.RLock()
        client, ok := tn.NetworkMembers[binary.BigEndian.Uint32(ip)]
        tn.NetworkMembersLock.RUnlock()
        if !ok {
            log.Println("No client with IP ", ip)
            return
        }
        _, _, err = tn.RemoteConn.WriteMsgUDP(pkt, nil, &client.RemoteAddress)
        if err != nil {
            log.Println("Unicast: RemoteConn.WriteMsgUDP() to ", client.RemoteAddress, " -> ", err)
        }
    }
}

func parseDestinationIPAddress(pkt []byte) net.IP {
    if len(pkt) < 20 {
        return nil
    }
    version := pkt[0] & 0xF
    if version != 4 {
        return nil
    }
    return pkt[16:20]
}