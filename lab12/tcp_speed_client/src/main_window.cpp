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

using namespace declarative;

MainWindow::MainWindow() : AWindow("TCP client", 600_dp, 400_dp) {
    setContents(Centered{Vertical{
        Centered{Label{"TCP client"} with_style{FontSize{50_dp} }},
        SpacerFixed(25_dp),
        Horizontal{
            _new<ATextArea>() with_style{FixedSize{200_dp, 35_dp}, FontSize {20_dp}, Border { 2_dp, 0x111111_rgb }},
            Label{"server IP-address"} with_style{FontSize{20_dp}}
        },
        Horizontal{
            _new<ATextArea>() with_style{FixedSize{200_dp, 35_dp}, FontSize {20_dp}, Border { 2_dp, 0x111111_rgb }},
            Label{"port"} with_style{FontSize{20_dp}}
        },
        Horizontal{
            _new<ATextArea>() with_style{FixedSize{200_dp, 35_dp}, FontSize {20_dp}, Border { 2_dp, 0x111111_rgb }},
            Label{"number of packets"} with_style{FontSize{20_dp}}
        },
        SpacerFixed(20_dp),
        Centered{
            _new<AButton>("Send")
            .connect(&AView::clicked, this, [] {
                APlatform::openUrl("https://github.com/aui-framework/aui");
            }) with_style{FixedSize{160_dp, 40_dp}, FontSize{20_dp}}
        }
    }} with_style { BackgroundSolid {AColor::WHITE}} );
  }
  
