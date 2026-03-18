package commands

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/p-arndt/dorcs/internal/github"
)

type mockBootstrapGitHubClient struct {
	files              map[string][]byte
	errs               map[string]error
	defaultBranch      string
	defaultBranchErr   error
	discoveredMarkdown []string
}

func (m *mockBootstrapGitHubClient) DiscoverMarkdownFiles(owner, repo, branch, rootPath string) ([]string, error) {
	return m.discoveredMarkdown, nil
}

func (m *mockBootstrapGitHubClient) ListDirectory(owner, repo, branch, dirPath string) ([]github.TreeEntry, error) {
	return nil, nil
}

func (m *mockBootstrapGitHubClient) FetchMarkdown(owner, repo, branch, filePath string) ([]byte, error) {
	if err, ok := m.errs[filePath]; ok {
		return nil, err
	}
	if content, ok := m.files[filePath]; ok {
		return content, nil
	}
	return nil, errors.New("file not found: " + filePath + " (404)")
}

func (m *mockBootstrapGitHubClient) GetDefaultBranch(owner, repo string) (string, error) {
	if m.defaultBranchErr != nil {
		return "", m.defaultBranchErr
	}
	if m.defaultBranch != "" {
		return m.defaultBranch, nil
	}
	return "main", nil
}

func TestLoadConfigFromRepoPrefersRoot(t *testing.T) {
	restore := stubBootstrapGitHubClient(t, &mockBootstrapGitHubClient{
		files: map[string][]byte{
			"dorcs.yaml":      []byte("site:\n  title: Root Config\n"),
			"docs/dorcs.yaml": []byte("site:\n  title: Docs Config\n"),
		},
	})
	defer restore()

	result, err := loadConfigFromRepo(t.TempDir(), "owner/repo/tree/main/docs", "./docs")
	if err != nil {
		t.Fatalf("loadConfigFromRepo() error = %v", err)
	}

	if result.Source != "dorcs.yaml" {
		t.Fatalf("loadConfigFromRepo() source = %q, want dorcs.yaml", result.Source)
	}
	if result.Config.Site.Title != "Root Config" {
		t.Fatalf("loadConfigFromRepo() title = %q, want Root Config", result.Config.Site.Title)
	}
	if !result.Config.GitHub.Enabled || result.Config.GitHub.Repository != "https://github.com/owner/repo/tree/main/docs" {
		t.Fatalf("loadConfigFromRepo() GitHub override = %+v", result.Config.GitHub)
	}
}

func TestLoadConfigFromRepoFallsBackToRepoPath(t *testing.T) {
	restore := stubBootstrapGitHubClient(t, &mockBootstrapGitHubClient{
		files: map[string][]byte{
			"docs/dorcs.yaml": []byte("site:\n  title: Docs Config\n"),
		},
	})
	defer restore()

	result, err := loadConfigFromRepo(t.TempDir(), "owner/repo/tree/main/docs", "./docs")
	if err != nil {
		t.Fatalf("loadConfigFromRepo() error = %v", err)
	}

	if result.Source != "docs/dorcs.yaml" {
		t.Fatalf("loadConfigFromRepo() source = %q, want docs/dorcs.yaml", result.Source)
	}
	if result.Config.Site.Title != "Docs Config" {
		t.Fatalf("loadConfigFromRepo() title = %q, want Docs Config", result.Config.Site.Title)
	}
}

func TestLoadConfigFromRepoReturnsDefaultsWhenMissing(t *testing.T) {
	restore := stubBootstrapGitHubClient(t, &mockBootstrapGitHubClient{
		files: map[string][]byte{},
	})
	defer restore()

	result, err := loadConfigFromRepo(t.TempDir(), "owner/repo/tree/main/docs", "./docs")
	if err != nil {
		t.Fatalf("loadConfigFromRepo() error = %v", err)
	}

	if result.Source != "defaults" {
		t.Fatalf("loadConfigFromRepo() source = %q, want defaults", result.Source)
	}
	if result.Config.Site.Title != "Documentation" {
		t.Fatalf("loadConfigFromRepo() title = %q, want Documentation", result.Config.Site.Title)
	}
	if !result.Config.GitHub.Enabled {
		t.Fatal("loadConfigFromRepo() should enable GitHub in repo mode")
	}
}

