package collectors

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/ImDuong/vola-auto/datastore"
	"github.com/ImDuong/vola-auto/plugins/volatility/info"
	"github.com/ImDuong/vola-auto/utils"
	"go.uber.org/zap"
)

type (
	MachinePlugin struct {
	}
)

func (colp *MachinePlugin) GetName() string {
	return "MACHINE INFO COLLECTION PLUGIN"
}

func (colp *MachinePlugin) GetArtifactsCollectionPath() string {
	return ""
}

func (colp *MachinePlugin) Run() error {
	correspPlg := info.InfoPlugin{}
	filescanArtifactFiles, err := os.Open(correspPlg.GetArtifactsExtractionPath())
	if err != nil {
		return err
	}
	defer filescanArtifactFiles.Close()
	scanner := bufio.NewScanner(filescanArtifactFiles)
	isProcessDataFound := false

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}
		if !isProcessDataFound {
			if strings.Contains(line, "Variable") && strings.Contains(line, "Value") {
				isProcessDataFound = true
			}
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		if parts[0] == "Is64Bit" {
			is64Bit, err := strconv.ParseBool(parts[1])
			if err != nil {
				return err
			}
			datastore.HostInfo.Is64Bit = is64Bit
			continue
		}

		if parts[0] == "NTBuildLab" {
			datastore.HostInfo.NTBuildLab = parts[1]
			if strings.Contains(datastore.HostInfo.NTBuildLab, "win7") {
				datastore.HostInfo.MainProfile = datastore.Win7Profile
			} else if strings.Contains(datastore.HostInfo.NTBuildLab, "win10") {
				datastore.HostInfo.MainProfile = datastore.Win10Profile
			}
			continue
		}

		if parts[0] == "NtMajorVersion" && len(datastore.HostInfo.MainProfile) == 0 {
			if parts[1] == "10" {
				datastore.HostInfo.MainProfile = datastore.Win10Profile
			}
		}
	}

	if err := scanner.Err(); err != nil {
		utils.Logger.Warn("Collecting artifacts", zap.String("plugin", colp.GetName()), zap.Error(err))
	}
	return nil
}
