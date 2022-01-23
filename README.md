### Описание
**gomocker** - это небольшой генератор кода для создания интерфейсов, используя структуры Golang (как правило, редакторы кода и IDE поддерживают автозаполнение полей структур, что упрощает создание мокированой функции)
### Установка
```shell
go install github.com/gomocker/gomocker@latest
```
### Создание тестового конфига
```shell
gomocker touch
```
Команда выше создаст файла *gomocker.yml* в рабочей директории.
### Структура конфигурационного файла
```yaml
# gomocker.yml

# название пакета, которое будет записано в сгенерированный файл
package: main

# файл, куда записать сгенерированный код
output: gomocker_output.go

# интерфейсы, для которых нужно сгенерировать моки
mocks:
  # указывается полный импорт пакета (например, github.com/minio/minio-go/v7, или в случае с io просто io)
  # и список интерфейсов из этого пакета, для которых нужно сгенерировать код.
  io:
    - Reader
    - Writer
    - ReadWriter
  math/rand:
    - Source

# опциональная секция для правильной настройки имортов в сгенерированном коде
# например, если у нас где-то используется структура minio.Client, то
# gomocker будет пытаться импортировать minio как github.com/minio/minio-go
# в то время как правильный импорт будет github.com/minio/minio-go/v7
imports:
  io: io
  rand: math/rand
  minio: github.com/minio/minio-go/v7
```
### Генерация
```shell
gomocker
```
Команда выше сгенерирует код для создания моков и запишет его в файл, указанный в *gomocker.yml*

Gomocker генерирует конструкторы для создания интерфейсов.

На вход конструктору передается структура, все поля которой являются фукнциями аналогичными фукнциям интерфейса
### Демонстрация
Представим проект **test**, со следующим содержимым
#### go.mod
```golang
module test

go 1.17
```
#### main.go
```golang
package main

import (
	"fmt"
)

type Auth interface {
	Login(login string, password string) (err error)
}
```
#### gomocker.yml
```yaml
package: main

output: output_gomocker.go

mocks:
  test:
    - Auth
```
#### Генерация
```shell
gomocker
```
#### Использование в main.go
```golang
package main

import (
	"fmt"
)

type Auth interface {
	Login(login string, password string) (err error)
}

func main() {
	auth := NewAuthMock(AuthBehavior{
		Login: func(login string, password string) (err error) {
			fmt.Println("Функия Login была вызана")
			return nil
		},
	})

	if err := auth.Login("admin", "admin"); err != nil {
		fmt.Println(err)
	}
}
```
### Плюшки
- [x] Поддержка алиасов
```golang
type Test1 interface {
  Method1()
}

type Test2 = Test1

type Test3 = Test2

type Test = Test3
```
- [x] Поддержка вложенных интерфейсов
```golang
type Test1 interface {
	Method1()
}

type Test2 interface {
	Method2()
}

type Test interface {
	Test1
	Test2
}
```
- [x] Не используется рефлексия и interface{}. Аргументы функций сохраняют свой тип
- [x] Если аргумент функции именован, то имя также сохранится.
