package main

import (
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	list := []string{"пятак", "пятка", "кот", "лиСток", "тяпка", "столик", "пятак", "слиток"}
	getAnagrams(&list)
}

func getAnagrams(list *[]string) map[string]*[]string {
	// Мапа для промежуточного хранения анаграмм.
	// Ключ - первое вхождение анаграммы в list, значения - все анаграммы в порядке очередности в list кроме первого
	anagrams := map[string]map[string]struct{}{}

	// Мапа для сопоставления ключа - отсортированная по возрастанию строка и значение - первое вхождение анаграммы строки в list
	orderedSymbolsToAnagramKeys := map[string]string{}

	// Заполнение anagrams. Значение anagrams это сет, поэтому одинаковые слова добавляться не будут
	for _, word := range *list {
		// приводим все символы в слове к строчным
		lowercaseWord := strings.ToLower(word)

		// Получение строки с упорядоченным порядком букв слова
		orderedLetters := getOrderedLetters(lowercaseWord)
		anagramKey, ok := orderedSymbolsToAnagramKeys[orderedLetters]
		if !ok {
			orderedSymbolsToAnagramKeys[orderedLetters] = lowercaseWord
			anagrams[lowercaseWord] = map[string]struct{}{}
		} else if word != anagramKey {
			anagrams[anagramKey][lowercaseWord] = struct{}{}
		}
	}

	removeKeyWithEmptyValues(anagrams)

	// Перенос анаграмм в переменную res
	res := map[string]*[]string{}
	for key, set := range anagrams {
		res[key] = &[]string{}
		for word := range set {
			*res[key] = append(*res[key], word)
		}
	}

	// Сортировка анаграмм по возрастанию
	for _, values := range res {
		sort.Strings(*values)
	}

	// Возвращаение объекта анаграмм
	return res
}

// Функция для получения строки с символами в отсортированном по возрастанию порядке
func getOrderedLetters(str string) string {
	letters := strings.Split(str, "")
	sort.Strings(letters)
	return strings.Join(letters, "")
}

// Удаление из списка анаграмм единичных элементов
func removeKeyWithEmptyValues(anagrams map[string]map[string]struct{}) {
	for key, value := range anagrams {
		if len(value) == 0 {
			delete(anagrams, key)
		}
	}
}
