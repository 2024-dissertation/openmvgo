package openmvs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/2024-dissertation/openmvgo/mocks"
	"github.com/2024-dissertation/openmvgo/pkg/openmvs"
	"go.uber.org/mock/gomock"
)

func TestRunDensifyPointCloud_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	expectedArgs := []string{
		"scene.mvs", "-o", "scene_dense.mvs",
		"-w", config.BuildDir,
		"--max-threads", fmt.Sprintf("%d", config.MaxThreads),
	}

	mockUtils.EXPECT().
		RunCommand("DensifyPointCloud", expectedArgs).
		Return(nil)

	service.RunDensifyPointCloud()
}

func TestRunDensifyPointCloud_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	expectedArgs := []string{
		"scene.mvs", "-o", "scene_dense.mvs",
		"-w", config.BuildDir,
		"--max-threads", fmt.Sprintf("%d", config.MaxThreads),
	}

	expectedErr := errors.New("command failed")

	mockUtils.EXPECT().
		RunCommand("DensifyPointCloud", expectedArgs).
		Return(expectedErr)

	mockUtils.EXPECT().
		Check(gomock.Any()).
		Do(func(err error) {
			if err == nil || err.Error() != fmt.Sprintf("failed to run DensifyPointCloud: %v", expectedErr) {
				t.Errorf("unexpected error passed to Check: %v", err)
			}
		})

	service.RunDensifyPointCloud()
}

func TestRunReconstructMesh_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	expectedArgs := []string{
		"scene_dense.mvs", "-o", "scene_mesh.ply",
		"-w", config.BuildDir,
	}

	mockUtils.EXPECT().
		RunCommand("ReconstructMesh", expectedArgs).
		Return(nil)

	service.RunReconstructMesh()
}

func TestRunReconstructMesh_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	expectedArgs := []string{
		"scene_dense.mvs", "-o", "scene_mesh.ply",
		"-w", config.BuildDir,
	}

	expectedErr := errors.New("command failed")

	mockUtils.EXPECT().
		RunCommand("ReconstructMesh", expectedArgs).
		Return(expectedErr)

	mockUtils.EXPECT().
		Check(gomock.Any()).
		Do(func(err error) {
			if err == nil || err.Error() != fmt.Sprintf("failed to run ReconstructMesh: %v", expectedErr) {
				t.Errorf("unexpected error passed to Check: %v", err)
			}
		})

	service.RunReconstructMesh()
}

func TestRunRefineMesh_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	expectedArgs := []string{
		"scene.mvs", "-m", "scene_mesh.ply", "-o", "scene_dense_mesh_refine.mvs",
		"-w", config.BuildDir,
		"--scales", "1", "--max-face-area", "16",
		"--max-threads", fmt.Sprintf("%d", config.MaxThreads),
	}

	mockUtils.EXPECT().
		RunCommand("RefineMesh", expectedArgs).
		Return(nil)

	service.RunRefineMesh()
}

func TestRunRefineMesh_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	expectedArgs := []string{
		"scene.mvs", "-m", "scene_mesh.ply", "-o", "scene_dense_mesh_refine.mvs",
		"-w", config.BuildDir,
		"--scales", "1", "--max-face-area", "16",
		"--max-threads", fmt.Sprintf("%d", config.MaxThreads),
	}

	expectedErr := errors.New("command failed")

	mockUtils.EXPECT().
		RunCommand("RefineMesh", expectedArgs).
		Return(expectedErr)

	mockUtils.EXPECT().
		Check(gomock.Any()).
		Do(func(err error) {
			if err == nil || err.Error() != fmt.Sprintf("failed to run RefineMesh: %v", expectedErr) {
				t.Errorf("unexpected error passed to Check: %v", err)
			}
		})

	service.RunRefineMesh()
}

func TestRunTextureMesh_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		OutputDir:  "/path/to/output",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	expectedArgs := []string{
		"scene_dense.mvs", "-m", "scene_dense_mesh_refine.ply",
		"-o", "scene_dense_mesh_refine_texture.mvs",
		"-w", config.BuildDir,
		"--export-type", "obj",
	}

	mockUtils.EXPECT().
		RunCommand("TextureMesh", expectedArgs).
		Return(nil)

	mockUtils.EXPECT().
		CopyFile(
			fmt.Sprintf("%s/scene_dense_mesh_refine_texture.mtl", config.BuildDir),
			fmt.Sprintf("%s/final.mtl", config.OutputDir),
		).
		Return(nil)

	mockUtils.EXPECT().
		CopyFile(
			fmt.Sprintf("%s/scene_dense_mesh_refine_texture.obj", config.BuildDir),
			fmt.Sprintf("%s/final.obj", config.OutputDir),
		).
		Return(nil)

	service.RunTextureMesh()
}

