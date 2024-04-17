package process

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	ProcessHandlesPlugin struct {
	}
)

func (volp *ProcessHandlesPlugin) GetName() string {
	return "PROCESS HANLDES PLUGIN"
}

func (volp *ProcessHandlesPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, "process_handles.txt")
}

func (volp *ProcessHandlesPlugin) Run() error {
	args := []string{"windows.handles.Handles"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	return nil
}
