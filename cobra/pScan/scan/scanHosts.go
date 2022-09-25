// Package scan provides types and functions to perform TCP port
// scans on a list of hosts
package scan

import (
	"fmt"
	"net"
	"time"
)

// Represents the state for a single TCP port
type PortState struct {
	Port int
	Open state
}

// Custom type state wrapper around bool type
// This allows us to define the String method to show open/closed
// instead of true/false when printing the value
type state bool

// String converts the boolean value of state to a human readable string
// The String method satisfies the Stringer interface, allowing you to use
// this type directly with print functions
func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

// scanPort performs a port scan on a single TCP port
func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	// DialTimeout attempts to connect to a network adddress within
	// a given time (1 sec in this case) then returns an error if it can't
	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)

	if err != nil {
		// assume port is closed, return PortState with default initial variable of false
		return p
	}

	// When connection succeeds, close connection
	scanConn.Close()
	p.Open = true
	return p
}

// Results represents the scan results for a single host
type Results struct {
	Host string
	// Can the host be resolved as a valid ip address
	NotFound   bool
	PortStates []PortState
}

// External function Run performs a port scan on the hosts list
// using private internal function scanPort to check each port
func Run(hl *HostsList, ports []int) []Results {
	// initialise slice of results with capacity as number of hosts in list
	res := make([]Results, 0, len(hl.Hosts))

	// loop through each host and define an instance of Results for each host
	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}

		// Resolve host name to a valid ip address
		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true
			res = append(res, r)

			// skip port scan on this host as unobtainable
			continue
		}

		// Execute port scan by looping through each port in ports slice
		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p))
		}

		res = append(res, r)
	}

	return res
}
