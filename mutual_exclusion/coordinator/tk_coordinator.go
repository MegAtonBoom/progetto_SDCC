package coordinator

import (
	u "main/utils"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"
	
)

var verbose bool

var delay bool

var working = false

var wg sync.WaitGroup

// procedure to call just to activate the coordinator for the centralized token case
func Coordinator(vFlag bool, dFlag bool) {

	initializeClock()
	verbose = vFlag
	delay = dFlag
	if verbose {
		u.Initialize(1)
	}
	coordinator := new(Coordination)
	coord := rpc.NewServer()
	err := coord.RegisterName("Coordinator", coordinator)
	u.LogError("Error registering coordinator: ", err, verbose)
	coord.HandleHTTP("/", "/debug")
	lis, err := net.Listen("tcp", u.CoordPort)
	u.LogError("Listen error: ", err, verbose)
	u.WriteInLog("Coordinator avaiable, listening to messages from peers", verbose)
	wg.Add(1)
	go peerListener()
	go http.Serve(lis, nil)
	u.LogError("Registrator listening error: ", err, verbose)
	wg.Wait()

}

// procedure that fills the vectorial clock with 0's
func initializeClock() {
	for i := 0; i < u.PeerNumber; i++ {
		vectorialClock = append(vectorialClock, 0)
	}
}

// procedure meant to be run in a syncronized goroutine- therefore waited before shutting down- that allowa the peers/coordinator to stay up until everyone finishes
func peerListener() {
	
	for true {
		time.Sleep(15 * time.Second)
		if !working {
			u.WriteInLog("No requests came from peers for a long time: the coordinator is shutting down", verbose)
			wg.Done()
			return
		} else {
			working = false
		}
	}
}
