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

// OvnACL holds ACL information.
type OvnSBGlobal struct {
	UUID        string `json:"uuid" yaml:"uuid"`
	ExternalIDs map[string]string
	NB_CFG 		int64
	Options 	map[string]string
}

// GetACL returns a list of OVN ACLs.
func (cli *OvnClient) GetSBGlobal() (*OvnSBGlobal, error) {
	sbglobal := &OvnSBGlobal{}
	// First, get basic information about OVN logical switches.
	query := "SELECT _uuid,external_ids,nb_cfg,options FROM SB_Global"
	result, err := cli.Database.Southbound.Client.Transact(cli.Database.Southbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Southbound.Name, "SB_Global", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: sb global is empty!", cli.Database.Southbound.Name)
	}
	if len(result.Rows) > 1 {
		return nil, fmt.Errorf("%s: sb global has more than 1 Rows!", cli.Database.Southbound.Name)
	}
	for _, row := range result.Rows {
		if r, dt, err := row.GetColumnValue("nb_cfg", result.Columns); err != nil {
			continue
		} else {
			if dt != "int64" {
				// continue
			}
			sbglobal.NB_CFG = r.(int64)
		}
		if r, dt, err := row.GetColumnValue("external_ids", result.Columns); err != nil {
			sbglobal.ExternalIDs = make(map[string]string)
		} else {
			if dt == "map[string]string" {
				sbglobal.ExternalIDs = r.(map[string]string)
			}
		}
		if r, dt, err := row.GetColumnValue("options", result.Columns); err != nil {
			sbglobal.Options = make(map[string]string)
		} else {
			if dt == "map[string]string" {
				sbglobal.Options = r.(map[string]string)
			}
		}
	}
	return sbglobal, nil
}
