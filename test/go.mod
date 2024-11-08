module github.com/Blank-Xu/sqlx-adapter-test

go 1.21

replace github.com/Blank-Xu/sqlx-adapter => ../.

require (
	github.com/Blank-Xu/sqlx-adapter v0.0.0-00010101000000-000000000000
	github.com/casbin/casbin/v2 v2.100.0
	github.com/go-sql-driver/mysql v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/microsoft/go-mssqldb v1.7.2
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.7.1 // indirect
	github.com/casbin/govaluate v1.2.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
