package server

import (
	"github.com/paradoxxl/gonnect/msg"
	"math/rand"
	"net"
	"strings"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	var client Client

	for {
		//Evaluate the command
		msg.ReadMessage(conn, client)

	}
}

func (cli *Client) OnCommand(commandType byte, message interface{}) {
	switch commandType {
	case msg.CreateNwType:
		data := msg.CreateNetworkCommand{}(message)
		cli.createNetwork(data.Networkname,data.Networkpass,data.Networkip,data.Peername,data.Peerport)
		cli.joinNetwork(data.Networkname,data.Networkpass,data.Peername,data.Peerport)
	case msg.JoinNwType:
		data := msg.JoinNetworkCommand(message)
		cli.joinNetwork(data.Networkname,data.Networkpass,data.Peername,data.Peerport)
	case msg.DisconnectNw:
		cli.disconnectNetwork()
	}
}
func (cli *Client) createNetwork(Networkname string,Networkpass string,Networkip net.IP,Peername string,Peerport string) {
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

	cli.nw = msg.Network{Networkname: Networkname}
}

func (cli *Client) joinNetwork(	Networkname string,Networkpass string,Peername string,Peerport string) {

	State.RLock()
	defer State.RUnlock()

	_, networkexists := State.m[Networkname]
	if !networkexists {
		return
	}

	//get Client IP Address
	cli.cli.Peeraddress = strings.Split(cli.conn.RemoteAddr().String()[0], ":") + ":" + Peerport

	//check free IP address. Select one at random, check for collisions. increment on collision
	ok := false
	var ip = rand.Intn(253) + 1
	for !ok {
		ok = true
		for _, v := range State.m[Networkname].Networkmembers {
			cliaddr := v.cli.VirtualAddress[3]
			if cliaddr == ip || ip == 0 {
				ok = false
				ip = (ip+1)%255
				break
			}
		}
	}


	cli.cli.VirtualAddress = State.m[Networkname].Networkip.To4()
	cli.cli.VirtualAddress[3] = ip

	append(State.m[Networkname].Networkmembers, cli.cli)
	cli.nw.Networkname = Networkname
	append(cli.nw.Networkmembers, cli.cli)

	data := msg.EncodeNetworkJoinNotification(msg.NetworkJoinNotification{
		VirtualIPAddress: cli.cli.VirtualAddress,
		Peers:            cli.nw.Networkmembers,
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

func (cli *Client) disconnectNetwork() {
	State.RLock()
	nw,exists := State.m[cli.nw.Networkname]
	State.RUnlock()

	if !exists {
		//TODO: Handle this properly. Should really not happen, but who knows?
		return
	}

	nw.Lock()
	_,exists = nw.Networkmembers[cli.cli.Peername]

	if !exists {
		//TODO: Handle this properly. Should really not happen, but who knows?
		return
	}
	delete(nw.Networkmembers,cli.cli.Peername)

	nw.Unlock()


	//check whether the there are no more peers in the Network. If so, delete it
	if len(nw.Networkmembers) == 0 {
		State.Lock()
		delete(State,cli.nw.Networkname)
		State.Unlock()
	}


}
