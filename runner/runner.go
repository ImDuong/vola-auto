package runner

import (
	"fmt"

	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/volatility/envars"
	"github.com/ImDuong/vola-auto/plugins/volatility/filescan"
	"github.com/ImDuong/vola-auto/plugins/volatility/help"
	"github.com/ImDuong/vola-auto/plugins/volatility/hivelist"
	"github.com/ImDuong/vola-auto/plugins/volatility/info"
	"github.com/ImDuong/vola-auto/plugins/volatility/netstat"
	"github.com/ImDuong/vola-auto/plugins/volatility/pe_version"
	"github.com/ImDuong/vola-auto/plugins/volatility/process"
	"github.com/alitto/pond"
)

func RunPlugins() error {
	err := runVolatilityPlugins()
	if err != nil {
		return err
	}
	err = runCollectorPlugins()
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
	fmt.Println("STARTING EXTRACTING")
	volPlgs := []plugins.VolPlugin{
		&help.HelpPlugin{},
		&info.InfoPlugin{},
		&process.ProcessPlugin{},
		&envars.EnvarsPlugin{},
		&pe_version.PEVersionPlugin{},
		&filescan.FilescanPlugin{},
		&netstat.NetstatPlugin{},
		&hivelist.HivelistPlugin{},
	}

	volPlgRunningPool := pond.New(5, 20)
	for _, plg := range volPlgs {
		if !plugins.IsRunRequired(plg.GetArtifactsExtractionPath()) {
			fmt.Printf("Skipping plugin %s\n", plg.GetName())
			continue
		}
		fmt.Printf("Start running plugin %s\n", plg.GetName())

		// if using the same plg variable for all tasks, the plg inside each task will change following the newest value of plg while looping
		// hence, copy the plugin inside each loop so each parallel task will have an indiviual plugin variable
		copiedPlg := plg
		volPlgRunningPool.Submit(func() {
			err := copiedPlg.Run()
			if err != nil {
				fmt.Printf("Running plugin %s got %s\n", copiedPlg.GetName(), err.Error())
				return
			}
			fmt.Printf("Finish running plugin %s\n", copiedPlg.GetName())
		})
	}
	volPlgRunningPool.StopAndWait()
	return nil
}
