package common

import "encoding/binary"

func Uint64ToByteArray(val uint64) []byte {
    bytes := make([]byte, UINT64_LEN)
    binary.LittleEndian.PutUint64(bytes, val)
    return bytes
}

func ByteArrayToUint64(bytes []byte) uint64 {
    return binary.LittleEndian.Uint64(bytes)
}
