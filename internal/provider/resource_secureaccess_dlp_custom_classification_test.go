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

func TestAccSecureAccessDLPCustomClassification(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_custom_classification.test", "name", "my_dlp_custom_classification"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_custom_classification.test", "description", "My custom classification for sensitive data."))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_custom_classification.test", "evaluation_strategy", "any"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccSecureAccessDLPCustomClassificationPrerequisitesConfig + testAccSecureAccessDLPCustomClassificationConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessDLPCustomClassificationPrerequisitesConfig + testAccSecureAccessDLPCustomClassificationConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_dlp_custom_classification.test",
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

const testAccSecureAccessDLPCustomClassificationPrerequisitesConfig = `
resource "secureaccess_dlp_custom_data_identifier" "test1" {
  name = "my_dlp_custom_classification_1"
  threshold_ranges = [{
    minimum = 1
    maximum = 10
  }]
  patterns = ["pattern"]
}

resource "secureaccess_dlp_custom_data_identifier" "test2" {
  name = "my_dlp_custom_classification_2"
  threshold_ranges = [{
    minimum = 1
    maximum = 10
  }]
  patterns = ["pattern"]
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal

func testAccSecureAccessDLPCustomClassificationConfig_minimum() string {
	config := `resource "secureaccess_dlp_custom_classification" "test" {` + "\n"
	config += `	name = "my_dlp_custom_classification"` + "\n"
	config += `	evaluation_strategy = "any"` + "\n"
	config += `	included_data_identifiers = [{` + "\n"
	config += `		id = secureaccess_dlp_custom_data_identifier.test1.id` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessDLPCustomClassificationConfig_all() string {
	config := `resource "secureaccess_dlp_custom_classification" "test" {` + "\n"
	config += `	name = "my_dlp_custom_classification"` + "\n"
	config += `	description = "My custom classification for sensitive data."` + "\n"
	config += `	evaluation_strategy = "any"` + "\n"
	config += `	included_data_identifiers = [{` + "\n"
	config += `		id = secureaccess_dlp_custom_data_identifier.test1.id` + "\n"
	config += `	}]` + "\n"
	config += `	excluded_data_identifiers = [{` + "\n"
	config += `		id = secureaccess_dlp_custom_data_identifier.test2.id` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll
