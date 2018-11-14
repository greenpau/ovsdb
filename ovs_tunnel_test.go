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

func TestOvsTunnelStringParse(t *testing.T) {
	testFailed := 0
	for i, test := range []struct {
		input         string
		id            uint64
		name          string
		encapsulation string
		localIP       string
		remoteIP      string
		key           string
		packetType    string
		ttl           uint64
		checksum      bool
		shouldFail    bool
		shouldErr     bool
	}{
		{
			input:         "port 3: ovn-450e9e-0 (vxlan: ::->10.77.88.11, key=flow, legacy_l2, dp port=3, ttl=64, csum=true)",
			id:            3,
			name:          "ovn-450e9e-0",
			encapsulation: "vxlan",
			localIP:       "0.0.0.0",
			remoteIP:      "10.77.88.11",
			key:           "flow",
			packetType:    "legacy_l2",
			ttl:           64,
			checksum:      true,
			shouldFail:    false,
			shouldErr:     false,
		},
		{
			input:         "port 2: ovn-0b9dcb-0 (geneve: ::->10.77.90.10, key=flow, legacy_l2, dp port=2, ttl=64, csum=true)",
			id:            2,
			name:          "ovn-0b9dcb-0",
			encapsulation: "geneve",
			localIP:       "0.0.0.0",
			remoteIP:      "10.77.90.10",
			key:           "flow",
			packetType:    "legacy_l2",
			ttl:           64,
			checksum:      true,
			shouldFail:    false,
			shouldErr:     false,
		},
		{
			input:         "port 3: ovn-e08372-1 (vxlan: ::->10.77.88.12, key=flow, legacy_l2, dp port=3, ttl=64, csum=true)",
			id:            3,
			name:          "ovn-e08372-1",
			encapsulation: "vxlan",
			localIP:       "0.0.0.0",
			remoteIP:      "10.77.88.12",
			key:           "flow",
			packetType:    "legacy_l2",
			ttl:           64,
			checksum:      true,
			shouldFail:    false,
			shouldErr:     false,
		},
	} {
		tunnel, err := NewOvsTunnelFromString(test.input)
		if err != nil {
			if !test.shouldErr {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but threw error: %v", i, test.input, err)
				testFailed++
				continue
			}
		} else {
			if test.shouldErr {
				t.Logf("FAIL: Test %d: input '%s', expected to throw error, but passed: %v", i, test.input, *tunnel)
				testFailed++
				continue
			}
		}

		if (tunnel.Encapsulation != test.encapsulation) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "encapsulation", *tunnel)
			testFailed++
			continue
		}

		if (tunnel.Name != test.name) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "name", *tunnel)
			testFailed++
			continue
		}

		if (tunnel.ID != test.id) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "id", *tunnel)
			testFailed++
			continue
		}

		if (tunnel.LocalIP != test.localIP) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "localIP", *tunnel)
			testFailed++
			continue
		}

		if (tunnel.RemoteIP != test.remoteIP) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "remoteIP", *tunnel)
			testFailed++
			continue
		}

		if (tunnel.Key != test.key) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "key", *tunnel)
			testFailed++
			continue
		}

		if (tunnel.PacketType != test.packetType) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "packetType", *tunnel)
			testFailed++
			continue
		}

		if (tunnel.TTL != test.ttl) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "ttl", *tunnel)
			testFailed++
			continue
		}

		if (tunnel.Checksum != test.checksum) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "checksum", *tunnel)
			testFailed++
			continue
		}

		if test.shouldFail {
			t.Logf("PASS: Test %d: input '%s', expected to fail, failed", i, test.input)
		} else {
			t.Logf("PASS: Test %d: input '%s', expected to pass, passed", i, test.input)
		}
	}
	if testFailed > 0 {
		t.Fatalf("Failed %d tests", testFailed)
	}
}
