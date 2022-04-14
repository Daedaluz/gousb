package main

import (
	"encoding/json"
	usb "github.com/daedaluz/gousb"
	"log"
)

func main() {
	usb.FindDevices(func(device *usb.Device) bool {
		if err := device.Open(); err != nil {
			return true
		}
		strDesc, err := device.GetDescriptor(usb.DescriptorTypeBOS, 0, 0)
		x, _ := device.GetDevicePTMStatus()
		data, _ := json.Marshal(x)
		log.Println(strDesc, err)
		log.Println(string(data))
		device.Close()
		return true
	})
}
