package main

type BitVector struct {
	vector []uint32
}

func NewBitVector(numBits uint32) BitVector {
	n := numBits / 32
	if numBits%32 > 0 {
		n++
	}

	return BitVector{make([]uint32, n)}
}

func (bv *BitVector) SetBit2(docId uint32) {
	bytepos := docId / 32
	bitpos := docId % 32
	bv.vector[bytepos] = bv.vector[bytepos] | 1<<uint8(bitpos)
}

func (bv *BitVector) SetBit(docId uint32) {
	bv.vector[docId/32] |= 1 << uint8(docId%32)
}

func (bv *BitVector) GetBit2(docId uint32) bool {
	bytepos := docId / 32
	bitpos := docId % 32

	if (bv.vector[bytepos] & (1 << uint8(bitpos))) > 0 {
		return true
	} else {
		return false
	}
}

func (bv *BitVector) GetBit(docId uint32) bool {
	return (bv.vector[docId/32] & (1 << uint8(docId%32))) > 0
}

func (bv *BitVector) ClearBit2(docId uint32) {
	bytepos := docId / 32
	bitpos := docId % 32
	bv.vector[bytepos] = bv.vector[bytepos] & ^(1 << uint8(bitpos))
}

func (bv *BitVector) ClearBit(docId uint32) {
	bv.vector[docId/32] &= ^(1 << uint8(docId%32))
}

func (bv *BitVector) Not() {
	for i := 0; i < len(bv.vector); i++ {
		bv.vector[i] = ^bv.vector[i]
	}
}

func (bv *BitVector) Intersect(other *BitVector) {
	if len(bv.vector) != len(other.vector) {
		panic("cannot intersect BitVectors of different sizes")
	}

	for i := 0; i < len(bv.vector); i++ {
		bv.vector[i] &= other.vector[i]
	}
}

func (bv *BitVector) Union(other *BitVector) {
	if len(bv.vector) != len(other.vector) {
		panic("cannot intersect BitVectors of different sizes")
	}

	for i := 0; i < len(bv.vector); i++ {
		bv.vector[i] |= other.vector[i]
	}
}

func (bv *BitVector) GetCount() int {
	count := 0
	for i := 0; i < len(bv.vector)*32; i++ {
		if bv.GetBit(uint32(i)) {
			count++
		}
	}
	return count
}

func (bv *BitVector) NumberOfSetBits() uint32 {
	var count uint32 = 0
	for i := 0; i < len(bv.vector); i++ {
		n := bv.vector[i]
		n = n - ((n >> 1) & 0x55555555)
		n = (n & 0x33333333) + ((n >> 2) & 0x33333333)
		count += (((n + (n >> 4)) & 0x0F0F0F0F) * 0x01010101) >> 24
	}

	return count
}

func (bv *BitVector) Vector() []uint32 {
	return bv.vector
}
