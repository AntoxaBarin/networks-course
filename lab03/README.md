# Практика 3. Прикладной уровень

Программы написаны на Go.

## Сборка сервера

```bash
cd server
go build .
```

## Запуск сервера

```bash
./server <server-port>           # single-thread
./server <server-port> -task=B   # multi-thread 
```
