package filescan

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	FilescanPlugin struct {
	}
)

const (
	PluginName = "FILESCAN PLUGIN"
)

var (
	artifactsExtractionFilename = "filescan.txt"
)

func (volp *FilescanPlugin) GetName() string {
	return PluginName
}

func (volp *FilescanPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, artifactsExtractionFilename)
}

func (volp *FilescanPlugin) SetArtifactsExtractionFilename(fileName string) {
	artifactsExtractionFilename = fileName
}

func (volp *FilescanPlugin) Run() error {
	args := []string{"windows.filescan.FileScan"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}
	return nil
}
