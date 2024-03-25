package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/runner"
	"github.com/ImDuong/vola-auto/utils"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "Vola Auto",
		Usage: "Auto streamline for Volatility 3",
		Commands: []*cli.Command{
			{
				Name:    "dumpfiles",
				Aliases: []string{"d"},
				Usage:   "dump files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Aliases: []string{"r"},
						Name:    "regex",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// TODO: support regex
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "vol",
				Aliases:  []string{"v"},
				Usage:    "Path to Volatility 3",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Path to memory dump file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Path to output folder",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "rerun",
				Aliases: []string{"r"},
				Usage:   "Force to re-run all plugins. Override old results",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			pythonRunner, pythonVersion, err := utils.GetPythonRunner()
			if err != nil {
				return fmt.Errorf("error when getting python version")
			}
			if len(pythonVersion) > 2 && pythonVersion[0] == '2' {
				return fmt.Errorf("volatility 2 is not supported yet")
			}
			config.Default.VolRunConfig.Runner = pythonRunner

			config.Default.VolRunConfig.Binary = filepath.Join(cmd.String("v"), "vol.py")
			config.Default.OutputFolder = cmd.String("o")
			config.Default.DumpFilesFolder = filepath.Join(config.Default.OutputFolder, "dump_files")
			config.Default.AnalyticFolder = filepath.Join(config.Default.OutputFolder, "analytics")
			config.Default.MemoryDumpPath = cmd.String("f")
			config.Default.IsForcedRerun = cmd.Bool("r")

			err = os.MkdirAll(config.Default.OutputFolder, 0755)
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

			err = runner.RunPlugins()
			if err != nil {
				return err
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
