### Third-Party Dependencies
* ImageMagick 6

### Generate Queries
Use [sqlc](https://github.com/sqlc-dev/sqlc) CLI:
```shell
sqlc generate
```

### Migration
Use [Goose](https://github.com/pressly/goose) CLI:
```shell
goose -dir=schemas sqlite3 <path/to/db> up
```
Create new migration file:
```shell
goose -dir=schemas create <brief_description> sql
```

### Build on Raspberry Pi 4

To build successfully, install proper compiler:
```shell
apt install g++-aarch64-linux-gnu gcc-aarch64-linux-gnu
```
and set env variables before build:

`CC=aarch64-linux-gnu-gcc` `CXX=aarch64-linux-gnu-g++`