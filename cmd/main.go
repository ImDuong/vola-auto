package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/info"
)

func main() {
	// Define command line flags
	volatilityPath := flag.String("v", "", "Path to Volatility 3")
	memDumpPath := flag.String("f", "", "Path to memory dump file")
	outputFolderPath := flag.String("o", "", "Path to output folder")
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

	plg := &info.InfoPlugin{}
	plg.Run()
}
