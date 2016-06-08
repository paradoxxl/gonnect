package client

import (
	"encoding/binary"
	"log"
)

func (tn *TunNetwork)addNewClient(client *Client) {
    if client == nil {
        return
    }
    if client.RemoteAddress.IP == nil {
        return
    }
    if client.RemoteAddress.Port == 0 {
        return
    }
    if client.VirtualAddress == nil {
        return
    }
    tn.NetworkMembersLock.Lock()
    tn.NetworkMembers[binary.BigEndian.Uint32(client.VirtualAddress)] = &Client{
        Nickname: client.Nickname,
        RemoteAddress: client.RemoteAddress,
        VirtualAddress: client.VirtualAddress,
    }
    tn.NetworkMembersLock.Unlock()
    log.Printf("Added client %+v", *client)
}