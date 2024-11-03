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
type OvnNBGlobal struct {
	UUID        string `json:"uuid" yaml:"uuid"`
	ExternalIDs map[string]string
	HV_CFG 		int64
	HV_CFG_Timestamp int64
	NB_CFG 		int64
	NB_CFG_Timestamp int64
	Options 	map[string]string
	SB_CFG 		int64
	SB_CFG_Timestamp int64
}

// GetACL returns a list of OVN ACLs.
func (cli *OvnClient) GetNBGlobal() (*OvnNBGlobal, error) {
	nbglobal := &OvnNBGlobal{}
	// First, get basic information about OVN logical switches.
	query := "SELECT _uuid,external_ids,hv_cfg,hv_cfg_timestamp,nb_cfg,nb_cfg_timestamp,options,sb_cfg,sb_cfg_timestamp FROM NB_Global"
	result, err := cli.Database.Northbound.Client.Transact(cli.Database.Northbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Northbound.Name, "NB_Global", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: nb global is empty!", cli.Database.Northbound.Name)
	}
	if len(result.Rows) > 1 {
		return nil, fmt.Errorf("%s: nb global has more than 1 Rows!", cli.Database.Northbound.Name)
	}
	for _, row := range result.Rows {
		if r, dt, err := row.GetColumnValue("nb_cfg", result.Columns); err != nil {
			continue
		} else {
			if dt != "int64" {
				// continue
			}
			nbglobal.NB_CFG = r.(int64)
		}
		if r, dt, err := row.GetColumnValue("nb_cfg_timestamp", result.Columns); err != nil {
			continue
		} else {
			if dt != "int64" {
				// continue
			}
			nbglobal.NB_CFG_Timestamp = r.(int64)
		}
		if r, dt, err := row.GetColumnValue("hv_cfg", result.Columns); err != nil {
			continue
		} else {
			if dt != "int64" {
				// continue
			}
			nbglobal.HV_CFG = r.(int64)
		}
		if r, dt, err := row.GetColumnValue("hv_cfg_timestamp", result.Columns); err != nil {
			continue
		} else {
			if dt != "int64" {
				// continue
			}
			nbglobal.HV_CFG_Timestamp = r.(int64)
		}
		if r, dt, err := row.GetColumnValue("sb_cfg", result.Columns); err != nil {
			continue
		} else {
			if dt != "int64" {
				// continue
			}
			nbglobal.SB_CFG = r.(int64)
		}
		if r, dt, err := row.GetColumnValue("sb_cfg_timestamp", result.Columns); err != nil {
			continue
		} else {
			if dt != "int64" {
				// continue
			}
			nbglobal.SB_CFG_Timestamp = r.(int64)
		}
		if r, dt, err := row.GetColumnValue("external_ids", result.Columns); err != nil {
			nbglobal.ExternalIDs = make(map[string]string)
		} else {
			if dt == "map[string]string" {
				nbglobal.ExternalIDs = r.(map[string]string)
			}
		}
		if r, dt, err := row.GetColumnValue("options", result.Columns); err != nil {
			nbglobal.Options = make(map[string]string)
		} else {
			if dt == "map[string]string" {
				nbglobal.Options = r.(map[string]string)
			}
		}
	}
	return nbglobal, nil
}
