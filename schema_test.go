// Copyright 2018 Paul Greenberg (greenpau@outlook.com)
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

package ovsdb

import (
	//"github.com/davecgh/go-spew/spew"
	"sort"
	"testing"
)

func TestGetSchemaMethod(t *testing.T) {
	dbs := make(map[string]string)
	dbs["Open_vSwitch"] = "unix:/var/run/openvswitch/db.sock"
	dbs["OVN_Southbound"] = "unix:/run/ovn/ovnsb_db.sock"
	dbs["OVN_Northbound"] = "unix:/run/ovn/ovnnb_db.sock"
	var keys []string
	for k := range dbs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for dbName := range dbs {
		dbSock := dbs[dbName]
		cli, err := NewClient(dbSock, 0)

		if err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		defer cli.Close()
		if err := cli.DatabaseExists(dbName); err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		schema, err := cli.GetSchema(dbName)
		if err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		if schema.Name != dbName {
			t.Fatalf("FAIL: %s (expected) vs. %s (actual)", dbName, schema.Name)
		}
		if len(schema.Tables) == 0 {
			t.Fatalf("FAIL: no tables found in %s", schema.Name)
		}
	}
	t.Logf("PASS: 'get_schema' method completed successfully")
}

func TestSchemaGetTables(t *testing.T) {
	dbs := make(map[string]string)
	dbs["Open_vSwitch"] = "unix:/var/run/openvswitch/db.sock"
	dbs["OVN_Southbound"] = "unix:/run/ovn/ovnsb_db.sock"
	dbs["OVN_Northbound"] = "unix:/run/ovn/ovnnb_db.sock"
	var keys []string
	for k := range dbs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for dbName := range dbs {
		dbSock := dbs[dbName]
		cli, err := NewClient(dbSock, 0)
		if err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		defer cli.Close()
		schema, err := cli.GetSchema(dbName)
		if err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		tables := schema.GetTables()
		if len(tables) > 0 {
			t.Logf("%s Tables:", dbName)
			for _, table := range tables {
				t.Logf("  - %s", table)
			}
		} else {
			t.Fatalf("FAIL: no tables found in %s", dbName)
		}
	}
	t.Logf("PASS: schema.GetTables")
}

func TestSchemaGetColumns(t *testing.T) {
	dbs := make(map[string]string)
	dbs["Open_vSwitch"] = "unix:/var/run/openvswitch/db.sock"
	dbs["OVN_Southbound"] = "unix:/run/ovn/ovnsb_db.sock"
	dbs["OVN_Northbound"] = "unix:/run/ovn/ovnnb_db.sock"
	var keys []string
	for k := range dbs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for dbName := range dbs {
		dbSock := dbs[dbName]
		cli, err := NewClient(dbSock, 0)
		if err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		defer cli.Close()
		schema, err := cli.GetSchema(dbName)
		if err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		tables := schema.GetTables()
		if len(tables) > 0 {
			for _, dbTable := range tables {
				columns := schema.GetColumns(dbTable)
				if len(columns) > 0 {
					t.Logf("The columns in %s table of %s database are:", dbTable, dbName)
					for _, column := range columns {
						t.Logf("  - %s", column)
					}
				} else {
					t.Fatalf("FAIL: no columns found in %s table of %s database", dbTable, dbName)
				}
			}
		} else {
			t.Fatalf("FAIL: no tables found in %s database", dbName)
		}
	}
	t.Logf("PASS: schema.GetColumns")
}

func TestSchemaGetColumnType(t *testing.T) {
	dbName := "Open_vSwitch"
	cli, err := NewClient("unix:/var/run/openvswitch/db.sock", 0)
	if err != nil {
		t.Fatalf("FAIL: %v", err)
	}
	defer cli.Close()
	if err := cli.DatabaseExists(dbName); err != nil {
		t.Fatalf("FAIL: %v", err)
	}
	testFailed := 0
	for i, test := range []struct {
		table      string
		column     string
		columnType string
		shouldFail bool
	}{
		{"Open_vSwitch", "UnknownColumn", "string", true},
		{"Open_vSwitch", "db_version", "string", false},
		{"Open_vSwitch", "next_cfg", "integer", false},
		{"Open_vSwitch", "bridges", "string", true},
		{"Open_vSwitch", "bridges", "map[string]uuid", false},
		{"Bridge", "ports", "map[string]uuid", false},
	} {
		schema, err := cli.GetSchema(dbName)
		if err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		columnType, err := schema.GetColumnType(test.table, test.column)
		if err != nil {
			if !test.shouldFail {
				t.Logf("FAIL: Test %d: table '%s', column '%s', type '%s', expected to pass, but failed with: %s", i, test.table, test.column, test.columnType, err)
				testFailed++
				continue
			}
		}
		if columnType != test.columnType {
			if !test.shouldFail {
				t.Logf("FAIL: Test %d: table '%s', column '%s', expected to pass, but failed: '%s' (expected) vs '%s' (actual)",
					i, test.table, test.column, test.columnType, columnType)
				testFailed++
				continue
			}
		}
		if test.shouldFail && err == nil && columnType == test.columnType {
			t.Logf("FAIL: Test %d: table '%s', column '%s', expected to fail, but passed, because expected and actual type are the same: %s",
				i, test.table, test.column, test.columnType)
			testFailed++
			continue
		}
		t.Logf("PASS: Test %d: table '%s', column '%s', type '%s', passed", i, test.table, test.column, test.columnType)
	}
	if testFailed > 0 {
		t.Fatalf("Failed %d tests", testFailed)
	}
}

func TestSchemaTypesAll(t *testing.T) {
	dbs := make(map[string]string)
	dbs["Open_vSwitch"] = "unix:/var/run/openvswitch/db.sock"
	dbs["OVN_Southbound"] = "unix:/run/ovn/ovnsb_db.sock"
	dbs["OVN_Northbound"] = "unix:/run/ovn/ovnnb_db.sock"
	var keys []string
	for k := range dbs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for dbName := range dbs {
		dbSock := dbs[dbName]
		cli, err := NewClient(dbSock, 0)
		if err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		defer cli.Close()
		if err := cli.DatabaseExists(dbName); err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		schema, err := cli.GetSchema(dbName)
		if err != nil {
			t.Fatalf("FAIL: %v", err)
		}
		tables := schema.GetTables()
		if len(tables) == 0 {
			t.Fatalf("FAIL: no tables found in %s", dbName)
		}
		testFail := 0
		for _, table := range tables {
			columns := schema.GetColumns(table)
			if len(columns) == 0 {
				t.Fatalf("FAIL: no columns found in %s table of %s database", table, dbName)
			}
			for _, column := range columns {
				columnType, err := schema.GetColumnType(table, column)
				if err != nil {
					t.Fatalf("FAIL: %s: %s: %s: %s: %v", dbName, table, column, columnType, err)
				}
				if columnType == "" {
					t.Logf("  - FAIL: %s: %s: %s: UNKNOWN", dbName, table, column)
					testFail++
					continue
				}
				t.Logf("  - %s: %s: %s: %s", dbName, table, column, columnType)
			}
		}
		if testFail > 0 {
			t.Fatalf("FAIL: schema.GetTables")
		}
	}
	t.Logf("PASS: schema.GetTables")
}
