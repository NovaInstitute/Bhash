package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashgraph/bhash/internal/fluree"
	"github.com/hashgraph/bhash/internal/tools"
)

var outputWriter io.Writer = os.Stdout

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
	case "fluree":
		runFluree(os.Args[2:])
	case "hedera":
		runHedera(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <install|shacl|sparql|fluree|hedera> [options]\n", filepath.Base(os.Args[0]))
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

func runFluree(args []string) {
	if len(args) == 0 {
		flureeUsage()
		os.Exit(1)
	}
	switch args[0] {
	case "create-dataset":
		runFlureeCreateDataset(args[1:])
	case "transact":
		runFlureeTransact(args[1:])
	case "generate-sparql":
		runFlureeGenerate(args[1:], "generate-sparql")
	case "generate-answer":
		runFlureeGenerate(args[1:], "generate-answer")
	case "generate-prompt":
		runFlureeGenerate(args[1:], "generate-prompt")
	default:
		flureeUsage()
		os.Exit(1)
	}
}

func flureeUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s fluree <create-dataset|transact|generate-sparql|generate-answer|generate-prompt> [options]\n", filepath.Base(os.Args[0]))
}

func runFlureeCreateDataset(args []string) {
	fs := flag.NewFlagSet("fluree create-dataset", flag.ExitOnError)
	apiToken := fs.String("api-token", "", "Fluree API token (defaults to $FLUREE_API_TOKEN)")
	tenant := fs.String("tenant", "", "Fluree tenant handle (defaults to $FLUREE_HANDLE)")
	baseURL := fs.String("base-url", "", "Fluree API base URL (defaults to $FLUREE_BASE_URL)")
	owner := fs.String("owner", "", "Owner handle responsible for the dataset")
	datasetName := fs.String("dataset-name", "", "Dataset name")
	storageType := fs.String("storage-type", "sparql", "Storage type")
	description := fs.String("description", "", "Dataset description")
	visibility := fs.String("visibility", "private", "Dataset visibility")
	tags := newStringSliceFlag()
	fs.Var(tags, "tag", "Tag to apply to the dataset (may be repeated)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg := mustFlureeConfig(*apiToken, *tenant, *baseURL)
	if *owner == "" || *datasetName == "" || *description == "" {
		fmt.Fprintln(os.Stderr, "owner, dataset-name, and description are required")
		os.Exit(1)
	}

	client := fluree.NewClient(cfg, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	result, err := client.CreateDataset(ctx, *owner, fluree.CreateDatasetRequest{
		DatasetName: *datasetName,
		StorageType: *storageType,
		Description: *description,
		Visibility:  *visibility,
		Tags:        tags.Values(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	printJSON(result)
}

func runFlureeTransact(args []string) {
	fs := flag.NewFlagSet("fluree transact", flag.ExitOnError)
	apiToken := fs.String("api-token", "", "Fluree API token (defaults to $FLUREE_API_TOKEN)")
	tenant := fs.String("tenant", "", "Fluree tenant handle (defaults to $FLUREE_HANDLE)")
	baseURL := fs.String("base-url", "", "Fluree API base URL (defaults to $FLUREE_BASE_URL)")
	ledger := fs.String("ledger", "", "Ledger identifier")
	insertPath := fs.String("insert", "", "Path to JSON file containing an array of insert statements")
	deletePath := fs.String("delete", "", "Path to JSON file containing an array of delete statements")
	wherePath := fs.String("where", "", "Path to JSON file containing a where clause array")
	contextPath := fs.String("context", "", "Path to JSON file containing a JSON-LD context object")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg := mustFlureeConfig(*apiToken, *tenant, *baseURL)
	if strings.TrimSpace(*ledger) == "" {
		fmt.Fprintln(os.Stderr, "ledger is required")
		os.Exit(1)
	}

	req := fluree.TransactionRequest{Ledger: *ledger}
	if *insertPath != "" {
		values, err := loadJSONArrayMap(*insertPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load insert payload: %v\n", err)
			os.Exit(1)
		}
		req.Insert = values
	}
	if *deletePath != "" {
		values, err := loadJSONArrayMap(*deletePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load delete payload: %v\n", err)
			os.Exit(1)
		}
		req.Delete = values
	}
	if *wherePath != "" {
		values, err := loadJSONArrayMap(*wherePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load where payload: %v\n", err)
			os.Exit(1)
		}
		req.Where = values
	}
	if *contextPath != "" {
		value, err := loadJSONMap(*contextPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load context payload: %v\n", err)
			os.Exit(1)
		}
		req.Context = value
	}

	client := fluree.NewClient(cfg, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	result, err := client.Transact(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	printJSON(result)
}

func runFlureeGenerate(args []string, endpoint string) {
	fs := flag.NewFlagSet("fluree "+endpoint, flag.ExitOnError)
	apiToken := fs.String("api-token", "", "Fluree API token (defaults to $FLUREE_API_TOKEN)")
	tenant := fs.String("tenant", "", "Fluree tenant handle (defaults to $FLUREE_HANDLE)")
	baseURL := fs.String("base-url", "", "Fluree API base URL (defaults to $FLUREE_BASE_URL)")
	owner := fs.String("owner", "", "Owner handle responsible for the datasets")
	prompt := fs.String("prompt", "", "Prompt/question to send to Fluree")
	datasets := newStringSliceFlag()
	fs.Var(datasets, "dataset", "Dataset identifier to include (may be repeated)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg := mustFlureeConfig(*apiToken, *tenant, *baseURL)
	if *owner == "" || *prompt == "" {
		fmt.Fprintln(os.Stderr, "owner and prompt are required")
		os.Exit(1)
	}

	request := fluree.PromptRequest{Datasets: datasets.Values(), Prompt: *prompt}
	client := fluree.NewClient(cfg, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var (
		result any
		err    error
	)
	switch endpoint {
	case "generate-sparql":
		result, err = client.GenerateSPARQL(ctx, *owner, request)
	case "generate-answer":
		result, err = client.GenerateAnswer(ctx, *owner, request)
	case "generate-prompt":
		result, err = client.GeneratePrompt(ctx, *owner, request)
	default:
		err = errors.New("unsupported endpoint")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	printJSON(result)
}

func mustFlureeConfig(apiToken, tenant, baseURL string) fluree.Config {
	cfg, err := fluree.EnvConfigFromLookup(func(key string) (string, bool) {
		value, ok := os.LookupEnv(key)
		return value, ok
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	cfg = cfg.WithOverrides(apiToken, tenant, baseURL)
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	return cfg
}

func loadJSONArrayMap(path string) ([]map[string]any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var payload []map[string]any
	if err := json.NewDecoder(file).Decode(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func loadJSONMap(path string) (map[string]any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	payload := make(map[string]any)
	if err := json.NewDecoder(file).Decode(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func printJSON(value any) {
	encoder := json.NewEncoder(outputWriter)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		fmt.Fprintf(os.Stderr, "encode JSON: %v\n", err)
		os.Exit(1)
	}
}

type stringSliceFlag struct {
	values []string
}

func newStringSliceFlag() *stringSliceFlag {
	return &stringSliceFlag{}
}

func (s *stringSliceFlag) String() string {
	return strings.Join(s.values, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	s.values = append(s.values, value)
	return nil
}

func (s *stringSliceFlag) Values() []string {
	return append([]string(nil), s.values...)
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
