package provider

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// IsolatedEnvironment creates a private provider runtime HOME and returns the
// small environment allowlist used by scanner-owned provider processes. An
// explicit environment credential gets a throwaway provider configuration. A
// login made through `everylast login` gets a dedicated Everylast-owned credential root;
// normal provider settings and project instructions are never loaded.
func IsolatedEnvironment(providerName, directory string) ([]string, map[string]string, error) {
	switch providerName {
	case "codex":
		home, err := privateDirectory(directory, "home")
		if err != nil {
			return nil, nil, err
		}
		config := ""
		if hasEnvironmentCredential("OPENAI_API_KEY", "CODEX_API_KEY") {
			config, err = privateDirectory(directory, "codex-home")
		} else {
			config, err = storedAuthenticationDirectory("codex")
		}
		if err != nil {
			return nil, nil, err
		}
		return []string{"CODEX_API_KEY", "LANG", "LC_ALL", "OPENAI_API_KEY", "PATH", "TMPDIR"}, map[string]string{
			"HOME": home, "CODEX_HOME": config,
		}, nil
	case "claude":
		if err := rejectClaudeManagedConfiguration(); err != nil {
			return nil, nil, err
		}
		home, err := privateDirectory(directory, "home")
		if err != nil {
			return nil, nil, err
		}
		config := ""
		if hasEnvironmentCredential("ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN", "CLAUDE_CODE_OAUTH_TOKEN") {
			config, err = privateDirectory(directory, "claude-config")
		} else {
			config, err = storedAuthenticationDirectory("claude")
		}
		if err != nil {
			return nil, nil, err
		}
		return []string{
			"ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN", "CLAUDE_CODE_OAUTH_TOKEN",
			"LANG", "LC_ALL", "PATH", "TMPDIR",
		}, map[string]string{
			"HOME":                                 home,
			"CLAUDE_CONFIG_DIR":                    config,
			"CLAUDE_CODE_DISABLE_AUTO_MEMORY":      "1",
			"CLAUDE_CODE_DISABLE_BACKGROUND_TASKS": "1",
			"CLAUDE_CODE_DISABLE_CLAUDE_MDS":       "1",
			"CLAUDE_CODE_DISABLE_CRON":             "1",
			"CLAUDE_CODE_ENABLE_PROMPT_SUGGESTION": "false",
			"DISABLE_AUTOUPDATER":                  "1",
		}, nil
	default:
		return nil, nil, fmt.Errorf("unsupported agent provider %q", providerName)
	}
}

// FilteredEnvironment materializes an allowlisted child environment. Override
// values win and are sorted so provider plan tests remain deterministic.
func FilteredEnvironment(names []string, overrides map[string]string) []string {
	allowed := map[string]bool{}
	for _, name := range names {
		allowed[name] = true
	}
	result := make([]string, 0, len(names)+len(overrides))
	for _, item := range os.Environ() {
		name, _, found := strings.Cut(item, "=")
		if found && allowed[name] {
			result = append(result, item)
		}
	}
	keys := make([]string, 0, len(overrides))
	for name := range overrides {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	for _, name := range keys {
		result = append(result, name+"="+overrides[name])
	}
	return result
}

func privateDirectory(root, name string) (string, error) {
	path := filepath.Join(root, name)
	if err := os.MkdirAll(path, 0o700); err != nil {
		return "", fmt.Errorf("create isolated provider configuration directory: %w", err)
	}
	if err := os.Chmod(path, 0o700); err != nil {
		return "", fmt.Errorf("protect isolated provider configuration directory: %w", err)
	}
	info, err := os.Lstat(path)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("isolated provider configuration path is not a regular directory")
	}
	return path, nil
}

func hasEnvironmentCredential(names ...string) bool {
	for _, name := range names {
		if strings.TrimSpace(os.Getenv(name)) != "" {
			return true
		}
	}
	return false
}

func rejectClaudeManagedConfiguration() error {
	paths := []string{}
	switch runtime.GOOS {
	case "darwin":
		paths = append(paths,
			"/Library/Application Support/ClaudeCode/managed-settings.json",
			"/Library/Application Support/ClaudeCode/managed-mcp.json",
			"/Library/Application Support/ClaudeCode/managed-settings.d",
			"/Library/Managed Preferences/com.anthropic.claudecode.plist",
			"/Library/Preferences/com.anthropic.claudecode.plist",
		)
	case "linux":
		paths = append(paths,
			"/etc/claude-code/managed-settings.json",
			"/etc/claude-code/managed-mcp.json",
			"/etc/claude-code/managed-settings.d",
		)
	case "windows":
		if programFiles := os.Getenv("ProgramFiles"); programFiles != "" {
			root := filepath.Join(programFiles, "ClaudeCode")
			paths = append(paths, filepath.Join(root, "managed-settings.json"), filepath.Join(root, "managed-mcp.json"), filepath.Join(root, "managed-settings.d"))
		}
	}
	for _, path := range paths {
		// #nosec G703 -- these are read-only, scanner-constructed operating-system
		// policy locations. The only environment-derived root is ProgramFiles on
		// Windows; no file content is read and any positive result fails closed.
		info, err := os.Stat(path)
		if err == nil && (info.Mode().IsRegular() && info.Size() > 0 || info.IsDir() && directoryHasJSON(path)) {
			return fmt.Errorf("claude managed configuration is present at %s; scanner isolation cannot override managed hooks, plugins, or MCP policy", path)
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) && !errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("inspect claude managed configuration path %s: %w", path, err)
		}
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("claude managed configuration path %s cannot be inspected; scanner isolation cannot be established", path)
		}
	}
	return nil
}

func directoryHasJSON(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return true
	}
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") && strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			return true
		}
	}
	return false
}
