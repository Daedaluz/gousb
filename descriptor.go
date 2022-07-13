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

	DescriptorTypeInterfacePower = DescriptorType(iota + 8)
	DescriptorTypeOTG
	DescriptorTypeDebug
	DescriptorTypeInterfaceAssociation

	DescriptorTypeBOS                                        = DescriptorType(15)
	DescriptorTypeDeviceCapability                           = DescriptorType(16)
	DescriptorTypeSuperSpeedUSBEndprointCompanion            = DescriptorType(48)
	DescriptorTypeSuperSpeedPlusIsochronousEndpointCompanion = DescriptorType(49)
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
	// DeviceDescriptor describes general information about a device.
	// It includes information that applies globally to the device and all of the devices configurations.
	// A device has only one DeviceDescriptor.
	//
	// All devices have a default control pipe. The maximum packet size of a device’s default
	// control pipe is described in the device descriptor. Endpoints specific to a configuration and
	// its interface(s) are described in the configuration descriptor. A conf iguration and its
	// interface(s) do not include an endpoint descriptor for the default control pipe. Other than
	// the maximum packet size, the characteristics of the default control pipe are defined by this
	// specification and are the same for all Enhanced SuperSpeed devices.
	DeviceDescriptor struct {
		DescriptorHeader
		// The bcdUSB field contains a BCD version number. The value of the bcdUSB field is 0xJJMN
		// for version JJ.M.N (JJ – major version number, M – minor version number, N – sub-minor
		// version number), e.g., version 2.1.3 is represented with value 0213H and version 3.0 is
		// represented with a value of 0300H.
		//
		// The device descriptor of an Enhanced SuperSpeed device shall have a version number of 3.1
		// (0310H). The device descriptor of an Enhanced SuperSpeed device operating in one of the
		// USB 2.0 modes shall have a version number of 2.1 (0210H).
		BcdUSB uint16

		// BDeviceClass is a class code assigned by the USB-IF.
		// If this field is reset to zero, each interface within a
		// configuration specifies its own class information
		// and the various interfaces operate independently.
		// If this field is set to FFH, the device class is vendor-specific.
		BDeviceClass ClassCode

		// BDeviceSubClass is a subclass code assigned by the USB-IF.
		// These codes are qualified by the value of the bDeviceClass field.
		// If the bDeviceClass field is reset to zero, this field shall also be reset to zero.
		// If the bDeviceClass field is not set to FFH, all values are reserved for assignment by the USB-IF.
		BDeviceSubClass SubClass

		// BDeviceProtocol (assigned by the USB-IF).
		// These codes are qualified by the value of the bDeviceClass and the bDeviceSubClass fields.
		// If a device supports class-specific protocols on a device basis as opposed to an interface
		// basis, this code identifies the protocols that the device uses as defined by the
		// specification of the device class.
		// If this field is reset to zero, the device does not use class-specific protocols on a device basis.
		// However, it may use class-specific protocols on an interface basis.
		// If this field is set to FFH, the device uses a vendor-specific protocol on a device basis.
		BDeviceProtocol uint8

		// BMaxPacketSize0 Maximum packet size for endpoint zero.
		// The BMaxPacketSize0 value is used as the exponent of 2, this means a value of 4 means
		// a max packet size of 16 (2^4 -> 16).
		//
		// 09H is the only valid value in this field when operating at Gen X speed.
		//
		// An Enhanced SuperSpeed device shall set the bMaxPacketSize0 field to 09H (see Table 9-11)
		// indicating a 512-byte maximum packet. An Enhanced SuperSpeed device shall not support
		// any other maximum packet sizes for the default control pipe (endpoint 0) control endpoint.
		BMaxPacketSize0 uint8

		// Vendor ID assigned by the USB-IF.
		IDVendor uint16

		// Product ID assigned by the manufacturer.
		IDProduct uint16

		// BcdDevice release number in binary-coded decimal.
		BcdDevice uint16

		// IManufacturer Index of string descriptor describing manufacturer
		IManufacturer uint8

		// IProduct Index of string descriptor describing product.
		IProduct uint8

		// ISerialNumber Index of string descriptor describing the devices serial number
		ISerialNumber uint8

		// The bNumConfigurations field indicates the number of configurations at the current
		// operating speed. Configurations for the other operating speed are not included in the count.
		// If there are specific configurations of the device for specific speeds, the bNumConfigurations
		// field only reflects the number of configurations for a single speed, not the total number of
		// configurations for both speeds.
		BNumConfigurations uint8
	}

	// BOSDescriptor defines a root descriptor that is similar to the configuration descriptor,
	// and is the base descriptor for accessing a family of related descriptors.
	// A host can read a BOS descriptor and learn from the wTotalLength field the entire size
	// of the device-level descriptor set, or it can read in the entire BOS descriptor
	// set of device capabilities.
	// The host accesses this descriptor using the GetDescriptor() request.
	// There is no way for a host to read individual device capability descriptors.
	// The entire set can only be accessed via reading the BOS descriptor with a
	// GetDescriptor() request and using the length reported in the WTotalLength field.
	BOSDescriptor struct {
		DescriptorHeader
		// Length of this descriptor and all of its sub descriptors
		WTotalLength uint16

		// The number of separate device capability descriptors in the BOS
		BNumDeviceCaps uint8
		/* DeviceCapabilityDescriptors */
	}

	// DeviceCapabilityDescriptor are always returned as part of the BOS information returned
	// by a GetDescriptor(BOS) request. A Device Capability cannot be directly accessed with a
	// GetDescriptor() or SetDescriptor() request.
	DeviceCapabilityDescriptor struct {
		DescriptorHeader
		// BDevCapabilityType capability code listed in Table 9-14.
		BDevCapabilityType Capability

		// Data is capability specific data.
		Data []byte
	}

	// ConfigurationDescriptor describes information about a specific device configuration.
	// The descriptor contains a BConfigurationValue field with a value that, when used as a parameter to
	// SetConfiguration() request, causes the device to assume that described configuration.
	//
	// When the host requests the configuration descriptor, all related interface, endpoint, and
	// endpoint companion descriptors are returned.
	ConfigurationDescriptor struct {
		DescriptorHeader
		// WTotalLength Total length of data returned for this configuration
		// Includes the combined length of all descriptors (configuration, interface,endpoint,
		// and class- or vendor-specific) returned for this configuration
		WTotalLength uint16

		// BNumInterfaces represents the number of interfaces supported by this configuration.
		BNumInterfaces uint8

		// BConfigurationValue Value to use as an argument to the SetConfiguration() request to select this configuration.
		BConfigurationValue uint8

		// IConfiguration Index of string descriptor describing this configuration.
		IConfiguration uint8

		// BmAttributes Configuration characteristics.
		//
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 7     | Reserved. Shall be set to 0 for historical reasons.    |
		// +----------------------------------------------------------------+
		// | 6     | Self-powered. A device configuration that uses power   |
		// |       | from the bus and a local source reports a non-zero     |
		// |       | value in BMaxPower to indicate the ammount of bus      |
		// |       | power required and sets this bit field.                |
		// +----------------------------------------------------------------+
		// | 5     | Remote wakeup. If a device configuration supports      |
		// |       | remote wakeup, this bit is set.                        |
		// +----------------------------------------------------------------+
		// | 4:0   | Reserved                                               |
		// +----------------------------------------------------------------+
		BmAttributes uint8

		// BMaxPower is the maximum power consumption of the device from the bus in this specific configuration
		// when the device is fully operational.
		//
		// Expressed in 2 mA units when the device is operating in high-speed mode and in 8 mA units when
		// operating at Gen X speed.
		// i.e:
		// 50 = 100 mA when operating at high-speed
		// 50 = 400 mA when operating at Gen X speed
		//
		// Note: A device configuration reports whether the configuration is bus-powered or self-powered.
		//       Device status reports whether the device is currently self-powered.
		//       If a device is disconnected from its external power source, it updates device status to indicate
		//       that it is no longer self-powered.
		//       A device may not increase its power draw from the bus, when it loses its external power
		//       source, beyond the amount reported by its configuration.
		//       If a device can continue to operate when disconnected from its external power source,
		//       it continues to do so.
		//       If the device cannot continue to operate, it shall return to the Powered state.
		BMaxPower uint8
	}

	// InterfaceAssociationDescriptor is used to describe that two or more interfaces are associated to the same function.
	// An “association” includes two or more interfaces and all of their alternate setting interfaces.
	// A device must use an Interface Association descriptor for each device function
	// that requires more than one interface.
	// An Interface Association descriptor is always returned as part of the configuration information returned by a
	// GetDescriptor(Configuration) request.
	// An InterfaceAssociationDescriptor cannot be directly accessed with a GetDescriptor() or SetDescriptor() request.
	// An interfaceAssociationDescriptor must be located before the set of interface descriptors
	// (including all alternate settings) for the interfaces it associates.
	// All of the interface numbers in the set of associated interfaces must be contiguous.
	//
	// The interface association descriptor includes func tion class,subclass, and protocol fields.
	// The values in these fields can be the same as the interface class, subclass, and protocol values from
	// any one of the associated interfaces.
	// The preferred implementation, for existing device classes, is to use th e interface class, subclass, and
	// protocol field values from the first interface in the list of associated interfaces.
	//
	//  Note: Since this particular feature was not included in earlier versions of the USB specification, there
	//        is an issue with how existing USB operating system implementations will support devices that
	//        use this descriptor. It is strongly recommended that dev ice implementations utilizing the
	//        interface association descriptor use the Multi-interface Function Class codes in the device
	//        descriptor. This allows simple and easy identification of these devices and allows on some
	//        operating systems, installation of an upgrade driver that can parse and enumerate
	//        configurations that include the Interface Association Descriptor. The Multi-interface Function.
	//        Class code is documented at http://www.usb.org/developers/docs.
	//        TODO: enumerate described class codes.
	InterfaceAssociationDescriptor struct {
		DescriptorHeader
		// BFirstInterface is the number of the first interface that is associated with this function.
		BFirstInterface uint8

		// BInterfaceCount is the number of contiguous interfaces that are associated with this function.
		BInterfaceCount uint8

		// BFunctionClass Class code (assigned by USB-IF).
		// A value of zero is not allowed in this descriptor.
		// If this field is FFH, the function class is vendor-specific.
		// All other values are reserved for assignment by the USB-IF.
		BFunctionClass ClassCode

		// BFunctionSubClass Subclass code (assigned by USB-IF).
		// If the bFunctionClass field is not set to FFH,
		// all values are reserved for assignment by the USB-IF.
		BFunctionSubClass SubClass

		// BFunctionProtocol Protocol code (assigned by USB-IF).
		// These codes are qualified by the values of the BFunctionClass and BFunctionSubClass fields.
		BFunctionProtocol uint8

		// IFunction is an index of a string descriptor describing this function.
		IFunction uint8
	}

	// InterfaceDescriptor describes a specific interface within a configuration.
	// An interface descriptor is always returned as part of a configuration descriptor.
	// Interface descriptors cannot be directly accessed with a GetDescriptor() or SetDescriptor() request.
	// An interface may include alternate settings that allow the endpoints and/or their
	// characteristics to be varied after the device has been configured.
	// The default setting for an interface is always alternate setting zero.
	// The SetInterface() request is used to select an alternate setting or to return to the default setting.
	// The GetInterface() request returns the selected alternate setting.
	//
	// Alternate settings allow a portion of the device configuration to be varied while other
	// interfaces remain in operation.
	// If a configuration has alternate settings for one or more of its interfaces,
	// a separate interface descriptor and its associated endpoint and endpoint
	// companion (when reporting its Enhanced SuperSpeed configuration) descriptors are included for each setting.
	InterfaceDescriptor struct {
		DescriptorHeader
		// BInterfaceNumber Number of this interface.
		// Zero-based value identifying the index in the array of concurrent
		// interfaces supported by this configuration.
		BInterfaceNumber uint8

		// BAlternateSetting Value used to select this alternate setting
		// for the interface identified in the prior field.
		BAlternateSetting uint8

		// BNumEndpoints Number of endpoints used by this interface (excluding the Default Control Pipe).
		// If this  value is zero, this interface only uses the Default Control Pipe.
		BNumEndpoints uint8

		// BInterfaceClass Class code (assigned by the USB-IF).
		// A value of zero is reserved for future standardization.
		// If this field is set to FFH, the interface class is vendor-specific.
		// All other values are reserved for assignment by the USB-IF.
		BInterfaceClass ClassCode

		// BInterfaceSubClass Subclass code (assigned by the USB-IF).
		// These codes are qualified by the value of the BInterfaceClass field.
		// If the bInterfaceClass field is reset to zero, this field shall also be reset to zero.
		// If the bInterfaceClass field is not set to FFH, all values are reserved for assignment by
		// the USB-IF.
		BInterfaceSubClass SubClass

		// BInterfaceProtocol Protocol code (assigned by the USB).
		// These codes are qualified by the value of the BInterfaceClass and the BInterfaceSubClass fields.
		// If an interface supports class-specific requests, this code identifies the protocols that the
		// device uses as defined by the specification of the device.
		//
		// If this field is set to zero, the device does not use a class-specific protocol on this interface.
		// If this field is set to FFH, the device uses a vendor-specifc protocol for this interface.
		BInterfaceProtocol uint8

		// IInterface Index of string descriptor describing this interface.
		IInterface uint8
	}

	// EndpointDescriptor contains the information required by the host to determine
	// the bandwidth requirements of each endpoint.
	// An endpoint descriptor cannot be directly accessed with a GetDescriptor() or SetDescriptor() request.
	// There is never an endpoint descriptor for endpoint zero.
	EndpointDescriptor struct {
		DescriptorHeader
		// BEndpointAddress The address of the endpoint on the device described by this descriptor.
		//
		// The address is encoded as follows:
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 3:0   | The endpoint number                                    |
		// +----------------------------------------------------------------+
		// | 6:4   | Reserved, reset to zero                                |
		// +----------------------------------------------------------------+
		// | 7     | Direction. Ignored for control endpoints.              |
		// |       |    0 - OUT endpoint                                    |
		// |       |    1 - IN  endpoint                                    |
		// +----------------------------------------------------------------+
		BEndpointAddress uint8

		// BmAttributes This field describes the endpoint’s attributes when it is
		// configured using the BConfigurationValue.
		//
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 1:0   | Transfer type.                                         |
		// |       |    00 - Control                                        |
		// |       |    01 - Isochronous                                    |
		// |       |    10 - Bulk                                           |
		// |       |    11 = Interrupt                                      |
		// +----------------------------------------------------------------+
		// | 3:2   | If an interrupt endpoint:                              |
		// |       |    Reserved. Shall be set to 0                         |
		// |       | If isochronous:                                        |
		// |       |   Synchronization type:                                |
		// |       |     00 - No synchronization                            |
		// |       |     01 - Asynchronous                                  |
		// |       |     10 - Adaptive                                      |
		// |       |     11 - Synchronous                                   |
		// +----------------------------------------------------------------+
		// | 5:4   | If an interrupt endpoint:                              |
		// |       |   Usage type:                                          |
		// |       |     00 - Periodic                                      |
		// |       |     01 - Notification                                  |
		// |       |     10 - Reserved                                      |
		// |       |     11 - Reserved                                      |
		// |       | If isochronous:                                        |
		// |       |   Usage type:                                          |
		// |       |     00 - Data endpoint                                 |
		// |       |     01 - Feedback endpoint                             |
		// |       |     10 - Implicit feedback data endpoint               |
		// |       |     11 - Reserved                                      |
		// +----------------------------------------------------------------+
		BmAttributes uint8

		// WMaxPacketSize Maximum packet size this endpoint is capable of sending or
		// receiving when this configuration is selected.
		// For control endpoints this field shall be set to 512.
		// For bulk endpoint types this field shall be set to 1024.
		//
		// For interrupt and isochronous endpoints this field shall be set to 1024 if this
		// endpoint defines a value in the BMaxBurst field greater than zero.
		// If the value in the bMaxBurst field is set to zero then this field can have any value from 0 to 1024.
		// for an isochronous endpoint and 1 to 1024 for an interrupt endpoint.
		WMaxPacketSize uint16

		// BInterval for servicing the endpoint for data transfers.
		// Expressed in 125 µs units.
		// For Enhanced SuperSpeed isochronous and interrupt endpoints, this value shall be in the range from 1 to 16.
		// However, the valid ranges are 8 to 16 for Notification type Interrupt endpoints.
		// The bInterval value is used as the exponent for a 2^(BInterval-1) value;
		// eg a BInterval of 4 means a period of 8 ( 2^(4-1) -> 2^3 -> 8 ).
		//
		// This field is reserved and shall not be used for Enhanced SuperSpeed bulk or control endpoints.
		BInterval uint8
	}

	// StringDescriptor are optional.
	// String descriptors use UNICODE UTF16LE encodings as defined by The Unicode Standard,
	// Worldwide Character Encoding, Version 5.0, The Unicode Consortium, Addison-Wesley
	// Publishing Company, Reading, Massachusetts (http://www.unicode.org).
	// The strings in a device may support multiple languages.
	// When requesting a string descr iptor, the requester specifies the desired language using a 16-bit
	// language ID (LANGID) defined by the USB-IF. TODO: Lookup and enumerate langIDs.
	// String index zero for all languages returns a string descriptor that contains an array of
	// 2-byte LANGID codes supported by the device.
	// A device may omit all string descriptors.
	// Devices that omit all string descriptors shall not return an array of LANGID codes.
	StringDescriptor struct {
		DescriptorHeader
		// If langID is zero, this field contains an array of []uint16 of supported languages.
		// else, this field contains the string of specified language.
		Data []byte
	}

	// SSEndpointCompanionDescriptor shall only be returned by Enhanced SuperSpeed devices that
	// are operating at Gen X Speed.
	// Each endpoint described in an interface is followed by a SuperSpeed Endpoint
	// Companion descriptor.
	// The Default Control Pipe does not have an Endpoint Companion descriptor.
	// The Endpoint Companion descriptor shall immediately follow the endpoint descriptor it is
	// associated with in the configuration information.
	SSEndpointCompanionDescriptor struct {
		DescriptorHeader
		// BMaxBurst The maximum number of packets the endpoint can send or receive as part of a burst.
		// Valid values are from 0 to 15.
		// A value of 0 indicates that the endpoint can only burst one packet at a time and a
		// value of 15 indicates that the endpoint can burst up to 16 packets at a time.
		// For endpoints of type control this shall be set to 0.
		BMaxBurst uint8

		// BmAttributes.
		// If this is a bulk endpoint:
		//
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 4:0   | Max streams. The maximum number of streams this        |
		// |       | endpoint supports. Valid values are from 0 to 16,      |
		// |       | where a value of 0 indicates that the endpoint does    |
		// |       | not define streams. for the values 1 to 16, the number |
		// |       | of streams supported equals 2^MaxStream.               |
		// +----------------------------------------------------------------+
		// | 7:5   | Reserved. these bits are reserved                      |
		// |       | and shall be set to zero.                              |
		// +----------------------------------------------------------------+
		//
		// If this is a control or interrupt endpoint:
		//
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 7:0   | Reserved.                                              |
		// +----------------------------------------------------------------+
		//
		// If this is an isochronous endpoint:
		//
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 1:0   | Mult. A zero based value that determines the maximum   |
		// |       | number of packets within a service interval that       |
		// |       | this endpoint supports                                 |
		// |       | Maximum number of packets = (BMaxBurst+1) * (Mult+1)   |
		// |       | The maximum value that can be set in this field is 2.  |
		// |       | This field shall be set to zero if BMaxBurst is zero   |
		// +----------------------------------------------------------------+
		// | 6:2   | Reserved.                                              |
		// | 7     | SSP ISO Companion.                                     |
		// |       | If this field is set to one then a                     |
		// |       | SuperSpeedPlus Isochronous Endpoint Companion          |
		// |       | descriptor shall immediately follow this descriptor    |
		// |       | and the value in the Mult field shall be ignored.      |
		// |       |                                                        |
		// |       | The actual Mult shall be determined as follows:        |
		// |       | DWBytesPerInterval / BMaxBurst / WMaxPacketSize        |
		// |       | rounded up to the nearest integer value.               |
		// +----------------------------------------------------------------+
		BmAttributes uint8

		// The total number of bytes this endpoint will transfer every service interval (SI).
		// This field is only valid for periodic endpoints.
		//
		// For isochronous endpoints:
		//   If the SSP ISO Companion bit in the bmAttributes
		//   field is set to zero, this value is used to reserve the
		//   bus time in the schedule, required for the frame data payloads per SI.
		//   The pipe may, on an ongoing basis, actually use less bandwidth than that reserved.
		//   The device reports, if necessary, the actual bandwidth used via its normal, non-USB defined mechanisms.
		//
		//   If the SSP ISO Companion bit in the BmAttributes field is set to one, this field shall be set to one, and
		//   the total number of bytes this endpoint will transfer shall be reported via the endpoint’s SuperSpeedPlus
		//   Isochronous Endpoint Companion descriptor.
		//
		// wBytesPerInterval is reserved and must be set to zero for control and bulk endpoints.
		WBytesPerInterval uint16
	}

	// SSPIsochronousEndpointCompanionDescriptor contains additional endpoint characteristics that are only defined for
	// endpoints of devices operating at above Gen 1 speed.
	// The SuperSpeedPlus Isochronous Endpoint Companion descriptor shall immediately follow
	// the SuperSpeed Endpoint Companion descriptor that follows the Isochronous endpoint
	// descriptor in the configuration information.
	//
	// When an alternate setting is selected that has an Isochronous endpoint that has a
	// SuperSpeedPlus Isochronous Endpoint Companion descriptor the endpoint shall operate
	// with the characteristics as described in the SuperSpeedPlus Isochronous Endpoint Companion descriptor.
	SSPIsochronousEndpointCompanionDescriptor struct {
		DescriptorHeader

		// Reserved.
		WReserved uint16

		// DWBytesPerInterval The total number of bytes this endpoint will transfer every
		// service interval (SI).
		//
		// This value is used to reserve the bus time in the schedule, required for the frame data payloads per SI.
		// The pipe may, on an ongoing basis, actually use less bandwidth than that reserved.
		// The device reports, if necessary, the actual bandwidth used via its normal, non-USB defined
		// mechanisms.
		//
		// The value in this field shall be less than
		// (MAX_ISO_BYTES_PER_BI_GEN1 x Number of Lanes x Lane Speed Mantissa/LANE_SPEED_MANTISSA_GEN1)
		//
		// MAX_ISO_BYTES_PER_BI_GEN1 = 48 * 1024 (49152):
		//    Maximum number of Isochronous bytes per bus interval when a device is operating at Gen 1x1 speed.
		// LANE_SPEED_MANTISSA_GEN1 = 5:
		//    Land Speed Mantissa for Gen 1.
		DWBytesPerInterval uint32
	}
)

