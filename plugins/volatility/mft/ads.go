package mft

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	MFTAdsPlugin struct {
	}
)

const (
	PluginName                  = "MFT_ADS PLUGIN"
	ArtifactsExtractionFilename = "mft_ads.txt"
)

func (volp *MFTAdsPlugin) GetName() string {
	return PluginName
}

func (volp *MFTAdsPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, ArtifactsExtractionFilename)
}

func (volp *MFTAdsPlugin) Run() error {
	args := []string{"windows.mftscan.ADS"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	return nil
}
