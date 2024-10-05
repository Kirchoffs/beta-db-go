package transaction

import (
    "beta-db-go/common"
    fm "beta-db-go/file_manager"
    lm "beta-db-go/log_manager"
    "fmt"
)

type SetStringRecord struct {
    txn     uint64
    blockId *fm.BlockId
    offset  uint64
    val     string
}

func NewSetStringRecord(page *fm.Page) *SetStringRecord {
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
    val := page.GetString(valPos)

    return &SetStringRecord{
        txn:     txn,
        blockId: blockId,
        offset:  offset,
        val:     val,
    }
}

func (setStringRecord *SetStringRecord) Op() RECORD_TYPE {
    return SETSTRING
}

func (setStringRecord *SetStringRecord) Txn() uint64 {
    return setStringRecord.txn
}

func (setStringRecord *SetStringRecord) Undo(tx TransactionInterface) {
    tx.Pin(setStringRecord.blockId)
    tx.SetString(setStringRecord.blockId, setStringRecord.offset, setStringRecord.val, false)
    tx.Unpin(setStringRecord.blockId)
}

func (setStringRecord *SetStringRecord) ToString() string {
    str := fmt.Sprintf("<SETSTRING %d %s %d %d %s>", setStringRecord.txn, setStringRecord.blockId.FileName(), setStringRecord.blockId.BlockNum(), setStringRecord.offset, setStringRecord.val)
    return str
}

func WriteSetStringRecord(logManager *lm.LogManager, txn uint64, blockId *fm.BlockId, offset uint64, val string) (uint64, error) {
    record := make([]byte, common.UINT64_LEN*4+fm.MaxLengthForString(blockId.FileName())+fm.MaxLengthForString(val))
    page := fm.NewPageByBytes(record)
    page.SetUInt(0, uint64(SETSTRING))
    page.SetUInt(common.UINT64_LEN, txn)
    page.SetString(common.UINT64_LEN*2, blockId.FileName())
    page.SetUInt(common.UINT64_LEN*2+fm.MaxLengthForString(blockId.FileName()), blockId.BlockNum())
    page.SetUInt(common.UINT64_LEN*3+fm.MaxLengthForString(blockId.FileName()), offset)
    page.SetString(common.UINT64_LEN*4+fm.MaxLengthForString(blockId.FileName()), val)

    return logManager.AppendRecord(record)
}
