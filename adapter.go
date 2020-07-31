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

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jmoiron/sqlx"
)

const DefaultTableName = "casbin_rule"

type CasbinRule struct {
	PType string `db:"p_type"`
	V0    string `db:"v0"`
	V1    string `db:"v1"`
	V2    string `db:"v2"`
	V3    string `db:"v3"`
	V4    string `db:"v4"`
	V5    string `db:"v5"`
}

type Adapter struct {
	db        *sqlx.DB
	tableName string

	sqlCreateTable   string
	sqlIsTableExist  string
	sqlTruncateTable string
	sqlInsertRow     string
	sqlDeleteAll     string
	sqlDeleteRow     string
	sqlDeleteByArgs  string
	sqlSelectAll     string
}

func NewAdapter(db *sqlx.DB, tableName string) (*Adapter, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	// check db connecting
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	if tableName == "" {
		tableName = DefaultTableName
	}

	adapter := &Adapter{
		db:        db,
		tableName: tableName,
	}

	adapter.genSql()

	if !adapter.isTableExist() {
		if err = adapter.createTable(); err != nil {
			return nil, err
		}
	}

	return adapter, nil
}

func (p *Adapter) genSql() {
	switch p.db.DriverName() {
	case "postgres", "pgx", "pq-timeouts", "cloudsqlpostgres":
		p.sqlCreateTable = fmt.Sprintf(sqlCreateTablePostgresql, p.tableName)
	case "mysql":
		p.sqlCreateTable = fmt.Sprintf(sqlCreatTableMysql, p.tableName)
	case "sqlite3":
		p.sqlCreateTable = fmt.Sprintf(sqlCreateTableSqlite3, p.tableName)
	case "oci8", "ora", "goracle":
		p.sqlCreateTable = fmt.Sprintf(sqlCreateTableOracle, p.tableName)
	case "sqlserver":
		p.sqlCreateTable = fmt.Sprintf(sqlCreateTableSqlserver, p.tableName)
	default:
		p.sqlCreateTable = fmt.Sprintf(sqlCreatTable, p.tableName)
	}

	p.sqlIsTableExist = fmt.Sprintf(sqlIsTableExist, p.tableName)
	p.sqlTruncateTable = fmt.Sprintf(sqlTruncateTable, p.tableName)
	p.sqlInsertRow = fmt.Sprintf(sqlInsertRow, p.tableName)
	p.sqlDeleteAll = fmt.Sprintf(sqlDeleteAll, p.tableName)
	p.sqlDeleteRow = fmt.Sprintf(sqlDeleteRow, p.tableName)
	p.sqlDeleteByArgs = fmt.Sprintf(sqlDeleteByArgs, p.tableName)

	p.sqlSelectAll = fmt.Sprintf(sqlSelectAll, p.tableName)
}

func (p *Adapter) createTable() error {
	_, err := p.db.Exec(p.sqlCreateTable)
	return err
}

func (p *Adapter) isTableExist() bool {
	_, err := p.db.Exec(p.sqlIsTableExist)
	return err == nil
}

func (p *Adapter) truncateTable() error {
	_, err := p.db.Exec(p.sqlTruncateTable)
	return err
}

func (p *Adapter) insertRow(line *CasbinRule) error {
	_, err := p.db.NamedExec(p.sqlInsertRow, line)
	return err
}

func (p *Adapter) deleteAll() error {
	_, err := p.db.Exec(p.sqlDeleteAll)
	return err
}

func (p *Adapter) deleteRow(line *CasbinRule) error {
	_, err := p.db.NamedExec(p.sqlDeleteRow, line)
	return err
}

func (p *Adapter) deleteByArgs(line *CasbinRule) error {
	args := make([]interface{}, 0, 4)
	args = append(args, line.PType)

	sqlBuf := bytes.NewBufferString(p.sqlDeleteByArgs)

	if line.V0 != "" {
		sqlBuf.WriteString(" AND `v0` = ?")
		args = append(args, line.V0)
	}
	if line.V1 != "" {
		sqlBuf.WriteString(" AND `v1` = ?")
		args = append(args, line.V1)
	}
	if line.V2 != "" {
		sqlBuf.WriteString(" AND `v2` = ?")
		args = append(args, line.V2)
	}
	if line.V3 != "" {
		sqlBuf.WriteString(" AND `v3` = ?")
		args = append(args, line.V3)
	}
	if line.V4 != "" {
		sqlBuf.WriteString(" AND `v4` = ?")
		args = append(args, line.V4)
	}
	if line.V5 != "" {
		sqlBuf.WriteString(" AND `v5` = ?")
		args = append(args, line.V5)
	}

	_, err := p.db.Exec(sqlBuf.String(), args...)

	return err
}

func (p *Adapter) truncateAndInsetRows(lines []*CasbinRule) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(p.sqlTruncateTable); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			err = fmt.Errorf("truncate err: %v, rollback err: %v", err, err1)
		}
		return err
	}

	for _, line := range lines {
		if _, err = tx.Exec(p.sqlInsertRow, line); err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				err = fmt.Errorf("insert err: %v, rollback err: %v", err, err1)
			}
			return err
		}
	}

	return tx.Commit()
}

func (p *Adapter) selectAll() (lines []*CasbinRule, err error) {
	lines = make([]*CasbinRule, 0, 50)
	err = p.db.Select(&lines, p.sqlSelectAll)
	return
}

func (p *Adapter) LoadPolicy(model model.Model) error {
	lines, err := p.selectAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

func (p *Adapter) SavePolicy(model model.Model) error {
	lines := make([]*CasbinRule, 0, 32)

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)

		}
	}

	return p.truncateAndInsetRows(lines)
}

func (p *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	return p.insertRow(line)
}

func (p *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	return p.deleteRow(line)
}

func (p *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := CasbinRule{PType: ptype}

	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}

	return p.deleteRow(&line)
}

func loadPolicyLine(line *CasbinRule, model model.Model) {
	var lineBuf bytes.Buffer
	lineBuf.Grow(32)

	lineBuf.WriteString(line.PType)
	if line.V0 != "" {
		lineBuf.WriteByte(',')
		lineBuf.WriteString(line.V0)
	}
	if line.V1 != "" {
		lineBuf.WriteByte(',')
		lineBuf.WriteString(line.V1)
	}
	if line.V2 != "" {
		lineBuf.WriteByte(',')
		lineBuf.WriteString(line.V2)
	}
	if line.V3 != "" {
		lineBuf.WriteByte(',')
		lineBuf.WriteString(line.V3)
	}
	if line.V4 != "" {
		lineBuf.WriteByte(',')
		lineBuf.WriteString(line.V4)
	}
	if line.V5 != "" {
		lineBuf.WriteByte(',')
		lineBuf.WriteString(line.V5)
	}

	persist.LoadPolicyLine(lineBuf.String(), model)
}

func savePolicyLine(ptype string, rule []string) *CasbinRule {
	line := CasbinRule{PType: ptype}

	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return &line
}
