package envars

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	EnvarsPlugin struct {
	}
)

const (
	PluginName         = "ENVARS PLUGIN"
	AnalyticResultPath = "envars.txt"
)

func (anp *EnvarsPlugin) GetName() string {
	return PluginName
}

func (anp *EnvarsPlugin) GetAnalyticResultPath() string {
	return filepath.Join(config.Default.AnalyticFolder, AnalyticResultPath)
}

func (anp *EnvarsPlugin) Run() error {
	args := []string{config.Default.VolRunConfig.Binary, "-f", config.Default.MemoryDumpPath, "windows.envars.Envars"}
	return plugins.RunVolatilityPluginAndWriteResult(args, anp.GetAnalyticResultPath(), config.Default.IsForcedRerun)
}
