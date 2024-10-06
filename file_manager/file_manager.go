package file_manager

import (
    "os"
    "path/filepath"
    "strings"
    "sync"
)

type FileManager struct {
    dbDirectory string
    blockSize   uint64
    isNew       bool
    openFiles   map[string]*os.File
    mu          *sync.Mutex
}

func NewFileManager(dbDirectory string, blockSize uint64) (*FileManager, error) {
    fileManager := &FileManager{
        dbDirectory: dbDirectory,
        blockSize:   blockSize,
        isNew:       false,
        openFiles:   make(map[string]*os.File),
        mu:          &sync.Mutex{},
    }

    if _, err := os.Stat(dbDirectory); os.IsNotExist(err) {
        fileManager.isNew = true
        if err := os.Mkdir(dbDirectory, 0755); err != nil {
            return nil, err
        }
    } else {
        err := filepath.Walk(dbDirectory, func(path string, info os.FileInfo, err error) error {
            mode := info.Mode()
            if mode.IsRegular() {
                name := info.Name()
                if strings.HasPrefix(name, "temp") {
                    os.Remove(filepath.Join(dbDirectory, name))
                }
            }

            return nil
        })

        if err != nil {
            return nil, err
        }
    }

    return fileManager, nil
}

func (fileManager *FileManager) Close() error {
    fileManager.mu.Lock()
    defer fileManager.mu.Unlock()

    for _, file := range fileManager.openFiles {
        if err := file.Close(); err != nil {
            return err
        }
    }
    fileManager.openFiles = make(map[string]*os.File)
    return nil
}

func (fileManager *FileManager) getFile(fileName string) (*os.File, error) {
    if file, ok := fileManager.openFiles[fileName]; ok {
        return file, nil
    }

    file, err := os.OpenFile(filepath.Join(fileManager.dbDirectory, fileName), os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return nil, err
    }

    fileManager.openFiles[fileName] = file
    return file, nil
}

func (fileManager *FileManager) Read(blockId *BlockId, page *Page) (int, error) {
    fileManager.mu.Lock()
    defer fileManager.mu.Unlock()

    file, err := fileManager.getFile(blockId.fileName)
    if err != nil {
        return 0, err
    }

    count, err := file.ReadAt(page.contents(), int64(blockId.blockNum*fileManager.blockSize))
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (fileManager *FileManager) Write(blockId *BlockId, page *Page) (int, error) {
    fileManager.mu.Lock()
    defer fileManager.mu.Unlock()

    file, err := fileManager.getFile(blockId.fileName)
    if err != nil {
        return 0, err
    }

    count, err := file.WriteAt(page.contents(), int64(blockId.blockNum*fileManager.blockSize))
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (fileManager *FileManager) size(fileName string) (uint64, error) {
    file, err := fileManager.getFile(fileName)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    info, err := file.Stat()
    if err != nil {
        return 0, err
    }

    return uint64(info.Size() / int64(fileManager.blockSize)), nil
}

func (fileManager *FileManager) Append(fileName string, page *Page) (*BlockId, error) {
    fileManager.mu.Lock()
    defer fileManager.mu.Unlock()

    blockNum, err := fileManager.size(fileName)
    if err != nil {
        return nil, err
    }

    blockId := NewBlockId(fileName, blockNum)
    file, err := fileManager.getFile(fileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    bytes := make([]byte, fileManager.blockSize)
    _, err = file.WriteAt(bytes, int64(blockNum*fileManager.blockSize))
    if err != nil {
        return nil, err
    }

    return blockId, nil
}

func (fileManager *FileManager) IsNew() bool {
    return fileManager.isNew
}

func (fileManager *FileManager) BlockSize() uint64 {
    return fileManager.blockSize
}
