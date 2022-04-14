package usb

import (
	"fmt"
	"reflect"
)

type Capability uint8

func (c Capability) String() string {
	if str, ok := capabilityStringMap[c]; ok {
		return str
	}
	return fmt.Sprintf("CapUnknown(0x%.2X)", uint8(c))
}

// Capability constants
const (
	CapWirelessUSB          = Capability(0x01)
	CapUSB20Extension       = Capability(0x02)
	CapSuperSpeedUSB        = Capability(0x03)
	CapContainerID          = Capability(0x04)
	CapPlatform             = Capability(0x05)
	CapPowerDelivery        = Capability(0x06)
	CapBatteryInfo          = Capability(0x07)
	CapPDConsumerPort       = Capability(0x08)
	CapPDProviderPort       = Capability(0x09)
	CapSuperSpeedPlus       = Capability(0x0A)
	CapPrecisionTime        = Capability(0x0B)
	CapWirelessUSBExt       = Capability(0x0C)
	CapBillboard            = Capability(0x0D)
	CapAuthentication       = Capability(0x0E)
	CapBillboardEx          = Capability(0x0F)
	CapConfigurationSummary = Capability(0x10)
)

type (
	// CapUSB20ExtensionDescriptor is an Enhanced SuperSpeed device shall include the
	// USB 2.0 Extension descriptor and shall support LPM when operating in USB 2.0 High-Speed mode.
	CapUSB20ExtensionDescriptor struct {
		DescriptorHeader
		// BDevCapabilityType Capability type: CapUSB20Extension
		BDevCapabilityType Capability
		// BMAttributes bitfield.
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 0     | reserved. Shall be set to 0.                           |
		// |----------------------------------------------------------------+
		// | 1     | LPM. A value of one in this bit                        |
		// |       | location indicates that this device                    |
		// |       | supports the Link Power Management protocol.           |
		// |       | Enhanced SuperSpeed devices shall set this bit to one. |
		// +----------------------------------------------------------------+
		// | 31:2  | Reserved. Shall be set to 0.                           |
		// +----------------------------------------------------------------+
		BMAttributes uint32
	}

	// CapSuperSpeedUSBDescriptor describes a device-level descriptor which shall be
	// implemented by all Enhanced SuperSpeed devices.
	CapSuperSpeedUSBDescriptor struct {
		DescriptorHeader
		// BDevCapabilityType Capability type: CapSuperSpeedUSB
		BDevCapabilityType Capability

		// BMAttributes bitfield.
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 0     | reserved. Shall be set to 0.                           |
		// |----------------------------------------------------------------+
		// | 1     | LPM. A value of one in this bit                        |
		// |       | location indicates that this device                    |
		// |       | supports the Link Power Management protocol.           |
		// |       | Enhanced SuperSpeed devices shall set this bit to one. |
		// +----------------------------------------------------------------+
		// | 7:2   | Reserved. Shall be set to 0.                           |
		// +----------------------------------------------------------------+
		BMAttributes uint8

		// WSpeedsSupported bitfield.
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 0     | low-Speed USB support                                  |
		// +----------------------------------------------------------------+
		// | 1     | full-Speed USB support                                 |
		// +----------------------------------------------------------------+
		// | 2     | high-Speed USB support                                 |
		// +----------------------------------------------------------------+
		// | 3     | Gen 1 speed support                                    |
		// +----------------------------------------------------------------+
		// | 15:4  | Reserved. Shall be set to zero                         |
		// +----------------------------------------------------------------+
		WSpeedsSupported uint16

		// BFunctionalitySupport is the minimal speed setting at which
		// all the functionality supported by the device is a vailable to the user.
		// For example, if the device supports all its functionality when connected
		// at full-speed and above, the device sets this value to 1.
		// Refer to WSpeedsSupported field for valid values.
		BFunctionalitySupport uint8

		// BU1DevExitLat is U1 Device Exit Latency.
		// Worst-case latency to transition from U1 to U0, assuming the latency is
		// limited only by the device and not the device’s link partner.
		// This field applies only to the exit latency associated
		// with an individual port, and does not apply to the
		// total latency through a hub (e.g., from downstream port to upstream port).
		//
		// Valid values are 0x00 through 0x0A (less than 10) uS.
		//
		// For a hub, this is the value for both its upstream and
		// downstream ports.
		BU1DevExitLat uint8

		// WU2DevExitLat is U2 Device Exit Latency.
		// Worst-case latency to transition from U2 to U0, assuming the latency is
		// limited only by the device and not the device’s link partner.
		// Applies to all ports on a device.
		//
		// Valid values are 0x0000 through 0x07FF (less than 2047) uS.
		//
		// For a hub, this is the value for both its upstream and
		// downstream ports.
		WU2DevExitLat uint16
	}

	// CapContainerIDDescriptor shall be implemented by all USB hubs, and is optional for other devices.
	// If this descriptor is provided when operating in one mode,
	// it shall be provided when operating in any mode.
	// T his descriptor may be used by a host in order to identify a unique
	// device instance across all operating modes.
	// If a device can also connect to a host through other technologies, the same Container ID value
	// contained in this descriptor should also be provided over those other technologies in a
	// technology specific manner.
	CapContainerIDDescriptor struct {
		DescriptorHeader
		// BDevCapabilityType Capability type: CapContainerID
		BDevCapabilityType Capability

		// Reserved.
		Reserved uint8

		// ContainerID UUID. This is a 128-bit number that is unique to a device
		// instance that is used to uniquely identify the device instance across all modes of operation.
		// This same value may be provided over other technologies as well to allow the host to identify
		// the device independent of means of connectivity.
		ContainerID [16]byte
	}

	// CapPlatformDescriptor contains a 128-bit UUID value that is defined and published
	// independently by the platform/operating system vendor, and is used to identify a unique
	// platform specific device capability.
	// The descriptor may also contain one or more bytes of data associated with the capability.
	CapPlatformDescriptor struct {
		DescriptorHeader
		// BDevCapabilityType Capability type: CapPlatform
		BDevCapabilityType Capability

		// Reserved.
		Reserved uint8

		// PlatformCapabilityUUID is a 128-bit number that uniquely identifies a platform specific
		// capability of the device
		PlatformCapabilityUUID [16]byte

		// CapabilityData is a variable length field containing data associated with the platform
		// specific capability. This field may be 0 bytes in length.
		CapabilityData []byte
	}

	CapSuperSpeedPlusUSBDescriptor struct {
		DescriptorHeader
		// BDevCapabilityType Capability type: CapSuperSpeedPlus
		BDevCapabilityType Capability

		// BReserved1.
		BReserved1 uint8

		// BMAttributes Bitmap
		// Bitmap encoding of supported SuperSpeedPlus features:
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 4:0   | Sublink Speed Attribute Count (SSAC).                  |
		// |       | The number of Sublink Speed Attribute bitmaps.         |
		// |       | A SuperSpeedPlus device shall report at least one SSAC.|
		// |       | The number of Sublink Speed Attribute Count = SSAC + 1.|
		// +----------------------------------------------------------------+
		// | 8:5   | Sublink Speed ID Count (SSIC).                         |
		// |       | The number of unique Sublink Speed IDs supported by    |
		// |       | the device.                                            |
		// |       | The number of Sublink Speed IDs = SSIC + 1.            |
		// +----------------------------------------------------------------+
		// | 31:9  | Reserved                                               |
		// +----------------------------------------------------------------+
		BMAttributes uint32

		// WFunctionalitySupport.
		// The device shall support full functionality at all reported bandwidths
		// at or above the minimum bandwidth described via this field.
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 3:0   | Sublink Speed Attribute ID (SSID).                     |
		// |       | This Field indicates the minimum lane speed            |
		// +----------------------------------------------------------------+
		// | 7:4   | Reserved. Shall be set to zero                         |
		// +----------------------------------------------------------------+
		// | 11:8  | Min RX Lane Count.                                     |
		// |       | This field indicates the minimum receive lane count.   |
		// +----------------------------------------------------------------+
		// | 15:12 | Min TX Lane Count.                                     |
		// |       | This field indicates the minimum transmit lane count.  |
		// +----------------------------------------------------------------+
		WFunctionalitySupport uint16

		// BReserved2.
		BReserved2 uint16

		// BMSublinkSpeedAttr is an array of Sublink Speed Attribute bitfield.
		// One Attribute:
		// +----------------------------------------------------------------+
		// | Bit   | Encoding                                               |
		// +----------------------------------------------------------------+
		// | 3:0   | Sublink Speed Attribute ID (SSID).                     |
		// |       | This field is an ID That uniquely identifies the speed |
		// |       | of the sublink. Note that a maximum of 16 unique SSIDs |
		// |       | may be defined.                                        |
		// +----------------------------------------------------------------+
		// | 5:4   | Lane Speed Exponent (LSE).                             |
		// |       | This field defines the base 10 exponent times 3, that  |
		// |       | shall be applied to the Lane Speed Mantissa (LSM) when |
		// |       | calculating the maximum bit rate represented by this   |
		// |       | Lane Speed Attribute.                                  |
		// +----------------------------------------------------------------+
		// | 7:6   | Sublink Type (ST).                                     |
		// |       | This field identifies whether the Sublink Speed        |
		// |       | Attribute defines a symmetric or asymmetric bit rate.  |
		// |       | This field also indicates if this Sublink Speed        |
		// |       | Attribute defines the receive or transmit bit rate.    |
		// |       | Note that the Sublink Speed Attributes shall be paired,|
		// |       | i.e. an Rx immediately followed by a Tx, and both      |
		// |       | Attributes shall define the same value for the SSID    |
		// |       |                                                        |
		// |       | Bit6:                                                  |
		// |       |   0 - Symmetric. Rx and Tx Sublinks have the same      |
		// |       |       number of lanes and operate at the same speed.   |
		// |       |   1 - Asymmetric. Rx and Tx Sublinks have different    |
		// |       |       number of lanes and/or operates at different     |
		// |       |       speeds.                                          |
		// |       | Bit7:                                                  |
		// |       |   0 - Sublink operates in Receive mode.                |
		// |       |   1 - Sublink operates in Transmit mode.               |
		// +----------------------------------------------------------------+
		// | 13:8  | Reserved.                                              |
		// +----------------------------------------------------------------+
		// | 15:14 | Link Protocol (LP).                                    |
		// |       | This field identifies the protocol supported           |
		// |       | by the link.                                           |
		// |       |                                                        |
		// |       | 0   - SuperSpeed                                       |
		// |       | 1   - SuperSpeedPlus                                   |
		// |       | 3-2 - Reserved                                         |
		// +----------------------------------------------------------------+
		// | 31:16 | Lane Speed Mantissa (LSM).                             |
		// |       | This field defines the mantissa that shall be applied  |
		// |       | to the LSE when calculating the maximum bitrate        |
		// |       | represented by Lane Speed Attribute.                   |
		// +----------------------------------------------------------------+
		BMSublinkSpeedAttr []uint32
	}

	// CapPrecisionTimeDescriptor defines the device-level capabilities which shall be implemented
	// by all hubs and devices that support PTM capability.
	CapPrecisionTimeDescriptor struct {
		DescriptorHeader
		// BDevCapabilityType Capability type: CapPrecisionTime
		BDevCapabilityType Capability
	}

	// CapConfigurationSummaryDescriptor may be implemented by a device with more than one
	// configuration, and identifies a single function presented by the device along with a list of the
	// configuration descriptor indices that include the function.
	// If implemented, each function presented by the device shall be represented by a
	// separate CapConfigurationSummaryDescriptor.
	// However, a function’s Configuration Summary Descriptor may be omitted if the
	// function is present in all possible configurations.
	// CapConfigurationSummaryDescriptor(s) should be included in the BOS descriptor in order of descending preference.
	CapConfigurationSummaryDescriptor struct {
		DescriptorHeader
		// BDevCapabilityType Capability type: CapConfigurationSummary
		BDevCapabilityType Capability

		// 0100H, the revision of the Configuration Summary Descriptor with this document.
		BCDVersion uint16

		// BClass Class code of the function
		BClass ClassCode

		// BSubClass Subclass code of the function
		BSubClass SubClass

		// BProtocol Protocol of the function
		BProtocol uint8

		// Number of configurations (N) that include this
		// class/subclass/protocol
		BConfigurationCount uint8

		BConfigurationIndex []uint8
	}
)

