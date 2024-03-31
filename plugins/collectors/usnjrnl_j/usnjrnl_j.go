package usnjrnl_j

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
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
		fmt.Println("edge case", err)
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
