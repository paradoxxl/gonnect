package server

import (
	"net"
	"fmt"
	"github.com/paradoxxl/gonnect/msg"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [1600]byte
	for {
		//fmt.Println("Trying to read")
		_, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
		}

		//Evaluate the command
		_,cmdtype,err := msg.CheckMsg(buf)
		if err != nil{
			fmt.Printf("Handleclient %v",err)
			return
		}

		switch cmdtype{
		case msg.CreateNwType:
			//create network
		case msg.JoinNwType:
			//join the network
		case msg.DisconnectNw:
			//Disconnect the Network
		}

		/*
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}*/
	}
}

func createNetwork(cmd msg.CreateNetworkCommand){

}

func joinNetwork(cmd msg.CreateNetworkCommand){

}

func disconnectNetwork(cmd msg.CreateNetworkCommand){

}