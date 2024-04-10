package collectors

import (
	"bufio"
	"fmt"
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
	return "FILES COLLECTION PLUGIN"
}

func (colp *ProcessesPlugin) GetArtifactsCollectionPath() string {
	return ""
}

// Run() only processes & stores info about files in memory, not dump files
func (colp *ProcessesPlugin) Run() error {
	correspPlg := process.ProcessCmdlinePlugin{}
	cmdlineArtifactFiles, err := os.Open(correspPlg.GetArtifactsExtractionPath())
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
			return fmt.Errorf("parse pid %s failed", parts[0])
		}

		// TODO: handle process name have spaces
		proc := datastore.Process{
			ImageName: parts[1],
			PID:       uint(parsedPID),
		}

		if proc.PID != 4 {
			proc.Args = strings.Join(parts[2:], " ")
		}

		datastore.ProcessList = append(datastore.ProcessList, proc)
	}

	if err := scanner.Err(); err != nil {
		utils.Logger.Warn("Collecting artifacts", zap.String("plugin", colp.GetName()), zap.Error(err))
	}
	return nil
}
