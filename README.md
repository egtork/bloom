bloom
=====
A Bloom filter written in Go.

## Initialization

Use `New` or `NewFnv32` to initialize a new Bloom filter.

`New` requires that you supply the hash functions used by the Bloom filter, whereas `NewFnv32` uses the FNV-1 and FNV-1a hash functions from the Go standard library.

#### func New
`func New(size uint32, k int, hashFuncs []hash.Hash) *Bloom`

where `size` is the number of bits in the bit array,

`k` (>=2) is the number of hashes generated for each element,

and `hashFuncs` is a slice containing two distinct hash functions satisfying the [Hash interface](http://golang.org/pkg/hash/). The two hash functions are used to create `k` unique hashes using the double hashing technique described in *"Less Hashing, Same Performance: Building a Better Bloom Filter"* by A. Kirsch and M. Mitzenmacher.

#### func NewFnv32(size uint32, k int) *Bloom
`func NewFnv32(size uint32, k int) *Bloom`

Parameters `size` and `k` are the same as above. `hashFuncs` is set internally to use the 32-bit FNV-1 and FNV-1a hash functions from the Go standard library.

## Usage

#### func Add
`func (blm *Bloom) Add(element []byte)`

Add an element to the set.

#### func Check
`func (blm *Bloom) Check(element []byte) bool`

Check whether the element has been added to the set. If the element is present, Check returns true. If the element is not present, Check is likely to return false but may return true (a false positive).

#### func Reset
`func (blm *Bloom) Reset()`

Reset clears the Bloom filter's bit array.

## Example

    bloom := NewFnv32(1024, 4)
    key := []byte("Sample element")
    bloom.Add(key)
    exists := bloom.Check(key)