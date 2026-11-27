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
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

// The Security Profiles API guide describes `aiSupplyChainCategories` as a string that is `0` if the AI Supply Chain Blocking is off and
// otherwise 54 hexadecimal digits, of which the first adds up the categories that are blocked and the other 53 are 0.
// A value other than `0` was never seen on a live organization, this reading is from the guide only.
const (
	aiSupplyChainMaskLength               = 54
	aiSupplyChainValueProhibitedSuppliers = 2
	aiSupplyChainValueCodeExecution       = 4
	aiSupplyChainValueCopyleftLicense     = 8
)

// aiSupplyChainOption returns whether the category with the given value is blocked in the mask, or null if the mask is missing or is
// neither `0` nor 54 hexadecimal digits. The digits after the first one are not read.
func aiSupplyChainOption(mask gjson.Result, value int) types.Bool {
	if mask.Type != gjson.String {
		return types.BoolNull()
	}
	s := mask.String()
	if s == "0" {
		return types.BoolValue(false)
	}
	if len(s) != aiSupplyChainMaskLength {
		return types.BoolNull()
	}
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return types.BoolNull()
		}
	}
	// the first digit is a single hexadecimal digit, ParseUint accepts it in both cases
	first, err := strconv.ParseUint(s[:1], 16, 8)
	if err != nil {
		return types.BoolNull()
	}
	return types.BoolValue(int(first)&value != 0)
}
