package tokencen

import (
	"encoding/json"
	"errors"
	u "main/utils"
)

// Service exported
type Peer int

// remote procedure to call when you want to send a program message to a peer
func (proc *Peer) SendProgramMessage(request u.Request, reply *u.Reply) error {
	working = true
	var validMessage = false
	for _, peer := range allPeers {
		if peer.Address == request.Address {
			validMessage = true
			break
		}
	}
	if validMessage {
		ec, _ := json.Marshal(request.VectorialClock)
		vc, _ := json.Marshal(vectorialClock)
		u.WriteInLog("Got a program message from peer with address "+request.Address+" with clock: "+string(ec)+" while mine is: "+string(vc), verbose)
		for i, time := range vectorialClock {

			if request.VectorialClock[i] > time {
				vectorialClock[i] = request.VectorialClock[i]
			}

		}

		vc, _ = json.Marshal(vectorialClock)
		u.WriteInLog("After the program message from "+request.Address+" the updated clock is "+string(vc), verbose)

	} else {
		u.WriteInLog("Got an unvalid program message prom the ip: "+request.Address, verbose)
	}

	return nil
}

// remote procedure to invoke when you want to send the token to a peer
func (proc *Peer) SendToken(extToken bool, reply *u.Reply) error {
	working = true
	reply = nil
	if extToken {
		token = true
		return nil
	} else {
		return errors.New("Invocation error")
	}
}

// procedure to call when the registrator wants to send the info about every registered peer
func (prog *Peer) SendRegisteredPeers(list []u.RegisterVoice, reply *u.Reply) error {
	working = true
	allPeers = list
	return nil
}
