package router

func BenchmarkDecide(b *testing.B) {
	r := NewRouter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Decide(0x08) // Test the "Premium" bit
	}
}