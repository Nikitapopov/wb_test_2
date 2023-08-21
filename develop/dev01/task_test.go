package main

import (
	mock_main "dev01/mocks"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

// TODO тесты запуска go файла
func Test_task_01_printTime_print_response(t *testing.T) {
	cmd := exec.Command("cmd", "go run task.go")
	// path, _ := exec.LookPath("task.go1")
	// fmt.Println("Path:", path)
	// cmd := exec.Command("ls")
	// cmd := exec.Command("go run /task.go")

	var out strings.Builder
	cmd.Stdout = &out
	var stderr strings.Builder
	cmd.Stderr = &stderr

	// Print the output
	if err := cmd.Run(); err != nil {
		t.Errorf("Ошибка во время запуска задачи: %v", err)
		return
	}
	fmt.Println("i cho:", out.String())
	fmt.Println("i cho2:", stderr.String())

	// got := out.String()

	// _, err = time.Parse(RFC3339Nano, got)
	// if err != nil {
	// 	t.Errorf("Формат возвращаемого значения не совпал с ожидаемым, ошибка парсинга: %v", err)
	// }
}

func Test_task_01_getNow_get_string(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNtpClient := mock_main.NewMockNtpClient(ctrl)

	nowTime := time.Now()

	mockNtpClient.EXPECT().Time().Return(nowTime, nil)
	res, err := getNow(mockNtpClient)
	if err != nil {
		t.Errorf("getNow вернул ошибку: %v", err)
	}

	if nowTime.Format(formatRFC3339Nano) != res {
		t.Errorf("Формат возвращаемого значения не совпал с ожидаемым, ошибка парсинга: %v", err)
	}
}

func Test_task_01_getNow_get_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNtpClient := mock_main.NewMockNtpClient(ctrl)

	mockNtpClient.EXPECT().Time().Return(time.Time{}, errors.New("err"))
	_, err := getNow(mockNtpClient)
	if err == nil {
		t.Error("getNow должен вернуть ошибку")
	}
}
