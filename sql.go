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
CREATE INDEX idx_casbin_rule_ptype ON %s (p_type, v0, v1);
`
	sqlIsTableExist = "SELECT 1 FROM `%s`;"
	sqlInsertRow    = "INSERT INTO `%s` (`p_type`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (?, ?, ?, ?, ?, ?, ?);"
	sqlDeleteAll    = "DELETE FROM `%s`;"
	sqlDeleteRow    = "DELETE FROM `%s` WHERE `p_type` = :p_type AND `v0` = :v0 AND `v1` = :v1 AND `v2` = :v2 AND `v3` = :v3 AND `v4` = :v4 AND `v5` = :v5;"
	sqlDeleteByArgs = "DELETE FROM `%s` WHERE `p_type` = :p_type"
	sqlSelectAll    = "SELECT * FROM `%s`;"
)

// for Sqlite3
const (
	sqlCreateTableSqlite3 = `
CREATE TABLE IF NOT EXISTS %s
(
    p_type varchar(32),
    v0     varchar(255),
    v1     varchar(255),
    v2     varchar(255),
    v3     varchar(255),
    v4     varchar(255),
    v5     varchar(255)
);
CREATE INDEX idx_casbin_rule_ptype ON %s (p_type, v0, v1);
`
)

// for Mysql
const (
	sqlCreatTableMysql = `
CREATE TABLE %s
(
    p_type varchar(32)  NOT NULL DEFAULT '',
    v0     varchar(255) NOT NULL DEFAULT '',
    v1     varchar(255) NOT NULL DEFAULT '',
    v2     varchar(255) NOT NULL DEFAULT '',
    v3     varchar(255) NOT NULL DEFAULT '',
    v4     varchar(255) NOT NULL DEFAULT '',
    v5     varchar(255) NOT NULL DEFAULT '',
    INDEX idx_casbin_rule_ptype (p_type, v0, v1)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
`
)

// for Postgresql
const (
	sqlCreateTablePostgresql = `
CREATE TABLE %s
(
    p_type varchar(32)  NOT NULL DEFAULT '',
    v0     varchar(255) NOT NULL DEFAULT '',
    v1     varchar(255) NOT NULL DEFAULT '',
    v2     varchar(255) NOT NULL DEFAULT '',
    v3     varchar(255) NOT NULL DEFAULT '',
    v4     varchar(255) NOT NULL DEFAULT '',
    v5     varchar(255) NOT NULL DEFAULT ''
);
CREATE INDEX idx_casbin_rule_ptype ON %s (p_type, v0, v1);
`
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
CREATE INDEX idx_casbin_rule_ptype ON %s (p_type, v0, v1);
`
)

// for Oracle
const (
	sqlCreateTableOracle = `
CREATE TABLE %s
(
    p_type NVARCHAR2(32),
    v0     NVARCHAR2(255),
    v1     NVARCHAR2(255),
    v2     NVARCHAR2(255),
    v3     NVARCHAR2(255),
    v4     NVARCHAR2(255),
    v5     NVARCHAR2(255)
);
CREATE INDEX idx_casbin_rule_ptype ON %s (p_type, v0, v1);
`
)
