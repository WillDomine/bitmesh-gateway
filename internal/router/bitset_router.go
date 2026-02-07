package router

import "sync"

type Router struct {
	requestPool sync.Pool
}

//Pool Object
type BitContext struct {
	Flags uint32
}

//Resets struct so it can be reused safely
func (b *BitContext) Reset() {
	b.Flags = 0
}

func NewRouter() *Router {
	return &Router{
		requestPool: sync.Pool{
			New: func() any { return new(BitContext) },
		},
	}
}

func (r *Router) Decide(flags uint32) (string, bool) {
	//Acquire generic object and type assert it
	ctx := r.requestPool.Get().(*BitContext)

	//Rewrite old data
	ctx.Reset()
	ctx.Flags = flags

	//Returns context to pool to reduce GC pressure
	defer r.requestPool.Put(ctx)

	// Bitwise for premium
	const PremiumUserBit = 1 << 3 // 0x08

	if flags&PremiumUserBit != 0 {
		return "premium-backend", true
	}
	return "standard-backend", true
}
