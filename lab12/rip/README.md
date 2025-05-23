# RIP
Программа написана на C++.
Для сборки программы нужна библиотека `json`, можно поставить с помощью пакетного менеджера `Conan` ([install Conan](https://docs.conan.io/2/installation.html)):

```bash
# Информация о системе сборки и библиотеке указана в conanfile.txt
conan install . --output-folder=build --build=missing
```

Собираем программу:

```bash
cmake -S . -B build -DCMAKE_TOOLCHAIN_FILE="conan_toolchain.cmake" -DCMAKE_BUILD_TYPE=Release
cmake --build build
```

Запускаем:

```bash
./build/rip
```

P.S. Путь до конфига захардкожен, конфиг `rip_config.json` двигать нельзя :)
