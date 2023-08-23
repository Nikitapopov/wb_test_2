package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	tcpNetwork = "tcp"
	address    = "127.0.0.1:8081"
)

func main() {
	fmt.Println("Start server...")

	// Устанавка прослушивания порта
	ln, err := net.Listen(tcpNetwork, address)
	if err != nil {
		fmt.Printf("listening: %v\n", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server started!")

	// Запуск цикла обработки входящих данных
	for {
		// Принятие входящего соединения
		conn, _ := ln.Accept()

		// Обработка событий соединения
		go process(conn)
	}
}

// Функция для обработки событий соединения
func process(conn net.Conn) {
	// Закрытие соединение перед выходом из функции
	defer conn.Close()

	// Переменная для чтения входящих данных
	connReader := bufio.NewReader(conn)
	for {
		msg, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("error while reading message: %v\n", err.Error())
			break
		}

		// Распечатывание полученого сообщения
		fmt.Print("Message Received:", msg)

		// Отправка нового сообщения обратно клиенту
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Printf("error while writing: %v\n", err.Error())
			break
		}
	}
}
