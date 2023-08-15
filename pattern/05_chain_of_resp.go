package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
	Применимость:
Применим когда, нужно обрабатывать разнообразные запросы, но неизвестно в каком порядке и какие именно.
Когда нужно, чтобы обработчики выполнялись один за одним в строгом порядке или только какой то один из них.
Также применим в случаях, когда обработчики должны задаваться динамически.

	Плюсы:
- Выделяет каждый случай обработки в отдельный класс.
- Уменьшает связность между клиентом и обработчиками.

	Минусы:
- Увеличение количества классов.

	Примеры:
- Авторизация, валидация и другие подобные проверки.
- Обработка событий в html дереве.
- Перевод статуса пользователя в нужное состояние
*/

// Структура запроса
type request struct {
	method      string
	url         string
	accessToken string
	params      string
}

// Обработчик запроса
type requestCheck interface {
	execute(request)
	setNext(requestCheck)
}

// Конкретный обработчик проверки на существование эндпоинта
type endpointExistenceCheck struct {
	next requestCheck
}

func (c *endpointExistenceCheck) execute(request request) {
	fmt.Printf("Cheching endpoint existence: %s %s\n", request.method, request.url)
}

func (c *endpointExistenceCheck) setNext(next requestCheck) {
	c.next = next
}

// Конкретный обработчик проверки авторизации
type authenticationCheck struct {
	next requestCheck
}

func (c *authenticationCheck) execute(request request) {
	fmt.Printf("Cheching authentication token: %s\n", request.accessToken)
}

func (c *authenticationCheck) setNext(next requestCheck) {
	c.next = next
}

// Конкретный обработчик проверки валидации параметров
type paramsValidationCheck struct {
	next requestCheck
}

func (c *paramsValidationCheck) execute(request request) {
	fmt.Printf("Params validation: %s\n", request.params)
}

func (c *paramsValidationCheck) setNext(next requestCheck) {
	c.next = next
}

func main() {
	// Пример запроса
	request := request{
		method:      "POST",
		url:         "comment",
		accessToken: "123123123",
		params:      "text:It's cool",
	}

	// Объявление первого обработчика цепочки
	endpointExistenceCheck := &endpointExistenceCheck{}

	// Объявление и установка второго обработчика цепочки
	authenticationCheck := &authenticationCheck{}
	endpointExistenceCheck.setNext(authenticationCheck)

	// Объявление и установка третьего обработчика цепочки
	paramsValidationCheck := &paramsValidationCheck{}
	authenticationCheck.setNext(paramsValidationCheck)

	// Запуск цепочки обработчиков
	endpointExistenceCheck.execute(request)
}
