package openmvs

import (
	"fmt"

	"github.com/2024-dissertation/openmvgo/pkg/mvgoutils"
)

// Config object
type OpenMVSConfig struct {
	MaxThreads int
	OutputDir  string
	BuildDir   string
}

// Helper function to create an OpenMVSConfig
func NewOpenMVSConfig(outputDir string, buildDir string, maxThreads int) *OpenMVSConfig {
	return &OpenMVSConfig{
		MaxThreads: maxThreads,
		OutputDir:  outputDir,
		BuildDir:   buildDir,
	}
}

// OpenMVSService interface implementation
type OpenMVSServiceImpl struct {
	Utils  mvgoutils.OpenmvgoUtilsInterface
	Config *OpenMVSConfig
}

// Helper function to create a new OpenMVSServiceImpl
func NewOpenMVSService(config *OpenMVSConfig, utils mvgoutils.OpenmvgoUtilsInterface) OpenMVSServiceImpl {
	if utils == nil {
		utils = mvgoutils.NewMvgoUtils()
	}

	if config.OutputDir == "" {
		utils.Check(fmt.Errorf("input directory must be specified"))
	}

	if err := utils.EnsureDir(config.OutputDir); err != nil {
		utils.Check(fmt.Errorf("failed to ensure input directory"))
	}

	return OpenMVSServiceImpl{
		Utils:  utils,
		Config: config,
	}
}

// RunPipeline runs the entire OpenMVS pipeline in sequence
func (s OpenMVSServiceImpl) RunPipeline() {
	s.RunDensifyPointCloud()
	s.RunReconstructMesh()
	s.RunRefineMesh()
	s.RunTextureMesh()
}

// RunDensifyPointCloud runs the DensifyPointCloud command with the configured parameters
func (s OpenMVSServiceImpl) RunDensifyPointCloud() {
	if err := s.Utils.RunCommand("DensifyPointCloud", []string{"scene.mvs", "-o", "scene_dense.mvs", "-w", s.Config.BuildDir, "--max-threads", fmt.Sprintf("%d", s.Config.MaxThreads)}); err != nil {
		s.Utils.Check(fmt.Errorf("failed to run DensifyPointCloud: %w", err))
	}
}

// RunReconstructMesh runs the ReconstructMesh command with the configured parameters
func (s OpenMVSServiceImpl) RunReconstructMesh() {
	if err := s.Utils.RunCommand("ReconstructMesh", []string{"scene_dense.mvs", "-o", "scene_mesh.ply", "-w", s.Config.BuildDir}); err != nil {
		s.Utils.Check(fmt.Errorf("failed to run ReconstructMesh: %w", err))
	}
}

// RunRefineMesh runs the RefineMesh command with the configured parameters
func (s OpenMVSServiceImpl) RunRefineMesh() {
	if err := s.Utils.RunCommand("RefineMesh", []string{"scene.mvs", "-m", "scene_mesh.ply", "-o", "scene_dense_mesh_refine.mvs", "-w", s.Config.BuildDir, "--scales", "1", "--max-face-area", "16", "--max-threads", fmt.Sprintf("%d", s.Config.MaxThreads)}); err != nil {
		s.Utils.Check(fmt.Errorf("failed to run RefineMesh: %w", err))
	}
}

// RunTextureMesh runs the TextureMesh command with the configured parameters
func (s OpenMVSServiceImpl) RunTextureMesh() {
	if err := s.Utils.RunCommand("TextureMesh", []string{"scene_dense.mvs", "-m", "scene_dense_mesh_refine.ply", "-o", "scene_dense_mesh_refine_texture.mvs", "-w", s.Config.BuildDir, "--export-type", "obj"}); err != nil {
		s.Utils.Check(fmt.Errorf("failed to run TextureMesh: %w", err))
	}
	s.Utils.CopyFile(
		fmt.Sprintf("%s/scene_dense_mesh_refine_texture.mtl", s.Config.BuildDir),
		fmt.Sprintf("%s/final.mtl", s.Config.OutputDir),
	)

	s.Utils.CopyFile(
		fmt.Sprintf("%s/scene_dense_mesh_refine_texture.obj", s.Config.BuildDir),
		fmt.Sprintf("%s/final.obj", s.Config.OutputDir),
	)
}
