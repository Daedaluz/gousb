package usb

import "fmt"

// From https://www.usb.org/defined-class-codes

type (
	ClassCode uint8
)

func (code ClassCode) String() string {
	if codeString, exist := classCodeMap[code]; exist {
		return codeString
	}
	return fmt.Sprintf("Unknown(%.2X)", uint8(code))
}

// Both class codes
const (
	ClassCodeCDCControl     = ClassCode(0x02)
	ClassCodeDiagnostic     = ClassCode(0xDC)
	ClassCodeMisc           = ClassCode(0xEF)
	ClassCodeVendorSpecific = ClassCode(0xFF)
)

// Interface class codes
const (
	ClassCodeInterfaceAudio               = ClassCode(0x01)
	ClassCodeInterfaceHID                 = ClassCode(0x03)
	ClassCodeInterfacePhysical            = ClassCode(0x05)
	ClassCodeInterfaceImage               = ClassCode(0x06)
	ClassCodeInterfacePrinter             = ClassCode(0x07)
	ClassCodeInterfaceMassStorage         = ClassCode(0x08)
	ClassCodeInterfaceCDCData             = ClassCode(0x0A)
	ClassCodeInterfaceSmartCard           = ClassCode(0x0B)
	ClassCodeInterfaceContentSecurity     = ClassCode(0x0D)
	ClassCodeInterfaceVideo               = ClassCode(0x0E)
	ClassCodeInterfacePersonalHealthcare  = ClassCode(0x0F)
	ClassCodeInterfaceAudioVideo          = ClassCode(0x10)
	ClassCodeInterfaceTypeCBridgeClass    = ClassCode(0x12)
	ClassCodeInterfaceWirelessController  = ClassCode(0xE0)
	ClassCodeInterfaceApplicationSpecific = ClassCode(0xFE)
)

const (
	ClassCodeDeviceHub       = ClassCode(0x09)
	ClassCodeDeviceBillBoard = ClassCode(0x11)
)

var (
	classCodeMap = map[ClassCode]string{
		0x00: "UseInterfaceDescriptors",
		0x01: "InterfaceAudio",
		0x03: "InterfaceHID",
		0x05: "InterfacePhysical",
		0x06: "InterfaceImage",
		0x07: "InterfacePrinter",
		0x08: "InterfaceMassStorage",
		0x0A: "InterfaceCDCData",
		0x0B: "InterfaceSmartCard",
		0x0D: "InterfaceContentSecurity",
		0x0E: "InterfaceVideo",
		0x0F: "InterfacePersonalHealthcare",
		0x10: "InterfaceAudioVideo",
		0x12: "InterfaceTypeCBridgeClass",
		0xE0: "InterfaceWirelessController",
		0xFE: "InterfaceApplicationSpecific",
		0x09: "DeviceHub",
		0x11: "DeviceBillBoard",
		0x02: "CDCControl",
		0xDC: "Diagnostic",
		0xEF: "Misc",
		0xFF: "VendorSpecific",
	}
)
