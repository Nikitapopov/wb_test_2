package pattern

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
	Применимость:
В случае необходимости параметризировать объекты выполняемыми действиями. Т.е. отделить вызывающую и выполняющую логику.
Также можно применять паттерн "Команда", когда необходима возможность накапливать в очереди операции.
Необходимость возможности отмены операций.

	Плюсы:
- Уменьшение связности между объектами, которые вызывают и выполняют операции.
- Облегачает использование композитных операций.
- Облегчение реализации операции отмены действия.
- Облегчение реализации очередности и накапливания операций.

	Минусы:
- Увеличение количества дополнительных классов

	Примеры:
- Реализация графического редактора с большим количеством операций, которые выполняются разными способами.
- Организация очередности загрузки видео на хостинг.
- Удаленное управление устройствами/сервисами с нескольких передатчиков.
*/

import (
	"fmt"
)

// Интерфейс команды
type command interface {
	execute()
}

// Структура умного дома - отправитель команд
type cleverHouse struct{}

// Конкретная команда заказа еды на количество денег money
type orderFoodCommand struct {
	money int
}

func (c *orderFoodCommand) execute() {
	fmt.Printf("Ordering food for %d rubles\n", c.money)
}

// Конкретная команда вызова горничной
type callMaidCommand struct{}

func (c *callMaidCommand) execute() {
	fmt.Println("Calling out of maid")
}

// Конкретная команда выброса мусора
type throwTrashCommand struct{}

func (c *throwTrashCommand) execute() {
	fmt.Println("Trash throwing")
}

// Метод умного дома для вызова команды заказа еды
func (ch *cleverHouse) orderFood(money int) command {
	return &orderFoodCommand{money: money}
}

// Метод умного дома для вызова команды вызова горничной
func (ch *cleverHouse) callMaid() command {
	return &callMaidCommand{}
}

// Метод умного дома для вызова команды выбрасывания мусора
func (ch *cleverHouse) throwTrash() command {
	return &throwTrashCommand{}
}

// Структура управляющего дома - получатель команд
type houseManager struct {
	commands []command
}

func (m *houseManager) executeCommands() {
	for _, c := range m.commands {
		c.execute()
	}
}

func main() {
	// Объявление получателя
	cleverhouse := cleverHouse{}

	// Получение команд
	tasks := []command{
		cleverhouse.throwTrash(),
		cleverhouse.orderFood(5000),
		cleverhouse.callMaid(),
		cleverhouse.throwTrash(),
	}

	// Загрузка команд в получателя
	houseManager := houseManager{
		commands: tasks,
	}

	// Выполнение команд получателем
	houseManager.executeCommands()
}
