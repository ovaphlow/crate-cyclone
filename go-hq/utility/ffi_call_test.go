package utility

import (
	"testing"
	"time"
)

// TestCallGenerateKsuid 测试 CallGenerateKsuid 函数
func TestCallGenerateKsuid(t *testing.T) {
	ksuid := CallGenerateKsuid()
	if len(ksuid) == 0 {
		t.Errorf("Expected non-empty KSUID, got empty string")
	}
	println(time.Now().Format(time.RFC3339), "Generated KSUID(ffi call):", ksuid)
}
