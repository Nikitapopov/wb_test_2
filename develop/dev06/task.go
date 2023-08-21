package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// Чтение аргументов командной строки
	args := cmdArgs{}
	err := args.parseArgs()
	if err != nil {
		errMsg := fmt.Sprintf("error in flags: %v", err)
		os.Stderr.WriteString(errMsg)
		os.Exit(1)
	}

	// Получение исходных строк из stdin
	lines := scanLines()

	// Получение запрошеных колонок
	cutlines := cut(lines, args)

	// Вывод запрошенных колонок
	printLines(cutlines)
}

// Структура флагов вырезания строк
type cmdArgs struct {
	// Выбрать поля (колонки)
	fields []int
	// Использовать другой разделитель
	delimiter string
	// Только строки с разделителем
	separated bool
}

// Метод для инициализирования аргументов
func (a *cmdArgs) parseArgs() error {
	// Привязка аргументов к полям структуры
	var fieldsAsStr string
	flag.StringVar(&fieldsAsStr, "f", "1", "Выбрать поля (колонки)")
	flag.StringVar(&a.delimiter, "d", "\t", "Использовать другой разделитель")
	flag.BoolVar(&a.separated, "s", false, "Только строки с разделителем")

	flag.Parse()

	// Парсинг переменной поля
	var err error
	a.fields, err = parseFields(fieldsAsStr)
	if err != nil {
		return err
	}

	return nil
}

// Функция для парсинга переменной поля
func parseFields(field string) ([]int, error) {
	fields := []int{}

	// Если поле содержит запятую
	if strings.Contains(field, ",") {
		// Разделение строки по запятым и приведение значений к числам
		splitFields := strings.Split(field, ",")
		for _, field := range splitFields {
			num, err := strconv.Atoi(field)
			if err != nil {
				return []int{}, errors.New("incorrect f flag")
			}

			fields = append(fields, num)
		}

		// Получение уникальных значений
		uniqueFields := getUnique(fields)

		// Сортировка значений
		sort.Ints(uniqueFields)

		return uniqueFields, nil
	}

	// Если поле содержит тире
	if strings.Contains(field, "-") {
		// Разделение строки по тире
		splitFields := strings.Split(field, "-")

		// Валидация значений
		if len(splitFields) != 2 || splitFields[0] > splitFields[1] {
			return []int{}, errors.New("incorrect f flag")
		}

		// Приведение значений к числам
		lowColNum, err := strconv.Atoi(splitFields[0])
		if err != nil {
			return []int{}, errors.New("incorrect f flag")
		}

		highColNum, err := strconv.Atoi(splitFields[1])
		if err != nil {
			return []int{}, errors.New("incorrect f flag")
		}

		// Получение диапазона значений
		for i := lowColNum; i <= highColNum; i++ {
			fields = append(fields, i)
		}

		return fields, nil
	}

	// Если поле это число, то приводим строку к числу
	num, err := strconv.Atoi(field)
	if err != nil {
		return []int{}, errors.New("incorrect f flag")
	}

	return []int{num}, nil
}

// Функция для получения уникальных значений в исходном порядке слайса values
func getUnique(values []int) []int {
	uniqueValues := make([]int, 0, len(values)/2)
	set := make(map[int]struct{}, len(values))
	for _, value := range values {
		_, ok := set[value]
		if !ok {
			set[value] = struct{}{}
			uniqueValues = append(uniqueValues, value)
		}
	}

	return uniqueValues
}

// Метод для сканирования построчного ввода
func scanLines() []string {
	lines := []string{}

	fmt.Println(`Введите строки. Для завершения ввода введите "\q"`)
	sc := bufio.NewScanner(os.Stdin)

	// Сканирование введенных строк до тех пор, пока не введена строка "\q"
	for sc.Scan() {
		txt := sc.Text()
		if txt == "\\q" {
			break
		}
		lines = append(lines, txt)
	}

	return lines
}

// Функция для вырезки пределенных колонок и строк lines
func cut(lines []string, args cmdArgs) []string {
	result := []string{}

	// Если планируется получать ровно все строки, то выполняется выделение памяти
	if !args.separated {
		result = make([]string, 0, len(lines))
	}

	// Итерация по строкам
	for _, line := range lines {
		fields := []string{}

		// Разделение строки по разделителю
		splitLine := strings.Split(line, args.delimiter)

		// Если строка не содержит разделителя, а это требуется, то переход на следующую итерацию
		if args.separated && len(splitLine) < 2 {
			continue
		}

		// Если искомые колонки имеются в строке, то добавляем, пока не закончатся колонки в полях или в строке
		for i := 0; i < len(args.fields) && args.fields[i] <= len(splitLine); i++ {
			fields = append(fields, splitLine[args.fields[i]-1])
		}

		// добавление искомых колонок в результат
		result = append(result, strings.Join(fields, args.delimiter))
	}

	return result
}

// Функция для печати строк
func printLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}
