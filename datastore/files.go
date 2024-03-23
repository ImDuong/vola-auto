package datastore

type (
	FileInfo struct {
		Path               string
		PhysicalAddrOffset string
		VirtualAddrOffset  string
	}
)

var FileList []FileInfo
