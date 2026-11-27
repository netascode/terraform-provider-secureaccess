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

type SecurityProfilePolicySetting struct {
	Id                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	DecryptionEnabled                types.Bool   `tfsdk:"decryption_enabled"`
	EndUserNotification              types.Bool   `tfsdk:"end_user_notification"`
	SingleSignOn                     types.Bool   `tfsdk:"single_sign_on"`
	SafeSearch                       types.Bool   `tfsdk:"safe_search"`
	AiSupplyChainProhibitedSuppliers types.Bool   `tfsdk:"ai_supply_chain_prohibited_suppliers"`
	AiSupplyChainCodeExecution       types.Bool   `tfsdk:"ai_supply_chain_code_execution"`
	AiSupplyChainCopyleftLicense     types.Bool   `tfsdk:"ai_supply_chain_copyleft_license"`
	AiSemanticInspection             types.Bool   `tfsdk:"ai_semantic_inspection"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data SecurityProfilePolicySetting) getPath() string {
	return "/policies/v2/policysettings"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *SecurityProfilePolicySetting) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("intelligentProxyHttps"); value.Exists() {
		data.DecryptionEnabled = types.BoolValue(value.Bool())
	} else {
		data.DecryptionEnabled = types.BoolNull()
	}
	if value := res.Get("swgDisplayBlockPage"); value.Exists() {
		data.EndUserNotification = types.BoolValue(value.Bool())
	} else {
		data.EndUserNotification = types.BoolNull()
	}
	if value := res.Get("saml"); value.Exists() {
		data.SingleSignOn = types.BoolValue(value.Bool())
	} else {
		data.SingleSignOn = types.BoolNull()
	}
	if value := res.Get("safeSearch"); value.Exists() {
		data.SafeSearch = types.BoolValue(value.Bool())
	} else {
		data.SafeSearch = types.BoolNull()
	}
	if value := res.Get("mcpSemanticInspection"); value.Exists() {
		data.AiSemanticInspection = types.BoolValue(value.Bool())
	} else {
		data.AiSemanticInspection = types.BoolNull()
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

// adjustRead reads the three AI Supply Chain Blocking options from `aiSupplyChainCategories`, a mask string. They are null if the mask
// is missing or is neither `0` nor 54 hexadecimal digits.
func (data *SecurityProfilePolicySetting) adjustRead(ctx context.Context, res gjson.Result, prior *SecurityProfilePolicySetting) {
	mask := res.Get("aiSupplyChainCategories")
	data.AiSupplyChainProhibitedSuppliers = aiSupplyChainOption(mask, aiSupplyChainValueProhibitedSuppliers)
	data.AiSupplyChainCodeExecution = aiSupplyChainOption(mask, aiSupplyChainValueCodeExecution)
	data.AiSupplyChainCopyleftLicense = aiSupplyChainOption(mask, aiSupplyChainValueCopyleftLicense)
}
