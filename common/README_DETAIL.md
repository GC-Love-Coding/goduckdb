### Components

-   `FileBuffer`
-   `FileHandle`
-   `FileSystem`



#### `FileBuffer`

`FileBuffer` is a buffer that is employed to store the content which is read from or written to the file. We need to specify the buffer size when we create a `FileBuffer` instance. Note that, the buffer size parameter must be multiple times of **4096**. Inside the `FileBuffer` constructor, 