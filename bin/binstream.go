package bin

import (
	"encoding/binary"
	"errors"
	"io"
)

type BinStream struct {
	r io.Reader
}

func NewBinStream(r io.Reader) *BinStream {
	return &BinStream{r: r}
}

func (b *BinStream) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := io.ReadFull(b.r, buf[:])
	return buf[0], err
}

func (b *BinStream) ReadVarInt() (int32, error) {
	var value uint32
	var position int32

	for {
		currentByte, err := b.ReadByte()
		if err != nil {
			return 0, err
		}

		value |= uint32(currentByte&0x7F) << uint32(position)

		if (currentByte & 0x80) == 0 {
			break
		}

		position += 7
		if position >= 32 {
			return 0, errors.New("VarInt terlalu besar")
		}
	}
	return int32(value), nil
}

func (b *BinStream) ReadDouble() (float64, error) {
	var val float64
	err := binary.Read(b.r, binary.BigEndian, &val)
	return val, err
}

func (b *BinStream) ReadFloat() (float32, error) {
	var val float32
	err := binary.Read(b.r, binary.BigEndian, &val)
	return val, err
}

func (b *BinStream) ReadBool() (bool, error) {
	val, err := b.ReadByte()
	return val != 0, err
}

func (b *BinStream) ReadString() (string, error) {
	length, err := b.ReadVarInt()
	if err != nil {
		return "", err
	}

	if length < 0 || length > 32767 {
		return "", errors.New("panjang string tidak valid")
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(b.r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (b *BinStream) ReadUnsignedShort() (uint16, error) {
	var val uint16
	err := binary.Read(b.r, binary.BigEndian, &val)
	return val, err
}

func (b *BinStream) ReadShort() (int16, error) {
	var val int16
	err := binary.Read(b.r, binary.BigEndian, &val)
	return val, err
}

func (b *BinStream) ReadLong() (int64, error) {
	var val int64
	err := binary.Read(b.r, binary.BigEndian, &val)
	return val, err
}

func (b *BinStream) ReadPosition() (int32, int32, int32, error) {
	val, err := b.ReadLong()
	if err != nil {
		return 0, 0, 0, err
	}

	x := int32(val >> 38)
	y := int32((val >> 26) & 0xFFF)
	z := int32(val << 38 >> 38)
	return x, y, z, nil
}

func (b *BinStream) ReadBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(b.r, buf)
	return buf, err
}

func (b *BinStream) ReadPacket() ([]byte, int32, error) {
	len, err := b.ReadVarInt()
	if err != nil {
		return nil, 0, err
	}
	data, err := b.ReadBytes(int(len))
	return data, len, err
}
