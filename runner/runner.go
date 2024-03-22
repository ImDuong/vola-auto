package runner

import (
	"fmt"

	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/volatility/envars"
	"github.com/ImDuong/vola-auto/plugins/volatility/help"
	"github.com/ImDuong/vola-auto/plugins/volatility/info"
	"github.com/ImDuong/vola-auto/plugins/volatility/pe_version"
	"github.com/ImDuong/vola-auto/plugins/volatility/process"
	"github.com/alitto/pond"
)

func RunPlugins() error {
	err := runVolatilityPlugins()
	if err != nil {
		return err
	}
	err = runAnalyticPlugins()
	if err != nil {
		return err
	}
	return nil
}

func runVolatilityPlugins() error {
	volPlgs := []plugins.VolPlugin{
		&help.HelpPlugin{},
		&info.InfoPlugin{},
		&process.ProcessPlugin{},
		&envars.EnvarsPlugin{},
		&pe_version.PEVersionPlugin{},
	}

	volPlgRunningPool := pond.New(5, 20)
	for _, plg := range volPlgs {
		if !plugins.IsRunRequired(plg.GetArtifactsExtractionPath()) {
			fmt.Printf("Skipping plugin %s\n", plg.GetName())
			continue
		}
		fmt.Printf("Start running plugin %s\n", plg.GetName())
		volPlgRunningPool.Submit(func() {
			err := plg.Run()
			if err != nil {
				fmt.Printf("Running plugin %s got %s\n", plg.GetName(), err.Error())
				return
			}
			fmt.Printf("Finish running plugin %s\n", plg.GetName())
		})
	}
	volPlgRunningPool.StopAndWait()
	return nil
}
