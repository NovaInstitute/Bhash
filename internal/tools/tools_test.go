package tools

import (
	"archive/zip"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewConfigSetsDefaults(t *testing.T) {
	repoRoot := t.TempDir()
	cfg := NewConfig(repoRoot)

	if cfg.RepoRoot != repoRoot {
		t.Fatalf("expected RepoRoot %q, got %q", repoRoot, cfg.RepoRoot)
	}

	expectedBuild := filepath.Join(repoRoot, "build")
	expectedTools := filepath.Join(expectedBuild, "tools")
	expectedBin := filepath.Join(expectedTools, "bin")

	if cfg.BuildDir != expectedBuild {
		t.Errorf("expected BuildDir %q, got %q", expectedBuild, cfg.BuildDir)
	}
	if cfg.ToolsDir != expectedTools {
		t.Errorf("expected ToolsDir %q, got %q", expectedTools, cfg.ToolsDir)
	}
	if cfg.BinDir != expectedBin {
		t.Errorf("expected BinDir %q, got %q", expectedBin, cfg.BinDir)
	}
	if cfg.RobotVersion != DefaultRobotVersion {
		t.Errorf("expected RobotVersion %q, got %q", DefaultRobotVersion, cfg.RobotVersion)
	}
	if cfg.ShaclVersion != DefaultShaclVersion {
		t.Errorf("expected ShaclVersion %q, got %q", DefaultShaclVersion, cfg.ShaclVersion)
	}
}

func TestConfigHelpers(t *testing.T) {
	repoRoot := t.TempDir()
	cfg := NewConfig(repoRoot)

	if got := cfg.RobotJarPath(); !strings.HasSuffix(got, filepath.Join("tools", "robot", "robot.jar")) {
		t.Errorf("unexpected RobotJarPath %q", got)
	}
	if got := cfg.RobotExecutable(); !strings.HasSuffix(got, filepath.Join("bin", "robot")) {
		t.Errorf("unexpected RobotExecutable %q", got)
	}
	expectedRobotURL := "https://github.com/ontodev/robot/releases/download/v" + DefaultRobotVersion + "/robot.jar"
	if got := cfg.RobotDownloadURL(); got != expectedRobotURL {
		t.Errorf("unexpected RobotDownloadURL %q", got)
	}

	if got := cfg.ShaclArchivePath(); !strings.HasSuffix(got, filepath.Join("downloads", "shacl-"+DefaultShaclVersion+"-bin.zip")) {
		t.Errorf("unexpected ShaclArchivePath %q", got)
	}
	expectedShaclURL := "https://repo1.maven.org/maven2/org/topbraid/shacl/" + DefaultShaclVersion + "/shacl-" + DefaultShaclVersion + "-bin.zip"
	if got := cfg.ShaclDownloadURL(); got != expectedShaclURL {
		t.Errorf("unexpected ShaclDownloadURL %q", got)
	}
	if got := cfg.ShaclInstallDir(); !strings.HasSuffix(got, filepath.Join("tools", "shacl")) {
		t.Errorf("unexpected ShaclInstallDir %q", got)
	}
	if got := cfg.ShaclVersionDir(); !strings.HasSuffix(got, filepath.Join("tools", "shacl", "shacl-"+DefaultShaclVersion)) {
		t.Errorf("unexpected ShaclVersionDir %q", got)
	}
	if got := cfg.ShaclValidateScript(); !strings.HasSuffix(got, filepath.Join("bin", "shaclvalidate")) {
		t.Errorf("unexpected ShaclValidateScript %q", got)
	}
}

func TestEnsureBaseDirs(t *testing.T) {
	repoRoot := t.TempDir()
	cfg := NewConfig(repoRoot)
	if err := cfg.EnsureBaseDirs(); err != nil {
		t.Fatalf("EnsureBaseDirs returned error: %v", err)
	}

	dirs := []string{
		cfg.BuildDir,
		cfg.ToolsDir,
		cfg.BinDir,
		filepath.Join(cfg.ToolsDir, "robot"),
		filepath.Join(cfg.ToolsDir, "downloads"),
	}
	for _, dir := range dirs {
		if info, err := os.Stat(dir); err != nil {
			t.Fatalf("expected directory %s to exist: %v", dir, err)
		} else if !info.IsDir() {
			t.Fatalf("expected %s to be a directory", dir)
		}
	}
}

func createFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", path, err)
	}
	if err := os.WriteFile(path, []byte("content"), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

func TestDatasetPaths(t *testing.T) {
	repoRoot := t.TempDir()
	cfg := NewConfig(repoRoot)

	baseDataset := filepath.Join(repoRoot, "ontology", "examples", "core-consensus.ttl")
	otherDataset := filepath.Join(repoRoot, "ontology", "examples", "hiero.ttl")
	createFile(t, baseDataset)
	createFile(t, otherDataset)

	fixtureA := filepath.Join(repoRoot, "tests", "fixtures", "datasets", "a.ttl")
	fixtureB := filepath.Join(repoRoot, "tests", "fixtures", "datasets", "b.ttl")
	fixtureOther := filepath.Join(repoRoot, "tests", "fixtures", "datasets", "ignore.txt")
	createFile(t, fixtureB)
	createFile(t, fixtureA)
	createFile(t, fixtureOther)

	paths, err := cfg.datasetPaths()
	if err != nil {
		t.Fatalf("datasetPaths returned error: %v", err)
	}

	expected := []string{baseDataset, otherDataset, fixtureA, fixtureB}
	if !reflect.DeepEqual(paths, expected) {
		t.Fatalf("unexpected dataset paths:\n got: %v\nwant: %v", paths, expected)
	}
}

func TestShapePaths(t *testing.T) {
	repoRoot := t.TempDir()
	cfg := NewConfig(repoRoot)

	shapeA := filepath.Join(repoRoot, "ontology", "shapes", "b.ttl")
	shapeB := filepath.Join(repoRoot, "ontology", "shapes", "a.ttl")
	other := filepath.Join(repoRoot, "ontology", "shapes", "ignore.txt")
	createFile(t, shapeA)
	createFile(t, shapeB)
	createFile(t, other)

	paths, err := cfg.shapePaths()
	if err != nil {
		t.Fatalf("shapePaths returned error: %v", err)
	}

	expected := []string{shapeB, shapeA}
	if !reflect.DeepEqual(paths, expected) {
		t.Fatalf("unexpected shape paths: got %v want %v", paths, expected)
	}
}

func TestQueryPaths(t *testing.T) {
	repoRoot := t.TempDir()
	cfg := NewConfig(repoRoot)

	queryA := filepath.Join(repoRoot, "tests", "queries", "b.rq")
	queryB := filepath.Join(repoRoot, "tests", "queries", "a.rq")
	other := filepath.Join(repoRoot, "tests", "queries", "ignore.txt")
	createFile(t, queryA)
	createFile(t, queryB)
	createFile(t, other)

	paths, err := cfg.queryPaths()
	if err != nil {
		t.Fatalf("queryPaths returned error: %v", err)
	}

	expected := []string{queryB, queryA}
	if !reflect.DeepEqual(paths, expected) {
		t.Fatalf("unexpected query paths: got %v want %v", paths, expected)
	}
}

func TestExistingFiles(t *testing.T) {
	dir := t.TempDir()
	present := filepath.Join(dir, "present.ttl")
	missing := filepath.Join(dir, "missing.ttl")
	createFile(t, present)

	filtered := existingFiles([]string{present, missing})
	expected := []string{present}
	if !reflect.DeepEqual(filtered, expected) {
		t.Fatalf("unexpected filtered paths: got %v want %v", filtered, expected)
	}
}

func TestEnsureNonEmpty(t *testing.T) {
	if err := ensureNonEmpty([]string{"a"}, "test"); err != nil {
		t.Fatalf("ensureNonEmpty returned unexpected error: %v", err)
	}
	err := ensureNonEmpty(nil, "test")
	if err == nil || !strings.Contains(err.Error(), "no test files located") {
		t.Fatalf("expected error about missing files, got %v", err)
	}
}

func TestFindRepoRoot(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.Mkdir(gitDir, 0o755); err != nil {
		t.Fatalf("Mkdir(%s): %v", gitDir, err)
	}
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", nested, err)
	}

	found, err := FindRepoRoot(nested)
	if err != nil {
		t.Fatalf("FindRepoRoot returned error: %v", err)
	}
	if found != root {
		t.Fatalf("expected repo root %s, got %s", root, found)
	}

	_, err = FindRepoRoot(filepath.Join(os.TempDir(), "nonexistent"))
	if err == nil {
		t.Fatalf("expected error when repo root not found")
	}
}

func TestDownloadFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "downloaded")
	}))
	defer server.Close()

	destDir := t.TempDir()
	dest := filepath.Join(destDir, "nested", "file.txt")
	if err := downloadFile(server.URL, dest); err != nil {
		t.Fatalf("downloadFile returned error: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", dest, err)
	}
	if string(data) != "downloaded" {
		t.Fatalf("unexpected file contents: %q", data)
	}

	badServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer badServer.Close()

	if err := downloadFile(badServer.URL, dest); err == nil {
		t.Fatalf("expected error for non-200 response")
	}
}

func createZip(t *testing.T, entries map[string]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "archive.zip")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Create zip: %v", err)
	}
	defer file.Close()

	zw := zip.NewWriter(file)
	for name, content := range entries {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("Create entry %s: %v", name, err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatalf("Write entry %s: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("Close zip: %v", err)
	}
	return path
}

func TestUnzip(t *testing.T) {
	zipPath := createZip(t, map[string]string{"dir/file.txt": "hello"})
	dest := filepath.Join(t.TempDir(), "out")
	if err := unzip(zipPath, dest); err != nil {
		t.Fatalf("unzip returned error: %v", err)
	}

	extracted := filepath.Join(dest, "dir", "file.txt")
	data, err := os.ReadFile(extracted)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", extracted, err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected extracted content: %q", data)
	}

	maliciousZip := createZip(t, map[string]string{"../evil.txt": "bad"})
	err = unzip(maliciousZip, dest)
	if err == nil {
		t.Fatalf("expected unzip to reject path traversal entry")
	}
	if !strings.Contains(err.Error(), "illegal file path") {
		t.Fatalf("unexpected error for malicious zip: %v", err)
	}
}

func TestExtractSelectHeader(t *testing.T) {
	queryFile := filepath.Join(t.TempDir(), "query.rq")
	query := `PREFIX ex: <http://example.com/>
SELECT ?a ?b ?c WHERE {
    ?a ex:related ?b .
}`
	if err := os.WriteFile(queryFile, []byte(query), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", queryFile, err)
	}

	header, err := extractSelectHeader(queryFile)
	if err != nil {
		t.Fatalf("extractSelectHeader returned error: %v", err)
	}
	if header != "a,b,c" {
		t.Fatalf("unexpected header %q", header)
	}

	badQuery := filepath.Join(t.TempDir(), "bad.rq")
	if err := os.WriteFile(badQuery, []byte("ASK {}"), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", badQuery, err)
	}
	_, err = extractSelectHeader(badQuery)
	if err == nil {
		t.Fatalf("expected error for non-select query")
	}
	if !errors.Is(err, os.ErrNotExist) && !strings.Contains(err.Error(), "unable to infer header") {
		t.Fatalf("unexpected error: %v", err)
	}
}
