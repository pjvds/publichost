package message

const (
	TypeExposeRequest    = iota
	TypeExposeReponse    = iota
	TypeProxyDataRequest = iota
	TypeNokResponse      = iota
)

type Header struct {
	// The id of the type of the message.
	TypeId byte

	// The length of the message content.
	Length int32
}

type ExposeRequest struct {
	LocalAddress string
}

type ExposeResponse struct {
	RouteId       int32
	RemoteAddress string
}

type NokResponse struct {
	Error string
}
