package collectors

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/ImDuong/vola-auto/datastore"
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
		}

		datastore.PIDToProcess[proc.PID] = &proc
	}

	if err := scanner.Err(); err != nil {
		utils.Logger.Warn("Collecting processes", zap.String("plugin", colp.GetName()), zap.Error(err))
	}

	// read pslist to construct parent-child relation
	pslistPlg := process.ProcessPsListPlugin{}
	pslistArtifactFiles, err := os.Open(pslistPlg.GetArtifactsExtractionPath())
	if err != nil {
		return err
	}
	defer pslistArtifactFiles.Close()
	scanner = bufio.NewScanner(pslistArtifactFiles)
	isProcessDataFound = false

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
			datastore.PIDToProcess[uint(parsedPPID)] = &datastore.Process{PID: uint(parsedPPID)}
		}

		datastore.PIDToProcess[uint(parsedPID)].ParentProc = datastore.PIDToProcess[uint(parsedPPID)]
	}

	if err := scanner.Err(); err != nil {
		utils.Logger.Warn("Constructing process relations", zap.String("plugin", colp.GetName()), zap.Error(err))
	}

	return nil
}
