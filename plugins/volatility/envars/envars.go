package envars

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	EnvarsPlugin struct {
	}
)

const (
	PluginName                  = "ENVARS PLUGIN"
	ArtifactsExtractionFilename = "envars.txt"
)

func (volp *EnvarsPlugin) GetName() string {
	return PluginName
}

func (volp *EnvarsPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *EnvarsPlugin) Run() error {
	args := []string{"windows.envars.Envars"}
	return plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
}
