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

func TestAccSecureAccessZTANetworkLocation(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_zta_network_location.test", "name", "my_network_location"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_zta_network_location.test", "egress_criteria.0.source_type", "manual"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_zta_network_location.test", "egress_criteria.0.manual_address", "203.0.113.10"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_zta_network_location.test", "egress_exclusion_criteria.0.source_type", "manual"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_zta_network_location.test", "egress_exclusion_criteria.0.manual_address", "203.0.113.200"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_zta_network_location.test", "private_subnet_criteria.0.source_type", "manual"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_zta_network_location.test", "private_subnet_criteria.0.manual_address", "10.10.0.0/16"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_zta_network_location.test", "revision"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccSecureAccessZTANetworkLocationConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessZTANetworkLocationConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_zta_network_location.test",
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

func testAccSecureAccessZTANetworkLocationConfig_minimum() string {
	config := `resource "secureaccess_zta_network_location" "test" {` + "\n"
	config += `	name = "my_network_location"` + "\n"
	config += `	egress_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessZTANetworkLocationConfig_all() string {
	config := `resource "secureaccess_zta_network_location" "test" {` + "\n"
	config += `	name = "my_network_location"` + "\n"
	config += `	egress_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `		manual_address = "203.0.113.10"` + "\n"
	config += `	}]` + "\n"
	config += `	egress_exclusion_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `		manual_address = "203.0.113.200"` + "\n"
	config += `	}]` + "\n"
	config += `	private_subnet_criteria = [{` + "\n"
	config += `		source_type = "manual"` + "\n"
	config += `		manual_address = "10.10.0.0/16"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll
