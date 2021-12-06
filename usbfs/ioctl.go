package usbfs

import (
	ioctl "github.com/daedaluz/goioctl"
	"strings"
	"unsafe"
)

type (
	usbdevfs_ctrltransfer struct {
		RequestType uint8
		Request     uint8
		Value       uint16
		Index       uint16
		Length      uint16
		Timeout     uint32
		Data        uintptr
	}
	usbdevfs_bulktransfer struct {
		Endpoint uint32
		Length   uint32
		Timeout  uint32
		Data     uintptr
	}

	usbdevfs_setinterface struct {
		Interface  uint32
		AltSetting uint32
	}

	usbdevfs_getdriver struct {
		Interface uint32
		Driver    [nUSBDEVFS_MAXDRIVERNAME + 1]byte
	}

	usbdevfs_urb struct {
		Type            uint8
		Endpoint        uint8
		Status          int32
		Flags           uint32
		Buffer          uintptr
		BufferLength    int32
		ActualLength    int32
		StartFrame      int32
		PacketsOrStream uint32 /* StreamID if bulk, number of packets if isoc */
		ErrorCount      int32
		SigNumber       uint32
		UserContext     uintptr
		/* iso_frame_desc... */
	}

	usbdevfs_disconnectsignal struct {
		Signr   uint32
		Context uintptr
	}

	usbdevfs_connectinfo struct {
		DevNum uint32
		Slow   uint8
	}

	usbdevfs_ioctl struct {
		Interface uint32
		IoctlCode uint32
		Data      uintptr
	}
	usbdevfs_hub_portinfo struct {
		NPorts uint8
		Port   [127]uint8
	}
	usbdevfs_disconnect_claim struct {
		Interface uint32
		Flags     uint32
		Driver    [nUSBDEVFS_MAXDRIVERNAME + 1]uint8
	}
	usbdevfs_streams struct {
		NumStreams   uint32
		NumEndpoints uint32
		/* Endpoints... */
	}
)

const (
	nUSBDEVFS_MAXDRIVERNAME = 255
)

//const (
//	USBDEVFS_CONTROL          = 0xC0185500
//	USBDEVFS_BULK             = 0xC0185502
//	USBDEVFS_RESETEP          = 0x80045503
//	USBDEVFS_SETINTERFACE     = 0x80085504
//	USBDEVFS_SETCONFIGURATION = 0x80045505
//	USBDEVFS_GETDRIVER        = 0x41045508
//	USBDEVFS_SUBMITURB        = 0x8000550A
//	USBDEVFS_DISCARDURB       = 0x0000550B
//	USBDEVFS_REAPURB          = 0x4008550C
//	USBDEVFS_REAPURBNDELAY    = 0x4008550D
//	USBDEVFS_DISCSIGNAL       = 0x8010550E
//	USBDEVFS_CLAIMINTERFACE   = 0x8004550F
//	USBDEVFS_RELEASEINTERFACE = 0x80045510
//	USBDEVFS_CONNECTINFO      = 0x40085511
//	USBDEVFS_IOCTL            = 0xC0105512
//	USBDEVFS_HUB_PORTINFO     = 0x80805513
//	USBDEVFS_RESET            = 0x00005514
//	USBDEVFS_CLEAR_HALT       = 0x80045515
const USBDEVFS_DISCONNECT = 0x00005516
const USBDEVFS_CONNECT = 0x00005517

//	USBDEVFS_CLAIM_PORT       = 0x80045518
//	USBDEVFS_RELEASE_PORT     = 0x80045519
//	USBDEVFS_GET_CAPABILITIES = 0x8004551A
//	USBDEVFS_DISCONNECT_CLAIM = 0x8108551B
//	USBDEVFS_ALLOC_STREAMS    = 0x8008551C
//	USBDEVFS_FREE_STREAMS     = 0x8008551D
//	USBDEVFS_DROP_PRIVILEGES  = 0x4004551E
//	USBDEVFS_GET_SPEED        = 0x0000551F
//)

