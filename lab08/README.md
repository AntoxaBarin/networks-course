# Практика 8. Транспортный уровень

Клиент и сервер написаны на Go.

### Сборка и запуск

```bash
# Сервер

cd src/server
go build .
echo "" > received.txt

./server -v 1 # task A
./server -v 2 -port <server port> -client-port <client port> -data <path to payload> -timeout <timeout duration> # task B

# Клиент

cd src/client
go build .
echo "" > received.txt

./client -server-port <server port> -data <path to payload> -timeout <timeout duration> # task A
./client -port <client port> -server-port <server port> -data <path to payload> -timeout <timeout duration> # task B
```
