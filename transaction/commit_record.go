package transaction

import (
    "beta-db-go/common"
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
    "fmt"
)

type CommitRecord struct {
    txn uint64
}

func NewCommitRecord(page *fm.Page) *CommitRecord {
    txn := page.GetUInt(common.UINT64_LEN)
    return &CommitRecord{
        txn: txn,
    }
}

func (commitRecord *CommitRecord) Op() RECORD_TYPE {
    return COMMIT
}

func (commitRecord *CommitRecord) Txn() uint64 {
    return commitRecord.txn
}

func (commitRecord *CommitRecord) Undo(tx TransactionInterface) {}

func (commitRecord *CommitRecord) ToString() string {
    return fmt.Sprintf("<COMMIT %d>", commitRecord.txn)
}

func WriteCommitRecord(logManager *lm.LogManager, txn uint64) (uint64, error) {
    record := make([]byte, common.UINT64_LEN*2)
    page := fm.NewPageByBytes(record)
    page.SetUInt(0, uint64(COMMIT))
    page.SetUInt(common.UINT64_LEN, txn)

    return logManager.AppendRecord(record)
}
