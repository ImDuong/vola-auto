package config

type (
	VolatilityRunnerConfig struct {
		Runner string
		Binary string
	}

	Configuration struct {
		VolRunConfig   VolatilityRunnerConfig
		MemoryDumpPath string
		OutputFolder   string
		AnalyticFolder string
		IsForcedRerun  bool
	}
)

var Default Configuration
