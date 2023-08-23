package main

import (
	"fmt"
	"net"
)

// Константы конфигурации сервера: тип сети и адрес
const (
	udpNetwork = "udp"
	address    = "127.0.0.1:8081"
)

func main() {
	fmt.Println("Start server...")

	// Устанавка прослушивания порта
	server, err := net.ListenPacket(udpNetwork, address)
	if err != nil {
		fmt.Printf("listening: %v\n", err)
		return
	}
	defer server.Close()

	fmt.Println("Server started!")

	// Чтение данных от клиента в буфер и вывод
	for {
		buf := make([]byte, 128)
		_, _, err := server.ReadFrom(buf)
		if err != nil {
			continue
		}
		fmt.Println(string(buf))
	}
}
