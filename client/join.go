package client

import (
	"net"
    "encoding/binary"
    "github.com/arcpop/tun"
)


func createUDPConnection(port int) (*net.UDPConn, error) {
    return net.ListenUDP("udp", &net.UDPAddr{Port: port})
}

func joinNetwork(remoteConn *net.UDPConn, network *Network) (*TunNetwork, error) {
    adapter, err := tun.New("")
    if err != nil {
        return nil, err
    }

    tn := &TunNetwork{
        RemoteConn: remoteConn,
        NetworkMembers: make(map[uint32]*Client),
        Stop: make(chan interface{}),
        TunAdapter: adapter,
    }
    
    for _, client := range network.Networkmembers {
        if client == nil {
            continue
        }
        if client.RemoteAddress.IP == nil {
            continue
        }
        if client.RemoteAddress.Port == 0 {
            continue
        }
        ip := binary.BigEndian.Uint32(client.VirtualAddress.To4())
        if ip == 0 {
            continue
        }
        tn.NetworkMembers[ip] = &Client{
            Nickname: client.Nickname,
            RemoteAddress: client.RemoteAddress,
            VirtualAddress: client.VirtualAddress.To4(),
        }
    }

    go tn.listenTUN()
    go tn.listenUDP()

    return tn, nil
}

func (tn *TunNetwork) Disconnect() {
    close(tn.Stop)

    
}
