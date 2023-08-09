Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Вывод: error
Функция test возвращает значение nil, но тип задан, как *customError. Из-за этого переменная err в строчке 23 будет иметь значение nil, но тип не nil.
Поэтому условие err != nil положительно.

Как исправить: возвращаемое значение функции test поменять на error или заменить условие 24 строки на "err != (*customError)(nil)"
```
