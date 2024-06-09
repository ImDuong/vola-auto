package datastore

import (
	"fmt"
	"time"
)

type (
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

func (nc *NetworkConnection) GetLocalSocketAddr() string {
	return nc.getSocketAddr(nc.LocalAddr, nc.LocalPort)
}

func (nc *NetworkConnection) GetForeignSocketAddr() string {
	return nc.getSocketAddr(nc.ForeignAddr, nc.ForeignPort)
}

func (nc *NetworkConnection) getSocketAddr(ipAddr string, port uint) string {
	if ipAddr == "::" {
		ipAddr += " "
	}
	return fmt.Sprintf("%s:%d", ipAddr, port)
}
