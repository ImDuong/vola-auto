package process

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	ProcessPsListPlugin struct {
	}
)

func (volp *ProcessPsListPlugin) GetName() string {
	return "PROCESS PSLIST PLUGIN"
}

func (volp *ProcessPsListPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, "process_pslist.txt")
}

func (volp *ProcessPsListPlugin) Run() error {
	args := []string{"windows.pslist.PsList"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	return nil
}
