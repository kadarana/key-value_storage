package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
)

type storageFile struct {
	Inner map[string]Value `json:"inner"`
	List  map[string][]any `json:"list"`
}

func (s *Storage) SaveToFile(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := storageFile{
		Inner: s.inner,
		List:  make(map[string][]any),
	}

	for k, v := range s.list {
		data.List[k] = v.Elem
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	s.logger.Info("Storage saved to file", zap.String("file", path))
	return nil
}

func (s *Storage) LoadFromFile(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %w", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var data storageFile
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	s.inner = data.Inner
	s.list = make(map[string]*List)

	for k, v := range data.List {
		s.list[k] = &List{Elem: v}
	}

	s.logger.Info("Storage loaded from file",
		zap.String("file", path),
		zap.Int("items_loaded", len(s.inner)+len(s.list)))
	return nil
}
