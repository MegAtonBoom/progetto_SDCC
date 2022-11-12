/*  SDCC PROJECT: MUTUAL EXCLUSION SIMULATION
 *	Author: Baba Adrian Petru
 *	Year: 2021/2022
 *	This distributed application is built to simulate 3 different distributed mutual exclusion algorithms:
 *	-Centralized token algorithm,
 *	-Decentralized token algorithm,
 *	-Ricart & agrawala decentralized algorithm
 *	Meant to be run with docker compose; read README file for further details
 */

package main

import (
	"errors"
	c "main/coordinator"
	r "main/registrator"
	ra "main/ricandagr"
	tc "main/tokencen"
	td "main/tokendec"
	u "main/utils"
	"os"
	"strconv"
)

//1st and only arg: service- 0:Registrator, 1:Coordinator (only for centralized token algorithm and therefore peers), 2: centralized token peer,
//3: decentralized token peer, 4:ricart & agrawala peer

func main() {

	var verbose bool = false
	var delay bool = false
	var service int
	var err error

	npeer := os.Getenv("NPEERS")
	u.PeerNumber, err = strconv.Atoi(npeer)
	u.LogError("Error in getting the number of involved peers: ", err, verbose)
	if u.PeerNumber < 1 {
		u.LogError("Error in getting the number of involved peers: ", errors.New("The number of peers can't be < 1!"), verbose)
	}

	if os.Getenv("VERBOSE") == "true" {
		verbose = true
	}

	if os.Getenv("DELAY") == "true" {
		delay = true
	}

	if len(os.Args) != 2 {
		u.LogError("Error in getting the input parameters: ", errors.New("There should be only one input parameter for the service requested."), verbose)
	}
	service, err = strconv.Atoi(os.Args[1])
	u.LogError("Error in getting the requested service by the numeric input: ", err, verbose)
	if service < 0 || service > 4 {
		u.LogError("Error in getting the requested service: ", errors.New("the value is not in the range 0-4."), verbose)
	}
	switch service {
	case 0:
		r.Registrator(verbose, delay)
		break
	case 1:
		c.Coordinator(verbose, delay)
		break
	case 2:
		tc.PeerSimulation(verbose, delay)
		break
	case 3:
		td.PeerSimulation(verbose, delay)
		break
	case 4:
		ra.PeerSimulation(verbose, delay)
		break
	default:
		break
	}

}
