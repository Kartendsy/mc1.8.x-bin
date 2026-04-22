package bin

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
	}
}

func (r *Reader) ReadByte() (byte, error) {
	return r.reader.ReadByte()
}

func (r *Reader) ReadVarInt() (int32, error) {
	var value uint32
	var pos int32
	for {
		currByte, err := r.reader.ReadByte()
		if err != nil {
			return 0, err
		}
		value |= uint32(currByte&0x7F) << uint32(pos)
		if (currByte & 0x80) == 0 {
			break
		}

		pos += 7
		if pos >= 32 {
			return 0, errors.New("VarInt too big")
		}
	}
	return int32(value), nil
}

func (r *Reader) ReadShort() (int16, error) {
	var bf [2]byte
	if _, err := io.ReadFull(r.reader, bf[:]); err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(bf[:])), nil
}

func (r *Reader) ReadUnsignedShort() (uint16, error) {
	var bf [2]byte
	if _, err := io.ReadFull(r.reader, bf[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(bf[:]), nil
}

func (r *Reader) ReadInt() (int32, error) {
	var bf [4]byte
	if _, err := io.ReadFull(r.reader, bf[:]); err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(bf[:])), nil
}

func (r *Reader) ReadLong() (int64, error) {
	var bf [8]byte
	if _, err := io.ReadFull(r.reader, bf[:]); err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(bf[:])), nil
}

func (r *Reader) ReadFloat() (float32, error) {
	var bf [4]byte
	if _, err := io.ReadFull(r.reader, bf[:]); err != nil {
		return 0, nil
	}
	return math.Float32frombits(binary.BigEndian.Uint32(bf[:])), nil
}

func (r *Reader) ReadDouble() (float64, error) {
	var bf [8]byte
	if _, err := io.ReadFull(r.reader, bf[:]); err != nil {
		return 0, err
	}
	return math.Float64frombits(binary.BigEndian.Uint64(bf[:])), nil
}

func (r *Reader) ReadBool() (bool, error) {
	val, err := r.ReadByte()
	return val != 0, err
}

func (r *Reader) ReadString() (string, error) {
	length, err := r.ReadVarInt()
	if err != nil {
		return "", err
	}
	if length < 0 || length > 32767 {
		return "", errors.New("too long")
	}
	bf := make([]byte, length)
	if _, err := io.ReadFull(r.reader, bf); err != nil {
		return "", err
	}
	return string(bf), nil
}

func (r *Reader) ReadPosition() (int32, int32, int32, error) {
	val, err := r.ReadLong()
	if err != nil {
		return 0, 0, 0, err
	}

	x := int32(val >> 38)
	y := int32((val >> 26) & 0xFFF)
	z := int32(val << 38 >> 38)
	return x, y, z, nil
}

func (r *Reader) ReadBytes(n int) ([]byte, error) {
	bf := make([]byte, n)
	_, err := io.ReadFull(r.reader, bf)
	return bf, err
}

func (r *Reader) ReadPacket() ([]byte, int32, error) {
	length, err := r.ReadVarInt()
	if err != nil {
		return nil, 0, err
	}
	data, err := r.ReadBytes(int(length))
	return data, length, err
}
