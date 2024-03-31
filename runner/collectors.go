package runner

import (
	"fmt"

	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/ImDuong/vola-auto/plugins/collectors/amcache"
	"github.com/ImDuong/vola-auto/plugins/collectors/eventlogs"
	"github.com/ImDuong/vola-auto/plugins/collectors/logfile"
	"github.com/ImDuong/vola-auto/plugins/collectors/mft"
	"github.com/ImDuong/vola-auto/plugins/collectors/prefetch"
	"github.com/ImDuong/vola-auto/plugins/collectors/sru"
	"github.com/ImDuong/vola-auto/plugins/collectors/system32_config_hive"
	"github.com/ImDuong/vola-auto/plugins/collectors/usnjrnl_j"
	"github.com/alitto/pond"
)

func runCollectorPlugins() error {
	fmt.Println("STARTING COLLECTING")
	colPlgs := []plugins.CollectorPlugin{
		&collectors.MachinePlugin{},
		&collectors.FilesPlugin{},
	}

	for _, plg := range colPlgs {
		fmt.Printf("Start running plugin %s\n", plg.GetName())
		err := plg.Run()
		if err != nil {
			fmt.Printf("Running plugin %s got %s\n", plg.GetName(), err.Error())
			continue
		}
		fmt.Printf("Finish running plugin %s\n", plg.GetName())
	}

	colPlgRunningPool := pond.New(15, 100)
	colPlgs = []plugins.CollectorPlugin{
		&prefetch.PrefetchPlugin{
			WorkerPool: colPlgRunningPool,
		},
		&eventlogs.EventLogsPlugin{
			WorkerPool: colPlgRunningPool,
		},
		&system32_config_hive.HivePlugin{
			WorkerPool: colPlgRunningPool,
		},
		&mft.MFTPlugin{
			WorkerPool: colPlgRunningPool,
		},
		&usnjrnl_j.UsnJrnlJPlugin{
			WorkerPool: colPlgRunningPool,
		},
		&logfile.LogFilePlugin{
			WorkerPool: colPlgRunningPool,
		},
		&amcache.AmCachePlugin{
			WorkerPool: colPlgRunningPool,
		},
		&sru.SRUPlugin{
			WorkerPool: colPlgRunningPool,
		},
	}

	mainTaskGroup := colPlgRunningPool.Group()
	for _, plg := range colPlgs {
		if !plugins.IsRunRequired(plg.GetArtifactsCollectionPath()) {
			fmt.Printf("Skipping plugin %s\n", plg.GetName())
			continue
		}
		fmt.Printf("Start running plugin %s\n", plg.GetName())
		copiedPlg := plg
		mainTaskGroup.Submit(func() {
			err := copiedPlg.Run()
			if err != nil {
				fmt.Printf("Running plugin %s got %s\n", copiedPlg.GetName(), err.Error())
				return
			}
			fmt.Printf("Finish running plugin %s\n", copiedPlg.GetName())
		})
	}

	mainTaskGroup.Wait()
	colPlgRunningPool.StopAndWait()
	return nil
}
