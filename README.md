# wuu

wuu is a simple pastebin service built in Go. Simply post any text data to
the server and it'll give you a paste link.

## Usage

Simply compile using the [go toolchain](https://golang.org/dl/) then specify
a directory for the leveldb database to be stored using the "-db" flag.
By default wuu listens on all interfaces, this can be changed using the "-addr"
flag.

Start wuu.
```
wuu -db="~/wuudb"
```

Start wuu and force listening to loopback interface.
```
wuu -db="~/wuudb" -addr="localhost:8080"
```