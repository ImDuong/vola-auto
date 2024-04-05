package sru

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	SRUPlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *SRUPlugin) GetName() string {
	return "SRU COLLECTION PLUGIN"
}

func (colp *SRUPlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "sru")
}

func (colp *SRUPlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	foundFiles, err := filePlg.FindFilesByRegex(`\\Windows\\System32\\sru\\SRU`)
	if err != nil {
		return err
	}

	err = filePlg.DumpFiles(foundFiles, colp.GetArtifactsCollectionPath())
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

	return nil
}
