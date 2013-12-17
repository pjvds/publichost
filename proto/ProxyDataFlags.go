package proto

type ProxyDataFlags byte

const (
	flag_Close = ProxyDataFlags(1 << iota)
)

func (p ProxyDataFlags) IsClose() bool {
	return byte(p)&byte(1<<flag_Close) != 0
}
