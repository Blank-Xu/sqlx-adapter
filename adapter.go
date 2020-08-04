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

	isFiltered bool

	sqlCreateTable   string
	sqlTruncateTable string
	sqlIsTableExist  string
	sqlInsertRow     string
	sqlDeleteAll     string
	sqlDeleteByArgs  string
	sqlSelectAll     string
	sqlSelectWhere   string
}

type Filter struct {
	PType []string
	V0    []string
	V1    []string
	V2    []string
	V3    []string
	V4    []string
	V5    []string
}

// NewAdapter is the constructor for Adapter.
// db should connected to database and controlled by user.
// If tableName == "", the Adapter will automatically create a table named "casbin_rule".
func NewAdapter(db *sqlx.DB, tableName string) (*Adapter, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	// check db connecting
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	switch db.DriverName() {
	case "oci8", "ora", "goracle":
		return nil, errors.New("sqlxadapter: please checkout 'oracle' branch")
	}

	if tableName == "" {
		tableName = DefaultTableName
	}

	adapter := Adapter{
		db:        db,
		tableName: tableName,
	}

	// prepare different database sql
	adapter.genSql()

	if !adapter.isTableExist() {
		if err = adapter.createTable(); err != nil {
			return nil, err
		}
	}

	return &adapter, nil
}

// genSql  generate sql based on db driver name.
func (p *Adapter) genSql() {
	p.sqlCreateTable = fmt.Sprintf(sqlCreatTable, p.tableName, p.tableName, p.tableName)
	p.sqlTruncateTable = fmt.Sprintf(sqlTruncateTable, p.tableName)
	p.sqlIsTableExist = fmt.Sprintf(sqlIsTableExist, p.tableName)

	p.sqlInsertRow = fmt.Sprintf(sqlInsertRow, p.tableName)
	p.sqlDeleteAll = fmt.Sprintf(sqlDeleteAll, p.tableName)
	p.sqlDeleteByArgs = fmt.Sprintf(sqlDeleteByArgs, p.tableName)

	p.sqlSelectAll = fmt.Sprintf(sqlSelectAll, p.tableName)
	p.sqlSelectWhere = fmt.Sprintf(sqlSelectWhere, p.tableName)

	switch p.db.DriverName() {
	case "postgres", "pgx", "pq-timeouts", "cloudsqlpostgres":
		p.sqlCreateTable = fmt.Sprintf(sqlCreateTablePostgres, p.tableName, p.tableName, p.tableName)
		p.sqlInsertRow = fmt.Sprintf(sqlInsertRowPostgres, p.tableName)
	case "mysql":
		p.sqlCreateTable = fmt.Sprintf(sqlCreatTableMysql, p.tableName, p.tableName)
	case "sqlite3":
		p.sqlCreateTable = fmt.Sprintf(sqlCreateTableSqlite3, p.tableName, p.tableName, p.tableName)
		p.sqlTruncateTable = fmt.Sprintf(sqlTruncateTableSqlite3, p.tableName, p.tableName, p.tableName, p.tableName)
	case "sqlserver":
		p.sqlCreateTable = fmt.Sprintf(sqlCreateTableSqlserver, p.tableName, p.tableName, p.tableName)
		p.sqlInsertRow = fmt.Sprintf(sqlInsertRowSqlserver, p.tableName)
	}
}

// createTable  create a not exists table.
func (p *Adapter) createTable() error {
	_, err := p.db.Exec(p.sqlCreateTable)

	return err
}

// truncateTable  clear the table.
func (p *Adapter) truncateTable() error {
	_, err := p.db.Exec(p.sqlTruncateTable)

	return err
}

// isTableExist  check the table exists.
func (p *Adapter) isTableExist() bool {
	_, err := p.db.Exec(p.sqlIsTableExist)

	return err == nil
}

// insertRow  insert one row to table.
func (p *Adapter) insertRow(line CasbinRule) error {
	// if line.PType == "" {
	// 	return errors.New("invalid params")
	// }

	_, err := p.db.Exec(p.sqlInsertRow, line.PType, line.V0, line.V1, line.V2, line.V3, line.V4, line.V5)

	return err
}

