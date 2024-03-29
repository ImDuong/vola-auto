package logfile

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	LogFilePlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *LogFilePlugin) GetName() string {
	return "LogFile COLLECTION PLUGIN"
}

func (colp *LogFilePlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "logfile")
}

func (colp *LogFilePlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	logFiles, err := filePlg.FindFilesByRegex(`\$LogFile`)
	if err != nil {
		return err
	}

	err = filePlg.DumpFiles(logFiles, colp.GetArtifactsCollectionPath())
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
