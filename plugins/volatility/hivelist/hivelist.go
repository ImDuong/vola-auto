package hivelist

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	HivelistPlugin struct {
	}
)

const (
	PluginName                  = "HIVELIST PLUGIN"
	ArtifactsExtractionFilename = "hivelist.txt"
)

func (volp *HivelistPlugin) GetName() string {
	return PluginName
}

func (volp *HivelistPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *HivelistPlugin) Run() error {
	args := []string{"windows.registry.hivelist.HiveList"}
	return plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
}
