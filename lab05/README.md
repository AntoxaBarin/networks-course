# Практика 5. Прикладной уровень

## Программирование сокетов.

Программа для задачи 1 написана на Go + [gomail](https://github.com/go-gomail/gomail).

### Сборка и запуск задания 1

```bash
cd src/simple_smtp
go mod tidy
go build .
./smtp --from <sender email> --to <recepient email> --pass <sender email password> --type <html | text>
```


