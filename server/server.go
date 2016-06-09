package server

import (
	"github.com/paradoxxl/gonnect/msg"
	"math/rand"
	"net"
	"strings"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	var client PeerState

	for {
		//Evaluate the command
		msg.ReadMessage(conn, client)

	}
}

func (peer *PeerState) OnCommand(commandType byte, message interface{}) {
	switch commandType {
	case msg.CreateNwType:
		data := msg.CreateNetworkCommand{}(message)
		peer.createNetwork(data.Networkname,data.Networkpass,data.Networkip,data.Peername,data.Peerport)
		peer.joinNetwork(data.Networkname,data.Networkpass,data.Peername,data.Peerport)
	case msg.JoinNwType:
		data := msg.JoinNetworkCommand(message)
		peer.joinNetwork(data.Networkname,data.Networkpass,data.Peername,data.Peerport)
	case msg.DisconnectNw:
		peer.disconnectNetwork()
	}
}
func (peer *PeerState) createNetwork(Networkname string,Networkpass string,Networkip net.IP,Peername string,Peerport string) {
	_, networkexists := State.m[Networkname]
	if networkexists {
		return
	}

	State.RLock()
	defer State.RUnlock()

	State.m[Networkname] = Network{
		Networkname:    Networkname,
		Networkmembers: make(map[string]msg.Peer),
		Networkpass:    Networkpass,
		Networkip:      Networkip,
	}

	peer.network = msg.Network{Networkname: Networkname}
}

func (peer *PeerState) joinNetwork(	Networkname string,Networkpass string,Peername string,Peerport string) {

	State.RLock()
	defer State.RUnlock()

	_, networkexists := State.m[Networkname]
	if !networkexists {
		return
	}

	//get Client IP Address
	peer.peer.Peeraddress = strings.Split(peer.conn.RemoteAddr().String()[0], ":") + ":" + Peerport

	//check free IP address. Select one at random, check for collisions. increment on collision
	ok := false
	var ip = rand.Intn(253) + 1
	for !ok {
		ok = true
		for _, v := range State.m[Networkname].Networkmembers {
			cliaddr := v.peer.VirtualAddress[3]
			if cliaddr == ip || ip == 0 {
				ok = false
				ip = (ip+1)%255
				break
			}
		}
	}


	peer.peer.VirtualAddress = State.m[Networkname].Networkip.To4()
	peer.peer.VirtualAddress[3] = ip

	append(State.m[Networkname].Networkmembers, peer.peer)
	peer.network.Networkname = Networkname
	append(peer.network.Networkmembers, peer.peer)

	data := msg.EncodeNetworkJoinNotification(msg.NetworkJoinNotification{
		VirtualIPAddress: peer.peer.VirtualAddress,
		Peers:            peer.network.Networkmembers,
	})

	//Notify all peers including the new

	for _, v := range State.m[Networkname].Networkmembers {
		n, err := v.conn.Write(data)
		if err {
			//TODO: Do something
		}
		if n != len(data) {
			//TODO: Try resend
		}
	}

}

func (peer *PeerState) disconnectNetwork() {
	defer 	peer.conn.Close()
	
	State.RLock()
	nw,exists := State.m[peer.network.Networkname]
	State.RUnlock()

	if !exists {
		//TODO: Handle this properly. Should really not happen, but who knows?
		return
	}

	nw.Lock()
	_,exists = nw.Networkmembers[peer.peer.Peername]

	if !exists {
		//TODO: Handle this properly. Should really not happen, but who knows?
		return
	}
	delete(nw.Networkmembers, peer.peer.Peername)

	nw.Unlock()


	//check whether the there are no more peers in the Network. If so, delete it
	if len(nw.Networkmembers) == 0 {
		State.Lock()
		delete(State, peer.network.Networkname)
		State.Unlock()
	}else{
		nw.RLock()
		data := msg.EncodeClientDisconnectNotification(msg.ClientDisconnectNotification{
			VirtualIPAddress: peer.peer.VirtualAddress,
		})

		for _, v := range nw.Networkmembers {
			n, err := v.conn.Write(data)
			if err {
				//TODO: Do something
			}
			if n != len(data) {
				//TODO: Try resend
			}
		}
		nw.Unlock()
	}



}
