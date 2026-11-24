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

func TestAccDataSourceSecureAccessZTAProfile(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "name", "my ztna profile"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_zta_profile.test", "revision"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "priority", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "operating_system_macos", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "operating_system_windows", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "operating_system_linux", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "operating_system_ios", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "operating_system_android", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "operating_system_android_knox", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "all_users", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "all_groups", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "private_access_use_trusted_networks", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "private_access_end_user_pause", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "private_access_end_user_pause_duration", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "internet_access_steering_mode", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "internet_access_use_trusted_networks", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "internet_access_end_user_pause", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile.test", "internet_access_end_user_pause_duration", "60"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessZTAProfileConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccNamedDataSourceSecureAccessZTAProfileConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessZTAProfileConfig() string {
	config := `resource "secureaccess_zta_profile" "test" {` + "\n"
	config += `	name = "my ztna profile"` + "\n"
	config += `	priority = 1` + "\n"
	config += `	operating_system_macos = true` + "\n"
	config += `	operating_system_windows = true` + "\n"
	config += `	operating_system_linux = false` + "\n"
	config += `	operating_system_ios = true` + "\n"
	config += `	operating_system_android = true` + "\n"
	config += `	operating_system_android_knox = false` + "\n"
	config += `	all_users = true` + "\n"
	config += `	all_groups = true` + "\n"
	config += `	private_access_use_trusted_networks = true` + "\n"
	config += `	private_access_end_user_pause = true` + "\n"
	config += `	private_access_end_user_pause_duration = 30` + "\n"
	config += `	internet_access_steering_mode = 1` + "\n"
	config += `	internet_access_use_trusted_networks = true` + "\n"
	config += `	internet_access_end_user_pause = true` + "\n"
	config += `	internet_access_end_user_pause_duration = 60` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_zta_profile" "test" {
			id = secureaccess_zta_profile.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessZTAProfileConfig() string {
	config := `resource "secureaccess_zta_profile" "test" {` + "\n"
	config += `	name = "my ztna profile"` + "\n"
	config += `	priority = 1` + "\n"
	config += `	operating_system_macos = true` + "\n"
	config += `	operating_system_windows = true` + "\n"
	config += `	operating_system_linux = false` + "\n"
	config += `	operating_system_ios = true` + "\n"
	config += `	operating_system_android = true` + "\n"
	config += `	operating_system_android_knox = false` + "\n"
	config += `	all_users = true` + "\n"
	config += `	all_groups = true` + "\n"
	config += `	private_access_use_trusted_networks = true` + "\n"
	config += `	private_access_end_user_pause = true` + "\n"
	config += `	private_access_end_user_pause_duration = 30` + "\n"
	config += `	internet_access_steering_mode = 1` + "\n"
	config += `	internet_access_use_trusted_networks = true` + "\n"
	config += `	internet_access_end_user_pause = true` + "\n"
	config += `	internet_access_end_user_pause_duration = 60` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_zta_profile" "test" {
			name = secureaccess_zta_profile.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
