package system32_config_hive

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	HivePlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *HivePlugin) GetName() string {
	return "SYSTEM32 CONFIG HIVE COLLECTION PLUGIN"
}

func (colp *HivePlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "system32_config_hive")
}

func (colp *HivePlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	hivePrefix := "\\\\Windows\\\\System32\\\\config\\\\"
	hiveRegexes := []string{
		hivePrefix + "SAM",
		hivePrefix + "SECURITY",
		hivePrefix + "SOFTWARE",
		hivePrefix + "SYSTEM",
	}
	for i := range hiveRegexes {
		foundFiles, err := filePlg.FindFilesByRegex(hiveRegexes[i])
		if err != nil {
			fmt.Println("[WARNING] Cannot find files by regex", err)
			continue
		}

		err = filePlg.DumpFiles(foundFiles, colp.GetArtifactsCollectionPath())
		if err != nil {
			fmt.Println("[WARNING] Cannot dump files", err)
			continue
		}

		err = filePlg.RenameDumpedFilesExtention(".dat", "", colp.GetArtifactsCollectionPath())
		if err != nil {
			fmt.Println("[WARNING] Cannot rename files", err)
			continue
		}
	}

	return nil
}
