package tokencen

import (
	"encoding/json"
	u "main/utils"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

var vectorialClock u.VectorialClock

var token u.Token

var myAddress string

var myPortNumber int

var identifier int

var allPeers []u.RegisterVoice

var verbose bool

var delay bool

var allReady = false

var working = false

var wg sync.WaitGroup

// procedure that starts the peer, sending the request to the registrator and waiting to get his response before starting the simulation
func PeerSimulation(vFlag bool, dFlag bool) {
	var tries = 0
	initializeClock()
	verbose = vFlag
	delay = dFlag
	if verbose {
		u.Initialize(2)
	}

	token = false

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

	//waiting for the registrator to send us the list of partecipating processes
	for len(allPeers) == 0 {
	}
	go peerListener()
	simulateProcess()

}

func simulateProcess() {
	var rep u.Reply
	var request u.Request
	u.WriteInLog("Peer got the list of peers and started the simulation", verbose)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < u.IterNumber; i++ {

		//Simulating some random operation by our procecc wich may take a while...before he needs to enter the critical section
		time.Sleep(time.Duration(1+(rand.Intn(u.MaxDelay))) * time.Second)

		//now he needs to enter the critical section
		//increment his vectorial clock
		vectorialClock[identifier]++

		request.Identificator = identifier
		request.VectorialClock = vectorialClock
		request.Address = myAddress

		u.WriteInLog("Requesting the token to the coordinator", verbose)
		go getToken(request)

		vc, _ := json.Marshal(vectorialClock)
		u.WriteInLog("Sending program messages to the other peers with the curren vectorial clock being: "+string(vc), verbose)

		go sendProgramMessages(request)

		//waiting for the token (if we already have the token, we'll not be waiting)
		u.WriteInLog("Waiting for the token from the coordinator", verbose)

		for !token {
		}

		u.WriteInLog("****START CRITICAL SECTION****", verbose)
		u.SimulateCriticalSection()
		u.WriteInLog("****END CRITICAL SECTION****", verbose)

		token = false
		u.WriteInLog("Returning token to the coordinator", verbose)
		u.SimulateDelay(delay)

		coord, err := rpc.DialHTTP("tcp", u.CoordAddress)
		defer coord.Close()
		u.LogError("Error in dialing: ", err, verbose)
		err = coord.Call("Coordinator.ReturnToken", identifier, &rep)
		u.LogError("Error in calling the remote procedure: ", err, verbose)

	}

	u.WriteInLog("This peer finished the simulation but will stay up waiting for the others to do the same", verbose)
	wg.Wait()

}

func getToken(request u.Request) {

	var tokenReply u.Reply

	u.SimulateDelay(delay)
	coord, err := rpc.DialHTTP("tcp", u.CoordAddress)
	defer coord.Close()
	u.LogError("Error in dialing: ", err, verbose)
	err = coord.Call("Coordinator.GetToken", request, &tokenReply)
	u.LogError("Error in calling the remote procedure: ", err, verbose)
	token = u.Token(tokenReply)

}

// sending the program messages
func sendProgramMessages(request u.Request) {

	var reply u.Reply
	var nTries = 1
	for _, peer := range allPeers {
		if peer.PeerId == identifier {
			continue
		}
		u.SimulateDelay(delay)
		if !allReady {
			u.WriteInLog("First connection: trying to connect to "+peer.Address+" for the "+strconv.Itoa(nTries)+" time", verbose)
		} else {
			u.WriteInLog("Trying to connect to "+peer.Address, verbose)
		}

		client, err := rpc.DialHTTP("tcp", peer.Address)
		defer client.Close()
		//just for the first time we send program messages, we try again 10 times- other peers might not be already running so we give them some time
		if err != nil && !allReady {
			for nTries = 2; nTries <= 10; nTries++ {
				u.SimulateDelay(delay)
				time.Sleep(150 * time.Millisecond)
				u.WriteInLog("Trying to connect to "+peer.Address+" for the "+strconv.Itoa(nTries)+" time", verbose)
				client, err = rpc.DialHTTP("tcp", peer.Address)
				if err == nil {
					break
				}
			}
		}
		u.LogError("Error in dialing ", err, verbose)
		err = client.Call("Peer.SendProgramMessage", request, &reply)
		u.WriteInLog("Sent program message to peer with address: "+peer.Address, verbose)
	}
	allReady = true
}

// just a procedure that creates the listening peer, for token exchange rpc calls
func registerPeerService() {

	peer := new(Peer)
	peerServ := rpc.NewServer()
	err := peerServ.RegisterName("Peer", peer)
	u.LogError("Error while registering the service name: ", err, verbose)
	peerServ.HandleHTTP("/", "/debug")
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(myPortNumber))
	u.LogError("Listening error", err, verbose)
	go http.Serve(lis, nil)
	u.WriteInLog("Centralized token peer registered with port: "+strconv.Itoa(myPortNumber), verbose)

}

func initializeClock() {
	for i := 0; i < u.PeerNumber; i++ {
		vectorialClock = append(vectorialClock, 0)
	}
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
