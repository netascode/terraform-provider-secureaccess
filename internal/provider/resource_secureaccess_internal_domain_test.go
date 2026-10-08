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

func TestAccSecureAccessInternalDomain(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_internal_domain.test", "domain", "mydomain.local"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_internal_domain.test", "description", "My Internal Domain"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_internal_domain.test", "include_all_mobile_devices", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_internal_domain.test", "include_all_virtual_appliances", "false"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccSecureAccessInternalDomainPrerequisitesConfig + testAccSecureAccessInternalDomainConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessInternalDomainPrerequisitesConfig + testAccSecureAccessInternalDomainConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_internal_domain.test",
		ImportState:  true,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccSecureAccessInternalDomainPrerequisitesConfig = `
resource "secureaccess_site" "test" {
  name = "my_internal_domain"
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal

func testAccSecureAccessInternalDomainConfig_minimum() string {
	config := `resource "secureaccess_internal_domain" "test" {` + "\n"
	config += `	domain = "mydomain.local"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessInternalDomainConfig_all() string {
	config := `resource "secureaccess_internal_domain" "test" {` + "\n"
	config += `	domain = "mydomain.local"` + "\n"
	config += `	description = "My Internal Domain"` + "\n"
	config += `	include_all_mobile_devices = false` + "\n"
	config += `	include_all_virtual_appliances = false` + "\n"
	config += `	site_ids = [secureaccess_site.test.id]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll
