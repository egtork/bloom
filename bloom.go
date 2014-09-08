// A Bloom filter implementation in Go

package bloom

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
	"log"
	"math"
)

type Bloom struct {
	size      uint32
	k         int
	hashFuncs []hash.Hash
	hashes    []uint
	bits      bitArray
}

// New creates a new Bloom filter of the specified size, using k hash functions derived from
// two distinct hash functions specified in the []hash.Hash slice.
func New(size uint32, k int, hashFuncs []hash.Hash) *Bloom {
	if len(hashFuncs) != 2 {
		log.Fatal("Supply a slice of two distinct hash functions")
	}
	return &Bloom{
		bits:      newBitArray(size),
		hashFuncs: hashFuncs,
		k:         k,
		size:      size,
		hashes:    make([]uint, k),
	}
}

// NewFnv32 creates a new Bloom filter of the specified size, using k hash functions derived from
// 32-bit FNV-1 and FNV-1a hash functions
func NewFnv32(size uint32, k int) *Bloom {
	hashFuncs := []hash.Hash{
		fnv.New64(),
		fnv.New64a(),
	}
	return &Bloom{
		bits:      newBitArray(size),
		hashFuncs: hashFuncs,
		k:         k,
		size:      size,
		hashes:    make([]uint, k),
	}
}

// Add the input element to the set.
func (blm *Bloom) Add(input []byte) {
	blm.doubleHash(input)
	for i := 0; i < blm.k; i++ {
		blm.bits.setBit(blm.hashes[i])
	}
}

// Check whether the input element has been added to the set. If the input is present,
// Check returns true. If the input is not present, Check is likely to return false but may return
// true (a false positive).
func (blm *Bloom) Check(input []byte) bool {
	blm.doubleHash(input)
	for i := 0; i < blm.k; i++ {
		set := blm.bits.bit(blm.hashes[i])
		if !set {
			return false
		}
	}
	return true
}

// Reset clears the Bloom filter's bit array.
func (blm *Bloom) Reset() {
	blm.bits.reset()
}

// doubleHash computes k hashes from the Bloom filter's two distinct hash functions.
func (blm *Bloom) doubleHash(input []byte) {
	x1 := hashToUint32(blm.hashFuncs[0], input)
	x2 := hashToUint32(blm.hashFuncs[1], input)
	for i := 0; i < blm.k; i++ {
		blm.hashes[i] = uint((x1 + uint32(i)*x2) % blm.size)
	}
}

// hashToUint32 computes the hash of the specified byte slice and returns it as a uint32.
func hashToUint32(h hash.Hash, input []byte) uint32 {
	h.Reset()
	h.Write(input)
	b := h.Sum(nil)
	return binary.BigEndian.Uint32(b)
}

// bitArray is a lightweight bit array that's slightly faster than big.Int in this context.
type bitArray []byte

func newBitArray(size uint32) bitArray {
	if !powerOfTwo(size) || size < 8 {
		log.Fatal("Bit array size must be a power of 2 and larger than 8")
	}
	return make(bitArray, size/8)
}

// setBit sets the nth bit to 1
func (b bitArray) setBit(n uint) {
	b[n>>3] |= 1 << (n & 7)
}

var mask = [8]byte{1, 2, 4, 8, 16, 32, 64, 128}

// bit returns true if the nth bit is 1
func (b bitArray) bit(n uint) bool {
	return (b[n>>3] & mask[n&7]) != 0
}

// reset sets all bits to 0
func (b bitArray) reset() {
	for i := range b {
		b[i] = 0
	}
}

// FalsePositiveRate computes the expected rate of false-positives given a Bloom filter of size m,
// k hash functions, and n set elements.
func FalsePositiveRate(m, k, n float64) float64 {
	var t, f float64
	t = math.Pow(1-1/m, k*n)
	f = math.Pow(1-t, k)
	return f
}

func FalsePositiveRateApprox(m, k, n float64) float64 {
	var t, f float64
	t = math.Exp(-k * n / m)
	f = math.Pow(1-t, k)
	return f
}

func powerOfTwo(n uint32) bool {
	return (n & (n - 1)) == 0
}
