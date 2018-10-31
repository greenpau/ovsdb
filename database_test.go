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
	"testing"
)

func TestListDatabasesMethod(t *testing.T) {
	cli, err := NewClient("unix:/var/run/openvswitch/db.sock", 0)
	if err != nil {
		t.Fatalf("FAIL: %v", err)
	}
	defer cli.Close()
	databases, err := cli.Databases()
	if err != nil {
		t.Fatalf("FAIL: expected to find databases, but received the error: %v", err)
	}
	if len(databases) < 1 {
		t.Fatalf("FAIL: expected to find a single database, but found: %d", len(databases))
	}
	dbName := "Open_vSwitch"
	if databases[0] != dbName {
		t.Fatalf("FAIL: expected to find '%s' database, but found: %s", dbName, databases[0])
	}
	t.Logf("PASS: 'list_dbs' method completed successfully")
}

func TestDatabaseExist(t *testing.T) {
	cli, err := NewClient("unix:/var/run/openvswitch/db.sock", 0)
	if err != nil {
		t.Fatalf("FAIL: %v", err)
	}
	defer cli.Close()
	dbName := "Open_vSwitch"
	if err := cli.DatabaseExists(dbName); err != nil {
		t.Fatalf("FAIL: expected to find '%s' database, but received the error: %v", dbName, err)
	}
	t.Logf("PASS: 'get_schema' method completed successfully")
}
