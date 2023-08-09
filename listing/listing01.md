Что выведет программа? Объяснить вывод программы.

```go
package main

import (
    "fmt"
)

func main() {
    a := [5]int{76, 77, 78, 79, 80}
    var b []int = a[1:4]
    fmt.Println(b)
}
```

Ответ:
``` Программа выведет [77, 78, 79], так как b - это срез массива a от 1 индекса до 4 не включительно
...

```