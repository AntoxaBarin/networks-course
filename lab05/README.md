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

### Сборка и запуск задания 2

```bash
cd src/smtp_client
go build .
./smtp_client --from <sender email> --to <recepient email> --pass <sender email password> --msg <one_word_message_like_this_,_yes_stdlib_cant_handle_slice_of_strings_as_flag_value>
```

### Сборка и запуск задания 3

```bash
cd src/smtp_client
go build .
./smtp_client --from <sender email> --to <recepient email> --pass <sender email password> --msg <message> --image <path to image>
```



