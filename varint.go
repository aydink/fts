package main

import (
	"encoding/binary"
	"errors"
)

type VIntArray struct {
	array         []byte
	cursor_writer int
	cursor_reader int
	buf           []byte
}

func NewVIntArray() VIntArray {
	return VIntArray{make([]byte, 0, 4096), 0, 0, make([]byte, 9, 9)}
}

func NewVIntArrayFromBytes(b []byte) VIntArray {
	return VIntArray{b, len(b), 0, make([]byte, 9, 9)}
}

func (va *VIntArray) Add(n uint64) {
	i := binary.PutUvarint(va.buf, n)
	va.cursor_writer += i

	va.array = append(va.array, va.buf[0:i]...)
}

func (va *VIntArray) Read() (uint64, error) {
	n, i := binary.Uvarint(va.array[va.cursor_reader:])
	if i > 0 {
		va.cursor_reader += i
		return n, nil
	}

	return n, errors.New("either buffer is too small or value overflows uint64")
}

func (va *VIntArray) Bytes() []byte {
	return va.array
}
