package transaction

import (
    "beta-db-go/common"
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
)

type CheckpointRecord struct{}

func NewCheckpointRecord() *CheckpointRecord {
    return &CheckpointRecord{}
}

func (checkpointRecord *CheckpointRecord) Op() RECORD_TYPE {
    return CHECKPOINT
}

func (checkpointRecord *CheckpointRecord) Txn() uint64 {
    return 0
}

func (checkpointRecord *CheckpointRecord) Undo(tx TransactionInterface) {}

func (checkpointRecord *CheckpointRecord) ToString() string {
    return "<CHECKPOINT>"
}

func WriteCheckpointRecord(logManager *lm.LogManager) (uint64, error) {
    record := make([]byte, common.UINT64_LEN)
    page := fm.NewPageByBytes(record)
    page.SetUInt(0, uint64(CHECKPOINT))

    return logManager.AppendRecord(record)
}
