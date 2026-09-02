package main

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

const maximumAuthenticationOutput = 256 * 1024

type authenticationOutput struct {
	data     []byte
	exceeded bool
}

func (output *authenticationOutput) Write(value []byte) (int, error) {
	remaining := maximumAuthenticationOutput - len(output.data)
	if remaining > 0 {
		count := len(value)
		if count > remaining {
			count = remaining
		}
		output.data = append(output.data, value[:count]...)
	}
	if len(value) > remaining {
		output.exceeded = true
	}
	return len(value), nil
}

func (output *authenticationOutput) String() string {
	return string(output.data)
}

var authenticationEnvironmentNames = []string{
	"APPDATA", "BROWSER", "DBUS_SESSION_BUS_ADDRESS", "DISPLAY", "HOME",
	"HOMEDRIVE", "HOMEPATH", "HTTP_PROXY", "HTTPS_PROXY", "LANG", "LC_ALL",
	"LOCALAPPDATA", "NO_PROXY", "PATH", "SHELL", "TEMP", "TERM", "TMP",
	"TMPDIR", "USERPROFILE", "WAYLAND_DISPLAY", "XDG_RUNTIME_DIR",
	"http_proxy", "https_proxy", "no_proxy",
}

func runAuthentication(operation string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		printAuthenticationUsage(operation, stdout)
		return nil
	}
	if operation == "auth" {
		if len(args) > 1 {
			return exitError(exitConfiguration, errors.New("auth accepts at most one provider: codex or claude"))
		}
		providers := []string{"codex", "claude"}
		if len(args) == 1 {
			providers = []string{args[0]}
		}
		for _, providerName := range providers {
			if err := validateAuthenticationProvider(providerName); err != nil {
				return exitError(exitConfiguration, err)
			}
			if err := printAuthenticationStatus(providerName, stdout); err != nil {
				return exitError(exitExecution, err)
			}
		}
		return nil
	}
	if len(args) != 1 {
		return exitError(exitConfiguration, fmt.Errorf("%s requires one provider: codex or claude", operation))
	}
	providerName := args[0]
	if err := validateAuthenticationProvider(providerName); err != nil {
		return exitError(exitConfiguration, err)
	}
	path, err := authenticationExecutable(providerName)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	overrides, err := provider.AuthenticationOverrides(providerName)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	arguments := authenticationArguments(providerName, operation)
	command := exec.Command(path, arguments...)
	command.Env = provider.FilteredEnvironment(authenticationEnvironmentNames, overrides)
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	if operation == "login" {
		printAuthenticationStart(stdout, newTerminalStyle("auto", stdout), providerName)
	}
	if err := command.Run(); err != nil {
		return exitError(exitExecution, fmt.Errorf("%s %s failed: %w", providerName, operation, err))
	}
	if operation == "login" {
		if err := provider.MarkStoredAuthentication(providerName); err != nil {
			return exitError(exitInternal, err)
		}
		printAuthenticationSuccess(stdout, newTerminalStyle("auto", stdout), providerName)
	} else if operation == "logout" {
		if err := provider.ClearStoredAuthentication(providerName); err != nil {
			return exitError(exitInternal, err)
		}
		fmt.Fprintf(stdout, "The scanner's %s login has been cleared.\n", authenticationProviderTitle(providerName))
	}
	return nil
}

func printAuthenticationUsage(operation string, output io.Writer) {
	switch operation {
	case "login":
		fmt.Fprintln(output, "Usage: prc login <codex|claude>")
		fmt.Fprintln(output, "Opens the provider's official sign-in flow in a private scanner credential folder.")
	case "logout":
		fmt.Fprintln(output, "Usage: prc logout <codex|claude>")
		fmt.Fprintln(output, "Removes the provider login stored for scanner runs.")
	case "auth":
		fmt.Fprintln(output, "Usage: prc auth [codex|claude]")
		fmt.Fprintln(output, "Shows whether the scanner can use each provider login.")
	}
}

func printAuthenticationStatus(providerName string, output io.Writer) error {
	path, err := authenticationExecutable(providerName)
	if err != nil {
		fmt.Fprintf(output, "%-11s not installed\n", authenticationProviderTitle(providerName)+":")
		return nil
	}
	overrides, err := provider.AuthenticationOverrides(providerName)
	if err != nil {
		return err
	}
	command := exec.Command(path, authenticationArguments(providerName, "auth")...)
	command.Env = provider.FilteredEnvironment(authenticationEnvironmentNames, overrides)
	var childOutput authenticationOutput
	command.Stdout = &childOutput
	command.Stderr = &childOutput
	err = command.Run()
	if childOutput.exceeded {
		return fmt.Errorf("%s authentication status exceeded its output limit", providerName)
	}
	message := terminalText(strings.TrimSpace(childOutput.String()))
	if message == "" {
		message = "status unavailable"
	}
	if err == nil {
		if markErr := provider.MarkStoredAuthentication(providerName); markErr != nil {
			return markErr
		}
		fmt.Fprintf(output, "%-11s %s\n", authenticationProviderTitle(providerName)+":", message)
		return nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		if clearErr := provider.ClearStoredAuthentication(providerName); clearErr != nil {
			return clearErr
		}
		fmt.Fprintf(output, "%-11s %s\n", authenticationProviderTitle(providerName)+":", message)
		return nil
	}
	return fmt.Errorf("check %s authentication: %w", providerName, err)
}

func authenticationExecutable(providerName string) (string, error) {
	path, err := exec.LookPath(providerName)
	if err != nil {
		return "", fmt.Errorf("%s CLI is not installed or is not on PATH", authenticationProviderTitle(providerName))
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve %s CLI: %w", providerName, err)
	}
	name := strings.TrimSuffix(strings.ToLower(filepath.Base(absolute)), ".exe")
	if name != providerName {
		return "", fmt.Errorf("%s CLI executable must be named %s", authenticationProviderTitle(providerName), providerName)
	}
	return absolute, nil
}

func authenticationArguments(providerName, operation string) []string {
	switch providerName {
	case "codex":
		switch operation {
		case "login":
			return []string{"login", "-c", `cli_auth_credentials_store="file"`}
		case "logout":
			return []string{"logout", "-c", `cli_auth_credentials_store="file"`}
		case "auth":
			return []string{"login", "status", "-c", `cli_auth_credentials_store="file"`}
		}
	case "claude":
		switch operation {
		case "login":
			return []string{"auth", "login"}
		case "logout":
			return []string{"auth", "logout"}
		case "auth":
			return []string{"auth", "status", "--text"}
		}
	}
	panic("authentication provider and operation were validated")
}

func validateAuthenticationProvider(providerName string) error {
	if providerName != "codex" && providerName != "claude" {
		return fmt.Errorf("unsupported provider %q; use codex or claude", providerName)
	}
	return nil
}

func authenticationProviderTitle(providerName string) string {
	if providerName == "codex" {
		return "Codex"
	}
	return "Claude Code"
}
