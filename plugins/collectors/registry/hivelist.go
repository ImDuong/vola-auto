package registry

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	HiveListPlugin struct {
	}
)

func (colp *HiveListPlugin) GetName() string {
	return "HIVELIST PLUGIN"
}

func (colp *HiveListPlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "hivelist")
}

func (colp *HiveListPlugin) Run() error {

	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	args := []string{config.Default.VolRunConfig.Binary,
		"-f", config.Default.MemoryDumpPath,
		"-o", colp.GetArtifactsCollectionPath(),
		"windows.registry.hivelist.HiveList",
		"--dump",
	}
	return plugins.RunVolatilityPluginAndWriteResult(args, "", true)
}
