#pragma once
#include "../3rdparty/magic_enum/magic_enum.hpp"
#include <assert.h>
#include <cstdlib>
#include <functional>
#include <iostream>
#include <memory>
#include <utility>
#include <vector>

template <typename Key, typename Value>
struct Node : public std::enable_shared_from_this<Node<Key, Value>>
{
    Node(Key &&key, Value &&value, size_t layer)
        : key_(key), value_(value), forward_(layer, nullptr)
    {
    }
    ~Node() = default;

    size_t height() const
    {
        return forward_.size();
    }

    void update(size_t h, std::shared_ptr<Node> node)
    {
        forward_.at(h - 1) = node;
    }

    std::shared_ptr<Node> forward(size_t h)
    {
        if (h > height())
            return nullptr;

        return forward_.at(h - 1);
    }

    Key key_;
    Value value_;

    std::vector<std::shared_ptr<Node>> forward_;
};

template <typename Key, typename Value>
class SkipList
{
  public:
    enum class CompareResult
    {
        Equal,
        Bigger,
        Smaller
    };

    enum class InsertType
    {
        Update,
        Insert
    };

    enum class EraseType
    {
        NoThisElement,
        HasThisElement
    };

    typedef Node<Key, Value> NodeType;
    typedef std::function<CompareResult(const Key &, const Key &)> Comparator;

  public:
    SkipList()
        : elementCount_(0),
          head_(std::make_shared<NodeType>(Key(), Value(), maxLayer_))
    {
    }

    SkipList(Comparator &&comparator)
        : comparator_(comparator), elementCount_(0),
          head_(std::make_shared<NodeType>(Key(), Value(), maxLayer_))
    {
    }

    SkipList(SkipList &) = delete;
    SkipList(SkipList &&) = delete;

  public:
    InsertType insert(Key &&key, Value &&value);
    std::shared_ptr<NodeType> find(const Key &key) const;
    EraseType erase(const Key &key);

    std::vector<std::pair<Key, Value>> serialize() const;

  private:
    std::size_t random_layer() const;

    CompareResult compare(const Key &k1, const Key &k2) const;

    Comparator comparator_;
    size_t elementCount_;
    std::shared_ptr<NodeType> head_;

    static constexpr std::size_t maxLayer_ = 10;
};

template <typename Key, typename Value>
std::vector<std::pair<Key, Value>> SkipList<Key, Value>::serialize() const
{
    std::vector<std::pair<Key, Value>> r;
    r.reserve(elementCount_);

    for (std::shared_ptr<NodeType> cur = head_->forward(1); cur;
         cur = cur->forward(1))
    {
        r.emplace_back(cur->key_, cur->value_);
    }

    return r;
}

template <typename Key, typename Value>
typename SkipList<Key, Value>::CompareResult SkipList<Key, Value>::compare(
    const Key &k1, const Key &k2) const
{
    if (comparator_)
        return comparator_(k1, k2);

    if (k1 > k2)
        return CompareResult::Bigger;
    else if (k1 == k2)
        return CompareResult::Equal;
    else
        return CompareResult::Smaller;
}

template <typename Key, typename Value>
typename SkipList<Key, Value>::EraseType SkipList<Key, Value>::erase(
    const Key &key)
{
    EraseType r = EraseType::NoThisElement;
    std::shared_ptr<NodeType> cur = head_;
    int h = cur->height();

    do
    {
        std::shared_ptr<NodeType> next = cur->forward(h);

        if (!next)
        {
            h--;
            continue;
        }

        switch (compare(key, next->key_))
        {
        case CompareResult::Equal:
            r = EraseType::HasThisElement;
            cur->update(h, next->forward(h));
            continue;
        case CompareResult::Bigger:
            cur = next;
            continue;
        case CompareResult::Smaller:
            h--;
            continue;
        default:
            assert(0);
        }
    } while (h > 0);

    return r;
}

template <typename Key, typename Value>
typename SkipList<Key, Value>::InsertType SkipList<Key, Value>::insert(
    Key &&key, Value &&value)
{
    std::shared_ptr<NodeType> node = find(key);
    if (node)
    {
        node->value_ = std::move(value);
        std::cout << "insert key=" << key << ", value=" << value
                  << ", height=" << node->height() << ", type \""
                  << magic_enum::enum_name(InsertType::Insert) << "\""
                  << std::endl;
        return InsertType::Update;
    }

    std::shared_ptr<NodeType> cur = head_;
    int h = random_layer();
    node = std::make_shared<NodeType>(std::move(key), std::move(value), h);
    do
    {
        std::shared_ptr<NodeType> next = cur->forward(h);
        if (!next)
        {
            cur->update(h, node);
            h--;
            continue;
        }

        switch (compare(key, next->key_))
        {
        case CompareResult::Bigger:
            cur = next;
            continue;
        case CompareResult::Smaller:
            node->update(h, next);
            cur->update(h, node);
            h--;
            continue;
        case CompareResult::Equal:
        default:
            assert(0);
        }
    } while (h > 0);

    std::cout << "insert key=" << key << ", value=" << value
              << ", height=" << node->height() << ", type \""
              << magic_enum::enum_name(InsertType::Insert) << "\"" << std::endl;
    return InsertType::Insert;
}

template <typename Key, typename Value>
std::shared_ptr<Node<Key, Value>> SkipList<Key, Value>::find(
    const Key &key) const
{
    std::shared_ptr<NodeType> cur = head_;
    int h = cur->height();

    do
    {
        std::shared_ptr<NodeType> next = cur->forward(h);
        if (not next)
        {
            h--;
            continue;
        }

        switch (compare(key, next->key_))
        {
        case CompareResult::Equal:
            return next;
        case CompareResult::Bigger:
            cur = next;
            continue;
        case CompareResult::Smaller:
            h--;
            break;
        default:
            assert(0);
        }

    } while (h > 0);

    return nullptr;
}

template <typename Key, typename Value>
std::size_t SkipList<Key, Value>::random_layer() const
{
    size_t r = std::rand() % maxLayer_;
    return r == 0 ? 1 : r;
}