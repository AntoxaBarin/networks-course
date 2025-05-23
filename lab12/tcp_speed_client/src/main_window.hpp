#pragma once

#include <AUI/Platform/AWindow.h>

#include <cstdint>

class MainWindow: public AWindow {
public:
    MainWindow();
    void sendPackets(const AString& ip, std::uint16_t port, std::uint32_t count);
};
