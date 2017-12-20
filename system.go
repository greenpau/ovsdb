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
	"bufio"
	"fmt"
	"os"
)

// GetSystemID TODO
func (cli *OvnClient) GetSystemID() error {
	var systemID string
	file, err := os.Open(cli.Database.Vswitch.File.SystemID.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		systemID = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(systemID) != 36 {
		return fmt.Errorf("system-id is not 32 characters in length, but %d", len(systemID))
	}
	cli.System.ID = systemID
	return nil
}

// GetSystemInfo returns a hash containing system information, e.g. `system_id`
// associated with the Open_vSwitch database.
func (cli *OvnClient) GetSystemInfo() error {
	systemInfo := make(map[string]string)
	var systemID string
	file, err := os.Open(cli.Database.Vswitch.File.SystemID.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		systemID = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(systemID) != 36 {
		return fmt.Errorf("system-id is not 32 characters in length, but %d", len(systemID))
	}
	cli.System.ID = systemID
	query := fmt.Sprintf("SELECT ovs_version, db_version, system_type, system_version, external_ids FROM %s", cli.Database.Vswitch.Name)
	result, err := cli.Database.Vswitch.Client.Transact(cli.Database.Vswitch.Name, query)
	if err != nil {
		return fmt.Errorf("The '%s' query failed: %s", query, err)
	}
	if len(result.Rows) == 0 {
		return fmt.Errorf("The '%s' query did not return any rows", query)
	}
	//spew.Dump(result)
	for _, row := range result.Rows {
		col := "external_ids"
		rowData, dataType, err := row.GetColumnValue(col, result.Columns)
		if err != nil {
			return fmt.Errorf("The '%s' query succeeded, but parsing '%s' failed: %s", query, col, err)
		}
		if dataType != "map[string]string" {
			return fmt.Errorf("The '%s' query succeeded, but data type '%s' for '%s' column is unexpected in this context", query, dataType, col)
		}
		systemInfo = rowData.(map[string]string)
		columns := []string{"ovs_version", "db_version", "system_type", "system_version"}
		for _, col := range columns {
			rowData, dataType, err = row.GetColumnValue(col, result.Columns)
			if err != nil {
				return fmt.Errorf("The '%s' query succeeded, but parsing '%s' failed: %s", query, col, err)
			}
			if dataType != "string" {
				return fmt.Errorf("The '%s' query succeeded, but data type '%s' for '%s' column is unexpected in this context", query, dataType, col)
			}
			systemInfo[col] = rowData.(string)
		}
		break
	}
	if dbSystemID, exists := systemInfo["system-id"]; exists {
		if dbSystemID != systemID {
			return fmt.Errorf("The '%s' query succeeded, but found 'system-id' mismatch %s (db) vs. %s (config)", query, dbSystemID, systemID)
		}
	} else {
		return fmt.Errorf("The '%s' query succeeded, but no 'system-id' found", query)
	}

	requiredKeys := []string{"rundir", "hostname", "ovs_version", "db_version", "system_type", "system_version"}
	for _, key := range requiredKeys {
		if _, exists := systemInfo[key]; exists == false {
			return fmt.Errorf("The '%s' query succeeded, but no '%s' found", query, key)
		}
	}

	cli.System.ID = systemInfo["system-id"]
	cli.System.RunDir = systemInfo["rundir"]
	cli.System.Hostname = systemInfo["hostname"]
	cli.System.Type = systemInfo["system_type"]
	cli.System.Version = systemInfo["system_version"]
	cli.Database.Vswitch.Version = systemInfo["ovs_version"]
	cli.Database.Vswitch.Schema.Version = systemInfo["db_version"]
	return nil
}
