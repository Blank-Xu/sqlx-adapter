# sqlx-adapter

[![Go Report Card](https://goreportcard.com/badge/github.com/Blank-Xu/sqlx-adapter)](https://goreportcard.com/report/github.com/Blank-Xu/sqlx-adapter)
[![Build Status](https://github.com/Blank-Xu/sqlx-adapter/actions/workflows/tests.yaml/badge.svg)](https://github.com/Blank-Xu/sqlx-adapter/actions)
[![Coverage Status](https://coveralls.io/repos/github/Blank-Xu/sqlx-adapter/badge.svg?branch=master)](https://coveralls.io/github/Blank-Xu/sqlx-adapter?branch=master)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Blank-Xu/sqlx-adapter)](https://pkg.go.dev/github.com/Blank-Xu/sqlx-adapter)
[![Release](https://img.shields.io/github/release/Blank-Xu/sqlx-adapter.svg)](https://github.com/Blank-Xu/sqlx-adapter/releases/latest)
[![Sourcegraph](https://sourcegraph.com/github.com/Blank-Xu/sqlx-adapter/-/badge.svg)](https://sourcegraph.com/github.com/Blank-Xu/sqlx-adapter?badge)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

---

This `sqlx-adapter` is a [sqlx](https://github.com/jmoiron/sqlx) adapter for [Casbin V2](https://github.com/casbin/casbin).

With this library, Casbin can load policy lines from sqlx supported databases or save policy lines.

## Tested Databases

### `master` branch

- SQLite3: [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- MySQL(v5.5): [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- MariaDB(v10.2): [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- PostgreSQL(v9.6): [github.com/lib/pq](https://github.com/lib/pq)
- SQL Server(v2008R2-SP3): [github.com/denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb)

### `oracle` branch

- Oracle(v11.2): [github.com/mattn/go-oci8](https://github.com/mattn/go-oci8)

## Installation

```sh
go get github.com/Blank-Xu/sqlx-adapter
```

## Simple Examples

### SQLite3

```go
package main

import (
    "log"

    sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
    "github.com/casbin/casbin/v2"
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // connect to the database first.
    db, err := sqlx.Connect("sqlite3", "file:test.db")
    if err != nil {
        panic(err)
    }

    // Initialize a Sqlx adapter and use it in a Casbin enforcer:
    // The adapter will use the Sqlite3 table name "casbin_rule_test",
    // the default table name is "casbin_rule".
    // If it doesn't exist, the adapter will create it automatically.
    a, err := sqlxadapter.NewAdapter(db, "casbin_rule_test")
    if err != nil {
        panic(err)
    }

    e, err := casbin.NewEnforcer("test/testdata/rbac_model.conf", a)
    if err != nil {
        panic(err)
    }

    // Load the policy from DB.
    if err = e.LoadPolicy(); err != nil {
        log.Println("LoadPolicy failed, err: ", err)
    }

    // Check the permission.
    has, err := e.Enforce("alice", "data1", "read")
    if err != nil {
        log.Println("Enforce failed, err: ", err)
    }
    if !has {
        log.Println("do not have permission")
    }

    // Modify the policy.
    // e.AddPolicy(...)
    // e.RemovePolicy(...)

    // Save the policy back to DB.
    if err = e.SavePolicy(); err != nil {
        log.Println("SavePolicy failed, err: ", err)
    }
}
```

### MySQL

```go
package main

import (
    "log"
    "runtime"
    "time"

    sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
    "github.com/casbin/casbin/v2"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

func finalizer(db *sqlx.DB) {
    err := db.Close()
    if err != nil {
        panic(err)
    }
}

func main() {
    // connect to the database first.
    db, err := sqlx.Connect("mysql", "root:YourPassword@tcp(127.0.0.1:3306)/YourDBName")
    if err != nil {
        panic(err)
    }

    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(time.Minute * 10)

    // need to control by user, not the package
    runtime.SetFinalizer(db, finalizer)

    // Initialize a Sqlx adapter and use it in a Casbin enforcer:
    // The adapter will use the Sqlite3 table name "casbin_rule_test",
    // the default table name is "casbin_rule".
    // If it doesn't exist, the adapter will create it automatically.
    a, err := sqlxadapter.NewAdapter(db, "casbin_rule_test")
    if err != nil {
        panic(err)
    }

    e, err := casbin.NewEnforcer("test/testdata/rbac_model.conf", a)
    if err != nil {
        panic(err)
    }

    // Load the policy from DB.
    if err = e.LoadPolicy(); err != nil {
        log.Println("LoadPolicy failed, err: ", err)
    }

    // Check the permission.
    has, err := e.Enforce("alice", "data1", "read")
    if err != nil {
        log.Println("Enforce failed, err: ", err)
    }
    if !has {
        log.Println("do not have permission")
    }

    // Modify the policy.
    // e.AddPolicy(...)
    // e.RemovePolicy(...)

    // Save the policy back to DB.
    if err = e.SavePolicy(); err != nil {
        log.Println("SavePolicy failed, err: ", err)
    }
}
```

## Getting Help

- [Casbin](https://github.com/casbin/casbin)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.
