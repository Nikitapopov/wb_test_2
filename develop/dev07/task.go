package main

import (
	"fmt"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

// Функция для генерации done-канала
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

// Функция для объединения несколькоих каналов-читателей в один
func or(channels ...<-chan interface{}) <-chan interface{} {
	// Возвращаемый канал
	resCh := make(chan interface{})

	// Запускаем для каждого канала горутину, в которой отслеживается закрытие канала.
	// Когда один из каналов закрывается, закрывается и возвращаемый канал.
	for _, ch := range channels {
		go func(ch <-chan interface{}) {
			for {
				select {
				case _, ok := <-ch:
					if !ok {
						close(resCh)
						return
					}
				case _, ok := <-resCh:
					if !ok {
						return
					}
				}
			}
		}(ch)
	}

	return resCh
}

func main() {
	start := time.Now()
	<-or(
		sig(2*time.Second),
		sig(3*time.Second),
		sig(1*time.Second),
		sig(2*time.Second),
		sig(4*time.Second),
	)

	fmt.Printf("fone after %v", time.Since(start))
}
