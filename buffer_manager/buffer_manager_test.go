package buffer_manager

import (
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestBufferManager(t *testing.T) {
    path := filepath.Join(os.TempDir(), "test_dir")
    t.Log("Testing path: ", path)
    defer os.RemoveAll(path)

    fileManager, _ := fm.NewFileManager(path, 4096)
    logManager, _ := lm.NewLogManager(fileManager, "test_log")
    bufferManager := NewBufferManager(fileManager, logManager, 3)

    firstBuffer, err := bufferManager.Pin(fm.NewBlockId("test_file", 1))
    require.Nil(t, err)
    firstPage := firstBuffer.Page()
    numToWrite := int64(42)
    firstPage.SetInt(89, numToWrite)
    firstBuffer.SetModified(1, 1)
    numToRead := firstPage.GetInt(89)
    require.Equal(t, numToWrite, numToRead)

    secondBuffer, err := bufferManager.Pin(fm.NewBlockId("test_file", 2))
    require.NotNil(t, secondBuffer)
    require.Nil(t, err)

    thirdBuffer, err := bufferManager.Pin(fm.NewBlockId("test_file", 3))
    require.NotNil(t, thirdBuffer)
    require.Nil(t, err)

    secondBufferAgain, err := bufferManager.Pin(fm.NewBlockId("test_file", 2))
    require.NotNil(t, secondBufferAgain)
    require.Nil(t, err)
}
