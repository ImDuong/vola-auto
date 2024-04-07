package process

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	ProcessPsScanPlugin struct {
	}
)

func (volp *ProcessPsScanPlugin) GetName() string {
	return "PROCESS PSSCAN PLUGIN"
}

func (volp *ProcessPsScanPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, "process_psscan.txt")
}

func (volp *ProcessPsScanPlugin) Run() error {
	args := []string{"windows.psscan.PsScan"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	return nil
}
