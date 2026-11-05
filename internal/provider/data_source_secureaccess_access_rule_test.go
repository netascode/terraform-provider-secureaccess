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

func TestAccDataSourceSecureAccessAccessRule(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_rule.test", "name", "my_internal_network"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_rule.test", "description", "Allows access to internal network"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_rule.test", "action", "allow"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_rule.test", "enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_rule.test", "log_level", "LOG_ALL"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_rule.test", "traffic_type", "PRIVATE_NETWORK"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_rule.test", "ips_profile_id", "524837"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessAccessRulePrerequisitesConfig + testAccDataSourceSecureAccessAccessRuleConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccDataSourceSecureAccessAccessRulePrerequisitesConfig = `
resource "secureaccess_network_object" test {
  name  = "my_access_rule_network_object"
  type  = "network"
  value = "192.168.1.0/24"
}

resource "secureaccess_service_object" test {
  name     = "my_access_rule_service_object"
  protocol = "tcp"
  ports    = ["8080"]
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessAccessRuleConfig() string {
	config := `resource "secureaccess_access_rule" "test" {` + "\n"
	config += `	name = "my_internal_network"` + "\n"
	config += `	description = "Allows access to internal network"` + "\n"
	config += `	action = "allow"` + "\n"
	config += `	priority = 1` + "\n"
	config += `	enabled = true` + "\n"
	config += `	source_network_object_ids = [secureaccess_network_object.test.id]` + "\n"
	config += `	source_private_literals = ["192.168.0.1"]` + "\n"
	config += `	source_operator = "OR"` + "\n"
	config += `	destination_network_object_ids = [secureaccess_network_object.test.id]` + "\n"
	config += `	destination_service_object_ids = [secureaccess_service_object.test.id]` + "\n"
	config += `	destination_operator = "AND"` + "\n"
	config += `	destination_operator_network_service = "AND"` + "\n"
	config += `	log_level = "LOG_ALL"` + "\n"
	config += `	traffic_type = "PRIVATE_NETWORK"` + "\n"
	config += `	ips_profile_id = 524837` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_access_rule" "test" {
			id = secureaccess_access_rule.test.id
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
