package client

import (
	"net"
    "encoding/binary"
    "github.com/arcpop/tun"
)
/*func JoinNetwork(networkName, networkPassword, nickName string, server string) (*TunNetwork, error) {
    var remoteConn *net.UDPConn
    var err error
    for {
        port := 3000 + rand.Intn(12000)
        remoteConn, err = net.ListenUDP("udp", &net.UDPAddr{ Port: port })
        if err == nil {
            break
        }
    }

    resp, err := http.Get(server + "/join?" + 
        "networkname=" + base64.URLEncoding.EncodeToString([]byte(networkName)) + 
        "&networkpass=" + base64.URLEncoding.EncodeToString([]byte(networkPassword)) + 
        "&nickname=" + base64.URLEncoding.EncodeToString([]byte(nickName)))
    if err != nil  {
        remoteConn.Close()
        return nil, err
    }
    if resp.StatusCode != 200 {
        remoteConn.Close()
        
    }
    defer resp.Body.Close()

    var network Network
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(&network)
    if err != nil {
        remoteConn.Close()
        return nil, err
    }
    nw, err := joinNetworkInternal(remoteConn, &network)
    if err != nil {
        remoteConn.Close()
        return nil, err
    }
    return nw, nil
}*/

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
