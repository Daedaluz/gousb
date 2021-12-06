package usbfs

import (
	"fmt"
	ioctl "github.com/daedaluz/goioctl"
	"syscall"
	"unsafe"
)

const (
	usbDevPath = "/dev/bus/usb"
)

func GetDriver(fd int, iface uint32) (string, error) {
	data := &usbdevfs_getdriver{
		Interface: iface,
	}
	e := ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_getdriver, uintptr(unsafe.Pointer(data)))
	if e != nil {
		return "", e
	}
	return data.String(), nil
}

func GetConnectInfo(fd int) (uint8, error) {
	info := &usbdevfs_connectinfo{}
	e := ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_connectionfo, uintptr(unsafe.Pointer(info)))
	if e != nil {
		return 0, e
	}
	return info.Slow, nil
}

func SetInterface(fd int, iface, setting uint32) error {
	data := &usbdevfs_setinterface{
		Interface:  iface,
		AltSetting: setting,
	}
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_setinterface, uintptr(unsafe.Pointer(data)))
}

func ClaimInterface(fd, iface int) error {
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_claiminterface, uintptr(iface))
}

func ReleaseInterface(fd, iface int) error {
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_releaseinterface, uintptr(iface))
}

func Disconnect(fd int, iface uint32) error {
	data := &usbdevfs_ioctl{
		Interface: iface,
		IoctlCode: USBDEVFS_DISCONNECT,
		Data:      0,
	}
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_ioctl, uintptr(unsafe.Pointer(data)))
}

func Connect(fd int, iface uint32) error {
	data := &usbdevfs_ioctl{
		Interface: iface,
		IoctlCode: USBDEVFS_CONNECT,
		Data:      0,
	}
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_ioctl, uintptr(unsafe.Pointer(data)))
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
	x, e := ioctl.IoctlX(uintptr(fd), ctl_usbdevfs_control, uintptr(unsafe.Pointer(data)))
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
	x, e := ioctl.IoctlX(uintptr(fd), ctl_usbdevfs_bulk, uintptr(unsafe.Pointer(data)))
	return int(x), e
}

func ResetDevice(fd int) error {
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_reset, 0)
}

func OpenDevice(busNumber, deviceNumber int) (int, error) {
	devPath := fmt.Sprintf("%s/%.3d/%.3d", usbDevPath, busNumber, deviceNumber)
	fd, err := syscall.Open(devPath, syscall.O_RDWR, 0)
	if err != nil {
		return -1, err
	}
	return fd, nil
}
