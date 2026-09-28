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

func TestAccDataSourceSecureAccessDestinationList(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_destination_list.test", "name", "my_destination_list"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_destination_list.test", "destinations.0.destination", "example.com"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_destination_list.test", "destinations.0.comment", "This is a comment for the destination."))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_destination_list.test", "destinations.0.type", "domain"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessDestinationListConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessDestinationListConfig() string {
	config := `resource "secureaccess_destination_list" "test" {` + "\n"
	config += `	name = "my_destination_list"` + "\n"
	config += `	destinations = [{` + "\n"
	config += `		destination = "example.com"` + "\n"
	config += `		comment = "This is a comment for the destination."` + "\n"
	config += `		type = "domain"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_destination_list" "test" {
			id = secureaccess_destination_list.test.id
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
