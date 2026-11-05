// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
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

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Manual tests of the access rule logical operators, not covered by the generated tests

const testAccSecureAccessAccessRuleConditionsPrerequisitesConfig = `
variable "roaming_computer_name" { default = null } // tests will set $TF_VAR_roaming_computer_name

data "secureaccess_roaming_computer" "test" {
  name = var.roaming_computer_name
}

resource "secureaccess_network_object" "test" {
  name  = "my_access_rule_conditions_network_object"
  type  = "network"
  value = "192.168.2.0/24"
}

resource "secureaccess_service_object" "test" {
  name     = "my_access_rule_conditions_service_object"
  protocol = "tcp"
  ports    = ["8443"]
}
`

func TestAccSecureAccessAccessRuleConditions(t *testing.T) {
	if os.Getenv("TF_VAR_roaming_computer_name") == "" {
		t.Skip("skipping test, set environment variable TF_VAR_roaming_computer_name")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Internet rule: source identity AND identity type, content category AND (network OR service objects)
			{
				Config: testAccSecureAccessAccessRuleConditionsPrerequisitesConfig + `
resource "secureaccess_access_rule" "test" {
  name                                 = "my_access_rule_conditions"
  action                               = "block"
  log_level                            = "LOG_ALL"
  traffic_type                          = "PUBLIC_INTERNET"
  source_identity_ids                  = [data.secureaccess_roaming_computer.test.origin_id]
  source_identity_type_ids             = [7]
  source_operator                      = "AND"
  destination_content_category_ids     = [72]
  destination_network_object_ids       = [secureaccess_network_object.test.id]
  destination_service_object_ids       = [secureaccess_service_object.test.id]
  destination_operator                 = "AND"
  destination_operator_network_service = "OR"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "source_operator", "AND"),
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_operator", "AND"),
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_operator_network_service", "OR"),
				),
			},
			{
				ResourceName:            "secureaccess_access_rule.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"priority", "description", "enabled"},
			},
			// Source AND with internet literals, destination OR with literals, updated in place
			{
				Config: testAccSecureAccessAccessRuleConditionsPrerequisitesConfig + `
resource "secureaccess_access_rule" "test" {
  name                           = "my_access_rule_conditions"
  action                         = "block"
  log_level                      = "LOG_ALL"
  traffic_type                   = "PUBLIC_INTERNET"
  source_identity_type_ids       = [7]
  source_internet_literals       = [{ ips = ["4.5.6.7", "1.2.3.4", "10.0.0.0/24"], ports = ["443", "8443", "8000-8010"] }]
  source_operator                = "AND"
  destination_content_category_ids = [72]
  destination_network_object_ids = [secureaccess_network_object.test.id]
  destination_literals = [
    { ips = ["192.168.1.1", "192.168.3.1", "192.168.2.1"], ports = ["333", "3344-3345"], protocol = "TCP" },
  ]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "source_operator", "AND"),
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "source_internet_literals.#", "1"),
				),
			},
			{
				ResourceName:            "secureaccess_access_rule.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"priority", "description", "enabled"},
			},
			// Destination literals with several TCP items and ICMP, updated in place
			{
				Config: testAccSecureAccessAccessRuleConditionsPrerequisitesConfig + `
resource "secureaccess_access_rule" "test" {
  name                             = "my_access_rule_conditions"
  action                           = "block"
  log_level                        = "LOG_ALL"
  traffic_type                     = "PUBLIC_INTERNET"
  source_identity_type_ids         = [7]
  destination_content_category_ids = [72]
  destination_network_object_ids   = [secureaccess_network_object.test.id]
  destination_literals = [
    { ips = ["192.168.1.1", "192.168.3.1", "192.168.2.1"], ports = ["333", "3344-3345"], protocol = "TCP" },
    { ips = ["192.168.1.1"], ports = ["443"], protocol = "TCP" },
    { ips = ["1.2.3.4"], ports = ["0-65535"], protocol = "ICMP" },
  ]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "source_operator", "OR"),
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_operator", "OR"),
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_operator_network_service", "AND"),
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_literals.#", "3"),
				),
			},
			{
				ResourceName:            "secureaccess_access_rule.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"priority", "description", "enabled"},
			},
		},
	})
}

// Internet block rule with destination geolocations
func TestAccSecureAccessAccessRuleDestinationGeolocations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "secureaccess_access_rule" "test" {
  name                     = "my_access_rule_destination_geolocations"
  action                   = "block"
  log_level                = "LOG_ALL"
  traffic_type             = "PUBLIC_INTERNET"
  source_identity_type_ids = [7]
  destination_geolocations = ["PL", "DE"]
}
`,
				Check: resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_geolocations.#", "2"),
			},
			{
				ResourceName:            "secureaccess_access_rule.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"priority", "description", "enabled"},
			},
		},
	})
}

