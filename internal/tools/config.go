package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultRobotVersion = "1.9.5"
	DefaultShaclVersion = "1.4.3"
)

type Config struct {
	RepoRoot     string
	BuildDir     string
	ToolsDir     string
	BinDir       string
	RobotVersion string
	ShaclVersion string
}

func NewConfig(repoRoot string) *Config {
	buildDir := filepath.Join(repoRoot, "build")
	toolsDir := filepath.Join(buildDir, "tools")
	binDir := filepath.Join(toolsDir, "bin")
	return &Config{
		RepoRoot:     repoRoot,
		BuildDir:     buildDir,
		ToolsDir:     toolsDir,
		BinDir:       binDir,
		RobotVersion: DefaultRobotVersion,
		ShaclVersion: DefaultShaclVersion,
	}
}

func (c *Config) RobotJarPath() string {
	return filepath.Join(c.ToolsDir, "robot", "robot.jar")
}

func (c *Config) RobotExecutable() string {
	return filepath.Join(c.BinDir, "robot")
}

func (c *Config) RobotDownloadURL() string {
	return fmt.Sprintf("https://github.com/ontodev/robot/releases/download/v%[1]s/robot.jar", c.RobotVersion)
}

func (c *Config) ShaclArchivePath() string {
	return filepath.Join(c.ToolsDir, "downloads", fmt.Sprintf("shacl-%s-bin.zip", c.ShaclVersion))
}

func (c *Config) ShaclDownloadURL() string {
	return fmt.Sprintf("https://repo1.maven.org/maven2/org/topbraid/shacl/%[1]s/shacl-%[1]s-bin.zip", c.ShaclVersion)
}

func (c *Config) ShaclInstallDir() string {
	return filepath.Join(c.ToolsDir, "shacl")
}

func (c *Config) ShaclVersionDir() string {
	return filepath.Join(c.ShaclInstallDir(), fmt.Sprintf("shacl-%s", c.ShaclVersion))
}

func (c *Config) ShaclValidateScript() string {
	return filepath.Join(c.BinDir, "shaclvalidate")
}

func (c *Config) EnsureBaseDirs() error {
	dirs := []string{c.BuildDir, c.ToolsDir, c.BinDir, filepath.Join(c.ToolsDir, "robot"), filepath.Join(c.ToolsDir, "downloads")}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}
