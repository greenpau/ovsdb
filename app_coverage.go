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
	"strconv"
	"strings"
)

// GetAppCoverageMetrics returns the counters of the the number of times particular events
// occur during a daemon's runtime. The counters include averaged per-second
// rates for the last few seconds, the last minute and the last hour, and the
// total counts of all of the coverage counters.
func (cli *OvnClient) GetAppCoverageMetrics(db string) (map[string]map[string]float64, error) {
	var app Client
	var err error
	cmd := "coverage/show"
	metrics := make(map[string]map[string]float64)
	switch db {
	case "ovsdb-server-northbound":
		app, err = NewClient(cli.Database.Northbound.Socket.Control, cli.Timeout)
	case "ovsdb-server-southbound":
		app, err = NewClient(cli.Database.Southbound.Socket.Control, cli.Timeout)
	default:
		return metrics, fmt.Errorf("The '%s' database is unsupported for '%s'", db, cmd)
	}
	if err != nil {
		app.Close()
		return metrics, fmt.Errorf("failed '%s' from %s: %s", cmd, db, err)
	}
	r, err := app.query(cmd, nil)
	if err != nil {
		app.Close()
		return metrics, fmt.Errorf("the '%s' command failed for %s: %s", cmd, db, err)
	}
	app.Close()
	response := r.String()
	if response == "" {
		return metrics, fmt.Errorf("the '%s' command return no data for %s", cmd, db)
	}
	lines := strings.Split(response, "\\n")
	for _, line := range lines {
		lineArray := strings.Split(strings.Join(strings.Fields(line), " "), " ")
		if len(lineArray) != 6 {
			continue
		}
		name := lineArray[0]
		if _, exists := metrics[name]; !exists {
			metrics[name] = make(map[string]float64)
		}
		if interval, err := strconv.ParseFloat(strings.TrimSuffix(lineArray[1], "/sec"), 64); err == nil {
			metrics[name]["5s"] = interval
		}
		if interval, err := strconv.ParseFloat(strings.TrimSuffix(lineArray[2], "/sec"), 64); err == nil {
			metrics[name]["5m"] = interval
		}
		if interval, err := strconv.ParseFloat(strings.TrimSuffix(lineArray[3], "/sec"), 64); err == nil {
			metrics[name]["1h"] = interval
		}
		if total, err := strconv.Atoi(lineArray[5]); err == nil {
			metrics[name]["total"] = float64(total)
		}
	}
	return metrics, nil
}
