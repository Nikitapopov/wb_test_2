package storage

import "dev11/calendar/internal/model"

type IStorage interface {
	Get() ([]model.Event, error)
	Save([]model.Event) error
}
