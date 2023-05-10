#include "skiplist.h"
#include <algorithm>
#include <gtest/gtest.h>
#include <iostream>
#include <spdlog/spdlog.h>

template <typename A, typename B>
void printVector(const std::vector<std::pair<A, B>> &v)
{
    std::cout << "serialization result ";
    for (auto &p : v)
    {
        std::cout << "[" << p.first << "," << p.second << "] ";
    }
    std::cout << std::endl;
}

std::vector<std::pair<int, int>> generateRandomData(size_t amount)
{
    std::vector<std::pair<int, int>> v;
    v.reserve(amount);

    for (size_t i = 0; i < amount; i++)
    {
        v.emplace_back(random() % amount, random() % amount);
    }

    return v;
}

bool isAscending(const std::vector<std::pair<int, int>> &v)
{
    for (auto iter = v.cbegin(); iter != v.cend(); iter++)
    {
        auto next = iter + 1;
        if (next == v.cend())
            break;

        if (iter->first > next->first)
            return false;
    }
    return true;
}

TEST(skiplist, find)
{
    spdlog::set_level(spdlog::level::debug);

    SkipList<int, int> skiplist;

    auto node = skiplist.find(1);
    EXPECT_TRUE(not node);

    skiplist.insert(1, 1);
    node = skiplist.find(1);
    EXPECT_TRUE(node);
    EXPECT_TRUE(node->value_ == 1);
}

TEST(skiplist, insert)
{
    SkipList<int, int> skiplist;

    auto v = generateRandomData(10000);

    for (auto &&p : v)
    {
        skiplist.insert(std::move(p.first), std::move(p.second));
    }

    auto s = skiplist.serialize();
    EXPECT_TRUE(isAscending(s));
}

TEST(skiplist, erase)
{
    SkipList<int, int> skiplist;

    auto v = generateRandomData(10000);

    for (auto &&p : v)
    {
        skiplist.insert(std::move(p.first), std::move(p.second));
    }

    auto s = skiplist.serialize();
    size_t beforeSize = s.size();

    skiplist.erase(v[999].first);

    s = skiplist.serialize();

    EXPECT_TRUE(beforeSize - 1 == s.size());
    EXPECT_TRUE(isAscending(s));
}
