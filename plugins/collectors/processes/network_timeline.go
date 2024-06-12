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

	var sortingConnList []*datastore.NetworkConnection
	for i := range datastore.PIDToProcess {
		if len(datastore.PIDToProcess[i].Connections) == 0 {
			continue
		}
		sortingConnList = append(sortingConnList, datastore.PIDToProcess[i].Connections...)
	}

	sort.Slice(sortingConnList, func(i, j int) bool {
		if sortingConnList[i].CreatedTime.IsZero() {
			return false
		}
		return sortingConnList[i].CreatedTime.After(sortingConnList[j].CreatedTime)
	})

	for _, conn := range sortingConnList {
		networkTimelineFileWriter.Write([]byte(fmt.Sprintf(
			"%-29s - %-4d - %-25s - %-8s - %-11s - %-44s => %-44s - %s\n",
			conn.GetCreatedTimeAsStr(),
			conn.OwnerProcess.PID,
			conn.OwnerProcess.ImageName,
			conn.Protocol,
			conn.State,
			conn.GetLocalSocketAddr(),
			conn.GetForeignSocketAddr(),
			conn.OwnerProcess.GetCmdline(),
		)))
	}

	networkTimelineFileWriter.Write([]byte("\nMissing Information Network Connection\n"))
	for i := range datastore.MissingInfoNetworkConnection {
		networkTimelineFileWriter.Write([]byte(getFormattedDataForMissingInfoConnection(datastore.MissingInfoNetworkConnection[i])))
	}

	return nil
}
