package network

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	NetscanPlugin struct {
	}
)

func (volp *NetscanPlugin) GetName() string {
	return "NETSCAN PLUGIN"
}

func (volp *NetscanPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, "netscan.txt")
}

func (volp *NetscanPlugin) Run() error {
	args := []string{"windows.netscan.NetScan"}
	return plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
}
