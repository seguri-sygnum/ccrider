package session

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ResolveWorkingDir determines the correct directory to start claude --resume
//
// CRITICAL: Always returns projectPath, NOT lastCwd.
// Claude --resume only finds sessions stored in the project directory.
// The resume prompt tells Claude where the session last was (lastCwd).
//
// DO NOT CHANGE THIS - see commits db2bc33 and 33050ea
func ResolveWorkingDir(projectPath, lastCwd string) string {
	return projectPath
}

// ValidateClaudeRunnable checks if 'claude' command can run from the given directory.
// This catches common issues like version manager (asdf/mise) failures before launching.
// Returns nil if claude is runnable, or a descriptive error with suggested fixes.
func ValidateClaudeRunnable(workDir string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	// Run 'claude --version' through a login shell to simulate actual launch conditions
	cmd := exec.Command(shell, "-l", "-c", "claude --version")
	if workDir != "" {
		cmd.Dir = workDir
	}
	cmd.Env = os.Environ()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return nil
	}

	// Analyze the failure and provide helpful diagnostics
	stderrStr := stderr.String()

	// Check for version manager issues (asdf, mise, nvm, etc.)
	if strings.Contains(stderrStr, "No version is set for command") ||
		strings.Contains(stderrStr, "not installed") ||
		strings.Contains(stderrStr, ".tool-versions") {
		return fmt.Errorf(`claude command failed in %s

%s
This is typically caused by a version manager (asdf/mise/nvm) configuration issue.

Suggested fixes:
  1. Install the required nodejs version: asdf install nodejs <version>
  2. Or set a global nodejs version: asdf global nodejs <version>
  3. Or remove/update the .tool-versions file in the project`, workDir, strings.TrimSpace(stderrStr))
	}

	// Check for command not found
	if strings.Contains(stderrStr, "command not found") ||
		strings.Contains(stderrStr, "not found") {
		return fmt.Errorf(`claude command not found in %s

%s
Suggested fixes:
  1. Install Claude Code: npm install -g @anthropic-ai/claude-code
  2. Ensure your PATH includes the claude binary
  3. Check that your shell profile loads correctly`, workDir, strings.TrimSpace(stderrStr))
	}

	// Generic failure
	if stderrStr != "" {
		return fmt.Errorf("claude command failed in %s:\n%s", workDir, strings.TrimSpace(stderrStr))
	}

	return fmt.Errorf("claude command failed in %s: %w", workDir, err)
}

// SessionFileExists checks if the Claude Code session file exists on disk.
// Claude stores sessions in ~/.claude/projects/<encoded-project-path>/<session-id>.jsonl
func SessionFileExists(sessionID, projectPath string) bool {
	filePath := GetSessionFilePath(sessionID, projectPath)
	if filePath == "" {
		return false
	}
	_, err := os.Stat(filePath)
	return err == nil
}

// GetSessionFilePath returns the expected path to a Claude Code session file.
// Returns empty string if home directory cannot be determined.
func GetSessionFilePath(sessionID, projectPath string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Claude encodes project paths by replacing / with -
	// e.g., /Users/neil/enaia/enaia -> -Users-neil-enaia-enaia
	encodedPath := strings.ReplaceAll(projectPath, "/", "-")

	return filepath.Join(home, ".claude", "projects", encodedPath, sessionID+".jsonl")
}
