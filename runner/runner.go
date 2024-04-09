package runner

import (
	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/volatility/envars"
	"github.com/ImDuong/vola-auto/plugins/volatility/filescan"
	"github.com/ImDuong/vola-auto/plugins/volatility/help"
	"github.com/ImDuong/vola-auto/plugins/volatility/hivelist"
	"github.com/ImDuong/vola-auto/plugins/volatility/info"
	"github.com/ImDuong/vola-auto/plugins/volatility/lsadump"
	"github.com/ImDuong/vola-auto/plugins/volatility/mft"
	"github.com/ImDuong/vola-auto/plugins/volatility/netstat"
	"github.com/ImDuong/vola-auto/plugins/volatility/pe_version"
	"github.com/ImDuong/vola-auto/plugins/volatility/process"
	"github.com/ImDuong/vola-auto/utils"
	"github.com/alitto/pond"
	"go.uber.org/zap"
)

func RunPlugins() error {
	utils.Logger.Info("Start extracting")
	err := runVolatilityPlugins()
	if err != nil {
		return err
	}

	utils.Logger.Info("Start collecting")
	err = runCollectorPlugins()
	if err != nil {
		return err
	}

	utils.Logger.Info("Start analyzing")
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
		&process.ProcessCmdlinePlugin{},
		&process.ProcessPsListPlugin{},
		&process.ProcessPsScanPlugin{},
		&process.ProcessPsTreePlugin{},
		&envars.EnvarsPlugin{},
		&pe_version.PEVersionPlugin{},
		&filescan.FilescanPlugin{},
		&netstat.NetstatPlugin{},
		&hivelist.HivelistPlugin{},
		&lsadump.LsadumpPlugin{},
		&mft.MFTAdsPlugin{},
	}

	volPlgRunningPool := pond.New(10, 20)
	for _, plg := range volPlgs {
		if !plugins.IsRunRequired(plg.GetArtifactsExtractionPath()) {
			utils.Logger.Warn("Skipping", zap.String("plugin", plg.GetName()))
			continue
		}
		utils.Logger.Info("Running", zap.String("plugin", plg.GetName()))

		// if using the same plg variable for all tasks, the plg inside each task will change following the newest value of plg while looping
		// hence, copy the plugin inside each loop so each parallel task will have an indiviual plugin variable
		copiedPlg := plg
		volPlgRunningPool.Submit(func() {
			err := copiedPlg.Run()
			if err != nil {
				utils.Logger.Error("Running", zap.String("plugin", copiedPlg.GetName()), zap.Error(err))
				return
			}
			utils.Logger.Info("Finish", zap.String("plugin", copiedPlg.GetName()))
		})
	}
	volPlgRunningPool.StopAndWait()
	return nil
}
