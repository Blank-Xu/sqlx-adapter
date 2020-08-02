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
CREATE TABLE %s
(
    p_type varchar(32),
    v0     varchar(255),
    v1     varchar(255),
    v2     varchar(255),
    v3     varchar(255),
    v4     varchar(255),
    v5     varchar(255)
);
CREATE INDEX idx_%s_p_type_v0_v1 ON %s (p_type, v0, v1);
`
	sqlTruncateTable = "TRUNCATE TABLE %s"
	sqlIsTableExist  = "SELECT 1 FROM %s"
	sqlInsertRow     = "INSERT INTO %s (p_type, v0, v1, v2, v3, v4, v5) VALUES (?, ?, ?, ?, ?, ?, ?)"
	sqlDeleteAll     = "DELETE FROM %s"
	sqlDeleteRow     = "DELETE FROM %s WHERE p_type = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ? AND v4 = ? AND v5 = ?"
	sqlDeleteByArgs  = "DELETE FROM %s WHERE p_type = ?"
	sqlSelectAll     = "SELECT * FROM %s"
	sqlSelectWhere   = "SELECT * FROM %s WHERE "
)

// for Sqlite3
const (
	sqlCreateTableSqlite3 = `
CREATE TABLE IF NOT EXISTS %s
(
    p_type varchar(32)  NOT NULL DEFAULT '',
    v0     varchar(255) NOT NULL DEFAULT '',
    v1     varchar(255) NOT NULL DEFAULT '',
    v2     varchar(255) NOT NULL DEFAULT '',
    v3     varchar(255) NOT NULL DEFAULT '',
    v4     varchar(255) NOT NULL DEFAULT '',
    v5     varchar(255) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_%s_p_type_v0_v1 ON %s (p_type, v0, v1);
`
	sqlTruncateTableSqlite3 = "DROP TABLE IF EXISTS %s;" + sqlCreateTableSqlite3
)

// for Mysql
const (
	sqlCreatTableMysql = `
CREATE TABLE IF NOT EXISTS %s
(
    p_type varchar(32)  NOT NULL DEFAULT '',
    v0     varchar(255) NOT NULL DEFAULT '',
    v1     varchar(255) NOT NULL DEFAULT '',
    v2     varchar(255) NOT NULL DEFAULT '',
    v3     varchar(255) NOT NULL DEFAULT '',
    v4     varchar(255) NOT NULL DEFAULT '',
    v5     varchar(255) NOT NULL DEFAULT '',
    INDEX idx_%s_p_type_v0_v1 (p_type, v0, v1)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
`
)

// for Postgres
const (
	sqlCreateTablePostgres = `
CREATE TABLE IF NOT EXISTS %s
(
    p_type varchar(32)  NOT NULL DEFAULT '',
    v0     varchar(255) NOT NULL DEFAULT '',
    v1     varchar(255) NOT NULL DEFAULT '',
    v2     varchar(255) NOT NULL DEFAULT '',
    v3     varchar(255) NOT NULL DEFAULT '',
    v4     varchar(255) NOT NULL DEFAULT '',
    v5     varchar(255) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_%s_p_type_v0_v1 ON %s (p_type, v0, v1);
`
	sqlInsertRowPostgres = "INSERT INTO %s (p_type, v0, v1, v2, v3, v4, v5) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	sqlDeleteRowPostgres = "DELETE FROM %s WHERE p_type = $1 AND v0 = $2 AND v1 = $3 AND v2 = $4 AND v3 = $5 AND v4 = $6 AND v5 = $7"
)

// for Sqlserver
const (
	sqlCreateTableSqlserver = `
CREATE TABLE %s
(
    p_type nvarchar(32)  NOT NULL DEFAULT '',
    v0     nvarchar(255) NOT NULL DEFAULT '',
    v1     nvarchar(255) NOT NULL DEFAULT '',
    v2     nvarchar(255) NOT NULL DEFAULT '',
    v3     nvarchar(255) NOT NULL DEFAULT '',
    v4     nvarchar(255) NOT NULL DEFAULT '',
    v5     nvarchar(255) NOT NULL DEFAULT ''
);
CREATE INDEX idx_%s_p_type_v0_v1 ON %s (p_type, v0, v1);
`
	sqlInsertRowSqlserver = "INSERT INTO %s (p_type, v0, v1, v2, v3, v4, v5) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)"
	sqlDeleteRowSqlserver = "DELETE FROM %s WHERE p_type = @p1 AND v0 = @p2 AND v1 = @p3 AND v2 = @p4 AND v3 = @p5 AND v4 = @p6 AND v5 = @p7"
)

// for Oracle
const (
	sqlCreateTableOracle = `
CREATE TABLE %s
(
    p_type NVARCHAR2(32)  NOT NULL DEFAULT '',
    v0     NVARCHAR2(255) NOT NULL DEFAULT '',
    v1     NVARCHAR2(255) NOT NULL DEFAULT '',
    v2     NVARCHAR2(255) NOT NULL DEFAULT '',
    v3     NVARCHAR2(255) NOT NULL DEFAULT '',
    v4     NVARCHAR2(255) NOT NULL DEFAULT '',
    v5     NVARCHAR2(255) NOT NULL DEFAULT ''
);
CREATE INDEX idx_%s_p_type_v0_v1 ON %s (p_type, v0, v1);
`
	sqlInsertRowOracle = "INSERT INTO %s (p_type, v0, v1, v2, v3, v4, v5) VALUES (:1, :2, :3, :4, :5, :6, :7)"
	sqlDeleteRowOracle = "DELETE FROM %s WHERE p_type = :1 AND v0 = :2 AND v1 = :3 AND v2 = :4 AND v3 = :5 AND v4 = :6 AND v5 = :7"
)
