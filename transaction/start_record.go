package transaction

import (
    "beta-db-go/common"
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
    "fmt"
)

type StartRecord struct {
    txn uint64
}

func NewStartRecord(page *fm.Page) *StartRecord {
    txn := page.GetUInt(common.UINT64_LEN)
    return &StartRecord{
        txn: txn,
    }
}

func (startRecord *StartRecord) Op() RECORD_TYPE {
    return START
}

func (startRecord *StartRecord) Txn() uint64 {
    return startRecord.txn
}

func (startRecord *StartRecord) Undo(tx TransactionInterface) {}

func (startRecord *StartRecord) ToString() string {
    str := fmt.Sprintf("<START %d>", startRecord.txn)
    return str
}

func WriteStartRecord(logManager *lm.LogManager, txn uint64) (uint64, error) {
    record := make([]byte, common.UINT64_LEN*2)
    page := fm.NewPageByBytes(record)
    page.SetUInt(0, uint64(START))
    page.SetUInt(common.UINT64_LEN, txn)

    return logManager.AppendRecord(record)
}
