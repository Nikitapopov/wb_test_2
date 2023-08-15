package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
	Применимость:
Избавление от перегруженного конструктора.
Потребность в создании разных типов объектов с похожими этапами построения.
Инкапсулирование логики построения сложных объектов.

	Плюсы:
- Инкапсуляция сложной логики построения от бизнес логики.
- Возможность построения сложных объектов последовательно.
- Использование единой точки входа для создания различных объектов.

	Минусы:
- Увеличивается количество классов.
- Клиент строителя должен знать логику построения строителей при отсутствии нужного метода директора.

	Примеры:
- Генерация отчета в разных форматах.
- Логистика, например, рассчитать маршруты для разных видов транспорта.
- Создание пользователей с разными правами.
*/

// Енам типа файлов
type fileType int16

const (
	docType fileType = iota
	xlsxType
	csvType
)

// Структура генератора отчетов
type reportGenerationDirector struct {
	builder iExportBuilder
}

func (e *reportGenerationDirector) setBuilder(builder iExportBuilder) {
	e.builder = builder
}

func (e *reportGenerationDirector) generateReport() error {
	err := e.builder.createFile()
	if err != nil {
		return err
	}

	err = e.builder.setHeaders()
	if err != nil {
		return err
	}

	err = e.builder.setBody()
	if err != nil {
		return err
	}

	err = e.builder.save()
	if err != nil {
		return err
	}

	return nil
}

func (e *reportGenerationDirector) getReportName() string {
	return e.builder.getFileName()
}

type iExportBuilder interface {
	createFile() error
	setHeaders() error
	setBody() error
	save() error
	getFileName() string
}

func getBuilder(bType fileType) iExportBuilder {
	switch bType {
	case docType:
		return newDocBuilder()
	case xlsxType:
		return newXlsxBuilder()
	case csvType:
		return newCsvBuilder()
	default:
		return nil
	}
}

type docBuilder struct {
	filename string
}

func newDocBuilder() iExportBuilder {
	return &docBuilder{}
}

func (b *docBuilder) createFile() error {
	b.filename = "7cfb429e-0220-42a9-919c-54d6886d5128.doc"
	fmt.Println("Created doc file with id = ", b.filename)
	return nil
}

func (b *docBuilder) setHeaders() error {
	fmt.Println("Headers set to doc file")
	return nil
}
func (b *docBuilder) setBody() error {
	fmt.Println("Body set to doc file")
	return nil
}

func (b *docBuilder) save() error {
	fmt.Println("Doc file saved")
	return nil
}

func (b *docBuilder) getFileName() string {
	return b.filename
}

type xlsxBuilder struct {
	filename string
}

func newXlsxBuilder() iExportBuilder {
	return &xlsxBuilder{}
}

func (b *xlsxBuilder) createFile() error {
	b.filename = "8cfb429e-0220-42a9-919c-54d6886d5128.xlsx"
	fmt.Println("Created xlsx file with id = ", b.filename)
	return nil
}

func (b *xlsxBuilder) setHeaders() error {
	fmt.Println("Headers set to xlsx file")
	return nil
}
func (b *xlsxBuilder) setBody() error {
	fmt.Println("Body set to xlsx file")
	return nil
}

func (b *xlsxBuilder) save() error {
	fmt.Println("Xlsx file saved")
	return nil
}

func (b *xlsxBuilder) getFileName() string {
	return b.filename
}

type csvBuilder struct {
	filename string
}

func newCsvBuilder() iExportBuilder {
	return &csvBuilder{}
}

func (b *csvBuilder) createFile() error {
	b.filename = "9cfb429e-0220-42a9-919c-54d6886d5128.csv"
	fmt.Println("Created csv file with id = ", b.filename)
	return nil
}

func (b *csvBuilder) setHeaders() error {
	fmt.Println("Headers set to csv file")
	return nil
}

func (b *csvBuilder) setBody() error {
	fmt.Println("Body set to csv file")
	return nil
}

func (b *csvBuilder) save() error {
	fmt.Println("Csv file saved")
	return nil
}

func (b *csvBuilder) getFileName() string {
	return b.filename
}

func main() {
	docBuilder := getBuilder(docType)

	director := reportGenerationDirector{}
	director.setBuilder(docBuilder)
	err := director.generateReport()
	if err != nil {
		// log err
		return
	}
	_ = director.getReportName()
}
