package transaction

import (
    bm "beta-db-go/buffer_manager"
    fm "beta-db-go/file_manager"
)

type BufferList struct {
    pinnedBuffers map[*fm.BlockId]*bm.Buffer
    bufferManager *bm.BufferManager
    pinnedCounts  map[*fm.BlockId]uint64
}

func NewBufferList(bufferManager *bm.BufferManager) *BufferList {
    return &BufferList{
        pinnedBuffers: make(map[*fm.BlockId]*bm.Buffer),
        bufferManager: bufferManager,
        pinnedCounts:  make(map[*fm.BlockId]uint64),
    }
}

func (bufferList *BufferList) GetBuffer(blockId *fm.BlockId) *bm.Buffer {
    return bufferList.pinnedBuffers[blockId]
}

func (bufferList *BufferList) Pin(blockId *fm.BlockId) error {
    buffer, err := bufferList.bufferManager.Pin(blockId)
    if err != nil {
        return nil
    }

    bufferList.pinnedBuffers[blockId] = buffer
    bufferList.pinnedCounts[blockId]++
    return nil
}

func (bufferList *BufferList) Unpin(blockId *fm.BlockId) {
    buffer, exists := bufferList.pinnedBuffers[blockId]
    if !exists {
        return
    }

    bufferList.bufferManager.Unpin(buffer)
    bufferList.pinnedCounts[blockId]--

    if bufferList.pinnedCounts[blockId] == 0 {
        delete(bufferList.pinnedBuffers, blockId)
    }
}

func (bufferList *BufferList) UnpinAll() {
    for blockId := range bufferList.pinnedBuffers {
        bufferList.Unpin(blockId)
    }

    bufferList.pinnedBuffers = make(map[*fm.BlockId]*bm.Buffer)
    bufferList.pinnedCounts = make(map[*fm.BlockId]uint64)
}
