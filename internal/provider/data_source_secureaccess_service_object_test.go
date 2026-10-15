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

func TestAccDataSourceSecureAccessServiceObject(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_service_object.test", "name", "my_service_object"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_service_object.test", "description", "This is my service object."))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_service_object.test", "protocol", "tcp"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_service_object.test", "ports.0", "8080"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessServiceObjectConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccNamedDataSourceSecureAccessServiceObjectConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessServiceObjectConfig() string {
	config := `resource "secureaccess_service_object" "test" {` + "\n"
	config += `	name = "my_service_object"` + "\n"
	config += `	description = "This is my service object."` + "\n"
	config += `	protocol = "tcp"` + "\n"
	config += `	ports = ["8080"]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_service_object" "test" {
			id = secureaccess_service_object.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessServiceObjectConfig() string {
	config := `resource "secureaccess_service_object" "test" {` + "\n"
	config += `	name = "my_service_object"` + "\n"
	config += `	description = "This is my service object."` + "\n"
	config += `	protocol = "tcp"` + "\n"
	config += `	ports = ["8080"]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_service_object" "test" {
			name = secureaccess_service_object.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
