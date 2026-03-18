package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/p-arndt/dorcs/internal/config"
	"github.com/p-arndt/dorcs/internal/github"
	"github.com/p-arndt/dorcs/internal/site"
)

var newGitHubClient = func(token string, cache *github.Cache, cacheTTL time.Duration) github.ClientAPI {
	return github.NewClient(token, cache, cacheTTL)
}

type configBootstrapResult struct {
	Config *config.Config
	Source string
}

func loadConfigWithBootstrap(rootDir, docsDir, configFile, repo string) (*configBootstrapResult, error) {
	switch {
	case strings.TrimSpace(configFile) != "":
		cfg, err := config.LoadFromFile(configFile)
		if err != nil {
			return nil, err
		}
		return &configBootstrapResult{Config: cfg, Source: configFile}, nil
	case strings.TrimSpace(repo) != "":
		return loadConfigFromRepo(rootDir, repo, docsDir)
	default:
		cfg, err := config.Load(docsDir)
		if err != nil {
			return nil, err
		}
		return &configBootstrapResult{Config: cfg}, nil
	}
}

func loadConfigFromRepo(rootDir, repo, docsDir string) (*configBootstrapResult, error) {
	token := bootstrapGitHubToken()
	client, repoInfo, _, err := buildGitHubClient(rootDir, repo, token, "")
	if err != nil {
		return nil, err
	}
	resolvedRepo := canonicalRepositoryURL(repoInfo)

	searchDirs := []string{""}
	if repoInfo.Path != "" {
		searchDirs = append(searchDirs, repoInfo.Path)
	}

	for _, dir := range searchDirs {
		for _, name := range []string{"dorcs.yaml", "dorcs.yml", "dorcs.json"} {
			candidate := name
			if dir != "" {
				candidate = strings.Trim(dir, "/") + "/" + name
			}

			data, err := client.FetchMarkdown(repoInfo.Owner, repoInfo.Repo, repoInfo.Branch, candidate)
			if err != nil {
				if isGitHubNotFoundError(err) {
					continue
				}
				return nil, fmt.Errorf("load config from %s: %w", candidate, err)
			}

			cfg, err := config.LoadFromBytes(name, data)
			if err != nil {
				return nil, fmt.Errorf("decode config from %s: %w", candidate, err)
			}
			applyRepoOverride(cfg, resolvedRepo, token, docsDir)
			return &configBootstrapResult{
				Config: cfg,
				Source: candidate,
			}, nil
		}
	}

	cfg := config.Default()
	applyRepoOverride(cfg, resolvedRepo, token, docsDir)
	return &configBootstrapResult{
		Config: cfg,
		Source: "defaults",
	}, nil
}

func setupGitHubIntegration(rootDir string, cfg *config.Config) (github.ClientAPI, *github.RepositoryInfo, error) {
	if cfg == nil || !cfg.GitHub.Enabled || strings.TrimSpace(cfg.GitHub.Repository) == "" {
		return nil, nil, nil
	}

	client, repoInfo, cacheDir, err := buildGitHubClient(rootDir, cfg.GitHub.Repository, cfg.GitHub.Token, cfg.GitHub.CacheTTL)
	if err != nil {
		return nil, nil, err
	}

	if cacheDir != "" {
		log.Printf("dorcs: using persistent cache at %s", cacheDir)
	}

	tokenState := "empty"
	switch {
	case strings.Contains(cfg.GitHub.Token, "${"):
		tokenState = "unexpanded"
	case strings.TrimSpace(cfg.GitHub.Token) != "":
		tokenState = "present"
	}
	log.Printf("dorcs: GitHub token state: %s (length=%d)", tokenState, len(strings.TrimSpace(cfg.GitHub.Token)))
	log.Printf("dorcs: GitHub integration enabled: %s/%s@%s/%s", repoInfo.Owner, repoInfo.Repo, repoInfo.Branch, repoInfo.Path)

	return client, repoInfo, nil
}

func buildGitHubClient(rootDir, repo, token, cacheTTLValue string) (github.ClientAPI, *github.RepositoryInfo, string, error) {
	repoInfo, err := github.ParseRepositoryURL(repo)
	if err != nil {
		return nil, nil, "", fmt.Errorf("parse GitHub repository URL: %w", err)
	}

	cacheTTL := time.Hour
	if cacheTTLValue != "" {
		parsed, err := time.ParseDuration(cacheTTLValue)
		if err != nil {
			return nil, nil, "", fmt.Errorf("parse cache_ttl %q: %w", cacheTTLValue, err)
		}
		cacheTTL = parsed
	}

	cacheDir := filepath.Join(rootDir, ".cache", "github")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Printf("dorcs: warning: failed to create cache directory %s: %v (using in-memory cache only)", cacheDir, err)
		cacheDir = ""
	}

	cache := github.NewCache(cacheDir)
	client := newGitHubClient(token, cache, cacheTTL)
	if repoInfo.Branch == "" {
		defaultBranch, err := client.GetDefaultBranch(repoInfo.Owner, repoInfo.Repo)
		if err != nil {
			return nil, nil, "", fmt.Errorf("get default branch: %w", err)
		}
		repoInfo.Branch = defaultBranch
	}

	return client, repoInfo, cacheDir, nil
}

