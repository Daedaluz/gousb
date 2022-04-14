package usbfs

import (
	"fmt"
	ioctl "github.com/daedaluz/goioctl"
	"syscall"
	"unsafe"
)

// Useful links:
// https://www.engineersgarage.com/usb-requests-and-stages-of-control-transfer-part-4-6/

// ControlTransfer
// fd: device file descriptor
// typ: consists of Direction | Type | Recipient
//      eg RequestDirectionIn | RequestTypeClass | RequestRecipientInterface
// request: specific request.
// value: message value, according to request.
// index: message index value, according to request.
// timeout: timeout in ms.
// payload: data to send.
func ControlTransfer(fd int, typ, request uint8, value, index uint16, timeout uint32, payload []byte) (int, error) {
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

func BulkTransfer(fd int, endpoint, timeout uint32, payload []byte) (int, error) {
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

func SetInterface(fd int, iface, setting uint32) error {
	data := &usbdevfs_setinterface{
		Interface:  iface,
		AltSetting: setting,
	}
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_setinterface, uintptr(unsafe.Pointer(data)))
}

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

func ClaimInterface(fd int, iface uint32) error {
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_claiminterface, uintptr(iface))
}

func ReleaseInterface(fd int, iface uint32) error {
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_releaseinterface, uintptr(iface))
}

func GetConnectInfo(fd int) (uint8, error) {
	info := &usbdevfs_connectinfo{}
	e := ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_connectionfo, uintptr(unsafe.Pointer(info)))
	if e != nil {
		return 0, e
	}
	return info.Slow, nil
}

// DriverIOCTL for talking directly with drivers
func DriverIOCTL(fd int, iface uint32, request, data uintptr) error {
	req := &usbdevfs_ioctl{
		Interface: int32(iface),
		IoctlCode: int32(request),
		Data:      data,
	}
	return ioctl.Ioctl(uintptr(fd), request, uintptr(unsafe.Pointer(req)))
}

func ResetDevice(fd int) error {
	return ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_reset, 0)
}

func GetCapabilities(fd int) (Capability, error) {
	res := Capability(0)
	if err := ioctl.Ioctl(uintptr(fd), ctl_usbdevfs_get_capabilities, uintptr(unsafe.Pointer(&res))); err != nil {
		return 0, err
	}
	return res, nil
}

func SetConfiguration(fd, config int) error {
	return nil
}

func OpenDevice(busNumber, deviceNumber int) (int, error) {
	devPath := fmt.Sprintf("%s/%.3d/%.3d", usbDevPath, busNumber, deviceNumber)
	fd, err := syscall.Open(devPath, syscall.O_RDWR, 0)
	if err != nil {
		return -1, err
	}
	return fd, nil
}

func CloseDevice(fd int) error {
	return syscall.Close(fd)
}

func Disconnect(fd int, iface uint32) error {
	return DriverIOCTL(fd, iface, ctl_usbdevfs_disconnect, 0)
}

func Connect(fd int, iface uint32) error {
	return DriverIOCTL(fd, iface, ctl_usbdevfs_connect, 0)
}
