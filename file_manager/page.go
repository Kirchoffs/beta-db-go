package file_manager

import (
    "beta-db-go/common"
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

func (page *Page) GetInt(offset uint64) int64 {
    return int64(page.GetUInt(offset))
}

func (page *Page) SetInt(offset uint64, val int64) {
    page.SetUInt(offset, uint64(val))
}

func (page *Page) GetUInt(offset uint64) uint64 {
    return binary.LittleEndian.Uint64(page.buffer[offset : offset+common.UINT64_LEN])
}

func (page *Page) SetUInt(offset uint64, val uint64) {
    copy(page.buffer[offset:], common.Uint64ToByteArray(val))
}

func (page *Page) GetBytes(offset uint64) []byte {
    len := page.GetUInt(offset)
    bytesBuffer := make([]byte, len)
    copy(bytesBuffer, page.buffer[offset+common.UINT64_LEN:offset+common.UINT64_LEN+len])
    return bytesBuffer
}

func (page *Page) SetBytes(offset uint64, bytes []byte) {
    page.SetUInt(offset, uint64(len(bytes)))
    copy(page.buffer[offset+common.UINT64_LEN:], bytes)
}

func (page *Page) GetString(offset uint64) string {
    return string(page.GetBytes(offset))
}

func (page *Page) SetString(offset uint64, str string) {
    page.SetBytes(offset, []byte(str))
}

func MaxLengthForString(str string) uint64 {
    bytes := []byte(str)
    return uint64(len(bytes)) + common.UINT64_LEN
}

func (page *Page) contents() []byte {
    return page.buffer
}
