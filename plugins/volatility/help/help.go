package help

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	HelpPlugin struct {
	}
)

const (
	PluginName                  = "HELP PLUGIN"
	ArtifactsExtractionFilename = "help.txt"
)

func (volp *HelpPlugin) GetName() string {
	return PluginName
}

func (volp *HelpPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *HelpPlugin) Run() error {
	args := []string{config.Default.VolRunConfig.Binary, "-h"}
	return plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
}
