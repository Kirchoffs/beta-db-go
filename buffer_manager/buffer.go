package buffer_manager

import (
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
)

type Buffer struct {
    fileManager *fm.FileManager
    logManager  *lm.LogManager
    page        *fm.Page
    blockId     *fm.BlockId
    pins        uint32
    txn         uint64
    lsn         uint64
}

func NewBuffer(fileManager *fm.FileManager, logManager *lm.LogManager) *Buffer {
    return &Buffer{
        fileManager: fileManager,
        logManager:  logManager,
        page:        fm.NewPageBySize(fileManager.BlockSize()),
        blockId:     nil,
        pins:        0,
        txn:         0,
        lsn:         0,
    }
}

func (buffer *Buffer) Page() *fm.Page {
    return buffer.page
}

func (buffer *Buffer) BlockId() *fm.BlockId {
    return buffer.blockId
}

func (buffer *Buffer) Txn() uint64 {
    return buffer.txn
}

func (buffer *Buffer) SetModified(txn uint64, lsn uint64) {
    buffer.txn = txn
    if lsn > 0 {
        buffer.lsn = lsn
    }
}

func (buffer *Buffer) IsPinned() bool {
    return buffer.pins > 0
}

func (buffer *Buffer) AssignToBlock(blockId *fm.BlockId) {
    buffer.Flush()
    buffer.blockId = blockId
    buffer.fileManager.Read(blockId, buffer.page)
    buffer.pins = 0
}

func (buffer *Buffer) Flush() {
    if buffer.txn > 0 {
        buffer.logManager.FlushByLSN(buffer.lsn)
        buffer.fileManager.Write(buffer.blockId, buffer.page)
        buffer.txn = 0
    }
}

func (buffer *Buffer) Pin() {
    buffer.pins++
}

func (buffer *Buffer) Unpin() {
    buffer.pins--
}
