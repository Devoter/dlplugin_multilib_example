#include <iostream>
#include <string>
#include <sstream>

#include "device.hpp"

Device::Device()
{
  m_value = 0;
}

int32_t Device::value() const
{
  return m_value;
}

void Device::set_value(int32_t v)
{
  m_value = v;
}

void Device::print() const
{
  std::cout << m_value << std::endl;
}

std::vector<char> Device::encode_binary() const
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

std::string Device::encode_json() const
{
  std::ostringstream data;
  data << "{\"val\":" << m_value << '}';

  return data.str();
}