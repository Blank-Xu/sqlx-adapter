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
	"log"
	"runtime"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	"github.com/jmoiron/sqlx"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	rbacModelFile  = "examples/rbac_model.conf"
	rbacPolicyFile = "examples/rbac_policy.csv"
)

var (
	dataSourceNames = map[string]string{
		"sqlite3":  ":memory:",
		"mysql":    "root:@tcp(127.0.0.1:3306)/sqlx_adapter_test",
		"postgres": "user=postgres host=127.0.0.1 port=5432 dbname=sqlx_adapter_test sslmode=disable",
		// "sqlserver": "sqlserver://sa:YourPassword@127.0.0.1:1433?database=sqlx_adapter_test&connection+timeout=30",
	}

	lines = []CasbinRule{
		{PType: "p", V0: "alice", V1: "data1", V2: "read"},
		{PType: "p", V0: "bob", V1: "data2", V2: "read"},
		{PType: "p", V0: "bob", V1: "data2", V2: "write"},
		{PType: "p", V0: "data2_admin", V1: "data1", V2: "read", V3: "test1", V4: "test2", V5: "test3"},
		{PType: "p", V0: "data2_admin", V1: "data2", V2: "write", V3: "test1", V4: "test2", V5: "test3"},
		{PType: "p", V0: "data1_admin", V1: "data2", V2: "write"},
		{PType: "g", V0: "alice", V1: "data2_admin"},
		{PType: "g", V0: "bob", V1: "data2_admin", V2: "test"},
		{PType: "g", V0: "bob", V1: "data1_admin", V2: "test2", V3: "test3", V4: "test4", V5: "test5"},
	}

	filter = Filter{
		PType: []string{"p"},
		V0:    []string{"bob", "data2_admin"},
		V1:    []string{"data1", "data2"},
		V2:    []string{"read", "write"},
		V3:    []string{"test1"},
		V4:    []string{"test2"},
		V5:    []string{"test3"},
	}
)

func finalizer(db *sqlx.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}

func TestAdapters(t *testing.T) {
	for key, value := range dataSourceNames {
		log.Printf("test [%s] start, dataSourceName: [%s]", key, value)

		db, err := sqlx.Connect(key, value)
		if err != nil {
			t.Fatalf("sqlx.Connect failed, err: %v", err)
		}

		// need to control by user, not the package
		runtime.SetFinalizer(db, finalizer)

		testTableName(t, db)

		testSql(t, db, "sqlxadapter_sql")

		testPolicyLine(t, db, "sqlxadapter_policy_line")

		testSaveLoad(t, db, "sqlxadapter_save_load")

		testAutoSave(t, db, "sqlxadapter_auto_save")

		testFilteredPolicy(t, db, "sqlxadapter_filtered_policy")
	}
}

func testTableName(t *testing.T, db *sqlx.DB) {
	_, err := NewAdapter(db, "")
	if err != nil {
		t.Fatalf("NewAdapter failed, err: %v", err)
	}
}

func testSql(t *testing.T, db *sqlx.DB, tableName string) {
	a, err := NewAdapter(db, tableName)
	if err != nil {
		t.Fatalf("NewAdapter failed, err: %v", err)
	}

	// CreateTable test has passed when adapter create
	// if err = a.CreateTable(); err != nil {
	// 	t.Fatal("CreateTable failed, err: ", err)
	// }

	if b := a.isTableExist(); b == false {
		t.Fatal("IsTableExist test failed")
	}

	if err = a.insertRow(lines[0]); err != nil {
		t.Fatal("InsertRow test failed, err: ", err)
	}

	if err = a.truncateAndInsertRows(lines); err != nil {
		t.Fatal("TruncateAndInsertRows test failed, err: ", err)
	}

	if err = a.deleteAll(); err != nil {
		t.Fatal("DeleteAll test failed, err: ", err)
	}

	if err = a.truncateAndInsertRows(lines); err != nil {
		t.Fatal("TruncateAndInsertRows test failed, err: ", err)
	}

	if err = a.deleteByArgs(CasbinRule{}); err != nil && err.Error() != "invalid delete args" {
		t.Fatal("DeleteByArgs test without args failed, err: ", err)
	}

	if err = a.deleteByArgs(lines[8]); err != nil {
		t.Fatal("DeleteByArgs test failed, err: ", err)
	}

	if err = a.truncateAndInsertRows(lines); err != nil {
		t.Fatal("TruncateAndInsertRows test failed, err: ", err)
	}

	records, err := a.selectAll()
	if err != nil {
		t.Fatal("SelectAll test failed, err: ", err)
	}
	for idx, record := range records {
		line := lines[idx]
		if record.PType != line.PType ||
			record.V0 != line.V0 ||
			record.V1 != line.V1 ||
			record.V2 != line.V2 ||
			record.V3 != line.V3 ||
			record.V4 != line.V4 ||
			record.V5 != line.V5 {
			t.Fatalf("SelectAll records test not equal, query record: %+v, need record: %+v", record, line)
		}
	}

	records, err = a.selectWhereIn(filter)
	if err != nil {
		t.Fatal("SelectWhereIn test failed, err: ", err)
	}
	i := 3
	for _, record := range records {
		line := lines[i]
		if record.PType != line.PType ||
			record.V0 != line.V0 ||
			record.V1 != line.V1 ||
			record.V2 != line.V2 ||
			record.V3 != line.V3 ||
			record.V4 != line.V4 ||
			record.V5 != line.V5 {
			t.Fatalf("SelectWhereIn records test not equal, query record: %+v, need record: %+v", record, line)
		}
		i += 1
	}

	if err = a.truncateTable(); err != nil {
		t.Fatal("TruncateTable test failed, err: ", err)
	}
}

