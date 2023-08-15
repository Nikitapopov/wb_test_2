package main

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

import (
	"os"
	"time"

	"github.com/beevik/ntp"
)

const (
	ntpServerAddress = "0.beevik-ntp.pool.ntp.org"
	RFC3339Nano      = "2006-01-02T15:04:05.999999999Z07:00"
)

func main() {
	ntpClient := newNtpClient()
	// printNow(ntpClient)
	time, err := getNow(ntpClient)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	os.Stdout.WriteString(time)
}

// // Функция для отображения текущего времени в формате RFC3339Nano c помощью клиента ntpClient в stdout
// func printNow(ntpClient iNtpClient) {
// 	time, err := getNow(ntpClient)
// 	if err != nil {
// 		os.Stderr.WriteString(err.Error())
// 		os.Exit(1)
// 	}
// 	os.Stdout.WriteString(time)
// }

// Функция для получения текущего времени в формате RFC3339Nano с помощью клиента ntpClient
func getNow(ntpClient iNtpClient) (string, error) {
	time, err := ntpClient.Time()
	if err != nil {
		return "", err
	}

	return time.Format(RFC3339Nano), nil
}

// ntp клиент

type iNtpClient interface {
	Time() (time.Time, error)
}

type ntpClient struct{}

func newNtpClient() iNtpClient {
	return &ntpClient{}
}

func (c *ntpClient) Time() (time.Time, error) {
	return ntp.Time(ntpServerAddress)
}
