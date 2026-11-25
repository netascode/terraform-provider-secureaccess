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
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

// The bits of `categoryBits` that the Web UI uses for the options of a Threat Category that the API does not return as fields.
const (
	threatCategoryBitPotentiallyHarmfulDomains = 109
	threatCategoryBitNewlySeenDomains          = 108
	threatCategoryBitDynamicDns                = 61
	threatCategoryBitDnsTunnelingVpn           = 110
	threatCategoryBitCryptomining              = 150
)

// threatCategoryOption returns whether bit n of `categoryBits` is set, or null if `categoryBits` is missing or is not a
// hexadecimal string.
func threatCategoryOption(categoryBits gjson.Result, n int) types.Bool {
	if categoryBits.Type != gjson.String {
		return types.BoolNull()
	}
	set, ok := threatCategoryBitSet(categoryBits.String(), n)
	if !ok {
		return types.BoolNull()
	}
	return types.BoolValue(set)
}

// threatCategoryBitSet reports whether bit n of the hexadecimal number s is set: bit n is `1 << n` of the whole number, so the
// last digit holds the bits 0 to 3. An empty string is the number 0, so no bit is set. ok is false if s is not a plain
// hexadecimal number: a sign, a `0x` prefix or an underscore is not accepted. ok is also false for a negative n, which big.Int
// would panic on. The numbers of the API are longer than 64 bits (the Cryptomining bit is 150), so the number is a big.Int.
func threatCategoryBitSet(s string, n int) (set, ok bool) {
	if n < 0 {
		return false, false
	}
	if s == "" {
		return false, true
	}
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false, false
		}
	}
	number, ok := new(big.Int).SetString(s, 16)
	if !ok {
		return false, false
	}
	return number.Bit(n) == 1, true
}
