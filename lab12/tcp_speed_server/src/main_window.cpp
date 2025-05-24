#include "main_window.hpp"

#include <AUI/ASS/Property/BackgroundSolid.h>
#include <AUI/ASS/Property/FontSize.h>
#include <AUI/ASS/Property/TextColor.h>
#include <AUI/Common/AColor.h>
#include <AUI/Common/AString.h>
#include <AUI/Logging/ALogger.h>
#include <AUI/Platform/APlatform.h>
#include <AUI/Util/UIBuildingHelpers.h>
#include <AUI/View/AButton.h>
#include <AUI/View/ALabel.h>
#include <AUI/View/ATextArea.h>
#include <AUI/View/ASpacerFixed.h>
#include <AUI/Network/ATcpServerSocket.h>
#include <AUI/Network/AInet4Address.h>

#include <chrono>
#include <utility>

using namespace declarative;

MainWindow::MainWindow() : AWindow("TCP server", 600_dp, 400_dp) {
    auto portTextArea = _new<ATextArea>() with_style{FixedSize{200_dp, 35_dp}, FontSize {20_dp}, Border { 2_dp, 0x111111_rgb}};
    auto bytesTextArea = _new<ATextArea>() with_style{FixedSize{200_dp, 35_dp}, FontSize {20_dp}, Border { 2_dp, 0x111111_rgb}};
    auto speedTextArea = _new<ATextArea>() with_style{FixedSize{200_dp, 35_dp}, FontSize {20_dp}, Border { 2_dp, 0x111111_rgb}};

    setContents(Centered{Vertical{
        Centered{Label{"TCP server"} with_style{FontSize{50_dp} }},
        SpacerFixed(25_dp),
        Horizontal{
            portTextArea,
            Label{"port"} with_style{FontSize{20_dp}}
        },
        Horizontal{
            speedTextArea,
            Label{"Speed (bytes/s)"} with_style{FontSize{20_dp}}
        },
        Horizontal{
            bytesTextArea,
            Label{"Bytes received"} with_style{FontSize{20_dp}}
        },
        SpacerFixed(20_dp),
        Centered{
            _new<AButton>("Listen")
            .connect(&AView::clicked, this, [this, speedTextArea, portTextArea, bytesTextArea] {
                auto [speed, bytes] = listen(portTextArea->getText().toInt().value());
                speedTextArea->setText(AString(std::to_string(speed)));
                bytesTextArea->setText(AString(std::to_string(bytes)));
            }) with_style{FixedSize{160_dp, 40_dp}, FontSize{20_dp}}
        }
    }} with_style { BackgroundSolid {AColor::WHITE}} );
}

std::pair<double, std::uint32_t> MainWindow::listen(std::uint16_t port) {
    std::vector<std::uint8_t> buffer(1024);
    ALogger::info("TCP SERVER") << "Start listening...";

    ATcpServerSocket socket(port);
    auto conn = socket.accept();

    ALogger::info("TCP SERVER") << "Start handling connection...";

    bool is_first_packet = true;
    std::uint64_t sending_timestamp = 0;
    std::uint32_t total_bytes_received = 0;

    while (true) {
        auto bytes_read = conn->read(reinterpret_cast<char*>(buffer.data()), 1024);
            if (bytes_read == 0) {
                break;
            }
            ALogger::info("TCP SERVER") << "Read " << bytes_read << " bytes";
            // Get timestamp from first packet
            if (is_first_packet) {
                std::memcpy(&sending_timestamp, buffer.data(), sizeof(sending_timestamp));
                is_first_packet = false;
            }
            total_bytes_received += bytes_read;
    }

    std::uint64_t timestamp = std::chrono::duration_cast<std::chrono::milliseconds>(
        std::chrono::system_clock::now().time_since_epoch()
    ).count();
    std::uint64_t delta = timestamp - sending_timestamp;
    ALogger::info("TCP SERVER") << "TIMESTAMP: " << timestamp;
    ALogger::info("TCP SERVER") << "Delta time: " << delta << " ms";
            
    return {total_bytes_received / delta * 1000, total_bytes_received};
}
