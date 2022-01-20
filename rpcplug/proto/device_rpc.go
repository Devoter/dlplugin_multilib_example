package proto

type GetDeviceRequest struct {
	Ptr     uint64
	UseJson bool
}

type DeviceSetValueRequest struct {
	Ptr   uint64
	Value int32
}
