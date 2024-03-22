package info

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	InfoPlugin struct {
	}
)

const (
	PluginName                  = "INFO PLUGIN"
	ArtifactsExtractionFilename = "info.txt"
)

func (volp *InfoPlugin) GetName() string {
	return PluginName
}

func (volp *InfoPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *InfoPlugin) Run() error {
	args := []string{config.Default.VolRunConfig.Binary, "-f", config.Default.MemoryDumpPath, "windows.info.Info"}
	return plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
}
