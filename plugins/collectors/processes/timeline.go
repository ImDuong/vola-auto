package processes

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/datastore"
	"github.com/ImDuong/vola-auto/plugins/volatility/process"
	"github.com/ImDuong/vola-auto/utils"
	"go.uber.org/zap"
)

type (
	TimelinePlugin struct {
	}
)

func (colp *TimelinePlugin) GetName() string {
	return "PROCESS TIMELINE COLLECTION PLUGIN"
}

func (colp *TimelinePlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(colp.GetArtifactsCollectionFolderpath(), "timeline.txt")
}

func (colp *TimelinePlugin) GetArtifactsCollectionFolderpath() string {
	return filepath.Join(config.Default.OutputFolder, ProcessCollectionFolderName)
}

// Run() collects created time from pslist plugin, and writes out under csv format
func (colp *TimelinePlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionFolderpath(), 0755)
	if err != nil {
		return err
	}

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

		// CreatedTime is contained in 9th and 10th column
		if len(parts) < 10 {
			continue
		}

		// TODO: handle process name have spaces
		createdTime := parts[8] + " " + parts[9]

		parsedPID, err := strconv.Atoi(parts[0])
		if err != nil {
			utils.Logger.Warn("parse pid failed", zap.String("pid", parts[0]), zap.String("plugin", colp.GetName()), zap.Error(err))
			continue
		}

		if _, ok := datastore.PIDToProcess[uint(parsedPID)]; !ok {
			utils.Logger.Warn("pid not found in current datastore", zap.Int("pid", parsedPID), zap.String("plugin", colp.GetName()))
			continue
		}

		layout := "2006-01-02 15:04:05.000000"
		parsedTime, err := time.Parse(layout, createdTime)
		if err != nil {
			utils.Logger.Warn("parse created time failed", zap.Int("pid", parsedPID), zap.String("time", createdTime), zap.String("plugin", colp.GetName()), zap.Error(err))
		} else {
			datastore.PIDToProcess[uint(parsedPID)].CreatedTime = parsedTime
		}
	}

	if err := scanner.Err(); err != nil {
		utils.Logger.Warn("Constructing process relations", zap.String("plugin", colp.GetName()), zap.Error(err))
	}

	timelineFileWriter, err := os.OpenFile(colp.GetArtifactsCollectionPath(), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer timelineFileWriter.Close()

	writer := csv.NewWriter(timelineFileWriter)
	defer writer.Flush()
	writer.Comma = ' '

	if err := writer.Write([]string{"CreatedTime", "PID", "ImageFileName", "Args"}); err != nil {
		return fmt.Errorf("write header to CSV failed: %w", err)
	}

	tempProcList := make([]*datastore.Process, len(datastore.PIDToProcess))
	idx := 0
	for i := range datastore.PIDToProcess {
		tempProcList[idx] = datastore.PIDToProcess[i]
		idx++
	}

	sort.Slice(tempProcList, func(i, j int) bool {
		return tempProcList[i].CreatedTime.After(tempProcList[j].CreatedTime)
	})

	for _, p := range tempProcList {
		record := []string{
			fmt.Sprintf("%-29s", p.CreatedTime.String()),
			fmt.Sprintf("%-7d", p.PID),
			fmt.Sprintf("%-25s", p.ImageName),
			p.Args,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
