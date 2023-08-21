package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const (
	// Файл по умолчанию для входных данных
	defaultInputFilename = "./data/input.txt"
	// Файл для вывода отсортированных строк
	outputFilename = "./data/output.txt"
)

func main() {
	// Чтение аргументов командной строки
	args := cmdArgs{}
	args.parseArgs()

	// Получение исходных строк из файла
	lines, err := readFile(args.f)
	if err != nil {
		errMsg := fmt.Sprintf("reading file: %v", err)
		os.Stderr.WriteString(errMsg)
		os.Exit(1)
	}

	// Сортировка строк
	sortedLines := sortLines(lines, args)

	// Конкатенация отсортированных строк
	result := strings.Join(sortedLines, "\n")

	// Сохранение строк в файл
	err = saveResult(result, outputFilename)
	if err != nil {
		msg := fmt.Sprintf("Saving result: %v", err)
		os.Stderr.WriteString(msg)
		os.Exit(1)
	}
}

// Структура флагов сортировки строк файла
type cmdArgs struct {
	// Колонка для сортировки
	k int
	// Сортировка по числовому значению
	n bool
	// Сортировка в обратном порядке
	r bool
	// Игнорирование повторяющихся строк
	u bool
	// Название файла для сортировки
	f string
}

// Метод для инициализирования аргументов
func (a *cmdArgs) parseArgs() {
	flag.IntVar(&a.k, "k", 1, "specifying the column to sort")
	flag.BoolVar(&a.n, "n", false, "sort by numeric value")
	flag.BoolVar(&a.r, "r", false, "sort in reverse order")
	flag.BoolVar(&a.u, "u", false, "do not output duplicate lines")
	flag.StringVar(&a.f, "f", defaultInputFilename, "file name")
	flag.Parse()
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

// Метод для сортировки строк lines согласно флагам args
func sortLines(lines []string, args cmdArgs) []string {
	// Если задан флаг u, то применяется отфильтровывание уникальных значений
	if args.u {
		lines = getUnique(lines)
	}

	// Функция для сравнения числовых или строковых значений
	var firstStringLess func(string, string) bool

	// Если задан флаг n, то применяется сортировка по числовым значениям
	if args.n {
		firstStringLess = func(value1, value2 string) bool {
			leftNum, lErr := strconv.Atoi(value1)
			rightNum, rErr := strconv.Atoi(value2)

			if lErr != nil && rErr != nil {
				return leftNum < rightNum
			}

			if lErr != nil || rErr != nil {
				return lErr == nil
			}

			return leftNum < rightNum
		}
	} else { // Иначе сортировка по строковым значениям
		firstStringLess = func(value1, value2 string) bool {
			return value1 < value2
		}
	}

	// Метод для сравнения значений по индексам в lines
	sortLess := func(i, j int) bool {
		// Если одна из строк нулевая, то при сравнении она меньше
		if len(lines[i]) == 0 {
			return true
		}
		if len(lines[j]) == 0 {
			return false
		}

		// Разделение строк по колонкам
		leftWords := strings.Split(lines[i], " ")
		rightWords := strings.Split(lines[j], " ")

		// Если в левой строке отсутствует количество колонок, заданное флагом k,
		// а в правой оно есть, то при сравнении левая меньше
		if len(leftWords) < args.k && len(rightWords) >= args.k {
			return true
		}

		// Если в правой строке отсутствует количество колонок, заданное флагом k,
		// а в левой оно есть, то при сравнении правая меньше
		if len(leftWords) >= args.k && len(rightWords) < args.k {
			return false
		}

		// Если левая и правая строки не имеют k колонок, то сравнение выполняется по всей строке
		if len(leftWords) < args.k && len(rightWords) < args.k {
			// Сравнение выполняется по ранее заданной функции, т.е. по числовому или строковому значению
			return firstStringLess(lines[i], lines[j])
		}

		// Если левая и правая строки имеют по крайней мере по k колонок,
		// то сравнение выполняется по k-ой (и последующим) колонкам по числовому или строковому значению
		leftWordsSinceKCol := strings.Join(leftWords[args.k-1:], " ")
		rightWordsSinceKCol := strings.Join(rightWords[args.k-1:], " ")
		return firstStringLess(leftWordsSinceKCol, rightWordsSinceKCol)
	}

	// Выполнение сортировки строк
	// Если задан флаг r, то сортировка выполняется в обратном порядке
	if args.r {
		sort.Slice(lines, func(i, j int) bool { return !sortLess(i, j) })
	} else {
		sort.Slice(lines, sortLess)
	}

	return lines
}

// Метод для получения уникальных строк
func getUnique(lines []string) []string {
	// Конечный слайс с уникальными значениями
	uniqueLines := make([]string, 0, len(lines)/2)

	// Сет для хранения уже имеющихся строк
	linesSet := make(map[string]struct{}, len(lines)/2)
	for _, line := range lines {
		_, ok := linesSet[line]
		if !ok {
			linesSet[line] = struct{}{}
			uniqueLines = append(uniqueLines, line)
		}
	}

	return uniqueLines
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
