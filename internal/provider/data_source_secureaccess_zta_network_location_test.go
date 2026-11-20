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

func TestAccDataSourceSecureAccessZTANetworkLocation(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_network_location.test", "name", "my_network_location"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_network_location.test", "egress_criteria.0.source_type", "manual"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_network_location.test", "egress_criteria.0.manual_address", "203.0.113.10"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_network_location.test", "egress_exclusion_criteria.0.source_type", "manual"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_network_location.test", "egress_exclusion_criteria.0.manual_address", "203.0.113.200"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_network_location.test", "private_subnet_criteria.0.source_type", "manual"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_network_location.test", "private_subnet_criteria.0.manual_address", "10.10.0.0/16"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_zta_network_location.test", "revision"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessZTANetworkLocationConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccNamedDataSourceSecureAccessZTANetworkLocationConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessZTANetworkLocationConfig() string {
	config := `resource "secureaccess_zta_network_location" "test" {` + "\n"
	config += `	name = "my_network_location"` + "\n"
	config += `	egress_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `		manual_address = "203.0.113.10"` + "\n"
	config += `	}]` + "\n"
	config += `	egress_exclusion_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `		manual_address = "203.0.113.200"` + "\n"
	config += `	}]` + "\n"
	config += `	private_subnet_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `		manual_address = "10.10.0.0/16"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_zta_network_location" "test" {
			id = secureaccess_zta_network_location.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessZTANetworkLocationConfig() string {
	config := `resource "secureaccess_zta_network_location" "test" {` + "\n"
	config += `	name = "my_network_location"` + "\n"
	config += `	egress_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `		manual_address = "203.0.113.10"` + "\n"
	config += `	}]` + "\n"
	config += `	egress_exclusion_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `		manual_address = "203.0.113.200"` + "\n"
	config += `	}]` + "\n"
	config += `	private_subnet_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `		manual_address = "10.10.0.0/16"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_zta_network_location" "test" {
			name = secureaccess_zta_network_location.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
