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

func TestAccSecureAccessDLPRealTimeRule(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "name", "my_dlp_real_time_rule"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "description", "My real time rule"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "severity", "CRITICAL"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_size_configuration_type", "RANGE"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_size_from", "1024"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_size_from_unit", "KB"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_size_to", "51200"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_size_to_unit", "KB"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_types_included.0.id", "250"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_types_included.0.name", "cdr"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_type_groups_included.0.id", "6"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_types_excluded.0.id", "156"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_types_excluded.0.name", "dbf"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_type_groups_excluded.0.id", "2"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "file_type_groups_excluded.0.name", "Compressed files"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "unclassified_files_scan", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "identities_included.0.type_id", "7"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "identities_excluded.0.type_id", "40"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "all_destinations_scope", "FILES_AND_FORMS"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "action", "BLOCK"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "override_block", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "notification_popup", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "email_notify_actor", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("secureaccess_dlp_real_time_rule.test", "email_notify_actors_manager", "true"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccSecureAccessDLPRealTimeRulePrerequisitesConfig + testAccSecureAccessDLPRealTimeRuleConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccSecureAccessDLPRealTimeRulePrerequisitesConfig + testAccSecureAccessDLPRealTimeRuleConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName: "secureaccess_dlp_real_time_rule.test",
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

const testAccSecureAccessDLPRealTimeRulePrerequisitesConfig = `
data "secureaccess_dlp_builtin_classification" "test" {
  name = "Built-in PCI Classification"
}
`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal

func testAccSecureAccessDLPRealTimeRuleConfig_minimum() string {
	config := `resource "secureaccess_dlp_real_time_rule" "test" {` + "\n"
	config += `	name = "my_dlp_real_time_rule"` + "\n"
	config += `	severity = "CRITICAL"` + "\n"
	config += `	enabled = true` + "\n"
	config += `	data_classification_scannable_contexts = ["CONTENT"]` + "\n"
	config += `	data_classifications = [{` + "\n"
	config += `		id = data.secureaccess_dlp_builtin_classification.test.id` + "\n"
	config += `	}]` + "\n"
	config += `	identities_included = [{` + "\n"
	config += `		type_id = 7` + "\n"
	config += `	}]` + "\n"
	config += `	all_destinations_scope = "FILES_AND_FORMS"` + "\n"
	config += `	action = "BLOCK"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll

func testAccSecureAccessDLPRealTimeRuleConfig_all() string {
	config := `resource "secureaccess_dlp_real_time_rule" "test" {` + "\n"
	config += `	name = "my_dlp_real_time_rule"` + "\n"
	config += `	description = "My real time rule"` + "\n"
	config += `	severity = "CRITICAL"` + "\n"
	config += `	enabled = true` + "\n"
	config += `	data_classification_scannable_contexts = ["CONTENT"]` + "\n"
	config += `	data_classifications = [{` + "\n"
	config += `		id = data.secureaccess_dlp_builtin_classification.test.id` + "\n"
	config += `		name = data.secureaccess_dlp_builtin_classification.test.name` + "\n"
	config += `	}]` + "\n"
	config += `	file_labels = ["Confidential"]` + "\n"
	config += `	file_size_configuration_type = "RANGE"` + "\n"
	config += `	file_size_from = 1024` + "\n"
	config += `	file_size_from_unit = "KB"` + "\n"
	config += `	file_size_to = 51200` + "\n"
	config += `	file_size_to_unit = "KB"` + "\n"
	config += `	file_types_included = [{` + "\n"
	config += `		id = 250` + "\n"
	config += `		name = "cdr"` + "\n"
	config += `	}]` + "\n"
	config += `	file_type_groups_included = [{` + "\n"
	config += `		id = 6` + "\n"
	config += `	}]` + "\n"
	config += `	file_types_excluded = [{` + "\n"
	config += `		id = 156` + "\n"
	config += `		name = "dbf"` + "\n"
	config += `	}]` + "\n"
	config += `	file_type_groups_excluded = [{` + "\n"
	config += `		id = 2` + "\n"
	config += `		name = "Compressed files"` + "\n"
	config += `	}]` + "\n"
	config += `	unclassified_files_scan = false` + "\n"
	config += `	identities_included = [{` + "\n"
	config += `		type_id = 7` + "\n"
	config += `	}]` + "\n"
	config += `	identities_excluded = [{` + "\n"
	config += `		type_id = 40` + "\n"
	config += `	}]` + "\n"
	config += `	all_destinations_scope = "FILES_AND_FORMS"` + "\n"
	config += `	action = "BLOCK"` + "\n"
	config += `	override_block = false` + "\n"
	config += `	notification_popup = true` + "\n"
	config += `	email_notify_actor = true` + "\n"
	config += `	email_notify_actors_manager = true` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll
