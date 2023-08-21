package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const (
	// Файл для вывода отсортированных строк
	outputFilename = "./data/output.txt"
)

func main() {
	// Чтение аргументов командной строки
	args := cmdArgs{}
	err := args.parseArgs()
	if err != nil {
		errMsg := fmt.Sprintf("error in flags: %v", err)
		os.Stderr.WriteString(errMsg)
		os.Exit(1)
	}

	// Получение исходных строк из файла
	lines, err := readFile(args.filename)
	if err != nil {
		errMsg := fmt.Sprintf("reading file: %v", err)
		os.Stderr.WriteString(errMsg)
		os.Exit(1)
	}

	// Фильтрация строк
	filteredLines, err := filter(lines, args)
	if err != nil {
		errMsg := fmt.Sprintf("filter: %v", err)
		os.Stderr.WriteString(errMsg)
		os.Exit(1)
	}

	// Объединение строк в одну
	result := strings.Join(filteredLines, "\n")

	// Сохранение строк в файл
	err = saveResult(result, outputFilename)
	if err != nil {
		msg := fmt.Sprintf("Saving result: %v", err)
		os.Stderr.WriteString(msg)
		os.Exit(1)
	}
}

// Структура флагов фильтрации строк файла
type cmdArgs struct {
	// Печатать +N строк после совпадения
	after int
	// Печатать +N строк до совпадения
	before int
	// Печатать ±N строк вокруг совпадения
	context int
	// Количество строк
	count bool
	// Игнорировать регистр
	ignoreCase bool
	// Вместо совпадения, исключать
	invert bool
	// Точное совпадение со строкой, не паттерн
	fixed bool
	// Печатать номер строки
	lineNum bool
	// Образец для поиска
	pattern string
	// Название файла для фильтрации
	filename string
}

// Метод для инициализирования аргументов
func (a *cmdArgs) parseArgs() error {
	// Привязка аргументов к полям структуры

	// Для флагов A и B проставлено значение по умолчанию -1, чтобы не путать случаи,
	// когда указано значение 0 и когда значение не указано.
	flag.IntVar(&a.after, "A", -1, "Печатать +N строк после совпадения")
	flag.IntVar(&a.before, "B", -1, "Печатать +N строк до совпадения")
	flag.IntVar(&a.context, "C", 0, "Печатать ±N строк вокруг совпадения")
	flag.BoolVar(&a.count, "c", false, "Количество строк")
	flag.BoolVar(&a.ignoreCase, "i", false, "Игнорировать регистр")
	flag.BoolVar(&a.invert, "v", false, "Вместо совпадения, исключать")
	flag.BoolVar(&a.fixed, "F", false, "Точное совпадение со строкой, не паттерн")
	flag.BoolVar(&a.lineNum, "n", false, "Печатать номер строки")

	// Загрузка значений
	flag.Parse()

	// Если количество строк после совпадения не задано, то берем значение из количества строк вокруг
	if a.after == -1 {
		a.after = a.context
	}

	// Если количество строк до совпадения не задано, то берем значение из количества строк вокруг
	if a.before == -1 {
		a.before = a.context
	}

	// 2 значения без флагов устанавливаем в поля образец и файл
	other := flag.Args()
	if len(other) < 2 {
		return errors.New("two required parameters are missing")
	}

	a.pattern = other[0]
	a.filename = other[1]

	return nil
}

