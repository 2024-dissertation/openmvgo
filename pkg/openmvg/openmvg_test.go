package openmvg_test

import (
	"testing"

	"github.com/2024-dissertation/openmvgo/mocks"
	"github.com/2024-dissertation/openmvgo/pkg/openmvg"
	"go.uber.org/mock/gomock"
)

func TestNewOpenMVGService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	mockUtils.EXPECT().EnsureDir(gomock.Any()).Return(nil).AnyTimes()

	config := openmvg.OpenMVGConfig{
		InputDir:  "input",
		OutputDir: "output",
	}

	service := openmvg.NewOpenMVGService(
		config,
		mockUtils,
	)

	if service.Config.InputDir != config.InputDir {
		t.Errorf("Expected InputDir %s, got %s", config.InputDir, service.Config.InputDir)
	}

	if service.Config.OutputDir != config.OutputDir {
		t.Errorf("Expected OutputDir %s, got %s", config.OutputDir, service.Config.OutputDir)
	}
}

func TestRunSfMInitImageListing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	mockUtils.EXPECT().EnsureDir(gomock.Any()).Return(nil).AnyTimes()

	cameraDBFile := "camera_db.txt"
	config := openmvg.OpenMVGConfig{
		InputDir:     "input",
		OutputDir:    "output",
		CameraDBFile: &cameraDBFile,
	}

	service := openmvg.NewOpenMVGService(
		config,
		mockUtils,
	)

	expectedArgs := []string{
		"-i", config.InputDir,
		"-o", config.MatchesDir,
		"-d", *config.CameraDBFile,
		"-f", "2304",
	}

	mockUtils.EXPECT().
		RunCommand("openMVG_main_SfMInit_ImageListing", expectedArgs).
		Return(nil)

	service.RunSfMInitImageListing()
}

func TestRunSfMComputeFeatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	mockUtils.EXPECT().EnsureDir(gomock.Any()).Return(nil).AnyTimes()

	config := openmvg.OpenMVGConfig{
		InputDir:  "input",
		OutputDir: "output",
	}

	service := openmvg.NewOpenMVGService(
		config,
		mockUtils,
	)

	expectedArgs := []string{
		"-i", config.MatchesDir + "/sfm_data.json",
		"-o", config.MatchesDir,
		"-m", "SIFT",
	}

	mockUtils.EXPECT().
		RunCommand("openMVG_main_ComputeFeatures", expectedArgs).
		Return(nil)

	service.RunSfMComputeFeatures()
}

func TestRunSfMComputeMatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	mockUtils.EXPECT().EnsureDir(gomock.Any()).Return(nil).AnyTimes()

	config := openmvg.OpenMVGConfig{
		InputDir:  "input",
		OutputDir: "output",
	}

	service := openmvg.NewOpenMVGService(
		config,
		mockUtils,
	)

	expectedArgs := []string{
		"-i", config.MatchesDir + "/sfm_data.json",
		"-p", config.MatchesDir + "/pairs.bin",
		"-o", config.MatchesDir + "/matches.putative.bin",
	}

	mockUtils.EXPECT().
		RunCommand("openMVG_main_ComputeMatches", expectedArgs).
		Return(nil)

	service.RunSfMComputeMatches()
}

func TestRunSfMReconstruction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	mockUtils.EXPECT().EnsureDir(gomock.Any()).Return(nil).AnyTimes()

	config := openmvg.OpenMVGConfig{
		InputDir:  "input",
		OutputDir: "output",
	}

	service := openmvg.NewOpenMVGService(
		config,
		mockUtils,
	)

	expectedArgs := []string{
		"--sfm_engine", "INCREMENTAL",
		"--input_file", config.MatchesDir + "/sfm_data.json",
		"--match_dir", config.MatchesDir,
		"--output_dir", config.ReconstructionDir,
	}

	mockUtils.EXPECT().
		RunCommand("openMVG_main_SfM", expectedArgs).
		Return(nil)

	service.RunSfMReconstruction()
}

func TestRunSfMComputeSfMDataColor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	mockUtils.EXPECT().EnsureDir(gomock.Any()).Return(nil).AnyTimes()

	config := openmvg.OpenMVGConfig{
		InputDir:  "input",
		OutputDir: "output",
	}

	service := openmvg.NewOpenMVGService(
		config,
		mockUtils,
	)

	expectedArgs := []string{
		"-i", config.ReconstructionDir + "/sfm_data.bin",
		"-o", config.ReconstructionDir + "/colorized.ply",
	}

	mockUtils.EXPECT().
		RunCommand("openMVG_main_ComputeSfM_DataColor", expectedArgs).
		Return(nil)

	service.RunSfMComputeSfMDataColor()
}

func TestRunOpenMVG2OpenMVS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	mockUtils.EXPECT().EnsureDir(gomock.Any()).Return(nil).AnyTimes()

	config := openmvg.OpenMVGConfig{
		InputDir:  "input",
		OutputDir: "output",
	}

	service := openmvg.NewOpenMVGService(
		config,
		mockUtils,
	)

	expectedArgs := []string{
		"-i", config.ReconstructionDir + "/sfm_data.bin",
		"-o", config.OutputDir + "/scene.mvs",
		"-d", config.OutputDir,
	}

	mockUtils.EXPECT().
		RunCommand("openMVG_main_openMVG2openMVS", expectedArgs).
		Return(nil)

	service.RunOpenMVG2OpenMVS()
}

func TestRunHealthCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	mockUtils.EXPECT().EnsureDir(gomock.Any()).Return(nil).AnyTimes()

	config := openmvg.OpenMVGConfig{
		InputDir:  "input",
		OutputDir: "output",
	}

	service := openmvg.NewOpenMVGService(
		config,
		mockUtils,
	)

	mockUtils.EXPECT().
		RunCommand("Tests", []string{}).
		Return(nil)

	service.RunHealthCheck()
}

func TestSfMSequentialPipeline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	mockUtils.EXPECT().EnsureDir(gomock.Any()).Return(nil).AnyTimes()

	cameraDBFile := "camera_db.txt"
	config := openmvg.OpenMVGConfig{
		InputDir:     "input",
		OutputDir:    "output",
		CameraDBFile: &cameraDBFile,
	}

	service := openmvg.NewOpenMVGService(
		config,
		mockUtils,
	)

	mockUtils.EXPECT().RunCommand("openMVG_main_SfMInit_ImageListing", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("openMVG_main_ComputeFeatures", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("openMVG_main_PairGenerator", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("openMVG_main_ComputeMatches", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("openMVG_main_GeometricFilter", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("openMVG_main_SfM", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("openMVG_main_ComputeSfM_DataColor", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("openMVG_main_openMVG2openMVS", gomock.Any()).Return(nil)

	service.SfMSequentialPipeline()
}