func TestRunTextureMesh_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		OutputDir:  "/path/to/output",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	expectedArgs := []string{
		"scene_dense.mvs", "-m", "scene_dense_mesh_refine.ply",
		"-o", "scene_dense_mesh_refine_texture.mvs",
		"-w", config.BuildDir,
		"--export-type", "obj",
	}

	expectedErr := errors.New("command failed")

	mockUtils.EXPECT().
		RunCommand("TextureMesh", expectedArgs).
		Return(expectedErr)

	mockUtils.EXPECT().
		CopyFile(
			fmt.Sprintf("%s/scene_dense_mesh_refine_texture.mtl", config.BuildDir),
			fmt.Sprintf("%s/final.mtl", config.OutputDir),
		).
		Return(nil)

	mockUtils.EXPECT().
		CopyFile(
			fmt.Sprintf("%s/scene_dense_mesh_refine_texture.obj", config.BuildDir),
			fmt.Sprintf("%s/final.obj", config.OutputDir),
		).
		Return(nil)

	mockUtils.EXPECT().
		Check(gomock.Any()).
		Do(func(err error) {
			if err == nil || err.Error() != fmt.Sprintf("failed to run TextureMesh: %v", expectedErr) {
				t.Errorf("unexpected error passed to Check: %v", err)
			}
		})

	service.RunTextureMesh()
}

func TestRunPipeline_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		OutputDir:  "/path/to/output",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	mockUtils.EXPECT().RunCommand("DensifyPointCloud", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("ReconstructMesh", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("RefineMesh", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("TextureMesh", gomock.Any()).Return(nil)

	mockUtils.EXPECT().CopyFile(gomock.Any(), gomock.Any()).Return(nil).Times(2)

	service.RunPipeline()
}

func TestRunPipeline_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := openmvs.OpenMVSConfig{
		BuildDir:   "/path/to/build",
		OutputDir:  "/path/to/output",
		MaxThreads: 4,
	}

	service := openmvs.OpenMVSServiceImpl{
		Utils:  mockUtils,
		Config: &config,
	}

	expectedErr := errors.New("command failed")

	mockUtils.EXPECT().RunCommand("DensifyPointCloud", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("ReconstructMesh", gomock.Any()).Return(nil)
	mockUtils.EXPECT().RunCommand("RefineMesh", gomock.Any()).Return(expectedErr)
	mockUtils.EXPECT().RunCommand("TextureMesh", gomock.Any()).Return(nil)

	mockUtils.EXPECT().CopyFile(gomock.Any(), gomock.Any()).Return(nil).Times(2)

	mockUtils.EXPECT().
		Check(gomock.Any()).
		Do(func(err error) {
			if err == nil || err.Error() != fmt.Sprintf("failed to run RefineMesh: %v", expectedErr) {
				t.Errorf("unexpected error passed to Check: %v", err)
			}
		})

	service.RunPipeline()
}

func TestNewOpenMVSService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := &openmvs.OpenMVSConfig{
		OutputDir:  "/path/to/output",
		BuildDir:   "/path/to/build",
		MaxThreads: 4,
	}

	mockUtils.EXPECT().EnsureDir(config.OutputDir).Return(nil)

	service := openmvs.NewOpenMVSService(config, mockUtils)

	if service.Utils != mockUtils {
		t.Errorf("expected Utils to be %v, got %v", mockUtils, service.Utils)
	}
	if service.Config != config {
		t.Errorf("expected Config to be %v, got %v", config, service.Config)
	}
}

func TestNewOpenMVSService_FailOutputDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := &openmvs.OpenMVSConfig{
		OutputDir:  "",
		BuildDir:   "/path/to/build",
		MaxThreads: 4,
	}

	mockUtils.EXPECT().Check(gomock.Any()).
		Do(func(err error) {
			panic(err)
		})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic due to empty OutputDir, but did not panic")
		} else {
			if err, ok := r.(error); ok {
				expected := "input directory must be specified"
				if err.Error() != expected {
					t.Errorf("expected panic error %q, got %q", expected, err.Error())
				}
			}
		}
	}()

	openmvs.NewOpenMVSService(config, mockUtils)
}

func TestNewOpenMVSService_FailEnsureDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtils := mocks.NewMockUtilsInterface(ctrl)

	config := &openmvs.OpenMVSConfig{
		OutputDir:  "/path/to/output",
		BuildDir:   "/path/to/build",
		MaxThreads: 4,
	}

	expectedErr := fmt.Errorf("failed to ensure input directory")

	mockUtils.EXPECT().EnsureDir(config.OutputDir).Return(expectedErr)

	mockUtils.EXPECT().Check(gomock.Any()).
		Do(func(err error) {
			panic(err)
		})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic due to EnsureDir failure, but did not panic")
		}
	}()

	openmvs.NewOpenMVSService(config, mockUtils)
}
