package tokendec

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

var allPeers []u.RegisterVoice

var ready = false

var myAddress string

var myPortNumber int

var identifier int

var token u.Token = false

var verbose bool

var delay bool

var nextReady = false

var working = false

var wg sync.WaitGroup

// procedure that starts the peer, sending the request to the registrator and waiting to get his response before starting the simulation
func PeerSimulation(vFlag bool, dFlag bool) {
	var tries = 0
	verbose = vFlag
	delay = dFlag
	if verbose {
		u.Initialize(3)
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

// the procedure that actually simulates the mutual exclusion process
func simulateProcess() {
	var nTries = 1
	var reply u.Reply
	u.WriteInLog("Peer got the list of peers and started the simulation", verbose)
	address, err := getNextAddress()
	u.LogError("Error in getting the next peer address: ", err, verbose)
	rand.Seed(time.Now().UnixNano())

	if identifier == allPeers[0].PeerId {
		token = true
	}

	for i := 0; i < u.IterNumber; i++ {
		//waiting to get the token

		u.WriteInLog("Waiting the token for the neighbor peer", verbose)

		for !token {
		}

		//in this simulation we don't have a delay but we don't know
		//if the current peer actually needs to get into CS: we simulate this situation
		//by getting a pseudorandom number 0<=n<10 and getting into CS if it's <5
		if rand.Intn(10) < 7 {

			u.WriteInLog("****START CRITICAL SECTION**** the peer needs to get into critical section", verbose)
			u.SimulateCriticalSection()
			u.WriteInLog("****END CRITICAL SECTION**** the peer needed to get into critical section", verbose)

		} else {
			u.WriteInLog("****START/END CRITICAL SECTION**** the peer doesn't need to get into critical section", verbose)
		}
		token = false
		u.SimulateDelay(delay)
		if !nextReady {
			u.WriteInLog("First connection: trying to connect to "+address+" for the "+strconv.Itoa(nTries)+" time", verbose)
		} else {
			u.WriteInLog("Trying to connect to "+address, verbose)
		}
		client, err := rpc.DialHTTP("tcp", address)
		defer client.Close()
		//just for the first time we send the token to the next peer in the ring, we try again till 10 times- the other peer might not be already running so we give him some time
		if err != nil && !nextReady {
			for nTries = 2; nTries <= 10; nTries++ {
				u.SimulateDelay(delay)
				time.Sleep(150 * time.Millisecond)
				u.WriteInLog("Trying to connect to "+address+" for the "+strconv.Itoa(nTries)+" time", verbose)

				client, err = rpc.DialHTTP("tcp", address)
				if err == nil {
					break
				}
			}
		}
		u.LogError("Error in dialing: ", err, verbose)

		u.WriteInLog("Sending token to the neighbor "+address, verbose)

		err = client.Call("Peer.SendToken", true, &reply)
		u.LogError("Error in calling the remote procedure: ", err, verbose)
		nextReady = true
	}

	u.WriteInLog("This peer finished the simulation but will stay up waiting for the others to do the same", verbose)
	wg.Wait()

}

// procedure called to find who's the next peer connected to the current
// one in the ring, based on the identifier value
func getNextAddress() (string, error) {

	u.WriteInLog("Peer retrieving his neighbor in the ring", verbose)

	if identifier == u.PeerNumber-1 {
		return allPeers[0].Address, nil
	} else {
		return allPeers[identifier+1].Address, nil
	}
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

}

// procedure meant to be run in a syncronized goroutine- therefore waited before shutting down- that allowa the peers/coordinator to stay up until everyone finishes
func peerListener() {
	wg.Add(1)
	for true {
		time.Sleep(time.Duration(u.PeerNumber) * 15 * time.Second)
		if !working {
			u.WriteInLog("No requests came from peers for a long time: this peer is shutting down", verbose)
			wg.Done()
			return
		} else {
			working = false
		}
	}
}
