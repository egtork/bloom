package bloom

import (
	"crypto/rand"
	"fmt"
	"testing"
)

// Test over a few keys to check against false negatives and false positives
func TestBloom(t *testing.T) {
	element1 := []byte("A stormy ocean")
	element2 := []byte("A towering wave")
	element3 := []byte("A rickety boat")
	bloom := NewFNV32(1024, 4)
	if bloom.Check(element1) || bloom.Check(element2) || bloom.Check(element3) {
		t.Errorf("Element detected before any elements were added")
	}
	bloom.Add(element1)
	if !bloom.Check(element1) {
		t.Errorf("Element1 not detected after it was added")
	}
	if bloom.Check(element2) {
		t.Errorf("Element2 detected before it was added")
	}
	bloom.Add(element2)
	if !bloom.Check(element2) {
		t.Errorf("Element2 not detected after it was added")
	}
	if bloom.Check(element3) {
		t.Errorf("Element3 detected before it was added")
	}
	bloom.Add(element3)
	if !bloom.Check(element3) {
		t.Errorf("Element3 not detected after it was added")
	}
	if !bloom.Check(element1) {
		t.Errorf("Element1 not detected after element1, element2, and element3 were added")
	}
}

// Generate a random element of `keyLen` bytes
func generateKey(keyLen int) []byte {
	key := make([]byte, keyLen)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return key
}

// Test over many random keys to check against false negatives
func TestBloomRandom(t *testing.T) {
	bloom := NewFNV32(1024, 4)
	nKeys := 100
	nBytesPerKey := 10
	for k := 0; k < nKeys; k++ {
		key := generateKey(nBytesPerKey)
		bloom.Add(key)
		found := bloom.Check(key)
		if !found {
			t.Errorf("Detected a false negative")
		}
	}
}

func BenchmarkBloomAdd(b *testing.B) {
	key := generateKey(10)
	bloom := NewFNV32(1024, 4)
	for n := 0; n < b.N; n++ {
		bloom.Add(key)
	}
}

func BenchmarkBloomCheck(b *testing.B) {
	key := generateKey(10)
	bloom := NewFNV32(1024, 4)
	for n := 0; n < b.N; n++ {
		bloom.Check(key)
	}
}
