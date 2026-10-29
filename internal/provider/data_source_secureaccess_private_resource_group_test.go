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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSource

func TestAccDataSourceSecureAccessPrivateResourceGroup(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_private_resource_group.test", "name", "my private resource group"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_private_resource_group.test", "description", "This is my private resource group"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessPrivateResourceGroupPrerequisitesConfig + testAccDataSourceSecureAccessPrivateResourceGroupConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccDataSourceSecureAccessPrivateResourceGroupPrerequisitesConfig + testAccNamedDataSourceSecureAccessPrivateResourceGroupConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccDataSourceSecureAccessPrivateResourceGroupPrerequisitesConfig = `
resource "secureaccess_private_resource" "test" {
  name = "my private resource group"
  endpoint_connection_methods = [
    {
      type                       = "client"
      client_reachable_addresses = ["192.168.1.0/24"]
    }
  ]
  addresses = [
    {
      addresses = ["192.168.1.0/24"]
      protocols_and_ports = [
        {
          protocol = "http/https"
          ports    = "443"
        }
      ]
    }
  ]
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessPrivateResourceGroupConfig() string {
	config := `resource "secureaccess_private_resource_group" "test" {` + "\n"
	config += `	name = "my private resource group"` + "\n"
	config += `	description = "This is my private resource group"` + "\n"
	config += `	private_resource_ids = [secureaccess_private_resource.test.id]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_private_resource_group" "test" {
			id = secureaccess_private_resource_group.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessPrivateResourceGroupConfig() string {
	config := `resource "secureaccess_private_resource_group" "test" {` + "\n"
	config += `	name = "my private resource group"` + "\n"
	config += `	description = "This is my private resource group"` + "\n"
	config += `	private_resource_ids = [secureaccess_private_resource.test.id]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_private_resource_group" "test" {
			name = secureaccess_private_resource_group.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
