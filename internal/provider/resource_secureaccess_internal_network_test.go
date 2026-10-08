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

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc

func TestAccSecureAccessInternalNetwork(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_internal_network.test", "name", "my_internal_network"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_internal_network.test", "prefix", "10.1.1.0"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_internal_network.test", "prefix_length", "24"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessInternalNetworkPrerequisitesConfig + testAccSecureAccessInternalNetworkConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_internal_network.test",
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

const testAccSecureAccessInternalNetworkPrerequisitesConfig = `
resource "secureaccess_site" "test" {
  name = "my_internal_network"
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessInternalNetworkConfig_all() string {
	config := `resource "secureaccess_internal_network" "test" {` + "\n"
	config += `	name = "my_internal_network"` + "\n"
	config += `	prefix = "10.1.1.0"` + "\n"
	config += `	prefix_length = 24` + "\n"
	config += `	site_id = secureaccess_site.test.id` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

func TestAccSecureAccessInternalNetwork_Sequential(t *testing.T) {

	step_01 := `resource "secureaccess_internal_network" "test" {` + "\n" +
		`	name = "my_internal_network_seq"` + "\n" +
		`	prefix = "10.1.1.0"` + "\n" +
		`	prefix_length = 24` + "\n" +
		`	site_id = secureaccess_site.test.id` + "\n" +
		`}` + "\n"

	step_02 := `resource "secureaccess_internal_network" "test" {` + "\n" +
		`	name = "my_internal_network_seq_rename"` + "\n" +
		`	prefix = "10.1.2.0"` + "\n" +
		`	prefix_length = 24` + "\n" +
		`	site_id = secureaccess_site.test.id` + "\n" +
		`}` + "\n"

	steps := []resource.TestStep{{
		Config: testAccSecureAccessInternalNetworkPrerequisitesConfig + step_01,
	}, {
		Config: testAccSecureAccessInternalNetworkPrerequisitesConfig + step_02,
	}}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}
