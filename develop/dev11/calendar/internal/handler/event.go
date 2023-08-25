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

type eventHandler struct {
	eventService service.IEventService
}

func NewEventHandler(eventService service.IEventService) IEventHandler {
	return &eventHandler{
		eventService: eventService,
	}
}

func (h *eventHandler) Register(router *http.ServeMux) {
	router.Handle("/create_event", middleware.Log(http.HandlerFunc(h.Insert)))
	router.Handle("/update_event", middleware.Log(http.HandlerFunc(h.Update)))
	router.Handle("/delete_event", middleware.Log(http.HandlerFunc(h.Remove)))
	router.Handle("/events_for_day", middleware.Log(http.HandlerFunc(h.GetForDay)))
	router.Handle("/events_for_week", middleware.Log(http.HandlerFunc(h.GetForWeek)))
	router.Handle("/events_for_month", middleware.Log(http.HandlerFunc(h.GetForMonth)))
}

func (h *eventHandler) Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var dto service.InsertEventDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = validateInsertDto(dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	id := h.eventService.Insert(dto)

	var payload api_helper.JsonResponse
	payload.Result = struct {
		Id int `json:"id"`
	}{Id: id}

	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var dto service.UpdateEventDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = validateUpdateDto(dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = h.eventService.Update(dto.UserId, dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	var payload api_helper.JsonResponse
	payload.Result = "ok"

	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) Remove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var dto service.RemoveEventDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = validateRemoveDto(dto)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = h.eventService.Remove(dto.ID)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	var payload api_helper.JsonResponse
	payload.Result = "ok"

	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) GetForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		api_helper.ErrorJSON(w, errors.New("date parameter is not defined"), http.StatusBadRequest)
		return
	}

	err := validateDate(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	events, err := h.eventService.GetForDay(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	var payload api_helper.JsonResponse
	payload.Result = struct {
		Events []model.Event `json:"events"`
	}{Events: events}

	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) GetForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		api_helper.ErrorJSON(w, errors.New("date parameter is not defined"), http.StatusBadRequest)
		return
	}

	err := validateDate(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	events, err := h.eventService.GetForWeek(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	var payload api_helper.JsonResponse
	payload.Result = struct {
		Events []model.Event `json:"events"`
	}{Events: events}

	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *eventHandler) GetForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		api_helper.ErrorJSON(w, errors.New("date parameter is not defined"), http.StatusBadRequest)
		return
	}

	err := validateDate(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	events, err := h.eventService.GetForMonth(date)
	if err != nil {
		api_helper.ErrorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	var payload api_helper.JsonResponse
	payload.Result = struct {
		Events []model.Event `json:"events"`
	}{Events: events}

	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func validateInsertDto(dto service.InsertEventDTO) error {
	if dto.UserId < 0 {
		return errors.New("user id parameter should be positive")
	}

	if len(dto.Description) == 0 {
		return errors.New("description parameter should not be empty")
	}

	if _, err := time.Parse(model.DateLayout, dto.Date); err != nil {
		return errors.New("date parameter should be in format 2006-01-02")
	}

	return nil
}

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

	if _, err := time.Parse(model.DateLayout, dto.Date); err != nil {
		return errors.New("date parameter should be in format 2006-01-02")
	}

	return nil
}

func validateRemoveDto(dto service.RemoveEventDTO) error {
	if dto.ID < 0 {
		return errors.New("id parameter should be positive")
	}
	return nil
}

func validateDate(date string) error {
	_, err := time.Parse(model.DateLayout, date)
	if err != nil {
		return errors.New("incorrect date format")
	}

	return nil
}
