package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/envars"
	"github.com/ImDuong/vola-auto/plugins/help"
	"github.com/ImDuong/vola-auto/plugins/info"
	"github.com/ImDuong/vola-auto/plugins/process"
	"github.com/alitto/pond"
)

func runPlugins() {
	volPlgs := []plugins.VolPlugin{
		&help.HelpPlugin{},
		&info.InfoPlugin{},
		&process.ProcessPlugin{},
		&envars.EnvarsPlugin{},
	}

	volPlgRunningPool := pond.New(5, 20)
	for _, plg := range volPlgs {
		if !plugins.IsRunRequired(plg.GetArtifactsExtractionPath()) {
			fmt.Printf("Skipping plugin %s\n", plg.GetName())
			continue
		}
		fmt.Printf("Start running plugin %s\n", plg.GetName())
		volPlgRunningPool.Submit(func() {
			err := plg.Run()
			if err != nil {
				fmt.Printf("Running plugin %s got %s\n", plg.GetName(), err.Error())
				return
			}
		})
		fmt.Printf("Finish running plugin %s\n", plg.GetName())
	}
	volPlgRunningPool.StopAndWait()
}

func main() {
	// Define command line flags
	volatilityPath := flag.String("v", "", "Path to Volatility 3")
	memDumpPath := flag.String("f", "", "Path to memory dump file")
	outputFolderPath := flag.String("o", "", "Path to output folder")
	isForcedRerun := flag.Bool("r", false, "Force to re-run all plugins. Override old results")
	flag.Parse()

	// Check if required flags are provided
	if *volatilityPath == "" || *memDumpPath == "" || *outputFolderPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Create the output folder if it doesn't exist
	err := os.MkdirAll(*outputFolderPath, 0755)
	if err != nil {
		log.Fatalf("Error creating output folder: %v\n", err)
	}

	config.Default.VolRunConfig.Runner = "python"
	config.Default.VolRunConfig.Binary = filepath.Join(*volatilityPath, "vol.py")
	config.Default.OutputFolder = *outputFolderPath
	config.Default.MemoryDumpPath = *memDumpPath
	config.Default.IsForcedRerun = *isForcedRerun

	runPlugins()
}
