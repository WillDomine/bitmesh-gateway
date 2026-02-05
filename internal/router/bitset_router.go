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
	//Acquire generic object and type assert it
	ctx := r.requestPool.Get().(*BitContext)

	//Rewrite old data
	ctx.Reset()
	ctx.Flags = flags

	//Object goes back to pool even if panic
	defer r.requestPool.Put(ctx)

	// Bitwise for premium
	const PremiumUserBit = 1 << 3 // 0x08

	if flags & PremiumUserBit != 0 {
		return "premium-backend", true
	}
	return "standard-backend", false
}