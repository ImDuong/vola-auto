package datastore

import "time"

type (
	Process struct {
		ImageName   string
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
