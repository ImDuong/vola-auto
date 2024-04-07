package process

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	ProcessCmdlinePlugin struct {
	}
)

const (
	PluginName                  = "PROCESS CMDLINE PLUGIN"
	ArtifactsExtractionFilename = "process_cmdline.txt"
)

func (volp *ProcessCmdlinePlugin) GetName() string {
	return PluginName
}

func (volp *ProcessCmdlinePlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *ProcessCmdlinePlugin) Run() error {
	args := []string{"windows.cmdline.CmdLine"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	return nil
}
