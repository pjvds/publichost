package proto

type Header struct {
	// The id of the type of the message.
	TypeId byte

	// The id of the tunnel this message is writen to or is send from.
	TunnelId TunnelId

	// The id of the request itself, or when the message
	// is a reply this id indicates to which request the
	// reply belongs.
	RequestId int16

	// The length of the message content.
	Length int32
}

func NewHeader(typeId byte, tunnelId int16, requestId int16, length int32) *Header {
	return &Header{
		TypeId:    typeId,
		TunnelId:  TunnelId(tunnelId),
		RequestId: requestId,
		Length:    length,
	}
}
