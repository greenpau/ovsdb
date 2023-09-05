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
type OvnACL struct {
	UUID        string `json:"uuid" yaml:"uuid"`
	ExternalIDs map[string]string
}

// GetACL returns a list of OVN ACLs.
func (cli *OvnClient) GetACL() ([]*OvnACL, error) {
	acls := []*OvnACL{}
	// First, get basic information about OVN logical switches.
	query := "SELECT _uuid, external_ids FROM ACL"
	result, err := cli.Database.Northbound.Client.Transact(cli.Database.Northbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Northbound.Name, "ACL", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: no acl found", cli.Database.Northbound.Name)
	}
	for _, row := range result.Rows {
		acl := &OvnACL{}
		if r, dt, err := row.GetColumnValue("_uuid", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			acl.UUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("external_ids", result.Columns); err != nil {
			acl.ExternalIDs = make(map[string]string)
		} else {
			if dt == "map[string]string" {
				acl.ExternalIDs = r.(map[string]string)
			}
		}
		acls = append(acls, acl)
	}
	return acls, nil
}
