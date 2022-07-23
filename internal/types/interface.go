package types

import "sync/atomic"

type Interface struct {
	value atomic.Value
}

func (v *Interface) Val() interface{} {
	return v.value.Load()
}
