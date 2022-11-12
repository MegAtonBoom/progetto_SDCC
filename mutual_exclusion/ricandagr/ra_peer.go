package ricandagr

import (
	u "main/utils"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

// variables related to the ricart & agrawala algorithm
var numReplies int = 0

var state u.State = 0

var requests []u.ReqMessage

var myLastRequest int = 0

var clock int = 0

var identifier int = 0

var allPeers []u.RegisterVoice

var myAddress string

var myPortNumber int

var verbose bool

var delay bool

var allReady = false

var working = false

var wg sync.WaitGroup

// procedure that starts the peer, sending the request to the registrator and waiting to get his response before starting the simulation
func PeerSimulation(vFlag bool, dFlag bool) {
	var tries = 0
	verbose = vFlag
	delay = dFlag
	if verbose {
		u.Initialize(4)
	}

	myPortNumber, myAddress = u.GeneratePeerAddr()

	registerPeerService()

	registrator, err := rpc.DialHTTP("tcp", u.RegAddress)
	if err!=nil {
                for ;tries<10;tries++{
                        u.SimulateDelay(delay)
                        time.Sleep(150*time.Millisecond)
                        registrator, err = rpc.DialHTTP("tcp", u.RegAddress)
                        if err == nil {
                        break
                        }
                }
    	}
	defer registrator.Close()
	u.LogError("Error in dialing: ", err, verbose)

	u.SimulateDelay(delay)
	err = registrator.Call("Registration.RegisterPeer", myAddress, &identifier)
	u.LogError("Error in calling the remote procedure: ", err, verbose)

	u.WriteInLog("Peer registered with identifier: "+strconv.Itoa(identifier), verbose)

	//waiting for the registrator to send us the list of partecipating peers
	for len(allPeers) == 0 {
	}
	go peerListener()

	simulateProcess()

}

// procedure that starts the simulation of the ricart & agrawala algorithm
func simulateProcess() {

	u.WriteInLog("Peer got the list of peers and started the simulation", verbose)
	for i := 0; i < u.IterNumber; i++ {
		//randomness to generate the delay between one CS request and another
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(1+(rand.Intn(u.MaxDelay))) * time.Second)

		//now we need to get in the CS- REQUESTING
		state = 1

		u.WriteInLog("The peer is now in state \"REQUESTING\": he will send the request messages", verbose)
		clock++
		myLastRequest = clock
		sendRequestMessages()

		u.WriteInLog("The peer is waiting to get all the "+strconv.Itoa(u.PeerNumber-1)+" replies", verbose)

		//waiting to have all the replies befor getting in critical section
		for numReplies != u.PeerNumber-1 {
		}

		//CS
		state = 2

		u.WriteInLog("****START CRITICAL SECTION**** state: \"CS\"", verbose)
		u.SimulateCriticalSection()
		u.WriteInLog("****END CRITICAL SECTION**** state: \"NCS\"", verbose)
		//outside of CS
		state = 0
		numReplies = 0

		sendReplyMessages()
	}

	u.WriteInLog("This peer finished the simulation but will stay up waiting for the others to do the same", verbose)
	wg.Wait()
}

// procedure that sends request messages to every other peer known to partecipate to the mutual exclusion process
func sendRequestMessages() {

	var request u.ReqMessage
	var nTries = 1
	request.PeerIdentifier = identifier
	request.Clock = clock
	var reply u.Reply

	for i, proc := range allPeers {
		//just skipping ourself
		if proc.Address == myAddress {
			continue
		}
		u.SimulateDelay(delay)
		if !allReady {
			u.WriteInLog("First connection: trying to connect to "+proc.Address+" for the "+strconv.Itoa(nTries)+" time", verbose)
		} else {
			u.WriteInLog("Trying to connect to "+proc.Address, verbose)
		}
		peer, err := rpc.DialHTTP("tcp", proc.Address)
		//just for the first time we send request messages, we try again 10 times- other peers might not be already running so we give them some time
		if err != nil && !allReady {
			for nTries = 2; nTries <= 10; nTries++ {
				u.SimulateDelay(delay)
				time.Sleep(150 * time.Millisecond)
				u.WriteInLog("Trying to connect to "+proc.Address+" for the "+strconv.Itoa(nTries)+" time", verbose)

				peer, err = rpc.DialHTTP("tcp", proc.Address)
				if err == nil {
					break
				}
			}
		}
		defer peer.Close()
		u.LogError("Error in dialing: ", err, verbose)
		err = peer.Call("Peer.SendRequestMessage", request, &reply)

		u.WriteInLog("Request message number "+strconv.Itoa(i)+" sent to peer with ip: "+proc.Address, verbose)

		u.LogError("Error in calling the remote procedure: ", err, verbose)
	}
	allReady = true
}

// function to send the reply messages to everyone else- after the caller has used the CS
func sendReplyMessages() {

	var reply u.Reply
	for i, req := range requests {
		address := allPeers[req.PeerIdentifier].Address
		peer, err := rpc.DialHTTP("tcp", address)
		defer peer.Close()
		u.LogError("Error in dialing: ", err, verbose)
		u.SimulateDelay(delay)
		err = peer.Call("Peer.SendReplyMessage", identifier, &reply)
		u.LogError("Error in calling the remote procedure: ", err, verbose)

		u.WriteInLog("The peer correctly sent the reply message number "+strconv.Itoa(i)+" to peer with ip: "+address, verbose)

	}

	//just emptying the request slice
	requests = make([]u.ReqMessage, 0)
}

// Just to create the peer service listening for RQUESTS and REPLIES- using the previous random generated port
func registerPeerService() {

	peer := new(Peer)
	peerServ := rpc.NewServer()
	err := peerServ.RegisterName("Peer", peer)
	u.LogError("Error while registering the service name: ", err, verbose)
	peerServ.HandleHTTP("/", "/debug")
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(myPortNumber))
	u.LogError("Listening error", err, verbose)
	go http.Serve(lis, nil)

}

// procedure meant to be run in a syncronized goroutine- therefore waited before shutting down- that allowa the peers/coordinator to stay up until everyone finishes
func peerListener() {
	wg.Add(1)
	for true {
		time.Sleep(15 * time.Second)
		if !working {
			u.WriteInLog("No requests came from peers for a long time: this peer is shutting down", verbose)
			wg.Done()
			return
		} else {
			working = false
		}
	}
}
