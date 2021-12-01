package usb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	sysfsDeviceDir = "/sys/bus/usb/devices"
)

func readSysfsAttrInt(devName, attrName string) (int, error) {
	fileName := fmt.Sprintf("%s/%s/%s", sysfsDeviceDir, devName, attrName)
	var err error
	var data []byte
	var value int64
	data, err = ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}
	strData := strings.Trim(string(data), "\n")
	value, err = strconv.ParseInt(strData, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(value), nil
}

func openSysfsAttr(devName, attrName string) (*os.File, error) {
	fileName := fmt.Sprintf("%s/%s/%s", sysfsDeviceDir, devName, attrName)
	return os.Open(fileName)
}

func getDeviceAddress(devName string) (int, int, error) {
	busNum, err := readSysfsAttrInt(devName, "busnum")
	if err != nil {
		return 0, 0, err
	}
	devNum, err := readSysfsAttrInt(devName, "devnum")
	if err != nil {
		return 0, 0, err
	}
	return busNum, devNum, nil
}

func readDescriptorHeader(i io.Reader) (DescriptorHeader, error) {
	header := DescriptorHeader{
		Length:         0,
		DescriptorType: 0,
	}
	err := binary.Read(i, binary.BigEndian, &header)
	return header, err
}

func parseDescriptor(devName string) ([]Descriptor, error) {
	var hdr DescriptorHeader
	var err error
	var x *os.File
	res := make([]Descriptor, 0, 10)
	x, err = openSysfsAttr(devName, "descriptors")
	if err != nil {
		return nil, err
	}
	defer x.Close()
	for hdr, err = readDescriptorHeader(x); err == nil; hdr, err = readDescriptorHeader(x) {
		// Create a separate input stream for descriptor to prevent overstepping descriptor boundary.
		descriptorData := make([]byte, hdr.Length-2)
		if _, err := io.ReadFull(x, descriptorData); err != nil {
			log.Println("Bad descriptor data:", err)
			continue
		}
		descriptorReader := bytes.NewReader(descriptorData)
		desc, descErr := createDescriptor(hdr, descriptorReader)
		if descErr != nil {
			return nil, descErr
		}
		res = append(res, desc)
	}
	if err != io.EOF {
		return nil, err
	}
	return res, nil
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
		descriptors, err := parseDescriptor(name)
		if err != nil {
			return nil, err
		}
		busNum, devNum, err := getDeviceAddress(name)
		if err != nil {
			return nil, err
		}
		device := &Device{
			BusNumber:    busNum,
			DeviceNumber: devNum,
			Descriptors:  descriptors,
			fd:           -1,
		}
		res = append(res, device)
	}
	return res, nil
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