// deleteAll  clear the table.
func (p *Adapter) deleteAll() error {
	_, err := p.db.Exec(p.sqlDeleteAll)

	return err
}

// deleteByArgs  delete eligible data.
func (p *Adapter) deleteByArgs(line CasbinRule) error {
	var sqlBuf bytes.Buffer

	sqlBuf.Grow(32)
	sqlBuf.WriteString(p.sqlDeleteByArgs)

	args := make([]interface{}, 0, 4)
	args = append(args, line.PType)

	if line.V0 != "" {
		sqlBuf.WriteString(" AND v0 = ?")
		args = append(args, line.V0)
	}
	if line.V1 != "" {
		sqlBuf.WriteString(" AND v1 = ?")
		args = append(args, line.V1)
	}
	if line.V2 != "" {
		sqlBuf.WriteString(" AND v2 = ?")
		args = append(args, line.V2)
	}
	if line.V3 != "" {
		sqlBuf.WriteString(" AND v3 = ?")
		args = append(args, line.V3)
	}
	if line.V4 != "" {
		sqlBuf.WriteString(" AND v4 = ?")
		args = append(args, line.V4)
	}
	if line.V5 != "" {
		sqlBuf.WriteString(" AND v5 = ?")
		args = append(args, line.V5)
	}

	query := sqlBuf.String()
	if query == p.sqlDeleteByArgs {
		return errors.New("invalid delete args")
	}

	switch p.db.DriverName() {
	case "sqlite3", "mysql":

	default:
		query = p.db.Rebind(query)
	}

	_, err := p.db.Exec(query, args...)

	return err
}

// truncateAndInsertRows  clear table and insert new rows.
func (p *Adapter) truncateAndInsertRows(lines []CasbinRule) error {
	err := p.truncateTable()
	if err != nil {
		return err
	}

	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	// if _, err = tx.Exec(p.sqlDeleteAll); err != nil {
	// 	if err1 := tx.Rollback(); err1 != nil {
	// 		err = fmt.Errorf("delete err: %v, rollback err: %v", err, err1)
	// 	}
	// 	return err
	// }

	stmt, err := tx.Preparex(p.sqlInsertRow)
	if err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			err = fmt.Errorf("preparex err: %v, rollback err: %v", err, err1)
		}
		return err
	}

	for _, line := range lines {
		if _, err = stmt.Exec(line.PType, line.V0, line.V1, line.V2, line.V3, line.V4, line.V5); err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				err = fmt.Errorf("exec err: %v, rollback err: %v", err, err1)
			}

			return err
		}
	}

	if err = stmt.Close(); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			err = fmt.Errorf("stmt close err: %v, rollback err: %v", err, err1)
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			err = fmt.Errorf("commit err: %v, rollback err: %v", err, err1)
		}
	}

	return err
}

// selectAll  select all data from the table.
func (p *Adapter) selectAll() ([]*CasbinRule, error) {
	// make a slice with capacity
	lines := make([]*CasbinRule, 0, 32)

	return lines, p.db.Select(&lines, p.sqlSelectAll)
}

