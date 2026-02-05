package router

import "testing"

var GlobalSink *BitContext

type SimpleRouter struct{}

//Simple Allocating Version of Router, shows the cost of allocation and GC
func (r *SimpleRouter) Decide(flags uint32) (string, bool) {
	ctx := &BitContext{Flags: flags}

	//Force allocate on heap instead of stack
	GlobalSink = ctx

	const PremiumUserBit = 1 << 3
	if ctx.Flags&PremiumUserBit != 0 {
		return "premium-backend", true
	}
	return "standard-backend", true
}

//Slow Standard Allocation, no pooling.
func BenchmarkDecide_Standard_Allocation(b *testing.B) {
	r := &SimpleRouter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Decide(0x08)
	}
}

//Fast uses pool
func BenchmarkDecide_Pooled_ZeroAlloc(b *testing.B) {
	r := NewRouter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Decide(0x08)
	}
}
