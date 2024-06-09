package processes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/datastore"
)

type (
	NetworkTimelinePlugin struct {
	}
)

func (colp *NetworkTimelinePlugin) GetName() string {
	return "PROCESS WITH NETWORK TIMELINE COLLECTION PLUGIN"
}

func (colp *NetworkTimelinePlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(colp.GetArtifactsCollectionFolderpath(), "network_timeline.txt")
}

func (colp *NetworkTimelinePlugin) GetArtifactsCollectionFolderpath() string {
	return filepath.Join(config.Default.OutputFolder, ProcessCollectionFolderName)
}

func (colp *NetworkTimelinePlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionFolderpath(), 0755)
	if err != nil {
		return err
	}

	networkTimelineFileWriter, err := os.OpenFile(colp.GetArtifactsCollectionPath(), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer networkTimelineFileWriter.Close()

	tempProcList := make([]*datastore.Process, len(datastore.PIDToProcess))
	idx := 0
	for i := range datastore.PIDToProcess {
		tempProcList[idx] = datastore.PIDToProcess[i]
		idx++
	}

	sort.Slice(tempProcList, func(i, j int) bool {
		if tempProcList[i].CreatedTime.IsZero() {
			return false
		}
		return tempProcList[i].CreatedTime.After(tempProcList[j].CreatedTime)
	})

	for _, proc := range tempProcList {
		if len(proc.Connections) == 0 {
			continue
		}

		for conIdx := range proc.Connections {
			networkTimelineFileWriter.Write([]byte(fmt.Sprintf(
				"%-29s - %-4d - %-25s - %-8s - %-11s - %-44s => %-44s - %s\n",
				proc.Connections[conIdx].GetCreatedTimeAsStr(),
				proc.PID,
				proc.ImageName,
				proc.Connections[conIdx].Protocol,
				proc.Connections[conIdx].State,
				proc.Connections[conIdx].GetLocalSocketAddr(),
				proc.Connections[conIdx].GetForeignSocketAddr(),
				proc.GetCmdline(),
			)))
		}

	}

	networkTimelineFileWriter.Write([]byte("\nMissing Information Network Connection\n"))
	for i := range datastore.MissingInfoNetworkConnection {
		networkTimelineFileWriter.Write([]byte(getFormattedDataForMissingInfoConnection(datastore.MissingInfoNetworkConnection[i])))
	}

	return nil
}
