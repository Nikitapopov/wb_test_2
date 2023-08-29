package handler

import (
	"dev11/calendar/internal/middleware"
	"dev11/calendar/internal/model"
	"dev11/calendar/internal/service"
	"dev11/calendar/pkg/api_helper"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Хэндлер событий
type eventHandler struct {
	eventService service.IEventService
}

// Конструктор хэндлера событий
func NewEventHandler(eventService service.IEventService) IEventHandler {
	return &eventHandler{
		eventService: eventService,
	}
}

// Регистрация конкретных обработчиков в роутере router
func (h *eventHandler) Register(router *http.ServeMux) {
	router.Handle("/create_event", middleware.Log(http.HandlerFunc(h.Insert)))
	router.Handle("/update_event", middleware.Log(http.HandlerFunc(h.Update)))
	router.Handle("/delete_event", middleware.Log(http.HandlerFunc(h.Remove)))
	router.Handle("/events_for_day", middleware.Log(http.HandlerFunc(h.GetForDay)))
	router.Handle("/events_for_week", middleware.Log(http.HandlerFunc(h.GetForWeek)))
	router.Handle("/events_for_month", middleware.Log(http.HandlerFunc(h.GetForMonth)))
}

// Добавление события
func (h *eventHandler) Insert(w http.ResponseWriter, r *http.Request) {
	// Обработка несоответствия метода запроса
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	// Десериализация параметров
	var dto service.InsertEventDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Валидация параметров
	err = validateInsertDto(dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Вставка события
	id := h.eventService.Insert(dto)

	// Возвращаемое значение
	var payload api_helper.JsonResponse
	payload.Result = struct {
		Id int `json:"id"`
	}{Id: id}

	// Оформление ответа
	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Обработка несоответствия метода запроса
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	// Десериализация параметров
	var dto service.UpdateEventDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Валидация параметров
	err = validateUpdateDto(dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Обновление события
	err = h.eventService.Update(dto.UserId, dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	// Возвращаемое значение
	var payload api_helper.JsonResponse
	payload.Result = "ok"

	// Оформление ответа
	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) Remove(w http.ResponseWriter, r *http.Request) {
	// Обработка несоответствия метода запроса
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	// Десериализация параметров
	var dto service.RemoveEventDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Валидация параметров
	err = validateRemoveDto(dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Удаление события
	err = h.eventService.Remove(dto.ID)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	// Возвращаемое значение
	var payload api_helper.JsonResponse
	payload.Result = "ok"

	// Оформление ответа
	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) GetForDay(w http.ResponseWriter, r *http.Request) {
	// Обработка несоответствия метода запроса
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Получение параметра date
	date := r.URL.Query().Get("date")
	if date == "" {
		api_helper.ErrorJSON(w, errors.New("date parameter is not defined"), http.StatusBadRequest)
		return
	}

	// Валидация даты
	err := validateDate(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Получение событий за дату date
	events, err := h.eventService.GetForDay(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	// Возвращаемое значение
	var payload api_helper.JsonResponse
	payload.Result = struct {
		Events []model.Event `json:"events"`
	}{Events: events}

	// Оформление ответа
	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) GetForWeek(w http.ResponseWriter, r *http.Request) {
	// Обработка несоответствия метода запроса
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Получение параметра date
	date := r.URL.Query().Get("date")
	if date == "" {
		api_helper.ErrorJSON(w, errors.New("date parameter is not defined"), http.StatusBadRequest)
		return
	}

	// Валидация даты
	err := validateDate(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Получение событий за неделю, в которой имеется дата date
	events, err := h.eventService.GetForWeek(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	// Возвращаемое значение
	var payload api_helper.JsonResponse
	payload.Result = struct {
		Events []model.Event `json:"events"`
	}{Events: events}

	// Оформление ответа
	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) GetForMonth(w http.ResponseWriter, r *http.Request) {
	// Обработка несоответствия метода запроса
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Получение параметра date
	date := r.URL.Query().Get("date")
	if date == "" {
		api_helper.ErrorJSON(w, errors.New("date parameter is not defined"), http.StatusBadRequest)
		return
	}

	// Валидация даты
	err := validateDate(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Получение событий за месяц, в котором имеется дата date
	events, err := h.eventService.GetForMonth(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	// Возвращаемое значение
	var payload api_helper.JsonResponse
	payload.Result = struct {
		Events []model.Event `json:"events"`
	}{Events: events}

	// Оформление ответа
	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

// Валидация параметров для вставки события
func validateInsertDto(dto service.InsertEventDTO) error {
	if dto.UserId < 0 {
		return errors.New("user id parameter should be positive")
	}

	if len(dto.Description) == 0 {
		return errors.New("description parameter should not be empty")
	}

	if err := validateDate(dto.Date); err != nil {
		return err
	}

	return nil
}

// Валидация параметров для обновления события
func validateUpdateDto(dto service.UpdateEventDTO) error {
	if dto.ID < 0 {
		return errors.New("id parameter should be positive")
	}
	if dto.UserId < 0 {
		return errors.New("user id parameter should be positive")
	}

	if len(dto.Description) == 0 {
		return errors.New("description parameter should not be empty")
	}

	if err := validateDate(dto.Date); err != nil {
		return err
	}

	return nil
}

// Валидация параметров для удаления события
func validateRemoveDto(dto service.RemoveEventDTO) error {
	if dto.ID < 0 {
		return errors.New("id parameter should be positive")
	}
	return nil
}

// Валидация даты
func validateDate(date string) error {
	if _, err := time.Parse(model.DateLayout, date); err != nil {
		return errors.New("date parameter should be in format 2006-01-02")
	}

	return nil
}
