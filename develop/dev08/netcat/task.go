package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
)

/*
Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const (
	tcpNetwork = "tcp"
	udpNetwork = "udp"
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

	// Запуск netcat клиента и работа с ним
	err = execute(args)
	if err != nil {
		errMsg := fmt.Sprintf("error in execution: %v", err)
		os.Stderr.WriteString(errMsg)
		os.Exit(1)
	}
}

// Структура флагов командной строки
type cmdArgs struct {
	// Использовать UDP вместо TCP
	useUDP bool
	// IP или доменное имя
	host string
	// Порт
	port string
}

// Метод для инициализирования аргументов
func (a *cmdArgs) parseArgs() error {
	flag.BoolVar(&a.useUDP, "u", false, "Использовать UDP вместо TCP")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		return errors.New("required arguments are missing")
	}

	a.host = args[0]
	a.port = args[1]

	return nil
}

// Интерфейс netcat-клиента
type iClient interface {
	initConnection() error
	closeConnection() error
	sendMsg(msg string) error
}

// Структура netcat-клиента
type client struct {
	typeNetwork string
	address     string
	conn        net.Conn
}

// Конструктор netcat-клиента
func newClient(typeNetwork string, host, port string) iClient {
	return &client{
		typeNetwork: typeNetwork,
		address:     net.JoinHostPort(host, port),
	}
}

// Инициализация соединения netcat-клиента и сервера
func (c *client) initConnection() error {
	conn, err := net.Dial(c.typeNetwork, c.address)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// Закрытие соединения netcat-клиента
func (c *client) closeConnection() error {
	return c.conn.Close()
}

// Отправка сообщения
func (c *client) sendMsg(msg string) error {
	if _, err := c.conn.Write([]byte(msg)); err != nil {
		return err
	}

	return nil
}

// Функция запуска netcat-клиента, чтения сообщений из потока ввода и отправка серверу
func execute(args cmdArgs) error {
	// Объявление типа сети
	typeNetwork := tcpNetwork
	if args.useUDP {
		typeNetwork = udpNetwork
	}

	// Создание netcat-клиента
	client := newClient(typeNetwork, args.host, args.port)

	// Инициализация соединение между клиентом и сервером
	err := client.initConnection()
	if err != nil {
		return fmt.Errorf("error while connection initialization: %v", err)
	}
	defer func() {
		err := client.closeConnection()
		if err != nil {
			fmt.Printf("error while connection closing: %v\n", err)
		}
	}()

	// Бесконечное чтение из потока ввода и отправка сообщений серверу
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		// При вводе "\q" выполняется выход из утилиты
		if input == "\\q" {
			break
		}

		// Отправка сообщений серверу
		err = client.sendMsg(input)
		if err != nil {
			fmt.Printf("error while writing message: %v\n", err)
		}
	}

	return nil
}
