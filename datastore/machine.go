package datastore

type (
	MachineInfo struct {
		Is64Bit    bool
		NTBuildLab string
		Profile    string // win10, or win7
	}
)

var HostInfo MachineInfo
