package provider

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const authenticationMarker = ".prc-authenticated"

var locateUserConfigDirectory = os.UserConfigDir

// AuthenticationDirectory returns the scanner's provider-specific credential root.
// It is deliberately separate from the provider's normal user configuration,
// so a scan can reuse login credentials without loading normal hooks, plugins,
// MCP servers, instructions, or session history.
func AuthenticationDirectory(providerName string) (string, error) {
	if providerName != "codex" && providerName != "claude" {
		return "", fmt.Errorf("unsupported agent provider %q", providerName)
	}
	root, err := locateUserConfigDirectory()
	if err != nil {
		return "", fmt.Errorf("locate scanner provider authentication directory: %w", err)
	}
	if root == "" {
		return "", fmt.Errorf("locate scanner provider authentication directory: empty user configuration path")
	}
	return filepath.Join(root, "prc", "provider-auth", providerName), nil
}

// PrepareAuthenticationDirectory creates a private provider credential root.
func PrepareAuthenticationDirectory(providerName string) (string, error) {
	path, err := AuthenticationDirectory(providerName)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(path, 0o700); err != nil {
		return "", fmt.Errorf("create %s authentication directory: %w", providerName, err)
	}
	if err := os.Chmod(path, 0o700); err != nil {
		return "", fmt.Errorf("protect %s authentication directory: %w", providerName, err)
	}
	information, err := os.Lstat(path)
	if err != nil || !information.IsDir() || information.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("%s authentication path is not a regular directory", providerName)
	}
	if runtime.GOOS != "windows" && information.Mode().Perm()&0o077 != 0 {
		return "", fmt.Errorf("%s authentication directory is accessible by other users", providerName)
	}
	return path, nil
}

// MarkStoredAuthentication records only that the provider's own login command
// succeeded. The marker contains no credential or account information.
func MarkStoredAuthentication(providerName string) error {
	path, err := PrepareAuthenticationDirectory(providerName)
	if err != nil {
		return err
	}
	marker := filepath.Join(path, authenticationMarker)
	if information, statErr := os.Lstat(marker); statErr == nil &&
		(!information.Mode().IsRegular() || information.Mode()&os.ModeSymlink != 0) {
		return fmt.Errorf("%s authentication marker is not a regular file", providerName)
	} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
		return fmt.Errorf("inspect %s authentication marker: %w", providerName, statErr)
	}
	if err := os.WriteFile(marker, []byte("prc-provider-login-v1\n"), 0o600); err != nil {
		return fmt.Errorf("record %s authentication: %w", providerName, err)
	}
	if err := os.Chmod(marker, 0o600); err != nil {
		return fmt.Errorf("protect %s authentication marker: %w", providerName, err)
	}
	return nil
}

// ClearStoredAuthentication removes the scanner's non-secret login marker after the
// provider's own logout command has removed its credentials.
func ClearStoredAuthentication(providerName string) error {
	path, err := AuthenticationDirectory(providerName)
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(path, authenticationMarker))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("clear %s authentication marker: %w", providerName, err)
	}
	return nil
}

func storedAuthenticationDirectory(providerName string) (string, error) {
	path, err := AuthenticationDirectory(providerName)
	if err != nil {
		return "", err
	}
	information, err := os.Lstat(path)
	if err != nil || !information.IsDir() || information.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("no private scanner %s login was found; run `prc login %s`", providerName, providerName)
	}
	if runtime.GOOS != "windows" && information.Mode().Perm()&0o077 != 0 {
		return "", fmt.Errorf("the scanner's %s login directory is accessible by other users", providerName)
	}
	marker, err := os.Lstat(filepath.Join(path, authenticationMarker))
	if err != nil || !marker.Mode().IsRegular() || marker.Mode()&os.ModeSymlink != 0 || marker.Size() > 128 {
		return "", fmt.Errorf("no private scanner %s login was found; run `prc login %s`", providerName, providerName)
	}
	if runtime.GOOS != "windows" && marker.Mode().Perm()&0o077 != 0 {
		return "", fmt.Errorf("the scanner's %s login marker is accessible by other users", providerName)
	}
	return path, nil
}

// AuthenticationOverrides binds a provider's own login/status/logout command
// to the scanner's private credential root. Normal HOME is left alone for the browser-
// based login UI; scans replace HOME separately.
func AuthenticationOverrides(providerName string) (map[string]string, error) {
	path, err := PrepareAuthenticationDirectory(providerName)
	if err != nil {
		return nil, err
	}
	switch providerName {
	case "codex":
		return map[string]string{"CODEX_HOME": path}, nil
	case "claude":
		return map[string]string{
			"CLAUDE_CONFIG_DIR":                    path,
			"CLAUDE_CODE_DISABLE_AUTO_MEMORY":      "1",
			"CLAUDE_CODE_DISABLE_BACKGROUND_TASKS": "1",
			"CLAUDE_CODE_DISABLE_CLAUDE_MDS":       "1",
			"CLAUDE_CODE_DISABLE_CRON":             "1",
			"CLAUDE_CODE_ENABLE_PROMPT_SUGGESTION": "false",
			"DISABLE_AUTOUPDATER":                  "1",
		}, nil
	default:
		panic("provider name was validated")
	}
}
