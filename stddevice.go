package usb

import "encoding/binary"

// Standard request codes
const (
	ReqGetStatus        = 0x00
	ReqClearFeature     = 0x01
	ReqSetFeature       = 0x03
	ReqSetAddress       = 0x05
	ReqGetDescriptor    = 0x06
	ReqSetDescriptor    = 0x07
	ReqGetConfiguration = 0x08
	ReqSetConfiguration = 0x09
	ReqGetInterface     = 0x0A
	ReqSetInterface     = 0x0B
	ReqSynchFrame       = 0x0C
	ReqSetSel           = 0x30
	ReqSetIsochDelay    = 0x31
)

// Suspend options
const (
	OptionSuspendNormalState        = 0b00
	OptionSuspendLowPower           = 0b01
	OptionSuspendRemoteWakeDisabled = 0b00
	OptionSuspendRemoteWakeEnabled  = 0b10
)

type Feature uint16

const (
	FeatureEndpointHalt             = Feature(0)
	FeatureInterfaceFunctionSuspend = Feature(0)
	FeatureDeviceRemoteWakeUp       = Feature(1)
	FeatureDeviceTestMode           = Feature(2)
	FeatureDeviceBHnpEnable         = Feature(3)
	FeatureDeviceAHnpSupport        = Feature(4)
	FeatureDeviceAAltHnpSupport     = Feature(5)
	FeatureDeviceWUSB               = Feature(6)
	FeatureDeviceU1Enable           = Feature(48)
	FeatureDeviceU2Enable           = Feature(49)
	FeatureDeviceLTMEnable          = Feature(50)
	FeatureDeviceB3NtfHostRel       = Feature(51)
	FeatureDeviceB3RspEnable        = Feature(52)
	FeatureDeviceLDMEnable          = Feature(53)
)

// ClearFeature request is used to clear or disable a specific feature.
//
// Feature selector values must be appropriate to the recipient. Only device feature selector values
// may be used when the recipient is a device, only interface feature selector values may be used when the
// recipient is an interface, and only endpoint feature selector values may be used when the recipient is an
// endpoint.
// See USB documentation Table 9-7 for appropriate values.
// A ClearFeature request that references a feature that cannot be cleared, that does not
// exist, or that references an interface or an endpoint that does not exist, will cause the device
// to respond with a Request Error.
//
//  Default state:
//     Device behavior when this request is received while the device is in the
//     Default state is not specified.
//  Address state:
//     This request is valid when the device is in the Address state; references
//     to interfaces, or to endpoints other than the Default Control Pipe, shall
//     cause the device to respond with a Request Error.
//  Configured state:
//     This request is valid when the device is in the Configured state.
func (d *Device) ClearFeature(recipient RequestType, feature Feature, idx uint8) error {
	_, err := d.Ctrl(RequestDirectionOut|RequestTypeStandard|recipient,
		ReqClearFeature, uint16(feature), uint16(idx), nil)
	return err
}

// GetConfiguration returns the current device configuration value.
//
// if returned value is zero, the device is not configured.
//
//  Default state:
//     Device behavior when this request is received while the device is in the
//     Default state is not specified.
//  Address state:
//     The value zero shall be returned.
//  Configured state:
//     The non-zero bConfigurationValue of the current configuration shall be returned.
func (d *Device) GetConfiguration() (int, error) {
	buff := make([]byte, 1)
	_, err := d.Ctrl(RequestDirectionIn|RequestTypeStandard|RequestRecipientDevice,
		ReqGetConfiguration, 0, 0, buff)
	return int(buff[0]), err
}

