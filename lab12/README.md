# TCP speed

Клиент и сервер написаны на C++.

Сборка:

```bash
cd tcp_speed_client
mkdir build
cmake ..
cmake --build .
```

Если не получилось, скорее всего, проблема в отсутствии библиотек:

```bash
sudo apt install pkg-config libfontconfig-dev libxcursor-dev libxi-dev libxrandr-dev libglew-dev libstdc++-static libpulse-dev libdbus-1-dev
```

Запуск:

```bash
./bin/tcp_speed_client
```

Точно так же собирается сервер:

```bash
cd tcp_speed_server
mkdir build
cmake ..
cmake --build .
./bin/tcp_speed_server
```
