package buffer_manager

import (
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
    "errors"
    "sync"
    "time"
)

const (
    MAX_WAIT_TIME = 3
)

type BufferManager struct {
    bufferPool   []*Buffer
    numAvailable uint32
    mu           sync.Mutex
    cond         *sync.Cond
}

func NewBufferManager(fileManager *fm.FileManager, logManager *lm.LogManager, numBuffers uint32) *BufferManager {
    bufferManager := &BufferManager{
        numAvailable: numBuffers,
    }

    bufferManager.cond = sync.NewCond(&bufferManager.mu)

    for i := uint32(0); i < numBuffers; i++ {
        bufferManager.bufferPool = append(bufferManager.bufferPool, NewBuffer(fileManager, logManager))
    }

    return bufferManager
}

func (bufferManager *BufferManager) NumAvailable() uint32 {
    bufferManager.mu.Lock()
    defer bufferManager.mu.Unlock()
    return bufferManager.numAvailable
}

func (bufferManager *BufferManager) FlushAll(txn uint64) {
    bufferManager.mu.Lock()
    defer bufferManager.mu.Unlock()

    for _, buffer := range bufferManager.bufferPool {
        if buffer.Txn() == txn {
            buffer.Flush()
        }
    }
}

func (bufferManager *BufferManager) Pin(blockId *fm.BlockId) (*Buffer, error) {
    bufferManager.mu.Lock()
    defer bufferManager.mu.Unlock()

    start := time.Now()
    buffer := bufferManager.tryPin(blockId)
    for buffer == nil && !waitingForTooLong(start) {
        bufferManager.cond.Wait()
        buffer = bufferManager.tryPin(blockId)
    }

    if buffer == nil {
        return nil, errors.New(`no available buffer found, deadlock may occur`)
    }

    return buffer, nil
}

func (bufferManager *BufferManager) Unpin(buffer *Buffer) {
    bufferManager.mu.Lock()
    defer bufferManager.mu.Unlock()

    if buffer == nil {
        return
    }

    buffer.Unpin()
    if !buffer.IsPinned() {
        bufferManager.numAvailable++
        bufferManager.cond.Signal()
    }
}

func waitingForTooLong(start time.Time) bool {
    return time.Since(start).Seconds() > MAX_WAIT_TIME
}

func (bufferManager *BufferManager) tryPin(blockId *fm.BlockId) *Buffer {
    buffer := bufferManager.findExistingBuffer(blockId)
    if buffer == nil {
        buffer = bufferManager.chooseUnpinnedBuffer()
        if buffer == nil {
            return nil
        }
        buffer.AssignToBlock(blockId)
    }

    if !buffer.IsPinned() {
        bufferManager.numAvailable--
    }

    buffer.Pin()
    return buffer
}

func (bufferManager *BufferManager) findExistingBuffer(blockId *fm.BlockId) *Buffer {
    for _, buffer := range bufferManager.bufferPool {
        if buffer.BlockId() != nil && buffer.BlockId().Equals(blockId) {
            return buffer
        }
    }

    return nil
}

func (bufferManager *BufferManager) chooseUnpinnedBuffer() *Buffer {
    for _, buffer := range bufferManager.bufferPool {
        if !buffer.IsPinned() {
            return buffer
        }
    }

    return nil
}
