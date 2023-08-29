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
	// Получение конфигураций
	conf := config.GetConfig()

	// Хранилище событий
	storage := storage.NewEventStorage(conf.StorageFilePath)

	// Репозиторий событий
	repo, err := repository.NewEventRepository(storage)
	if err != nil {
		log.Printf("error while init event repository: %v", err)
		panic(err)
	}

	// Сервис событий (бизнес логика)
	service := service.NewEventService(repo)

	// Перед выходом из программы выполняется сохранение событий в хранилище
	defer func() {
		err = service.SaveEvents()
		if err != nil {
			log.Printf("error while saving events: %v", err)
		}
	}()

	// Хэндлер событий
	handler := handler.NewEventHandler(service)

	// Роутер сервера
	mux := http.NewServeMux()

	// Регистрация методов событий в роутере
	handler.Register(mux)

	// Сервер
	srv := &http.Server{
		Addr:    net.JoinHostPort(conf.Host, conf.Port),
		Handler: mux,
	}

	// Запуск сервера
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("error while init event repository: %v", err)
			panic(err)
		}
	}()

	// Отлов сигналов об окончании работы
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// При поступлении сигнала программа завершается
	<-osSignal
}
