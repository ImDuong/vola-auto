package amcache

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	AmCachePlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *AmCachePlugin) GetName() string {
	return "AmCache COLLECTION PLUGIN"
}

func (colp *AmCachePlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "amcache")
}

func (colp *AmCachePlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	foundFiles, err := filePlg.FindFilesByRegex(`\\Windows\\appcompat\\Programs\\Amcache.hve`)
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

	return filePlg.ValidateDumpedFolder(colp.GetArtifactsCollectionPath())
}
