package usbfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	usbDevPath = "/dev/bus/usb"
)

func ioctl(fd int, ioc uint32, arg interface{}) (int, error) {
	b := bytes.Buffer{}
	if err := binary.Write(&b, binary.LittleEndian, arg); err != nil {
		return -1, err
	}
	buff := b.Bytes()
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(ioc), uintptr(unsafe.Pointer(&buff[0])))
	if e != syscall.Errno(0) {
		return int(r), e
	}
	return int(r), nil
}

func GetDriver(fd int, iface uint32) (string, error) {
	data := &usbdevfs_getdriver{
		Interface: iface,
	}
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_GETDRIVER, uintptr(unsafe.Pointer(data)))
	if e == syscall.Errno(0) {
		return data.String(), nil
	}
	return "", e
}

func GetConnectInfo(fd int) (uint8, error) {
	info := &usbdevfs_connectinfo{}
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_CONNECTINFO, uintptr(unsafe.Pointer(info)))
	if e == syscall.Errno(0) {
		return info.Slow, nil
	}
	return 0, e
}

func SetInterface(fd int, iface, setting uint32) error {
	data := &usbdevfs_setinterface{
		Interface:  iface,
		AltSetting: setting,
	}
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_SETINTERFACE, uintptr(unsafe.Pointer(data)))
	if e == syscall.Errno(0) {
		return nil
	}
	return e
}

func ClaimInterface(fd, iface int) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_CLAIMINTERFACE, uintptr(iface))
	if e == syscall.Errno(0) {
		return nil
	}
	return e
}

func ReleaseInterface(fd, iface int) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_RELEASEINTERFACE, uintptr(iface))
	if e == syscall.Errno(0) {
		return nil
	}
	return e
}

func Disconnect(fd int, iface uint32) error {
	data := usbdevfs_ioctl{
		Interface: iface,
		IoctlCode: USBDEVFS_DISCONNECT,
		Data:      0,
	}
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_IOCTL, uintptr(unsafe.Pointer(&data)))
	if e == syscall.Errno(0) {
		return nil
	}
	return e
}

func Connect(fd int, iface uint32) error {
	data := usbdevfs_ioctl{
		Interface: iface,
		IoctlCode: USBDEVFS_CONNECT,
		Data:      0,
	}
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_IOCTL, uintptr(unsafe.Pointer(&data)))
	if e == syscall.Errno(0) {
		return nil
	}
	return e
}

func ControlTransfer(fd int, typ uint8, request uint8, value uint16, index uint16, timeout uint32, payload []byte) (int, error) {
	data := &usbdevfs_ctrltransfer{
		RequestType: typ,
		Request:     request,
		Value:       value,
		Index:       index,
		Timeout:     timeout,
	}
	if payload != nil {
		data.Length = uint16(len(payload))
		data.Data = slicePtr(payload)
	}
	x, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_CONTROL, uintptr(unsafe.Pointer(data)))
	if e == syscall.Errno(0) {
		return int(x), nil
	}
	return int(x), e
}

func BulkTransfer(fd int, endpoint uint32, timeout uint32, payload []byte) (int, error) {
	data := &usbdevfs_bulktransfer{
		Endpoint: endpoint,
		Timeout:  timeout,
	}
	if payload != nil {
		data.Length = uint32(len(payload))
		data.Data = slicePtr(payload)
	}
	x, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_BULK, uintptr(unsafe.Pointer(data)))
	if e == syscall.Errno(0) {
		return int(x), nil
	}
	return int(x), e
}

func ResetDevice(fd int) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), USBDEVFS_RESET, uintptr(0))
	if e == syscall.Errno(0) {
		return nil
	}
	return e
}

func OpenDevice(busNumber, deviceNumber int) (int, error) {
	devPath := fmt.Sprintf("%s/%.3d/%.3d", usbDevPath, busNumber, deviceNumber)
	fd, err := syscall.Open(devPath, syscall.O_RDWR, 0)
	if err != nil {
		return -1, err
	}
	return fd, nil
}
