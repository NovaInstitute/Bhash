package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RunSparql(cfg *Config) error {
	datasets, err := cfg.datasetPaths()
	if err != nil {
		return err
	}
	if err := ensureNonEmpty(datasets, "dataset"); err != nil {
		return err
	}

	queries, err := cfg.queryPaths()
	if err != nil {
		return err
	}
	if err := ensureNonEmpty(queries, "query"); err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp(cfg.BuildDir, "sparql-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	dataFile := filepath.Join(tempDir, "data.ttl")
	if err := mergeWithRobot(cfg.RobotExecutable(), datasets, dataFile); err != nil {
		return err
	}

	outputDir := filepath.Join(cfg.BuildDir, "queries")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	resultsDir := filepath.Join(cfg.RepoRoot, "tests", "fixtures", "results")
	success := true
	var failures []string

	for _, query := range queries {
		name := filepath.Base(query)
		fmt.Printf("Running %s...\n", name)
		output := filepath.Join(outputDir, strings.TrimSuffix(name, filepath.Ext(name))+".csv")
		if err := runRobotQuery(cfg.RobotExecutable(), dataFile, query, output); err != nil {
			return err
		}
		expected := filepath.Join(resultsDir, filepath.Base(output))
		match, err := compareCSV(expected, output)
		if err != nil {
			return err
		}
		if !match {
			success = false
			failures = append(failures, name)
		}
	}

	if !success {
		return fmt.Errorf("sparql regression failures: %s", strings.Join(failures, ", "))
	}
	return nil
}

func compareCSV(expectedPath, actualPath string) (bool, error) {
	expected, err := readCSVLines(expectedPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No expected results for %s; skipping comparison.\n", filepath.Base(actualPath))
			return true, nil
		}
		return false, err
	}
	actual, err := readCSVLines(actualPath)
	if err != nil {
		return false, err
	}
	if len(expected) != len(actual) {
		printCSVDelta(expected, actual)
		return false, nil
	}
	for i := range expected {
		if expected[i] != actual[i] {
			printCSVDelta(expected, actual)
			return false, nil
		}
	}
	return true, nil
}

func readCSVLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			filtered = append(filtered, normalizeCSVLine(trimmed))
		}
	}
	return filtered, nil
}

func printCSVDelta(expected, actual []string) {
	fmt.Println("Expected:")
	for _, line := range expected {
		fmt.Printf("  %s\n", line)
	}
	fmt.Println("Actual:")
	for _, line := range actual {
		fmt.Printf("  %s\n", line)
	}
}

func normalizeCSVLine(line string) string {
	// Normalise time zone suffixes so that "Z" and "+00:00" compare equal.
	line = strings.ReplaceAll(line, "+00:00", "Z")
	// ROBOT may produce "-00:00" for unknown offsets; keep consistent with fixtures.
	line = strings.ReplaceAll(line, "-00:00", "Z")
	return line
}