func (d *usbdevfs_getdriver) String() string {
	buff := strings.Builder{}
	for _, x := range d.Driver {
		if x == 0 {
			break
		}
		buff.WriteByte(x)
	}
	return buff.String()
}

func slicePtr(s []byte) uintptr {
	return uintptr(unsafe.Pointer(&s[0]))
}

var (
	ctl_usbdevfs_control          = ioctl.IOWR('U', 0, unsafe.Sizeof(usbdevfs_ctrltransfer{}))
	ctl_usbdevfs_bulk             = ioctl.IOWR('U', 2, unsafe.Sizeof(usbdevfs_bulktransfer{}))
	ctl_usbdevfs_resetep          = ioctl.IOR('U', 3, unsafe.Sizeof(uint32(0)))
	ctl_usbdevfs_setinterface     = ioctl.IOR('U', 4, unsafe.Sizeof(usbdevfs_setinterface{}))
	ctl_usbdevfs_setconfiguration = ioctl.IOR('U', 5, unsafe.Sizeof(uint32(0)))
	ctl_usbdevfs_getdriver        = ioctl.IOW('U', 8, unsafe.Sizeof(usbdevfs_getdriver{}))
	ctl_usbdevfs_submiturb        = ioctl.IOR('U', 10, unsafe.Sizeof(usbdevfs_urb{}))
	ctl_usbdevfs_discardurb       = ioctl.IO('U', 11)
	ctl_usbdevfs_reapurb          = ioctl.IOW('U', 12, unsafe.Sizeof(uintptr(0)))
	ctl_usbdevfs_reapurbndelay    = ioctl.IOW('U', 13, unsafe.Sizeof(uintptr(0)))
	ctl_usbdevfs_discsignal       = ioctl.IOR('U', 14, unsafe.Sizeof(usbdevfs_disconnectsignal{}))
	ctl_usbdevfs_claiminterface   = ioctl.IOR('U', 15, unsafe.Sizeof(uint32(0)))
	ctl_usbdevfs_releaseinterface = ioctl.IOR('U', 16, unsafe.Sizeof(uint32(0)))
	ctl_usbdevfs_connectionfo     = ioctl.IOW('U', 17, unsafe.Sizeof(usbdevfs_connectinfo{}))
	ctl_usbdevfs_ioctl            = ioctl.IOWR('U', 18, unsafe.Sizeof(usbdevfs_ioctl{}))
	ctl_usbdevfs_portinfo         = ioctl.IOR('U', 19, unsafe.Sizeof(usbdevfs_hub_portinfo{}))
	ctl_usbdevfs_reset            = ioctl.IO('U', 20)
	ctl_usbdevfs_clear_halt       = ioctl.IOR('U', 21, unsafe.Sizeof(uint32(0)))
	ctl_usbdevfs_disconnect       = ioctl.IO('U', 22)
	ctl_usbdevfs_connect          = ioctl.IO('U', 23)
	ctl_usbdevfs_claim_port       = ioctl.IOR('U', 24, unsafe.Sizeof(uint32(0)))
	ctl_usbdevfs_release_port     = ioctl.IOR('U', 25, unsafe.Sizeof(uint32(0)))
	ctl_usbdevfs_get_capabilities = ioctl.IOR('U', 26, unsafe.Sizeof(uint32(0)))
	ctl_usbdevfs_disconnect_claim = ioctl.IOR('U', 27, unsafe.Sizeof(usbdevfs_disconnect_claim{}))
	ctl_usbdevfs_alloc_streams    = ioctl.IOR('U', 28, unsafe.Sizeof(usbdevfs_streams{}))
	ctl_usbdevfs_free_streams     = ioctl.IOR('U', 29, unsafe.Sizeof(usbdevfs_streams{}))
	ctl_usbdevfs_drop_privileges  = ioctl.IOW('U', 30, unsafe.Sizeof(uint32(0)))
	ctl_usbdevfs_get_speed        = ioctl.IO('U', 31)
)
