package service

import (
	"dev11/calendar/internal/model"
)

type IEventService interface {
	// GetById(id int) (model.Event, error)
	SaveEvents() error
	Insert(dto InsertEventDTO) int
	Update(id int, dto UpdateEventDTO) error
	Remove(id int) error
	GetForDay(day string) ([]model.Event, error)
	GetForWeek(day string) ([]model.Event, error)
	GetForMonth(day string) ([]model.Event, error)
}
