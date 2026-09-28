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

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSource

func TestAccDataSourceSecureAccessNetwork(t *testing.T) {
	if os.Getenv("TF_VAR_NETWORKS") == "" {
		t.Skip("skipping test, set environment variable TF_VAR_NETWORKS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network.test", "name", "my_network"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network.test", "status", "OPEN"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network.test", "is_dynamic", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network.test", "prefix_length", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network.test", "prefix", "192.168.1.0"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessNetworkConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccNamedDataSourceSecureAccessNetworkConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessNetworkConfig() string {
	config := `resource "secureaccess_network" "test" {` + "\n"
	config += `	name = "my_network"` + "\n"
	config += `	status = "OPEN"` + "\n"
	config += `	is_dynamic = true` + "\n"
	config += `	prefix_length = 30` + "\n"
	config += `	prefix = "192.168.1.0"` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_network" "test" {
			id = secureaccess_network.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessNetworkConfig() string {
	config := `resource "secureaccess_network" "test" {` + "\n"
	config += `	name = "my_network"` + "\n"
	config += `	status = "OPEN"` + "\n"
	config += `	is_dynamic = true` + "\n"
	config += `	prefix_length = 30` + "\n"
	config += `	prefix = "192.168.1.0"` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_network" "test" {
			name = secureaccess_network.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
