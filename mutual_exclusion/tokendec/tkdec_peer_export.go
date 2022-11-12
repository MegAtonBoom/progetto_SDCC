package tokendec

import (
	"errors"
	u "main/utils"
)

// Service exported
type Peer int

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
