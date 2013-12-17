package message

type Header struct {
	// The id of the type of the message.
	TypeId byte

	// The length of the message content.
	Length int32
}

func NewHeader(typeId byte, length int32) *Header {
	return &Header{
		TypeId: typeId,
		Length: length,
	}
}
