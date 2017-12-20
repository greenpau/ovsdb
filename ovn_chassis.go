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
)

// OvnChassis represent an OVN chassis.
type OvnChassis struct {
	UUID      string
	Name      string
	IPAddress string
	Encaps    struct {
		UUID  string
		Proto string
	}
	Up int
}

// GetChassis returns a list of OVN chassis.
func (cli *OvnClient) GetChassis() ([]*OvnChassis, error) {
	chassis := []*OvnChassis{}
	// First, get the names and UUIDs of chassis.
	query := "SELECT _uuid, name, encaps FROM Chassis"
	result, err := cli.Database.Southbound.Client.Transact(cli.Database.Southbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Southbound.Name, "Chassis", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: no chassis found", cli.Database.Southbound.Name)
	}
	for _, row := range result.Rows {
		c := &OvnChassis{}
		if r, dt, err := row.GetColumnValue("_uuid", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			c.UUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("name", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			c.Name = r.(string)
		}
		if r, dt, err := row.GetColumnValue("encaps", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			c.Encaps.UUID = r.(string)
		}
		chassis = append(chassis, c)
	}

	// Second, get the IP addresses of the chassis
	query = "SELECT _uuid, chassis_name, ip, type FROM Encap"
	result, err = cli.Database.Southbound.Client.Transact(cli.Database.Southbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Southbound.Name, "Encap", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: no chassis found", cli.Database.Southbound.Name)
	}
	for _, row := range result.Rows {
		var encapUUID string
		var encapProto string
		var chassisName string
		var chassisIPAddress string
		if r, dt, err := row.GetColumnValue("_uuid", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			encapUUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("type", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			encapProto = r.(string)
		}
		if r, dt, err := row.GetColumnValue("chassis_name", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			chassisName = r.(string)
		}
		if r, dt, err := row.GetColumnValue("ip", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			chassisIPAddress = r.(string)
		}
		for _, c := range chassis {
			if c.Encaps.UUID != encapUUID {
				continue
			}
			if c.Name != chassisName {
				continue
			}
			c.IPAddress = chassisIPAddress
			c.Encaps.Proto = encapProto
			break
		}
	}
	return chassis, nil
}