// Private rule with Security Group Tags, ORed with a network object, requires at least two `sgt_id` identities in the org
func TestAccSecureAccessAccessRuleSecurityGroupTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "secureaccess_identities" "sgt" {
  identity_types = ["sgt_id"]
}

resource "secureaccess_network_object" "test" {
  name  = "my_access_rule_sgt_network_object"
  type  = "network"
  value = "192.168.6.0/24"
}

resource "secureaccess_access_rule" "test" {
  name                               = "my_access_rule_security_group_tags"
  action                             = "allow"
  log_level                          = "LOG_ALL"
  traffic_type                       = "PRIVATE_NETWORK"
  destination_network_object_ids     = [secureaccess_network_object.test.id]
  destination_security_group_tag_ids = slice(sort([for v in values(data.secureaccess_identities.sgt.items["sgt_id"]) : tostring(v.id)]), 0, 2)
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_security_group_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_network_object_ids.#", "1"),
				),
			},
			{
				ResourceName:            "secureaccess_access_rule.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"priority", "description", "enabled"},
			},
		},
	})
}

// Private rule with Private Resources and Private Resource Groups created in the same apply
func TestAccSecureAccessAccessRulePrivateResources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "secureaccess_private_resource" "test" {
  name = "my access rule private resource"
  endpoint_connection_methods = [
    {
      type                       = "client"
      client_reachable_addresses = ["192.168.5.0/24"]
    }
  ]
  addresses = [
    {
      addresses = ["192.168.5.0/24"]
      protocols_and_ports = [
        {
          protocol = "http/https"
          ports    = "443"
        }
      ]
    }
  ]
}

resource "secureaccess_private_resource_group" "test" {
  name                 = "my access rule private resource group"
  private_resource_ids = [secureaccess_private_resource.test.id]
}

resource "secureaccess_access_rule" "test" {
  name                                   = "my_access_rule_private_resources"
  action                                 = "allow"
  log_level                              = "LOG_ALL"
  traffic_type                           = "PRIVATE_NETWORK"
  source_identity_type_ids               = [7]
  destination_private_resource_ids       = [secureaccess_private_resource.test.id]
  destination_private_resource_group_ids = [secureaccess_private_resource_group.test.id]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_private_resource_ids.#", "1"),
					resource.TestCheckResourceAttr("secureaccess_access_rule.test", "destination_private_resource_group_ids.#", "1"),
				),
			},
			{
				ResourceName:            "secureaccess_access_rule.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"priority", "description", "enabled"},
			},
		},
	})
}

