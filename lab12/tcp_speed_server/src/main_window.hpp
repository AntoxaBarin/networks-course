#pragma once

#include <AUI/Platform/AWindow.h>

#include <cstdint>

class MainWindow: public AWindow {
public:
    MainWindow();
    std::pair<double, std::uint32_t> listen(std::uint16_t port);
};
