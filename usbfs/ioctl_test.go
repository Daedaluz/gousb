package usbfs

import (
	"testing"
	"unsafe"
)

const (
	iocNrBits   = 8
	iocTypeBits = 8
	iocSizeBits = 14
	iocDirBits  = 2

	iocNrShift   = 0
	iocTypeShift = iocNrShift + iocNrBits
	iocSizeShift = iocTypeShift + iocTypeBits
	iocDirShift  = iocSizeShift + iocSizeBits

	iocNone  = 0
	iocWrite = 1
	iocRead  = 2
)

func _IO(t, nr uintptr) uintptr {
	return _IOC(iocNone, t, nr, 0)
}

func _IOR(t, nr, size uintptr) uintptr {
	return _IOC(iocRead, t, nr, size)
}

func _IOW(t, nr, size uintptr) uintptr {
	return _IOC(iocWrite, t, nr, size)
}

func _IOWR(t, nr, size uintptr) uintptr {
	return _IOC(iocRead|iocWrite, t, nr, size)
}

func _IOC(dir, t, nr, size uintptr) uintptr {
	return (dir << iocDirShift) | (t << iocTypeShift) | (nr << iocNrShift) | (size << iocSizeShift)
}

type ioctlstruct struct {
	name   string
	number uintptr
	target uintptr
}

var ioctls = []ioctlstruct{
	{"USBDEVFS_CONTROL", _IOWR('U', 0, unsafe.Sizeof(usbdevfs_ctrltransfer{})), 0xC0185500},
	{"USBDEVFS_BULK", _IOWR('U', 2, unsafe.Sizeof(usbdevfs_bulktransfer{})), 0xC0185502},
	{"USBDEVFS_RESETEP", _IOR('U', 3, unsafe.Sizeof(uint32(0))), 0x80045503},
	{"USBDEVFS_SETINTERFACE", _IOR('U', 4, unsafe.Sizeof(usbdevfs_setinterface{})), 0x80085504},
	{"USBDEVFS_SETCONFIGURATION", _IOR('U', 5, unsafe.Sizeof(uint32(0))), 0x80045505},
	{"USBDEVFS_GETDRIVER", _IOW('U', 8, unsafe.Sizeof(usbdevfs_getdriver{})), 0x41045508},
	{"USBDEVFS_SUBMITURB", _IOR('U', 10, unsafe.Sizeof(usbdevfs_urb{})), 0x8038550A},
	{"USBDEVFS_DISCARDURB", _IO('U', 11), 0x0000550B},
	{"USBDEVFS_REAPURB", _IOW('U', 12, unsafe.Sizeof(uintptr(0))), 0x4008550C},
	{"USBDEVFS_REAPURBNDELAY", _IOW('U', 13, unsafe.Sizeof(uintptr(0))), 0x4008550D},
	{"USBDEVFS_DISCSIGNAL", _IOR('U', 14, unsafe.Sizeof(usbdevfs_disconnectsignal{})), 0x8010550E},
	{"USBDEVFS_CLAIMINTERFACE", _IOR('U', 15, unsafe.Sizeof(uint32(0))), 0x8004550F},
	{"USBDEVFS_RELEASEINTERFACE", _IOR('U', 16, unsafe.Sizeof(uint32(0))), 0x80045510},
	{"USBDEVFS_CONNECTINFO", _IOW('U', 17, unsafe.Sizeof(usbdevfs_connectinfo{})), 0x40085511},
	{"USBDEVFS_IOCTL", _IOWR('U', 18, unsafe.Sizeof(usbdevfs_ioctl{})), 0xC0105512},
	{"USBDEVFS_HUB_PORTINFO", _IOR('U', 19, unsafe.Sizeof(usbdevfs_hub_portinfo{})), 0x80805513},
	{"USBDEVFS_RESET", _IO('U', 20), 0x00005514},
	{"USBDEVFS_CLEAR_HALT", _IOR('U', 21, unsafe.Sizeof(uint32(0))), 0x80045515},
	{"USBDEVFS_DISCONNECT", _IO('U', 22), 0x00005516},
	{"USBDEVFS_CONNECT", _IO('U', 23), 0x00005517},
	{"USBDEVFS_CLAIM_PORT", _IOR('U', 24, unsafe.Sizeof(uint32(0))), 0x80045518},
	{"USBDEVFS_RELEASE_PORT", _IOR('U', 25, unsafe.Sizeof(uint32(0))), 0x80045519},
	{"USBDEVFS_GET_CAPABILITIES", _IOR('U', 26, unsafe.Sizeof(uint32(0))), 0x8004551A},
	{"USBDEVFS_DISCONNECT_CLAIM", _IOR('U', 27, unsafe.Sizeof(usbdevfs_disconnect_claim{})), 0x8108551B},
	{"USBDEVFS_ALLOC_STREAMS", _IOR('U', 28, unsafe.Sizeof(usbdevfs_streams{})), 0x8008551C},
	{"USBDEVFS_FREE_STREAMS", _IOR('U', 29, unsafe.Sizeof(usbdevfs_streams{})), 0x8008551D},
	{"USBDEVFS_DROP_PRIVILEGES", _IOW('U', 30, unsafe.Sizeof(uint32(0))), 0x4004551E},
	{"USBDEVFS_GET_SPEED", _IO('U', 31), 0x0000551F},
}

