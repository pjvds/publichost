package proto

import (
	"io"
)

type ProxyDataRequest struct {
	RouteId int32

	Flags ProxyDataFlags

	DataLength int32
	Data       io.Reader
}
