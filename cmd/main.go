package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/ImDuong/vola-auto/plugins/volatility/filescan"
	"github.com/ImDuong/vola-auto/runner"
	"github.com/ImDuong/vola-auto/utils"
	"github.com/alitto/pond"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

func main() {
	cmd := &cli.Command{
		Name:  "Vola Auto",
		Usage: "Auto streamline for Volatility 3",
		Commands: []*cli.Command{
			{
				Name:    "dumpfiles",
				Aliases: []string{"d"},
				Usage:   "Dump files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Aliases: []string{"reg"},
						Name:    "regex",
					},
					&cli.StringFlag{
						Aliases: []string{"fs"},
						Name:    "filescan",
						Usage:   "Path to filescan plugin's output. If this flag is empty, auto find filescan.txt. If no file exists, auto run filescan plugin. However, if flag is not empty and file does not exist, program will exit",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var err error
					filescanResultPath := cmd.String("filescan")

					// backup output folder path to store new output folder for running filescan plugin
					backupOutputPath := config.Default.OutputFolder

					filescanPlg := filescan.FilescanPlugin{}
					fileCollectorPlg := collectors.FilesPlugin{}
					if len(filescanResultPath) == 0 {
						filescanResultPath = filescanPlg.GetArtifactsExtractionPath()

						// run filescan plugin and store output to filescan.txt if filescan.txt does not exist
						_, err = os.Stat(filescanResultPath)
						if err != nil && !os.IsNotExist(err) {
							return err
						}
						if os.IsNotExist(err) {
							err = filescanPlg.Run()
							if err != nil {
								return err
							}
						}
					} else {
						_, err = os.Stat(filescanResultPath)
						if err != nil {
							return err
						}

						// set new value for output folder, because filescan.FilescanPlugin use this new path to read filescan result path
						config.Default.OutputFolder = filepath.Dir(filescanResultPath)
						filescanPlg.SetArtifactsExtractionFilename(filepath.Base(filescanResultPath))
					}

					// construct file lists
					fileCollectorPlg.Run()

					// restore original output folder path for dumping files
					config.Default.OutputFolder = backupOutputPath

					foundFiles, err := fileCollectorPlg.FindFilesByRegex(cmd.String("regex"))
					if err != nil {
						return err
					}

					dumpFilesPool := pond.New(20, 100)
					var aggregatedError error
					var aggregateErrorMutex sync.Mutex
					for i := range foundFiles {
						copiedIdx := i
						dumpFilesPool.Submit(func() {
							err := fileCollectorPlg.DumpFile(foundFiles[copiedIdx], config.Default.OutputFolder)
							if err != nil {
								aggregateErrorMutex.Lock()
								aggregatedError = fmt.Errorf("%w;%w", aggregatedError, err)
								aggregateErrorMutex.Unlock()
							}
						})
					}
					dumpFilesPool.StopAndWait()
					if aggregatedError != nil {
						return aggregatedError
					}
					return nil
				},
			},
			{
				Name:    "batch",
				Aliases: []string{"b"},
				Usage:   "Run multiples commands parallely",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Aliases: []string{"f"},
						Name:    "file",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					commandFile, err := os.Open(cmd.String("file"))
					if err != nil {
						return err
					}
					defer commandFile.Close()

					commandPool := pond.New(20, 100)
					scanner := bufio.NewScanner(commandFile)
					for scanner.Scan() {
						line := strings.TrimSpace(scanner.Text())
						commandArgs := strings.Split(line, " ")
						args := []string{config.Default.VolRunConfig.Binary,
							"-f", config.Default.MemoryDumpPath,
						}
						args = append(args, commandArgs...)
						commandPool.Submit(func() {
							err = plugins.RunVolatilityPluginAndWriteResult(args, "", true)
							if err != nil {
								utils.Logger.Error("Running", zap.String("cmd", line), zap.Error(err))
								return
							}
						})
					}
					commandPool.StopAndWait()
					return scanner.Err()
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "vol",
				Aliases:  []string{"v"},
				Usage:    "Path to Volatility 3",
				Required: true,
				Action: func(ctx context.Context, c *cli.Command, s string) error {
					config.Default.VolRunConfig.Binary = filepath.Join(s, "vol.py")
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Path to memory dump file",
				Required: true,
				Action: func(ctx context.Context, c *cli.Command, s string) error {
					config.Default.MemoryDumpPath = s
					config.Default.OutputFolder = c.String("output")

					if len(config.Default.OutputFolder) == 0 {
						config.Default.OutputFolder = filepath.Join(filepath.Dir(config.Default.MemoryDumpPath), "artifacts")
					}
					config.Default.AnalyticFolder = filepath.Join(config.Default.OutputFolder, "0_analytics")
					config.Default.DumpFilesFolder = filepath.Join(config.Default.OutputFolder, "1_dump_files")
					err := os.MkdirAll(config.Default.OutputFolder, 0755)
					if err != nil {
						return fmt.Errorf("error creating output folder: %w", err)
					}

					err = os.MkdirAll(config.Default.DumpFilesFolder, 0755)
					if err != nil {
						return fmt.Errorf("error creating dump files folder: %w", err)
					}

					err = os.MkdirAll(config.Default.AnalyticFolder, 0755)
					if err != nil {
						return fmt.Errorf("error creating analytic folder: %w", err)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Path to output folder. If empty, output folder is set to artifacts folder in path having memory dump file",
				Action: func(ctx context.Context, c *cli.Command, s string) error {

					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "rerun",
				Aliases: []string{"r"},
				Usage:   "Force to re-run all plugins. Override old results",
				Action: func(ctx context.Context, c *cli.Command, b bool) error {
					config.Default.IsForcedRerun = b
					return nil
				},
			},
		},
		Before: func(ctx context.Context, c *cli.Command) error {
			pythonRunner, pythonVersion, err := utils.GetPythonRunner()
			if err != nil {
				return fmt.Errorf("error when getting python version")
			}
			if len(pythonVersion) > 2 && pythonVersion[0] == '2' {
				return fmt.Errorf("volatility 2 is not supported yet")
			}
			config.Default.VolRunConfig.Runner = pythonRunner
			return nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			err := runner.RunPlugins()
			if err != nil {
				return err
			}
			return nil
		},
		After: func(ctx context.Context, c *cli.Command) error {
			utils.Logger.Sync()
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
