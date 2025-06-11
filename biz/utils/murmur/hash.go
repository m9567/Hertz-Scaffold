package murmur

import (
	"encoding/binary"
	"math/bits"
)

var (
	c1_32 = uint32(0xcc9e2d51)
	c2_32 = uint32(0x1b873593)
	c1_64 = uint64(0x87c37b91114253d5)
	c2_64 = uint64(0x4cf5ad432745937f)
)

func fmix32(h uint32) uint32 {
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

func Hash32(data []byte, seed uint32) uint32 {
	h := seed
	nblocks := len(data) / 4
	for i := 0; i < nblocks; i++ {
		k := binary.LittleEndian.Uint32(data[i*4:])
		k *= c1_32
		k = bits.RotateLeft32(k, 15)
		k *= c2_32

		h ^= k
		h = bits.RotateLeft32(h, 13)
		h = h*5 + 0xe6546b64
	}

	tail := data[nblocks*4:]
	var k1 uint32
	switch len(tail) {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= c1_32
		k1 = bits.RotateLeft32(k1, 15)
		k1 *= c2_32
		h ^= k1
	}

	h ^= uint32(len(data))
	h = fmix32(h)
	return h
}

func Hash64(data []byte, seed uint32) uint64 {
	h := uint64(seed)
	nblocks := len(data) / 8
	for i := 0; i < nblocks; i++ {
		k := binary.LittleEndian.Uint64(data[i*8:])
		k *= c1_64
		k = bits.RotateLeft64(k, 31)
		k *= c2_64
		h ^= k
		h = bits.RotateLeft64(h, 27)*5 + 0x52dce729
	}
	tail := data[nblocks*8:]
	var k1 uint64
	switch len(tail) {
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0])
		k1 *= c1_64
		k1 = bits.RotateLeft64(k1, 31)
		k1 *= c2_64
		h ^= k1
	}

	h ^= uint64(len(data))
	return fmix64(h)
}

func Hash128(data []byte, seed uint32) (uint64, uint64) {
	h1 := uint64(seed)
	h2 := uint64(seed)
	nblocks := len(data) / 16

	for i := 0; i < nblocks; i++ {
		offset := i * 16
		k1 := binary.LittleEndian.Uint64(data[offset:])
		k2 := binary.LittleEndian.Uint64(data[offset+8:])

		k1 *= c1_64
		k1 = bits.RotateLeft64(k1, 31)
		k1 *= c2_64
		h1 ^= k1

		h1 = bits.RotateLeft64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2_64
		k2 = bits.RotateLeft64(k2, 33)
		k2 *= c1_64
		h2 ^= k2

		h2 = bits.RotateLeft64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}
	tail := data[nblocks*16:]
	var k1, k2 uint64
	switch len(tail) {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8])
		k2 *= c2_64
		k2 = bits.RotateLeft64(k2, 33)
		k2 *= c1_64
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0])
		k1 *= c1_64
		k1 = bits.RotateLeft64(k1, 31)
		k1 *= c2_64
		h1 ^= k1
	}

	h1 ^= uint64(len(data))
	h2 ^= uint64(len(data))

	h1 += h2
	h2 += h1

	h1 = fmix64(h1)
	h2 = fmix64(h2)

	h1 += h2
	h2 += h1

	return h1, h2
}
