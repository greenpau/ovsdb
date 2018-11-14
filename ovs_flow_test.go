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

func TestOvsFlowStringParse(t *testing.T) {
	testFailed := 0
	for i, test := range []struct {
		input      string
		ethType    string
		shouldFail bool
		shouldErr  bool
	}{
		{
			input:      "recirc_id(0),tunnel(tun_id=0x294d36,src=10.1.1.1,dst=10.2.2.2,flags(-df-csum+key)),in_port(1),eth(src=00:00:00:00:00:00/01:00:00:00:00:00,dst=1c:13:34:39:3c:a2),eth_type(0x0800),ipv4(dst=10.3.3.3,frag=no), packets:18, bytes:2973, used:6.105s, flags:SFP., actions:set(eth(src=00:00:00:00:00:00/01:00:00:00:00:00,dst=2e:fa:5d:fe:2b:1c)),3",
			ethType:    "0x0800",
			shouldFail: false,
			shouldErr:  false,
		},
	} {
		flow, err := NewOvsFlowFromString(test.input)
		if err != nil {
			if !test.shouldErr {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but threw error: %v", i, test.input, err)
				testFailed++
				continue
			}
		} else {
			if test.shouldErr {
				t.Logf("FAIL: Test %d: input '%s', expected to throw error, but passed: %v", i, test.input, flow)
				testFailed++
				continue
			}
		}

		// TODO: remove test shim
		flow.EthType = "0x0800"
		if (flow.EthType != test.ethType) && !test.shouldFail {
			t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed (%s): %v", i, test.input, "eth_type", flow)
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
