package usb

import (
	"fmt"
	"github.com/daedaluz/gousb/usbfs"
	"syscall"
)

const (
	usbDevPath = "/dev/bus/usb"
)

type (
	Device struct {
		fd           int
		BusNumber    int
		DeviceNumber int
		Name         string
	}
)

func (d *Device) Open() error {
	if d.fd != -1 {
		return fmt.Errorf("device already open")
	}
	fd, err := usbfs.OpenDevice(d.BusNumber, d.DeviceNumber)
	if err != nil {
		return err
	}
	d.fd = fd
	return nil
}

func (d *Device) IsOpen() bool {
	return d.fd != -1
}

func (d *Device) GetDriver(iface uint32) (string, error) {
	return usbfs.GetDriver(d.fd, iface)
}

func (d *Device) DetachKernel(iface uint32) error {
	return usbfs.Disconnect(d.fd, iface)
}

func (d *Device) AttachKernel(iface uint32) error {
	return usbfs.Connect(d.fd, iface)
}

func (d *Device) Ctrl(typ RequestType, req uint8, value uint16, index uint16, payload []byte) (int, error) {
	return usbfs.ControlTransfer(d.fd, uint8(typ), req, value, index, 1000, payload)
}

func (d *Device) CtrlTimeout(typ RequestType, req uint8, value uint16, index uint16, payload []byte, timeout uint32) (int, error) {
	return usbfs.ControlTransfer(d.fd, uint8(typ), req, value, index, timeout, payload)
}

func (d *Device) Bulk(ep uint8, data []byte) (int, error) {
	return usbfs.BulkTransfer(d.fd, uint32(ep)&0xFF, 1000, data)
}

func (d *Device) BulkTimeout(ep uint8, data []byte, timeout uint32) (int, error) {
	return usbfs.BulkTransfer(d.fd, uint32(ep)&0xFF, timeout, data)
}

func (d *Device) Close() error {
	e := syscall.Close(d.fd)
	d.fd = -1
	return e
}
