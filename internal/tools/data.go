package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func (c *Config) datasetPaths() ([]string, error) {
	base := []string{
		filepath.Join(c.RepoRoot, "ontology", "examples", "core-consensus.ttl"),
		filepath.Join(c.RepoRoot, "ontology", "examples", "token-compliance.ttl"),
		filepath.Join(c.RepoRoot, "ontology", "examples", "smart-contracts.ttl"),
		filepath.Join(c.RepoRoot, "ontology", "examples", "file-schedule.ttl"),
		filepath.Join(c.RepoRoot, "ontology", "examples", "mirror-analytics.ttl"),
		filepath.Join(c.RepoRoot, "ontology", "examples", "hiero.ttl"),
                filepath.Join(c.RepoRoot, "ontology", "examples", "alignment-impact.ttl"),
	}
	fixturesDir := filepath.Join(c.RepoRoot, "tests", "fixtures", "datasets")
	entries, err := os.ReadDir(fixturesDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".ttl" {
			base = append(base, filepath.Join(fixturesDir, entry.Name()))
		}
	}
	if len(base) > 7 {
		sort.Strings(base[7:])
	}
	return existingFiles(base), nil
}

func (c *Config) shapePaths() ([]string, error) {
	shapesDir := filepath.Join(c.RepoRoot, "ontology", "shapes")
	entries, err := os.ReadDir(shapesDir)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".ttl" {
			paths = append(paths, filepath.Join(shapesDir, entry.Name()))
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func (c *Config) queryPaths() ([]string, error) {
	queriesDir := filepath.Join(c.RepoRoot, "tests", "queries")
	entries, err := os.ReadDir(queriesDir)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".rq" {
			paths = append(paths, filepath.Join(queriesDir, entry.Name()))
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func existingFiles(paths []string) []string {
	filtered := paths[:0]
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			filtered = append(filtered, path)
		}
	}
	return filtered
}

func ensureNonEmpty(paths []string, label string) error {
	if len(paths) == 0 {
		return fmt.Errorf("no %s files located", label)
	}
	return nil
}
