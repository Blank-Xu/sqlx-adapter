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

// Test  for test Adapter private methods
type Test struct {
	Adapter *Adapter
}

// func (p *Test) CreateTable() error {
// 	return p.Adapter.createTable()
// }

func (p *Test) TruncateTable() error {
	return p.Adapter.truncateTable()
}

func (p *Test) IsTableExist() bool {
	return p.Adapter.isTableExist()
}

func (p *Test) InsertRow(line CasbinRule) error {
	return p.Adapter.insertRow(line)
}

func (p *Test) DeleteAll() error {
	return p.Adapter.deleteAll()
}

func (p *Test) DeleteByArgs(line CasbinRule) error {
	return p.Adapter.deleteByArgs(line)
}

func (p *Test) TruncateAndInsertRows(lines []CasbinRule) error {
	return p.Adapter.truncateAndInsertRows(lines)
}

func (p *Test) SelectAll() (lines []*CasbinRule, err error) {
	return p.Adapter.selectAll()
}

func (p *Test) SelectWhereIn(filter Filter) (lines []*CasbinRule, err error) {
	return p.Adapter.selectWhereIn(filter)
}

func (p *Test) LoadPolicyLine(line *CasbinRule, model model.Model) {
	p.Adapter.loadPolicyLine(line, model)
}

func (p *Test) GenPolicyLine(ptype string, rule []string) CasbinRule {
	return p.Adapter.genPolicyLine(ptype, rule)
}
