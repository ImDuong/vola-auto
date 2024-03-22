package info

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
)

type (
	InfoPlugin struct {
		name string
	}
)

func (ip *InfoPlugin) GetName() string {
	return "INFO PLUGIN"
}

func (ip *InfoPlugin) Run() error {
	command := config.Default.VolRunConfig.Runner
	args := []string{config.Default.VolRunConfig.Binary, "-f", config.Default.MemoryDumpPath, "windows.info.Info"}
	cmd := exec.Command(command, args...)

	outputFile := filepath.Join(config.Default.OutputFolder, "info.txt")
	outputFileWriter, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outputFileWriter.Close()

	cmd.Stdout = outputFileWriter
	cmd.Stderr = outputFileWriter

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
