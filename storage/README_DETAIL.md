### Components

-   `StorageManager`
    -   `SingleFileBlockManager`
        -   `MainHeader` and `DatabaseHeader`
        -   `MetaBlockReader`
    -   `WriteAheadLog`
-   `Transaction`
-   `Catalog`
-   `FileSystem`
-   `FileHandle`
-   
-   

`FileSystem` and `FileHandle` may need to be done first since the `SingleFileBlockManager` component replies on these two components to manipulate the disk file.







