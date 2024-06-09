package collectors

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ImDuong/vola-auto/datastore"
	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/volatility/network"
	"github.com/ImDuong/vola-auto/plugins/volatility/process"
	"github.com/ImDuong/vola-auto/utils"
	"github.com/alitto/pond"
	"go.uber.org/zap"
)

type (
	ProcessesPlugin struct {
		WorkerPool *pond.WorkerPool
	}
)

func (colp *ProcessesPlugin) GetName() string {
	return "PROCESSES COLLECTION PLUGIN"
}

func (colp *ProcessesPlugin) GetArtifactsCollectionPath() string {
	return ""
}

// Run() only processes & stores info about files in memory, not dump files
// 1. Read list of processes from cmdline
// 2. Based on pslist, construct process relation from parent to child
func (colp *ProcessesPlugin) Run() error {
	// read processes from cmdline
	cmdlinePlg := process.ProcessCmdlinePlugin{}
	cmdlineArtifactFiles, err := os.Open(cmdlinePlg.GetArtifactsExtractionPath())
	if err != nil {
		return err
	}
	defer cmdlineArtifactFiles.Close()
	scanner := bufio.NewScanner(cmdlineArtifactFiles)
	isProcessDataFound := false

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}
		if !isProcessDataFound {
			if strings.Contains(line, "PID") && strings.Contains(line, "Process") && strings.Contains(line, "Args") {
				isProcessDataFound = true
			}
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		parsedPID, err := strconv.Atoi(parts[0])
		if err != nil {
			utils.Logger.Warn("parse pid failed", zap.String("pid", parts[0]), zap.String("plugin", colp.GetName()), zap.Error(err))
			continue
		}

		// TODO: handle process name have spaces
		proc := datastore.Process{
			ImageName: parts[1],
			PID:       uint(parsedPID),
		}

		if proc.PID != 4 {
			proc.Args = strings.Join(parts[2:], " ")

			proc.ParseFullPathByArgs()
			if len(proc.FullPath) == 0 {
				proc.FullPath = proc.ImageName
			}
		} else {
			proc.FullPath = `C:\Windows\System32\ntoskrnl.exe`
		}

		datastore.PIDToProcess[proc.PID] = &proc
	}

	if err := scanner.Err(); err != nil {
		utils.Logger.Warn("Collecting processes", zap.String("plugin", colp.GetName()), zap.Error(err))
	}

	if err := colp.constructProcessRelation(); err != nil {
		utils.Logger.Warn("Collecting processes relation", zap.String("plugin", colp.GetName()), zap.Error(err))
	}

	if err := colp.retrieveNetworkObjects(&network.NetstatPlugin{}); err != nil {
		utils.Logger.Warn("Collecting processes network objects", zap.String("plugin", colp.GetName()), zap.Error(err))
	}

	if err := colp.retrieveNetworkObjects(&network.NetscanPlugin{}); err != nil {
		utils.Logger.Warn("Collecting processes network objects", zap.String("plugin", colp.GetName()), zap.Error(err))
	}

	return nil
}

// constructProcessRelation() read pslist to construct parent-child relation
func (colp *ProcessesPlugin) constructProcessRelation() error {
	pslistPlg := process.ProcessPsListPlugin{}
	pslistArtifactFiles, err := os.Open(pslistPlg.GetArtifactsExtractionPath())
	if err != nil {
		return err
	}
	defer pslistArtifactFiles.Close()
	scanner := bufio.NewScanner(pslistArtifactFiles)
	isProcessDataFound := false

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}
		if !isProcessDataFound {
			if strings.Contains(line, "PID") && strings.Contains(line, "PPID") && strings.Contains(line, "ImageFileName") {
				isProcessDataFound = true
			}
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		parsedPID, err := strconv.Atoi(parts[0])
		if err != nil {
			utils.Logger.Warn("parse pid failed", zap.String("pid", parts[0]), zap.String("plugin", colp.GetName()), zap.Error(err))
			continue
		}

		parsedPPID, err := strconv.Atoi(parts[1])
		if err != nil {
			utils.Logger.Warn("parse ppid failed", zap.String("ppid", parts[1]), zap.String("plugin", colp.GetName()), zap.Error(err))
			continue
		}

		if _, ok := datastore.PIDToProcess[uint(parsedPID)]; !ok {
			utils.Logger.Warn("pid not found in current datastore", zap.Int("pid", parsedPID), zap.String("plugin", colp.GetName()))
			continue
		}
		if _, ok := datastore.PIDToProcess[uint(parsedPPID)]; !ok {
			parentProc := datastore.Process{
				PID: uint(parsedPPID),
			}
			if parentProc.PID == 0 {
				parentProc.ImageName = "System Idle Process"
			}
			datastore.PIDToProcess[uint(parsedPPID)] = &parentProc
		}

		datastore.PIDToProcess[uint(parsedPID)].ParentProc = datastore.PIDToProcess[uint(parsedPPID)]
	}

	if err := scanner.Err(); err != nil {
		utils.Logger.Warn("Constructing process relations", zap.String("plugin", colp.GetName()), zap.Error(err))
	}
	return nil
}

