package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/2024-dissertation/openmvgo/pkg/mvgoutils"
	"github.com/2024-dissertation/openmvgo/pkg/openmvg"
	"github.com/2024-dissertation/openmvgo/pkg/openmvs"
	"github.com/2024-dissertation/openmvgo/pkg/services"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Couldn't load .env file: %v", err)
	}

	bucketPath := os.Getenv("BUCKET_PATH")

	storageService := services.NewKatapultStorageService()

	inputDir, err := storageService.DownloadFolder(bucketPath)
	if err != nil {
		log.Fatalf("Failed to download folder: %v", err)
	}

	outputDir := os.Getenv("OUTPUT_PATH")

	fmt.Printf("Input Directory: %s\n", inputDir)
	fmt.Printf("Output Directory: %s\n", outputDir)

	utils := mvgoutils.NewMvgoUtils()

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
			nil,
		),
		utils,
	)

	// Configure openmvs service
	openmvsService := openmvs.NewOpenMVSService(
		openmvs.NewOpenMVSConfig(
			outputDir,
			buildDir,
			0,
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

	fmt.Println("OpenMVGO pipeline completed successfully!")
}
