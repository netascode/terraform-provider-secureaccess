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

func TestAccDataSourceSecureAccessNetworkObjectGroup(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network_object_group.test", "name", "my_network_object"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network_object_group.test", "description", "This is my network object group."))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network_object_group.test", "literals.0.type", "network"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network_object_group.test", "literals.0.value", "192.168.2.0/24"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessNetworkObjectGroupPrerequisitesConfig + testAccDataSourceSecureAccessNetworkObjectGroupConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccDataSourceSecureAccessNetworkObjectGroupPrerequisitesConfig + testAccNamedDataSourceSecureAccessNetworkObjectGroupConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccDataSourceSecureAccessNetworkObjectGroupPrerequisitesConfig = `
resource "secureaccess_network_object" test {
  name = "my_network_object_group"
  type = "network"
  value = "192.168.1.0/24"
}

resource "secureaccess_network_object_group" test_prereq {
  name = "my_network_object_group"
  literals = [
    {
      type = "network"
      value = "192.168.10.0/24"
    }
  ]
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessNetworkObjectGroupConfig() string {
	config := `resource "secureaccess_network_object_group" "test" {` + "\n"
	config += `	name = "my_network_object"` + "\n"
	config += `	description = "This is my network object group."` + "\n"
	config += `	network_object_groups = [{` + "\n"
	config += `		id = secureaccess_network_object_group.test_prereq.id` + "\n"
	config += `	}]` + "\n"
	config += `	network_objects = [{` + "\n"
	config += `		id = secureaccess_network_object.test.id` + "\n"
	config += `	}]` + "\n"
	config += `	literals = [{` + "\n"
	config += `		type = "network"` + "\n"
	config += `		value = "192.168.2.0/24"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_network_object_group" "test" {
			id = secureaccess_network_object_group.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessNetworkObjectGroupConfig() string {
	config := `resource "secureaccess_network_object_group" "test" {` + "\n"
	config += `	name = "my_network_object"` + "\n"
	config += `	description = "This is my network object group."` + "\n"
	config += `	network_object_groups = [{` + "\n"
	config += `		id = secureaccess_network_object_group.test_prereq.id` + "\n"
	config += `	}]` + "\n"
	config += `	network_objects = [{` + "\n"
	config += `		id = secureaccess_network_object.test.id` + "\n"
	config += `	}]` + "\n"
	config += `	literals = [{` + "\n"
	config += `		type = "network"` + "\n"
	config += `		value = "192.168.2.0/24"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_network_object_group" "test" {
			name = secureaccess_network_object_group.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
