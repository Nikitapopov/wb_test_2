package pattern

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/*
	Применимость:
Применим при наличие объекта, поведение которого меняется в зависимости от внутреннего состояния, причем этих состояний много и они часто меняются.
Когда в программе используется машина состояний и нужно исключить дублирование для похожих состояний.

	Плюсы:
- Вынос логики каждого отдельного состояния в отдельный класс.
- Избавление от множества больших условных операторов машины состояния.
- Упрощение кода контекста.

	Минусы:
- При малом количестве состояний или их редком изменении может неоправданно усложнить код.

	Примеры:
- Логика работы вендингового автомата.
- Описание сценария пользования приложением.
- Операции с заказом на разных этапах доставки.
*/

import (
	"fmt"
	"math/rand"
)

// Тип енама состояния автомата
type stateType int16

// Константы енама состояния автомата
const (
	soldOutStateType stateType = iota
	noQuarterStateType
	hasQuarterStateType
	soldStateType
	winnerStateType
)

// Интерфейс состояния
type state interface {
	insertQuarter()
	ejectQuarter()
	turnCrank()
	dispense()
	refill()
}

// Интерфейс автомата
type iGumballMachine interface {
	InsertQuarter()
	EjectQuarter()
	TurnCrank()
	Refill(int)
	setState(stateType)
	releaseGum()
	getCount() int
}

// Структура автомата
type gumballMachine struct {
	soldOutState    state
	noQuarterState  state
	hasQuarterState state
	soldState       state
	winnerState     state
	state           state
	count           int
}

// Конструктор автомата
func NewGumnallMachine(count int) iGumballMachine {
	// Инициализация автомата
	machine := &gumballMachine{}

	// Инициализация состояний
	soldOutState := &soldOutState{
		machine: machine,
	}

	noQuarterState := &noQuarterState{
		machine: machine,
	}

	hasQuarterState := &hasQuarterState{
		machine: machine,
	}

	soldState := &soldState{
		machine: machine,
	}

	winnerState := &winnerState{
		machine: machine,
	}

	// Заполнение полей автомата состояниями
	machine.soldOutState = soldOutState
	machine.noQuarterState = noQuarterState
	machine.hasQuarterState = hasQuarterState
	machine.soldState = soldState
	machine.winnerState = winnerState

	// В зависимости от параметра count заполняется поле count и текущий статус
	if count > 0 {
		machine.count = count
		machine.state = noQuarterState
	} else {
		machine.count = 0
		machine.state = soldOutState
	}

	return machine
}

// Метод автомата для вставки монеты
func (m *gumballMachine) InsertQuarter() {
	m.state.insertQuarter()
}

// Метод автомата для возврата монеты
func (m *gumballMachine) EjectQuarter() {
	m.state.ejectQuarter()
}

// Метод автомата для проворота рукоятки
func (m *gumballMachine) TurnCrank() {
	m.state.turnCrank()
	m.state.dispense()
}

// Метод автомата для пополнения автомата жевачками
func (m *gumballMachine) Refill(count int) {
	m.count += count
	m.state.refill()
}

// Метод автомата для смены состояния
func (m *gumballMachine) setState(stateType stateType) {
	switch stateType {
	case soldOutStateType:
		m.state = m.soldOutState
	case noQuarterStateType:
		m.state = m.noQuarterState
	case hasQuarterStateType:
		m.state = m.hasQuarterState
	case soldStateType:
		m.state = m.soldState
	case winnerStateType:
		m.state = m.winnerState
	}
}

// Метод автомата декрементирования счетчика жевачек (выдача жевачки автоматом)
func (m *gumballMachine) releaseGum() {
	if m.count > 0 {
		fmt.Println("A gumball comes rolling out the slot...")
		m.count--
	}
}

// Метод автомата для получения количества оставшихся жевачек
func (m *gumballMachine) getCount() int {
	return m.count
}

// Состояние автомата "все продано"
type soldOutState struct {
	machine *gumballMachine
}

func (s *soldOutState) insertQuarter() {
	fmt.Println("You can't insert a quarter, the machine is sold out")
}

func (s *soldOutState) ejectQuarter() {
	fmt.Println("You can't eject, you haven't inserted a quarter yet")
}

