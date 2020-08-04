module tests

go 1.11

require (
	github.com/Blank-Xu/sqlx-adapter v0.0.0
	github.com/casbin/casbin/v2 v2.8.7
	github.com/denisenkom/go-mssqldb v0.0.0-20200620013148-b91950f658ec
	github.com/go-sql-driver/mysql v1.4.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.0.0
	github.com/mattn/go-sqlite3 v1.9.0
)

replace (
	github.com/Blank-Xu/sqlx-adapter => ../.
	google.golang.org/appengine v1.6.6 => github.com/golang/appengine v1.6.6
)
