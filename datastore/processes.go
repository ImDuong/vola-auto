package datastore

type (
	Process struct {
		ImageName  string
		PID        uint
		ParentProc *Process
		Args       string
	}
)

var ProcessList []Process
