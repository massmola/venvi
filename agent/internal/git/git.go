package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// HasChanges checks if there are any uncommitted changes in the repo.
func HasChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// StageAll stages all changes (git add .).
func StageAll() error {
	cmd := exec.Command("git", "add", ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %s: %w", string(output), err)
	}
	return nil
}

// Commit commits the staged changes with the given message.
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", "agent: "+message)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %s: %w", string(output), err)
	}
	return nil
}
