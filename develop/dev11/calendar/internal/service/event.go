package service

import (
	"dev11/calendar/internal/model"
	"dev11/calendar/internal/repository"
	"errors"
	"log"
	"time"
)

type eventService struct {
	repo repository.IEventRepository
}

func NewEventService(repo repository.IEventRepository) IEventService {
	return &eventService{
		repo: repo,
	}
}

func (s *eventService) SaveEvents() error {
	return s.repo.SaveEvents()
}

func (s *eventService) Insert(dto InsertEventDTO) int {
	event := model.Event{
		UserId:      dto.UserId,
		Date:        dto.Date,
		Description: dto.Description,
	}

	return s.repo.Insert(event)
}

func (s *eventService) Update(id int, dto UpdateEventDTO) error {
	event := model.Event{
		UserId:      dto.UserId,
		Date:        dto.Date,
		Description: dto.Description,
	}

	err := s.repo.Update(id, event)
	if err != nil {
		log.Printf("error while updating event: %v", err)
	}
	return err
}

func (s *eventService) Remove(id int) error {
	err := s.repo.Remove(id)
	if err != nil {
		log.Printf("error while removing event: %v", err)
	}
	return err
}

func (s *eventService) GetForDay(date string) ([]model.Event, error) {
	dateAsTime, err := time.Parse(model.DateLayout, date)
	if err != nil {
		return []model.Event{}, errors.New("incorrect date format")
	}

	events, err := s.repo.GetForDay(dateAsTime)
	if err != nil {
		log.Printf("error while getting events for day: %v", err)
	}

	return events, err
}

func (s *eventService) GetForWeek(date string) ([]model.Event, error) {
	dateAsTime, err := time.Parse(model.DateLayout, date)
	if err != nil {
		return []model.Event{}, errors.New("incorrect date format")
	}

	events, err := s.repo.GetForWeek(dateAsTime)
	if err != nil {
		log.Printf("error while getting events for week: %v", err)
	}

	return events, err
}

func (s *eventService) GetForMonth(date string) ([]model.Event, error) {
	dateAsTime, err := time.Parse(model.DateLayout, date)
	if err != nil {
		return []model.Event{}, errors.New("incorrect date format")
	}

	events, err := s.repo.GetForMonth(dateAsTime)
	if err != nil {
		log.Printf("error while getting events for month: %v", err)
	}

	return events, err
}