// GetDescriptor returns the specified device descriptor if the descriptor exists.
//
// The descriptor index is used to select a specific descriptor (only for configuration and string descriptors)
// when several descriptors of the same type are implemented in a device.
// eg. a device can implement several configuration descriptors.
// For other standard descriptors that can be retrieved via GetDescriptor request,
// a descriptor index of zero shall be used.
//
// LanguageID specifies which language to get for string descriptors or should be set to zero for other descriptors.
// A LanguageID of 0 will return supported languages when requesting StringDescriptor.
//
// The standard request to a device supports four types of descriptors: DeviceDescriptor, ConfigurationDescriptor,
// BOSDescriptor, StringDescriptor.
//
// As noted in USB documentation section 9.2.6.6, a device operating at Gen X speed reports the other speeds it supports
// via the BOS descriptor and shall not support the device_qualifier and other_speed_configuration descriptors.
//
// A request for a configuration descriptor returns the ConfigurationDescriptor and all InterfaceDescriptor,
// EndpointDescriptor, and EndpointCompanionDescriptors (when operating at Gen X speed) for all of the interfaces in a
// single request. The first InterfaceDescriptor follows the ConfigurationDescriptor. The EndpointDescriptor(s) for the
// first interface follow the first InterfaceDescriptor. In addition, Enhanced SuperSpeed devices shall return
// EndpointCompanionDescriptor(s) for each of the endpoints in that interface to return the endpoint capabilities
// required for Enhanced SuperSpeed devices, which would not fit inside the existing EndpointDescriptor footprint.
// If there are additional interfaces, their InterfaceDescriptor, EndpointDescriptor(s) and
// EndpointCompanionDescriptor(s) (when operating at Gen X speed) follow the first interface's
// endpoint and endpoint companion (when operating at Gen X speed) descriptors.
//
// The BOSDescriptor defines a root descriptor that is similar to the configuration descriptor,
// and is the base descriptor for accessing a family of related descriptors.
// A host can read a BOSDescriptor and learn from the TotalLength field the entire size of device-level descriptor set,
// or it can read in the entire BOSDescriptor set of device capabilities.
// The Entire set can only be accessed via reading the BOSDescriptor with a GetDescriptor request.
//
// Class-specific and/or vendor-specific descriptors follow the standard descriptors they extend or modify.
//
// All devices shall provide a DeviceDescriptor and at least one ConfigurationDescriptor.
//
// If a device does not support a requested descriptor, it responds with a request error.
//
//  Default state:
//    This is a valid request when the device is in the default state.
//  Address state:
//    This is a valid request when the device is in the address state.
//  Configured state:
//    This is a valid request when the device is in the configured state.
func (d *Device) GetDescriptor(descriptorType DescriptorType, idx uint8, languageID uint16) ([]byte, error) {
	buff := make([]byte, 256)
	n, err := d.Ctrl(RequestDirectionIn|RequestTypeStandard|RequestRecipientDevice,
		ReqGetDescriptor, (uint16(descriptorType)<<8)|uint16(idx), languageID, buff)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	return buff[0:n], nil
}

// GetInterface returns the selected alternate setting for the specified interface.
//
// Some devices have configurations with interfaces that have mutually exclusive settings.
// This request allows the host to determine the currently selected alternate setting.
//
//  Default state:
//    Device behaviour when this request is received while the device is in the default state is not specified.
//  Address state:
//    A Request error response is given by the device.
//  Configured state:
//    This is a valid request when the device is in the configured state.
func (d *Device) GetInterface(interfaceIndex uint8) (uint8, error) {
	data := make([]byte, 1)
	_, err := d.Ctrl(RequestDirectionIn|RequestTypeStandard|RequestRecipientInterface,
		ReqGetInterface, 0, uint16(interfaceIndex), data)
	return data[0], err
}

