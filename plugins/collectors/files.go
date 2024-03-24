package collectors

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ImDuong/vola-auto/datastore"
	"github.com/ImDuong/vola-auto/plugins/volatility/filescan"
)

type (
	FilesPlugin struct {
	}
)

const (
	PluginName = "FILES COLLECTION PLUGIN"
)

func (colp *FilesPlugin) GetName() string {
	return PluginName
}

func (colp *FilesPlugin) GetArtifactsCollectionPath() string {
	return ""
}

// this plugin only process & store info about files in memory, not dump files
func (colp *FilesPlugin) Run() error {
	correspPlg := filescan.FilescanPlugin{}
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
			if strings.Contains(line, "Offset") && strings.Contains(line, "Name") && strings.Contains(line, "Size") {
				isProcessDataFound = true
			}
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		fileInfo := datastore.FileInfo{
			Path: parts[1],
		}

		if datastore.HostInfo.Profile == "win10" {
			fileInfo.VirtualAddrOffset = parts[0]
		} else {
			fileInfo.PhysicalAddrOffset = parts[0]
		}

		datastore.FileList = append(datastore.FileList, fileInfo)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(PluginName, ":got some errors when collecting artifacts")
	}
	return nil
}

// TODO: dump all files and put them in original folder structure
func (colp *FilesPlugin) DumpAllFiles() error {
	return nil
}
