package file_manager

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestSetAndGetInt(t *testing.T) {
    page := NewPageBySize(256)
    val := uint64(65535)
    offset := uint64(42)
    page.SetInt(offset, val)

    valRead := page.GetInt(offset)
    require.Equal(t, val, valRead)
}

func TestSetAndGetByByteArray(t *testing.T) {
    page := NewPageBySize(256)
    val := []byte{0x01, 0x02, 0x03, 0x04}
    offset := uint64(42)
    page.SetBytes(offset, val)

    valRead := page.GetBytes(offset)
    require.Equal(t, val, valRead)
}

func TestSetAndGetByString(t *testing.T) {
    page := NewPageBySize(256)
    val := "hello, world!"
    offset := uint64(42)
    page.SetString(offset, val)

    valRead := page.GetString(offset)
    require.Equal(t, val, valRead)
}
