package repository

import (
	"dev11/calendar/internal/model"
	"dev11/calendar/internal/storage"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type eventRepository struct {
	storage storage.IStorage
	events  map[int]model.Event
	mtx     sync.RWMutex
	counter int
}

func NewEventRepository(storage storage.IStorage) (IEventRepository, error) {
	events, err := storage.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get events: %v", err)
	}

	var maxID int
	eventsMap := make(map[int]model.Event)
	for _, event := range events {
		if _, ok := eventsMap[event.ID]; ok {
			return nil, errors.New("incorrect event storage")
		}

		eventsMap[event.ID] = event
		if event.ID > maxID {
			maxID = event.ID
		}
	}

	repo := &eventRepository{
		events:  eventsMap,
		storage: storage,
		mtx:     sync.RWMutex{},
		counter: maxID + 1,
	}

	return repo, nil
}

func (repo *eventRepository) SaveEvents() error {
	eventSlc := make([]model.Event, 0, len(repo.events))
	for _, event := range repo.events {
		eventSlc = append(eventSlc, event)
	}

	err := repo.storage.Save(eventSlc)
	if err != nil {
		return fmt.Errorf("can't save events: %v", err)
	}

	return nil
}

func (repo *eventRepository) Insert(event model.Event) int {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	event.ID = repo.counter
	repo.counter++

	repo.events[event.ID] = event

	return event.ID
}

func (repo *eventRepository) Update(id int, event model.Event) error {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	updatingEvent, ok := repo.events[id]
	if !ok || updatingEvent.RemoveDate != "" {
		return errors.New("event not found")
	}

	repo.events[event.UserId] = model.Event{
		ID:          updatingEvent.ID,
		UserId:      event.UserId,
		Description: event.Description,
		Date:        event.Date,
		RemoveDate:  updatingEvent.RemoveDate,
	}

	return nil
}

func (repo *eventRepository) Remove(id int) error {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	deletingEvent, ok := repo.events[id]
	if !ok || deletingEvent.RemoveDate != "" {
		return errors.New("event not found")
	}

	deletingEvent.RemoveDate = time.Now().Format(model.DateLayout)

	repo.events[id] = deletingEvent

	return nil
}

func (repo *eventRepository) GetForDay(day time.Time) ([]model.Event, error) {
	repo.mtx.RLock()
	defer repo.mtx.RUnlock()

	eventsForDay := []model.Event{}
	for _, event := range repo.events {
		if event.RemoveDate != "" {
			continue
		}

		eventDate, err := time.Parse(model.DateLayout, event.Date)
		if err != nil {
			log.Printf("error while parsing event's date: %v", err)
			continue
		}

		if eventDate == day {
			event = model.Event{
				ID:          event.ID,
				UserId:      event.UserId,
				Date:        event.Date,
				Description: event.Description,
			}
			eventsForDay = append(eventsForDay, event)
		}
	}

	return eventsForDay, nil
}

func (repo *eventRepository) GetForWeek(day time.Time) ([]model.Event, error) {
	repo.mtx.RLock()
	defer repo.mtx.RUnlock()

	eventsForWeek := []model.Event{}

	year, week := day.ISOWeek()

	for _, event := range repo.events {
		if event.RemoveDate != "" {
			continue
		}

		eventDate, err := time.Parse(model.DateLayout, event.Date)
		if err != nil {
			log.Printf("error while parsing event's date: %v", err)
			continue
		}

		eventYear, eventWeek := eventDate.ISOWeek()

		if eventYear == year && eventWeek == week {
			eventsForWeek = append(eventsForWeek, event)
		}
	}

	return eventsForWeek, nil
}

func (repo *eventRepository) GetForMonth(day time.Time) ([]model.Event, error) {
	repo.mtx.RLock()
	defer repo.mtx.RUnlock()

	eventsForMonth := []model.Event{}

	month := day.Month()

	for _, event := range repo.events {
		if event.RemoveDate != "" {
			continue
		}

		eventDate, err := time.Parse(model.DateLayout, event.Date)
		if err != nil {
			log.Printf("error while parsing event's date: %v", err)
			continue
		}

		if month == eventDate.Month() {
			eventsForMonth = append(eventsForMonth, event)
		}
	}

	return eventsForMonth, nil
}
