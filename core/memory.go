package core

import "unsafe"

type Word int64
type Byte uint8

type Memory []Byte

func AllocateMemory(size Word) Memory {
	return make(Memory, size)
}

func (m Memory) ReadWord(address Word) Word {
	return *(*Word)(unsafe.Pointer(&m[address]))
}

func (m Memory) WriteWord(address Word, value Word) {
	*(*Word)(unsafe.Pointer(&m[address])) = value
}

func (m Memory) Size() Word {
	return Word(len(m))
}