func applyRepoToSite(s *site.Site, client github.ClientAPI, repoInfo *github.RepositoryInfo, githubPath string) {
	if s == nil || client == nil || repoInfo == nil {
		return
	}
	s.SetGitHubConfig(client, repoInfo.Owner, repoInfo.Repo, repoInfo.Branch, githubPath)
}

func canonicalRepositoryURL(repoInfo *github.RepositoryInfo) string {
	if repoInfo == nil {
		return ""
	}

	base := fmt.Sprintf("https://github.com/%s/%s", repoInfo.Owner, repoInfo.Repo)
	if repoInfo.Branch == "" {
		return base
	}
	if repoInfo.Path == "" {
		return fmt.Sprintf("%s/tree/%s", base, repoInfo.Branch)
	}
	return fmt.Sprintf("%s/tree/%s/%s", base, repoInfo.Branch, strings.Trim(repoInfo.Path, "/"))
}

func applyRepoOverride(cfg *config.Config, repo, token, docsDir string) {
	cfg.GitHub.Enabled = true
	cfg.GitHub.Repository = mergeRepoOverride(repo, cfg.GitHub.Repository, docsDir)
	if strings.TrimSpace(cfg.GitHub.Token) == "" {
		cfg.GitHub.Token = token
	}
}

func mergeRepoOverride(repo, configRepo, docsDir string) string {
	repo = strings.TrimSpace(repo)
	configRepo = strings.TrimSpace(configRepo)
	if repo == "" {
		return configRepo
	}
	if configRepo == "" {
		repoInfo, err := github.ParseRepositoryURL(repo)
		if err != nil {
			return repo
		}
		path := normalizeRepoDocsPath(docsDir)
		if path == "" {
			return repo
		}
		branch := repoInfo.Branch
		if branch == "" {
			return repo
		}
		return fmt.Sprintf("https://github.com/%s/%s/tree/%s/%s", repoInfo.Owner, repoInfo.Repo, branch, path)
	}

	repoInfo, err := github.ParseRepositoryURL(repo)
	if err != nil {
		return repo
	}
	configInfo, err := github.ParseRepositoryURL(configRepo)
	if err != nil {
		return repo
	}
	if repoInfo.Owner != configInfo.Owner || repoInfo.Repo != configInfo.Repo {
		return repo
	}
	if repoInfo.Path != "" {
		return repo
	}
	if configInfo.Path == "" {
		return repo
	}

	branch := configInfo.Branch
	if strings.Contains(repo, "/tree/") {
		branch = repoInfo.Branch
	}
	if branch == "" {
		branch = repoInfo.Branch
	}
	if branch == "" {
		return fmt.Sprintf("https://github.com/%s/%s", repoInfo.Owner, repoInfo.Repo)
	}

	path := strings.Trim(configInfo.Path, "/")
	if path == "" {
		path = normalizeRepoDocsPath(docsDir)
	}
	if path == "" {
		return fmt.Sprintf("https://github.com/%s/%s/tree/%s", repoInfo.Owner, repoInfo.Repo, branch)
	}

	return fmt.Sprintf("https://github.com/%s/%s/tree/%s/%s", repoInfo.Owner, repoInfo.Repo, branch, path)
}

func normalizeRepoDocsPath(docsDir string) string {
	trimmed := strings.TrimSpace(docsDir)
	if trimmed == "" {
		return ""
	}
	if filepath.IsAbs(trimmed) {
		return ""
	}

	cleaned := filepath.ToSlash(filepath.Clean(trimmed))
	switch cleaned {
	case ".", "/":
		return ""
	}

	return strings.Trim(cleaned, "/")
}

func bootstrapGitHubToken() string {
	config.LoadEnvFile()
	for _, name := range []string{"DORCS_GITHUB_TOKEN", "GITHUB_TOKEN", "GH_TOKEN"} {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			return value
		}
	}
	return ""
}

func isGitHubNotFoundError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "(404)") || strings.Contains(strings.ToLower(err.Error()), "not found"))
}
