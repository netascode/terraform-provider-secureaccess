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
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

// threatCategoryOptionState is "true", "false" or "null" for the value of an option.
func threatCategoryOptionState(b types.Bool) string {
	switch {
	case b.IsNull():
		return "null"
	case b.ValueBool():
		return "true"
	}
	return "false"
}

// threatCategoryOptionOf returns what threatCategoryOption gives for the JSON value of categoryBits (the value is given as JSON
// text, so that a string and other types can be tested). The numbers in the tests are built with threatCategoryHex, which is
// written independently of the code under test.
func threatCategoryOptionOf(rawValue string, n int) string {
	return threatCategoryOptionState(threatCategoryOption(gjson.Get(`{"categoryBits":`+rawValue+`}`, "categoryBits"), n))
}

func TestThreatCategoryOption(t *testing.T) {
	// the numbers that were seen on a live organization (see the header of model_secureaccess_threat_category_test.go)
	t.Run("the values seen live", func(t *testing.T) {
		for _, tt := range []struct {
			hex  string
			bits []int
		}{
			{"5000000000002000000000000000", []int{108, 61, 110}},
			{"4000000000000000000000000000", []int{110}},
			{"40000000003000000000002000000000000000", []int{109, 108, 61, 150}},
			{"d000000000000000", []int{60, 62, 63}},
			{"0", nil},
		} {
			if got := threatCategoryHex(tt.bits...); got != tt.hex {
				t.Fatalf("test helper: the bits %v give %s, want %s", tt.bits, got, tt.hex)
			}
			set := map[int]bool{}
			for _, b := range tt.bits {
				set[b] = true
			}
			for _, n := range []int{0, 60, 61, 62, 63, 64, 108, 109, 110, 150} {
				want := strconv.FormatBool(set[n])
				if got := threatCategoryOptionOf(`"`+tt.hex+`"`, n); got != want {
					t.Errorf("bit %d of %s = %s, want %s", n, tt.hex, got, want)
				}
			}
		}
	})

	t.Run("every bit alone sets exactly that bit, up to the bits that do not fit into 64 bits", func(t *testing.T) {
		for _, set := range []int{0, 1, 3, 4, 7, 8, 60, 61, 62, 63, 64, 65, 108, 109, 110, 111, 150, 151, 255} {
			raw := `"` + threatCategoryHex(set) + `"`
			for n := 0; n <= 260; n++ {
				want := "false"
				if n == set {
					want = "true"
				}
				if got := threatCategoryOptionOf(raw, n); got != want {
					t.Fatalf("bit %d of the number with only bit %d set (%s) = %s, want %s", n, set, raw, got, want)
				}
			}
		}
	})

	t.Run("the bit 150 needs more than 64 bits", func(t *testing.T) {
		hex := threatCategoryHex(150)
		if len(hex) != 38 {
			t.Fatalf("test helper: bit 150 gives %d hexadecimal digits, want 38", len(hex))
		}
		if got := threatCategoryOptionOf(`"`+hex+`"`, 150); got != "true" {
			t.Errorf("bit 150 = %s", got)
		}
		// the same number cut to the lower 64 bits, as a uint64 parser would see it
		if got := threatCategoryOptionOf(`"`+hex[len(hex)-16:]+`"`, 150); got != "false" {
			t.Errorf("bit 150 of a number without it = %s", got)
		}
	})

	t.Run("nothing set", func(t *testing.T) {
		for _, raw := range []string{`"0"`, `"00"`, `"0000000000000000000000000000"`, `""`} {
			for _, n := range []int{0, 61, 108, 150} {
				if got := threatCategoryOptionOf(raw, n); got != "false" {
					t.Errorf("bit %d of %s = %s, want false", n, raw, got)
				}
			}
		}
	})

	t.Run("leading zeros and the case of the digits do not matter", func(t *testing.T) {
		for _, raw := range []string{`"0000d000000000000000"`, `"00d000000000000000"`, `"D000000000000000"`, `"d000000000000000"`} {
			for n, want := range map[int]string{60: "true", 62: "true", 63: "true", 61: "false", 64: "false"} {
				if got := threatCategoryOptionOf(raw, n); got != want {
					t.Errorf("bit %d of %s = %s, want %s", n, raw, got, want)
				}
			}
		}
		// the digit a is 1010 and the digit A is the same number
		for _, raw := range []string{`"a"`, `"A"`} {
			for n, want := range map[int]string{0: "false", 1: "true", 2: "false", 3: "true"} {
				if got := threatCategoryOptionOf(raw, n); got != want {
					t.Errorf("bit %d of %s = %s, want %s", n, raw, got, want)
				}
			}
		}
	})

	// the YAML says that every other bit is ignored: a number with bits far above the mapped ones (256 and more, more than 64 hexadecimal
	// digits) must still give the options, not null
	t.Run("bits far above the mapped ones are ignored", func(t *testing.T) {
		for _, extra := range []int{256, 300, 1000, 4096} {
			raw := `"` + threatCategoryHex(extra, bitNewlySeenDomains, bitCryptomining) + `"`
			for n, want := range map[int]string{108: "true", 150: "true", 109: "false", 61: "false", 110: "false"} {
				if got := threatCategoryOptionOf(raw, n); got != want {
					t.Errorf("bit %d of the number with the bits %d, 108 and 150 = %s, want %s", n, extra, got, want)
				}
			}
		}
	})

	t.Run("a bit beyond the number is not set", func(t *testing.T) {
		for _, n := range []int{4, 64, 1000, 100000} {
			if got := threatCategoryOptionOf(`"f"`, n); got != "false" {
				t.Errorf("bit %d of f = %s, want false", n, got)
			}
		}
	})

	// A value that is not a plain hexadecimal number must not read as "no bit set" and must not read as "all bits set": the
	// option is null. big.Int.SetString alone would accept a sign, and a negative number reads every bit as 1.
	t.Run("a value that is not a hexadecimal string is null", func(t *testing.T) {
		for _, raw := range []string{
			`null`, `5`, `5.5`, `true`, `false`, `[]`, `{}`, `["5"]`,
			`"zz"`, `"g"`, `"-1"`, `"+5"`, `"0x5"`, `"0X5"`, `"5_0"`, `" 5"`, `"5 "`, `"5\n"`, `"5.0"`, `"٥"`,
		} {
			for _, n := range []int{0, 61, 108, 150} {
				if got := threatCategoryOptionOf(raw, n); got != "null" {
					t.Errorf("bit %d of %s = %s, want null", n, raw, got)
				}
			}
		}
	})

	t.Run("a key that does not exist is null", func(t *testing.T) {
		if got := threatCategoryOptionState(threatCategoryOption(gjson.Get(`{"name":"x"}`, "categoryBits"), 108)); got != "null" {
			t.Errorf("a missing key gives %s, want null", got)
		}
	})
}

