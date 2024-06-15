package runner

import (
	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/ImDuong/vola-auto/plugins/collectors/amcache"
	"github.com/ImDuong/vola-auto/plugins/collectors/eventlogs"
	"github.com/ImDuong/vola-auto/plugins/collectors/logfile"
	"github.com/ImDuong/vola-auto/plugins/collectors/mft"
	"github.com/ImDuong/vola-auto/plugins/collectors/notifications"
	"github.com/ImDuong/vola-auto/plugins/collectors/prefetch"
	"github.com/ImDuong/vola-auto/plugins/collectors/processes"
	"github.com/ImDuong/vola-auto/plugins/collectors/registry"
	"github.com/ImDuong/vola-auto/plugins/collectors/sru"
	"github.com/ImDuong/vola-auto/plugins/collectors/system32_config_hive"
	"github.com/ImDuong/vola-auto/plugins/collectors/usnjrnl_j"
	"github.com/ImDuong/vola-auto/utils"
	"github.com/alitto/pond"
	"go.uber.org/zap"
)

func runCollectorPlugins() error {
	colPlgs := []plugins.CollectorPlugin{
		&collectors.MachinePlugin{},
		&collectors.FilesPlugin{},
		&collectors.ProcessesPlugin{},
	}

	for _, plg := range colPlgs {
		utils.Logger.Info("Starting", zap.String("plugin", plg.GetName()))
		err := plg.Run()
		if err != nil {
			utils.Logger.Error("Starting", zap.String("plugin", plg.GetName()), zap.Error(err))
			continue
		}
		utils.Logger.Info("Finished", zap.String("plugin", plg.GetName()))
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
		&notifications.NotificationsPlugin{
			WorkerPool: colPlgRunningPool,
		},
		&processes.TreePlugin{},
		&processes.TimelinePlugin{},
		&processes.NetworkPlugin{},
		&processes.NetworkTimelinePlugin{},
		&registry.HiveListPlugin{},
		&registry.NTUserDatPlugin{
			WorkerPool: colPlgRunningPool,
		},
	}

	// empty file collector plugin to validate dumped folder
	filePlg := collectors.FilesPlugin{}

	mainTaskGroup := colPlgRunningPool.Group()
	for _, plg := range colPlgs {
		if !plugins.IsRunRequired(plg.GetArtifactsCollectionPath()) {
			utils.Logger.Warn("Skipping", zap.String("plugin", plg.GetName()))
			continue
		}
		utils.Logger.Info("Starting", zap.String("plugin", plg.GetName()))
		copiedPlg := plg
		mainTaskGroup.Submit(func() {
			err := copiedPlg.Run()
			if err != nil {
				utils.Logger.Error("Starting", zap.String("plugin", copiedPlg.GetName()), zap.Error(err))
				return
			}

			err = filePlg.ValidateDumpedFiles(copiedPlg.GetArtifactsCollectionPath())
			if err != nil {
				utils.Logger.Error("Validate dumped folder", zap.String("plugin", copiedPlg.GetName()), zap.Error(err))
				return
			}

			utils.Logger.Info("Finished", zap.String("plugin", copiedPlg.GetName()))
		})
	}

	mainTaskGroup.Wait()
	colPlgRunningPool.StopAndWait()
	return nil
}
