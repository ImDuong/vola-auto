package datastore

import (
	"strings"
	"time"

	"github.com/ImDuong/vola-auto/utils"
)

type (
	Process struct {
		ImageName   string
		FullPath    string
		PID         uint
		ParentProc  *Process
		Args        string
		CreatedTime time.Time
		Connections []*NetworkConnection
	}
)

var PIDToProcess map[uint]*Process = make(map[uint]*Process)

type ProcessByPID []*Process

func (a ProcessByPID) Len() int           { return len(a) }
func (a ProcessByPID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProcessByPID) Less(i, j int) bool { return a[i].PID < a[j].PID }

func (p *Process) ParseFullPathByArgs() {
	fullPath := strings.TrimSpace(p.Args)
	p.FullPath = fullPath
	if len(fullPath) == 0 {
		return
	}

	if fullPath[0] == '"' {
		endIdx := strings.Index(fullPath[1:], "\"") + 1
		if endIdx > 0 {
			p.FullPath = fullPath[1:endIdx]
			return
		}
	} else {
		endIdx := strings.Index(fullPath, " ")
		if endIdx > 0 {
			p.FullPath = fullPath[:endIdx]
			return
		}
	}
}

func (p *Process) GetFullPath() string {
	return utils.GetPathInCamelCase(p.FullPath)
}

func (p *Process) GetCmdline() string {
	if len(p.Args) == 0 {
		return p.FullPath
	}
	return p.Args
}

func (p *Process) AddConn(conn *NetworkConnection) {
	if conn == nil {
		return
	}
	p.Connections = append(p.Connections, conn)
}

func (p *Process) IsConnExisted(conn *NetworkConnection) bool {
	if conn == nil {
		return false
	}
	for i := range p.Connections {
		if p.Connections[i].GetSocketPair() == conn.GetSocketPair() {
			return true
		}
	}
	return false
}
