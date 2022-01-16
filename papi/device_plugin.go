package papi

/*
#include <stddef.h>
#include <stdint.h>
#include <string.h>

extern void GetDeviceCallback(uintptr_t h, char *data, size_t size);

static uintptr_t create_device(uintptr_t r)
{
	return ((uintptr_t (*)())r)();
}

static int free_device(uintptr_t r, uintptr_t ptr)
{
	return ((int (*)(uintptr_t))r)(ptr);
}

static void handle_callback(uintptr_t cbId, char* data, size_t size)
{
	GetDeviceCallback(cbId, data, size);
}

static int get_device(uintptr_t r, uintptr_t ptr, uintptr_t callback)
{
	typedef void (*get_device_callback_t)(uintptr_t h, char *, size_t);

	return ((int (*)(uintptr_t, uintptr_t, get_device_callback_t))r)(ptr, callback, handle_callback);
}

static int device__print(uintptr_t r, uintptr_t self)
{
	return ((int (*)(uintptr_t))r)(self);
}

static int device__value(uintptr_t r, uintptr_t self, int32_t* value)
{
	return ((int (*)(uintptr_t, int32_t*))r)(self, value);
}

static int device__set_value(uintptr_t r, uintptr_t self, int32_t value)
{
	return ((int (*)(uintptr_t, int32_t))r)(self, value);
}
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"

	"github.com/Devoter/dlplugin_multilib_example/cerror"
)

type getDeviceCallbackFn func(data *C.char, size C.size_t)

type DevicePlugin struct {
	createDevice    func() uintptr
	freeDevice      func(ptr uintptr) error
	getDevice       func(ptr uintptr) (encoded []byte, err error)
	device_Print    func(self uintptr) error
	device_Value    func(self uintptr) (value int32, err error)
	device_SetValue func(self uintptr, value int32) error
}

func (dev *DevicePlugin) CreateDevice() uintptr {
	return dev.createDevice()
}

func (dev *DevicePlugin) FreeDevice(ptr uintptr) error {
	return dev.freeDevice(ptr)
}

func (dev *DevicePlugin) GetDevice(ptr uintptr) (encoded []byte, err error) {
	return dev.getDevice(ptr)
}

func (dev *DevicePlugin) Device_Print(self uintptr) error {
	return dev.device_Print(self)
}

func (dev *DevicePlugin) Device_Value(self uintptr) (value int32, err error) {
	return dev.device_Value(self)
}

func (dev *DevicePlugin) Device_SetValue(self uintptr, value int32) error {
	return dev.device_SetValue(self, value)
}

func (dev *DevicePlugin) Init(lookup func(symName string) (uintptr, error)) error {
	createDevicePtr, err := lookup("create_device")
	if err != nil {
		return err
	}

	freeDevicePtr, err := lookup("free_device")
	if err != nil {
		return err
	}

	dev.createDevice = func() uintptr {
		return uintptr(C.create_device(C.uintptr_t(createDevicePtr)))
	}

	getDevicePtr, err := lookup("get_device")
	if err != nil {
		return err
	}

	devicePrintPtr, err := lookup("device__print")
	if err != nil {
		return err
	}

	deviceValuePtr, err := lookup("device__value")
	if err != nil {
		return err
	}

	deviceSetValuePtr, err := lookup("device__set_value")
	if err != nil {
		return err
	}

	dev.freeDevice = func(ptr uintptr) error {
		return cerror.WrapCError(int(C.free_device(C.uintptr_t(freeDevicePtr), C.uintptr_t(ptr))))
	}

	dev.getDevice = func(ptr uintptr) ([]byte, error) {
		var encoded []byte

		var cb getDeviceCallbackFn = func(data *C.char, size C.size_t) {
			encoded = make([]byte, size)

			C.memcpy(unsafe.Pointer(&encoded[0]), unsafe.Pointer(data), size)
		}

		cbHandle := cgo.NewHandle(cb)
		defer cbHandle.Delete()

		cErr := C.get_device(C.uintptr_t(getDevicePtr), C.uintptr_t(ptr), C.uintptr_t(cbHandle))

		if cErr < 0 { // size is an error code
			return nil, cerror.CError(cErr)
		}

		return encoded, nil
	}

	dev.device_Print = func(self uintptr) error {
		return cerror.WrapCError(int(C.device__print(C.uintptr_t(devicePrintPtr), C.uintptr_t(self))))
	}

	dev.device_Value = func(self uintptr) (int32, error) {
		var value C.int32_t

		err := cerror.WrapCError(int(C.device__value(C.uintptr_t(deviceValuePtr), C.uintptr_t(self), &value)))

		return int32(value), err
	}

	dev.device_SetValue = func(self uintptr, value int32) error {
		return cerror.WrapCError(int(C.device__set_value(C.uintptr_t(deviceSetValuePtr), C.uintptr_t(self), C.int32_t(value))))
	}

	return nil
}

//export GetDeviceCallback
func GetDeviceCallback(h C.uintptr_t, data *C.char, size C.size_t) {
	callback := cgo.Handle(uintptr(h)).Value().(getDeviceCallbackFn)

	callback(data, size)
}
