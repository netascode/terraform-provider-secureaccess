// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package provider

// Section below is generated&owned by "gen/generator.go". //template:begin imports
import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc

func TestAccSecureAccessNetwork(t *testing.T) {
	if os.Getenv("TF_VAR_NETWORKS") == "" {
		t.Skip("skipping test, set environment variable TF_VAR_NETWORKS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network.test", "name", "my_network"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network.test", "status", "OPEN"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network.test", "is_dynamic", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network.test", "prefix_length", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network.test", "prefix", "192.168.1.0"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccSecureAccessNetworkConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessNetworkConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_network.test",
		ImportState:  true,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal

func testAccSecureAccessNetworkConfig_minimum() string {
	config := `resource "secureaccess_network" "test" {` + "\n"
	config += `	name = "my_network"` + "\n"
	config += `	status = "OPEN"` + "\n"
	config += `	is_dynamic = true` + "\n"
	config += `	prefix_length = 30` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessNetworkConfig_all() string {
	config := `resource "secureaccess_network" "test" {` + "\n"
	config += `	name = "my_network"` + "\n"
	config += `	status = "OPEN"` + "\n"
	config += `	is_dynamic = true` + "\n"
	config += `	prefix_length = 30` + "\n"
	config += `	prefix = "192.168.1.0"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll
