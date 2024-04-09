package utils

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger       *zap.Logger
	loggingLevel = zapcore.InfoLevel
)

func init() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	if strings.TrimSpace(os.Getenv("DEBUG")) == "1" {
		loggingLevel = zapcore.DebugLevel
	}
	Logger = zap.New(zapcore.NewCore(consoleEncoder, os.Stdout, loggingLevel))
	if loggingLevel == zapcore.DebugLevel {
		Logger.Debug("DEBUG MODE IS ON. BE PREPARED...")
	}
}
