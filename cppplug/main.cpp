#include <memory>
#include <map>
#include <shared_mutex>
#include <vector>

#include "libcppplug.h"
#include "device.hpp"

std::map<std::uintptr_t, std::unique_ptr<Device>> devices;
std::shared_mutex devmx;

uintptr_t create_device()
{
  auto dev = std::make_unique<Device>();
  auto ptr = reinterpret_cast<std::uintptr_t>(&(*dev));

  devmx.lock();
  devices[ptr] = std::move(dev);
  devmx.unlock();

  return ptr;
}

int free_device(uintptr_t ptr)
{
  devmx.lock();

  if (devices.count(ptr) <= 0)
  {
    devmx.unlock();

    return -1;
  }

  devices.erase(ptr);
  devmx.unlock();

  return 0;
}

int get_device(uintptr_t ptr, uintptr_t cb_id, char use_json, get_device_callback_t callback)
{
  devmx.lock_shared();

  if (devices.count(ptr) <= 0)
  {
    devmx.unlock_shared();

    return -1;
  }

  if (use_json)
  {
    auto encoded = devices[ptr]->encode_json();

    devmx.unlock_shared();
    callback(cb_id, static_cast<char *>(&encoded[0]), encoded.size());
  }
  else
  {
    auto encoded = devices[ptr]->encode_binary();

    devmx.unlock_shared();
    callback(cb_id, static_cast<char *>(&encoded[0]), encoded.size());
  }

  return 0;
}

int device__print(uintptr_t self)
{
  devmx.lock_shared();

  if (devices.count(self) <= 0)
  {
    devmx.unlock_shared();

    return -1;
  }

  devices[self]->print();
  devmx.unlock_shared();

  return 0;
}

int device__value(uintptr_t self, int32_t *value)
{
  devmx.lock_shared();

  if (devices.count(self) <= 0)
  {
    devmx.unlock_shared();

    return -1;
  }

  *value = devices[self]->value();
  devmx.unlock_shared();

  return 0;
}

int device__set_value(uintptr_t self, int32_t value)
{
  devmx.lock();

  if (devices.count(self) <= 0)
  {
    devmx.unlock();

    return -1;
  }

  devices[self]->set_value(value);
  devmx.unlock();

  return 0;
}