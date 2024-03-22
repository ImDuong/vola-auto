package pe_version

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	PEVersionPlugin struct {
	}
)

const (
	PluginName                  = "PE_VERSION PLUGIN"
	ArtifactsExtractionFilename = "pe_version.txt"
)

func (volp *PEVersionPlugin) GetName() string {
	return PluginName
}

func (volp *PEVersionPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *PEVersionPlugin) Run() error {
	args := []string{config.Default.VolRunConfig.Binary, "-f", config.Default.MemoryDumpPath, "windows.verinfo.VerInfo"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	return nil
}
