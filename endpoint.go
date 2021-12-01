package usb

type (
	TransferType        uint8
	SynchronizationType uint8
	UsageType           uint8
)

const (
	TransferTypeControl = TransferType(iota)
	TransferTypeIsochronous
	TransferTypeBulk
	TransferTypeInterrupt
)

const (
	SynchronizationTypeNoSync = SynchronizationType(iota)
	SynchronizationTypeAsynchronous
	SynchronizationTypeAdaptive
	SynchronizationTypeSynchronous
)

const (
	UsageTypeData = UsageType(iota)
	UsageTypeFeedback
	UsageTypeExplicitFeedbackData
	UsageTypeReserved
)

const (
	EndpointDirectionIn  = 0x80
	EndpointDirectionOut = 0x00
)

func (ep *EndpointDescriptor) TransferType() TransferType {
	return TransferType(ep.BmAttributes & 0b00000011)
}

func (ep *EndpointDescriptor) SynchronizationType() SynchronizationType {
	return SynchronizationType((ep.BmAttributes & 0b00001100) >> 2)
}

func (ep *EndpointDescriptor) UsageType() UsageType {
	return UsageType((ep.BmAttributes & 0b00110000) >> 4)
}
