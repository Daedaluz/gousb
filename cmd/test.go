package main

import (
	usb "github.com/daedaluz/gousb"
	"log"
)

func main() {
	usb.FindDevices(func(device *usb.Device) bool {
		descriptors, err := device.GetSysfsDescriptors()
		if err != nil {
			return false
		}
		log.Println(descriptors.Manufacturer, descriptors.Product, descriptors.DeviceDescriptor.IDVendor, descriptors.DeviceDescriptor.IDProduct)
		return true
	})
}
