package transaction

import (
    bm "beta-db-go/buffer_manager"
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
)

type RecoveryManager struct {
    fileManager   *fm.FileManager
    logManager    *lm.LogManager
    bufferManager *bm.BufferManager
    tx            *Transaction
    txn           uint64
}

func NewRecoveryManager(tx *Transaction, txn uint64, logManager *lm.LogManager, bufferManager *bm.BufferManager) *RecoveryManager {
    recoveryManger := &RecoveryManager{
        tx:            tx,
        txn:           txn,
        logManager:    logManager,
        bufferManager: bufferManager,
    }

    WriteStartRecord(logManager, txn)

    return recoveryManger
}

func (recoveryManager *RecoveryManager) Commit() error {
    recoveryManager.bufferManager.FlushAll(recoveryManager.txn)
    lsn, err := WriteCommitRecord(recoveryManager.logManager, recoveryManager.txn)
    if err != nil {
        return err
    }
    recoveryManager.logManager.FlushByLSN(lsn)
    return nil
}

func (recoveryManager *RecoveryManager) Rollback() error {
    recoveryManager.doRollback()
    recoveryManager.bufferManager.FlushAll(recoveryManager.txn)
    lsn, err := WriteRollbackRecord(recoveryManager.logManager, recoveryManager.txn)
    if err != nil {
        return err
    }
    recoveryManager.logManager.FlushByLSN(lsn)
    return nil
}

func (recoveryManager *RecoveryManager) doRollback() {

}

func (recoveryManager *RecoveryManager) Recover() error {
    recoveryManager.doRecover()
    return nil
}

func (recoveryManager *RecoveryManager) doRecover() {
}