// selectWhereIn  select for eligible data from the table.
func (p *Adapter) selectWhereIn(filter Filter) (lines []*CasbinRule, err error) {
	var sqlBuf bytes.Buffer

	sqlBuf.Grow(32)
	sqlBuf.WriteString(p.sqlSelectWhere)

	params := make([]interface{}, 0, 4)
	if len(filter.PType) > 0 {
		// if sqlBuf.Bytes()[sqlBuf.Len()-1] == ')' {
		// 	sqlBuf.WriteString(" AND ")
		// }
		sqlBuf.WriteString("p_type IN (?)")
		params = append(params, filter.PType)
	}
	if len(filter.V0) > 0 {
		if sqlBuf.Bytes()[sqlBuf.Len()-1] == ')' {
			sqlBuf.WriteString(" AND ")
		}
		sqlBuf.WriteString("v0 IN (?)")
		params = append(params, filter.V0)
	}
	if len(filter.V1) > 0 {
		if sqlBuf.Bytes()[sqlBuf.Len()-1] == ')' {
			sqlBuf.WriteString(" AND ")
		}
		sqlBuf.WriteString("v1 IN (?)")
		params = append(params, filter.V1)
	}
	if len(filter.V2) > 0 {
		if sqlBuf.Bytes()[sqlBuf.Len()-1] == ')' {
			sqlBuf.WriteString(" AND ")
		}
		sqlBuf.WriteString("v2 IN (?)")
		params = append(params, filter.V2)
	}
	if len(filter.V3) > 0 {
		if sqlBuf.Bytes()[sqlBuf.Len()-1] == ')' {
			sqlBuf.WriteString(" AND ")
		}
		sqlBuf.WriteString("v3 IN (?)")
		params = append(params, filter.V3)
	}
	if len(filter.V4) > 0 {
		if sqlBuf.Bytes()[sqlBuf.Len()-1] == ')' {
			sqlBuf.WriteString(" AND ")
		}
		sqlBuf.WriteString("v4 IN (?)")
		params = append(params, filter.V4)
	}
	if len(filter.V5) > 0 {
		if sqlBuf.Bytes()[sqlBuf.Len()-1] == ')' {
			sqlBuf.WriteString(" AND ")
		}
		sqlBuf.WriteString("v5 IN (?)")
		params = append(params, filter.V5)
	}

	query, args, err := sqlx.In(sqlBuf.String(), params...)
	if err != nil {
		return
	}

	switch p.db.DriverName() {
	case "sqlite3", "mysql":

	default:
		query = p.db.Rebind(query)
	}

	lines = make([]*CasbinRule, 0, 32)
	err = p.db.Select(&lines, query, args...)

	return
}

// LoadPolicy  load all policy rules from the storage.
func (p *Adapter) LoadPolicy(model model.Model) error {
	lines, err := p.selectAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		p.loadPolicyLine(line, model)
	}

	return nil
}

// SavePolicy  save policy rules to the storage.
func (p *Adapter) SavePolicy(model model.Model) error {
	lines := make([]CasbinRule, 0, 32)

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := p.genPolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := p.genPolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	return p.truncateAndInsertRows(lines)
}

// AddPolicy  add one policy rule to the storage.
func (p *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := p.genPolicyLine(ptype, rule)

	return p.insertRow(line)
}

// RemovePolicy  remove policy rules from the storage.
func (p *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := p.genPolicyLine(ptype, rule)

	return p.deleteByArgs(line)
}

// RemoveFilteredPolicy  remove policy rules that match the filter from the storage.
func (p *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := CasbinRule{
		PType: ptype,
	}

	l := fieldIndex + len(fieldValues)
	if fieldIndex <= 0 && 0 < l {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < l {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < l {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < l {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < l {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < l {
		line.V5 = fieldValues[5-fieldIndex]
	}

	return p.deleteByArgs(line)
}

// LoadFilteredPolicy  load policy rules that match the filter.
func (p *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	filterValue, ok := filter.(Filter)
	if !ok {
		return errors.New("invalid filter type")
	}

	lines, err := p.selectWhereIn(filterValue)
	if err != nil {
		return err
	}

	for _, line := range lines {
		p.loadPolicyLine(line, model)
	}

	p.isFiltered = true

	return nil
}

// IsFiltered  returns true if the loaded policy rules has been filtered.
func (p *Adapter) IsFiltered() bool {
	return p.isFiltered
}

func (Adapter) loadPolicyLine(line *CasbinRule, model model.Model) {
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

func (Adapter) genPolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{
		PType: ptype,
	}

	l := len(rule)
	if l > 0 {
		line.V0 = rule[0]
	}
	if l > 1 {
		line.V1 = rule[1]
	}
	if l > 2 {
		line.V2 = rule[2]
	}
	if l > 3 {
		line.V3 = rule[3]
	}
	if l > 4 {
		line.V4 = rule[4]
	}
	if l > 5 {
		line.V5 = rule[5]
	}

	return line
}
