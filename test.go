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
	"github.com/casbin/casbin/v2/model"
)

// for test all Adapter methods
type Test struct {
	adapter *Adapter
}

func (p *Test) CreateTable() error {
	return p.adapter.createTable()
}

func (p *Test) TruncateTable() error {
	return p.adapter.truncateTable()
}

func (p *Test) IsTableExist() bool {
	return p.adapter.isTableExist()
}

func (p *Test) InsertRow(line *CasbinRule) error {
	return p.adapter.insertRow(line)
}

func (p *Test) DeleteAll() error {
	return p.adapter.deleteAll()
}

func (p *Test) DeleteRow(line *CasbinRule) error {
	return p.adapter.deleteRow(line)
}

func (p *Test) DeleteByArgs(line *CasbinRule) error {
	return p.adapter.deleteByArgs(line)
}

func (p *Test) TruncateAndInsetRows(lines []CasbinRule) error {
	return p.adapter.truncateAndInsetRows(lines)
}

func (p *Test) SelectAll() (lines []*CasbinRule, err error) {
	return p.adapter.selectAll()
}

func (p *Test) selectWhereIn(filter Filter) (lines []*CasbinRule, err error) {
	return p.adapter.selectWhereIn(filter)
}

func (p *Test) LoadPolicy(model model.Model) error {
	return p.adapter.LoadPolicy(model)
}

func (p *Test) SavePolicy(model model.Model) error {
	return p.adapter.SavePolicy(model)
}

func (p *Test) AddPolicy(sec string, ptype string, rule []string) error {
	return p.adapter.AddPolicy(sec, ptype, rule)
}

func (p *Test) RemovePolicy(sec string, ptype string, rule []string) error {
	return p.adapter.RemovePolicy(sec, ptype, rule)
}

func (p *Test) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return p.adapter.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...)
}

func (p *Test) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	return p.adapter.LoadFilteredPolicy(model, filter)
}

func (p *Test) IsFiltered() bool {
	return p.adapter.IsFiltered()
}

func (p *Test) LoadPolicyLine(line *CasbinRule, model model.Model) {
	p.adapter.loadPolicyLine(line, model)
}

func (p *Test) GenPolicyLine(ptype string, rule []string) CasbinRule {
	return p.adapter.genPolicyLine(ptype, rule)
}