func testPolicyLine(t *testing.T, db *sqlx.DB, tableName string) {
	a, err := NewAdapter(db, tableName)
	if err != nil {
		t.Fatalf("NewAdapter failed, err: %v", err)
	}

	testLine := CasbinRule{
		PType: "p",
		V0:    "test0",
		V1:    "test1",
		V2:    "test2",
		V3:    "test3",
		V4:    "test4",
		V5:    "test5",
	}
	rule := []string{"test0", "test1", "test2", "test3", "test4", "test5"}

	line := a.genPolicyLine("p", rule)

	if testLine.PType != line.PType ||
		testLine.V0 != line.V0 ||
		testLine.V1 != line.V1 ||
		testLine.V2 != line.V2 ||
		testLine.V3 != line.V3 ||
		testLine.V4 != line.V4 ||
		testLine.V5 != line.V5 {
		t.Fatalf("GenPolicyLine records test not equal, query record: %+v, need record: %+v", line, testLine)
	}

	e, _ := casbin.NewEnforcer(rbacModelFile, rbacPolicyFile)
	a.loadPolicyLine(&line, e.GetModel())
}

func initPolicy(t *testing.T, db *sqlx.DB, tableName string) {
	// Because the DB is empty at first,
	// so we need to load the policy from the file adapter (.CSV) first.
	e, _ := casbin.NewEnforcer(rbacModelFile, rbacPolicyFile)

	a, err := NewAdapter(db, tableName)
	if err != nil {
		t.Fatal("NewAdapter test failed, err: ", err)
	}

	// This is a trick to save the current policy to the DB.
	// We can't call e.SavePolicy() because the adapter in the enforcer is still the file adapter.
	// The current policy means the policy in the Casbin enforcer (aka in memory).
	err = a.SavePolicy(e.GetModel())
	if err != nil {
		t.Fatal("SavePolicy test failed, err: ", err)
	}

	// Clear the current policy.
	e.ClearPolicy()
	testGetPolicy(t, e, [][]string{})

	// Load the policy from DB.
	err = a.LoadPolicy(e.GetModel())
	if err != nil {
		t.Fatal("LoadPolicy test failed, err: ", err)
	}
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func testSaveLoad(t *testing.T, db *sqlx.DB, tableName string) {
	// Initialize some policy in DB.
	initPolicy(t, db, tableName)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, _ := NewAdapter(db, tableName)
	e, _ := casbin.NewEnforcer(rbacModelFile, a)
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func testAutoSave(t *testing.T, db *sqlx.DB, tableName string) {
	// Initialize some policy in DB.
	initPolicy(t, db, tableName)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, _ := NewAdapter(db, tableName)
	e, _ := casbin.NewEnforcer(rbacModelFile, a)

	// AutoSave is enabled by default.
	// Now we disable it.
	e.EnableAutoSave(false)

	// Because AutoSave is disabled, the policy change only affects the policy in Casbin enforcer,
	// it doesn't affect the policy in the storage.
	e.AddPolicy("alice", "data1", "write")
	// Reload the policy from the storage to see the effect.
	e.LoadPolicy()
	// This is still the original policy.
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Now we enable the AutoSave.
	e.EnableAutoSave(true)

	// Because AutoSave is enabled, the policy change not only affects the policy in Casbin enforcer,
	// but also affects the policy in the storage.
	e.AddPolicy("alice", "data1", "write")
	// Reload the policy from the storage to see the effect.
	e.LoadPolicy()
	// The policy has a new rule: {"alice", "data1", "write"}.
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}, {"alice", "data1", "write"}})

	// Remove the added rule.
	e.RemovePolicy("alice", "data1", "write")
	e.LoadPolicy()
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Remove "data2_admin" related policy rules via a filter.
	// Two rules: {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"} are deleted.
	e.RemoveFilteredPolicy(0, "data2_admin")
	e.LoadPolicy()
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}})
}

func testFilteredPolicy(t *testing.T, db *sqlx.DB, tableName string) {
	// Initialize some policy in DB.
	initPolicy(t, db, tableName)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, _ := NewAdapter(db, tableName)
	e, _ := casbin.NewEnforcer(rbacModelFile, a)
	// Now set the adapter
	e.SetAdapter(a)

	// Load only alice's policies
	e.LoadFilteredPolicy(Filter{V0: []string{"alice"}})
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}})

	// Load only bob's policies
	e.LoadFilteredPolicy(Filter{V0: []string{"bob"}})
	testGetPolicy(t, e, [][]string{{"bob", "data2", "write"}})

	// Load policies for data2_admin
	e.LoadFilteredPolicy(Filter{V0: []string{"data2_admin"}})
	testGetPolicy(t, e, [][]string{{"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Load policies for alice and bob
	e.LoadFilteredPolicy(Filter{V0: []string{"alice", "bob"}})
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}})

	e.AddPolicy("bob", "data1", "write", "test1", "test2", "test3")
	e.LoadFilteredPolicy(filter)
	testGetPolicy(t, e, [][]string{{"bob", "data1", "write", "test1", "test2", "test3"}})
}

func testGetPolicy(t *testing.T, e *casbin.Enforcer, res [][]string) {
	t.Helper()
	myRes := e.GetPolicy()
	t.Log("Policy: ", myRes)

	if !util.Array2DEquals(res, myRes) {
		t.Error("Policy: ", myRes, ", supposed to be ", res)
	}
}
