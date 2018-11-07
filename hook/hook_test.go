package hook

import (
	"testing"
	"time"
)

func BenchmarkSha256_Sha256Sum(b *testing.B) {

	b.ReportAllocs()

	s := Sha256{}
	now := time.Now().Unix()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_ = s.Sha256Sum("foobarbaz", now)
	}
}
