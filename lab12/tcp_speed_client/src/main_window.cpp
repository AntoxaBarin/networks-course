#include "main_window.hpp"

#include <AUI/ASS/Property/BackgroundSolid.h>
#include <AUI/ASS/Property/FontSize.h>
#include <AUI/ASS/Property/TextColor.h>
#include <AUI/Common/AColor.h>
#include <AUI/Platform/APlatform.h>
#include <AUI/Util/UIBuildingHelpers.h>
#include <AUI/View/AButton.h>
#include <AUI/View/ALabel.h>
#include <AUI/View/ATextArea.h>
#include <AUI/View/ASpacerFixed.h>
#include <AUI/Network/ATcpSocket.h>
#include <AUI/Network/AInet4Address.h>

#include <random>

using namespace declarative;

MainWindow::MainWindow() : AWindow("TCP client", 600_dp, 400_dp) {
    const auto ipTextArea = _new<ATextArea>() with_style{FixedSize{200_dp, 35_dp}, FontSize {20_dp}, Border { 2_dp, 0x111111_rgb}};
    auto portTextArea = _new<ATextArea>() with_style{FixedSize{200_dp, 35_dp}, FontSize {20_dp}, Border { 2_dp, 0x111111_rgb}};
    auto countTextArea = _new<ATextArea>() with_style{FixedSize{200_dp, 35_dp}, FontSize {20_dp}, Border { 2_dp, 0x111111_rgb}};

    setContents(Centered{Vertical{
        Centered{Label{"TCP client"} with_style{FontSize{50_dp} }},
        SpacerFixed(25_dp),
        Horizontal{
            ipTextArea,
            Label{"server IP-address"} with_style{FontSize{20_dp}}
        },
        Horizontal{
            portTextArea,
            Label{"port"} with_style{FontSize{20_dp}}
        },
        Horizontal{
            countTextArea,
            Label{"number of packets"} with_style{FontSize{20_dp}}
        },
        SpacerFixed(20_dp),
        Centered{
            _new<AButton>("Send")
            .connect(&AView::clicked, this, [this, ipTextArea, portTextArea, countTextArea] {
                sendPackets(ipTextArea->getText(), portTextArea->getText().toInt().value(), countTextArea->getText().toInt().value());
            }) with_style{FixedSize{160_dp, 40_dp}, FontSize{20_dp}}
        }
    }} with_style { BackgroundSolid {AColor::WHITE}} );
}

void MainWindow::sendPackets(const AString& ip, std::uint16_t port, std::uint32_t count) {
    AInet4Address address(ip, port);
    ATcpSocket socket(address);
    
    std::random_device rd;
    std::uniform_int_distribution<uint8_t> dist(0, 255);
    
    for (auto i = 0; i < count; ++i) {
        std::vector<std::uint8_t> data(1024, 0);
        for (auto& byte : data) {
            byte = dist(rd);
        }
        socket.write(reinterpret_cast<const char*>(data.data()), 1024);
        ALogger::info("Client") << "Packet #" << i << " sent";
    }
}
