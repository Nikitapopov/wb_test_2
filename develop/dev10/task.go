package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

const (
	tcpNetwork = "tcp"
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

	// Запуск telnet-клиента
	startClient(args)
}

// Структура флагов командной строки
type cmdArgs struct {
	// Таймаут на подключение к серверу в секундах
	timeout int
	// IP или доменное имя
	host string
	// Порт
	port string
}

// Метод для инициализирования аргументов
func (a *cmdArgs) parseArgs() error {
	// Привязка аргументов к полям структуры
	flag.IntVar(&a.timeout, "timeout", 10, "Таймаут на подключение к серверу в секундах")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		return errors.New("required arguments are missing")
	}

	a.host = args[0]
	a.port = args[1]

	return nil
}

// Интерфейс telnet-клиента
type iClient interface {
	initConnection() error
	closeConnection() error
	receiveMsg() (string, error)
	sendMsg(msg string) error
}

// Структура telnet-клиента
type client struct {
	address    string
	timeout    int
	conn       net.Conn
	connReader *bufio.Reader
}

// Конструктор telnet-клиента
func newClient(address string, timeout int) iClient {
	return &client{
		address: address,
		timeout: timeout,
	}
}

// Инициализация соединения telnet-клиента и сервера
func (c *client) initConnection() error {
	conn, err := net.DialTimeout(tcpNetwork, c.address, time.Duration(c.timeout)*time.Second)
	if err != nil {
		return err
	}

	c.conn = conn
	c.connReader = bufio.NewReader(conn)

	return nil
}

// Закрытие соединения telnet-клиента
func (c *client) closeConnection() error {
	return c.conn.Close()
}

// Чтение сообщения
func (c *client) receiveMsg() (string, error) {
	msg, err := c.connReader.ReadString('\n')
	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
			err = fmt.Errorf("connection closed")
		}
		return "", err
	}

	msg = strings.TrimSuffix(msg, "\n")
	return msg, nil
}

// Отправка сообщения
func (c *client) sendMsg(msg string) error {
	if _, err := c.conn.Write([]byte(msg)); err != nil {
		return err
	}

	return nil
}

// Запуск telnet-клиента
func startClient(args cmdArgs) {
	// Создание клиента
	telnetClient := newClient(net.JoinHostPort(args.host, args.port), args.timeout)

	// Подсоединение к серверу
	err := telnetClient.initConnection()
	if err != nil {
		fmt.Printf("error while connection initialization: %v\n", err)
		return
	}
	defer func() {
		err = telnetClient.closeConnection()
		if err != nil {
			fmt.Printf("error while connection closing: %v\n", err)
		}
	}()

	// Канал, закрытие которого влияет на закрытие клиента
	shutdown := make(chan struct{}, 1)

	// Канал для отслеживания сигнала об окончании работы
	signalToShutdown := make(chan os.Signal, 1)
	signal.Notify(signalToShutdown, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	defer close(signalToShutdown)

	// При получении сигнала об окончании выполняется закрытие клиента
	go func() {
		<-signalToShutdown
		fmt.Println("Client is closing...")
		close(shutdown)
	}()

	// Горутина для чтения сообщений
	go func() {
		for {
			select {
			// При закрытии клиента возврат из горутины
			case <-shutdown:
				return
			default:
				// Получение сообщения
				msg, err := telnetClient.receiveMsg()
				if err != nil {
					// Если ошибка указывает на остановку работы сервера, то отправка сигнала об окончании работы клиента и выход из горутины
					if err.Error() == "connection closed" {
						fmt.Println("Server closed")
						signalToShutdown <- syscall.SIGQUIT
						return
					}

					fmt.Printf("Error in receiving message: %v\n", err)
					continue
				}

				fmt.Printf("Message from server: %s\n", msg)
			}
		}
	}()

	// Горутина для отправки сообщений серверу
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			// При закрытии клиента возврат из горутины
			case <-shutdown:
				return
			default:
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println(err)
				}

				if err = telnetClient.sendMsg(input); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	<-shutdown
}
