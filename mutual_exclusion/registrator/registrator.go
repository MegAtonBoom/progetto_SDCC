package registrator

import (
	u "main/utils"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

var verbose bool
var delay bool

// procedure used to activate the registrator
func Registrator(vFlag bool, dFlag bool) {

	verbose = vFlag
	delay = dFlag
	if verbose {
		u.Initialize(0)
	}
	registration := new(Registration)
	register := rpc.NewServer()
	err := register.RegisterName("Registration", registration)
	u.LogError("Error registering the registrator: ", err, verbose)
	register.HandleHTTP("/", "/debug")
	lis, err := net.Listen("tcp", u.RegPort)
	u.LogError("Listen error: ", err, verbose)
	u.WriteInLog("The registerer is up and ready to accept next "+strconv.Itoa(u.PeerNumber)+" peers requests", verbose)
	err = http.Serve(lis, nil)
	u.LogError("Registrator listening error: ", err, verbose)
}
