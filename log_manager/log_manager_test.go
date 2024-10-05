package log_manager

import (
    fm "beta-db-go/file_manager"
    "fmt"
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/require"
)

func makeRecords(logManager *LogManager, start uint64, end uint64) {
    for i := start; i <= end; i++ {
        record := []byte(fmt.Sprintf("record-%d", i))
        logManager.AppendRecord(record)
    }
}

func TestLogManager(t *testing.T) {
    path := filepath.Join(os.TempDir(), "test_dir")
    t.Log("Testing path: ", path)
    defer os.RemoveAll(path)

    fileManager, _ := fm.NewFileManager(path, 512)
    logManager, err := NewLogManager(fileManager, "test_log")
    require.Nil(t, err)

    makeRecords(logManager, 1, 42)

    iterator := logManager.Iterator()
    recordNum := uint64(42)
    for iterator.HasNext() {
        record := iterator.Next()
        t.Log(string(record))
        require.Equal(t, []byte(fmt.Sprintf("record-%d", recordNum)), record)
        recordNum--
    }

    makeRecords(logManager, 43, 84)
    logManager.FlushByLSN(72)
    iterator = logManager.Iterator()
    recordNum = uint64(84)
    for iterator.HasNext() {
        record := iterator.Next()
        t.Log(string(record))
        require.Equal(t, []byte(fmt.Sprintf("record-%d", recordNum)), record)
        recordNum--
    }
}
