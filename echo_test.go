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

func TestEchoMethod(t *testing.T) {
	cli, err := NewClient("unix:/var/run/openvswitch/db.sock", 0)
	if err != nil {
		t.Fatalf("FAIL: %v", err)
	}
	defer cli.Close()
	if err := cli.Echo("test message"); err != nil {
		t.Fatalf("FAIL: %v", err)
	}
	t.Logf("PASS: 'echo' method completed successfully")
}
