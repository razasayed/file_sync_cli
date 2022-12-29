# Description

A CLI written in Go for synchronizing a source directory to a destination directory.

# Running the program

```
go run fs_cli.go sync -s path_to_source_folder -d path_to_target_folder
```

Or

```
go build fs_cli.go

./fs_cli sync -s path_to_source_folder -d path_to_target_folder
```

# Running the unit tests

```
go test -v
```
