#ifndef LIBCPPPLUG_H
#define LIBCPPPLUG_H

#include <cstddef>
#include <cstdint>

typedef void (*get_device_callback_t)(uintptr_t, char *, size_t);

#ifdef __cplusplus
extern "C"
{
#endif

  extern uintptr_t create_device();
  extern int free_device(uintptr_t ptr);
  extern int get_device(uintptr_t ptr, uintptr_t cbId, get_device_callback_t callback);
  extern int device__print(uintptr_t self);
  extern int device__value(uintptr_t self, int32_t *value);
  extern int device__set_value(uintptr_t self, int32_t value);

#ifdef __cplusplus
}
#endif

#endif
