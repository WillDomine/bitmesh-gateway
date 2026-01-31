package router

import "sync"

type Router struct {
	requestPool sync.Pool
}

type BitContext struct {
	Flags uint32
}

func NewRouter() *Router {
	return &Router{
		requestPool: sync.Pool{
			New: func() any {return new(BitContext)},
		},
	}
}

func (r *Router) Decide(flags uint32) (string, bool) {
	const PremiumUserBit = 1 << 3

	if flags & PremiumUserBit != 0 {
		return "premium-backend", true
	}
	return "standard-backend", false
}