func TestThreatCategoryBitSet(t *testing.T) {
	tests := []struct {
		s       string
		n       int
		set, ok bool
	}{
		{"", 0, false, true},
		{"0", 0, false, true},
		{"1", 0, true, true},
		{"2", 0, false, true},
		{"2", 1, true, true},
		{"80", 7, true, true},
		{"80", 6, false, true},
		{"fF", 7, true, true},
		{"zz", 0, false, false},
		{"-1", 150, false, false},
		{"+1", 0, false, false},
		{"0x1", 0, false, false},
		{"1_0", 4, false, false},
		{" 1", 0, false, false},
		{"5", -1, false, false},
		{"", -1, false, false},
		{"zz", -1, false, false},
	}
	for _, tt := range tests {
		set, ok := threatCategoryBitSet(tt.s, tt.n)
		if set != tt.set || ok != tt.ok {
			t.Errorf("threatCategoryBitSet(%q, %d) = %t, %t, want %t, %t", tt.s, tt.n, set, ok, tt.set, tt.ok)
		}
	}
}

// Every hexadecimal digit, in both cases, sets exactly the bits of its value. The digits 8 and 9 and the first and the last digit of
// each of the three ranges (0-9, a-f, A-F) are where an off-by-one in the check of the digits would show. The value of a digit
// comes from strconv, not from threatCategoryBitSet.
func TestThreatCategoryBitSetEveryDigit(t *testing.T) {
	for _, d := range "0123456789abcdefABCDEF" {
		value, err := strconv.ParseUint(string(d), 16, 8)
		if err != nil {
			t.Fatalf("test helper: the digit %q: %v", d, err)
		}
		for k := 0; k < 8; k++ {
			want := k < 4 && value>>k&1 == 1
			if set, ok := threatCategoryBitSet(string(d), k); !ok || set != want {
				t.Errorf("threatCategoryBitSet(%q, %d) = %t, %t, want %t, true", string(d), k, set, ok, want)
			}
			// the same digit as the second digit from the end holds the bits 4 to 7
			wantHigh := k >= 4 && k < 8 && value>>(k-4)&1 == 1
			if set, ok := threatCategoryBitSet(string(d)+"0", k); !ok || set != wantHigh {
				t.Errorf("threatCategoryBitSet(%q, %d) = %t, %t, want %t, true", string(d)+"0", k, set, ok, wantHigh)
			}
		}
	}
	// the characters just outside of the three ranges are not digits
	for _, c := range []string{"/", ":", "`", "g", "@", "G"} {
		if set, ok := threatCategoryBitSet(c, 0); ok || set {
			t.Errorf("threatCategoryBitSet(%q, 0) = %t, %t, want false, false", c, set, ok)
		}
	}
}

