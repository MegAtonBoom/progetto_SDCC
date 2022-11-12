package coordinator

import (
	u "main/utils"
	"net/rpc"
	"strconv"
)

var requests []u.Request

// Service exported
type Coordination int

var coordinatorToken u.Token = true

var vectorialClock u.VectorialClock

// remote procedure to call when a peer wants to ask the coordinator for the token
func (coord *Coordination) GetToken(request u.Request, clientToken *bool) error {
	working = true
	eligible := eligibleRequest(request)
	if eligible {
		u.WriteInLog("Got a token request from peer with ip: "+request.Address+" and the request is eligible", verbose)
	} else {
		u.WriteInLog("Got a token request from peer with ip: "+request.Address+" and the request is not eligible", verbose)
	}
	if coordinatorToken == true && eligible {
		coordinatorToken = false
		*clientToken = true

		updateTimestamp(request)

		u.WriteInLog("The coordinator has also the token: sending token to the requester", verbose)

	} else {
		requests = append(requests, request)
		*clientToken = false
		u.WriteInLog("The coordinator has not the token or the request is not eligible: it was appended on the request list wich is now size "+strconv.Itoa(len(requests)), verbose)

	}
	return nil
}

// remote procedure called to return the token to the coordinator
func (coord *Coordination) ReturnToken(peer int, reply *u.Reply) error {
	working = true
	coordinatorToken = true

	u.WriteInLog("Got the token back: checking for pending eligible requests", verbose)

	nextRequest()
	return nil
}

// procedure that updates the timestamp
func updateTimestamp(request u.Request) {

	u.WriteInLog("The coordinator updated his timestamp", verbose)

	vectorialClock[request.Identificator] = request.VectorialClock[request.Identificator]

}

// procedure that checks if the procedure is eligible based on its timestamp
func eligibleRequest(request u.Request) bool {

	eligible := true
	for i, time := range vectorialClock {
		if request.VectorialClock[i] > time && i != request.Identificator {
			eligible = false
		}
	}

	return eligible
}

// the procedure that find the next eligible request and send the token- when the coordinator gets the token back
func nextRequest() {

	for i, req := range requests {
		if eligible := eligibleRequest(req); eligible == true && coordinatorToken {
			requests = append(requests[:i], requests[i+1:]...)
			updateTimestamp(req)
			coordinatorToken = false
			sendToken(req)
			break
		}
	}
}

// procedure that sends the token to one of the peers
func sendToken(request u.Request) {

	// Try to connect to addr using HTTP protocol
	client, err := rpc.DialHTTP("tcp", request.Address)
	defer client.Close()
	u.LogError("Error in dialing: ", err, verbose)
	u.SimulateDelay(delay)
	err = client.Call("Peer.SendToken", true, nil)

	u.WriteInLog("Found an eligible request: sending token to peer with ip: "+request.Address, verbose)
	u.LogError("Error in calling the remote procedure: ", err, verbose)
}
