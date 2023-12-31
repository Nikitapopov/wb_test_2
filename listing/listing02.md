Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
В функции test объявляется переменная x с нулевым значением, равным 0. 
В строчке 12 в стек отложенные функций кладется функция, которая инкрементирует x.
Далее x приравнивается к 1 в строчке 15. Перед возвратом значения из функции вызывается функция из стека отложенных функций, где значение инкрементируется. Так как отложенные функции могут менять возвращаемые значения, если они заданы, как именованные возвращаемые параметры, то значение возвращается инкрементированным.
В функции anotherTest объявляется переменная x с нулевым значением, равным 0.
В строчке 12 в стек отложенные функций кладется функция, которая инкрементирует x.
Далее x приравнивается к 1 в строчке 25. Перед возвратом значения из функции вызывается функция из стека defer, где значение инкрементируется, но на возвращаемое значение это уже не влияет. Возвращается значение 1.

Принцип работы defer-ов:
Функция вызванная с ключевым словом defer "запоминает" параметры, с которыми она вызывается и кладется в стек отложенных функций. Перед тем, как функция вернет возвращаемое значение, сработают все функции из стека в порядке LIFO.
```
