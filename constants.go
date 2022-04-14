package usb

type StatusType uint8

const (
	StatusStandard = StatusType(0x00)
	StatusPTM      = StatusType(0x01)
)

type RequestType uint8

const (
	RequestDirectionIn  = RequestType(0b10000000)
	RequestDirectionOut = RequestType(0b00000000)

	RequestTypeStandard = RequestType(0b00000000)
	RequestTypeClass    = RequestType(0b00100000)
	RequestTypeVendor   = RequestType(0b01000000)
	RequestTypeReserved = RequestType(0b01100000)

	RequestRecipientDevice    = RequestType(0b00000000)
	RequestRecipientInterface = RequestType(0b00000001)
	RequestRecipientEndpoint  = RequestType(0b00000010)
	RequestRecipientOther     = RequestType(0b00000011)
	// From Wireless USB 1.0
	RequestRecipientPort  = RequestType(0x04)
	RequestRecipientRPipe = RequestType(0x05)
)
