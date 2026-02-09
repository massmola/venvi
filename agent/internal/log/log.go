package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single log entry in a session.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Role      string    `json:"role"` // "User", "Agent", "System"
	Content   string    `json:"content"`
}

// Session represents a recording of an agent's task execution.
type Session struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	Entries   []Entry   `json:"entries"`
}

// Logger manages session recording.
type Logger struct {
	DataDir string
}

// NewLogger creates a new logger.
func NewLogger(dataDir string) *Logger {
	return &Logger{DataDir: dataDir}
}

// getSessionPath returns the path to the current session file.
// It validates the sessionID to prevent path traversal.
func (l *Logger) getSessionPath(sessionID string) (string, error) {
	cleanID := filepath.Base(sessionID)
	if cleanID != sessionID {
		return "", fmt.Errorf("invalid session ID: %s", sessionID)
	}
	// Additional strict validation (optional but recommended based on user prompt)
	for _, r := range sessionID {
		isValid := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_'
		if !isValid {
			return "", fmt.Errorf("invalid characters in session ID: %s", sessionID)
		}
	}
	return filepath.Join(l.DataDir, "logs", cleanID+".json"), nil
}

// StartSession initializes a new session log.
func (l *Logger) StartSession(sessionID string) error {
	path, err := l.getSessionPath(sessionID)
	if err != nil {
		return err
	}

	// Check if session already exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("session '%s' already exists", sessionID)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check session existence: %w", err)
	}

	// Create logs directory
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	session := Session{
		ID:        sessionID,
		StartTime: time.Now(),
		Entries:   []Entry{},
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Write (overwrite if exists)
	// Write (overwrite if exists)
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename session file: %w", err)
	}
	return nil
}

// AppendEntry adds a log entry to the specified session.
func (l *Logger) AppendEntry(sessionID, role, content string) error {
	path, err := l.getSessionPath(sessionID)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return fmt.Errorf("failed to parse session file: %w", err)
	}

	entry := Entry{
		Timestamp: time.Now(),
		Role:      role,
		Content:   content,
	}

	session.Entries = append(session.Entries, entry)

	updatedData, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated session: %w", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated session file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename updated session file: %w", err)
	}
	return nil
}

// GetSession returns the content of a session.
func (l *Logger) GetSession(sessionID string) (*Session, error) {
	path, err := l.getSessionPath(sessionID)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to parse session file: %w", err)
	}

	return &session, nil
}
