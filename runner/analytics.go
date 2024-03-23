package runner

import (
	"fmt"

	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/analytics/envars"
	"github.com/alitto/pond"
)

func runAnalyticPlugins() error {
	fmt.Println("STARTING ANALYZING")
	volPlgs := []plugins.AnalyticPlugin{
		&envars.EnvarsPlugin{},
	}

	volPlgRunningPool := pond.New(5, 20)
	for _, plg := range volPlgs {
		if !plugins.IsRunRequired(plg.GetAnalyticResultPath()) {
			fmt.Printf("Skipping plugin %s\n", plg.GetName())
			continue
		}
		fmt.Printf("Start running plugin %s\n", plg.GetName())
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
