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

func RunVolatilityPluginAndWriteResult(args []string, resultFilepath string) error {
	cmd := exec.Command(config.Default.VolRunConfig.Runner, args...)

	perms := os.O_CREATE | os.O_WRONLY
	if !config.Default.IsForcedRerun {
		perms = perms | os.O_APPEND
	} else {
		perms = perms | os.O_TRUNC
	}
	outputFileWriter, err := os.OpenFile(resultFilepath, perms, 0644)
	if err != nil {
		return err
	}
	defer outputFileWriter.Close()

	cmd.Stdout = outputFileWriter
	cmd.Stderr = outputFileWriter

	fmt.Println("Executing", cmd.Args, "and writing to", resultFilepath)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
