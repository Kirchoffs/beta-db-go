package log_manager

import (
    "sync"

    "beta-db-go/common"
    fm "beta-db-go/file_manager"
)

type LogManager struct {
    fileManager       *fm.FileManager
    logFile           string
    logPage           *fm.Page
    currentLogBlockId *fm.BlockId
    latestLSN         uint64
    lastSavedLSN      uint64
    mu                sync.Mutex
}

func NewLogManager(fileManager *fm.FileManager, logFile string) (*LogManager, error) {
    logManager := LogManager{
        fileManager:  fileManager,
        logFile:      logFile,
        logPage:      fm.NewPageBySize(fileManager.BlockSize()),
        lastSavedLSN: 0,
        latestLSN:    0,
    }
    // Valid lsn will be greater than 0

    logSize, err := fileManager.Size(logFile)
    if err != nil {
        return nil, err
    }

    if logSize == 0 {
        newBlockId, err := logManager.appendNewBlock()
        if err != nil {
            return nil, err
        }
        logManager.currentLogBlockId = newBlockId
    } else {
        logManager.currentLogBlockId = fm.NewBlockId(logFile, logSize-1)
        fileManager.Read(logManager.currentLogBlockId, logManager.logPage)
    }

    return &logManager, nil
}

// logManager.logPage should be flushed before calling this function
func (logManager *LogManager) appendNewBlock() (*fm.BlockId, error) {
    fileManager := logManager.fileManager
    blockId, err := fileManager.Append(logManager.logFile)
    if err != nil {
        return nil, err
    }

    logManager.logPage.SetUInt(0, uint64(fileManager.BlockSize()))
    fileManager.Write(blockId, logManager.logPage)

    return blockId, nil
}

func (logManager *LogManager) FlushByLSN(lsn uint64) error {
    if lsn > logManager.lastSavedLSN {
        err := logManager.Flush()
        if err != nil {
            return err
        }

        logManager.lastSavedLSN = lsn
    }

    return nil
}

func (logManager *LogManager) Flush() error {
    fileManager := logManager.fileManager
    _, err := fileManager.Write(logManager.currentLogBlockId, logManager.logPage)
    if err != nil {
        return err
    }

    return nil
}

func (logManager *LogManager) AppendRecord(logRecord []byte) (uint64, error) {
    logManager.mu.Lock()
    defer logManager.mu.Unlock()

    boundary := logManager.logPage.GetUInt(0)
    logRecordSize := uint64(len(logRecord))
    bytesNeed := common.UINT64_LEN + logRecordSize
    var err error
    if int(boundary)-int(bytesNeed) < int(common.UINT64_LEN) {
        err = logManager.Flush()
        if err != nil {
            return logManager.latestLSN, err
        }

        logManager.currentLogBlockId, err = logManager.appendNewBlock()
        if err != nil {
            return logManager.latestLSN, err
        }

        boundary = logManager.logPage.GetUInt(0)
    }

    logRecordPos := boundary - bytesNeed
    logManager.logPage.SetBytes(logRecordPos, logRecord)
    logManager.logPage.SetUInt(0, logRecordPos)
    logManager.latestLSN++

    return logManager.latestLSN, nil
}

func (logManager *LogManager) Iterator() *LogIterator {
    logManager.Flush()
    return NewLogIterator(logManager.fileManager, logManager.currentLogBlockId)
}
