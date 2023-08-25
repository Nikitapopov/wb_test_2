package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"dev11/calendar/internal/model"
)

type jsonStorage struct {
	fileName string
}

func NewJsonStorage(fileName string) IStorage {
	return &jsonStorage{
		fileName: fileName,
	}
}

func (s *jsonStorage) Get() ([]model.Event, error) {
	file, err := s.openStorageFile()
	if err != nil {
		return nil, fmt.Errorf("getting opening file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if _, err := decoder.Token(); err != nil {
		return []model.Event{}, nil
	}

	events := []model.Event{}
	for decoder.More() {
		var event model.Event

		if err := decoder.Decode(&event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *jsonStorage) Save(events []model.Event) error {
	file, err := s.openStorageFile()
	if err != nil {
		return fmt.Errorf("getting opening file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(events)
	if err != nil {
		return fmt.Errorf("encoding events to file: %w", err)
	}

	return nil
}

// Функция для открытия файла. Требует закрытия возвращаего значения!
func (s *jsonStorage) openStorageFile() (*os.File, error) {
	dir := filepath.Dir(s.fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return nil, fmt.Errorf("making dir: %w", err)
		}
	}

	file, err := os.OpenFile(s.fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return file, nil
}
