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

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc

func TestAccSecureAccessResourceConnector(t *testing.T) {
	if os.Getenv("TF_VAR_resource_connector_instance_id") == "" {
		t.Skip("skipping test, set environment variable TF_VAR_resource_connector_instance_id")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_resource_connector.test", "enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_resource_connector.test", "confirmed", "true"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_resource_connector.test", "hostname"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_resource_connector.test", "status"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_resource_connector.test", "resource_connector_group_id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_resource_connector.test", "version"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_resource_connector.test", "delete_on_destroy", "false"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessResourceConnectorPrerequisitesConfig + testAccSecureAccessResourceConnectorConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccSecureAccessResourceConnectorPrerequisitesConfig = `
variable "resource_connector_instance_id" { default = null } // tests will set $TF_VAR_RESOURCE_CONNECTOR_INSTANCE_ID
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessResourceConnectorConfig_all() string {
	config := `resource "secureaccess_resource_connector" "test" {` + "\n"
	config += `	instance_id = var.resource_connector_instance_id` + "\n"
	config += `	enabled = true` + "\n"
	config += `	confirmed = true` + "\n"
	config += `	delete_on_destroy = false` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll
