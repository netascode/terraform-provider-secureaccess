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
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc

func TestAccSecureAccessDestinationList(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_destination_list.test", "name", "my_destination_list"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_destination_list.test", "destinations.0.destination", "example.com"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_destination_list.test", "destinations.0.comment", "This is a comment for the destination."))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_destination_list.test", "destinations.0.type", "domain"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccSecureAccessDestinationListConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessDestinationListConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_destination_list.test",
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

func testAccSecureAccessDestinationListConfig_minimum() string {
	config := `resource "secureaccess_destination_list" "test" {` + "\n"
	config += `	name = "my_destination_list"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessDestinationListConfig_all() string {
	config := `resource "secureaccess_destination_list" "test" {` + "\n"
	config += `	name = "my_destination_list"` + "\n"
	config += `	destinations = [{` + "\n"
	config += `		destination = "example.com"` + "\n"
	config += `		comment = "This is a comment for the destination."` + "\n"
	config += `		type = "domain"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

func TestAccSecureAccessDestinationList_Sequential(t *testing.T) {

	step_01 := `resource "secureaccess_destination_list" "test" {` + "\n" +
		`	name = "my_destination_list"` + "\n" +
		`	destinations = [{` + "\n" +
		`           destination = "example.com",` + "\n" +
		`           comment = "This is a comment for the destination.",` + "\n" +
		`			type = "domain",` + "\n" +
		`		}]` + "\n" +
		`}` + "\n"

	var checks_step01 []resource.TestCheckFunc
	checks_step01 = append(checks_step01, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "id"))
	checks_step01 = append(checks_step01, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.0.id"))

	step_02 := `resource "secureaccess_destination_list" "test" {` + "\n" +
		`	name = "my_destination_list"` + "\n" +
		`	destinations = [{` + "\n" +
		`			destination = "example.com",` + "\n" +
		`			comment = "This is a comment for the destination.",` + "\n" +
		`			type = "domain",` + "\n" +
		`		}, {` + "\n" +
		`			destination = "example2.com",` + "\n" +
		`			comment = "This is a comment for 2nd destination.",` + "\n" +
		`			type = "domain",` + "\n" +
		`		}, {` + "\n" +
		`			comment = "Second warning url managed by TF",` + "\n" +
		`			destination = "127.0.0.2",` + "\n" +
		`		}, {` + "\n" +
		`			comment = "Second warning url managed by TF",` + "\n" +
		`			destination = "http://foo.bar/blockwarn",` + "\n" +
		`		}]` + "\n" +
		`}` + "\n"

	var checks_step02 []resource.TestCheckFunc
	checks_step02 = append(checks_step02, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "id"))
	checks_step02 = append(checks_step02, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.0.id"))
	checks_step02 = append(checks_step02, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.1.id"))
	checks_step02 = append(checks_step02, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.2.id"))
	checks_step02 = append(checks_step02, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.3.id"))

	step_03 := `resource "secureaccess_destination_list" "test" {` + "\n" +
		`	name = "my_destination_list"` + "\n" +
		`	destinations = [{` + "\n" +
		`			destination = "example.com",` + "\n" +
		`			comment = "This is a comment for the destination.",` + "\n" +
		`			type = "domain",` + "\n" +
		`		}, {` + "\n" +
		`			destination = "example2.com",` + "\n" +
		`			comment = "This is a comment for 2nd destination.",` + "\n" +
		`			type = "domain",` + "\n" +
		`		}, {` + "\n" +
		`			comment = "Second warning url managed by TF",` + "\n" +
		`			destination = "127.0.0.89",` + "\n" +
		`		}, {` + "\n" +
		`			comment = "Second warning url managed by TF",` + "\n" +
		`			destination = "http://foo.bar/blockwarn",` + "\n" +
		`		}]` + "\n" +
		`}` + "\n"

	var checks_step03 []resource.TestCheckFunc
	checks_step03 = append(checks_step03, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "id"))
	checks_step03 = append(checks_step03, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.0.id"))
	checks_step03 = append(checks_step03, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.1.id"))
	checks_step03 = append(checks_step03, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.2.id"))
	checks_step03 = append(checks_step03, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.3.id"))

	step_04 := `resource "secureaccess_destination_list" "test" {` + "\n" +
		`	name = "my_destination_list"` + "\n" +
		`	destinations = [{` + "\n" +
		`			destination = "example.com",` + "\n" +
		`			comment = "This is a comment for the destination.",` + "\n" +
		`			type = "domain",` + "\n" +
		`		}, {` + "\n" +
		`			destination = "example2.com",` + "\n" +
		`			comment = "This is a comment for 2nd destination.",` + "\n" +
		`			type = "domain",` + "\n" +
		`		}, {` + "\n" +
		`			comment = "Updated warning comment",` + "\n" +
		`			destination = "127.0.0.189",` + "\n" +
		`		}]` + "\n" +
		`}` + "\n"

	var checks_step04 []resource.TestCheckFunc
	checks_step04 = append(checks_step04, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "id"))
	checks_step04 = append(checks_step04, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.0.id"))
	checks_step04 = append(checks_step04, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.1.id"))
	checks_step04 = append(checks_step04, resource.TestCheckResourceAttrSet("secureaccess_destination_list.test", "destinations.2.id"))

	step_05 := `resource "secureaccess_destination_list" "test" {` + "\n" +
		`	name = "my_destination_list"` + "\n" +
		`	destinations = [{` + "\n" +
		`			destination = "example.com",` + "\n" +
		`			comment = "This is updated comment for the destination.",` + "\n" +
		`			type = "domain",` + "\n" +
		`		}, {` + "\n" +
		`			destination = "example2.com",` + "\n" +
		`			comment = "This is a comment for 2nd destination.",` + "\n" +
		`			type = "domain",` + "\n" +
		`		}, {` + "\n" +
		`			comment = "Updated warning comment",` + "\n" +
		`			destination = "127.0.0.189",` + "\n" +
		`		}]` + "\n" +
		`}` + "\n"

	steps := []resource.TestStep{{
		Config: step_01,
		Check:  resource.ComposeTestCheckFunc(checks_step01...),
	}, {
		Config: step_02,
		Check:  resource.ComposeTestCheckFunc(checks_step02...),
	}, {
		Config: step_03,
		Check:  resource.ComposeTestCheckFunc(checks_step03...),
	}, {
		Config: step_04,
		Check:  resource.ComposeTestCheckFunc(checks_step04...),
	}, {
		Config: step_05,
	}}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func TestAccSecureAccessDestinationList_600destinations(t *testing.T) {
	step_01 := `resource "secureaccess_destination_list" "test_dl1_600" { ` + "\n" +
		`	name = "destination_list_600_destinations"` + "\n" +
		`	destinations = [` + "\n"
	for i := range 600 {
		step_01 += fmt.Sprintf(`	{ destination = "dest_%03d.example.com" },`+"\n", i)
	}
	step_01 += `	]` + "\n" + `}` + "\n"

	step_02 := `resource "secureaccess_destination_list" "test_dl1_600" { ` + "\n" +
		`	name = "destination_list_600_destinations"` + "\n" +
		`	destinations = [` + "\n"
	for i := 0; i < 600; i += 2 {
		step_02 += fmt.Sprintf(`	{ destination = "dest_%03d.example.com" },`+"\n", i)
	}
	step_02 += `	]` + "\n" + `}` + "\n"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: step_01,
			},
			{
				Config: step_02,
			},
			{
				Config: step_01,
			},
		},
	})
}
