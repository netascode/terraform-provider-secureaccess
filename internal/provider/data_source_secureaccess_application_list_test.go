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

func TestAccDataSourceSecureAccessApplicationList(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_application_list.test", "name", "my_application_list"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_application_list.test", "is_default", "false"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessApplicationListPrerequisitesConfig + testAccDataSourceSecureAccessApplicationListConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccDataSourceSecureAccessApplicationListPrerequisitesConfig + testAccNamedDataSourceSecureAccessApplicationListConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccDataSourceSecureAccessApplicationListPrerequisitesConfig = `
data "secureaccess_application_categories" "test" {}
data "secureaccess_applications" "test" {
    items = {
        "Google Docs" = {}
    }
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessApplicationListConfig() string {
	config := `resource "secureaccess_application_list" "test" {` + "\n"
	config += `	name = "my_application_list"` + "\n"
	config += `	application_ids = [data.secureaccess_applications.test.items["Google Docs"].id]` + "\n"
	config += `	application_category_ids = [data.secureaccess_application_categories.test.items["Search"].id]` + "\n"
	config += `	is_default = false` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_application_list" "test" {
			id = secureaccess_application_list.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessApplicationListConfig() string {
	config := `resource "secureaccess_application_list" "test" {` + "\n"
	config += `	name = "my_application_list"` + "\n"
	config += `	application_ids = [data.secureaccess_applications.test.items["Google Docs"].id]` + "\n"
	config += `	application_category_ids = [data.secureaccess_application_categories.test.items["Search"].id]` + "\n"
	config += `	is_default = false` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_application_list" "test" {
			name = secureaccess_application_list.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
