package server

import (
	"github.com/paradoxxl/gonnect/msg"
	"crypto/rand"
	"net"
	"crypto/tls"
	"fmt"
	"log"
	"strconv"
)

type GonnectServer struct{
	listener net.Listener
}


func NewGonnectServer(pubkPath,privkPath string,port int) GonnectServer{

	cert,err := tls.LoadX509KeyPair(pubkPath,privkPath)
	if err != nil{
		fmt.Printf("NewServer - Load Cert %v\n",err)
		return GonnectServer{}
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		Rand: rand.Reader,
	}

	service := "0.0.0.0:"+strconv.Itoa(port)

	fmt.Printf("Goint to listen on %v\n",service)
	listener,err:=tls.Listen("tcp",service,&config)
	if err != nil {log.Fatalf("server: listen: %s", err)}
	log.Printf("server is listening: %v\n" ,listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		go handleClient(conn)
	}




	return  GonnectServer{listener:listener}
}


func (srv GonnectServer) Close(){
	//TODO: Handle everything
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var client *PeerState

	for {
		//Evaluate the command
		msg.ReadMessage(conn, client)

	}
}

