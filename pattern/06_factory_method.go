package pattern

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
	Применимость:
В случаях, когда классу заранее неизвестно, объекты каких типо ему нужно создавать. Обязанности создания делегируются подклассам.

	Плюсы:
- Уменьшает связность между классом и конкретными классами продуктов.
- Помещает весь код, отвечающий за производство продуктов, в одно место.
- Упрощает добавление новых продуктов.

	Минусы:
- Увеличение числа классов, реализующих отдельные продукты.

	Примеры:
- Создание разных типов пользователей.
- Организация логистика разным типом транспорта.
- Генерация отчетов с разными расширениями.
*/

import (
	"fmt"
	"time"
)

// енам пользователя
type userType int16

const (
	staffUserType userType = iota
	visitorUserType
	clientUserType
)

// Интерфейс продукта
type iUser interface {
	// Метод установки персональных данных
	setPersonalInfo(any) error
	// Метод установки доступов accesses к локации с id locationId.
	giveAccessToLocation(locationId string, accesses []string) error
	// Метод для получения времени проведенного на локации с id locationId в промежутке от from до to.
	getTimeSpendingOnLocation(locationId string, from time.Time, to time.Time) ([]byte, error)
}

// Структура продукта
type user struct {
	name       string
	department string
}

// Структура конкретного продукта сотрудник
type staff struct {
	user
	position string
}

// Структура персональных данных сотрудника
type staffPersonalInfo struct {
	position string
}

func (s *staff) setPersonalInfo(info any) error {
	fmt.Println("Set personal staff info")
	return nil
}

func (s *staff) giveAccessToLocation(locationId string, accesses []string) error {
	fmt.Println("Giving an access to location to staff")
	return nil
}

func (s *staff) getTimeSpendingOnLocation(locationId string, from time.Time, to time.Time) ([]byte, error) {
	fmt.Println("Exporting time spending on location by staff")
	return []byte{}, nil
}

// Структура конкретного продукта посетитель
type visitor struct {
	user
	accompanyingPerson string
}

// Структура персональных данных посетителя
type visitorPersonalInfo struct {
	accompanyingPersonId string
}

func (v *visitor) setPersonalInfo(info any) error {
	fmt.Println("Set personal visitor info")
	return nil
}

func (v *visitor) giveAccessToLocation(locationId string, accesses []string) error {
	fmt.Println("Giving an access to location to visitor")
	return nil
}

func (v *visitor) getTimeSpendingOnLocation(locationId string, from time.Time, to time.Time) ([]byte, error) {
	fmt.Println("Exporting time spending on location by visitor")
	return []byte{}, nil
}

// Структура конкретного продукта клиент
type client struct {
	user
}

func (c *client) setPersonalInfo(info any) error {
	fmt.Println("Set personal client info")
	return nil
}

func (c *client) giveAccessToLocation(locationId string, accesses []string) error {
	fmt.Println("Giving an access to location to client")
	return nil
}

func (c *client) getTimeSpendingOnLocation(locationId string, from time.Time, to time.Time) ([]byte, error) {
	fmt.Println("Exporting time spending on location by client")
	return []byte{}, nil
}

// Фабричный метод
func getUser(userType userType) iUser {
	switch userType {
	case staffUserType:
		return &staff{}
	case visitorUserType:
		return &visitor{}
	case clientUserType:
		return &client{}
	default:
		return nil
	}
}

func main() {
	// Переменные для фунции getTimeSpendingOnLocation всех трех продуктов
	from, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	to, _ := time.Parse(time.RFC3339, "2007-01-02T15:04:05Z")

	// Создание продуктов из фабричного метода и запуск его методов
	staffUser := getUser(staffUserType)
	staffUser.setPersonalInfo(staffPersonalInfo{position: "programmer"})
	staffUser.giveAccessToLocation("1234", []string{"1234", "5678"})
	staffUser.getTimeSpendingOnLocation("1234", from, to)

	visitorUser := getUser(visitorUserType)
	visitorUser.setPersonalInfo(visitorPersonalInfo{accompanyingPersonId: "1234"})
	visitorUser.giveAccessToLocation("1234", []string{"1234", "5678"})
	visitorUser.getTimeSpendingOnLocation("1234", from, to)

	clientUser := getUser(clientUserType)
	clientUser.setPersonalInfo(visitorPersonalInfo{accompanyingPersonId: "1234"})
	clientUser.giveAccessToLocation("1234", []string{"1234", "5678"})
	clientUser.getTimeSpendingOnLocation("1234", from, to)
}
