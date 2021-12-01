package hid

import (
	"github.com/daedaluz/gousb"
	"log"
)

type (
	Device struct {
		*usb.Device
		HidDescriptor *Descriptor
		EpIn          *usb.EndpointDescriptor
		EpOut         *usb.EndpointDescriptor
	}

	Descriptor struct {
		usb.DescriptorHeader
		BcdHID                   uint16
		CountryCode              uint8
		NumDescriptors           uint8
		DescriptorType           uint8
		DescriptorLength         uint16
		OptionalDescriptorType   uint8
		OptionalDescriptorLength uint16
	}
)

const (
	DescriptorTypeHID      = usb.DescriptorType(0x21)
	DescriptorTypeReport   = usb.DescriptorType(0x22)
	DescriptorTypePhysical = usb.DescriptorType(0x23)
)

const (
	GetReport   = 0x01
	GetIdle     = 0x02
	GetProtocol = 0x03
	SetReport   = 0x09
	SetIdle     = 0x0A
	SetProtocol = 0x0B
)

func init() {
	usb.RegisterDescriptorType(DescriptorTypeHID, Descriptor{})
}

func NewHIDDevice(dev *usb.Device) *Device {
	var hidDesc *Descriptor
	var inEp *usb.EndpointDescriptor
	var outEp *usb.EndpointDescriptor

	for _, d := range dev.Descriptors {
		switch desc := d.(type) {
		case *Descriptor:
			hidDesc = desc
		case *usb.EndpointDescriptor:
			if (desc.BEndpointAddress & usb.EndpointDirectionIn) > 0 {
				inEp = desc
			} else {
				outEp = desc
			}
		}
	}
	return &Device{
		Device:        dev,
		HidDescriptor: hidDesc,
		EpIn:          inEp,
		EpOut:         outEp,
	}
}

func (dev *Device) ReadMax() ([]byte, error) {
	size := dev.EpIn.WMaxPacketSize
	buffer := make([]byte, size)
	x, err := dev.Device.BulkTimeout(dev.EpIn.BEndpointAddress, buffer, 100)
	if err != nil {
		return nil, err
	}
	return buffer[0:x], nil
}

func (dev *Device) Read(buff []byte) (int, error) {
	x, err := dev.Device.BulkTimeout(dev.EpIn.BEndpointAddress, buff, 100)
	return x, err
}

func (dev *Device) Write(data []byte) (int, error) {
	x, err := dev.Device.BulkTimeout(dev.EpOut.BEndpointAddress, data, 1000)
	return x, err
}

// Work in progress....
func (dev *Device) GetReportDescriptor() {
	//	buf := make([]byte, dev.HidDescriptor.DescriptorLength)
	x, err := dev.Device.GetDescriptorData(DescriptorTypeReport, 0, dev.HidDescriptor.DescriptorLength)
	log.Printf("%X %v\n ", x, err)
}

func (dev *Device) GetReport() ([]byte, error) {
	data := make([]byte, dev.HidDescriptor.DescriptorLength)
	reqType := usb.RequestDirectionIn | usb.RequestTypeClass | usb.RequestRecipientInterface
	value := uint16(0x01)<<8 | uint16(0)
	idx := uint16(1)
	_, err := dev.Device.Ctrl(reqType, GetReport, value, idx, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (dev *Device) SetReport() {
	panic("Implement me")
}

func (dev *Device) GetIdle(interfaceIdx, reportId uint8) (int, error) {
	data := []byte{0}
	reqType := usb.RequestDirectionIn | usb.RequestTypeClass | usb.RequestRecipientInterface
	_, err := dev.Device.Ctrl(reqType, GetIdle, uint16(reportId), uint16(interfaceIdx), data)
	if err != nil {
		return 0, err
	}
	return int(data[0]), nil
}

func (dev *Device) SetIdle(interfaceIdx, reportId, duration uint8) error {
	reqType := usb.RequestDirectionOut | usb.RequestTypeClass | usb.RequestRecipientInterface
	value := (uint16(duration) << 8) | uint16(reportId)
	idx := uint16(interfaceIdx)
	_, err := dev.Device.Ctrl(reqType, SetIdle, value, idx, nil)
	return err
}

func hidUSBFilter(device *usb.Device) bool {
	for _, desc := range device.Descriptors {
		if _, ok := desc.(*Descriptor); ok {
			return true
		}
	}
	return false
}

func FindHIDDevices() ([]*usb.Device, error) {
	return usb.FindDevices(hidUSBFilter)
}
