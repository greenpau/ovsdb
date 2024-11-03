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
	"fmt"
	//"github.com/davecgh/go-spew/spew"
)

// OvnLFlow holds Logical_Flow information.
type OvnLFlow struct {
	UUID        string `json:"uuid" yaml:"uuid"`
}

// GetLFlow returns a list of OVN LFlows.
func (cli *OvnClient) GetLFlow() ([]*OvnLFlow, error) {
	lflows := []*OvnLFlow{}
	query := "SELECT _uuid FROM Logical_Flow"
	result, err := cli.Database.Southbound.Client.Transact(cli.Database.Southbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Southbound.Name, "Logical_Flow", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: no lflow found", cli.Database.Southbound.Name)
	}
	for _, row := range result.Rows {
		lflow := &OvnLFlow{}
		if r, dt, err := row.GetColumnValue("_uuid", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			lflow.UUID = r.(string)
		}
		lflows = append(lflows, lflow)
	}
	return lflows, nil
}
