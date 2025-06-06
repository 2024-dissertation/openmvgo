package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/2024-dissertation/openmvgo/pkg/openmvg"
	"github.com/2024-dissertation/openmvgo/pkg/openmvs"
	"github.com/2024-dissertation/openmvgo/pkg/utils"
	"github.com/urfave/cli/v3"
)

func main() {
	var inputDir string
	var outputDir string
	var cameraDBFile string
	var maxThreads int

	cmd := &cli.Command{
		Name:  "OpenMVGO",
		Usage: "A CLI tool for OpenMVG and OpenMVS operations", // go run cmd/cli/main.go <input> <output>
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "maxThreads",
				Value:       1,
				Destination: &maxThreads,
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "input",
				Destination: &inputDir,
			},
			&cli.StringArg{
				Name:        "output",
				Destination: &outputDir,
			},
			&cli.StringArg{
				Name:        "cameraDB",
				Destination: &cameraDBFile,
			},
		},
		Action: func(context.Context, *cli.Command) error {
			if inputDir == "" || outputDir == "" {
				return cli.Exit("input and output directories must be specified", 1)
			}

			fmt.Printf("Input Directory: %s\n", inputDir)
			fmt.Printf("Output Directory: %s\n", outputDir)

			// Setup Utils
			utils := utils.NewUtils()

			timestamp := time.Now().Unix()

			// Middle directory creation
			buildDir, err := os.MkdirTemp("", fmt.Sprintf("%dbuild", timestamp))
			utils.Check(err)
			defer os.RemoveAll(buildDir)

			// Configure openmvg service

			openmvgService := openmvg.NewOpenMVGService(
				openmvg.NewOpenMVGConfig(
					inputDir,
					buildDir,
					&cameraDBFile,
				),
				utils,
			)

			// Configure openmvs service
			openmvsService := openmvs.NewOpenMVSService(
				openmvs.NewOpenMVSConfig(
					outputDir,
					buildDir,
					maxThreads,
				),
				utils,
			)

			// Populate and Run Pipelines
			openmvgService.PopulateTmpDir()
			defer os.Remove(*openmvgService.Config.CameraDBFile)
			defer os.RemoveAll(openmvgService.Config.MatchesDir)
			defer os.RemoveAll(openmvgService.Config.ReconstructionDir)

			openmvgService.SfMSequentialPipeline()
			openmvsService.RunPipeline()

			// Complete
			fmt.Println("OpenMVGO pipeline completed successfully!")

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
