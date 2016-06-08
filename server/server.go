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
		cli.createNetwork(message)
		cli.joinNetwork(message)
	case msg.JoinNwType:
		cli.joinNetwork(message)
	case msg.DisconnectNw:
		cli.disconnectNetwork()
	}
}
func (cli *Client) createNetwork(cmd msg.CreateNetworkCommand) {
	_, networkexists := State.m[cmd.Networkname]
	if networkexists {
		return
	}

	State.RLock()
	defer State.RUnlock()

	State.m[cmd.Networkname] = Network{
		Networkname:    cmd.Networkname,
		Networkmembers: make(map[string]msg.Peer),
		Networkpass:    cmd.Networkpass,
		Networkip:      cmd.Networkip,
	}

	cli.nw = msg.Network{Networkname: cmd.Networkname}
}

func (cli *Client) joinNetwork(cmd msg.JoinNetworkCommand) {
	_, networkexists := State.m[cmd.Networkname]
	if !networkexists {
		return
	}

	//get Client IP Address
	cli.cli.Peeraddress = strings.Split(cli.conn.RemoteAddr().String()[0], ":") + ":" + cmd.Peerport

	//check free IP address
	//TODO: Make it better
	ok := false
	var ip int
	for !ok {
		ok = true
		ip = rand.Intn(253) + 1

		for _, v := range State.m[cmd.Networkname].Networkmembers {
			cliaddr := v.cli.VirtualAddress[15]
			if cliaddr == ip {
				ok = false
				return
			}
		}
	}

	cli.cli.VirtualAddress = State.m[cmd.Networkname].Networkip
	cli.cli.VirtualAddress[15] = ip

	append(State.m[cmd.Networkname].Networkmembers, cli.cli)
	cli.nw.Networkname = cmd.Networkname
	append(cli.nw.Networkmembers, cli.cli)

	data := msg.EncodeNetworkJoinNotification(msg.NetworkJoinNotification{
		VirtualIPAddress: cli.cli.VirtualAddress,
		Peers:            cli.nw.Networkmembers,
	})

	//Notify all peers including the new

	for _, v := range State.m[cmd.Networkname].Networkmembers {
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

}
