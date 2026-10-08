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

func TestAccDataSourceSecureAccessInternalDomain(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_internal_domain.test", "domain", "mydomain.local"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_internal_domain.test", "description", "My Internal Domain"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_internal_domain.test", "include_all_mobile_devices", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_internal_domain.test", "include_all_virtual_appliances", "false"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessInternalDomainPrerequisitesConfig + testAccDataSourceSecureAccessInternalDomainConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccDataSourceSecureAccessInternalDomainPrerequisitesConfig = `
resource "secureaccess_site" "test" {
  name = "my_internal_domain"
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessInternalDomainConfig() string {
	config := `resource "secureaccess_internal_domain" "test" {` + "\n"
	config += `	domain = "mydomain.local"` + "\n"
	config += `	description = "My Internal Domain"` + "\n"
	config += `	include_all_mobile_devices = false` + "\n"
	config += `	include_all_virtual_appliances = false` + "\n"
	config += `	site_ids = [secureaccess_site.test.id]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_internal_domain" "test" {
			id = secureaccess_internal_domain.test.id
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
