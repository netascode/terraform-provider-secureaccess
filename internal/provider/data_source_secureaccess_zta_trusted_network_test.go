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

func TestAccDataSourceSecureAccessZTATrustedNetwork(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_trusted_network.test", "name", "my zta trusted network"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_trusted_network.test", "inspect", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_trusted_network.test", "dns_servers.0.ip", "1.1.1.1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_trusted_network.test", "dns_domains.0.name", "*.example.com"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_trusted_network.test", "trusted_servers.0.url", "https://trusted.example.com:8443"))
	checks = append(checks, resource.TestCheckResourceAttr("data.secureaccess_zta_trusted_network.test", "trusted_servers.0.certificate_hash", "x2v1b7d9p0e3f6g8a4w5z2q1k9m7r5n8l6j0t3u2s9y8v5c1b6h4d7w2x9g3p5a1"))
	checks = append(checks, resource.TestCheckResourceAttrSet("data.secureaccess_zta_trusted_network.test", "revision"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessZTATrustedNetworkConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: testAccNamedDataSourceSecureAccessZTATrustedNetworkConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig

func testAccDataSourceSecureAccessZTATrustedNetworkConfig() string {
	config := `resource "secureaccess_zta_trusted_network" "test" {` + "\n"
	config += `	name = "my zta trusted network"` + "\n"
	config += `	inspect = 1` + "\n"
	config += `	dns_servers = [{` + "\n"
	config += `		ip = "1.1.1.1"` + "\n"
	config += `	}]` + "\n"
	config += `	dns_domains = [{` + "\n"
	config += `		name = "*.example.com"` + "\n"
	config += `	}]` + "\n"
	config += `	trusted_servers = [{` + "\n"
	config += `		url = "https://trusted.example.com:8443"` + "\n"
	config += `		certificate_hash = "x2v1b7d9p0e3f6g8a4w5z2q1k9m7r5n8l6j0t3u2s9y8v5c1b6h4d7w2x9g3p5a1"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_zta_trusted_network" "test" {
			id = secureaccess_zta_trusted_network.test.id
		}
	`
	return config
}

func testAccNamedDataSourceSecureAccessZTATrustedNetworkConfig() string {
	config := `resource "secureaccess_zta_trusted_network" "test" {` + "\n"
	config += `	name = "my zta trusted network"` + "\n"
	config += `	inspect = 1` + "\n"
	config += `	dns_servers = [{` + "\n"
	config += `		ip = "1.1.1.1"` + "\n"
	config += `	}]` + "\n"
	config += `	dns_domains = [{` + "\n"
	config += `		name = "*.example.com"` + "\n"
	config += `	}]` + "\n"
	config += `	trusted_servers = [{` + "\n"
	config += `		url = "https://trusted.example.com:8443"` + "\n"
	config += `		certificate_hash = "x2v1b7d9p0e3f6g8a4w5z2q1k9m7r5n8l6j0t3u2s9y8v5c1b6h4d7w2x9g3p5a1"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "secureaccess_zta_trusted_network" "test" {
			name = secureaccess_zta_trusted_network.test.name
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig
