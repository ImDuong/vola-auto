package datastore

type (
	SystemMainProfile string

	MachineInfo struct {
		Is64Bit     bool
		NTBuildLab  string
		MainProfile SystemMainProfile
	}
)

const (
	Win10Profile SystemMainProfile = "win10"
	Win8Profile  SystemMainProfile = "win8"
	Win7Profile  SystemMainProfile = "win7"
)

var HostInfo MachineInfo
