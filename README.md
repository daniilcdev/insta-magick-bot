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
goose -dir=schemas postgres "user=<db/user> password=<db/pass> dbname=<db/name> host=<db/host> port=<db/port> sslmode=disable" up
```
Create new migration file:
```shell
goose -dir=schemas create <brief_description> sql
```
