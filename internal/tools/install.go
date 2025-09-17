package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func InstallRobot(cfg *Config) error {
	if err := cfg.EnsureBaseDirs(); err != nil {
		return err
	}
	jarPath := cfg.RobotJarPath()
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		if err := downloadFile(cfg.RobotDownloadURL(), jarPath); err != nil {
			return err
		}
	}
	wrapper := cfg.RobotExecutable()
	if err := writeScript(wrapper, fmt.Sprintf("#!/usr/bin/env bash\nexec java -jar '%s' \"$@\"\n", jarPath)); err != nil {
		return err
	}
	return nil
}

func InstallShacl(cfg *Config) error {
	if err := cfg.EnsureBaseDirs(); err != nil {
		return err
	}

	versionDir := cfg.ShaclVersionDir()
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		archivePath := cfg.ShaclArchivePath()
		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			if err := downloadFile(cfg.ShaclDownloadURL(), archivePath); err != nil {
				return err
			}
		}
		if err := unzip(archivePath, cfg.ShaclInstallDir()); err != nil {
			return err
		}
	}

	scriptSrc := filepath.Join(versionDir, "bin", "shaclvalidate.sh")
	if err := os.Chmod(scriptSrc, 0o755); err != nil && !os.IsPermission(err) {
		return err
	}

	wrapper := cfg.ShaclValidateScript()
	script := fmt.Sprintf("#!/usr/bin/env bash\n'%s' \"$@\"\n", scriptSrc)
	if err := writeScript(wrapper, script); err != nil {
		return err
	}
	return nil
}

func writeScript(path, contents string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(contents), 0o755); err != nil {
		return err
	}
	return nil
}
