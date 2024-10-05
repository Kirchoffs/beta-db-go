package transaction

import fm "beta-db-go/file_manager"

type TransactionInterface interface {
    Commit()
    Rollback()
    Recover()
    Pin(blockId *fm.BlockId)
    Unpin(blockId *fm.BlockId)
    GetInt(blockId *fm.BlockId, offset uint64) int64
    GetString(blockId *fm.BlockId, offset uint64) string
    SetInt(blockId *fm.BlockId, offset uint64, val int64, okToLog bool)
    SetString(blockId *fm.BlockId, offset uint64, val string, okToLog bool)
    AvailableBuffers() uint32
    Size(fileName string) uint64
    BlockSize() uint64
    Append(fileName string)
}

type RECORD_TYPE uint64

const (
    CHECKPOINT RECORD_TYPE = iota
    START
    COMMIT
    ROLLBACK
    SETINT
    SETSTRING
)

const (
    END_OF_FILE = -1
)

type LogRecordInterface interface {
    Op() RECORD_TYPE
    Txn() uint64
    Undo(tx TransactionInterface)
    ToString() string
}