func TestIOCTLNumbers(t *testing.T) {
	for _, ctl := range ioctls {
		if ctl.number != ctl.target {
			t.Logf("WRONG NUMBER - %s, %.8X != %.8X\n",ctl.name, ctl.number, ctl.target)
			t.Fail()
		}
		t.Logf("%s = 0x%.8X\n", ctl.name, ctl.number)
	}
}

/* usbdevice_fs.h
#define USBDEVFS_CONTROL           _IOWR('U', 0, struct usbdevfs_ctrltransfer)
#define USBDEVFS_CONTROL32         _IOWR('U', 0, struct usbdevfs_ctrltransfer32)
#define USBDEVFS_BULK              _IOWR('U', 2, struct usbdevfs_bulktransfer)
#define USBDEVFS_BULK32            _IOWR('U', 2, struct usbdevfs_bulktransfer32)
#define USBDEVFS_RESETEP           _IOR('U', 3, unsigned int)
#define USBDEVFS_SETINTERFACE      _IOR('U', 4, struct usbdevfs_setinterface)
#define USBDEVFS_SETCONFIGURATION  _IOR('U', 5, unsigned int)
#define USBDEVFS_GETDRIVER         _IOW('U', 8, struct usbdevfs_getdriver)
#define USBDEVFS_SUBMITURB         _IOR('U', 10, struct usbdevfs_urb)
#define USBDEVFS_SUBMITURB32       _IOR('U', 10, struct usbdevfs_urb32)
#define USBDEVFS_DISCARDURB        _IO('U', 11)
#define USBDEVFS_REAPURB           _IOW('U', 12, void *)
#define USBDEVFS_REAPURB32         _IOW('U', 12, __u32)
#define USBDEVFS_REAPURBNDELAY     _IOW('U', 13, void *)
#define USBDEVFS_REAPURBNDELAY32   _IOW('U', 13, __u32)
#define USBDEVFS_DISCSIGNAL        _IOR('U', 14, struct usbdevfs_disconnectsignal)
#define USBDEVFS_DISCSIGNAL32      _IOR('U', 14, struct usbdevfs_disconnectsignal32)
#define USBDEVFS_CLAIMINTERFACE    _IOR('U', 15, unsigned int)
#define USBDEVFS_RELEASEINTERFACE  _IOR('U', 16, unsigned int)
#define USBDEVFS_CONNECTINFO       _IOW('U', 17, struct usbdevfs_connectinfo)
#define USBDEVFS_IOCTL             _IOWR('U', 18, struct usbdevfs_ioctl)
#define USBDEVFS_IOCTL32           _IOWR('U', 18, struct usbdevfs_ioctl32)
#define USBDEVFS_HUB_PORTINFO      _IOR('U', 19, struct usbdevfs_hub_portinfo)
#define USBDEVFS_RESET             _IO('U', 20)
#define USBDEVFS_CLEAR_HALT        _IOR('U', 21, unsigned int)
#define USBDEVFS_DISCONNECT        _IO('U', 22)
#define USBDEVFS_CONNECT           _IO('U', 23)
#define USBDEVFS_CLAIM_PORT        _IOR('U', 24, unsigned int)
#define USBDEVFS_RELEASE_PORT      _IOR('U', 25, unsigned int)
#define USBDEVFS_GET_CAPABILITIES  _IOR('U', 26, __u32)
#define USBDEVFS_DISCONNECT_CLAIM  _IOR('U', 27, struct usbdevfs_disconnect_claim)
#define USBDEVFS_ALLOC_STREAMS     _IOR('U', 28, struct usbdevfs_streams)
#define USBDEVFS_FREE_STREAMS      _IOR('U', 29, struct usbdevfs_streams)
#define USBDEVFS_DROP_PRIVILEGES   _IOW('U', 30, __u32)
#define USBDEVFS_GET_SPEED         _IO('U', 31)
*/

