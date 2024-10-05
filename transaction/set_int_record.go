package transaction

import (
    "beta-db-go/common"
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
    "fmt"
)

type SetIntRecord struct {
    txn     uint64
    blockId *fm.BlockId
    offset  uint64
    val     int64
}

func NewSetIntRecord(page *fm.Page) *SetIntRecord {
    txnPos := common.UINT64_LEN
    txn := page.GetUInt(txnPos)

    fileNamePos := txnPos + common.UINT64_LEN
    fileName := page.GetString(fileNamePos)

    blockNumPos := fileNamePos + fm.MaxLengthForString(fileName)
    blockNum := page.GetUInt(blockNumPos)

    blockId := fm.NewBlockId(fileName, blockNum)

    offsetPos := blockNumPos + common.UINT64_LEN
    offset := page.GetUInt(offsetPos)

    valPos := offsetPos + common.UINT64_LEN
    val := page.GetInt(valPos)

    return &SetIntRecord{
        txn:     txn,
        blockId: blockId,
        offset:  offset,
        val:     val,
    }
}

func (setIntRecord *SetIntRecord) Op() RECORD_TYPE {
    return SETINT
}

func (setIntRecord *SetIntRecord) Txn() uint64 {
    return setIntRecord.txn
}

func (setIntRecord *SetIntRecord) Undo(tx TransactionInterface) {
    tx.Pin(setIntRecord.blockId)
    tx.SetInt(setIntRecord.blockId, setIntRecord.offset, setIntRecord.val, false)
    tx.Unpin(setIntRecord.blockId)
}

func (setIntRecord *SetIntRecord) ToString() string {
    str := fmt.Sprintf("<SETINT %d %s %d %d %d>", setIntRecord.txn, setIntRecord.blockId.FileName(), setIntRecord.blockId.BlockNum(), setIntRecord.offset, setIntRecord.val)
    return str
}

func WriteSetIntRecord(logManager *lm.LogManager, txn uint64, blockId *fm.BlockId, offset uint64, val int64) (uint64, error) {
    record := make([]byte, common.UINT64_LEN*5+fm.MaxLengthForString(blockId.FileName()))
    page := fm.NewPageByBytes(record)
    page.SetUInt(0, uint64(SETINT))
    page.SetUInt(common.UINT64_LEN, txn)
    page.SetString(common.UINT64_LEN*2, blockId.FileName())
    page.SetUInt(common.UINT64_LEN*2+fm.MaxLengthForString(blockId.FileName()), blockId.BlockNum())
    page.SetUInt(common.UINT64_LEN*3+fm.MaxLengthForString(blockId.FileName()), offset)
    page.SetInt(common.UINT64_LEN*4+fm.MaxLengthForString(blockId.FileName()), val)

    return logManager.AppendRecord(record)
}
