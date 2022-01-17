package device

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

type Device struct {
	val int32
}

func NewDevice() *Device {
	return &Device{}
}

func (d *Device) Value() int32 {
	return d.val
}

func (d *Device) SetValue(v int32) {
	d.val = v
}

func (d *Device) Print() {
	fmt.Println(d.val)
}

func (d Device) MarshalBinary() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(d.val))

	binary.LittleEndian.PutUint32(b, uint32(d.val))

	return b, nil
}

func (d *Device) UnmarshalBinary(data []byte) error {
	if len(data) != int(unsafe.Sizeof(d.val)) {
		return errors.New("incompatible data size")
	}

	d.val = int32(binary.LittleEndian.Uint32(data))

	return nil
}

func (d *Device) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"val":%d}`, d.val)), nil
}

func (d *Device) UnmarshalJSON(data []byte) error {
	type tmpt struct {
		Val int32 `json:"val"`
	}

	var tmp tmpt

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	d.val = tmp.Val

	return nil
}
