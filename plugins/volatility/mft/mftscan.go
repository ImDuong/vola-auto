package mft

import (
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
)

type (
	MFTScanPlugin struct {
	}
)

func (volp *MFTScanPlugin) GetName() string {
	return "MFT_SCAN PLUGIN"
}

func (volp *MFTScanPlugin) GetArtifactsExtractionPath() string {
	return filepath.Join(config.Default.OutputFolder, "mft_scan.txt")
}

func (volp *MFTScanPlugin) Run() error {
	args := []string{"windows.mftscan.MFTScan"}
	err := plugins.RunVolatilityPluginAndWriteResult(args, volp.GetArtifactsExtractionPath(), config.Default.IsForcedRerun)
	if err != nil {
		return err
	}

	return nil
}
