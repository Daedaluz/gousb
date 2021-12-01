package usb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

type (
	DescriptorType uint8

	Descriptor interface {
		Type() DescriptorType
	}

	DescriptorHeader struct {
		Length         uint8
		DescriptorType DescriptorType
	}

	UnknownDescriptor struct {
		DescriptorHeader
		Data []byte
	}

	DescriptorParser interface {
		ReadUSBDescriptor(hdr DescriptorHeader, i io.Reader) error
	}

	DescriptorFieldParser interface {
		ReadUSBDescriptorField(i io.Reader) (int, error)
	}
)

const (
	DescriptorTypeDevice = DescriptorType(iota + 1)
	DescriptorTypeConfig
	DescriptorTypeString
	DescriptorTypeInterface
	DescriptorTypeEndpoint
)

var (
	descriptorMap = map[DescriptorType]reflect.Type{
		DescriptorTypeDevice:    reflect.TypeOf(DeviceDescriptor{}),
		DescriptorTypeConfig:    reflect.TypeOf(ConfigurationDescriptor{}),
		DescriptorTypeInterface: reflect.TypeOf(InterfaceDescriptor{}),
		DescriptorTypeEndpoint:  reflect.TypeOf(EndpointDescriptor{}),
		DescriptorTypeString:    reflect.TypeOf(StringDescriptor{}),
	}
)

func (h DescriptorHeader) Type() DescriptorType {
	return h.DescriptorType
}

func (t DescriptorType) String() string {
	if typ, exist := descriptorMap[t]; exist {
		return typ.String()
	}
	return fmt.Sprintf("Unknown(0x%.2X)", uint8(t))
}

type (
	DeviceDescriptor struct {
		DescriptorHeader
		BcdUSB             uint16
		BDeviceClass       ClassCode
		BDeviceSubClass    uint8
		BDeviceProtocol    uint8
		BMaxPacketSize0    uint8
		IdVendor           uint16
		IdProduct          uint16
		BcdDevice          uint16
		IManufacturer      uint8
		IProduct           uint8
		ISerialNumber      uint8
		BNumConfigurations uint8
	}

	ConfigurationDescriptor struct {
		DescriptorHeader
		WTotalLength        uint16
		BNumInterfaces      uint8
		BConfigurationValue uint8
		IConfiguration      uint8
		BmAttributes        uint8
		MaxPower            uint8
	}

	InterfaceDescriptor struct {
		DescriptorHeader
		BInterfaceNumber   uint8
		BAlternateSetting  uint8
		BNumEndpoints      uint8
		BInterfaceClass    ClassCode
		BInterfaceSubClass uint8
		BInterfaceProtocol uint8
		IInterface         uint8
	}

	EndpointDescriptor struct {
		DescriptorHeader
		BEndpointAddress uint8
		BmAttributes     uint8
		WMaxPacketSize   uint16
		BInterval        uint8
		BRefresh         uint8
		BSynchAddress    uint8
	}

	StringDescriptor struct {
		DescriptorHeader
		Data []byte
	}

	EndpointCompanionDescriptor struct {
		DescriptorHeader
		bMaxBurst         uint8
		mbAttributes      uint8
		wBytesPerInterval uint16
	}
)

func RegisterDescriptorType(typ DescriptorType, desc Descriptor) {
	descriptorMap[typ] = reflect.TypeOf(desc)
}

func ParseDescriptor(data []byte) (Descriptor, error) {
	reader := bytes.NewReader(data)
	hdr, err := readDescriptorHeader(reader)
	if err != nil {
		return nil, err
	}
	return createDescriptor(hdr, reader)
}

func newDescriptor(hdr DescriptorHeader) (interface{}, reflect.Value) {
	if descriptor, exist := descriptorMap[hdr.DescriptorType]; exist {
		x := reflect.New(descriptor)
		x.Elem().Field(0).Set(reflect.ValueOf(hdr))
		return x.Interface(), x
	}
	x := reflect.New(reflect.TypeOf(UnknownDescriptor{}))
	x.Elem().Field(0).Set(reflect.ValueOf(hdr))
	return x.Interface(), x
}

func createDescriptor(header DescriptorHeader, i io.Reader) (Descriptor, error) {
	descriptor, ptrVal := newDescriptor(header)
	if customReader, implements := descriptor.(DescriptorParser); implements {
		if err := customReader.ReadUSBDescriptor(header, i); err != nil {
			return nil, err
		}
		return descriptor.(Descriptor), nil
	}
	elem := ptrVal.Elem()

loop:
	for elemIndex := 1; elemIndex < elem.NumField(); elemIndex++ {
		field := elem.Field(elemIndex)
		dest := field.Addr().Interface()

		switch field.Kind() {
		case reflect.Slice:
			switch field.Type() {
			case reflect.TypeOf([]uint8{}):
				excessiveData, err := ioutil.ReadAll(i)
				field.Set(reflect.ValueOf(excessiveData))
				if err != nil {
					return nil, err
				}
			default:
				if err := binary.Read(i, binary.LittleEndian, dest); err != nil {
					break loop
				}
			}
		default:
			if err := binary.Read(i, binary.LittleEndian, dest); err != nil {
				break loop
			}
		}
	}
	return descriptor.(Descriptor), nil
}
