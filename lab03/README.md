# Практика 3. Прикладной уровень

Программы написаны на Go.

## Сборка сервера

```bash
cd server
go build .
```

## Запуск сервера

```bash
./server <server-port>                             # single-thread
./server <server-port> -task=B                     # not limited multi-thread 
./server <server-port> <concurrency-level> -task=D # limited multi-thread 
```

## Сборка клиента

```bash
cd client
go build .
```

## Запуск клиента
```bash
./client <server-host> <server-port> <filename>
```
