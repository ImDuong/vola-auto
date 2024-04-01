package netstat

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	NetstatPlugin struct {
	}
)

const (
	PluginName = "NETSTAT PLUGIN"
)

var (
	artifactsExtractionFilename = "netstat.txt"
)

func (volp *NetstatPlugin) GetName() string {
	return PluginName
}

func (volp *NetstatPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, artifactsExtractionFilename)
}

func (volp *NetstatPlugin) SetArtifactsExtractionFilename(fileName string) {
	artifactsExtractionFilename = fileName
}

func (volp *NetstatPlugin) Run() error {
	args := []string{"windows.netstat.NetStat"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	args = []string{"windows.netscan.NetScan"}
	err = plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), false)
	if err != nil {
		return err
	}
	return nil
}
