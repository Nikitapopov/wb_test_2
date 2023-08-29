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

// Репозиторий событий
type eventRepository struct {
	storage storage.IStorage
	events  map[int]model.Event
	mtx     sync.RWMutex
	counter int
}

// Конструктор репозитория событий
func NewEventRepository(storage storage.IStorage) (IEventRepository, error) {
	// Получение событий
	events, err := storage.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get events: %v", err)
	}

	// Наибольший ID
	var maxID int

	// Мапа событий
	eventsMap := make(map[int]model.Event, len(events))

	// Заполнение мапы
	for _, event := range events {
		// Если найдено событие c дублирующимся ID, то возврат ошибки
		if _, ok := eventsMap[event.ID]; ok {
			return nil, errors.New("incorrect event storage")
		}

		// Добавление события в мапу
		eventsMap[event.ID] = event

		// Обновление максимального ID
		if event.ID > maxID {
			maxID = event.ID
		}
	}

	// Создание объекта репозитория
	repo := &eventRepository{
		events:  eventsMap,
		storage: storage,
		mtx:     sync.RWMutex{},
		counter: maxID + 1,
	}

	return repo, nil
}

// Сохранение событий в хранилище
func (repo *eventRepository) SaveEvents() error {
	// Преобразование мапы событий в слайс
	eventSlc := make([]model.Event, 0, len(repo.events))
	for _, event := range repo.events {
		eventSlc = append(eventSlc, event)
	}

	// Сохранение событий в хранилище
	err := repo.storage.Save(eventSlc)
	if err != nil {
		return fmt.Errorf("can't save events: %v", err)
	}

	return nil
}

// Добавление события
func (repo *eventRepository) Insert(event model.Event) int {
	// Использование мьютекса для избежания гонки данных
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	// Использование счетчика для назначения ID
	event.ID = repo.counter
	repo.counter++

	// Добавление события в локальную мапу
	repo.events[event.ID] = event

	return event.ID
}

// Обновление события
func (repo *eventRepository) Update(id int, event model.Event) error {
	// Использование мьютекса для избежания гонки данных
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	// Поиск неудаленного события
	updatingEvent, ok := repo.events[id]
	if !ok || updatingEvent.RemoveDate != "" {
		return errors.New("event not found")
	}

	// Обновление события
	repo.events[event.UserId] = model.Event{
		ID:          updatingEvent.ID,
		UserId:      event.UserId,
		Description: event.Description,
		Date:        event.Date,
		RemoveDate:  updatingEvent.RemoveDate,
	}

	return nil
}

// Удаление события
func (repo *eventRepository) Remove(id int) error {
	// Использование мьютекса для избежания гонки данных
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	// Поиск неудаленного события
	deletingEvent, ok := repo.events[id]
	if !ok || deletingEvent.RemoveDate != "" {
		return errors.New("event not found")
	}

	// Назначение удаленной даты
	deletingEvent.RemoveDate = time.Now().Format(model.DateLayout)
	repo.events[id] = deletingEvent

	return nil
}

// Получение событий за дату date
func (repo *eventRepository) GetForDay(date time.Time) ([]model.Event, error) {
	// Использование мьютекса для избежания гонки данных
	repo.mtx.RLock()
	defer repo.mtx.RUnlock()

	// События за день
	eventsForDay := []model.Event{}

	// Итерация по событиям
	for _, event := range repo.events {
		// Если событие удалено, то продолжаем итерацию
		if event.RemoveDate != "" {
			continue
		}

		// Парсинг даты к формату
		eventDate, err := time.Parse(model.DateLayout, event.Date)
		if err != nil {
			log.Printf("error while parsing event's date: %v", err)
			continue
		}

		// Если дата события совпадает, то выполняется добавление в результат
		if eventDate == date {
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

// Получение событий за неделю, в которой имеется дата date
func (repo *eventRepository) GetForWeek(date time.Time) ([]model.Event, error) {
	// Использование мьютекса для избежания гонки данных
	repo.mtx.RLock()
	defer repo.mtx.RUnlock()

	// События за неделю
	eventsForWeek := []model.Event{}

	// Распарсинг даты date на год и неделю
	year, week := date.ISOWeek()

	// Итерация по событиям
	for _, event := range repo.events {
		// Если событие удалено, то продолжаем итерацию
		if event.RemoveDate != "" {
			continue
		}

		// Парсинг даты к формату
		eventDate, err := time.Parse(model.DateLayout, event.Date)
		if err != nil {
			log.Printf("error while parsing event's date: %v", err)
			continue
		}

		// Распарсинг даты события на год и неделю
		eventYear, eventWeek := eventDate.ISOWeek()

		// Если год и неделя совпадают, то выполняется добавление в результат
		if eventYear == year && eventWeek == week {
			eventsForWeek = append(eventsForWeek, event)
		}
	}

	return eventsForWeek, nil
}

// Получение событий за месяц, в котором имеется дата date
func (repo *eventRepository) GetForMonth(date time.Time) ([]model.Event, error) {
	// Использование мьютекса для избежания гонки данных
	repo.mtx.RLock()
	defer repo.mtx.RUnlock()

	// События за месяц
	eventsForMonth := []model.Event{}

	// Парсинг месяца из даты
	month := date.Month()

	// Итерация по событиям
	for _, event := range repo.events {
		// Если событие удалено, то продолжаем итерацию
		if event.RemoveDate != "" {
			continue
		}

		// Парсинг даты к формату
		eventDate, err := time.Parse(model.DateLayout, event.Date)
		if err != nil {
			log.Printf("error while parsing event's date: %v", err)
			continue
		}

		// Если месяц совпадает, то выполняется добавление в результат
		if month == eventDate.Month() {
			eventsForMonth = append(eventsForMonth, event)
		}
	}

	return eventsForMonth, nil
}
