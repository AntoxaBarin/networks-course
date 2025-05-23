cmake_minimum_required(VERSION 3.16)
project(tcp_speed_client)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

set(AUI_VERSION v7.0.1)

file(
        DOWNLOAD
        https://raw.githubusercontent.com/aui-framework/aui/${AUI_VERSION}/aui.boot.cmake
        ${CMAKE_CURRENT_BINARY_DIR}/aui.boot.cmake)
include(${CMAKE_CURRENT_BINARY_DIR}/aui.boot.cmake)

auib_import(aui https://github.com/aui-framework/aui
        COMPONENTS core views network
        VERSION ${AUI_VERSION})

aui_executable(${PROJECT_NAME})

aui_link(${PROJECT_NAME} PRIVATE aui::core aui::views aui::network)
