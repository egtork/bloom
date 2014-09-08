package bloom

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func generateKey(keyLen int) []byte {
	key := make([]byte, keyLen)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return key
}

// Test to confirm that there are no false negatives
func TestBloom(t *testing.T) {
	bloom := NewFnv32(1024, 4)
	nKeys := 100
	nBytesPerKey := 10
	for k := 0; k < nKeys; k++ {
		key := generateKey(nBytesPerKey)
		bloom.Add(key)
		found := bloom.Check(key)
		if !found {
			t.Errorf("Key missing")
		}
	}
}

func BenchmarkBloomAdd(b *testing.B) {
	key := generateKey(10)
	bloom := NewFnv32(1024, 4)
	for n := 0; n < b.N; n++ {
		bloom.Add(key)
	}
}

func BenchmarkBloomCheck(b *testing.B) {
	key := generateKey(10)
	bloom := NewFnv32(1024, 4)
	for n := 0; n < b.N; n++ {
		bloom.Check(key)
	}
}
