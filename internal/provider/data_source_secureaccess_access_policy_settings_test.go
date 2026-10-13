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

func TestAccDataSourceSecureAccessAccessPolicySettings(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_policy_settings.test", "zero_trust_access_authentication", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_policy_settings.test", "zero_trust_access_authentication_timeout", "485"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_policy_settings.test", "editing_session_duration", "55"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_policy_settings.test", "editing_session_extension", "6"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_policy_settings.test", "microsoft_365_compatibility", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_policy_settings.test", "decryption_global", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_policy_settings.test", "decryption_disable_for_sources.0", "10.10.0.0/24"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_policy_settings.test", "log_decrypted_internet_destinations", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_access_policy_settings.test", "log_decrypted_private_resources", "true"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessAccessPolicySettingsConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessAccessPolicySettingsConfig() string {
	config := `resource "secureaccess_access_policy_settings" "test" {` + "\n"
	config += `	zero_trust_access_authentication = true` + "\n"
	config += `	zero_trust_access_authentication_timeout = 485` + "\n"
	config += `	editing_session_duration = 55` + "\n"
	config += `	editing_session_extension = 6` + "\n"
	config += `	microsoft_365_compatibility = true` + "\n"
	config += `	decryption_global = true` + "\n"
	config += `	decryption_disable_for_sources = ["10.10.0.0/24"]` + "\n"
	config += `	log_decrypted_internet_destinations = true` + "\n"
	config += `	log_decrypted_private_resources = true` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_access_policy_settings" "test" {
			depends_on = [secureaccess_access_policy_settings.test]
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