// DeviceStatus ref. USB documentation figure 9-4.
type DeviceStatus struct {
	// The LTMEnable field indicates whether the device is currently enabled to send Latency Tolerance Messages.
	// The LTMEnable field can be modified by the SetFeature() and ClearF eature() requests using the
	// FeatureDeviceLTMEnable feature.
	// This field is reset to zero when the device is reset.
	LTMEnable bool

	// The U2Enable field indicates whether the device is currently enabled to initiate U2 entry.
	// The U2Enable field can be modified by the SetFeature() and ClearFeature() requests using the
	// FeatureDeviceU2Enable feature.
	// This field is reset to zero when the device is reset.
	U2Enable bool

	// The U1Enable field indicates whether the device is currently enabled to init iate U1 entry.
	// The U1Enable field can be modified by the SetFeature() and ClearFeature() requests using the
	// FeatureDeviceU1Enable feature.
	// This field is reset to zero when the device is reset.
	U1Enable bool

	// The RemoteWakeup field is reserved and must be set to zero by Enhanced SuperSpeed
	// devices. Enhanced SuperSpeed devices use the FunctionRemoteWake enable/disable field
	// to indicate whether they are enabled for Remote Wake.
	// TODO: Field must be zero for SS devices but is used to indicate whether device is enabled for remote wake??
	RemoteWakeup bool

	// The SelfPowered field indicates whether the device is currently self-powered.
	// False = bus-powered.
	SelfPowered bool
}

// InterfaceStatus ref. USB documentation figure 9-5
// The FunctionRemoteWakeup field can be modified by the SetFeature() requests
// using the FeatureInterfaceFunctionSuspend feature.
type InterfaceStatus struct {
	// The FunctionRemoteWakeup field indicates whether the function is currently enabled to request remote wakeup.
	FunctionRemoteWakeup bool

	// The FunctionRemoteWakeCapable field indicates whether the function supports remote wake up
	FunctionRemoteWakeCapable bool
}

// EndpointStatus ref. USB documentation figure 9-6
type EndpointStatus struct {
	// The Halt feature is required to be implemented for all interrupt and bulk endpoint types.
	// The Halt feature may optionally be set with the SetFeature() request using the
	// FeatureEndpointHalt feature.
	// The Halt feature is reset to zero after either a SetConfiguration() or SetInterface() request even if the
	// requested configuration or interface is the same as the current configuration or interface.
	//
	// Enhanced SuperSpeed devices do not support functional stall on control endpoints and
	// hence do not require the Halt feature be implemented for any control endpoints.
	Halt bool
}

// DevicePTMStatus ref. USB documentation figure 9-7
type DevicePTMStatus struct {
	// LDMValid field indicates whether the LDM Link Delay is valid.
	// LDMValid shall be false if LDM Enabled is false.
	LDMValid bool

	// LDMEnabled flag indicates whether the device is currently enabled to participate in
	// Precision Time Measurement (PTM).
	// If LDM Enabled flag is set to false, the device is disabled from executing the LDM protocol and providing
	// a local bus interval boundary reference, otherwise; it is enabled to execute the LDM protocol.
	//
	// The LDMEnabled flag can be modified by the SetFeature() and ClearFeature() requests using the
	// FeatureDeviceLDMEnable feature.
	//
	// This field shall be set to one when the device is reset, allowing a PTM capable device to automatically
	// attempt to participate in LDM with its upstream partner. If a Requestor is unable to successfully establish
	// LDM Timestamp Exchanges in its Responder, then the LDMEnabled field shall be cleared to zero.
	LDMEnabled bool

	// LDMLinkDelay field is in tIsochTimestampGranularity units.
	// If LDM Valid is true, then the LDM Link Delay field defines the link delay value measured by the PTM LDM mechanism.
	// If LDM Valid is false, then the LDM Link Delay field shall be set to zero.
	LDMLinkDelay uint16
}

// GetDeviceStatus returns current device status.
func (d *Device) GetDeviceStatus() (*DeviceStatus, error) {
	data := make([]byte, 2)
	_, err := d.Ctrl(RequestDirectionIn|RequestTypeStandard|RequestRecipientDevice,
		ReqGetStatus, uint16(StatusStandard), 0, data)
	if err != nil {
		return nil, err
	}
	res := &DeviceStatus{
		LTMEnable:    (data[0] & (1 << 4)) > 0,
		U2Enable:     (data[0] & (1 << 3)) > 0,
		U1Enable:     (data[0] & (1 << 2)) > 0,
		RemoteWakeup: (data[0] & (1 << 1)) > 0,
		SelfPowered:  (data[0] & (1 << 0)) > 0,
	}
	return res, err
}

