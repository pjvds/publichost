package proto

import (
	"strconv"
)

type TunnelId int16

func (id TunnelId) String() string {
	return strconv.Itoa(int(id))
}
