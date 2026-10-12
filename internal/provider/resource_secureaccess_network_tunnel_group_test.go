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

func TestAccSecureAccessNetworkTunnelGroup(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network_tunnel_group.test", "name", "my network tunnel group"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network_tunnel_group.test", "device_type", "ASA"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network_tunnel_group.test", "region", "us-east-1"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_network_tunnel_group.test", "routing_type", "static"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_network_tunnel_group.test", "hubs.#"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_network_tunnel_group.test", "hubs.0.id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_network_tunnel_group.test", "hubs.0.is_primary"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_network_tunnel_group.test", "hubs.0.tunnel_id"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_network_tunnel_group.test", "hubs.0.datacenter_name"))
	checks = append(checks, resource.TestCheckResourceAttrSet("secureaccess_network_tunnel_group.test", "hubs.0.datacenter_ip"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessNetworkTunnelGroupConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_network_tunnel_group.test",
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

func testAccSecureAccessNetworkTunnelGroupConfig_all() string {
	config := `resource "secureaccess_network_tunnel_group" "test" {` + "\n"
	config += `	name = "my network tunnel group"` + "\n"
	config += `	device_type = "ASA"` + "\n"
	config += `	region = "us-east-1"` + "\n"
	config += `	tunnel_id_format_email = "my_prefix"` + "\n"
	config += `	preshared_key = "t3stingTunn3lNow"` + "\n"
	config += `	routing_type = "static"` + "\n"
	config += `	routing_network_cidrs = ["192.168.0.0/24"]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

func TestAccSecureAccessNetworkTunnelGroup_Sequential(t *testing.T) {

	step_01 := `resource "secureaccess_network_tunnel_group" "test" {` + "\n" +
		`	name = "my network tunnel group seq"` + "\n" +
		`	device_type = "ASA"` + "\n" +
		`	region = "us-east-1"` + "\n" +
		`	tunnel_id_format_email = "my_prefix"` + "\n" +
		`	preshared_key = "t3stingTunn3lNow"` + "\n" +
		`	routing_type = "static"` + "\n" +
		`	routing_network_cidrs = ["192.168.10.0/24"]` + "\n" +
		`}` + "\n"

	step_02 := `resource "secureaccess_network_tunnel_group" "test" {` + "\n" +
		`	name = "my network tunnel group seq rename"` + "\n" +
		`	device_type = "ASA"` + "\n" +
		`	region = "us-east-1"` + "\n" +
		`	tunnel_id_format_email = "my_prefix_new"` + "\n" +
		`	preshared_key = "t3stingTunn3lNow2"` + "\n" +
		`	routing_type = "static"` + "\n" +
		`	routing_network_cidrs = ["192.168.10.0/24", "192.168.11.0/24"]` + "\n" +
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
