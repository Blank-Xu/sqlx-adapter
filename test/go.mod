module github.com/Blank-Xu/sqlx-adapter-test

go 1.23.0

replace github.com/Blank-Xu/sqlx-adapter => ../.

require (
	github.com/Blank-Xu/sqlx-adapter v0.0.0-00010101000000-000000000000
	github.com/casbin/casbin/v2 v2.120.0
	github.com/go-sql-driver/mysql v1.9.3
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/microsoft/go-mssqldb v1.9.2
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.7.1 // indirect
	github.com/casbin/govaluate v1.3.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)
