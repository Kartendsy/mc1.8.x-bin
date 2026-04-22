package bin

import (
	"bytes"
	"encoding/binary"
	"math"
	"sync"
)

type Writer struct {
	buf *bytes.Buffer
}

var bufferPool = sync.Pool{
	New: func() any {
		return &Writer{
			buf: new(bytes.Buffer),
		}
	},
}

func NewWriter() *Writer {
	w := bufferPool.Get().(*Writer)
	w.buf.Reset()
	return w
}

func ReleaseWriter(w *Writer) {
	// Reset buffer sebelum dikembalikan agar memori bisa digunakan ulang
	w.buf.Reset()
	bufferPool.Put(w)
}

func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *Writer) WriteByte(val byte) error {
	return w.buf.WriteByte(val)
}

func (w *Writer) WriteVarInt(val int32) error {
	uval := uint32(val)
	for {
		temp := byte(uval & 0x7F)
		uval >>= 7
		if uval != 0 {
			temp |= 0x80
		}
		if err := w.buf.WriteByte(temp); err != nil {
			return err
		}
		if uval == 0 {
			break
		}
	}
	return nil
}

func (w *Writer) WriteString(val string) error {
	strBytes := []byte(val)
	if err := w.WriteVarInt(int32(len(strBytes))); err != nil {
		return err
	}
	_, err := w.buf.Write(strBytes)
	return err
}

func (w *Writer) WriteShort(val int16) error {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], uint16(val))
	_, err := w.buf.Write(b[:])
	return err
}

func (w *Writer) WriteUnsignedShort(val uint16) error {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], val)
	_, err := w.buf.Write(b[:])
	return err
}

func (w *Writer) WriteInt(val int32) error {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(val))
	_, err := w.buf.Write(b[:])
	return err
}

func (w *Writer) WriteLong(val int64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(val))
	_, err := w.buf.Write(b[:])
	return err
}

func (w *Writer) WriteFloat(val float32) error {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], math.Float32bits(val))
	_, err := w.buf.Write(b[:])
	return err
}

func (w *Writer) WriteDouble(val float64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], math.Float64bits(val))
	_, err := w.buf.Write(b[:])
	return err
}

func (w *Writer) WriteBool(val bool) error {
	if val {
		return w.buf.WriteByte(1)
	}
	return w.buf.WriteByte(0)
}

func (w *Writer) WritePosition(x, y, z int32) error {
	// Bit-packing untuk koordinat
	val := (int64(x&0x3FFFFFF) << 38) | (int64(y&0xFFF) << 26) | int64(z&0xFFFFFF)
	return w.WriteLong(val)
}

func (w *Writer) WriteBytes(data []byte) error {
	_, err := w.buf.Write(data)
	return err
}
