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

func TestAccSecureAccessNetworkObject(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network_object.test", "name", "my_network_object"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network_object.test", "description", "This is my network object."))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network_object.test", "type", "network"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network_object.test", "value", "192.168.1.0/24"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccSecureAccessNetworkObjectConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessNetworkObjectConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_network_object.test",
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

func testAccSecureAccessNetworkObjectConfig_minimum() string {
	config := `resource "secureaccess_network_object" "test" {` + "\n"
	config += `	name = "my_network_object"` + "\n"
	config += `	type = "network"` + "\n"
	config += `	value = "192.168.1.0/24"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessNetworkObjectConfig_all() string {
	config := `resource "secureaccess_network_object" "test" {` + "\n"
	config += `	name = "my_network_object"` + "\n"
	config += `	description = "This is my network object."` + "\n"
	config += `	type = "network"` + "\n"
	config += `	value = "192.168.1.0/24"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll
