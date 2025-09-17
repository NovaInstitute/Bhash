package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunShacl(cfg *Config) error {
	datasets, err := cfg.datasetPaths()
	if err != nil {
		return err
	}
	if err := ensureNonEmpty(datasets, "dataset"); err != nil {
		return err
	}
	shapes, err := cfg.shapePaths()
	if err != nil {
		return err
	}
	if err := ensureNonEmpty(shapes, "shape"); err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp(cfg.BuildDir, "shacl-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	dataFile := filepath.Join(tempDir, "data.ttl")
	if err := mergeWithRobot(cfg.RobotExecutable(), datasets, dataFile); err != nil {
		return err
	}

	shapesFile := filepath.Join(tempDir, "shapes.ttl")
	if err := mergeWithRobot(cfg.RobotExecutable(), shapes, shapesFile); err != nil {
		return err
	}

	reportDir := filepath.Join(cfg.BuildDir, "reports")
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		return err
	}
	reportPath := filepath.Join(reportDir, "shacl-report.ttl")

	cmd := exec.Command(cfg.ShaclValidateScript(), "-datafile", dataFile, "-shapesfile", shapesFile)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	output := stdout.Bytes()
	fmt.Print(string(output))

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if writeErr := os.WriteFile(reportPath, output, 0o644); writeErr != nil {
				return fmt.Errorf("shacl validation failed: %v (additional error writing report: %w)", exitErr, writeErr)
			}
			return fmt.Errorf("shacl validation failed; report written to %s", reportPath)
		}
		return err
	}

	if _, err := os.Stat(reportPath); err == nil {
		os.Remove(reportPath)
	}
	return nil
}