// GetDevicePTMStatus returns current PTM status.
// This is a ReqGetStatus with StatusPTM as type.
func (d *Device) GetDevicePTMStatus() (*DevicePTMStatus, error) {
	data := make([]byte, 4)
	_, err := d.Ctrl(RequestDirectionIn|RequestTypeStandard|RequestRecipientDevice,
		ReqGetStatus, uint16(StatusPTM), 0, data)
	if err != nil {
		return nil, err
	}
	res := &DevicePTMStatus{
		LDMValid:     (data[0] & (1 << 0)) > 0,
		LDMEnabled:   (data[0] & (1 << 0)) > 0,
		LDMLinkDelay: (uint16(data[3]) & 0xFF << 8) | uint16(data[2])&0xFF,
	}
	return res, nil
}

// GetInterfaceStatus returns current status for specified interface.
func (d *Device) GetInterfaceStatus(interfaceIndex uint8) (*InterfaceStatus, error) {
	data := make([]byte, 2)
	_, err := d.Ctrl(RequestDirectionIn|RequestTypeStandard|RequestRecipientInterface,
		ReqGetStatus, uint16(StatusStandard), uint16(interfaceIndex), data)
	if err != nil {
		return nil, err
	}
	res := &InterfaceStatus{
		FunctionRemoteWakeup:      (data[0] & (1 << 1)) > 0,
		FunctionRemoteWakeCapable: (data[0] & (1 << 0)) > 0,
	}
	return res, err
}

// GetEndpointStatus returns current status for specified endpoint.
func (d *Device) GetEndpointStatus(endpoint uint8) (*EndpointStatus, error) {
	data := make([]byte, 2)
	_, err := d.Ctrl(RequestDirectionIn|RequestTypeStandard|RequestRecipientEndpoint,
		ReqGetStatus, uint16(StatusStandard), uint16(endpoint), data)
	if err != nil {
		return nil, err
	}
	res := &EndpointStatus{
		Halt: (data[0] & (1 << 0)) > 0,
	}
	return res, err
}

// SetAddress request sets the device address for all future device accesses.
//
// The Status stage after the initial Setup packet assumes the same device address as the Setup packet.
// The device does not change its device address until after the Status stage of this request is completed successfully.
// Note that this is a difference between this request and all other requests.
// For all other requests, the operation indicated shall be completed before the Status stage.
//
// If the specified device address is greater than 127, then the behavior of the device is not specified.
//
//  Default state:
//    If the address specified is non-zero, then the device shall enter the address state;
//    otherwise, the device remains in the Default state.
//  Address state:
//    If the address specified is zero, then the device shall enter the default state;
//    otherwise, the device remains in the Address state but uses the newly-specified address.
//  Configured state:
//    Device behavior when this request is received while the device is in the configured state is not specified.
func (d *Device) SetAddress(address uint16) error {
	_, err := d.Ctrl(RequestDirectionOut|RequestTypeStandard|RequestRecipientDevice,
		ReqSetAddress, address, 0, nil)
	return err
}

// SetConfiguration request sets the device configuration.
// configurationValue shall be 0 or match a configuration value from a configuration descriptor.
// If the configuration value is 0, the device is placed in its address state.
//
//  Default state:
//    Device behavior when this request is received while the device is in the default state is not specified.
//  Address state:
//    If the specified configuration value is zero, then the device remains in the Address state.
//    If the specified configuration value matches the configuration value from a configuration descriptor, then that
//    configuration is selected and the device enters the configured state.
//    Otherwise, the device responds with a Request Error.
//  Configured state:
//    If the specified configuration value is zero, then the device enters the address state.
//    If the specified configuration value matches the configuration value from a configuration descriptor, then that
//    configuration is selected and the device remains in the configured state.
//    Otherwise, the device responds with a Request Error.
func (d *Device) SetConfiguration(configurationValue int) error {
	_, err := d.Ctrl(RequestDirectionOut|RequestTypeStandard|RequestRecipientDevice,
		ReqSetConfiguration, uint16(configurationValue), 0, nil)
	return err
}

