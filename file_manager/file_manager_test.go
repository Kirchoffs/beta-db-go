package file_manager

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestFileManager(t *testing.T) {
    path := filepath.Join(os.TempDir(), "test_dir")
    defer os.RemoveAll(path)
    fileManager, err := NewFileManager(path, 256)
    require.NoError(t, err)
    defer fileManager.Close()

    blockId := NewBlockId("test_file", 2)
    page := NewPageBySize(fileManager.BlockSize())

    pos1 := uint64(42)
    content1 := "hello, world!"
    page.SetString(pos1, content1)
    size := page.MaxLengthForString(content1)

    pos2 := pos1 + size
    content2 := uint64(89)
    page.SetInt(pos2, uint64(content2))

    fileManager.Write(blockId, page)

    page = NewPageBySize(fileManager.blockSize)
    fileManager.Read(blockId, page)
    content1Read := page.GetString(pos1)
    content2Read := page.GetInt(pos2)

    require.Equal(t, content1, content1Read)
    require.Equal(t, content2, content2Read)
}
