# shrew
A little shrew that stores backups and makes them accessible in a web page.


## Types
### Status
| Field | Type    | Description                                  |
|-------|---------|----------------------------------------------|
| ok    | boolean | Represent the status of a request.           |
| error | string  | Optional. Contains the error message if any. |

### Item
Represent a generic object (can be a file or an archive).
| Field   | Type   | Description                 |
|---------|--------|-----------------------------|
| name    | string | Name of the item.           |
| archive | string | Name of the archive.        |
| path    | string | Optional. Path of the file. |

### Archive
| Field | Type         | Description                             |
|-------|--------------|-----------------------------------------|
| name  | string       | Name of the archive.                    |
| files | string array | All the files contained in the archive. |


## Endpoints
### /
This returns all the archives present in the folder and all the files contained in each archive.

### /upload
| Parameter | Type   | Required | Description                        |
|-----------|--------|----------|------------------------------------|
| archive   | string | true     | The name of the archive to upload. |

### /delete
If the files field is present Shrew will remove only the files that belong to the specified archive and with the provided filename.
Instead, if in the request is present only the archive field, the entire archive will be removed.

| Parameter | Type         | Required | Description                                                            |
|-----------|--------------|----------|------------------------------------------------------------------------|
| archive   | string       | true     | The name of the archive to delete.                                     |
| files     | string array | false    | If present, shrew will delete only the specified files of the archive. |

### /download
Coming soon...
