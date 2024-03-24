package runner

import (
	"fmt"

	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/ImDuong/vola-auto/plugins/collectors/prefetch"
	"github.com/alitto/pond"
)

func runCollectorPlugins() error {
	fmt.Println("STARTING COLLECTING")
	volPlgs := []plugins.CollectorPlugin{
		&collectors.MachinePlugin{},
		&collectors.FilesPlugin{},
	}

	for _, plg := range volPlgs {
		fmt.Printf("Start running plugin %s\n", plg.GetName())
		err := plg.Run()
		if err != nil {
			fmt.Printf("Running plugin %s got %s\n", plg.GetName(), err.Error())
			continue
		}
		fmt.Printf("Finish running plugin %s\n", plg.GetName())
	}

	volPlgs = []plugins.CollectorPlugin{
		&prefetch.PrefetchPlugin{},
	}

	volPlgRunningPool := pond.New(5, 20)
	for _, plg := range volPlgs {
		if !plugins.IsRunRequired(plg.GetArtifactsCollectionPath()) {
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
