package transaction

import (
    bm "beta-db-go/buffer_manager"
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
    "errors"
    "fmt"
    "sync"
)

var txn_mu sync.Mutex
var next_txn uint64

func nextTxn() uint64 {
    txn_mu.Lock()
    defer txn_mu.Unlock()
    next_txn++
    return next_txn
}

type Transaction struct {
    fileManager   *fm.FileManager
    logManager    *lm.LogManager
    bufferManager *bm.BufferManager
    txn           uint64
    buffers       *BufferList
}

func NewTransaction(fileManager *fm.FileManager, logManager *lm.LogManager, bufferManager *bm.BufferManager) *Transaction {
    tx := &Transaction{
        fileManager:   fileManager,
        logManager:    logManager,
        bufferManager: bufferManager,
        txn:           nextTxn(),
        buffers:       NewBufferList(bufferManager),
    }

    return tx
}

func (tx *Transaction) Commit() {
    fmt.Println(fmt.Sprintf("transaction %d committed", tx.txn))
    tx.buffers.UnpinAll()
}

func (tx *Transaction) Rollback() {
    fmt.Println(fmt.Sprintf("transaction %d rolled back", tx.txn))
    tx.buffers.UnpinAll()
}

func (tx *Transaction) Recover() {
    fmt.Println(fmt.Sprintf("transaction %d recovered", tx.txn))
}

func (tx *Transaction) Pin(blockId *fm.BlockId) error {
    return tx.buffers.Pin(blockId)
}

func (tx *Transaction) Unpin(blockId *fm.BlockId) {
    tx.buffers.Unpin(blockId)
}

func (tx *Transaction) GetInt(blockId *fm.BlockId, offset uint64) (int64, error) {
    buffer := tx.buffers.GetBuffer(blockId)
    if buffer == nil {
        return 0, errors.New("buffer not found")
    }

    return buffer.Page().GetInt(offset), nil
}

func (tx *Transaction) GetString(blockId *fm.BlockId, offset uint64) (string, error) {
    buffer := tx.buffers.GetBuffer(blockId)
    if buffer == nil {
        return "", errors.New("buffer not found")
    }

    return buffer.Page().GetString(offset), nil
}

func (tx *Transaction) SetInt(blockId *fm.BlockId, offset uint64, val int64, okToLog bool) error {
    buffer := tx.buffers.GetBuffer(blockId)
    if buffer == nil {
        return errors.New("buffer not found")
    }

    var lsn uint64
    var err error
    if okToLog {
        lsn, err = tx.recoveryManager.SetInt(buffer, offset, val)
        if err != nil {
            return err
        }
    }

    page := buffer.Page()
    page.SetInt(offset, val)
    buffer.SetModified(tx.txn, lsn)

    return nil
}

func (tx *Transaction) SetString(blockId *fm.BlockId, offset uint64, val string, okToLog bool) error {
    buffer := tx.buffers.GetBuffer(blockId)
    if buffer == nil {
        return errors.New("buffer not found")
    }

    var lsn uint64
    var err error
    if okToLog {
        lsn, err = tx.recoveryManager.SetString(buffer, offset, val)
        if err != nil {
            return err
        }
    }

    page := buffer.Page()
    page.SetString(offset, val)
    buffer.SetModified(tx.txn, lsn)

    return nil
}

func (tx *Transaction) AvailableBuffers() uint32 {
    return tx.bufferManager.NumAvailable()
}

func (tx *Transaction) Size(fileName string) (uint64, error) {
    return tx.fileManager.Size(fileName)
}

func (tx *Transaction) Append(fileName string) *fm.BlockId {
    blockId, err := tx.fileManager.Append(fileName)
    if err != nil {
        return nil
    }
    return blockId
}

func (tx *Transaction) BlockSize() uint64 {
    return tx.fileManager.BlockSize()
}
