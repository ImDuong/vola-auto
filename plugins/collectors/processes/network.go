package processes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/datastore"
)

type (
	NetworkPlugin struct {
	}

	PathToProcesses map[string]datastore.ProcessByPID
)

func (colp *NetworkPlugin) GetName() string {
	return "PROCESS WITH NETWORK COLLECTION PLUGIN"
}

func (colp *NetworkPlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(colp.GetArtifactsCollectionFolderpath(), "network.txt")
}

func (colp *NetworkPlugin) GetArtifactsCollectionFolderpath() string {
	return filepath.Join(config.Default.OutputFolder, ProcessCollectionFolderName)
}

func (colp *NetworkPlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionFolderpath(), 0755)
	if err != nil {
		return err
	}

	// group processes by process path
	procGroupedByPath := colp.groupProcessByPath()

	networkFileWriter, err := os.OpenFile(colp.GetArtifactsCollectionPath(), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer networkFileWriter.Close()

	// construct written data for a process group
	for procPath, procGroup := range procGroupedByPath {
		formattedData := colp.getFormattedDataForProcessGroup(procPath, procGroup)
		networkFileWriter.Write([]byte(formattedData))
	}

	return nil
}

func (colp *NetworkPlugin) groupProcessByPath() PathToProcesses {
	groupedProcesses := make(PathToProcesses)
	for _, proc := range datastore.PIDToProcess {
		groupedProcesses[proc.FullPath] = append(groupedProcesses[proc.FullPath], proc)
	}
	return groupedProcesses
}

// getFormattedDataForProcessGroup() constructs data following this schema:
//
//	{process_path}
//		PID: {PID_number} - {process_args}
//			{protocol} - {state} - {local_socket_addr} => {foreign_socket_addr} - {create_time}
//		PID:....
func (colp *NetworkPlugin) getFormattedDataForProcessGroup(processPath string, processGroup datastore.ProcessByPID) string {
	sort.Sort(processGroup)

	var result strings.Builder
	result.WriteString(processPath + "\n")

	for _, process := range processGroup {
		result.WriteString(fmt.Sprintf("\tPID: %d - %s\n", process.PID, process.Args))
		if process.Conn != nil {
			result.WriteString(fmt.Sprintf(
				"\t\t%s - %s - %s => %s - %s\n",
				process.Conn.Protocol,
				process.Conn.State,
				process.Conn.GetLocalSocketAddr(),
				process.Conn.GetForeignSocketAddr(),
				process.CreatedTime.Format(time.DateTime),
			))
		}
	}

	result.WriteString("\n")
	return result.String()
}