func (colp *ProcessesPlugin) retrieveNetworkObjects(netPlg plugins.VolPlugin) error {
	netstatArtifactFiles, err := os.Open(netPlg.GetArtifactsExtractionPath())
	if err != nil {
		return err
	}
	defer netstatArtifactFiles.Close()
	scanner := bufio.NewScanner(netstatArtifactFiles)
	isNetworkObjectDataFound := false

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}
		if !isNetworkObjectDataFound {
			if strings.Contains(line, "Offset") && strings.Contains(line, "Proto") && strings.Contains(line, "LocalAddr") {
				isNetworkObjectDataFound = true
			}
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 9 {
			continue
		}

		parsedProtocol := parts[1]

		isTCPConnection := false
		if strings.Contains(strings.ToLower(parsedProtocol), "tcp") {
			if len(parts) < 10 {
				continue
			}
			isTCPConnection = true
		}

		parsedLocalIP := parts[2]
		parsedLocalPort, err := strconv.Atoi(parts[3])
		if err != nil {
			utils.Logger.Warn("parse local port failed", zap.String("local_port", parts[3]), zap.String("plugin", colp.GetName()), zap.Error(err))
			continue
		}

		parsedForeignIP := parts[4]
		parsedForeignPort, err := strconv.Atoi(parts[5])
		if err != nil {
			utils.Logger.Warn("parse foreign port failed", zap.String("foreign_port", parts[5]), zap.String("plugin", colp.GetName()), zap.Error(err))
			continue
		}

		var state string
		var pidIdx int
		if isTCPConnection {
			state = parts[6]
			pidIdx = 7

			if !datastore.IsValidTCPConnectionState(state) {
				utils.Logger.Warn("invalid TCP connection state", zap.String("state", state), zap.String("plugin", colp.GetName()), zap.Error(err))
			}
		} else {
			state = "stateless"
			pidIdx = 6
		}

		netObj := datastore.NetworkConnection{
			Protocol:    parsedProtocol,
			LocalAddr:   parsedLocalIP,
			LocalPort:   uint(parsedLocalPort),
			ForeignAddr: parsedForeignIP,
			ForeignPort: uint(parsedForeignPort),
			State:       state,
		}

		parsedPID, err := strconv.Atoi(parts[pidIdx])
		if err != nil {
			datastore.MissingInfoNetworkConnection = append(datastore.MissingInfoNetworkConnection, &netObj)
			continue
		}

		parsedOwnerProcessName := parts[pidIdx+1]

		rawCreatedTime := strings.TrimSpace(strings.Join(parts[pidIdx+2:], " "))
		if len(rawCreatedTime) > 0 && rawCreatedTime != "N/A" && rawCreatedTime != "-" {
			layout := "2006-01-02 15:04:05.000000"
			parsedCreatedTime, err := time.Parse(layout, rawCreatedTime)
			if err != nil {
				utils.Logger.Warn("parse created time failed", zap.Int("pid", parsedPID), zap.String("time", rawCreatedTime), zap.String("plugin", colp.GetName()), zap.Error(err))
			} else {
				netObj.CreatedTime = parsedCreatedTime
			}
		}

		if _, ok := datastore.PIDToProcess[uint(parsedPID)]; !ok {
			datastore.PIDToProcess[uint(parsedPID)] = &datastore.Process{
				PID:       uint(parsedPID),
				ImageName: parsedOwnerProcessName,
			}
		}

		datastore.PIDToProcess[uint(parsedPID)].Conn = &netObj
		netObj.OwnerProcess = datastore.PIDToProcess[uint(parsedPID)]

		if !strings.EqualFold(datastore.PIDToProcess[uint(parsedPID)].ImageName, parsedOwnerProcessName) {
			utils.Logger.Warn("parsed process name mismatch", zap.Int("pid", parsedPID),
				zap.String("stored_proc_name", datastore.PIDToProcess[uint(parsedPID)].ImageName),
				zap.String("parsed_proc_name", parsedOwnerProcessName),
				zap.String("plugin", colp.GetName()), zap.Error(err))
		}
	}

	if err := scanner.Err(); err != nil {
		utils.Logger.Warn("Constructing process relations", zap.String("plugin", colp.GetName()), zap.Error(err))
	}
	return nil
}
