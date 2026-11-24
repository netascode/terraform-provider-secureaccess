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

func TestAccSecureAccessZTAProfilePrivateAccessSteeringPrivateResourceMappings(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_zta_profile_private_access_steering_private_resource_mappings.test", "revision"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessZTAProfilePrivateAccessSteeringPrivateResourceMappingsPrerequisitesConfig + testAccSecureAccessZTAProfilePrivateAccessSteeringPrivateResourceMappingsConfig_all(),
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

const testAccSecureAccessZTAProfilePrivateAccessSteeringPrivateResourceMappingsPrerequisitesConfig = `
resource "secureaccess_zta_profile" "test" {
  name     = "my ztna profile private resource mappings"
  priority = 1
}

resource "secureaccess_private_resource" "test" {
  name = "my ztna profile private resource mappings"
  endpoint_connection_methods = [
    {
      type                       = "client"
      client_reachable_addresses = ["192.168.1.0/24"]
    }
  ]
  addresses = [
    {
      addresses = ["192.168.1.0/24"]
      protocols_and_ports = [
        {
          protocol = "http/https"
          ports    = "443"
        }
      ]
    }
  ]
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessZTAProfilePrivateAccessSteeringPrivateResourceMappingsConfig_all() string {
	config := `resource "secureaccess_zta_profile_private_access_steering_private_resource_mappings" "test" {` + "\n"
	config += `	zta_profile_id = secureaccess_zta_profile.test.id` + "\n"
	config += `	private_resources = [{` + "\n"
	config += `		id = secureaccess_private_resource.test.id` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll
