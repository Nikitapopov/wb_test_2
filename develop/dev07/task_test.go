package main

import (
	"testing"
	"time"
)

type testCase struct {
	// Слайс чисел. Каждое число - это количество миллисекунд для запуска done-канала
	input []int
	// Нижняя граница времени работы
	lowerLimit int64
	// Верхняя граница времени работы
	upperLimit int64
}

func getDoneChannel(ms time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(ms * time.Millisecond)
	}()
	return c
}

func Test_or_3_sec_succeed_three_ch(t *testing.T) {
	tc := testCase{
		input:      []int{1000, 2000, 3000},
		lowerLimit: 1000,
		upperLimit: 1100,
	}

	doneChannels := make([]<-chan interface{}, 0, len(tc.input))
	for _, msNum := range tc.input {
		doneChannels = append(doneChannels, getDoneChannel(time.Duration(msNum)))
	}

	start := time.Now()
	<-or(doneChannels...)

	execTime := time.Since(start)
	if time.Duration(tc.lowerLimit)*time.Millisecond > execTime || time.Duration(tc.upperLimit)*time.Millisecond < execTime {
		t.Errorf("expected execution time in range: (%v; %v), got: %v", tc.lowerLimit, tc.upperLimit, execTime)
	}
}

func Test_or_3_sec_succeed_one_ch(t *testing.T) {
	tc := testCase{
		input:      []int{1000},
		lowerLimit: 1000,
		upperLimit: 1100,
	}

	doneChannels := make([]<-chan interface{}, 0, len(tc.input))
	for _, secNum := range tc.input {
		doneChannels = append(doneChannels, getDoneChannel(time.Duration(secNum)))
	}

	start := time.Now()
	<-or(doneChannels...)

	execTime := time.Since(start)
	if time.Duration(tc.lowerLimit)*time.Millisecond > execTime || time.Duration(tc.upperLimit)*time.Millisecond < execTime {
		t.Errorf("expected execution time in range: (%v; %v), got: %v", tc.lowerLimit, tc.upperLimit, execTime)
	}
}
