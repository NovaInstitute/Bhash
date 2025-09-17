package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: unexpected status %s", url, resp.Status)
	}

	tmp, err := os.CreateTemp("", "bhash-download-*")
	if err != nil {
		return err
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), dest)
}
