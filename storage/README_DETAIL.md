### Components

-   `StorageManager`
    -   `SingleFileBlockManager`
        -   `MainHeader` and `DatabaseHeader`
        -   `MetaBlockReader`
    -   `WriteAheadLog`
-   `Transaction`
-   `Catalog`
-   
-   `Block`
    -   `BlockHeader`
    -   `DataBlock`

-   `BlockManager`
-   `SingleFileBlockManager`





##### 2. Write ahead log

-   Why need WAL
    -   Provide durability guarantee without the storage data structures to be flushed to disk, by persisting every state change as a command to the append only log.
-   WAL format

