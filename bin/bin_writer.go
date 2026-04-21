package bin

import (
	"bytes"
	"encoding/binary"
)

type BinWriter struct {
	buf *bytes.Buffer
}

func NewBinWriter() *BinWriter {
	return &BinWriter{buf: new(bytes.Buffer)}
}

func (b *BinWriter) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *BinWriter) WriteByte(val byte) error {
	return b.buf.WriteByte(val)
}

func (b *BinWriter) WriteVarInt(value int32) error {
	uval := uint32(value)
	for {
		temp := byte(uval & 0x7F)
		uval >>= 7
		if uval != 0 {
			temp |= 0x80
		}
		if err := b.buf.WriteByte(temp); err != nil {
			return err
		}

		if uval == 0 {
			break
		}
	}
	return nil
}

func (b *BinWriter) WriteString(val string) error {
	strBytes := []byte(val)
	if err := b.WriteVarInt(int32(len(strBytes))); err != nil {
		return err
	}
	_, err := b.buf.Write(strBytes)
	return err
}

func (b *BinWriter) WriteUnsignedShort(val uint16) error {
	return binary.Write(b.buf, binary.BigEndian, val)
}

func (b *BinWriter) WriteShort(val int16) error {
	return binary.Write(b.buf, binary.BigEndian, val)
}

func (b *BinWriter) WriteLong(val int64) error {
	return binary.Write(b.buf, binary.BigEndian, val)
}

func (b *BinWriter) WriteInt(val int32) error {
	return binary.Write(b.buf, binary.BigEndian, val)
}

func (b *BinWriter) WriteDouble(val float64) error {
	return binary.Write(b.buf, binary.BigEndian, val)
}

func (b *BinWriter) WriteFloat(val float32) error {
	return binary.Write(b.buf, binary.BigEndian, val)
}

func (b *BinWriter) WriteBool(val bool) error {
	if val {
		return b.WriteByte(1)
	}
	return b.WriteByte(0)
}

func (b *BinWriter) WritePosition(x, y, z int32) error {
	val := (int64(x&0x3FFFFFF) << 38) | (int64(y&0xFFF) << 26) | int64(z&0xFFFFFF)
	return b.WriteLong(val)
}

func (b *BinWriter) WriteBytes(data []byte) error {
	_, err := b.buf.Write(data)
	return err
}

func (b *BinWriter) Reset() {
	b.buf.Reset()
}
