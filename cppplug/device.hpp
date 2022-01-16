#ifndef DEVICE_HPP
#define DEVICE_HPP

#include <vector>
#include <cstdint>

class Device
{
  int32_t m_value;

public:
  Device();
  int32_t value();
  void set_value(int32_t v);
  void print();
  std::vector<char> encode();
};

#endif