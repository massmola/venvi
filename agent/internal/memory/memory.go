package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Skill represents a learned capability or lesson.
type Skill struct {
	ID        string    `json:"id"`
	Topic     string    `json:"topic"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

// Store manages the persistence of skills.
type Store struct {
	DataDir string
}

// NewStore creates a new memory store.
func NewStore(dataDir string) *Store {
	return &Store{DataDir: dataDir}
}

// getFilePath returns the path to the memory file.
func (s *Store) getFilePath() string {
	return filepath.Join(s.DataDir, "memory.json")
}

// Load retrieves all skills from storage.
func (s *Store) Load() ([]Skill, error) {
	path := s.getFilePath()
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return []Skill{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat memory file: %w", err)
	}

	// #nosec G304
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read memory file: %w", err)
	}

	var skills []Skill
	if len(data) == 0 {
		return []Skill{}, nil
	}

	if err := json.Unmarshal(data, &skills); err != nil {
		return nil, fmt.Errorf("failed to parse memory file: %w", err)
	}

	return skills, nil
}

// Save persists a new skill to storage.
func (s *Store) Save(skill Skill) error {
	skills, err := s.Load()
	if err != nil {
		return err
	}

	// Simple ID generation if not provided
	if skill.ID == "" {
		skill.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if skill.CreatedAt.IsZero() {
		skill.CreatedAt = time.Now()
	}

	skills = append(skills, skill)

	data, err := json.MarshalIndent(skills, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal skills: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(s.DataDir, 0750); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	tmpPath := s.getFilePath() + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write temp memory file: %w", err)
	}

	if err := os.Rename(tmpPath, s.getFilePath()); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to rename memory file: %w", err)
	}

	return nil
}

// Search finds skills matching the query (case-insensitive substring in Topic, Content, or Tags).
func (s *Store) Search(query string) ([]Skill, error) {
	skills, err := s.Load()
	if err != nil {
		return nil, err
	}

	var results []Skill
	query = strings.ToLower(query)

	for _, skill := range skills {
		match := false
		switch {
		case strings.Contains(strings.ToLower(skill.Topic), query):
			match = true
		case strings.Contains(strings.ToLower(skill.Content), query):
			match = true
		default:
			for _, tag := range skill.Tags {
				if strings.Contains(strings.ToLower(tag), query) {
					match = true
					break
				}
			}
		}

		if match {
			results = append(results, skill)
		}
	}

	return results, nil
}
