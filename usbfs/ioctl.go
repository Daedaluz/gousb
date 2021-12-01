package usbfs

import (
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
		Driver    [USBDEVFS_MAXDRIVERNAME + 1]byte
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
		Driver    [USBDEVFS_MAXDRIVERNAME + 1]uint8
	}
	usbdevfs_streams struct {
		NumStreams   uint32
		NumEndpoints uint32
		/* Endpoints... */
	}
)

const (
	USBDEVFS_MAXDRIVERNAME = 255
)

const (
	USBDEVFS_CONTROL          = 0xC0185500
	USBDEVFS_BULK             = 0xC0185502
	USBDEVFS_RESETEP          = 0x80045503
	USBDEVFS_SETINTERFACE     = 0x80085504
	USBDEVFS_SETCONFIGURATION = 0x80045505
	USBDEVFS_GETDRIVER        = 0x41045508
	USBDEVFS_SUBMITURB        = 0x8000550A
	USBDEVFS_DISCARDURB       = 0x0000550B
	USBDEVFS_REAPURB          = 0x4008550C
	USBDEVFS_REAPURBNDELAY    = 0x4008550D
	USBDEVFS_DISCSIGNAL       = 0x8010550E
	USBDEVFS_CLAIMINTERFACE   = 0x8004550F
	USBDEVFS_RELEASEINTERFACE = 0x80045510
	USBDEVFS_CONNECTINFO      = 0x40085511
	USBDEVFS_IOCTL            = 0xC0105512
	USBDEVFS_HUB_PORTINFO     = 0x80805513
	USBDEVFS_RESET            = 0x00005514
	USBDEVFS_CLEAR_HALT       = 0x80045515
	USBDEVFS_DISCONNECT       = 0x00005516
	USBDEVFS_CONNECT          = 0x00005517
	USBDEVFS_CLAIM_PORT       = 0x80045518
	USBDEVFS_RELEASE_PORT     = 0x80045519
	USBDEVFS_GET_CAPABILITIES = 0x8004551A
	USBDEVFS_DISCONNECT_CLAIM = 0x8108551B
	USBDEVFS_ALLOC_STREAMS    = 0x8008551C
	USBDEVFS_FREE_STREAMS     = 0x8008551D
	USBDEVFS_DROP_PRIVILEGES  = 0x4004551E
	USBDEVFS_GET_SPEED        = 0x0000551F
)

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
