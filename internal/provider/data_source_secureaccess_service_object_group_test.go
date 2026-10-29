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

func TestAccDataSourceSecureAccessServiceObjectGroup(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_service_object_group.test", "name", "my_service_object_group"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_service_object_group.test", "description", "This is my service object group."))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_service_object_group.test", "literals.0.protocol", "tcp"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_service_object_group.test", "literals.0.ports.0", "8443"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessServiceObjectGroupPrerequisitesConfig + testAccDataSourceSecureAccessServiceObjectGroupConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccDataSourceSecureAccessServiceObjectGroupPrerequisitesConfig + testAccNamedDataSourceSecureAccessServiceObjectGroupConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccDataSourceSecureAccessServiceObjectGroupPrerequisitesConfig = `
resource "secureaccess_service_object" test {
  name     = "my_service_object_group"
  protocol = "tcp"
  ports    = ["8080"]
}

resource "secureaccess_service_object_group" test_prereq {
  name = "my_service_object_group_prereq"
  literals = [
    {
      protocol = "tcp"
      ports    = ["9090"]
    }
  ]
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessServiceObjectGroupConfig() string {
	config := `resource "secureaccess_service_object_group" "test" {` + "\n"
	config += `	name = "my_service_object_group"` + "\n"
	config += `	description = "This is my service object group."` + "\n"
	config += `	service_object_groups = [{` + "\n"
	config += `		id = secureaccess_service_object_group.test_prereq.id` + "\n"
	config += `	}]` + "\n"
	config += `	service_objects = [{` + "\n"
	config += `		id = secureaccess_service_object.test.id` + "\n"
	config += `	}]` + "\n"
	config += `	literals = [{` + "\n"
	config += `		protocol = "tcp"` + "\n"
	config += `		ports = ["8443"]` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_service_object_group" "test" {
			id = secureaccess_service_object_group.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessServiceObjectGroupConfig() string {
	config := `resource "secureaccess_service_object_group" "test" {` + "\n"
	config += `	name = "my_service_object_group"` + "\n"
	config += `	description = "This is my service object group."` + "\n"
	config += `	service_object_groups = [{` + "\n"
	config += `		id = secureaccess_service_object_group.test_prereq.id` + "\n"
	config += `	}]` + "\n"
	config += `	service_objects = [{` + "\n"
	config += `		id = secureaccess_service_object.test.id` + "\n"
	config += `	}]` + "\n"
	config += `	literals = [{` + "\n"
	config += `		protocol = "tcp"` + "\n"
	config += `		ports = ["8443"]` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_service_object_group" "test" {
			name = secureaccess_service_object_group.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
