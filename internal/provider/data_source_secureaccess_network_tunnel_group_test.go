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

func TestAccDataSourceSecureAccessNetworkTunnelGroup(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network_tunnel_group.test", "name", "my network tunnel group"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network_tunnel_group.test", "device_type", "ASA"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network_tunnel_group.test", "region", "us-east-1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_network_tunnel_group.test", "routing_type", "static"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_network_tunnel_group.test", "hubs.#"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_network_tunnel_group.test", "hubs.0.id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_network_tunnel_group.test", "hubs.0.is_primary"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_network_tunnel_group.test", "hubs.0.tunnel_id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_network_tunnel_group.test", "hubs.0.datacenter_name"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_network_tunnel_group.test", "hubs.0.datacenter_ip"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessNetworkTunnelGroupConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessNetworkTunnelGroupConfig() string {
	config := `resource "secureaccess_network_tunnel_group" "test" {` + "\n"
	config += `	name = "my network tunnel group"` + "\n"
	config += `	device_type = "ASA"` + "\n"
	config += `	region = "us-east-1"` + "\n"
	config += `	tunnel_id_format_email = "my_prefix"` + "\n"
	config += `	preshared_key = "t3stingTunn3lNow"` + "\n"
	config += `	routing_type = "static"` + "\n"
	config += `	routing_network_cidrs = ["192.168.0.0/24"]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_network_tunnel_group" "test" {
			id = secureaccess_network_tunnel_group.test.id
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
