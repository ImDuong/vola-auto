package network

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	NetstatPlugin struct {
	}
)

func (volp *NetstatPlugin) GetName() string {
	return "NETSTAT PLUGIN"
}

func (volp *NetstatPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, "netstat.txt")
}

func (volp *NetstatPlugin) Run() error {
	args := []string{"windows.netstat.NetStat"}
	return plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
}
