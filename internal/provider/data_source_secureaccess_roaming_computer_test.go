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

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSource

func TestAccDataSourceSecureAccessRoamingComputer(t *testing.T) {
	if os.Getenv("TF_VAR_roaming_computer_name") == "" {
		t.Skip("skipping test, set environment variable TF_VAR_roaming_computer_name")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "device_id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "anyconnect_device_id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "origin_id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "applied_bundle"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "type"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "os_version"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "os_version_name"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "has_ip_blocking"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "status"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "secure_web_gateway_status"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_roaming_computer.test", "secure_client_version"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessRoamingComputerPrerequisitesConfig + testAccDataSourceSecureAccessRoamingComputerConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccDataSourceSecureAccessRoamingComputerPrerequisitesConfig + testAccNamedDataSourceSecureAccessRoamingComputerConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccDataSourceSecureAccessRoamingComputerPrerequisitesConfig = `
variable "roaming_computer_name" { default = null } // tests will set $TF_VAR_ROAMING_COMPUTER_NAME
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessRoamingComputerConfig() string {
	config := `resource "secureaccess_roaming_computer" "test" {` + "\n"
	config += `	name = var.roaming_computer_name` + "\n"
	config += `	delete_on_destroy = false` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_roaming_computer" "test" {
			id = secureaccess_roaming_computer.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessRoamingComputerConfig() string {
	config := `resource "secureaccess_roaming_computer" "test" {` + "\n"
	config += `	name = var.roaming_computer_name` + "\n"
	config += `	delete_on_destroy = false` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_roaming_computer" "test" {
			name = secureaccess_roaming_computer.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
