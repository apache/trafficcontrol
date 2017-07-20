//
// atshash.go
// ConsistentHash() is a function which hashes the given string and returns a hashed integer
//
// atsconsistenthash.go
// ATSConsistentHash is an object which uses ConsistentHash(). It allows inserting parents ("ATSConsistentHashNode") via `Insert()`, and then hashing arbitrary strings (URLs) to them via `Lookup()`.
//
// atsparentconsistenthash.go
// ATSParentConsistentHash is an object representing all the parents of a given remap rule (delivery service). It is constructed via `NewSimpleATSParentConsistentHash`, and then given HTTP requests (`RequestData`) can get the parent which they should go to (`ATSParentResult`) via `SelectParent`.
//
// atsparentresult.go
// ATSParentResult is an object containing all the information needed to redirect a request to the proper parent. Used by ATSParentConsistentHash
//
// atsorderedmap.go
// OrderedMapUint64Node is an ordered map[uint64]ATSConsistentHashNode.
// It is used by ATSConsistentHash to store hashes mapped to parents. An ordered map is necessary, because we only create a few hashes of parents, and need to get the numerically next hash from arbitrary strings, and Go doesn't include an ordered map.
//

package grove

import (
	"encoding/binary"
)

const SipBlockSize = 8

func ConsistentHash(s string) uint64 {
	h := NewATSHash64Sip24()
	h.Update([]byte(s))
	h.Final()
	return h.Get()
}

type ATSHash interface {
	Update(data []byte)
	Final()
	Get() uint64
	Clear()
}

func NewATSHash64Sip24() ATSHash {
	return &ATSHash64Sip24{}
}

type ATSHash64Sip24 struct {
	BlockBuffer    [8]byte
	BlockBufferLen uint64

	Key0 uint64
	Key1 uint64
	V0   uint64
	V1   uint64
	V2   uint64
	V3   uint64

	HFinal   uint64
	TotalLen uint64

	Finalized bool
}

func ROTL64(a, b uint64) uint64 {
	return (a << b) | (a >> (64 - b))
}

func U8To64LE(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func (s *ATSHash64Sip24) SipCompress() {
	s.V0 += s.V1
	s.V2 += s.V3
	s.V1 = ROTL64(s.V1, 13)
	s.V3 = ROTL64(s.V3, 16)
	s.V1 ^= s.V0
	s.V3 ^= s.V2
	s.V0 = ROTL64(s.V0, 32)
	s.V2 += s.V1
	s.V0 += s.V3
	s.V1 = ROTL64(s.V1, 17)
	s.V3 = ROTL64(s.V3, 21)
	s.V1 ^= s.V2
	s.V3 ^= s.V0
	s.V2 = ROTL64(s.V2, 32)
}

func (self *ATSHash64Sip24) Update(data []byte) {
	if self.Finalized {
		return
	}

	len := uint64(len(data))
	m := data
	self.TotalLen += len
	if len+self.BlockBufferLen < SipBlockSize {
		copy(self.BlockBuffer[self.BlockBufferLen:], m[:len])
		self.BlockBufferLen += len
	} else {
		blockOff := uint64(0)
		if self.BlockBufferLen > 0 {
			blockOff = SipBlockSize - self.BlockBufferLen
			copy(self.BlockBuffer[self.BlockBufferLen:], m[:blockOff])

			mi := U8To64LE(self.BlockBuffer[:])
			self.V3 ^= mi
			self.SipCompress()
			self.SipCompress()
			self.V0 ^= mi
		}

		blocks := uint64(0)
		for i, blocks := blockOff, ((len - blockOff) & ^(uint64(SipBlockSize - 1))); i < blocks; i += SipBlockSize {
			mi := U8To64LE(m[i:])
			self.V3 ^= mi
			self.SipCompress()
			self.SipCompress()
			self.V0 ^= mi
		}

		self.BlockBufferLen = (len - blockOff) & (SipBlockSize - 1)
		copy(self.BlockBuffer[:], m[blockOff+blocks:blockOff+blocks+self.BlockBufferLen])
	}
}

func (self *ATSHash64Sip24) Final() {
	if self.Finalized {
		return
	}

	last7 := uint64((self.TotalLen & 0xff) << 56)

	for i := int(self.BlockBufferLen) - 1; i >= 0; i-- {
		last7 |= uint64(self.BlockBuffer[i] << (uint(i) * 8))
	}

	self.V3 ^= last7
	self.SipCompress()
	self.SipCompress()
	self.V0 ^= last7
	self.V2 ^= 0xff
	self.SipCompress()
	self.SipCompress()
	self.SipCompress()
	self.SipCompress()
	self.HFinal = self.V0 ^ self.V1 ^ self.V2 ^ self.V3
	self.Finalized = true
}

func (self *ATSHash64Sip24) Get() uint64 {
	if self.Finalized {
		return self.HFinal
	}
	return 0
}

func (self *ATSHash64Sip24) Clear() {
	self.V0 = self.Key0 ^ 0x736f6d6570736575
	self.V1 = self.Key1 ^ 0x646f72616e646f6d
	self.V2 = self.Key0 ^ 0x6c7967656e657261
	self.V3 = self.Key1 ^ 0x7465646279746573
	self.Finalized = false
	self.TotalLen = 0
	self.BlockBufferLen = 0
}
