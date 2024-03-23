package plugins

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ImDuong/vola-auto/config"
)

type (
	VolPlugin interface {
		Run() error
		GetName() string
		GetArtifactsExtractionPath() string
	}

	CollectorPlugin interface {
		Run() error
		GetName() string
		GetArtifactsCollectionPath() string
	}

	AnalyticPlugin interface {
		Run() error
		GetName() string
		GetAnalyticResultPath() string
	}
)

func IsRunRequired(artifactExtractionFilepath string) bool {
	if config.Default.IsForcedRerun {
		return true
	}
	_, err := os.Stat(artifactExtractionFilepath)
	return os.IsNotExist(err)
}

func GetPermissionsToWriteResult(isOverride bool) int {
	perms := os.O_CREATE | os.O_WRONLY
	if !isOverride {
		perms = perms | os.O_APPEND
	} else {
		perms = perms | os.O_TRUNC
	}
	return perms
}

func RunVolatilityPluginAndWriteResult(args []string, resultFilepath string, isOverride bool) error {
	outputFileWriter, err := os.OpenFile(resultFilepath, GetPermissionsToWriteResult(isOverride), 0644)
	if err != nil {
		return err
	}
	defer outputFileWriter.Close()

	args = append([]string{config.Default.VolRunConfig.Binary, "-f", config.Default.MemoryDumpPath}, args...)
	cmd := exec.Command(config.Default.VolRunConfig.Runner, args...)
	cmd.Stdout = outputFileWriter
	cmd.Stderr = outputFileWriter

	fmt.Println("Executing", cmd.Args, "and writing to", resultFilepath)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
