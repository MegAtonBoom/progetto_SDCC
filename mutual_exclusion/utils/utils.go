package utils

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

/*
 *
 * CONSTANTS- not needed to be changed, are just dummy values for simulation purpose
 *
 */

// max duration of the delay between a CS completion and the next CS request-1
const MaxDelay = 4

// Max duration of the CS operation in seconds-1, for the simulation
const MaxCSDuration = 9

// number of accesses on CS per peer, if needed
const IterNumber = 5

// min time of delay- for the network congestion simulation
const MinCallDelay = 500

/*
 *
 * STRUCTS
 */

// request construct for token algorithms
type Request struct {
	Identificator  int
	VectorialClock VectorialClock
	Address        string
}

// the element of the registrator array that stores the registered peers
type RegisterVoice struct {
	PeerId  int
	Address string
}

// Needed for ricard & agrawala, for the management of the 3 states 
const (
        NCS State = iota
        Requesting
        CS
)

func (s State) String() string {
        switch s {
        case NCS:
                return "NCS"
        case Requesting:
                return "Requesting"
        case CS:
                return "CS"
        }
        return "unknown"
}

// request message structure, containing and identifier and the current clock of the sender
type ReqMessage struct {
        PeerIdentifier int
        Clock          int
}

/*
 *
 * OTHER VARIABLES
 *
 */

//still for ricar&agrawala
type State int

// dummy construct for the 2 variable -output- of an rpc procedure, when we don't need anything as output
type Reply bool

// the token construct for token algorithms
type Token bool

type VectorialClock []int

// exact number of peers involved in the mutual exclusion process (excluded the eventual coordinator)
var PeerNumber int

/*
 *
 *ADDRESSES AND PORTS
 */

// coordinator and registrator address and port
var CoordAddress string = "coordinator" + CoordPort
var CoordPort string = ":8100"

var RegAddress string = "registerer" + RegPort
var RegPort string = ":8200"

/*
 *
 * PROCEDURES
 *
 */
// global utility procedure that randomply generates a usable port in range 10k-11k
func GeneratePeerAddr() (int, string) {
	rand.Seed(time.Now().UnixNano())
	portNumber := rand.Intn(1000) + 10000
	myAddress := GetIP() + ":" + strconv.Itoa(portNumber)
	return portNumber, myAddress
}

// global procedure that retrieves the local ip
func GetIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// global utility procedure that logs an error on screen and eventually on the log file
func LogError(message string, err error, verbose bool) {
	if err != nil {
		WriteInLog("FATAL ERROR CAUSING SHUTDOWN ->"+message+err.Error(), verbose)
		time.Sleep(1 * time.Second)
		log.Fatal(message, err)
	}

}

// global utility procedure that simulates a really computing expensive critical section...
func SimulateCriticalSection() {
	time.Sleep(time.Duration(1+(rand.Intn(MaxCSDuration))) * time.Second)
}

// global utility procedure that simulates an eventual delay for the trasmission of any message
func SimulateDelay(delay bool) {
	if delay {
		time.Sleep(((time.Duration(rand.Intn(1000))) + MinCallDelay) * time.Millisecond)
	}
}
