### Components

-   `FileSystem`
-   `FileBuffer`
-   `FileHandle`



##### 1. `FileSystem`

For `FileSystem` in duckdb, as its name implies, is an abtraction of file system. This `FileSystem` class provides method to manipulate file and directory. And the detailed methods signature are listed below.

```cpp
unique_ptr<FileHandle> OpenFile(const char* path, uint8_t flags, FileLockType lock);
unique_ptr<FileHandle> OpenFile(string& path, uint8_t flags, FileLockType lock);
void Read(FileHandle& handle, void* buffer, int64_t nr_bytes, index_t location);
void Write(FileHandle& handle, void* buffer, int64_t nr_bytes, index_t location);
int64_t Read(FileHandle& handle, void* buffer, int64_t nr_bytes);
int64_t Write(FileHandle& handle, void* buffer, int64_t nr_bytes);
int64_t GetFileSize(FileHandle& handle);
bool DirectoryExists(const string& directory);
void CreateDirectory(const string& directory);
void RemoveDirectory(const string& directory);
bool ListFiles(const string& directory, std::function<void(string)> callback);
void MoveFile(const string& source, const string& target)
bool FileExists(const string& filename);
void RemoveFile(const string& filename);
string PathSeparator();
string JoinPath(const string& a, const string& path);
void FileSync(FileHandle &handle);
```

Note that, the `OpenFile` method returns a `FileHandle`. The detailed implementation will be dependent on the operating system. For Unix/Linux based os, the implementation is a class wraps file descriptor (fd). 

##### 2. `FileBuffer`

`FileBuffer` is a buffer that is employed to store the content which is read from or written to the file. We need to specify the buffer size when we create a `FileBuffer` instance. Note that, the buffer size parameter must be multiple times of **4096**. Inside the `FileBuffer` constructor, 





#### Reference

1.  
    https://pkg.go.dev/syscall#Errno