package usb

func (d *Device) GetStringDescriptor(idx uint8) (string, error) {
	desc, err := d.GetDescriptor(DescriptorTypeString, idx)
	if err != nil {
		return "", err
	}
	strDesc := desc.(*StringDescriptor)
	str := strDesc.Data[0 : strDesc.Length-2]
	return string(str), nil
}

func (d *Device) GetDescriptor(descriptorType DescriptorType, idx uint8) (Descriptor, error) {
	buff := make([]byte, 256)
	_, err := d.Ctrl(RequestDirectionIn, RequestDeviceGetDescriptor, (uint16(descriptorType)<<8)|uint16(idx), 0, buff)
	if err != nil {
		return nil, err
	}
	return ParseDescriptor(buff)
}

func (d *Device) GetDescriptorData(descriptorType DescriptorType, idx, size uint16) ([]byte, error) {
	buff := make([]byte, size)
	_, err := d.Ctrl(RequestDirectionIn, RequestDeviceGetDescriptor, (uint16(descriptorType)<<8)|idx, 0, buff)
	return buff, err
}

func (d *Device) GetConfiguration() (int, error) {
	buff := make([]byte, 1)
	_, err := d.Ctrl(RequestDirectionIn, RequestDeviceGetConfiguration, 0, 0, buff)
	return int(buff[0]), err
}

func (d *Device) SetConfiguration(ci int) error {
	_, err := d.Ctrl(RequestDirectionOut, RequestDeviceSetConfiguration, uint16(ci), 0, nil)
	return err
}

func (d *Device) GetAltInterface(interfaceIndex int) (int, error) {
	data := make([]byte, 1)
	_, err := d.Ctrl(RequestDirectionIn|RequestRecipientInterface, RequestInterfaceGetInterface, 0, uint16(interfaceIndex), data)
	return int(data[0]), err
}

func (d *Device) SetAltInterface(interfaceIndex, setting int) error {
	_, err := d.Ctrl(RequestDirectionOut|RequestRecipientInterface, RequestInterfaceSetInterface, uint16(setting), uint16(interfaceIndex), nil)
	return err
}