func TestLoadConfigFromRepoMalformedConfig(t *testing.T) {
	restore := stubBootstrapGitHubClient(t, &mockBootstrapGitHubClient{
		files: map[string][]byte{
			"dorcs.yaml": []byte("nav:\n  items:\n    - Broken: {}\n"),
		},
	})
	defer restore()

	_, err := loadConfigFromRepo(t.TempDir(), "owner/repo", "./docs")
	if err == nil {
		t.Fatal("loadConfigFromRepo() expected error for malformed config")
	}
}

func TestLoadConfigFromRepoFetchFailure(t *testing.T) {
	restore := stubBootstrapGitHubClient(t, &mockBootstrapGitHubClient{
		errs: map[string]error{
			"dorcs.yaml": errors.New("authentication failed: 401"),
		},
	})
	defer restore()

	_, err := loadConfigFromRepo(t.TempDir(), "owner/repo", "./docs")
	if err == nil {
		t.Fatal("loadConfigFromRepo() expected fetch error")
	}
}

func TestLoadConfigWithBootstrapConfigFileWinsOverRepo(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "dorcs.yaml")
	if err := os.WriteFile(configPath, []byte("site:\n  title: Local Config\n"), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	restore := stubBootstrapGitHubClient(t, nil)
	defer restore()

	result, err := loadConfigWithBootstrap(dir, dir, configPath, "owner/repo")
	if err != nil {
		t.Fatalf("loadConfigWithBootstrap() error = %v", err)
	}

	if result.Source != configPath {
		t.Fatalf("loadConfigWithBootstrap() source = %q, want %q", result.Source, configPath)
	}
	if result.Config.Site.Title != "Local Config" {
		t.Fatalf("loadConfigWithBootstrap() title = %q, want Local Config", result.Config.Site.Title)
	}
}

func TestLoadConfigWithBootstrapInvalidRepo(t *testing.T) {
	restore := stubBootstrapGitHubClient(t, &mockBootstrapGitHubClient{})
	defer restore()

	_, err := loadConfigWithBootstrap(t.TempDir(), t.TempDir(), "", "https://github.com/owner")
	if err == nil {
		t.Fatal("loadConfigWithBootstrap() expected invalid repo error")
	}
}

func TestLoadConfigFromRepoPreservesConfiguredDocsPathForRepoRoot(t *testing.T) {
	restore := stubBootstrapGitHubClient(t, &mockBootstrapGitHubClient{
		files: map[string][]byte{
			"dorcs.yaml": []byte("github:\n  repository: https://github.com/owner/repo/tree/main/docs\n"),
		},
	})
	defer restore()

	result, err := loadConfigFromRepo(t.TempDir(), "https://github.com/owner/repo", "./docs")
	if err != nil {
		t.Fatalf("loadConfigFromRepo() error = %v", err)
	}

	want := "https://github.com/owner/repo/tree/main/docs"
	if result.Config.GitHub.Repository != want {
		t.Fatalf("loadConfigFromRepo() repository = %q, want %q", result.Config.GitHub.Repository, want)
	}
}

func TestLoadConfigFromRepoUsesDirHintWhenConfigHasNoGitHubRepository(t *testing.T) {
	restore := stubBootstrapGitHubClient(t, &mockBootstrapGitHubClient{
		files: map[string][]byte{
			"dorcs.yaml": []byte("languages:\n  default: en\n  enabled:\n    - code: en\n      name: English\n"),
		},
	})
	defer restore()

	result, err := loadConfigFromRepo(t.TempDir(), "https://github.com/owner/repo", "./docs")
	if err != nil {
		t.Fatalf("loadConfigFromRepo() error = %v", err)
	}

	want := "https://github.com/owner/repo/tree/main/docs"
	if result.Config.GitHub.Repository != want {
		t.Fatalf("loadConfigFromRepo() repository = %q, want %q", result.Config.GitHub.Repository, want)
	}
}

func stubBootstrapGitHubClient(t *testing.T, client github.ClientAPI) func() {
	t.Helper()

	original := newGitHubClient
	newGitHubClient = func(token string, cache *github.Cache, cacheTTL time.Duration) github.ClientAPI {
		if client == nil {
			t.Fatal("newGitHubClient should not have been called")
		}
		return client
	}

	return func() {
		newGitHubClient = original
	}
}
