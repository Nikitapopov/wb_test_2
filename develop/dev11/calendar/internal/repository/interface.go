package repository

import (
	"dev11/calendar/internal/model"
	"time"
)

type IEventRepository interface {
	SaveEvents() error
	Insert(event model.Event) int
	Update(id int, event model.Event) error
	Remove(id int) error
	GetForDay(day time.Time) ([]model.Event, error)
	GetForWeek(day time.Time) ([]model.Event, error)
	GetForMonth(day time.Time) ([]model.Event, error)
}
