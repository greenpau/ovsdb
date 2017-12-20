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
	"bytes"
	"encoding/json"
	//	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestNewOperation(t *testing.T) {
	sock := "unix:/var/run/openvswitch/db.sock"
	cli, err := NewClient(sock, 0)
	if err != nil {
		t.Fatalf("FAIL: expected to connect to %s, but failed with: %v", sock, err)
	}
	defer cli.Close()
	testFailed := 0
	for i, test := range []struct {
		query      string
		shouldFail bool
	}{
		{query: "DOSELECT * FROM Open_vSwitch", shouldFail: true},
		{query: "SELECT * TO Open_vSwitch", shouldFail: true},
		{query: "SELECT * FROM Open_vSwitch", shouldFail: false},
		{query: "SELECT db_version, ovs_version FROM Open_vSwitch", shouldFail: false},
		{query: "SELECT db_version, ovs_version FROM Open_vSwitch WHERE db_version==7.3.0", shouldFail: false},
	} {
		tr, err := NewOperation(test.query)
		if err != nil {
			if !test.shouldFail {
				t.Logf("FAIL: Test %d: query '%s', expected to pass, but failed with: %v", i, test.query, err)
				testFailed++
				continue
			}
			t.Logf("PASS: Test %d: query '%s', expected to fail, failed with: %v", i, test.query, err)
			continue
		}
		if test.shouldFail {
			t.Logf("FAIL: Test %d: query '%s', expected to fail, but passed: %v", i, test.query, tr)
			testFailed++
			continue
		}
		t.Logf("PASS: Test %d: query '%s', expected to pass, passed", i, test.query)
	}
	if testFailed > 0 {
		t.Fatalf("Failed %d tests", testFailed)
	}
}

func TestMarshalOperation(t *testing.T) {
	sock := "unix:/var/run/openvswitch/db.sock"
	cli, err := NewClient(sock, 0)
	if err != nil {
		t.Fatalf("FAIL: expected to connect to %s, but failed with: %v", sock, err)
	}
	defer cli.Close()
	testFailed := 0
	for i, test := range []struct {
		query      string
		response   []byte
		shouldFail bool
	}{
		{
			query:      "SELECT * FROM Open_vSwitch",
			response:   []byte(`{"op":"select","table":"Open_vSwitch","where":[]}`),
			shouldFail: false,
		},
		{
			query:      "SELECT db_version, ovs_version FROM Open_vSwitch WHERE db_version==\"7.3.0\"",
			response:   []byte(`{"op":"select","table":"Open_vSwitch","where":[["db_version","==","7.3.0"]],"columns":["db_version","ovs_version"]}`),
			shouldFail: false,
		},
	} {
		op, err := NewOperation(test.query)
		if err != nil {
			t.Logf("FAIL: Test %d: query '%s', expected to create an operation, but failed: %v", i, test.query, err)
			testFailed++
			continue
		}
		response, err := json.Marshal(op)
		if err != nil {
			t.Logf("FAIL: Test %d: query '%s', expected to marshal, but failed: %v", i, test.query, err)
			testFailed++
			continue
		}
		if !bytes.Equal(response, test.response) {
			t.Logf("FAIL: Test %d: query '%s', the expected and actual responses do not match: '%s' vs. '%s'", i, test.query, test.response, response)
			//spew.Dump(response)
			testFailed++
			continue
		}
		t.Logf("PASS: Test %d: query '%s', expected to pass, passed", i, test.query)
	}
	if testFailed > 0 {
		t.Fatalf("Failed %d tests", testFailed)
	}
}
