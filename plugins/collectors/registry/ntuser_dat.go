package registry

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	NTUserDatPlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *NTUserDatPlugin) GetName() string {
	return "NTUSER.DAT PLUGIN"
}

func (colp *NTUserDatPlugin) GetArtifactsCollectionPath() string {
	// if there is no dumped files, can look for /collectors/registry/hivelist if such registries are present in memory image
	return filepath.Join(config.Default.OutputFolder, "ntuser_dat")
}

func (colp *NTUserDatPlugin) Run() error {

	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	ntuserFiles, err := filePlg.FindFilesByRegex(`\\NTUSER\.DAT`)
	if err != nil {
		return err
	}

	err = filePlg.DumpFiles(ntuserFiles, colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}

	err = filePlg.RenameDumpedFilesExtention(".dat", "", colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}

	return nil
}
