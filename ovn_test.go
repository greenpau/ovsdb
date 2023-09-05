// Copyright 2020 Paul Greenberg greenpau@outlook.com
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

func TestOvnClientUpdateRefs(t *testing.T) {
	client := NewOvnClient()

	client.Database.Vswitch.Process.ID = 200
	client.Service.Vswitchd.Process.ID = 201
	client.Service.Northd.Process.ID = 202
	client.System.RunDir = "/tmp/random-path"

	client.updateRefs()
	expectedDatabaseCtrl := "unix:/tmp/random-path/ovsdb-server.200.ctl"
	if client.Database.Vswitch.Socket.Control != expectedDatabaseCtrl {
		t.Errorf("UpdateRefs fail. Expected: %s Ctrl: %s", expectedDatabaseCtrl, client.Database.Vswitch.Socket.Control)
	}

	expectedVswitchdCtrl := "unix:/tmp/random-path/ovs-vswitchd.201.ctl"
	if client.Service.Vswitchd.Socket.Control != expectedVswitchdCtrl {
		t.Errorf("UpdateRefs fail. Expected: %s Ctrl: %s", expectedVswitchdCtrl, client.Service.Vswitchd.Socket.Control)
	}

	expectedNorthdCtrl := "unix:/tmp/random-path/ovn-northd.202.ctl"
	if client.Service.Northd.Socket.Control != expectedNorthdCtrl {
		t.Errorf("UpdateRefs fail. Expected: %s Ctrl: %s", expectedNorthdCtrl, client.Service.Vswitchd.Socket.Control)
	}
}
