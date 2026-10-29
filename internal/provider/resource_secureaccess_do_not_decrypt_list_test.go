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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc

func TestAccSecureAccessDoNotDecryptList(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "name", "my_do_not_decrypt_list"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "destinations.0.destination", "example.com"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "destinations.0.comment", "This is a comment for the destination."))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccSecureAccessDoNotDecryptListConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessDoNotDecryptListConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_do_not_decrypt_list.test",
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

func testAccSecureAccessDoNotDecryptListConfig_minimum() string {
	config := `resource "secureaccess_do_not_decrypt_list" "test" {` + "\n"
	config += `	name = "my_do_not_decrypt_list"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessDoNotDecryptListConfig_all() string {
	config := `resource "secureaccess_do_not_decrypt_list" "test" {` + "\n"
	config += `	name = "my_do_not_decrypt_list"` + "\n"
	config += `	destinations = [{` + "\n"
	config += `		destination = "example.com"` + "\n"
	config += `		comment = "This is a comment for the destination."` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

// Manual test: adding, changing and removing destinations is done with separate add/remove lists in the API,
// as it does not return the destinations, only the totals
func TestAccSecureAccessDoNotDecryptList_Sequential(t *testing.T) {
	config := func(name string, destinations ...string) string {
		cfg := `resource "secureaccess_do_not_decrypt_list" "test" {` + "\n" +
			`	name = "` + name + `"` + "\n"
		if len(destinations) > 0 {
			cfg += `	destinations = [` + "\n"
			for i := 0; i+1 < len(destinations); i += 2 {
				cfg += `		{` + "\n" +
					`			destination = "` + destinations[i] + `"` + "\n"
				if destinations[i+1] != "" {
					cfg += `			comment = "` + destinations[i+1] + `"` + "\n"
				}
				cfg += `		},` + "\n"
			}
			cfg += `	]` + "\n"
		}
		cfg += `}` + "\n"
		return cfg
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// create with a comment and without one
				Config: config("my_do_not_decrypt_list_seq", "example.com", "first comment", "10.10.10.10", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("secureaccess_do_not_decrypt_list.test", "id"),
					resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "destinations.#", "2"),
				),
			},
			{
				// change a comment, remove a destination and add two others
				Config: config("my_do_not_decrypt_list_seq", "example.com", "changed comment", "example2.com", "second", "https://foo.bar/path", "url"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "destinations.#", "3"),
				),
			},
			{
				// rename only
				Config: config("my_do_not_decrypt_list_seq_renamed", "example.com", "changed comment", "example2.com", "second", "https://foo.bar/path", "url"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "name", "my_do_not_decrypt_list_seq_renamed"),
					resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "destinations.#", "3"),
				),
			},
			{
				// remove all destinations
				Config: config("my_do_not_decrypt_list_seq_renamed"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "destinations.#", "0"),
				),
			},
		},
	})
}

// Manual test: the same destination cannot be listed twice, the config is rejected before any request is sent
func TestAccSecureAccessDoNotDecryptList_DuplicateDestination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "secureaccess_do_not_decrypt_list" "test" {` + "\n" +
					`	name = "my_do_not_decrypt_list_dup"` + "\n" +
					`	destinations = [` + "\n" +
					`		{ destination = "example.com", comment = "first" },` + "\n" +
					`		{ destination = "example.com", comment = "second" },` + "\n" +
					`	]` + "\n" +
					`}` + "\n",
				ExpectError: regexp.MustCompile(`Duplicate destination`),
			},
		},
	})
}
