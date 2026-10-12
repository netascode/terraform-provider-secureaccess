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

func TestAccDataSourceSecureAccessPrivateResource(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_private_resource.test", "name", "my private resource"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_private_resource.test", "description", "This is my private resource"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_private_resource.test", "endpoint_connection_methods.0.type", "client"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_private_resource.test", "endpoint_connection_methods.0.client_reachable_addresses.0", "192.168.1.0/24"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_private_resource.test", "addresses.0.addresses.0", "192.168.1.0/24"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_private_resource.test", "addresses.0.protocols_and_ports.0.protocol", "http/https"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_private_resource.test", "addresses.0.protocols_and_ports.0.ports", "443"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessPrivateResourceConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessPrivateResourceConfig() string {
	config := `resource "secureaccess_private_resource" "test" {` + "\n"
	config += `	name = "my private resource"` + "\n"
	config += `	description = "This is my private resource"` + "\n"
	config += `	endpoint_connection_methods = [{` + "\n"
	config += `		type = "client"` + "\n"
	config += `		client_reachable_addresses = ["192.168.1.0/24"]` + "\n"
	config += `	}]` + "\n"
	config += `	addresses = [{` + "\n"
	config += `		addresses = ["192.168.1.0/24"]` + "\n"
	config += `		protocols_and_ports = [{` + "\n"
	config += `			protocol = "http/https"` + "\n"
	config += `			ports = "443"` + "\n"
	config += `		}]` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_private_resource" "test" {
			id = secureaccess_private_resource.test.id
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
