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

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

// This file is not generated. It tests adjustRead(), the part of this data source that is written by hand: it reads the three AI Supply
// Chain Blocking options from the mask aiSupplyChainCategories. How the mask is decoded is tested in
// helpers_security_profile_policy_setting_test.go. The encoding of the mask (the first digit adds up 2, 4 and 8) is from the Security Profiles
// API guide only: no value other than "0" was seen on a live organization (see the YAML), so the expected values here are from the guide.
// The schema, fromBody() and Read() are generated from the YAML and are not tested for themselves. A test of the schema would repeat the YAML and
// fail whenever the YAML changes, and a test with a fake of the API only shows that Read() understands what we think the API answers, so it
// cannot notice that the API changed. fromBody() is used below only as a way to reach adjustRead(), which the generated fromBody() calls.
// Whether the data source works against a real organization is for an acceptance test. That needs a Security Profile resource, because the API
// creates a Policy Setting only for a Security Profile.

// The three AI Supply Chain Blocking options of data as three letters, in the order Prohibited Suppliers, Code Execution and Copyleft
// License: T for true, F for false and N for null.
func aiSupplyChainOptionsOf(data SecurityProfilePolicySetting) string {
	var out string
	for _, v := range []types.Bool{data.AiSupplyChainProhibitedSuppliers, data.AiSupplyChainCodeExecution, data.AiSupplyChainCopyleftLicense} {
		switch {
		case v.IsNull():
			out += "N"
		case v.ValueBool():
			out += "T"
		default:
			out += "F"
		}
	}
	return out
}

func TestSecurityProfilePolicySettingAdjustRead(t *testing.T) {
	// The other keys of the body are true: the options are read from the mask and from no other key.
	maskBody := func(mask string) string {
		return fmt.Sprintf(`{"aiSupplyChainCategories":%q,"mcpSemanticInspection":true,"intelligentProxy":true,"saml":1}`, mask)
	}

	for _, tt := range []struct{ name, body, want string }{
		{"the value 0 is everything off", maskBody("0"), "FFF"},
		{"only Prohibited Suppliers", maskBody(policySettingMask("2")), "TFF"},
		{"only Code Execution", maskBody(policySettingMask("4")), "FTF"},
		{"only Copyleft License", maskBody(policySettingMask("8")), "FFT"},
		{"Prohibited Suppliers and Code Execution", maskBody(policySettingMask("6")), "TTF"},
		{"Prohibited Suppliers and Copyleft License", maskBody(policySettingMask("A")), "TFT"},
		{"Code Execution and Copyleft License", maskBody(policySettingMask("C")), "FTT"},
		{"all three", maskBody(policySettingMask("E")), "TTT"},
		// null is not false: false would claim that nothing is blocked
		{"a missing mask is unknown", `{"mcpSemanticInspection":true,"intelligentProxy":true,"saml":1}`, "NNN"},
		{"a mask that is null is unknown", `{"aiSupplyChainCategories":null}`, "NNN"},
		{"a mask that is not 0 or 54 digits is unknown", maskBody("2"), "NNN"},
		{"a mask that is a number is unknown", `{"aiSupplyChainCategories":2}`, "NNN"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got SecurityProfilePolicySetting
			got.fromBody(context.Background(), gjson.Parse(tt.body))
			if options := aiSupplyChainOptionsOf(got); options != tt.want {
				t.Errorf("Prohibited Suppliers, Code Execution and Copyleft License = %s, want %s", options, tt.want)
			}
		})
	}

	// adjustRead() must not touch any attribute but the three options, whatever the mask is and whatever the other attributes and the other keys of
	// the body are. The whole struct is compared with the three options left out, so the test names no other attribute and needs no change when one
	// is added or removed.
	t.Run("it sets the three options and no other attribute", func(t *testing.T) {
		withoutOptions := func(d SecurityProfilePolicySetting) SecurityProfilePolicySetting {
			d.AiSupplyChainProhibitedSuppliers, d.AiSupplyChainCodeExecution, d.AiSupplyChainCopyleftLicense = types.BoolNull(), types.BoolNull(), types.BoolNull()
			return d
		}
		for _, set := range []bool{true, false} {
			for _, first := range []string{"0", "2", "E"} {
				var data SecurityProfilePolicySetting
				v := reflect.ValueOf(&data).Elem()
				for i := 0; i < v.NumField(); i++ {
					if !v.Field(i).CanSet() {
						continue
					}
					// the types that a field of another type keeps as its zero value, which DeepEqual compares too
					switch f := v.Field(i).Addr().Interface().(type) {
					case *types.Bool:
						*f = types.BoolValue(set)
					case *types.String:
						*f = types.StringValue("set")
					case *types.Int64:
						*f = types.Int64Value(7)
					}
				}
				before := data
				// the keys of the body that are not the mask have the opposite value of the attributes
				body := fmt.Sprintf(`{"aiSupplyChainCategories":%q,"mcpSemanticInspection":%t,"intelligentProxy":%t}`, policySettingMask(first), !set, !set)
				data.adjustRead(context.Background(), gjson.Parse(body), nil)

				if !reflect.DeepEqual(withoutOptions(before), withoutOptions(data)) {
					t.Errorf("mask %s, attributes set to %t: adjustRead() changed an attribute other than the three options:\nbefore %+v\nafter  %+v", first, set, before, data)
				}
			}
		}
	})
}
