package config

type (
	VolatilityRunnerConfig struct {
		Runner string // python3
		Binary string // filepath of vol.py
	}

	Configuration struct {
		VolRunConfig    VolatilityRunnerConfig
		MemoryDumpPath  string
		OutputFolder    string
		AnalyticFolder  string
		DumpFilesFolder string
		BatchCmdFolder  string
		IsForcedRerun   bool
	}
)

var Default Configuration

const (
	DefaultArtifactFolderName = "artifacts"
	AnalyticsFoldername       = "0_analytics"
	DumpFilesFoldername       = "1_dump_files"
	BatchCmdResultFilename    = "2_batch_cmd_results"
)
