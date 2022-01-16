#include <iostream>

#include "device.hpp"

Device::Device()
{
  m_value = 0;
}

int32_t Device::value()
{
  return m_value;
}

void Device::set_value(int32_t v)
{
  m_value = v;
}

void Device::print()
{
  std::cout << m_value << std::endl;
}

std::vector<char> Device::encode()
{
  auto data = std::vector<char>(sizeof(int32_t));
  uint32_t v = m_value;

  for (std::size_t i = 0; i < sizeof(int32_t); ++i)
  {
    data[i] = static_cast<unsigned char>(v);
    v >>= 8;
  }

  return data;
}