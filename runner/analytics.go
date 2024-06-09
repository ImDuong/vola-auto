package runner

import (
	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/analytics/envars"
	"github.com/ImDuong/vola-auto/utils"
	"github.com/alitto/pond"
	"go.uber.org/zap"
)

func runAnalyticPlugins() error {
	volPlgs := []plugins.AnalyticPlugin{
		&envars.EnvarsPlugin{},
	}

	volPlgRunningPool := pond.New(5, 20)
	for _, plg := range volPlgs {
		if !plugins.IsRunRequired(plg.GetAnalyticResultPath()) {
			utils.Logger.Warn("Skipping", zap.String("plugin", plg.GetName()))
			continue
		}
		utils.Logger.Info("Starting", zap.String("plugin", plg.GetName()))
		copiedPlg := plg
		volPlgRunningPool.Submit(func() {
			err := copiedPlg.Run()
			if err != nil {
				utils.Logger.Error("Starting", zap.String("plugin", copiedPlg.GetName()), zap.Error(err))
				return
			}
			utils.Logger.Info("Finished", zap.String("plugin", copiedPlg.GetName()))
		})
	}
	volPlgRunningPool.StopAndWait()
	return nil
}
