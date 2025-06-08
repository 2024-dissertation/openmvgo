package mvgoutils

//go:generate mockgen -source=./utils.go -destination=../../mocks/mock_utils.go -package=mocks
type OpenmvgoUtilsInterface interface {
	Check(e error)
	RunCommand(name string, args []string) error
	EnsureDir(path string) error
	DownloadFile(url string) (string, error)
	CopyFile(src, dst string) error
}
