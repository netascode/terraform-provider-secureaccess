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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type ThreatCategory struct {
	Id                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Malware                    types.Bool   `tfsdk:"malware"`
	CommandAndControlCallbacks types.Bool   `tfsdk:"command_and_control_callbacks"`
	PhishingAttacks            types.Bool   `tfsdk:"phishing_attacks"`
	PotentiallyHarmfulDomains  types.Bool   `tfsdk:"potentially_harmful_domains"`
	NewlySeenDomains           types.Bool   `tfsdk:"newly_seen_domains"`
	DynamicDns                 types.Bool   `tfsdk:"dynamic_dns"`
	DnsTunnelingVpn            types.Bool   `tfsdk:"dns_tunneling_vpn"`
	Cryptomining               types.Bool   `tfsdk:"cryptomining"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ThreatCategory) getPath() string {
	return "/policies/v2/threatCategories"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ThreatCategory) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("malwareProtection"); value.Exists() {
		data.Malware = types.BoolValue(value.Bool())
	} else {
		data.Malware = types.BoolNull()
	}
	if value := res.Get("botnetProtection"); value.Exists() {
		data.CommandAndControlCallbacks = types.BoolValue(value.Bool())
	} else {
		data.CommandAndControlCallbacks = types.BoolNull()
	}
	if value := res.Get("phishingProtection"); value.Exists() {
		data.PhishingAttacks = types.BoolValue(value.Bool())
	} else {
		data.PhishingAttacks = types.BoolNull()
	}
	data.adjustRead(ctx, res, nil)
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody

// adjustRead reads the options that the API does not return as fields from `categoryBits`, a hexadecimal number of which
// bit N is `1 << N`. They are null if `categoryBits` is missing or is not a hexadecimal string. An empty string counts as 0, like "0":
// no option is set.
func (data *ThreatCategory) adjustRead(ctx context.Context, res gjson.Result, prior *ThreatCategory) {
	categoryBits := res.Get("categoryBits")
	data.PotentiallyHarmfulDomains = threatCategoryOption(categoryBits, threatCategoryBitPotentiallyHarmfulDomains)
	data.NewlySeenDomains = threatCategoryOption(categoryBits, threatCategoryBitNewlySeenDomains)
	data.DynamicDns = threatCategoryOption(categoryBits, threatCategoryBitDynamicDns)
	data.DnsTunnelingVpn = threatCategoryOption(categoryBits, threatCategoryBitDnsTunnelingVpn)
	data.Cryptomining = threatCategoryOption(categoryBits, threatCategoryBitCryptomining)
}
