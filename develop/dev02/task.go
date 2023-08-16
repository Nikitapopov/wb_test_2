package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// Исходная строка
	str := "a4bc2d5e"

	// Выполнение парсинга строки
	res, err := parseStr(str)
	if err != nil {
		fmt.Printf("Ошибка в parseStr: %v", err)
		return
	}

	fmt.Println(res)
}

// Функция для парсинга строки str.
// Если после очередного символа имеется натуральное число, то символ записывается n раз вместо одного.
// Предшествующий экранирующий символ '\' обрабатывает последущий символ '\' или цифру, как обычный символ.
// Обычные символы записываются в результат без изменений.
func parseStr(str string) (string, error) {
	// При итерации по строке содержит символ предшествующий текущему
	var prevElem rune

	// Слайс рун, в который помещается распарсенная по правилам парсинга данной функции строка str
	res := make([]rune, 0, len(str))

	// Слайс рун, в котором находится число повторений предшествующего символа
	var curNum []rune

	// Флаг, что предыдущий символ '\'
	isPrevSymbolEscape := false

	// Итерация по строке
	for i, value := range str {
		// Если предыдущий символ был '\', то запись любого символа, как обычного. Продолжение итерации.
		if isPrevSymbolEscape {
			prevElem = value
			isPrevSymbolEscape = false
			continue
		}

		// Если текущее значение цифра, то добавление этого значения в curNum. Продолжение итерации.
		if unicode.IsDigit(value) {
			if prevElem == 0 {
				return "", errors.New("некорректная строка")
			}
			curNum = append(curNum, value)
			continue
		}

		// Если текущее значение '\', то проставление флага isPrevSymbolEscape
		if value == '\\' {
			isPrevSymbolEscape = true
		}

		// Сохранение символа происходит на шаг позже, когда он становится предшествующим текущему, поэтому пропускаем первый символ.
		if i > 0 {
			// Определяем сколкьо раз необходимо записать символ
			appendingNum, err := getAppendingNumber(curNum)
			if err != nil {
				return "", err
			}

			// Запись символа в строку количеством appendingNum раз
			for i := 0; i < appendingNum; i++ {
				res = append(res, prevElem)
			}

			// Обнуление числа повторений
			curNum = []rune{}
		}

		// Назначение текущего элемента на предыдущий
		prevElem = value
	}

	// Строка не может закончиться одинарным '\'
	if isPrevSymbolEscape {
		return "", errors.New("некорректная строка")
	}

	// Обработка последнего символа
	if prevElem != 0 {
		appendingNum, err := getAppendingNumber(curNum)
		if err != nil {
			return "", err
		}

		for i := 0; i < appendingNum; i++ {
			res = append(res, prevElem)
		}
	}

	return string(res), nil
}

// Функция для конвертирования слайса рун slc, заполненного цифрами в целочисленное число. Если длина slc = 0, то возвращает 1
func getAppendingNumber(slc []rune) (int, error) {
	if len(slc) == 0 {
		return 1, nil
	}

	num, err := strconv.Atoi(string(slc))
	if err != nil {
		return 0, errors.New("неправильная логика обработки чисел")
	}
	return num, nil
}
