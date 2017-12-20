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
	"testing"
)

func TestTransactMethod(t *testing.T) {
	sock := "unix:/var/run/openvswitch/db.sock"
	cli, err := NewClient(sock, 0)
	if err != nil {
		t.Fatalf("FAIL: expected to connect to %s, but failed with: %v", sock, err)
	}
	defer cli.Close()
	testsFailed := 0
	for i, test := range []struct {
		db         string
		query      string
		column     string
		shouldFail bool
	}{
		{db: "Open_vSwitch", query: "SELECT * FROM Open_vSwitch", column: "_uuid", shouldFail: false},
		{db: "Open_vSwitch", query: "SELECT _uuid, external_ids FROM Open_vSwitch", column: "_uuid", shouldFail: false},
		{db: "Open_vSwitch", query: "SELECT * FROM Open_vSwitch WHERE db_version==\"7.3.0\"", column: "_uuid", shouldFail: false},
		{db: "Open_vSwitch", query: "SELECT _uuid, db_version FROM Open_vSwitch WHERE db_version==\"7.3.0\"", column: "db_version", shouldFail: false},
		{db: "Open_vSwitch", query: "SELECT _uuid, external_ids FROM Open_vSwitch", column: "external_ids", shouldFail: false},
		{db: "Open_vSwitch", query: "SELECT _uuid, statistics FROM Interface", column: "statistics", shouldFail: false},
		{db: "Open_vSwitch", query: "SELECT _uuid, bridges FROM Open_vSwitch", column: "bridges", shouldFail: false},
		{db: "Open_vSwitch", query: "SELECT _uuid, ports FROM Bridge", column: "ports", shouldFail: false},
	} {
		testFailed := false
		result, err := cli.Transact(test.db, test.query)
		if err != nil && !test.shouldFail {
			t.Logf("FAIL: Test %d: database '%s', query '%s', expected to pass, but failed with: %v", i, test.db, test.query, err)
			testsFailed++
			continue
		}
		for j, row := range result.Rows {
			value, valueType, err := row.GetColumnValue(test.column, result.Columns)
			switch valueType {
			case "string":
				t.Logf("INFO: Test %d: database '%s', query '%s', row: %d, value: %v", i, test.db, test.query, j, value.(string))
			case "uuid":
				t.Logf("INFO: Test %d: database '%s', query '%s', row: %d, value: %v", i, test.db, test.query, j, value.(string))
			case "map[string]string":
				t.Logf("INFO: Test %d: database '%s', query '%s', row: %d, value: %v", i, test.db, test.query, j, value.(map[string]string))
			case "map[string]integer":
				t.Logf("INFO: Test %d: database '%s', query '%s', row: %d, value: %v", i, test.db, test.query, j, value.(map[string]int))
			case "[]string":
				t.Logf("INFO: Test %d: database '%s', query '%s', row: %d, value: %v", i, test.db, test.query, j, value.([]string))
			default:
				t.Logf("FAIL: Test %d: database '%s', query '%s', row: %d, unsupported data type '%s' for value: %v", i, test.db, test.query, j, valueType, value)
				testFailed = true
				testsFailed++
				break
			}
			if err != nil {
				if !test.shouldFail {
					t.Logf("FAIL: Test %d: database '%s', query '%s', expected to pass, but failed with: %v", i, test.db, test.query, err)
					testFailed = true
					testsFailed++
					break
				}
			} else {
				if test.shouldFail {
					t.Logf("FAIL: Test %d: database '%s', query '%s', expected to fail, but passed: %v", i, test.db, test.query, err)
					testFailed = true
					testsFailed++
					break
				}
			}
		}
		if !testFailed {
			if test.shouldFail {
				t.Logf("PASS: Test %d: database '%s', query '%s', expected to fail, failed", i, test.db, test.query)

			} else {
				t.Logf("PASS: Test %d: database '%s', query '%s', expected to pass, passed", i, test.db, test.query)
			}
		}
	}
	if testsFailed > 0 {
		t.Fatalf("Failed %d tests", testsFailed)
	}
}
