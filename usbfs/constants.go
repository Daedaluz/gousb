package usbfs

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
