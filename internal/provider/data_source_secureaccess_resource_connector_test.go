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

func TestAccDataSourceSecureAccessResourceConnector(t *testing.T) {
	if os.Getenv("TF_VAR_resource_connector_instance_id") == "" {
		t.Skip("skipping test, set environment variable TF_VAR_resource_connector_instance_id")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_resource_connector.test", "enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_resource_connector.test", "confirmed", "true"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector.test", "hostname"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector.test", "status"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector.test", "resource_connector_group_id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_resource_connector.test", "version"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessResourceConnectorPrerequisitesConfig + testAccDataSourceSecureAccessResourceConnectorConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccDataSourceSecureAccessResourceConnectorPrerequisitesConfig + testAccNamedByInstanceIdDataSourceSecureAccessResourceConnectorConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccDataSourceSecureAccessResourceConnectorPrerequisitesConfig = `
variable "resource_connector_instance_id" { default = null } // tests will set $TF_VAR_RESOURCE_CONNECTOR_INSTANCE_ID
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessResourceConnectorConfig() string {
	config := `resource "secureaccess_resource_connector" "test" {` + "\n"
	config += `	instance_id = var.resource_connector_instance_id` + "\n"
	config += `	enabled = true` + "\n"
	config += `	confirmed = true` + "\n"
	config += `	delete_on_destroy = false` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_resource_connector" "test" {
			id = secureaccess_resource_connector.test.id
		}
	`
	return config
}

func testAccNamedByInstanceIdDataSourceSecureAccessResourceConnectorConfig() string {
	config := `resource "secureaccess_resource_connector" "test" {` + "\n"
	config += `	instance_id = var.resource_connector_instance_id` + "\n"
	config += `	enabled = true` + "\n"
	config += `	confirmed = true` + "\n"
	config += `	delete_on_destroy = false` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_resource_connector" "test" {
			instance_id = secureaccess_resource_connector.test.instance_id
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
