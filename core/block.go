package core

import (
	"MyChain/types"
	"encoding/binary"
	"io"
)

type Header struct {
	Version   uint32
	PrevBlock types.Hash
	Timestamp int64
	Height    uint32
	Nonce     uint64
}

// EncodeBinary
//
//	@Description: 将h编码入w
//	@receiver h
//	@param w
//	@return error
func (h *Header) EncodeBinary(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, &h.Version); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.PrevBlock); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Timestamp); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Height); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, &h.Nonce)
}

// DecodeBinary
//
//	@Description: 将r中的数据解码入h
//	@receiver h
//	@param r
//	@return error
func (h *Header) DecodeBinary(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.PrevBlock); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.Timestamp); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.Height); err != nil {
		return err
	}
	return binary.Read(r, binary.LittleEndian, &h.Nonce)
}

type Block struct {
	Header
	Transactions []Transaction
}
