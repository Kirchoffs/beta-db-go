package transaction

import (
    "beta-db-go/common"
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
    "encoding/binary"
    "fmt"
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestStartRecord(t *testing.T) {
    path := filepath.Join(os.TempDir(), "test_dir")
    t.Log("Testing path: ", path)
    defer os.RemoveAll(path)

    fileManager, _ := fm.NewFileManager(path, 4096)
    logManager, _ := lm.NewLogManager(fileManager, "test_log")

    txn := uint64(42)
    page := fm.NewPageBySize(1024)
    page.SetUInt(0, uint64(START))
    page.SetUInt(common.UINT64_LEN, txn)

    startRecord := NewStartRecord(page, logManager)
    expetedStartRecordStr := fmt.Sprintf("<START %d>", txn)
    require.Equal(t, expetedStartRecordStr, startRecord.ToString())

    _, err := WriteStartRecord(logManager, txn)
    require.Nil(t, err)

    iter := logManager.Iterator()
    record := iter.Next()
    recordOp := binary.LittleEndian.Uint64(record[:common.UINT64_LEN])
    recordTxn := binary.LittleEndian.Uint64(record[common.UINT64_LEN:])
    require.Equal(t, uint64(START), recordOp)
    require.Equal(t, txn, recordTxn)
}

func TestSetStringRecord(t *testing.T) {
    path := filepath.Join(os.TempDir(), "test_dir")
    t.Log("Testing path: ", path)
    defer os.RemoveAll(path)

    fileManager, _ := fm.NewFileManager(path, 4096)
    logManager, _ := lm.NewLogManager(fileManager, "test_log")

    txn := uint64(42)
    fileName := "test_file"
    blockNum := uint64(1)
    offset := uint64(89)
    val := "test_val"

    page := fm.NewPageBySize(1024)
    page.SetUInt(0, uint64(SETSTRING))
    page.SetUInt(common.UINT64_LEN, txn)
    page.SetString(common.UINT64_LEN*2, fileName)
    page.SetUInt(common.UINT64_LEN*2+fm.MaxLengthForString(fileName), blockNum)
    page.SetUInt(common.UINT64_LEN*3+fm.MaxLengthForString(fileName), offset)
    page.SetString(common.UINT64_LEN*4+fm.MaxLengthForString(fileName), val)

    setStringRecord := NewSetStringRecord(page)
    expetedSetStringRecordStr := fmt.Sprintf("<SETSTRING %d %s %d %d %s>", txn, fileName, blockNum, offset, val)
    require.Equal(t, expetedSetStringRecordStr, setStringRecord.ToString())

    _, err := WriteSetStringRecord(logManager, txn, fm.NewBlockId(fileName, blockNum), offset, val)
    require.Nil(t, err)

    iter := logManager.Iterator()
    record := iter.Next()

    recordOp := binary.LittleEndian.Uint64(record[:common.UINT64_LEN])
    recordTxn := binary.LittleEndian.Uint64(record[common.UINT64_LEN:])
    recordFileName := string(record[common.UINT64_LEN*3 : common.UINT64_LEN*3+uint64(len([]byte(fileName)))])
    recordBlockNum := binary.LittleEndian.Uint64(record[common.UINT64_LEN*3+uint64(len([]byte(fileName))) : common.UINT64_LEN*4+uint64(len([]byte(fileName)))])
    recordOffset := binary.LittleEndian.Uint64(record[common.UINT64_LEN*4+uint64(len([]byte(fileName))) : common.UINT64_LEN*5+uint64(len([]byte(fileName)))])
    recordVal := string(record[common.UINT64_LEN*6+uint64(len([]byte(fileName))) : common.UINT64_LEN*6+uint64(len([]byte(fileName)))+uint64(len([]byte(val)))])

    require.Equal(t, uint64(SETSTRING), recordOp)
    require.Equal(t, txn, recordTxn)
    require.Equal(t, fileName, recordFileName)
    require.Equal(t, blockNum, recordBlockNum)
    require.Equal(t, offset, recordOffset)
    require.Equal(t, val, recordVal)
}

func TestSetIntRecord(t *testing.T) {
    path := filepath.Join(os.TempDir(), "test_dir")
    t.Log("Testing path: ", path)
    defer os.RemoveAll(path)

    fileManager, _ := fm.NewFileManager(path, 4096)
    logManager, _ := lm.NewLogManager(fileManager, "test_log")

    txn := uint64(42)
    fileName := "test_file"
    blockNum := uint64(1)
    offset := uint64(89)
    val := int64(6174)

    page := fm.NewPageBySize(1024)
    page.SetUInt(0, uint64(SETINT))
    page.SetUInt(common.UINT64_LEN, txn)
    page.SetString(common.UINT64_LEN*2, fileName)
    page.SetUInt(common.UINT64_LEN*2+fm.MaxLengthForString(fileName), blockNum)
    page.SetUInt(common.UINT64_LEN*3+fm.MaxLengthForString(fileName), offset)
    page.SetInt(common.UINT64_LEN*4+fm.MaxLengthForString(fileName), val)

    setIntRecord := NewSetIntRecord(page)
    expetedSetIntRecordStr := fmt.Sprintf("<SETINT %d %s %d %d %d>", txn, fileName, blockNum, offset, val)
    require.Equal(t, expetedSetIntRecordStr, setIntRecord.ToString())

    _, err := WriteSetIntRecord(logManager, txn, fm.NewBlockId(fileName, blockNum), offset, val)
    require.Nil(t, err)

    iter := logManager.Iterator()
    record := iter.Next()

    recordOp := binary.LittleEndian.Uint64(record[:common.UINT64_LEN])
    recordTxn := binary.LittleEndian.Uint64(record[common.UINT64_LEN:])
    recordFileName := string(record[common.UINT64_LEN*3 : common.UINT64_LEN*3+uint64(len([]byte(fileName)))])
    recordBlockNum := binary.LittleEndian.Uint64(record[common.UINT64_LEN*3+uint64(len([]byte(fileName))) : common.UINT64_LEN*4+uint64(len([]byte(fileName)))])
    recordOffset := binary.LittleEndian.Uint64(record[common.UINT64_LEN*4+uint64(len([]byte(fileName))) : common.UINT64_LEN*5+uint64(len([]byte(fileName)))])
    recordVal := binary.LittleEndian.Uint64(record[common.UINT64_LEN*5+uint64(len([]byte(fileName))) : common.UINT64_LEN*6+uint64(len([]byte(fileName)))])

    require.Equal(t, uint64(SETINT), recordOp)
    require.Equal(t, txn, recordTxn)
    require.Equal(t, fileName, recordFileName)
    require.Equal(t, blockNum, recordBlockNum)
    require.Equal(t, offset, recordOffset)
    require.Equal(t, val, int64(recordVal))
}

func TestCommitRecord(t *testing.T) {
    path := filepath.Join(os.TempDir(), "test_dir")
    t.Log("Testing path: ", path)
    defer os.RemoveAll(path)

    fileManager, _ := fm.NewFileManager(path, 4096)
    logManager, _ := lm.NewLogManager(fileManager, "test_log")

    txn := uint64(42)
    page := fm.NewPageBySize(1024)
    page.SetUInt(0, uint64(COMMIT))
    page.SetUInt(common.UINT64_LEN, txn)

    commitRecord := NewCommitRecord(page)
    expetedCommitRecordStr := fmt.Sprintf("<COMMIT %d>", txn)
    require.Equal(t, expetedCommitRecordStr, commitRecord.ToString())

    _, err := WriteCommitRecord(logManager, txn)
    require.Nil(t, err)

    iter := logManager.Iterator()
    record := iter.Next()
    recordOp := binary.LittleEndian.Uint64(record[:common.UINT64_LEN])
    recordTxn := binary.LittleEndian.Uint64(record[common.UINT64_LEN:])
    require.Equal(t, uint64(COMMIT), recordOp)
    require.Equal(t, txn, recordTxn)
}

func TestRollbackRecord(t *testing.T) {
    path := filepath.Join(os.TempDir(), "test_dir")
    t.Log("Testing path: ", path)
    defer os.RemoveAll(path)

    fileManager, _ := fm.NewFileManager(path, 4096)
    logManager, _ := lm.NewLogManager(fileManager, "test_log")

    txn := uint64(42)
    page := fm.NewPageBySize(1024)
    page.SetUInt(0, uint64(ROLLBACK))
    page.SetUInt(common.UINT64_LEN, txn)

    rollbackRecord := NewRollbackRecord(page)
    expetedRollbackRecordStr := fmt.Sprintf("<ROLLBACK %d>", txn)
    require.Equal(t, expetedRollbackRecordStr, rollbackRecord.ToString())

    _, err := WriteRollbackRecord(logManager, txn)
    require.Nil(t, err)

    iter := logManager.Iterator()
    record := iter.Next()
    recordOp := binary.LittleEndian.Uint64(record[:common.UINT64_LEN])
    recordTxn := binary.LittleEndian.Uint64(record[common.UINT64_LEN:])
    require.Equal(t, uint64(ROLLBACK), recordOp)
    require.Equal(t, txn, recordTxn)
}

func TestCheckpointRecord(t *testing.T) {
    path := filepath.Join(os.TempDir(), "test_dir")
    t.Log("Testing path: ", path)
    defer os.RemoveAll(path)

    fileManager, _ := fm.NewFileManager(path, 4096)
    logManager, _ := lm.NewLogManager(fileManager, "test_log")

    page := fm.NewPageBySize(1024)
    page.SetUInt(0, uint64(CHECKPOINT))

    checkpointRecord := NewCheckpointRecord()
    expetedCheckpointRecordStr := "<CHECKPOINT>"
    require.Equal(t, expetedCheckpointRecordStr, checkpointRecord.ToString())

    _, err := WriteCheckpointRecord(logManager)
    require.Nil(t, err)

    iter := logManager.Iterator()
    record := iter.Next()
    recordOp := binary.LittleEndian.Uint64(record[:common.UINT64_LEN])
    require.Equal(t, uint64(CHECKPOINT), recordOp)
}
