package eventlogs

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	EventLogsPlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *EventLogsPlugin) GetName() string {
	return "EVENT LOGS COLLECTION PLUGIN"
}

func (colp *EventLogsPlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "evtx")
}

func (colp *EventLogsPlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	prefetchFiles, err := filePlg.FindFilesByRegex("\\.evtx$")
	if err != nil {
		return err
	}

	err = filePlg.DumpFiles(prefetchFiles, colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}

	err = filePlg.RenameDumpedFilesExtention(".evtx.dat", "", colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}
	err = filePlg.RenameDumpedFilesExtention(".evtx.vacb", "", colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}

	return filePlg.ValidateDumpedFolder(colp.GetArtifactsCollectionPath())
}
