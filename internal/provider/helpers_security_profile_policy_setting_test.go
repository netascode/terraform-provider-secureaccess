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
	"strings"
	"testing"

	"github.com/tidwall/gjson"
)

// policySettingMask is a mask of 54 digits, with the given first digit and 53 zeros.
func policySettingMask(first string) string { return first + strings.Repeat("0", 53) }

// The mask of aiSupplyChainCategories is read by the first digit only when it is "0" or 54 hexadecimal digits, anything else is
// unknown and gives null for the three options, not false: false would claim that nothing is blocked.
func TestAiSupplyChainOption(t *testing.T) {
	pattern := func(raw string) string {
		mask := gjson.Get(`{"m":`+raw+`}`, "m")
		var out string
		for _, value := range []int{aiSupplyChainValueProhibitedSuppliers, aiSupplyChainValueCodeExecution, aiSupplyChainValueCopyleftLicense} {
			switch v := aiSupplyChainOption(mask, value); {
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
	quote := strconv.Quote

	t.Run("0 is off", func(t *testing.T) {
		if got := pattern(`"0"`); got != "FFF" {
			t.Errorf(`"0" = %s, want FFF`, got)
		}
	})

	t.Run("every value of the first digit gives its categories", func(t *testing.T) {
		// the categories are the bits 2, 4 and 8 of the digit, the bit 1 is not one of them
		for digit := 0; digit < 16; digit++ {
			hex := strconv.FormatInt(int64(digit), 16)
			want := ""
			for _, bit := range []int{2, 4, 8} {
				if digit&bit != 0 {
					want += "T"
				} else {
					want += "F"
				}
			}
			for _, d := range []string{hex, strings.ToUpper(hex)} {
				if got := pattern(quote(policySettingMask(d))); got != want {
					t.Errorf("the first digit %s = %s, want %s", d, got, want)
				}
			}
		}
	})

	t.Run("the digits after the first one are not read", func(t *testing.T) {
		for _, tail := range []string{"1", "F", "A", "9"} {
			if got := pattern(quote(policySettingMask("4")[:53] + tail)); got != "FTF" {
				t.Errorf("the last digit %s = %s, want FTF", tail, got)
			}
		}
		if got := pattern(quote("8" + strings.Repeat("F", 53))); got != "FFT" {
			t.Errorf("all other digits F = %s, want FFT", got)
		}
	})

	t.Run("a value that is not 0 or 54 hexadecimal digits is unknown", func(t *testing.T) {
		for name, raw := range map[string]string{
			"null":                             `null`,
			"a number 0":                       `0`,
			"a number 2":                       `2`,
			"true":                             `true`,
			"a list":                           `["2"]`,
			"an object":                        `{}`,
			"an empty string":                  `""`,
			"00":                               `"00"`,
			"2":                                `"2"`,
			"53 digits":                        quote(policySettingMask("2")[:53]),
			"55 digits":                        quote(policySettingMask("2") + "0"),
			"a letter that is not a hex digit": quote("G" + strings.Repeat("0", 53)),
			"a sign":                           quote("+" + strings.Repeat("0", 53)),
			"a space":                          quote(" " + strings.Repeat("0", 53)),
			"0x":                               quote("0x" + strings.Repeat("0", 52)),
			"a non-hex digit after the first":  quote("2" + strings.Repeat("0", 52) + "g"),
			"a unicode digit":                  quote("٢" + strings.Repeat("0", 53)),
		} {
			if got := pattern(raw); got != "NNN" {
				t.Errorf("%s (%s) = %s, want NNN", name, raw, got)
			}
		}
		if got := aiSupplyChainOption(gjson.Result{}, aiSupplyChainValueCodeExecution); !got.IsNull() {
			t.Errorf("a missing key = %s, want null", got)
		}
	})

	// A valid value with whitespace around it is not trimmed, it is unknown like any other value that is not "0" or 54 hexadecimal digits.
	t.Run("whitespace around a valid value makes it unknown", func(t *testing.T) {
		for name, raw := range map[string]string{
			"a space before 0":        quote(" 0"),
			"a space after 0":         quote("0 "),
			"a new line before 0":     quote("\n0"),
			"a space before a mask":   quote(" " + policySettingMask("2")),
			"a new line after a mask": quote(policySettingMask("2") + "\n"),
			"a tab after a mask":      quote(policySettingMask("4") + "\t"),
		} {
			if got := pattern(raw); got != "NNN" {
				t.Errorf("%s (%s) = %s, want NNN", name, raw, got)
			}
		}
	})

	// The characters next to the ranges 0-9, A-F and a-f in the ASCII table and two that are in no range (_ and x), in every position: a
	// comparison that is off by one or a check that skips a position would accept one of them. The first digit is rejected by ParseUint
	// anyway, so the later positions are the ones that count.
	t.Run("the characters next to the hexadecimal ranges are not digits, in any position", func(t *testing.T) {
		for _, bad := range []string{"/", ":", "@", "G", "`", "g", "_", "x"} {
			for pos := 0; pos < aiSupplyChainMaskLength; pos++ {
				mask := []byte(policySettingMask("2"))
				mask[pos] = bad[0]
				if got := pattern(quote(string(mask))); got != "NNN" {
					t.Errorf("%q at the position %d = %s, want NNN", bad, pos, got)
				}
			}
		}
	})
}
