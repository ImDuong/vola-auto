package processes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

	networkFileWriter, err := os.OpenFile(colp.GetArtifactsCollectionPath(), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer networkFileWriter.Close()

	// construct written data for a process group
	for procPath, procGroup := range procGroupedByPath {
		formattedData := colp.getFormattedDataForProcessGroup(procPath, procGroup)
		networkFileWriter.Write([]byte(formattedData))
	}

	networkFileWriter.Write([]byte("\nMissing Information Network Connection\n"))
	for i := range datastore.MissingInfoNetworkConnection {
		networkFileWriter.Write([]byte(getFormattedDataForMissingInfoConnection(datastore.MissingInfoNetworkConnection[i])))
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

	if len(strings.TrimSpace(processPath)) == 0 {
		processPath = "cannot_parse_process"
	}

	var result strings.Builder
	isNetConnAvail := false
	for _, process := range processGroup {
		if process.Conn == nil {
			continue
		}

		if !isNetConnAvail {
			// write only once to the start of the string
			result.WriteString(processPath + "\n")
			isNetConnAvail = true
		}

		result.WriteString(fmt.Sprintf("\tPID: %-4d - %s\n", process.PID, process.GetCmdline()))
		result.WriteString(fmt.Sprintf(
			"\t\t%s - %s - %s => %s - %s\n",
			process.Conn.Protocol,
			process.Conn.State,
			process.Conn.GetLocalSocketAddr(),
			process.Conn.GetForeignSocketAddr(),
			process.Conn.GetCreatedTimeAsStr(),
		))
	}

	if isNetConnAvail {
		result.WriteString("\n")
	}
	return result.String()
}

func getFormattedDataForMissingInfoConnection(nc *datastore.NetworkConnection) string {
	return fmt.Sprintf(
		"\t\t%-8s - %-11s - %-44s => %-44s\n",
		nc.Protocol,
		nc.State,
		nc.GetLocalSocketAddr(),
		nc.GetForeignSocketAddr(),
	)
}
