# Notes
## Issue
```
func (fileManager *FileManager) Write(blockId *BlockId, page *Page) (int, error) {
    fileManager.mu.Lock()
    defer fileManager.mu.Unlock()

    file, err := fileManager.getFile(blockId.fileName)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    count, err := file.WriteAt(page.contents(), int64(blockId.blockNum*fileManager.blockSize))
    if err != nil {
        return 0, err
    }

    return count, nil
}
```

```
func (fileManager *FileManager) Read(blockId *BlockId, page *Page) (int, error) {
    fileManager.mu.Lock()
    defer fileManager.mu.Unlock()

    file, err := fileManager.getFile(blockId.fileName)
    if err != nil {
        return 0, err
    }
    defer file.Close()
    count, err := file.ReadAt(page.contents(), int64(blockId.blockNum*fileManager.blockSize))
    if err != nil {
        return 0, err
    }

    return count, nil
}
```

Above code will not pass the test `TestFileManager`.  

The issue stems from how file systems and operating systems handle file operations, especially in the context of buffering and caching. Let's break it down:

1. __Buffering__: When you write to a file, the operating system often doesn't immediately write the data to the physical disk. Instead, it stores the data in a buffer in memory for efficiency. This is because disk operations are much slower than memory operations.

2. __Caching__: Similarly, when you read from a file, the operating system may cache the data in memory to speed up subsequent reads.


3. __File closing__: When you close a file, it signals to the operating system that you're done with the file. This can trigger several actions:
- Any buffered writes are flushed to disk.
- File handles are released.
- Caches associated with the file might be __cleared__ or __invalidated__.

4. __File opening__: When you open a file, the operating system needs to set up new file handles and potentially read fresh data from the disk (not from cache because the file was closed and caches might have been cleared).

What happens in the code:
- [1]. You write data to the file.
- [2]. You immediately close the file.
- [3]. You then open the file again to read.
- [4]. You read the data and close the file again.

The problem occurs between steps 2 and 3. When you close the file after writing, there's no guarantee that the data has actually been written to the disk. The operating system might still be in the process of flushing its buffers.

When you immediately reopen the file to read, you will not read from the cache because cache is cleared after the closing the file, and you might be reading from a state of the file that doesn't yet reflect your recent write operation because the write hasn't been fully committed to disk. This is especially true if the read operation happens very quickly after the write operation.
