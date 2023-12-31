package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"dev11/calendar/internal/model"
)

// Хранилище событий
type eventStorage struct {
	fileName string
}

// Конструктор хранилища событий
func NewEventStorage(fileName string) IStorage {
	return &eventStorage{
		fileName: fileName,
	}
}

// Получение событий из хранилища
func (s *eventStorage) Get() ([]model.Event, error) {
	// Открытие файла с событиями
	file, err := s.openStorageFile()
	if err != nil {
		return nil, fmt.Errorf("getting opening file: %w", err)
	}
	defer file.Close()

	// Валидация файла на корректность json-а и возврат пустого слайса, если файл не корректен
	decoder := json.NewDecoder(file)
	if _, err := decoder.Token(); err != nil {
		return []model.Event{}, nil
	}

	// Итерационная запись событий в events
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

// Сохранение событий в хранилище
func (s *eventStorage) Save(events []model.Event) error {
	// Открытие файла с событиями
	file, err := s.openStorageFile()
	if err != nil {
		return fmt.Errorf("getting opening file: %w", err)
	}
	defer file.Close()

	// Сериализация событий и запись в файл
	encoder := json.NewEncoder(file)
	err = encoder.Encode(events)
	if err != nil {
		return fmt.Errorf("encoding events to file: %w", err)
	}

	return nil
}

// Функция для открытия файла. Требует закрытия возвращаего значения!
// Если файл не существует, то он создается
func (s *eventStorage) openStorageFile() (*os.File, error) {
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
