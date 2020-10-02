package benchmark

import (
	"github.com/sonh/qs"
	"testing"
)

type Primitive struct {
	String  string
	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
}

func BenchmarkEncodePrimitive(b *testing.B) {
	encoder := qs.NewEncoder()
	s := Primitive{
		String: "abc",
		Bool:   true,
		Int:    12,
		Int8:   int8(8),
		Int16:  int16(16),
		Int32:  int32(32),
		Int64:  int64(64),
		Uint:   24,
		Uint8:  uint8(8),
		Uint16: uint16(16),
		Uint32: uint32(32),
		Uint64: uint64(64),
	}
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := encoder.Values(&s); err != nil {
			b.Error(err)
		}
	}
}
