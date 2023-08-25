package calendar

import (
	"dev11/calendar/internal/config"
	"dev11/calendar/internal/handler"
	"dev11/calendar/internal/repository"
	"dev11/calendar/internal/service"
	"dev11/calendar/internal/storage"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	conf := config.GetConfig()
	jsonStorage := storage.NewJsonStorage(conf.StorageFilePath)

	repo, err := repository.NewEventRepository(jsonStorage)
	if err != nil {
		log.Printf("error while init event repository: %v", err)
		panic(err)
	}
	service := service.NewEventService(repo)
	defer func() {
		err = service.SaveEvents()
		if err != nil {
			log.Printf("error while saving events: %v", err)
		}
	}()

	handler := handler.NewEventHandler(service)

	mux := http.NewServeMux()
	handler.Register(mux)

	srv := &http.Server{
		Addr:    net.JoinHostPort(conf.Host, conf.Port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("error while init event repository: %v", err)
			panic(err)
		}
	}()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-osSignal
}
