package handler

import "net/http"

type IEventHandler interface {
	Register(routes *http.ServeMux)
	// GetById(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Remove(w http.ResponseWriter, r *http.Request)
	GetForDay(w http.ResponseWriter, r *http.Request)
	GetForWeek(w http.ResponseWriter, r *http.Request)
	GetForMonth(w http.ResponseWriter, r *http.Request)
}
