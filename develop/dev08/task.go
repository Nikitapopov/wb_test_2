package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// TODO конвеер на пайпах
// TODO тесты
func main() {
	startShell()
}

// Функция для запуска шелла
func startShell() {
	// Бесконечое чтение из потока ввода с последующей обработкой ввода в execInput
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// Функция для обработки команды в виде строки input
func execInput(input string) error {
	// Удаление лишних символов
	input = strings.TrimSuffix(input, "\r\n")
	input = strings.TrimPrefix(input, "\r\n")

	// Разделение строки по аргументам
	args := strings.Split(input, " ")

	isFork := false
	if args[len(args)-1] == "&" {
		args = args[:len(args)-1]
		isFork = true
	}

	// Обработка команды cd и выхода из утилиты отдельно от других команд
	switch args[0] {
	case "cd":
		if len(args) != 2 {
			return errors.New("required argument for cd is missing")
		}
		return os.Chdir(args[1])
	case "\\q":
		os.Exit(0)
	}

	// Формирование объекта для представления внешнего процесса команды
	cmd := exec.Command(args[0], args[1:]...)

	// Перенаправление потоков вывода и ошибок
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Выполнение команды
	if isFork {
		return cmd.Start()
	}

	return cmd.Run()
}
