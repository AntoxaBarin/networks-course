# Практика 7. Транспортный уровень

UDP сервер и клиент написаны на Go.

### Сборка и запуск

`port` и `server-port` должны совпадать, значение порта по умолчанию`8080`.

```bash
cd server
go run . [-port]

cd ../client
go run . [-server-port]
```
