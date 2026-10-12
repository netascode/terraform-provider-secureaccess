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

func TestAccSecureAccessPrivateResource(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_private_resource.test", "name", "my private resource"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_private_resource.test", "description", "This is my private resource"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_private_resource.test", "endpoint_connection_methods.0.type", "client"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_private_resource.test", "endpoint_connection_methods.0.client_reachable_addresses.0", "192.168.1.0/24"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_private_resource.test", "addresses.0.addresses.0", "192.168.1.0/24"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_private_resource.test", "addresses.0.protocols_and_ports.0.protocol", "http/https"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_private_resource.test", "addresses.0.protocols_and_ports.0.ports", "443"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccSecureAccessPrivateResourceConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessPrivateResourceConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_private_resource.test",
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

func testAccSecureAccessPrivateResourceConfig_minimum() string {
	config := `resource "secureaccess_private_resource" "test" {` + "\n"
	config += `	name = "my private resource"` + "\n"
	config += `	endpoint_connection_methods = [{` + "\n"
	config += `		type = "client"` + "\n"
	config += `		client_reachable_addresses = ["192.168.1.0/24"]` + "\n"
	config += `	}]` + "\n"
	config += `	addresses = [{` + "\n"
	config += `		addresses = ["192.168.1.0/24"]` + "\n"
	config += `		protocols_and_ports = [{` + "\n"
	config += `			protocol = "http/https"` + "\n"
	config += `			ports = "443"` + "\n"
	config += `		}]` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessPrivateResourceConfig_all() string {
	config := `resource "secureaccess_private_resource" "test" {` + "\n"
	config += `	name = "my private resource"` + "\n"
	config += `	description = "This is my private resource"` + "\n"
	config += `	endpoint_connection_methods = [{` + "\n"
	config += `		type = "client"` + "\n"
	config += `		client_reachable_addresses = ["192.168.1.0/24"]` + "\n"
	config += `	}]` + "\n"
	config += `	addresses = [{` + "\n"
	config += `		addresses = ["192.168.1.0/24"]` + "\n"
	config += `		protocols_and_ports = [{` + "\n"
	config += `			protocol = "http/https"` + "\n"
	config += `			ports = "443"` + "\n"
	config += `		}]` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll
