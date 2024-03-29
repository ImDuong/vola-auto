package process

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	ProcessPlugin struct {
	}
)

const (
	PluginName                  = "PROCESS PLUGIN"
	ArtifactsExtractionFilename = "process.txt"
)

func (volp *ProcessPlugin) GetName() string {
	return PluginName
}

func (volp *ProcessPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *ProcessPlugin) Run() error {
	args := []string{"windows.cmdline.CmdLine"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	args = []string{"windows.pstree.PsTree"}
	err = plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), false)
	if err != nil {
		return err
	}

	return nil
}
