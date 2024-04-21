package notifications

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	NotificationsPlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *NotificationsPlugin) GetName() string {
	return "NOTIFICATION COLLECTION PLUGIN"
}

func (colp *NotificationsPlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "notifications")
}

func (colp *NotificationsPlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	foundFiles, err := filePlg.FindFilesByRegex(`\\AppData\\Local\\Microsoft\\Windows\\Notifications`)
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
