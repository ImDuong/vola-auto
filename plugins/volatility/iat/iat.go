package iat

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	IATPlugin struct {
	}
)

const (
	PluginName                  = "IAT PLUGIN"
	ArtifactsExtractionFilename = "iat.txt"
)

func (volp *IATPlugin) GetName() string {
	return PluginName
}

func (volp *IATPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *IATPlugin) Run() error {
	args := []string{"windows.iat.IAT"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	return nil
}