// SetDescriptor request is optional, and may be used to update existing descriptors or new descriptors may be added.
//
// The descriptor index is used to select a specific descriptor (only for configuration and string descriptors) when
// several descriptors of the same type are implemented in a device. For other standard descriptors that can be set via
// a SetDescriptor request, a descriptor index of 0 shall be used.
//
// The only allowed values for descriptorType are DeviceDescriptor, ConfigurationDescriptor and StringDescriptor types.
//
//  Default state:
//    Device behaviour when this request is received while the device is in the default state is not specified.
//  Address state:
//    If supported, this is a valid request when the device is in the address state.
//  Configured state:
//    If supported, this is a valid request when the device is in the configured state.
func (d *Device) SetDescriptor(descriptorType DescriptorType, idx uint8, languageID uint16, descriptorData []byte) error {
	_, err := d.Ctrl(RequestDirectionOut|RequestTypeStandard|RequestRecipientDevice,
		ReqSetDescriptor, (uint16(descriptorType)<<8)|uint16(idx), languageID, descriptorData)
	return err
}

// SetFeature request is used to set or enable a specific feature.
//
// Feature values shall be appropriate to the recipient.
// Only device feature values may be used when recipient is a device;
// only interface feature values may be used when the recipient is an interface;
// and only endpoint feature values may be used when the recipient is an endpoint.
//
// The FeatureInterfaceFunctionSuspend is only defined for RequestRecipientInterface. idx parameter should be 0.
// Options are one or more OptionSuspend* constants.
//
// The FeatureDeviceU1Enable / FeatureDeviceU2Enable features are only defined for
// RequestRecipientDevice and options/idx shall be set to 0.
// A device shall support the U1/U2ENABLE feature when in the Configured state only.
// System software must not enable the device to initiate U1 if the time for U1 System Exit Latency initiated by
// Host plus one Bus Interval time is greater than the minimum of the service intervals of any periodic endpoints
// in the device. In addition, system software must not enable the device to initiate U2 if the time for U2 System
// Exit Latency initiated by Host plus one Bus Interval time is greater than the minimum of the service intervals of any
// periodic endpoints in the device.
// See USB documentation Table 9-7 for appropriate values.
//
// The FeatureDeviceLTMEnable is only defined for RequestRecipientDevice and options / idx shall be set to 0.
// Setting the FeatureDeviceLTMEnable feature allows the device to send Latency Tolerance Messages.
// A device shall support the FeatureDeviceLTMEnable feature if it is in the configured state
// and supports the LTM capability.
//
// The FeatureDeviceLDMEnable feature is only defined for RequestRecipientDevice and options / idx shall be set to 0.
// Setting the FeatureDeviceLDMEnable feature allows the device to execute the LDM protocol.
// A device shall support the FeatureDeviceLDMEnable feature if it is in the Address or Configured states
// and supports the PTM capability.
//
// A SetFeature request that references a feature that cannot be set or that does not exist causes a
// Stall Transaction packet to be returned in the status stage of the request.
func (d *Device) SetFeature(recipient RequestType, feature Feature, options, idx uint8) error {
	wIndex := uint16(options)<<8 | uint16(idx)
	_, err := d.Ctrl(RequestDirectionOut|RequestTypeStandard|recipient, ReqSetFeature, uint16(feature), wIndex, nil)
	return err
}

// SetInterface request allows the host to select an alternate setting for the specified interface.
//
// Some devices have configurations with interfaces that have mutually exclusive settings.
// This request allows the host to select the desired alternate setting.
// If a device only supports a default setting for the specified interface, then a STALL Transaction Packet
// may be returned in the Status stage of the request.
//
// If the interface or the alternate setting does not exist, then the device responds with a Request Error.
//
//  Default state:
//    Device behaviour when this request is received while the device is in the default state is not specified.
//  Address state:
//    The device shall respond with a request error.
//  Configured state:
//    This is a valid request when the device is in the configured state.
func (d *Device) SetInterface(interfaceIndex uint8, setting int) error {
	_, err := d.Ctrl(RequestDirectionOut|RequestTypeStandard|RequestRecipientInterface,
		ReqSetInterface, uint16(setting), uint16(interfaceIndex), nil)
	return err
}

