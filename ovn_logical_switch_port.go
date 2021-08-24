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
	"net"
	"strings"
)

type OvnLogicalSwitchPortAddress struct {
	MacAddress  net.HardwareAddr
	IpAddresses []net.IP
	Dynamic     bool
	Router      bool
	Unknown     bool
}

// OvnLogicalSwitchPort holds a consolidated record from both NB and SB
// databases about a logical switch port and the workload attached to it.
type OvnLogicalSwitchPort struct {
	UUID              string
	Name              string
	Addresses         []OvnLogicalSwitchPortAddress
	ExternalIDs       map[string]string
	Encapsulation     string
	TunnelKey         uint64
	Up                bool
	PortBindingUUID   string
	ChassisUUID       string
	ChassisIPAddress  net.IP
	DatapathUUID      string
	LogicalSwitchUUID string
	LogicalSwitchName string
}

func parseLogicalPortAddress(s string) OvnLogicalSwitchPortAddress {
	addrs := strings.Split(s, " ")

	// Check if address line is "router"
	if addrs[0] == "router" {
		return OvnLogicalSwitchPortAddress{Router: true}
	}

	// Check if address line is "unknown"
	if addrs[0] == "unknown" {
		return OvnLogicalSwitchPortAddress{Unknown: true}
	}

	// Check if address line is "dynamic"
	if addrs[0] == "dynamic" {
		portAddress := OvnLogicalSwitchPortAddress{Unknown: true}
		// Sometimes "dynamic" may be followed by an IP address(es)
		if len(addrs) > 1 {
			for _, v := range addrs[1:] {
				portAddress.IpAddresses = append(portAddress.IpAddresses, net.ParseIP(v))
			}
		}
		return portAddress
	}

	// In all other cases first entry should be a MAC address
	macAddr, _ := net.ParseMAC(addrs[0])
	portAddress := OvnLogicalSwitchPortAddress{MacAddress: macAddr}

	// Process any IP addresses
	if len(addrs) > 1 {
		if addrs[1] == "dynamic" {
			portAddress.Dynamic = true
		} else {
			for _, v := range addrs[1:] {
				portAddress.IpAddresses = append(portAddress.IpAddresses, net.ParseIP(v))
			}
		}
	}

	return portAddress
}

// GetLogicalSwitchPorts returns a list of OVN logical switch ports.
func (cli *OvnClient) GetLogicalSwitchPorts() ([]*OvnLogicalSwitchPort, error) {
	// First, fetch logical switch ports.
	ports := []*OvnLogicalSwitchPort{}
	query := "SELECT _uuid, addresses, external_ids, name, up FROM Logical_Switch_Port"
	result, err := cli.Database.Northbound.Client.Transact(cli.Database.Northbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Northbound.Name, "Logical_Switch_Port", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: no logical switch port found", cli.Database.Northbound.Name)
	}
	for _, row := range result.Rows {
		port := OvnLogicalSwitchPort{}
		if r, dt, err := row.GetColumnValue("_uuid", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			port.UUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("name", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			port.Name = r.(string)
		}
		if r, dt, err := row.GetColumnValue("up", result.Columns); err == nil {
			if dt == "bool" {
				port.Up = r.(bool)
			}
		}
		if r, dt, err := row.GetColumnValue("external_ids", result.Columns); err == nil {
			if dt == "map[string]string" {
				port.ExternalIDs = r.(map[string]string)
			}
		} else {
			port.ExternalIDs = make(map[string]string)
		}
		if r, dt, err := row.GetColumnValue("addresses", result.Columns); err == nil {
			switch dt {
			case "string":
				port.Addresses = append(port.Addresses, parseLogicalPortAddress(r.(string)))
			case "[]string":
				for _, s := range r.([]string) {
					port.Addresses = append(port.Addresses, parseLogicalPortAddress(s))
				}
			}
		}
		ports = append(ports, &port)
	}

	// Next, gather tunnel ids and other details about the logical ports.
	query = "SELECT _uuid, chassis, datapath, logical_port, tunnel_key FROM Port_Binding"
	result, err = cli.Database.Southbound.Client.Transact(cli.Database.Southbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Southbound.Name, "Port_Binding", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: no port binding found", cli.Database.Southbound.Name)
	}
	for _, row := range result.Rows {
		var portBindingUUID string
		var portBindingChassisUUID string
		var portBindingDatapathUUID string
		var portBindingLogicalPortName string
		var portBindingTunnelKey uint64
		if r, dt, err := row.GetColumnValue("_uuid", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			portBindingUUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("chassis", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			portBindingChassisUUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("datapath", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			portBindingDatapathUUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("logical_port", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			portBindingLogicalPortName = r.(string)
		}
		if r, dt, err := row.GetColumnValue("tunnel_key", result.Columns); err != nil {
			continue
		} else {
			if dt != "integer" {
				continue
			}
			portBindingTunnelKey = uint64(r.(int64))
		}
		for _, port := range ports {
			if port.Name == portBindingLogicalPortName {
				port.PortBindingUUID = portBindingUUID
				port.ChassisUUID = portBindingChassisUUID
				port.DatapathUUID = portBindingDatapathUUID
				port.TunnelKey = portBindingTunnelKey
				break
			}
		}
	}
	return ports, nil
}