var capabilityStringMap = map[Capability]string{
	CapWirelessUSB:          "Wireless USB",
	CapUSB20Extension:       "USB 2.0 Extension",
	CapSuperSpeedUSB:        "SuperSpeed USB",
	CapContainerID:          "Container ID",
	CapPlatform:             "Platform",
	CapPowerDelivery:        "Power Delivery",
	CapBatteryInfo:          "BatteryInfo",
	CapPDConsumerPort:       "Power Delivery Consumer Port",
	CapPDProviderPort:       "Power Delivery Provider Port",
	CapSuperSpeedPlus:       "SuperSpeed Plus",
	CapPrecisionTime:        "Precision Time",
	CapWirelessUSBExt:       "Wireless USB Extension",
	CapBillboard:            "Billboard",
	CapAuthentication:       "Authentication",
	CapBillboardEx:          "Billboard Extension",
	CapConfigurationSummary: "Configuration summary",
}

var capabilityMap = map[Capability]reflect.Type{
	CapUSB20Extension:       reflect.TypeOf(CapUSB20ExtensionDescriptor{}),
	CapSuperSpeedUSB:        reflect.TypeOf(CapSuperSpeedUSBDescriptor{}),
	CapContainerID:          reflect.TypeOf(CapContainerIDDescriptor{}),
	CapPlatform:             reflect.TypeOf(CapPlatformDescriptor{}),
	CapSuperSpeedPlus:       reflect.TypeOf(CapSuperSpeedPlusUSBDescriptor{}),
	CapPrecisionTime:        reflect.TypeOf(CapPrecisionTimeDescriptor{}),
	CapConfigurationSummary: reflect.TypeOf(CapConfigurationSummaryDescriptor{}),
}
