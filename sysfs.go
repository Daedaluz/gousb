package usb

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	sysfsDeviceDir = "/sys/bus/usb/devices"
)

func formatAttrFileName(devName, attrName string) string {
	return fmt.Sprintf("%s/%s/%s", sysfsDeviceDir, devName, attrName)
}

func readSysfsAttrInt(devName, attrName string, base, bitSize int) (int64, error) {
	fileName := formatAttrFileName(devName, attrName)
	var err error
	var data []byte
	var value int64
	data, err = ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}
	strData := strings.Trim(string(data), "\n")
	value, err = strconv.ParseInt(strData, base, bitSize)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func readSysfsAttrString(devName, attrName string) (string, error) {
	fileName := formatAttrFileName(devName, attrName)
	var err error
	var data []byte
	data, err = ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	strData := strings.Trim(string(data), "\n")
	return strData, nil
}

func openSysfsAttr(devName, attrName string) (*os.File, error) {
	fileName := formatAttrFileName(devName, attrName)
	file, err := os.Open(fileName)
	return file, err
}

func getDeviceAddress(devName string) (int, int, error) {
	busNum, err := readSysfsAttrInt(devName, "busnum", 10, 64)
	if err != nil {
		return 0, 0, err
	}
	devNum, err := readSysfsAttrInt(devName, "devnum", 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(busNum), int(devNum), nil
}

func EnumerateDevices() ([]*Device, error) {
	dirs, err := ioutil.ReadDir(sysfsDeviceDir)
	if err != nil {
		return nil, err
	}

	res := make([]*Device, 0, 10)

	for _, dir := range dirs {
		name := dir.Name()
		if strings.HasPrefix(name, "usb") ||
			strings.Contains(name, ":") {
			continue
		}
		busNum, devNum, err := getDeviceAddress(name)
		if err != nil {
			return nil, err
		}
		device := &Device{
			Name:         name,
			BusNumber:    busNum,
			DeviceNumber: devNum,
			fd:           -1,
		}
		res = append(res, device)
	}
	return res, nil
}

func (d *Device) ReadSysfsAttrInt(attrName string, base, bitSize int) (int64, error) {
	return readSysfsAttrInt(d.Name, attrName, base, bitSize)
}

func (d *Device) ReadSysfsString(attrName string) (string, error) {
	return readSysfsAttrString(d.Name, attrName)
}

func FindDevices(filter func(device *Device) bool) ([]*Device, error) {
	allDevices, err := EnumerateDevices()
	if err != nil {
		return nil, err
	}
	res := make([]*Device, 0, len(allDevices))
	for _, dev := range allDevices {
		if filter(dev) {
			res = append(res, dev)
		}
	}
	return res, nil
}

type SysfsDescriptors struct {
	Manufacturer     string
	Product          string
	DeviceDescriptor *DeviceDescriptor
	Interfaces       []*InterfaceDescriptor
	Endpoints        map[*InterfaceDescriptor][]*EndpointDescriptor
	OtherDescriptors []Descriptor
}

func (d *Device) GetSysfsDescriptors() (*SysfsDescriptors, error) {
	res := &SysfsDescriptors{
		DeviceDescriptor: nil,
		Interfaces:       make([]*InterfaceDescriptor, 0, 10),
		Endpoints:        make(map[*InterfaceDescriptor][]*EndpointDescriptor, 20),
		OtherDescriptors: make([]Descriptor, 0, 20),
	}
	if x, err := d.ReadSysfsString("manufacturer"); err == nil {
		res.Manufacturer = x
	}
	if x, err := d.ReadSysfsString("product"); err == nil {
		res.Product = x
	}
	if file, err := openSysfsAttr(d.Name, "descriptors"); err == nil {
		var lastInterfaceDescriptor *InterfaceDescriptor
		defer file.Close()
		err := ReadDescriptors(file, func(d Descriptor) {
			switch x := d.(type) {
			case *DeviceDescriptor:
				res.DeviceDescriptor = x
			case *InterfaceDescriptor:
				lastInterfaceDescriptor = x
				res.Interfaces = append(res.Interfaces, x)
			case *EndpointDescriptor:
				ep, exist := res.Endpoints[lastInterfaceDescriptor]
				if !exist {
					ep = make([]*EndpointDescriptor, 0, 4)
				}
				ep = append(ep, x)
				res.Endpoints[lastInterfaceDescriptor] = ep
			default:
				res.OtherDescriptors = append(res.OtherDescriptors, x)
			}
		})
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
