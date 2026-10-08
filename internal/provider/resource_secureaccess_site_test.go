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

func TestAccSecureAccessSite(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_site.test", "name", "my_site"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_site.test", "origin_id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_site.test", "is_default"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessSiteConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_site.test",
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
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessSiteConfig_all() string {
	config := `resource "secureaccess_site" "test" {` + "\n"
	config += `	name = "my_site"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

func TestAccSecureAccessSite_Sequential(t *testing.T) {

	step_01 := `resource "secureaccess_site" "test" {` + "\n" +
		`	name = "my_site_seq"` + "\n" +
		`}` + "\n"

	step_02 := `resource "secureaccess_site" "test" {` + "\n" +
		`	name = "my_site_seq_rename"` + "\n" +
		`}` + "\n"

	steps := []resource.TestStep{{
		Config: step_01,
	}, {
		Config: step_02,
	}}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}
