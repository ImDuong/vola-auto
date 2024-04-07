package process

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	ProcessPsTreePlugin struct {
	}
)

func (volp *ProcessPsTreePlugin) GetName() string {
	return "PROCESS PSTREE PLUGIN"
}

func (volp *ProcessPsTreePlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, "process_pstree.txt")
}

func (volp *ProcessPsTreePlugin) Run() error {
	args := []string{"windows.pstree.PsTree"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	return nil
}
