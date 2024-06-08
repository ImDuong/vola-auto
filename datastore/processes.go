package datastore

import (
	"strings"
	"time"
)

type (
	Process struct {
		ImageName   string
		FullPath    string
		PID         uint
		ParentProc  *Process
		Args        string
		CreatedTime time.Time
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
