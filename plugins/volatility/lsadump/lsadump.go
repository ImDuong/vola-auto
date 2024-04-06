package lsadump

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	LsadumpPlugin struct {
	}
)

const (
	PluginName                  = "LSADUMP PLUGIN"
	ArtifactsExtractionFilename = "lsadump.txt"
)

func (volp *LsadumpPlugin) GetName() string {
	return PluginName
}

func (volp *LsadumpPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *LsadumpPlugin) Run() error {
	args := []string{"windows.lsadump.Lsadump"}
	return plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
}
