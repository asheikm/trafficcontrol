package orttest

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"encoding/json"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/testing/ort-tests/tcdata"
	"testing"
)

type Package struct {
	Name    *string `json:"name"`
	Version *string `json:"version"`
}

func TestTORequester(t *testing.T) {
	fmt.Println("------------- Starting TestTORequester tests ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// chkconfig test
		output, err := runTORequester("atlanta-edge-03", "chkconfig")
		if err != nil {
			t.Fatalf("ERROR: to_requester exec failed: %v\n", err)
		}
		var chkConfig []map[string]interface{}
		err = json.Unmarshal([]byte(output), &chkConfig)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if chkConfig[0]["name"] != "trafficserver" {
			t.Fatal("ERROR unexpected result, expected 'trafficserver'")
		}

		// get system-info test
		output, err = runTORequester("atlanta-edge-03", "system-info")
		if err != nil {
			t.Fatalf("ERROR: to_requester exec failed: %v\n", err)
		}
		var sysInfo map[string]interface{}
		err = json.Unmarshal([]byte(output), &sysInfo)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if sysInfo["tm.instance_name"] != "Traffic Ops ORT Tests" {
			t.Fatalf("ERROR: unexpected 'tm.instance_name'")
		}

		// statuses test
		output, err = runTORequester("atlanta-edge-03", "statuses")
		if err != nil {
			t.Fatalf("ERROR: to_requester exec failed: %v\n", err)
		}
		// should parse json to an array of 'tc.Status'
		var statuses []tc.Status
		err = json.Unmarshal([]byte(output), &statuses)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}

		// packages test
		output, err = runTORequester("atlanta-edge-03", "packages")
		if err != nil {
			t.Fatalf("ERROR: to_requester exec failed: %v\n", err)
		}
		// should parse to an array of 'Package'
		var packages []Package
		err = json.Unmarshal([]byte(output), &packages)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if *packages[0].Name != "trafficserver" {
			t.Fatal("ERROR unexpected result, expected 'trafficserver'")
		}

		// update-status test
		output, err = runTORequester("atlanta-edge-03", "update-status")
		if err != nil {
			t.Fatalf("ERROR: to_requester exec failed: %v\n", err)
		}
		var serverStatus tc.ServerUpdateStatus
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if serverStatus.HostName != "atlanta-edge-03" {
			t.Fatal("ERROR unexpected result, expected 'atlanta-edge-03'")
		}

	})
	fmt.Println("------------- End of TestTORequester tests ---------------")
}