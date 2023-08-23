package main

import (
	"bufio"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_telnetClient(t *testing.T) {
	const (
		host    = "127.0.0.1"
		port    = 8081
		timeout = 1 * time.Second
		msg     = "Hello, world!"
	)

	var got1 string
	var got2 string

	// Запуск сервера
	ln, err := net.Listen(tcpNetwork, net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		t.Errorf("listening: %v\n", err)
		return
	}
	ln.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	var serverConn net.Conn
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		serverConn, _ = ln.Accept()

		serverConnReader := bufio.NewReader(serverConn)
		got2, err = serverConnReader.ReadString('\n')
		if err != nil {
			t.Errorf("error while reading message: %v\n", err.Error())
			return
		}

		_, err = serverConn.Write([]byte(msg + "\n"))
		if err != nil {
			t.Errorf("error while writing: %v\n", err.Error())
			return
		}
	}(&wg)

	defer func(serverConn *net.Conn) {
		(*serverConn).Close()
	}(&serverConn)

	time.Sleep(100 * time.Millisecond)

	// Запуск клиента
	telnetClient := newClient(net.JoinHostPort(host, strconv.Itoa(port)), int(timeout))
	err = telnetClient.initConnection()
	if err != nil {
		t.Errorf("error while connection initialization: %v\n", err)
		return
	}
	defer func() {
		err = telnetClient.closeConnection()
		if err != nil {
			t.Errorf("error while connection closing: %v\n", err)
			return
		}
	}()

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		receivedMsg, err := telnetClient.receiveMsg()
		if err != nil {
			t.Errorf("Error in receiving message: %v\n", err)
			return
		}

		got1 = receivedMsg
	}(&wg)

	time.Sleep(100 * time.Millisecond)

	if err = telnetClient.sendMsg(msg + "\n"); err != nil {
		t.Errorf("sending message: %v\n", err)
		return
	}

	wg.Wait()

	if got1 != msg {
		t.Errorf("Receiving message in client. Expected: %s, got: %s", msg, got1)
		return
	}

	if got2 != msg+"\n" {
		t.Errorf("Receiving message in server. Expected: %s, got: %s", msg, got2)
		return
	}
}
