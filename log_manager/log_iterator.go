package log_manager

import (
    "beta-db-go/common"
    fm "beta-db-go/file_manager"
)

type LogIterator struct {
    fileManager    *fm.FileManager
    currentBlockId *fm.BlockId
    logPage        *fm.Page
    currentPos     uint64
    boundary       uint64
}

func NewLogIterator(fileManager *fm.FileManager, blockId *fm.BlockId) *LogIterator {
    iterator := LogIterator{
        fileManager:    fileManager,
        currentBlockId: blockId,
    }

    iterator.logPage = fm.NewPageBySize(fileManager.BlockSize())
    err := iterator.moveToBlock(blockId)
    if err != nil {
        return nil
    }

    return &iterator
}

func (iterator *LogIterator) moveToBlock(blockId *fm.BlockId) error {
    _, err := iterator.fileManager.Read(blockId, iterator.logPage)
    if err != nil {
        return err
    }

    iterator.boundary = iterator.logPage.GetUInt(0)
    iterator.currentPos = iterator.boundary

    return nil
}

func (iterator *LogIterator) moveToCurrentBlock() error {
    return iterator.moveToBlock(iterator.currentBlockId)
}

func (iterator *LogIterator) Next() []byte {
    if iterator.currentPos == iterator.fileManager.BlockSize() {
        if iterator.currentBlockId.BlockNum() == 0 {
            return nil
        }
        iterator.currentBlockId = fm.NewBlockId(iterator.currentBlockId.FileName(), iterator.currentBlockId.BlockNum()-1)
        iterator.moveToCurrentBlock()
    }

    logRecord := iterator.logPage.GetBytes(iterator.currentPos)
    iterator.currentPos += common.UINT64_LEN + uint64(len(logRecord))

    return logRecord
}

func (iterator *LogIterator) HasNext() bool {
    return iterator.currentPos < iterator.fileManager.BlockSize() || iterator.currentBlockId.BlockNum() > 0
}
