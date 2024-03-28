package envars

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/volatility/envars"
)

type (
	EnvarsPlugin struct {
	}

	processEnvVar struct {
		PID      string
		Process  string
		Block    string
		Variable string
		Value    string
	}
)

const (
	PluginName         = "ENVARS ANALYTIC PLUGIN"
	AnalyticResultPath = "envars.txt"
)

func (anp *EnvarsPlugin) GetName() string {
	return PluginName
}

func (anp *EnvarsPlugin) GetAnalyticResultPath() string {
	return filepath.Join(config.Default.AnalyticFolder, AnalyticResultPath)
}

func (anp *EnvarsPlugin) IsWhitelistedVariable(checkingVariable string) bool {
	whitelist := []string{
		"ALLUSERSPROFILE",
		"APPDATA",
		"CommonProgramFiles",
		"CommonProgramFiles(x86)",
		"CommonProgramW6432",
		"COMPUTERNAME",
		"ComSpec",
		"HOMEDRIVE",
		"HOMEPATH",
		"LOCALAPPDATA",
		"LOGONSERVER",
		"NUMBER_OF_PROCESSORS",
		"OS",
		"Path",
		"PATHEXT",
		"PROCESSOR_ARCHITECTURE",
		"PROCESSOR_IDENTIFIER",
		"PROCESSOR_LEVEL",
		"PROCESSOR_REVISION",
		"ProgramData",
		"ProgramFiles",
		"ProgramFiles(x86)",
		"ProgramW6432",
		"PUBLIC",
		"SystemDrive",
		"SystemRoot",
		"TEMP",
		"TMP",
		"USERDOMAIN",
		"USERNAME",
		"USERPROFILE",
		"windir",
	}

	for _, v := range whitelist {
		if checkingVariable == v {
			return true
		}
	}
	return false
}

func (anp *EnvarsPlugin) Run() error {
	correspPlg := envars.EnvarsPlugin{}
	file, err := os.Open(correspPlg.GetArtifactsExtractionPath())
	if err != nil {
		return err
	}
	defer file.Close()

	var suspiciousProcesses []processEnvVar

	scanner := bufio.NewScanner(file)
	isProcessDataFound := false

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}
		if !isProcessDataFound {
			if strings.Contains(line, "PID") && strings.Contains(line, "Process") && strings.Contains(line, "Block") && strings.Contains(line, "Variable") && strings.Contains(line, "Value") {
				isProcessDataFound = true
			}
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		if anp.IsWhitelistedVariable(parts[3]) {
			continue
		}

		var variableValue string
		if len(parts) == 5 {
			variableValue = parts[4]
		}

		suspiciousProcesses = append(suspiciousProcesses, processEnvVar{
			PID:      parts[0],
			Process:  parts[1],
			Block:    parts[2],
			Variable: parts[3],
			Value:    variableValue,
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(PluginName, ":got some errors when analyzing")
	}

	resultFile, err := os.OpenFile(anp.GetAnalyticResultPath(), plugins.GetFileOpenFlag(config.Default.IsForcedRerun), 0644)
	if err != nil {
		return err
	}
	defer resultFile.Close()

	writer := csv.NewWriter(resultFile)
	defer writer.Flush()

	if err := writer.Write([]string{"PID", "Process", "Block", "Variable", "Value"}); err != nil {
		log.Fatalf("error writing header to CSV: %s", err)
	}

	for _, p := range suspiciousProcesses {
		record := []string{p.PID, p.Process, p.Block, p.Variable, p.Value}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
