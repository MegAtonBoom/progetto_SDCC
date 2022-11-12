package registrator

import (
	"errors"
	u "main/utils"
	"net/rpc"
	"os"
)

// service exported
type Registration int

// array containing registered peers data
var registrations []u.RegisterVoice

// remote procedure called when a peer wants to register himself using the registrator
func (reg *Registration) RegisterPeer(address string, identifier *int) error {

	var insert bool
	var id int
	u.WriteInLog("Got a request to register from a peer with address "+address, verbose)
	if len(registrations) == u.PeerNumber {
		u.WriteInLog("Registerer already full; shutting down", verbose)
		return errors.New("Register full capacity!")

	}
	if insert, id = insertPeer(address); !insert {
		u.WriteInLog("Got a request from peer with ip "+address+" but he is already in the registered peers", verbose)
		return errors.New("Peer already inserted!")
	}
	if len(registrations) == u.PeerNumber {
		go sendRegistered()
	}
	*identifier = id
	return nil
}

// local procedure that tries to insert the caller peer in the data structure of the registrator, if it is not inserted yet
func insertPeer(address string) (bool, int) {
	for i := 0; i < len(registrations); i++ {
		if registrations[i].Address == address {
			return false, 0
		}
	}
	var number = len(registrations)
	var voice u.RegisterVoice

	voice.PeerId = number
	voice.Address = address
	registrations = append(registrations, voice)
	return true, number
}

// local procedure called when the registrator want to send the registrated perrs to everyone of them- the registrator then shuts down
func sendRegistered() {
	u.WriteInLog("Register full", verbose)
	for _, proc := range registrations {
		sendToPeer(proc)
	}
	u.WriteInLog("registerer sent data to peers and now is shutting down", verbose)
	os.Exit(0)
}

// local procedure that sends to a specific peer the data with every registered peer
func sendToPeer(proc u.RegisterVoice) {
	var reply u.Reply
	u.WriteInLog("Registerer sent data to peer with ip: "+proc.Address, verbose)
	client, err := rpc.DialHTTP("tcp", proc.Address)
	defer client.Close()
	u.LogError("Error in dialing: ", err, verbose)
	u.SimulateDelay(delay)
	err = client.Call("Peer.SendRegisteredPeers", registrations, &reply)
	u.LogError("Error in calling the remote procedure: ", err, verbose)
}
