package usnjrnl_j

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/ImDuong/vola-auto/utils"
	"github.com/alitto/pond"
	"go.uber.org/zap"
)

type (
	UsnJrnlJPlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *UsnJrnlJPlugin) GetName() string {
	return "UsnJrnl:$J COLLECTION PLUGIN"
}

func (colp *UsnJrnlJPlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "usnjrnl_j")
}

func (colp *UsnJrnlJPlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	foundFiles, err := filePlg.FindFilesByRegex(`\$UsnJrnl:\$J`)
	if err != nil {
		return err
	}

	err = filePlg.DumpFiles(foundFiles, colp.GetArtifactsCollectionPath())
	if err != nil {
		// edge case when vol3 try to name the dump file with `:`
		utils.Logger.Warn("Collecting artifacts", zap.String("plugin", colp.GetName()), zap.Error(err))
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
