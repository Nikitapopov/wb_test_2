package main

import (
	"net"
	"sync"
	"testing"
	"time"
)

const (
	host = "127.0.0.1"
	port = "8081"
)

func Test_netcatClient(t *testing.T) {
	expected := "Hello, world!"

	client := newClient(udpNetwork, host, port)
	err := client.initConnection()
	if err != nil {
		t.Errorf("error while connection initialization: %v", err)
		return
	}
	defer func() {
		err := client.closeConnection()
		if err != nil {
			t.Errorf("error while connection closing: %v\n", err)
			return
		}
	}()

	// Установка прослушивания порта
	server, err := net.ListenPacket(udpNetwork, net.JoinHostPort(host, port))
	if err != nil {
		t.Errorf("listening: %v\n", err)
		return
	}
	defer server.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	var got []byte
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		got = make([]byte, 128)
		_, _, err = server.ReadFrom(got)
		if err != nil {
			t.Errorf("Unexpected error while server reading message: %v", err)
		}
	}(&wg)

	time.Sleep(100 * time.Millisecond)

	err = client.sendMsg(expected)
	if err != nil {
		t.Errorf("error while writing message: %v\n", err)
	}

	wg.Wait()

	expectedRunesLen := len([]rune(expected))
	gotInRunes := []rune(string(got))
	gotInStr := string(gotInRunes[:expectedRunesLen])

	if gotInStr != expected {
		t.Errorf("Expected: %s, got: %s\n", expected, got)
	}
}
