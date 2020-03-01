# Shrew [![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/) [![Go Report Card](https://goreportcard.com/badge/github.com/NicoNex/shrew)](https://goreportcard.com/report/github.com/NicoNex/shrew) [![License](https://img.shields.io/badge/license-GPL3-green.svg?style=flat)](https://github.com/NicoNex/shrew/blob/master/LICENSE)
A little shrew that stores backups and makes them accessible in a web page.


## Types
### Status
| Field | Type    | Description                                  |
|-------|---------|----------------------------------------------|
| ok    | boolean | Represent the status of a request.           |
| error | string  | Optional. Contains the error message if any. |

### Item
Represent a generic object (can be a file or an archive).
| Field   | Type   | Description                      |
|---------|--------|----------------------------------|
| name    | string | Name of the item.                |
| archive | string | Name of the archive.             |
| path    | string | Path of the file. Optional.      |
| sum     | string | Sha256sum of the file. Optional. |

### Archive
| Field | Type         | Description                             |
|-------|--------------|-----------------------------------------|
| name  | string       | Name of the archive.                    |
| files | string array | All the files contained in the archive. |


## Endpoints
### /
This returns all the archives present in the folder and all the files contained in each archive.

### /put
| Parameter | Type   | Required | Description                        |
|-----------|--------|----------|------------------------------------|
| archive   | string | true     | The name of the archive to upload. |

### /del
If the files field is present Shrew will remove only the files that belong to the specified archive and with the provided filename.
Instead, if in the request is present only the archive field, the entire archive will be removed.

| Parameter | Type         | Required | Description                                                            |
|-----------|--------------|----------|------------------------------------------------------------------------|
| archive   | string       | true     | The name of the archive to delete.                                     |
| files     | string array | false    | If present, shrew will delete only the specified files of the archive. |

### /get
| Paramenter  | Type   | Required | Description                                                              |
|-------------|--------|----------|--------------------------------------------------------------------------|
| archive     | string | true     | Name of the archive to download.                                         |
| compression | string | false    | Name of the compression to use. Available ones are: zip, targz, tarzstd. |
