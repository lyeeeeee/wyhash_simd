package main

import (
	"fmt"
	"io/ioutil"
	"math/bits"
	"time"
	"unsafe"
)

const (
	DefaultSeed = 0xa0761d6478bd642f // s0
	s1          = 0xe7037ed1a0b428db
	s2          = 0x8ebc6af09c88c6e3
	s3          = 0x589965cc75374cc3
	s4          = 0x1d8e4e27c47d124f
)

var seeds = [...]byte{
	0xe7, 0x03, 0x7e, 0xd1, 0xa0, 0xb4, 0x28, 0xdb,
	0x8e, 0xbc, 0x6a, 0xf0, 0x9c, 0x88, 0xc6, 0xe3,
	0x58, 0x99, 0x65, 0xcc, 0x75, 0x37, 0x4c, 0xc3,
	0x1d, 0x8e, 0x4e, 0x27, 0xc4, 0x7d, 0x12, 0x4f,
}

func Read8(p unsafe.Pointer) uint64 {
	q := (*[8]byte)(p)
	return uint64(q[0]) | uint64(q[1])<<8 | uint64(q[2])<<16 | uint64(q[3])<<24 | uint64(q[4])<<32 | uint64(q[5])<<40 | uint64(q[6])<<48 | uint64(q[7])<<56
}

func Read4(p unsafe.Pointer) uint64 {
	q := (*[4]byte)(p)
	return uint64(uint32(q[0]) | uint32(q[1])<<8 | uint32(q[2])<<16 | uint32(q[3])<<24)
}

func _wymix(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

func Sum641(data []byte, seed uint64) uint64 {
	var (
		a uint64
		b uint64
	)

	length := len(data)
	i := uintptr(len(data))
	paddr := *(*unsafe.Pointer)(unsafe.Pointer(&data))

	if i > 64 {
		var see1 = seed
		for i > 64 {
			seed = _wymix(Read8(paddr)^s1, Read8(add(paddr, 8))^seed) ^ _wymix(Read8(add(paddr, 16))^s2, Read8(add(paddr, 24))^seed)
			see1 = _wymix(Read8(add(paddr, 32))^s3, Read8(add(paddr, 40))^see1) ^ _wymix(Read8(add(paddr, 48))^s4, Read8(add(paddr, 56))^see1)
			paddr = add(paddr, 64)
			i -= 64
		}
		seed ^= see1
	}

	for i > 16 {
		seed = _wymix(Read8(paddr)^s1, Read8(add(paddr, 8))^seed)
		paddr = add(paddr, 16)
		i -= 16
	}

	// i <= 16
	switch {
	case i == 0:
		return _wymix(s1, _wymix(s1, seed))
	case i < 4:
		a = uint64(*(*byte)(paddr))<<16 | uint64(*(*byte)(add(paddr, uintptr(i>>1))))<<8 | uint64(*(*byte)(add(paddr, uintptr(i-1))))
		// b = 0
		return _wymix(s1^uint64(length), _wymix(a^s1, seed))
	case i == 4:
		a = Read4(paddr)
		// b = 0
		return _wymix(s1^uint64(length), _wymix(a^s1, seed))
	case i < 8:
		a = Read4(paddr)
		b = Read4(add(paddr, i-4))
		return _wymix(s1^uint64(length), _wymix(a^s1, b^seed))
	case i == 8:
		a = Read8(paddr)
		// b = 0
		return _wymix(s1^uint64(length), _wymix(a^s1, seed))
	default: // 8 < i <= 16
		a = Read8(paddr)
		b = Read8(add(paddr, i-8))
		return _wymix(s1^uint64(length), _wymix(a^s1, b^seed))
	}
}

func Sum642(data []byte, seed uint64) uint64 {
	var (
		a uint64
		b uint64
	)

	length := len(data)
	i := uintptr(len(data))
	paddr := *(*unsafe.Pointer)(unsafe.Pointer(&data))

	if i > 64 {
		var see1 = seed
		for i > 64 {
			a := Read8(paddr) ^ s1
			b := Read8(add(paddr, 8)) ^ seed
			c := Read8(add(paddr, 16)) ^ s2
			d := Read8(add(paddr, 24)) ^ seed

			e := Read8(add(paddr, 32)) ^ s3
			f := Read8(add(paddr, 40)) ^ see1
			g := Read8(add(paddr, 48)) ^ s4
			h := Read8(add(paddr, 56)) ^ see1

			seed = _wymix(a, b) ^ _wymix(c, d)
			see1 = _wymix(e, f) ^ _wymix(g, h)
			paddr = add(paddr, 64)
			i -= 64
		}
		seed ^= see1
	}

	for i > 16 {
		seed = _wymix(Read8(paddr)^s1, Read8(add(paddr, 8))^seed)
		paddr = add(paddr, 16)
		i -= 16
	}

	// i <= 16
	switch {
	case i == 0:
		return _wymix(s1, _wymix(s1, seed))
	case i < 4:
		a = uint64(*(*byte)(paddr))<<16 | uint64(*(*byte)(add(paddr, uintptr(i>>1))))<<8 | uint64(*(*byte)(add(paddr, uintptr(i-1))))
		// b = 0
		return _wymix(s1^uint64(length), _wymix(a^s1, seed))
	case i == 4:
		a = Read4(paddr)
		// b = 0
		return _wymix(s1^uint64(length), _wymix(a^s1, seed))
	case i < 8:
		a = Read4(paddr)
		b = Read4(add(paddr, i-4))
		return _wymix(s1^uint64(length), _wymix(a^s1, b^seed))
	case i == 8:
		a = Read8(paddr)
		// b = 0
		return _wymix(s1^uint64(length), _wymix(a^s1, seed))
	default: // 8 < i <= 16
		a = Read8(paddr)
		b = Read8(add(paddr, i-8))
		return _wymix(s1^uint64(length), _wymix(a^s1, b^seed))
	}
}

func main() {
	dat, _ := ioutil.ReadFile("output.dat")
	start := time.Now()
	res0 := Sum641(dat, DefaultSeed)
	duration0 := time.Since(start)
	fmt.Println("zyh", len(dat), res0, duration0)

	start = time.Now()
	res1 := Sum642(dat, DefaultSeed)
	duration1 := time.Since(start)
	fmt.Println("l y", len(dat), res1, duration1)

}