// adjustRead sets the five options and nothing else, and does not keep a value that is already in the model.
func TestThreatCategoryAdjustRead(t *testing.T) {
	ctx := context.Background()

	t.Run("it sets the five options and leaves the other attributes alone", func(t *testing.T) {
		data := ThreatCategory{
			Id:                         types.StringValue("keep id"),
			Name:                       types.StringValue("keep name"),
			Malware:                    types.BoolValue(true),
			CommandAndControlCallbacks: types.BoolValue(false),
			PhishingAttacks:            types.BoolValue(true),
		}
		res := gjson.Parse(`{"malwareProtection": false, "categoryBits": "` + threatCategoryHex(109, 108, 61, 150) + `"}`)
		data.adjustRead(ctx, res, nil)

		want := map[string]string{
			"potentially_harmful_domains": "true", "newly_seen_domains": "true", "dynamic_dns": "true",
			"dns_tunneling_vpn": "false", "cryptomining": "true",
			"malware": "true", "command_and_control_callbacks": "false", "phishing_attacks": "true",
		}
		for attr, v := range threatCategoryBoolsOf(data) {
			if got := threatCategoryOptionState(v); got != want[attr] {
				t.Errorf("%s = %s, want %s", attr, got, want[attr])
			}
		}
		if data.Id.ValueString() != "keep id" || data.Name.ValueString() != "keep name" {
			t.Errorf("id %s and name %s must not change", data.Id, data.Name)
		}
	})

	t.Run("a value that was in the model does not survive a missing or invalid categoryBits", func(t *testing.T) {
		for _, body := range []string{`{}`, `{"categoryBits": null}`, `{"categoryBits": 5}`, `{"categoryBits": "zz"}`} {
			data := ThreatCategory{
				PotentiallyHarmfulDomains: types.BoolValue(true),
				NewlySeenDomains:          types.BoolValue(true),
				DynamicDns:                types.BoolValue(true),
				DnsTunnelingVpn:           types.BoolValue(true),
				Cryptomining:              types.BoolValue(true),
			}
			data.adjustRead(ctx, gjson.Parse(body), nil)
			for _, v := range []types.Bool{data.PotentiallyHarmfulDomains, data.NewlySeenDomains, data.DynamicDns, data.DnsTunnelingVpn, data.Cryptomining} {
				if !v.IsNull() {
					t.Errorf("%s: an option is %s, want null", body, v)
				}
			}
		}
	})

	t.Run("a value that was in the model does not survive when the bit is not set", func(t *testing.T) {
		data := ThreatCategory{
			PotentiallyHarmfulDomains: types.BoolValue(true),
			NewlySeenDomains:          types.BoolValue(true),
			DynamicDns:                types.BoolValue(true),
			DnsTunnelingVpn:           types.BoolValue(true),
			Cryptomining:              types.BoolValue(true),
		}
		data.adjustRead(ctx, gjson.Parse(`{"categoryBits": "0"}`), nil)
		for _, v := range []types.Bool{data.PotentiallyHarmfulDomains, data.NewlySeenDomains, data.DynamicDns, data.DnsTunnelingVpn, data.Cryptomining} {
			if v.IsNull() || v.ValueBool() {
				t.Errorf("an option is %s, want false", v)
			}
		}
	})
}
