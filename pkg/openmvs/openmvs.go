package openmvs

// OpenMVSServiceInterface defines the methods for running OpenMVS commands in sequence.
//
//go:generate mockgen -source=./openmvs.go -destination=../../mocks/mock_openmvs.go -package=mocks
type OpenMVSServiceInterface interface {
	RunPipeline()
	RunDensifyPointCloud()
	RunReconstructMesh()
	RunRefineMesh()
	RunTextureMesh()
}
