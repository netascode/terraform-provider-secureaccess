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

func TestAccDataSourceSecureAccessResourceConnectorGroup(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_resource_connector_group.test", "name", "my resource connector group"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_resource_connector_group.test", "region", "us-west-2"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_resource_connector_group.test", "environment", "container"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector_group.test", "provisioning_key_expires_at"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector_group.test", "status"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector_group.test", "connectors_count"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector_group.test", "connected_connectors_count"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector_group.test", "disconnected_connectors_count"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector_group.test", "unable_to_sync_connectors_count"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector_group.test", "need_update_connectors_count"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector_group.test", "failed_update_agents_attempts"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessResourceConnectorGroupConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccNamedDataSourceSecureAccessResourceConnectorGroupConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessResourceConnectorGroupConfig() string {
	config := `resource "secureaccess_resource_connector_group" "test" {` + "\n"
	config += `	name = "my resource connector group"` + "\n"
	config += `	region = "us-west-2"` + "\n"
	config += `	environment = "container"` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_resource_connector_group" "test" {
			id = secureaccess_resource_connector_group.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessResourceConnectorGroupConfig() string {
	config := `resource "secureaccess_resource_connector_group" "test" {` + "\n"
	config += `	name = "my resource connector group"` + "\n"
	config += `	region = "us-west-2"` + "\n"
	config += `	environment = "container"` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_resource_connector_group" "test" {
			name = secureaccess_resource_connector_group.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
