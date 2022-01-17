#ifndef DEVICE_HPP
#define DEVICE_HPP

#include <vector>
#include <string>
#include <cstdint>

class Device
{
  int32_t m_value;

public:
  Device();
  int32_t value() const;
  void set_value(int32_t v);
  void print() const;
  std::vector<char> encode_binary() const;
  std::string encode_json() const;
};

#endif