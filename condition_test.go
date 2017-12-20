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

func TestConditionParse(t *testing.T) {
	testFailed := 0
	for i, test := range []struct {
		condition  string
		column     string
		function   string
		value      string
		shouldFail bool
		shouldErr  bool
	}{
		{condition: "db_version==7.3.0", column: "db_version", function: "==", value: "7.4.0", shouldFail: true, shouldErr: false},
		{condition: "db_version==", column: "db_version", function: "==", value: "", shouldFail: true, shouldErr: true},
		{condition: "db_version==7.3.0", column: "db_version", function: "==", value: "7.3.0", shouldFail: false, shouldErr: false},
	} {
		condition, err := NewCondition([]string{test.condition})
		if err != nil {
			if !test.shouldErr {
				t.Logf("FAIL: Test %d: condition '%s', expected to pass, but threw error: %v", i, test.condition, err)
				testFailed++
				continue
			}
		} else {
			if test.shouldErr {
				t.Logf("FAIL: Test %d: condition '%s', expected to throw error, but passed: %v", i, test.condition, condition)
				testFailed++
				continue
			}
		}
		if condition.Column != test.column || condition.Function != test.function || condition.Value != test.value {
			if !test.shouldFail {
				t.Logf("FAIL: Test %d: condition '%s', expected to fail, but passed: %v", i, test.condition, condition)
			}
		}
		if test.shouldFail {
			t.Logf("PASS: Test %d: condition '%s', expected to fail, failed", i, test.condition)
		} else {
			t.Logf("PASS: Test %d: condition '%s', expected to pass, passed", i, test.condition)
		}
	}
	if testFailed > 0 {
		t.Fatalf("Failed %d tests", testFailed)
	}
}
