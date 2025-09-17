package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashgraph/bhash/internal/tools"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		runInstall(os.Args[2:])
	case "shacl":
		runShacl(os.Args[2:])
	case "sparql":
		runSparql(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <install|shacl|sparql> [options]\n", filepath.Base(os.Args[0]))
}

func runInstall(args []string) {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	robotVersion := fs.String("robot-version", tools.DefaultRobotVersion, "ROBOT version to install")
	shaclVersion := fs.String("shacl-version", tools.DefaultShaclVersion, "TopBraid SHACL distribution version")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg := loadConfig()
	cfg.RobotVersion = *robotVersion
	cfg.ShaclVersion = *shaclVersion

	if err := tools.InstallRobot(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "install robot: %v\n", err)
		os.Exit(1)
	}
	if err := tools.InstallShacl(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "install shacl: %v\n", err)
		os.Exit(1)
	}
}

func runShacl(args []string) {
	fs := flag.NewFlagSet("shacl", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cfg := loadConfig()
	if err := tools.RunShacl(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func runSparql(args []string) {
	fs := flag.NewFlagSet("sparql", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cfg := loadConfig()
	if err := tools.RunSparql(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func loadConfig() *tools.Config {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "determine cwd: %v\n", err)
		os.Exit(1)
	}
	root, err := tools.FindRepoRoot(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "locate repo root: %v\n", err)
		os.Exit(1)
	}
	return tools.NewConfig(root)
}
