package usbfs

import "strings"

const (
	usbDevPath = "/dev/bus/usb"
)

const (
	nUSBDEVFS_MAXDRIVERNAME = 255
)

type Capability uint32

const (
	CapZeroPacket          = Capability(0x01)
	CapBulkContinuation    = Capability(0x02)
	CapNoPacketSizeLim     = Capability(0x04)
	CapBulkScatterGather   = Capability(0x08)
	CapReapAfterDisconnect = Capability(0x10)
	CapNMAP                = Capability(0x20)
	CapDropPrivileges      = Capability(0x40)
	CapConnInfoEx          = Capability(0x80)
	CapSuspend             = Capability(0x100)
)

func (c Capability) String() string {
	capStrings := make([]string, 0, 9)
	for i := CapZeroPacket; i <= CapSuspend; i <<= 1 {
		str := ""
		switch c & i {
		case CapZeroPacket:
			str = "CapZeroPacket"
		case CapBulkContinuation:
			str = "CapBulkContinuation"
		case CapNoPacketSizeLim:
			str = "CapNoPacketSizeLim"
		case CapBulkScatterGather:
			str = "CapBulkScatterGather"
		case CapReapAfterDisconnect:
			str = "CapReapAfterDisconnect"
		case CapNMAP:
			str = "CapNMAP"
		case CapDropPrivileges:
			str = "CapDropPrivileges"
		case CapConnInfoEx:
			str = "CapConnInfoEx"
		case CapSuspend:
			str = "CapSuspend"
		}
		capStrings = append(capStrings, str)
	}
	return strings.Join(capStrings, "|")
}
