cmake_minimum_required(VERSION 3.0)
project(algorithms)

set(CMAKE_CXX_FLAGS
    " -std=c++2a -g -Wall -Wextra -Werror -Wno-unused-parameter -Wshadow ")

include_directories(${PROJECT_SOURCE_DIR})
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

find_package(GTest REQUIRED)  
include_directories(${GTEST_INCLUDE_DIRS}) 

find_package(spdlog REQUIRED)  
include_directories(${SPDLOG_INCLUDE_DIRS})  

enable_testing()

add_subdirectory(skiplist)