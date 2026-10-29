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
	"context"
	"encoding/json"
	"os"
	"regexp"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/netascode/go-secureaccess"
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

// Manual test: adding, changing and removing destinations is done with separate add/remove lists in the API, as it does
// not return the destinations, only the totals. IP addresses and CIDR blocks are removed through the destination list
// behind the object, as the update request of the API did not find some of them.
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
	totals := testAccDoNotDecryptListTotals

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// create with and without comments: a domain, a private and a public IP address, a CIDR block
				Config: config("my_do_not_decrypt_list_seq", "example.com", "first comment", "10.10.10.10", "", "203.0.113.5", "ip comment", "192.0.2.0/24", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("secureaccess_do_not_decrypt_list.test", "id"),
					totals("4", "1", "3"),
				),
			},
			{
				// change a comment, set a comment, add a domain, remove an IP address (through the destination list)
				Config: config("my_do_not_decrypt_list_seq", "example.com", "changed comment", "example2.com", "second", "10.10.10.10", "now with a comment", "192.0.2.0/24", ""),
				Check:  totals("4", "2", "2"),
			},
			{
				// remove a comment (the destination is removed and added again), remove an IP address and a CIDR block
				Config: config("my_do_not_decrypt_list_seq", "example.com", "", "example2.com", "second"),
				Check:  totals("2", "2", "0"),
			},
			{
				// rename only
				Config: config("my_do_not_decrypt_list_seq_renamed", "example.com", "", "example2.com", "second"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "name", "my_do_not_decrypt_list_seq_renamed"),
					totals("2", "2", "0"),
				),
			},
			{
				// remove all destinations
				Config: config("my_do_not_decrypt_list_seq_renamed"),
				Check:  totals("0", "0", "0"),
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

// Manual test: destinations that the API cannot store as they are are rejected before any request is sent
func TestAccSecureAccessDoNotDecryptList_InvalidDestination(t *testing.T) {
	config := func(destination string) string {
		return `resource "secureaccess_do_not_decrypt_list" "test" {` + "\n" +
			`	name = "my_do_not_decrypt_list_invalid"` + "\n" +
			`	destinations = [{ destination = "` + destination + `" }]` + "\n" +
			`}` + "\n"
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("https://foo.bar/path"), ExpectError: regexp.MustCompile(`Invalid destination`)},
			{Config: config("Example.com"), ExpectError: regexp.MustCompile(`Invalid destination`)},
		},
	})
}

func testAccDoNotDecryptListTotals(destinations, domains, ips string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "destinations.#", destinations),
		resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "total_domains", domains),
		resource.TestCheckResourceAttr("secureaccess_do_not_decrypt_list.test", "total_ips", ips),
	)
}

// removeDestinationsOutsideTerraform removes destinations from a Do Not Decrypt List the way someone in the portal would:
// through the destination list behind the object.
func removeDestinationsOutsideTerraform(t *testing.T, listName string, destinations ...string) func() {
	return func() {
		url := os.Getenv("SECUREACCESS_URL")
		if url == "" {
			url = "https://api.sse.cisco.com"
		}
		client, err := secureaccess.NewClient(url, os.Getenv("SECUREACCESS_API_KEY"), os.Getenv("SECUREACCESS_API_KEY_SECRET"))
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()
		lists, err := client.GetAll(ctx, "/policies/v2/doNotDecryptLists")
		if err != nil {
			t.Fatal(err)
		}
		listId := ""
		for _, list := range lists.Array() {
			if list.Get("name").String() == listName {
				listId = list.Get("exceptiondomainlistId").String()
			}
		}
		if listId == "" {
			t.Fatalf("list %q not found", listName)
		}
		path := "/policies/v2/destinationlists/" + listId + "/destinations"
		items, err := client.GetAll(ctx, path)
		if err != nil {
			t.Fatal(err)
		}
		var ids []int64
		for _, item := range items.Array() {
			if slices.Contains(destinations, item.Get("destination").String()) {
				ids = append(ids, item.Get("id").Int())
			}
		}
		if len(ids) != len(destinations) {
			t.Fatalf("found %d of the %d destinations %v in list %s", len(ids), len(destinations), destinations, listId)
		}
		body, _ := json.Marshal(ids)
		if res, err := client.DeleteWithBody(ctx, path+"/remove", string(body)); err != nil {
			t.Fatalf("%s, %s", err, res.String())
		}
	}
}

// Manual test: destinations that were removed outside of Terraform (the API does not return them, so Terraform cannot see
// it) can still be removed from the configuration: removing a destination that does not exist is an error of the API,
// which the provider checks against the destination list and treats as removed.
func TestAccSecureAccessDoNotDecryptList_RemovedOutsideTerraform(t *testing.T) {
	config := func(destinations ...string) string {
		cfg := `resource "secureaccess_do_not_decrypt_list" "test" {` + "\n" +
			`	name = "my_do_not_decrypt_list_gone"` + "\n" +
			`	destinations = [` + "\n"
		for _, destination := range destinations {
			cfg += `		{ destination = "` + destination + `" },` + "\n"
		}
		return cfg + `	]` + "\n" + `}` + "\n"
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("a.example.com", "b.example.com", "10.10.10.10"),
				Check:  testAccDoNotDecryptListTotals("3", "2", "1"),
			},
			{
				PreConfig: removeDestinationsOutsideTerraform(t, "my_do_not_decrypt_list_gone", "b.example.com", "10.10.10.10"),
				Config:    config("a.example.com"),
				Check:     testAccDoNotDecryptListTotals("1", "1", "0"),
			},
		},
	})
}
