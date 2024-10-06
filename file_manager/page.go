package file_manager

import (
    "encoding/binary"
)

type Page struct {
    buffer []byte
}

func NewPageBySize(size uint64) *Page {
    return &Page{make([]byte, size)}
}

func NewPageByBytes(bytes []byte) *Page {
    return &Page{bytes}
}

func (page *Page) GetInt(offset uint64) uint64 {
    return binary.LittleEndian.Uint64(page.buffer[offset : offset+8])
}

func uint64ToByteArray(val uint64) []byte {
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, val)
    return bytes
}

func (page *Page) SetInt(offset uint64, val uint64) {
    copy(page.buffer[offset:], uint64ToByteArray(val))
}

func (page *Page) GetBytes(offset uint64) []byte {
    len := page.GetInt(offset)
    bytesBuffer := make([]byte, len)
    copy(bytesBuffer, page.buffer[offset+8:offset+8+len])
    return bytesBuffer
}

func (page *Page) SetBytes(offset uint64, bytes []byte) {
    page.SetInt(offset, uint64(len(bytes)))
    copy(page.buffer[offset+8:], bytes)
}

func (page *Page) GetString(offset uint64) string {
    return string(page.GetBytes(offset))
}

func (page *Page) SetString(offset uint64, str string) {
    page.SetBytes(offset, []byte(str))
}

func (page *Page) MaxLengthForString(str string) uint64 {
    bytes := []byte(str)
    return uint64(len(bytes)) + 8
}

func (page *Page) contents() []byte {
    return page.buffer
}
