package server



import (
	"github.com/paradoxxl/gonnect/msg"
	"net"
	"strings"
	"math/rand"
)

func (peer *PeerState) OnCommand(commandType byte, message interface{}) {
	switch commandType {
	case msg.CreateNwType:
		data,_ := message.(msg.CreateNetworkCommand)
		peer.createNetwork(data.Networkname,data.Networkpass,data.Networkip,data.Peername,data.Peerport)
		peer.joinNetwork(data.Networkname,data.Networkpass,data.Peername,data.Peerport)
	case msg.JoinNwType:
		data,_ := message.(msg.JoinNetworkCommand)
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
		Networkmembers: make(map[string]PeerState),
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
	ipadd := strings.Split(peer.conn.RemoteAddr().String(), ":")[0]
	udpaddr,_ := net.ResolveUDPAddr("tcp",ipadd+":"+Peerport)
	peer.peer.Peeraddress =  *udpaddr

	//check free IP address. Select one at random, check for collisions. increment on collision
	ok := false
	var ip = rand.Intn(253) + 1
	for !ok {
		ok = true
		for _, v := range State.m[Networkname].Networkmembers {
			cliaddr := v.peer.VirtualAddress[3]
			if int(cliaddr) == ip || ip == 0 {
				ok = false
				ip = (ip+1)%255
				break
			}
		}
	}


	peer.peer.VirtualAddress = State.m[Networkname].Networkip.To4()
	peer.peer.VirtualAddress[3] = byte(ip)
	peer.network.Networkname = Networkname

	State.RLock()
	nw := State.m[Networkname]
	State.RUnlock()

	nw.Lock()
	nw.Networkmembers[peer.peer.Peername] = *peer
	nw.Unlock()


	peers := []*msg.Peer{}
	for _,k := range nw.Networkmembers {
		peers = append(peers, &k.peer)
	}

	data := msg.EncodeNetworkJoinNotification(&msg.NetworkJoinNotification{
		VirtualIPAddress: peer.peer.VirtualAddress,
		Peers:            peers,
	})

	//Notify all peers including the new

	for _, v := range State.m[Networkname].Networkmembers {
		n, err := v.conn.Write(data)
		if err != nil {
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
		delete(State.m, peer.network.Networkname)
		State.Unlock()
	}else{
		nw.RLock()
		data := msg.EncodeClientDisconnectNotification(&msg.ClientDisconnectNotification{
			VirtualIPAddress: peer.peer.VirtualAddress,
		})

		for _, v := range nw.Networkmembers {
			n, err := v.conn.Write(data)
			if err != nil {
				//TODO: Do something
			}
			if n != len(data) {
				//TODO: Try resend
			}
		}
		nw.Unlock()
	}



}