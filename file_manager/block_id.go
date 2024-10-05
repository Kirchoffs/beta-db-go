package file_manager

import (
    "crypto/sha256"
    "fmt"
)

type BlockId struct {
    fileName string
    blockNum uint64
}

func NewBlockId(fileName string, blockNum uint64) *BlockId {
    return &BlockId{fileName, blockNum}
}

func (blockId *BlockId) Equals(other *BlockId) bool {
    return blockId.fileName == other.fileName && blockId.blockNum == other.blockNum
}

func asSha256(obj interface{}) string {
    hash := sha256.New()
    hash.Write([]byte(fmt.Sprintf("%v", obj)))
    return fmt.Sprintf("%x", hash.Sum(nil))
}

func (blockId *BlockId) HashCode() string {
    return asSha256(*blockId)
}

func (blockId *BlockId) FileName() string {
    return blockId.fileName
}

func (blockId *BlockId) BlockNum() uint64 {
    return blockId.blockNum
}
