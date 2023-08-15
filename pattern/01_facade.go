package pattern

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
	Применимость:
В случаях, когда нужно скрыть большую сложную логику и предоставить только упрощенный интерфейс взаимодействия. Таким образом
можно инкапсулировать группу классов, разделять разные слои.

	Плюсы:
- Инкапсуляция сложной логики за простым интерфейсом.

	Минусы:
- Риск большой связности с другими классами.
- Взятие на себя большой ответственности.

	Примеры:
- Регистрация заявки на подключение услуги
- Авторизация
- Обработка данных, например, загрузка изображения в соцсеть
- Отправка оповещений
*/

import (
	"fmt"
	"time"
)

// Загрузчик изображений, реализованный паттерном "Фасад"

// Интерфейс загрузчика
type iImageDowloaderFacade interface {
	download(userId string, img []byte) error
}

// Структура загрузчика
type imageDowloaderFacade struct {
	dbClient    dbClient
	cacheClient cacheClient
	improver    improver
	notifier    notifier
	checker     checker
	logger      logger
}

// Конструктор загрузчика
func NewImageDowloaderFacade(
	dbClient dbClient,
	cacheClient cacheClient,
	improver improver,
	notifier notifier,
	checker checker,
	logger logger,
) iImageDowloaderFacade {
	return &imageDowloaderFacade{
		dbClient:    dbClient,
		cacheClient: cacheClient,
		improver:    improver,
		notifier:    notifier,
		checker:     checker,
		logger:      logger,
	}
}

/*
Метод загрузчика для загрузки изображения data пользователя userId.
Сопровождается загрузкой в бд, кэш, обработкой изображения, проверкой на соблюдение правил размещения,
рассылкой уведомлений подписчиков.
*/
func (f *imageDowloaderFacade) download(userId string, data []byte) error {
	// Дата и время загрузки
	downloadTimestamp := time.Now()

	// Добавление в бд
	err := f.dbClient.add(userId, downloadTimestamp, data)
	if err != nil {
		logRecord := fmt.Sprintf("db adding: %v", err)
		f.logger.write(logRecord)
	}

	// Проверка на соблюдение правил размещения
	err = f.checker.check(userId, data)
	if err != nil {
		logRecord := fmt.Sprintf("checking: %v", err)
		f.logger.write(logRecord)
	}

	// Обработка изображения
	improvedData, err := f.improver.handle(data)
	if err != nil {
		logRecord := fmt.Sprintf("improver: %v", err)
		f.logger.write(logRecord)
	}

	// Добавление в кэш
	imgId, err := f.cacheClient.add(userId, improvedData)
	if err != nil {
		logRecord := fmt.Sprintf("cache adding: %v", err)
		f.logger.write(logRecord)
	}

	// Уведомление подписчиков
	err = f.notifier.sendNotificationsToSubscribers(userId, imgId)
	if err != nil {
		logRecord := fmt.Sprintf("sending notifications: %v", err)
		f.logger.write(logRecord)
	}

	f.logger.write("downloading image")

	return nil
}

// Сервис бд
type dbClient struct{}

func (c *dbClient) add(userId string, downloadTimestamp time.Time, data []byte) error {
	fmt.Printf("Adding image by user id = %s to db\n", userId)
	return nil
}

// Сервис кэша
type cacheClient struct{}

func (c *cacheClient) add(userId string, data []byte) (string, error) {
	fmt.Printf("Adding image by user id = %s to cache\n", userId)
	return "", nil
}

// Сервис улучшения изображений
type improver struct{}

func (i *improver) handle(data []byte) ([]byte, error) {
	fmt.Println("Handling image in improver service")
	return []byte{}, nil
}

// Сервис проверки контента
type checker struct{}

func (c *checker) check(userId string, data []byte) error {
	fmt.Println("Checking image in checker service")
	return nil
}

// Сервис уведомлений
type notifier struct{}

func (n *notifier) sendNotificationsToSubscribers(userId string, imgId string) error {
	fmt.Println("Sending notifications about image in notifier service")
	return nil
}

// Сервис логирования
type logger struct{}

func (l *logger) write(record string) error {
	fmt.Printf("Writing log record: %s\n", record)
	return nil
}

func main() {
	// Инициализирование сервисов
	dbClient := dbClient{}
	cacheClient := cacheClient{}
	improver := improver{}
	notifier := notifier{}
	checker := checker{}
	logger := logger{}

	// Данные изображения
	userId := "550e8400-e29b-41d4-a716-446655440000"
	imgData := []byte("asdfasdfasdfasdfasdf")

	// Использование фасада загрузчика изображений для загрузки изображения
	downloader := NewImageDowloaderFacade(dbClient, cacheClient, improver, notifier, checker, logger)
	downloader.download(userId, imgData)
}
