package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
	Применимость:
В случаях, когда нужно менять логику алгоритма внутри одного объекта. Или когда имеется ножество похожих алгоритмов, чье поведение различается незначительно.

	Плюсы:
- Смена алгоритма на лету.
- Инкапсуляция деталей каждого алгоритма в отдельном классе.
- Замена наследования на делегирование.

	Минусы:
- Увеличение количества классов.
- Клиент должен обладать дополнительной логикой для определения какую стратегию выбирать.

	Примеры:
- Выбор построения маршрута в зависимости от выбранного транспорта.
- Выбор алгоритма сортировки.
- Получение инструкции готовки различных напитков.
*/

// Интерфейс стратегии сортировки
type sortingStrategy interface {
	// Метод для сортировки слайса arr на месте
	sort(arr []int)
}

// Структура стратегии сортировки пузырьком
type bubbleSortingStrategy struct{}

func (s *bubbleSortingStrategy) sort(arr []int) {
	fmt.Println("Sorting by bubble method")
}

// Структура стратегии сортировки объединением
type mergeSortingStrategy struct{}

func (s *mergeSortingStrategy) sort(arr []int) {
	fmt.Println("Sorting by merging method")
}

// Структура стратегии быстрой сортировки
type quickSortingStrategy struct{}

func (s *quickSortingStrategy) sort(arr []int) {
	fmt.Println("Sorting by quick method")
}

// Структура обработчика данных
type dataHandler struct {
	data         []int
	sortStrategy sortingStrategy
}

// Конструктор обработчика данных
func newDataHandler(data []int, strategy sortingStrategy) dataHandler {
	return dataHandler{
		data:         data,
		sortStrategy: strategy,
	}
}

// Метод для установки стратегии сортировки
func (h *dataHandler) setSortStrategy(strategy sortingStrategy) {
	h.sortStrategy = strategy
}

// Метод для обработки данных
func (h *dataHandler) handle() {
	h.sortStrategy.sort(h.data)
}

func main() {
	// Исходные данные
	data := []int{5, 2, 4, 1, 3}

	// Создание переменной стратегии сортировки пузырьком
	bubbleSortStrategy := &bubbleSortingStrategy{}

	// Создание и инициализация обработчика данных
	handler := newDataHandler(data, bubbleSortStrategy)

	// Обработка данных
	handler.handle()
}
