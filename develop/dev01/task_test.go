package main

import (
	mock_main "dev01/mocks"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

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