// Метод для чтения файла построчно
func readFile(fileName string) ([]string, error) {
	// Открытие файла
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := []string{}

	// Построчное чтение файла
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Функция для сохранения данных data в файл с именем filename
func saveResult(data string, fileName string) error {
	// Создание файла
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Запись в файл
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func filter(lines []string, args cmdArgs) ([]string, error) {
	// поиск индексов строк совпадений
	matchesIndexes, err := getMatches(lines, args.pattern, args.ignoreCase, args.invert, args.fixed)
	if err != nil {
		return []string{}, err
	}

	// Если требуется вернуть только количество строк
	if args.count {
		return []string{strconv.Itoa(len(matchesIndexes))}, nil
	}

	// Добавление дополнительных данных к отфильтрованным строкам
	res := formResult(lines, matchesIndexes, args.after, args.before, args.lineNum)

	return res, nil
}

// Функция для поиска строк в слайсе lines, в которых присутствует подстрока pattern
// i - Поиск регистронезависимый
// v - Получение исключений вместо поиска
// F - Точное совпадение со строкой, а не паттерн
// Возвращаемое значение - индексы строк, в которых было найдено соответствие
func getMatches(lines []string, pattern string, i, v, F bool) ([]int, error) {
	doesMatchSearchOption := getIDKWrapper(pattern, i, v, F)

	matchesLinesIndexes := []int{}
	for lineIndex, line := range lines {
		matched, err := doesMatchSearchOption(line)
		if err != nil {
			return []int{}, err
		}

		if matched {
			matchesLinesIndexes = append(matchesLinesIndexes, lineIndex)
		}
	}

	return matchesLinesIndexes, nil
}

// Обертка функции определения подходит ли строка str под фильтрацию по подстроке subStr с флагами i, v, F.
// Описание параметров-флагов представлено в функции @link{getMatches}
func getIDKWrapper(subStr string, i, v, F bool) func(string) (bool, error) {
	// Получение функции применения флага v инвертирования
	invertIfNeed := getInvertionWrapper(v)

	// возвращаемая функия
	return func(str string) (bool, error) {
		// Локальная копия подстроки поиска
		pattern := subStr

		// Если выполняется фильтрация с помощью точного совпадения со строкой
		if F {
			// Если фильтрация регистронезависимая, то приводим строку и подстроку к строчным символам
			if i {
				str = strings.ToLower(str)
				pattern = strings.ToLower(subStr)
			}

			// Поиск подстроки в строке и применение опции инвертирования
			return invertIfNeed(strings.Contains(str, pattern)), nil
		}

		// Если выполняется фильтрация с помощью совпадения по паттерну

		// Если фильтрация регистронезависимая, то добавляется флаг регистронезависимости в подстроку
		if i {
			pattern = "(?i)" + subStr
		}

		// Выполняется поиск подстроки регулярным выражением
		matched, err := regexp.MatchString(pattern, str)
		if err != nil {
			return false, err
		}

		// Применение опции инвертирования
		return invertIfNeed(matched), nil
	}
}

// Обертка функции определения подходит ли значение value под фильтрацию согласно необходимости инвертирования результата
func getInvertionWrapper(needExclude bool) func(bool) bool {
	return func(value bool) bool {
		return value != needExclude
	}
}

// Формирование результата фильтрации строк lines по индексам indexes согласно флагам A, B и n
func formResult(lines []string, indexes []int, A, B int, n bool) []string {
	// Слайс для хранения отфильтрованых и дополнительных строк.
	result := make([]string, 0, len(indexes))

	// Функция записи строки в result
	write := getWriterWrapper(&result, lines, n)

	// Итерация по офильтрованным индексам
	for i := 0; i < len(indexes); i++ {
		// Индекс одной из отфильтрованных строк
		matchedIndex := indexes[i]

		// lowAdditionalIndex и highAdditionalIndex диапазон индексов, строки которого необходимо внести в результат,
		// так как они находятся вокруг i-го совпадения
		lowAdditionalIndex := matchedIndex - B
		// Итерация по индексам в количестве B до строки совпадения
		for j := lowAdditionalIndex; j < matchedIndex; j++ {
			// Если совпадающая строка не имеет перед собой достаточное B количество строк
			if j < 0 {
				continue
			}

			// Проверка что текущий цикл не вылезет за пределы предыдущего цикла
			if i > 0 {
				prevMaxHighAdditionalIndex := indexes[i-1] + A
				if prevMaxHighAdditionalIndex >= j {
					continue
				}
			}

			// Запись j-ой доп строки
			write(j, false)
		}

		// Запись строки совпадения
		write(matchedIndex, true)

		highAdditionalIndex := matchedIndex + A + 1
		for j := matchedIndex + 1; j < highAdditionalIndex && j < len(lines); j++ {
			if i < len(indexes)-1 {
				nextIndex := indexes[i+1]
				if nextIndex <= j {
					break
				}
			}

			// Запись j-ой доп строки
			write(j, false)
		}
	}

	return result
}

// Обертка функции записи строки с индексом index из lines в качестве строки совпадения (isMain=true)
// или строки окружения совпадения (isMain=false).
// Запись выполняется в слайс result, используя правила форматирование, определенные флагом n.
func getWriterWrapper(result *[]string, lines []string, n bool) func(index int, isMain bool) {
	return func(index int, isMain bool) {
		// Добавляемая в результат строка
		record := ""

		// Если необходимо писать в результате номер строки, то записываем
		if n {
			var splitter string
			if isMain {
				splitter = ":"
			} else {
				splitter = "-"
			}

			record += fmt.Sprintf("%d%s", index+1, splitter)
		}

		// Также выполняется запись самой строки
		record += lines[index]
		*result = append((*result), record)
	}
}