func (s *soldOutState) turnCrank() {
	fmt.Println("You turned, but there are no gumballs")
}

func (s *soldOutState) dispense() {
	fmt.Println("No gumball dispensed")
}

func (s *soldOutState) refill() {
	s.machine.setState(noQuarterStateType)
}

// Состояние автомата: "монета не вставлена"
type noQuarterState struct {
	machine *gumballMachine
}

func (s *noQuarterState) insertQuarter() {
	fmt.Println("You inserted a quarter")
	s.machine.setState(hasQuarterStateType)
}

func (s *noQuarterState) ejectQuarter() {
	fmt.Println("You haven't inserted a quarter")
}

func (s *noQuarterState) turnCrank() {
	fmt.Println("You turned but there's no quarter")
}

func (s *noQuarterState) dispense() {
	fmt.Println("You need to pay first")
}

func (s *noQuarterState) refill() {
}

// Состояние автомата "монета вставлена"
type hasQuarterState struct {
	machine *gumballMachine
}

func (s *hasQuarterState) insertQuarter() {
	fmt.Println("You can't insert another quarter")
}

func (s *hasQuarterState) ejectQuarter() {
	fmt.Println("Quarter returned")
	s.machine.setState(noQuarterStateType)
}

func (s *hasQuarterState) turnCrank() {
	fmt.Println("You turned...")
	// Если количество жечавек в автомате больше 1, то с 10% вероятностью автомат переходит в состояние "выигрыш"
	if s.machine.getCount() > 1 && rand.Intn(10) > 8 {
		s.machine.setState(winnerStateType)
	} else { // Иначе переход в состояние "жевачка продана"
		s.machine.setState(soldStateType)
	}
}

func (s *hasQuarterState) dispense() {
	fmt.Println("You need to turn the crank")
}

func (s *hasQuarterState) refill() {}

// Состояние автомата "жевачка продана"
type soldState struct {
	machine *gumballMachine
}

func (s *soldState) insertQuarter() {
	fmt.Println("Please wait, we're already giving you a gumball")
}

func (s *soldState) ejectQuarter() {
	fmt.Println("Sorry, you already turned the crank")
}

func (s *soldState) turnCrank() {
	fmt.Println("Turning twice doesn't get you another gumball!")
}

func (s *soldState) dispense() {
	s.machine.releaseGum()
	if s.machine.getCount() == 0 {
		fmt.Println("Oops, out of gumballs!")
		s.machine.setState(soldOutStateType)
	} else {
		s.machine.setState(noQuarterStateType)
	}
}

func (s *soldState) refill() {}

// Состояние автомата "выигрыш"
type winnerState struct {
	machine *gumballMachine
}

func (s *winnerState) insertQuarter() {
	fmt.Println("Please wait, we're already giving you a gumball")
}

func (s *winnerState) ejectQuarter() {
	fmt.Println("Sorry, you already turned the crank")
}

func (s *winnerState) turnCrank() {
	fmt.Println("Turning twice doesn't get you another gumball!")
}

func (s *winnerState) dispense() {
	// Выдача жевачки
	s.machine.releaseGum()
	// Если количество жевачек в автомате равно нулю, то выводим сообщение об ошибке и переходим в состояние "Все продано"
	if s.machine.getCount() == 0 {
		fmt.Println("Oops, out of gumballs!")
		s.machine.setState(soldOutStateType)
	} else { // Иначе выдача еще одной жевачки
		s.machine.releaseGum()
		fmt.Println("Winner!")
		// Если количество жевачек в автомате равно нулю, то выводим сообщение об ошибке и переходим в состояние "Все продано"
		if s.machine.getCount() == 0 {
			fmt.Println("Oops, out of gumballs!")
			s.machine.setState(soldOutStateType)
		} else { // Иначе переход в сотояние "монета не вставлена"
			s.machine.setState(noQuarterStateType)
		}
	}
}

func (s *winnerState) refill() {}

func main() {
	// Изначальное количество жевачек
	count := 2

	// Создание автомата
	gumballMachine := NewGumnallMachine(count)

	// Операции вставки монеты и поворота рычага для получения жевачки
	gumballMachine.InsertQuarter()
	gumballMachine.TurnCrank()
}
