package transaction

import (
    "beta-db-go/common"
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
    "fmt"
)

type RollbackRecord struct {
    txn uint64
}

func NewRollbackRecord(page *fm.Page) *RollbackRecord {
    txn := page.GetUInt(common.UINT64_LEN)
    return &RollbackRecord{
        txn: txn,
    }
}

func (rollbackRecord *RollbackRecord) Op() RECORD_TYPE {
    return ROLLBACK
}

func (rollbackRecord *RollbackRecord) Txn() uint64 {
    return rollbackRecord.txn
}

func (rollbackRecord *RollbackRecord) Undo(tx TransactionInterface) {}

func (rollbackRecord *RollbackRecord) ToString() string {
    return fmt.Sprintf("<ROLLBACK %d>", rollbackRecord.txn)
}

func WriteRollbackRecord(logManager *lm.LogManager, txn uint64) (uint64, error) {
    record := make([]byte, common.UINT64_LEN*2)
    page := fm.NewPageByBytes(record)
    page.SetUInt(0, uint64(ROLLBACK))
    page.SetUInt(common.UINT64_LEN, txn)

    return logManager.AppendRecord(record)
}
