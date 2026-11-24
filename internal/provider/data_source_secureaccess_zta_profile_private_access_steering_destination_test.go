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

func TestAccDataSourceSecureAccessZTAProfilePrivateAccessSteeringDestination(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_profile_private_access_steering_destination.test", "destination", "internal.example.com"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_zta_profile_private_access_steering_destination.test", "destination_type"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessZTAProfilePrivateAccessSteeringDestinationPrerequisitesConfig + testAccDataSourceSecureAccessZTAProfilePrivateAccessSteeringDestinationConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccDataSourceSecureAccessZTAProfilePrivateAccessSteeringDestinationPrerequisitesConfig + testAccNamedByDestinationDataSourceSecureAccessZTAProfilePrivateAccessSteeringDestinationConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites

const testAccDataSourceSecureAccessZTAProfilePrivateAccessSteeringDestinationPrerequisitesConfig = `
resource "secureaccess_zta_profile" "test" {
  name     = "my ztna profile private steering destination"
  priority = 1
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessZTAProfilePrivateAccessSteeringDestinationConfig() string {
	config := `resource "secureaccess_zta_profile_private_access_steering_destination" "test" {` + "\n"
	config += `	zta_profile_id = secureaccess_zta_profile.test.id` + "\n"
	config += `	destination = "internal.example.com"` + "\n"
	config += `	exclusions = ["app.internal.example.com"]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_zta_profile_private_access_steering_destination" "test" {
			id = secureaccess_zta_profile_private_access_steering_destination.test.id
			zta_profile_id = secureaccess_zta_profile.test.id
		}
	`
	return config
}

func testAccNamedByDestinationDataSourceSecureAccessZTAProfilePrivateAccessSteeringDestinationConfig() string {
	config := `resource "secureaccess_zta_profile_private_access_steering_destination" "test" {` + "\n"
	config += `	zta_profile_id = secureaccess_zta_profile.test.id` + "\n"
	config += `	destination = "internal.example.com"` + "\n"
	config += `	exclusions = ["app.internal.example.com"]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_zta_profile_private_access_steering_destination" "test" {
			zta_profile_id = secureaccess_zta_profile.test.id
			destination = secureaccess_zta_profile_private_access_steering_destination.test.destination
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
