add_executable(skipListTest skiplist_test.cc)
target_link_libraries(skipListTest gtest gtest_main spdlog::spdlog)

add_test(skiplist::test skipListTest)