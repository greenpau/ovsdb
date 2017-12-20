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

func TestNewClient(t *testing.T) {
	testFailed := 0
	for i, test := range []struct {
		socket     string
		shouldFail bool
	}{
		{socket: "", shouldFail: true},
		{socket: "127.0.0.1", shouldFail: true},
		{socket: "127.0.0.1:12", shouldFail: true},
		{socket: "127.0.0.1:123", shouldFail: true},
		{socket: "127.0.0.1:1234", shouldFail: true},
		{socket: "127.0.0.1:98765", shouldFail: true},
		{socket: "unix:/var/run/openvswitch/db.sock", shouldFail: false},
		{socket: "unixd:/var/run/openvswitch/db.sock", shouldFail: true},
	} {
		cli, err := NewClient(test.socket, 0)
		if err != nil {
			if !test.shouldFail {
				t.Logf("FAIL: Test %d: socket '%s', expected to pass, but failed with: %v", i, test.socket, err)
				testFailed++
				continue
			}
			t.Logf("PASS: Test %d: socket '%s', expected to fail: failed with: %v", i, test.socket, err)
			continue
		}
		defer cli.Close()

		if test.shouldFail {
			t.Logf("FAIL: Test %d: socket '%s', expected to fail, but passed", i, test.socket)
			testFailed++
			continue
		}
		t.Logf("PASS: Test %d: socket '%s', expected to pass: passed", i, test.socket)
	}
	if testFailed > 0 {
		t.Fatalf("Failed %d tests", testFailed)
	}
}
