package datastore

import (
	"fmt"
	"strings"
	"time"
)

type (
	TCPConnectionState string

	NetworkConnection struct {
		Protocol     string
		LocalAddr    string
		LocalPort    uint
		ForeignAddr  string
		ForeignPort  uint
		State        string
		CreatedTime  time.Time
		OwnerProcess *Process
	}
)

// TCPConnectionStates
// Ref: https://github.com/volatilityfoundation/volatility3/blob/771ed10b44573a7f8baa32822f3bc524195fe0c9/volatility3/framework/symbols/windows/netscan/netscan-win10-x64.json#L364
var TCPConnectionStates = [...]TCPConnectionState{
	"CLOSED",
	"LISTENING",
	"SYN_SENT",
	"SYN_RCVD",
	"ESTABLISHED",
	"FIN_WAIT1",
	"FIN_WAIT2",
	"CLOSE_WAIT",
	"CLOSING",
	"LAST_ACK",
	"TIME_WAIT",
	"DELETE_TCB",
}

var MissingInfoNetworkConnection []*NetworkConnection

func (nc *NetworkConnection) GetLocalSocketAddr() string {
	return nc.getSocketAddr(nc.LocalAddr, nc.LocalPort)
}

func (nc *NetworkConnection) GetForeignSocketAddr() string {
	return nc.getSocketAddr(nc.ForeignAddr, nc.ForeignPort)
}

func (nc *NetworkConnection) getSocketAddr(ipAddr string, port uint) string {
	if ipAddr == "::" {
		ipAddr = "[" + ipAddr + "]"
	}
	return fmt.Sprintf("%s:%d", ipAddr, port)
}

func (nc *NetworkConnection) GetCreatedTimeAsStr() string {
	var createdTime string = "N/A"
	if !nc.CreatedTime.IsZero() {
		createdTime = nc.CreatedTime.Format(time.DateTime)
	}
	return createdTime
}

func IsValidTCPConnectionState(checkingState string) bool {
	if len(checkingState) == 0 {
		return false
	}
	for i := range TCPConnectionStates {
		if strings.EqualFold(string(TCPConnectionStates[i]), checkingState) {
			return true
		}
	}
	return false
}
