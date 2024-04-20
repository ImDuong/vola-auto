package mft

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	MFTPlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *MFTPlugin) GetName() string {
	return "MFT COLLECTION PLUGIN"
}

func (colp *MFTPlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "mft")
}

func (colp *MFTPlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	mftFiles, err := filePlg.FindFilesByRegex("mft")
	if err != nil {
		return err
	}

	err = filePlg.DumpFiles(mftFiles, colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}

	err = filePlg.RenameDumpedFilesExtention(".dat", "", colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}
	err = filePlg.RenameDumpedFilesExtention(".vacb", "", colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}
	err = filePlg.RenameDumpedFilesExtention(".img", "", colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}

	return nil
}
