package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/runner"
	"github.com/ImDuong/vola-auto/utils"
)

func main() {
	volatilityPath := flag.String("v", "", "Path to Volatility 3")
	memDumpPath := flag.String("f", "", "Path to memory dump file")
	outputFolderPath := flag.String("o", "", "Path to output folder")
	isForcedRerun := flag.Bool("r", false, "Force to re-run all plugins. Override old results")
	flag.Parse()

	if *volatilityPath == "" || *memDumpPath == "" || *outputFolderPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	pythonRunner, pythonVersion, err := utils.GetPythonRunner()
	if err != nil {
		log.Fatalf("Error when getting python version")
	}
	if len(pythonVersion) > 2 && pythonVersion[0] == '2' {
		log.Fatalf("volatility 2 is not supported yet")
	}
	config.Default.VolRunConfig.Runner = pythonRunner

	config.Default.VolRunConfig.Binary = filepath.Join(*volatilityPath, "vol.py")
	config.Default.OutputFolder = *outputFolderPath
	config.Default.DumpFilesFolder = filepath.Join(config.Default.OutputFolder, "dump_files")
	config.Default.AnalyticFolder = filepath.Join(config.Default.OutputFolder, "analytics")
	config.Default.MemoryDumpPath = *memDumpPath
	config.Default.IsForcedRerun = *isForcedRerun

	err = os.MkdirAll(config.Default.OutputFolder, 0755)
	if err != nil {
		log.Fatalf("Error creating output folder: %v\n", err)
	}

	err = os.MkdirAll(config.Default.DumpFilesFolder, 0755)
	if err != nil {
		log.Fatalf("Error creating dump files folder: %v\n", err)
	}

	err = os.MkdirAll(config.Default.AnalyticFolder, 0755)
	if err != nil {
		log.Fatalf("Error creating analytic folder: %v\n", err)
	}

	err = runner.RunPlugins()
	if err != nil {
		log.Fatal(err)
	}
}
