package openmvg

//go:generate mockgen -source=./openmvg.go -destination=../../mocks/mock_openmvg.go -package=mocks
type OpenMVGServiceInterface interface {
	RunHealthCheck()
	SfMSequentialPipeline()
	RunSfMInitImageListing()
	RunSfMComputeFeatures()
	RunSfMPairGenerator()
	RunSfMComputeMatches()
	RunSfMGeometricFilter()
	RunSfMReconstruction()
	RunSfMComputeSfMDataColor()
	PopulateTmpDir()
}
