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

// Manual test: the generator does not create tests for data sources without a resource (no_resource: true)

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSecureAccessRegions(t *testing.T) {
	const name = "data.secureaccess_regions.test"
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttrWith(name, "regions.#", func(value string) error {
		if n, err := strconv.Atoi(value); err != nil || n < 1 {
			return fmt.Errorf("expected at least one region, got %q", value)
		}
		return nil
	}))
	for _, attribute := range []string{"name", "region", "description", "continent"} {
		checks = append(checks, resource.TestCheckResourceAttrSet(name, "regions.0."+attribute))
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureAccessRegionsConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

func testAccDataSourceSecureAccessRegionsConfig() string {
	return `
		data "secureaccess_regions" "test" {
		}
	`
}