// SetIsochronousDelay request informs the device of the delay from the time a host transmits a packet to
// the time it is received by the device.
//
// The delay field specifies a delay from 0 to 65535 ns.
// This delay represents the time from when the host starts transmitting the first framing symbol of the
// packet to when the device receives the first framing symbol of that packet.
//
// Delay should be calculated as follows:
// `delay = (sum of wHubDelay values) + (tTPTransmissionDelay * (number of hubs + 1))`
// where a wHubDelay value is provided by the Enhanced SuperSpeed Hub Descriptor of each
// hub in the path, respectively, and tTPTransmissionDelay is defined in USB documentation Table 8-35.
// TODO: Lookup Table 8-35
//
//  Default state:
//    This is a valid request when the device is in the default state.
//  Address state:
//    This is a valid request when the device is in the address state.
//  Configured state:
//    This is a valid request when the device is in the configured state.
func (d *Device) SetIsochronousDelay(delay uint16) error {
	_, err := d.Ctrl(RequestDirectionOut|RequestTypeStandard|RequestRecipientDevice, ReqSetIsochDelay, delay, 0, nil)
	return err
}

// SetSEL request sets both the U1 and U2 System Exit Latency and the U1 or U2 exit latency for
// all the links between a device and a root port on the host.
//
// u1sel: Time in us for U1 System Exit Latency.
// u1pel: Time in us for U1 Device to Host Exit Latency.
// u2sel: Time in us for U2 System Exit Latency.
// u2pel: Time in us for U2 Device to Host Exit Latency.
//
//  Default state:
//    Device behavior when this request is received while the device is in the default state is not specified.
//  Address state:
//    This is a valid request when the device is in the address state.
//  Configured state:
//    This is a valid request when the device is in the configured state.
func (d *Device) SetSEL(u1sel, u1pel uint8, u2sel, u2pel uint16) error {
	data := []byte{u1sel, u1pel, 0, 0, 0, 0}
	binary.LittleEndian.PutUint16(data[2:], u2sel)
	binary.LittleEndian.PutUint16(data[4:], u2pel)
	_, err := d.Ctrl(RequestDirectionOut|RequestTypeStandard|RequestRecipientDevice, ReqSetSel, 0, 0, data)
	return err
}

// SyncFrame request is used to set and then report an endpoint’s synchronization frame.
//
// When an endpoint supports isochronous transfers, the endpoint may also require per-frame transfers to vary
// in size according to a specific pattern. The host and the endpoint must agree
// on which frame the repeating pattern begins.
// The number of the frame in which the pattern began is returned to the host.
//
// If an Enhanced SuperSpeed device supports the Sync Frame request, it shall internally
// synchronize itself to the zeroth microframe and have a time notion of classic frame.
// Only the frame number is used to synchronize and reported by the device endpoint (i.e., no microframe number).
// The endpoint must synchronize to the zeroth microframe.
//
// This value is only used for isochronous data transfers using implicit pattern synchronization.
//
// If the specified endpoint does not support this request, then the device will respond with a Request Error.
//
//  Default state:
//    Device behavior when this request is received while the device is in the default state is not specified.
//  Address state:
//    The device shall respond with a Request Error.
//  Configured state:
//    This is a valid request when the device is in the configured state.
func (d *Device) SyncFrame(endpoint uint8) error {
	data := []byte{0, 0}
	_, err := d.Ctrl(RequestDirectionIn|RequestTypeStandard|RequestRecipientEndpoint,
		ReqSynchFrame, 0, uint16(endpoint), data)
	return err
}
