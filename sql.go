// Copyright 2020 by Blank-Xu. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlxadapter

// general sql
const (
	sqlCreatTable = `
CREATE TABLE %s(
    p_type varchar(32),
    v0     varchar(255),
    v1     varchar(255),
    v2     varchar(255),
    v3     varchar(255),
    v4     varchar(255),
    v5     varchar(255)
)`
	sqlCreateIndex   = "CREATE INDEX idx_%s_ptype ON %s (p_type, v0, v1)"
	sqlTruncateTable = "TRUNCATE TABLE %s"
	sqlIsTableExist  = "SELECT 1 FROM %s"
	sqlInsertRow     = "INSERT INTO %s (%s) VALUES (%s)"
	sqlDeleteAll     = "DELETE FROM %s"
	sqlDeleteByArgs  = "DELETE FROM %s WHERE p_type = ?"
	sqlSelectAll     = "SELECT * FROM %s"
	sqlSelectWhere   = "SELECT * FROM %s WHERE "
)

// for Sqlite3
const (
	sqlCreateTableSqlite3 = `
CREATE TABLE IF NOT EXISTS %s(
    p_type varchar(32)  DEFAULT '' NOT NULL,
    v0     varchar(255) DEFAULT '' NOT NULL,
    v1     varchar(255) DEFAULT '' NOT NULL,
    v2     varchar(255) DEFAULT '' NOT NULL,
    v3     varchar(255) DEFAULT '' NOT NULL,
    v4     varchar(255) DEFAULT '' NOT NULL,
    v5     varchar(255) DEFAULT '' NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_%s_ptype ON %s (p_type, v0, v1);`
	sqlTruncateTableSqlite3 = "DROP TABLE IF EXISTS %s;" + sqlCreateTableSqlite3
)

// for Mysql
const (
	sqlCreatTableMysql = `
CREATE TABLE IF NOT EXISTS %s(
    p_type varchar(32)  DEFAULT '' NOT NULL,
    v0     varchar(255) DEFAULT '' NOT NULL,
    v1     varchar(255) DEFAULT '' NOT NULL,
    v2     varchar(255) DEFAULT '' NOT NULL,
    v3     varchar(255) DEFAULT '' NOT NULL,
    v4     varchar(255) DEFAULT '' NOT NULL,
    v5     varchar(255) DEFAULT '' NOT NULL,
    INDEX idx_%s_ptype (p_type, v0, v1)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;`
)

// for Postgres
const (
	sqlCreateTablePostgres = `
CREATE TABLE IF NOT EXISTS %s(
    p_type varchar(32)  DEFAULT '' NOT NULL,
    v0     varchar(255) DEFAULT '' NOT NULL,
    v1     varchar(255) DEFAULT '' NOT NULL,
    v2     varchar(255) DEFAULT '' NOT NULL,
    v3     varchar(255) DEFAULT '' NOT NULL,
    v4     varchar(255) DEFAULT '' NOT NULL,
    v5     varchar(255) DEFAULT '' NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_%s_ptype ON %s (p_type, v0, v1);`
)

// for Sqlserver
const (
	sqlCreateTableSqlserver = `
CREATE TABLE %s(
    p_type nvarchar(32)  DEFAULT '' NOT NULL,
    v0     nvarchar(255) DEFAULT '' NOT NULL,
    v1     nvarchar(255) DEFAULT '' NOT NULL,
    v2     nvarchar(255) DEFAULT '' NOT NULL,
    v3     nvarchar(255) DEFAULT '' NOT NULL,
    v4     nvarchar(255) DEFAULT '' NOT NULL,
    v5     nvarchar(255) DEFAULT '' NOT NULL
);
CREATE INDEX idx_%s_ptype ON %s (p_type, v0, v1);`
)

// for Oracle
// const (
// 	sqlCreateTableOracle = `
// CREATE TABLE %s(
// 	"p_type" NVARCHAR2(32)  DEFAULT '',
//     "v0"     NVARCHAR2(255) DEFAULT '',
//     "v1"     NVARCHAR2(255) DEFAULT '',
//     "v2"     NVARCHAR2(255) DEFAULT '',
//     "v3"     NVARCHAR2(255) DEFAULT '',
//     "v4"     NVARCHAR2(255) DEFAULT '',
//     "v5"     NVARCHAR2(255) DEFAULT ''
// )`
// 	sqlCreateIndexOracle = `CREATE INDEX IDX_%s_PTYPE ON %s ("p_type", "v0", "v1")`
// )
