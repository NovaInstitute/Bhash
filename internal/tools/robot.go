package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func mergeWithRobot(robot string, inputs []string, output string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("merge: no input files provided")
	}
	args := []string{"merge"}
	for _, path := range inputs {
		args = append(args, "--input", path)
	}
	args = append(args, "--output", output)

	cmd := exec.Command(robot, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("robot merge: %w", err)
	}
	return nil
}

func runRobotQuery(robot, dataFile, queryFile, outputFile string) error {
	if err := os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil {
		return err
	}
	args := []string{"query", "--input", dataFile, "--query", queryFile, outputFile}
	cmd := exec.Command(robot, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("robot query (%s): %w", filepath.Base(queryFile), err)
	}
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		header, err := extractSelectHeader(queryFile)
		if err != nil {
			return err
		}
		if err := os.WriteFile(outputFile, []byte(header+"\n"), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func extractSelectHeader(queryFile string) (string, error) {
	data, err := os.ReadFile(queryFile)
	if err != nil {
		return "", err
	}
	lower := bytes.ToLower(data)
	idx := bytes.Index(lower, []byte("select"))
	if idx == -1 {
		return "", fmt.Errorf("query %s: unable to infer header", queryFile)
	}
	fragment := lower[idx:]
	end := bytes.IndexByte(fragment, '{')
	if end == -1 {
		end = len(fragment)
	}
	fields := bytes.Fields(fragment[:end])
	var vars []string
	for _, field := range fields[1:] {
		if len(field) > 0 && field[0] == '?' {
			vars = append(vars, string(field[1:]))
		}
	}
	if len(vars) == 0 {
		return "", fmt.Errorf("query %s: unable to infer header", queryFile)
	}
	return strings.Join(vars, ","), nil
}