func RegisterDescriptorType(typ DescriptorType, desc Descriptor) {
	descriptorMap[typ] = reflect.TypeOf(desc)
}

func readDescriptorHeader(i io.Reader) (*DescriptorHeader, error) {
	header := DescriptorHeader{
		Length:         0,
		DescriptorType: 0,
	}
	err := binary.Read(i, binary.BigEndian, &header)
	return &header, err
}

func newDescriptor(hdr DescriptorHeader) (any, reflect.Value) {
	if descriptor, exist := descriptorMap[hdr.DescriptorType]; exist {
		x := reflect.New(descriptor)
		x.Elem().Field(0).Set(reflect.ValueOf(hdr))
		return x.Interface(), x
	}
	x := reflect.New(reflect.TypeOf(UnknownDescriptor{}))
	x.Elem().Field(0).Set(reflect.ValueOf(hdr))
	return x.Interface(), x
}

func readDescriptor(header *DescriptorHeader, i io.Reader) (Descriptor, error) {
	descriptor, ptrVal := newDescriptor(*header)
	if customReader, implements := descriptor.(DescriptorParser); implements {
		if err := customReader.ReadUSBDescriptor(*header, i); err != nil {
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

func ReadDescriptors(i io.Reader, descriptorCB func(d Descriptor)) error {
	var err error
	var hdr *DescriptorHeader
	for hdr, err = readDescriptorHeader(i); err == nil; hdr, err = readDescriptorHeader(i) {
		descriptor, err := readDescriptor(hdr, i)
		if err != nil {
			return err
		}
		descriptorCB(descriptor)
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func ParseDescriptor(data []byte) (Descriptor, error) {
	reader := bytes.NewReader(data)
	hdr, err := readDescriptorHeader(reader)
	if err != nil {
		return nil, err
	}
	return readDescriptor(hdr, reader)
}
