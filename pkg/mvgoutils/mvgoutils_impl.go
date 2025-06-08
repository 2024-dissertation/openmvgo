package mvgoutils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type MvgoUtilsImpl struct{}

func NewMvgoUtils() OpenmvgoUtilsInterface {
	return &MvgoUtilsImpl{}
}

func (u *MvgoUtilsImpl) Check(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		panic(e)
	}
}

// RunCommand runs a command with arguments and prints its stdout/stderr in real-time.
func (u *MvgoUtilsImpl) RunCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("→ Running: %s %v\n", name, args)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command failed: %s %v: %w", name, args, err)
	}
	return nil
}

// EnsureDir creates the directory (and parent dirs) if it doesn't exist.
func (u *MvgoUtilsImpl) EnsureDir(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("could not resolve absolute path: %w", err)
	}

	if _, err := os.Stat(abs); os.IsNotExist(err) {
		fmt.Printf("→ Creating directory: %s\n", abs)
		if err := os.MkdirAll(abs, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", abs, err)
		}
	}
	return nil
}

// DownloadFile downloads a file from the given URL and saves it to a temporary file.
func (u *MvgoUtilsImpl) DownloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: %s returned status code %d", url, resp.StatusCode)
	}

	fileName := filepath.Base(url)
	out, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", fileName, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to file %s: %w", fileName, err)
	}

	fmt.Printf("→ Downloaded file: %s\n", fileName)
	return out.Name(), nil
}

func (u *MvgoUtilsImpl) CopyFile(src, dst string) error {
	fmt.Printf("→ Copying file from %s to %s\n", src, dst)
	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", src, err)
	}
	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to destination file %s: %w", dst, err)
	}
	fmt.Printf("→ Copied file from %s to %s\n", src, dst)
	return nil
}
