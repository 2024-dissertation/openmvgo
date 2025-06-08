package openmvg

import (
	"fmt"
	"os"
	"time"

	"github.com/2024-dissertation/openmvgo/pkg/mvgoutils"
)

// Config for running the OpenMVG pipeline
type OpenMVGConfig struct {
	InputDir          string
	OutputDir         string
	MatchesDir        string
	ReconstructionDir string
	CameraDBFile      *string
}

// Create an OpenMVG config. Handles creating tempory directories for mvs and reconstruction
func NewOpenMVGConfig(inputDir string, outputDir string, cameraDBFile *string) OpenMVGConfig {
	return OpenMVGConfig{
		InputDir:     inputDir,
		OutputDir:    outputDir,
		CameraDBFile: cameraDBFile,
	}
}

type AppFileServiceImpl struct {
	Utils  mvgoutils.OpenmvgoUtilsInterface
	Config OpenMVGConfig
}

func NewOpenMVGService(config OpenMVGConfig, utils mvgoutils.OpenmvgoUtilsInterface) AppFileServiceImpl {
	if utils == nil {
		utils = mvgoutils.NewMvgoUtils()
	}

	// Pre run folder checks
	if config.InputDir == "" || config.OutputDir == "" {
		utils.Check(fmt.Errorf("input and output directories must be specified"))
	}

	if err := utils.EnsureDir(config.InputDir); err != nil {
		utils.Check(fmt.Errorf("failed to ensure input directory: %w", err))
	}

	if err := utils.EnsureDir(config.OutputDir); err != nil {
		utils.Check(fmt.Errorf("failed to ensure output directory: %w", err))
	}

	return AppFileServiceImpl{
		Utils:  utils,
		Config: config,
	}
}

func (s *AppFileServiceImpl) PopulateTmpDir() {
	// Ensure the camera database file is set, if not download it
	if s.Config.CameraDBFile == nil || *s.Config.CameraDBFile == "" {
		f, err := s.Utils.DownloadFile("https://raw.githubusercontent.com/openMVG/openMVG/refs/heads/develop/src/openMVG/exif/sensor_width_database/sensor_width_camera_database.txt")
		s.Utils.Check(err)
		s.Config.CameraDBFile = &f
	}

	timestamp := time.Now().Unix()

	// Matches dir
	matchesDir, err := os.MkdirTemp("", fmt.Sprintf("%dmatches", timestamp))
	s.Utils.Check(err)

	s.Config.MatchesDir = matchesDir

	// Reconstruction dir
	reconstructionDir, err := os.MkdirTemp("", fmt.Sprintf("%dreconstruction", timestamp))
	s.Utils.Check(err)

	s.Config.ReconstructionDir = reconstructionDir

}

func (s *AppFileServiceImpl) SfMSequentialPipeline() {
	s.RunSfMInitImageListing()
	s.RunSfMComputeFeatures()
	s.RunSfMPairGenerator()
	s.RunSfMComputeMatches()
	s.RunSfMGeometricFilter()
	s.RunSfMReconstruction()
	s.RunSfMComputeSfMDataColor()
	s.RunOpenMVG2OpenMVS()
}

func (s *AppFileServiceImpl) RunHealthCheck() {
	s.Utils.RunCommand("Tests", []string{})
}

func (s *AppFileServiceImpl) RunSfMInitImageListing() {
	args := []string{
		"-i", s.Config.InputDir,
		"-o", s.Config.MatchesDir,
		"-d", *s.Config.CameraDBFile,
		"-f", "2304", // 2304 is the focal length, adjust as necessary
	}

	s.Utils.RunCommand("openMVG_main_SfMInit_ImageListing", args)
}

func (s *AppFileServiceImpl) RunSfMComputeFeatures() {
	args := []string{
		"-i", s.Config.MatchesDir + "/sfm_data.json",
		"-o", s.Config.MatchesDir,
		"-m", "SIFT",
	}

	s.Utils.RunCommand("openMVG_main_ComputeFeatures", args)
}

func (s *AppFileServiceImpl) RunSfMPairGenerator() {
	args := []string{
		"-i", s.Config.MatchesDir + "/sfm_data.json",
		"-o", s.Config.MatchesDir + "/pairs.bin",
	}

	s.Utils.RunCommand("openMVG_main_PairGenerator", args)
}

func (s *AppFileServiceImpl) RunSfMComputeMatches() {
	args := []string{
		"-i", s.Config.MatchesDir + "/sfm_data.json",
		"-p", s.Config.MatchesDir + "/pairs.bin",
		"-o", s.Config.MatchesDir + "/matches.putative.bin",
	}

	s.Utils.RunCommand("openMVG_main_ComputeMatches", args)
}

func (s *AppFileServiceImpl) RunSfMGeometricFilter() {
	args := []string{
		"-i", s.Config.MatchesDir + "/sfm_data.json",
		"-m", s.Config.MatchesDir + "/matches.putative.bin",
		"-g", "f",
		"-o", s.Config.MatchesDir + "/matches.f.bin",
	}

	s.Utils.RunCommand("openMVG_main_GeometricFilter", args)
}

func (s *AppFileServiceImpl) RunSfMReconstruction() {
	args := []string{
		"--sfm_engine", "INCREMENTAL",
		"--input_file", s.Config.MatchesDir + "/sfm_data.json",
		"--match_dir", s.Config.MatchesDir,
		"--output_dir", s.Config.ReconstructionDir,
	}

	s.Utils.RunCommand("openMVG_main_SfM", args)
}

func (s *AppFileServiceImpl) RunSfMComputeSfMDataColor() {
	args := []string{
		"-i", s.Config.ReconstructionDir + "/sfm_data.bin",
		"-o", s.Config.ReconstructionDir + "/colorized.ply",
	}

	s.Utils.RunCommand("openMVG_main_ComputeSfM_DataColor", args)
}

func (s *AppFileServiceImpl) RunOpenMVG2OpenMVS() {
	args := []string{
		"-i", s.Config.ReconstructionDir + "/sfm_data.bin",
		"-o", s.Config.OutputDir + "/scene.mvs",
		"-d", s.Config.OutputDir,
	}

	s.Utils.RunCommand("openMVG_main_openMVG2openMVS", args)
}
