package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
	Применимость:
В случае необходимости проведения какой либо логики над элементами сложной структуры объектов
без добавления новой функциональности в сами исходные классы.

	Плюсы:
- Удобное добавление новых операций над группами объектов.
- Объединение родственных операций в одном классе.
- Возможность накапливания состояния при обходе структуры объектов.

	Минусы:
- Ухудшение поддержки кода, в случае изменения структуры объектов.
- Потенциальные нарушения инкапсуляции классов.

	Примеры:
- Суммирование значений узлов дерева.
- Экспорт данных из сложной структуры.
- Отрисовка тени в компьютерной графике (затемнение пикселей в определенной области).
*/

// Интерфейс посетителя
type visitor interface {
	visitForFrontender(*frontender)
	visitForBackender(*backender)
	visitForDevopser(*devopser)
}

// Интерфейс разработчика
type iDeveloper interface {
	accept(visitor)
}

// Структура разработчика
type developer struct {
	is_fired bool
}

// Структура разработчика фронтендера
type frontender struct {
	developer
	fixedButtonsNumber int
	fixedFormsNumber   int
}

// Метод фронтендера для посетителя
func (f *frontender) accept(visitor visitor) {
	visitor.visitForFrontender(f)
}

// Структура разработчика бэкендера
type backender struct {
	developer
	developedApiMethodsNumber int
}

// Метод бэкендера для посетителя
func (b *backender) accept(visitor visitor) {
	visitor.visitForBackender(b)
}

// Структура разработчика девопсера
type devopser struct {
	developer
	startedMicroservicesInstancesNumber int
}

// Метод девопсера для посетителя
func (d *devopser) accept(visitor visitor) {
	visitor.visitForDevopser(d)
}

// Посетитель
type manager struct{}

// Метод посетителя для фронтендера
func (m *manager) visitForFrontender(f *frontender) {
	if f.fixedButtonsNumber+(f.fixedFormsNumber*3) < 50 {
		f.is_fired = true
	}
}

// Метод посетителя для бэкендера
func (m *manager) visitForBackender(b *backender) {
	if b.developedApiMethodsNumber > 20 {
		b.is_fired = true
	}
}

// Метод посетителя для девопсера
func (m *manager) visitForDevopser(d *devopser) {
	if d.startedMicroservicesInstancesNumber > 1 {
		d.is_fired = true
	}
}

func main() {
	// Слайс разработчиков и его заполнение
	team := []iDeveloper{}
	team = append(team, &frontender{fixedButtonsNumber: 100, fixedFormsNumber: 10})
	team = append(team, &frontender{fixedButtonsNumber: 5, fixedFormsNumber: 1})
	team = append(team, &backender{developedApiMethodsNumber: 20})
	team = append(team, &backender{developedApiMethodsNumber: 15})
	team = append(team, &backender{developedApiMethodsNumber: 10})
	team = append(team, &devopser{startedMicroservicesInstancesNumber: 3})
	team = append(team, &devopser{startedMicroservicesInstancesNumber: 10})

	// Объявление посетителя
	teamlead := manager{}

	// Итерация по разработчикам, оценивание эффективности и увольнение малоэффективных.
	for _, member := range team {
		member.accept(&teamlead)
	}
}
