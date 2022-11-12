package ricandagr

import (
	u "main/utils"
	"net/rpc"
	"strconv"
)

// Service exported
type Peer int

// remote procedure to call when the caller want to send a request message to someone
func (prog *Peer) SendRequestMessage(request u.ReqMessage, reply *u.Reply) error {
	working = true
	var rep u.Reply
	//checking if the caller clock/identifier is greater, if so, or if we're in CS, just appending the request
	//else send him immediately the reply message
	v := isBigger(request)
	if state.String() == "CS" || (v == true && state.String() == "Requesting") {
		requests = append(requests, request)
		if state.String() == "CS" {
			u.WriteInLog("The peer got a request but he's in CS: he will append the request in the list and send the reply later", verbose)
		} else {
			u.WriteInLog("The peer got a request but he is already requesting and his request has a smaller timestamp or id: he will append the request in the list and send the reply later", verbose)
		}
	} else {
		address := allPeers[request.PeerIdentifier].Address
		peer, err := rpc.DialHTTP("tcp", address)
		defer peer.Close()
		u.LogError("Error in dialing: ", err, verbose)
		u.SimulateDelay(delay)
		err = peer.Call("Peer.SendReplyMessage", identifier, &rep)
		u.LogError("Error in calling the remote function: ", err, verbose)

		u.WriteInLog("The peer got a request and he can send the reply message", verbose)
	}
	//checking if we need to update our clock with the caller one
	if clock < request.Clock {

		u.WriteInLog("The peer is updating its old timestamp ("+strconv.Itoa(clock)+") with the requester one ("+strconv.Itoa(request.Clock)+")", verbose)

		clock = request.Clock
	}
	return nil
}

// procedure to call when the registrator wants to send the info about every registered peer
func (prog *Peer) SendRegisteredPeers(list []u.RegisterVoice, reply *u.Reply) error {
	working = true
	allPeers = list
	return nil
}

// procedure to call when the caller wants to send us his reply message
func (prog *Peer) SendReplyMessage(number int, reply *u.Reply) error {
	working = true

	u.WriteInLog("The peer got a reply message and is updating his old reply number ("+strconv.Itoa(numReplies)+") with ("+strconv.Itoa(numReplies+1)+")", verbose)

	numReplies++
	return nil
}

// just returns true if the requester clock- or eventually identifier- is greater than ours
func isBigger(request u.ReqMessage) bool {
	if request.Clock > myLastRequest {
		return true
	} else if request.Clock == myLastRequest && request.PeerIdentifier > identifier {
		return true
	}
	return false
}