func TestAccSecureAccessAccessRuleConditionsValidation(t *testing.T) {
	config := func(attributes string) string {
		return `
resource "secureaccess_access_rule" "test" {
  name      = "my_access_rule_conditions_validation"
  log_level = "LOG_ALL"
` + attributes + `
}
`
	}
	steps := []struct {
		attributes string
		err        string
	}{
		{`
  action                           = "block"
  traffic_type                      = "PRIVATE_NETWORK"
  destination_content_category_ids = [72]`, `destination_content_category_ids is only supported with traffic_type "PUBLIC_INTERNET"`},
		{`
  action                           = "allow"
  traffic_type                      = "PUBLIC_INTERNET"
  destination_private_resource_ids = [1]`, `destination_private_resource_ids is only supported with traffic_type "PRIVATE_NETWORK"`},
		{`
  action      = "isolate"
  traffic_type = "PRIVATE_NETWORK"`, `action "isolate" is only supported with traffic_type "PUBLIC_INTERNET"`},
		{`
  action                               = "allow"
  traffic_type                          = "PRIVATE_NETWORK"
  destination_network_object_ids       = [1]
  destination_network_object_group_ids = [2]
  destination_operator                 = "AND"`, `destination_operator "AND" requires conditions of at least two of`},
		{`
  action                               = "allow"
  traffic_type                          = "PRIVATE_NETWORK"
  destination_network_object_ids       = [1]
  destination_service_object_ids       = [2]
  destination_operator                 = "AND"
  destination_operator_network_service = "OR"`, `destination_operator_network_service "OR" requires network objects, service objects and content categories or applications`},
		{`
  action                           = "allow"
  traffic_type                      = "PRIVATE_NETWORK"
  destination_network_object_ids   = [1]
  destination_private_resource_ids = [2]
  destination_operator             = "AND"`, `destination_operator "AND" does not support destination_literals, destination_geolocations, Destination Lists, Private Resources, Private Resource Groups and Security Group Tags`},
		{`
  action                                 = "allow"
  traffic_type                           = "PRIVATE_NETWORK"
  destination_network_object_ids         = [1]
  destination_private_resource_group_ids = [2]
  destination_operator                   = "AND"`, `destination_operator "AND" does not support destination_literals, destination_geolocations, Destination Lists, Private Resources, Private Resource Groups and Security Group Tags`},
		{`
  action                                 = "block"
  traffic_type                           = "PUBLIC_INTERNET"
  destination_private_resource_group_ids = [1]`, `destination_private_resource_group_ids is only supported with traffic_type "PRIVATE_NETWORK"`},
		{`
  action                             = "allow"
  traffic_type                       = "PRIVATE_NETWORK"
  destination_network_object_ids     = [1]
  destination_security_group_tag_ids = [2]
  destination_operator               = "AND"`, `destination_operator "AND" does not support destination_literals, destination_geolocations, Destination Lists, Private Resources, Private Resource Groups and Security Group Tags`},
		{`
  action                             = "block"
  traffic_type                       = "PUBLIC_INTERNET"
  destination_security_group_tag_ids = [1]`, `destination_security_group_tag_ids is only supported with traffic_type "PRIVATE_NETWORK"`},
		{`
  action                         = "allow"
  traffic_type                   = "PRIVATE_NETWORK"
  destination_network_object_ids = [1]
  destination_literals           = [{ ips = ["192.168.1.1"], ports = ["443"], protocol = "TCP" }]
  destination_operator           = "AND"`, `destination_operator "AND" does not support destination_literals`},
		{`
  action               = "allow"
  traffic_type         = "PRIVATE_NETWORK"
  destination_literals = [{ ips = ["192.168.1.1"], ports = ["443"], protocol = "GRE" }]`, `value must be one of`},
		{`
  action                    = "allow"
  traffic_type               = "PRIVATE_NETWORK"
  source_identity_ids       = [1]
  source_network_object_ids = [2]
  source_operator           = "AND"`, `source_operator "AND" does not support source_private_literals, source network objects`},
		{`
  action              = "allow"
  traffic_type        = "PRIVATE_NETWORK"
  source_identity_ids = [1]
  source_private_literals = ["192.168.0.1"]
  source_operator     = "AND"`, `source_operator "AND" does not support source_private_literals`},
		{`
  action                  = "block"
  traffic_type            = "PUBLIC_INTERNET"
  source_private_literals = ["192.168.0.1"]`, `source_private_literals is only supported with traffic_type "PRIVATE_NETWORK"`},
		{`
  action                   = "allow"
  traffic_type             = "PRIVATE_NETWORK"
  source_internet_literals = [{ ips = ["192.168.0.1"], ports = ["443"] }]`, `source_internet_literals is only supported with traffic_type "PUBLIC_INTERNET"`},
		{`
  action                   = "block"
  traffic_type             = "PUBLIC_INTERNET"
  source_internet_literals = [{ ips = ["192.168.0.1"], ports = ["443"] }, { ips = ["192.168.0.2"], ports = ["80"] }]`, `list must contain at most 1 elements`},
		{`
  action                   = "block"
  traffic_type             = "PRIVATE_NETWORK"
  destination_geolocations = ["PL"]`, `destination_geolocations is only supported with traffic_type "PUBLIC_INTERNET"`},
		{`
  action                   = "allow"
  traffic_type             = "PUBLIC_INTERNET"
  destination_geolocations = ["PL"]`, `destination_geolocations is only supported with action "block"`},
		{`
  action                           = "block"
  traffic_type                     = "PUBLIC_INTERNET"
  destination_geolocations         = ["PL"]
  destination_content_category_ids = [72]`, `destination_geolocations cannot be combined with other destinations`},
		{`
  action                   = "block"
  traffic_type             = "PUBLIC_INTERNET"
  source_geolocations      = ["DE"]
  destination_geolocations = ["PL"]`, `destination_geolocations cannot be combined with source_geolocations`},
		{`
  action                         = "block"
  traffic_type                   = "PUBLIC_INTERNET"
  source_internet_literals       = [{ ips = ["192.168.0.1"], ports = ["443"] }]
  destination_service_object_ids = [1]`, `source_internet_literals cannot be combined with destination service objects or service object groups`},
		{`
  action                   = "block"
  traffic_type             = "PUBLIC_INTERNET"
  source_internet_literals = [{ ips = ["192.168.0.1"], ports = ["443"] }]
  destination_literals     = [{ ips = ["192.168.1.1"], ports = ["0-65535"], protocol = "ICMP" }]`, `destination_literals protocol "ICMP" cannot be combined with source_internet_literals`},
		{`
  action                   = "block"
  traffic_type             = "PUBLIC_INTERNET"
  source_internet_literals = [{ ips = ["192.168.0.1"], ports = ["443"] }]
  destination_literals     = [{ ips = ["192.168.1.1"], ports = ["53"], protocol = "UDP" }, { ips = ["192.168.1.1"], ports = ["443"], protocol = "TCP" }]`, `destination_literals cannot have more than one item with source_internet_literals`},
		{`
  action              = "allow"
  traffic_type         = "PRIVATE_NETWORK"
  source_identity_ids = [1]
  source_operator     = "AND"`, `source_operator "AND" requires at least two of`},
		{`
  action              = "block"
  traffic_type         = "PUBLIC_INTERNET"
  source_identity_ids = [1]
  source_geolocations = ["PL"]`, `source_geolocations combined with source identities requires source_operator "AND"`},
		{`
  action              = "allow"
  traffic_type         = "PRIVATE_NETWORK"
  source_identity_ids = []`, `source_identity_ids must not be empty`},
		{`
  action                               = "block"
  traffic_type                         = "PUBLIC_INTERNET"
  destination_application_list_ids     = [1]
  destination_application_category_ids = [2]`, `destination_application_category_ids cannot be combined with destination_application_list_ids`},
		{`
  action                           = "block"
  traffic_type                     = "PUBLIC_INTERNET"
  destination_application_list_ids = [1]
  destination_application_ids      = [2]`, `destination_application_ids cannot be combined with destination_application_list_ids`},
		{`
  action                                = "block"
  traffic_type                          = "PUBLIC_INTERNET"
  destination_content_category_list_ids = [1]
  destination_content_category_ids      = [2]`, `destination_content_category_ids cannot be combined with destination_content_category_list_ids`},
	}
	var testSteps []resource.TestStep
	for _, s := range steps {
		testSteps = append(testSteps, resource.TestStep{
			Config:   config(s.attributes),
			PlanOnly: true,
			// Terraform wraps long error messages
			ExpectError: regexp.MustCompile(strings.ReplaceAll(regexp.QuoteMeta(s.err), " ", `\s+`)),
		})
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    testSteps,
	})
}
