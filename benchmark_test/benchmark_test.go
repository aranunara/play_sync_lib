package benchmark_test

import (
	"testing"

	"play_sync_lib/chan_pattern"
	"play_sync_lib/sync_pattern"
)

// BenchmarkSyncPattern-2                 1        10001583003 ns/op
func BenchmarkSyncPattern(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sync_pattern.RunSyncPattern()
	}
}

// BenchmarkChanPattern-2                 1        2001147909 ns/op
func BenchmarkChanPattern(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chan_pattern.RunChanPattern()
	}
}
