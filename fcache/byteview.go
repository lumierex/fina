package fcache

// ByteView implement Value and shouldn't have the pointer receive
// TODO explain
type ByteView struct {
	b []byte
}

// Len Cache Value required
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice avoid ByteView's inner data was revise by other program
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

// cloneBytes
func cloneBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
