package tools

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func unzip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		targetPath := filepath.Join(dest, file.Name)
		if !strings.HasPrefix(targetPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path %s", targetPath)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		out, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		if _, err := io.Copy(out, rc); err != nil {
			out.Close()
			rc.Close()
			return err
		}
		out.Close()
		rc.Close()
	}
	return nil
}
