package collectors

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ImDuong/vola-auto/datastore"
	"github.com/ImDuong/vola-auto/plugins/volatility/info"
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
				datastore.HostInfo.Profile = "win7"
			} else if strings.Contains(datastore.HostInfo.NTBuildLab, "win10") {
				datastore.HostInfo.Profile = "win10"
			}
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(colp.GetName(), ":got some errors when collecting artifacts")
	}
	return nil
}
