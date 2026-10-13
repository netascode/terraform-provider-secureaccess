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
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type AccessPolicySettings struct {
	Id                                   types.String `tfsdk:"id"`
	ZeroTrustAccessAuthentication        types.Bool   `tfsdk:"zero_trust_access_authentication"`
	ZeroTrustAccessAuthenticationTimeout types.Int64  `tfsdk:"zero_trust_access_authentication_timeout"`
	EditingSessionDuration               types.Int64  `tfsdk:"editing_session_duration"`
	EditingSessionExtension              types.Int64  `tfsdk:"editing_session_extension"`
	GlobalIpsProfileId                   types.String `tfsdk:"global_ips_profile_id"`
	Microsoft365Compatibility            types.Bool   `tfsdk:"microsoft_365_compatibility"`
	DecryptionGlobal                     types.Bool   `tfsdk:"decryption_global"`
	DecryptionDisableForSources          types.List   `tfsdk:"decryption_disable_for_sources"`
	LogDecryptedInternetDestinations     types.Bool   `tfsdk:"log_decrypted_internet_destinations"`
	LogDecryptedPrivateResources         types.Bool   `tfsdk:"log_decrypted_private_resources"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data AccessPolicySettings) getPath() string {
	return "/policies/v2/settings"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data AccessPolicySettings) toBody(ctx context.Context, state AccessPolicySettings) string {
	body := "[]"
	if !data.ZeroTrustAccessAuthentication.IsNull() {
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "sse.ztaAuthnTimeoutEnabled", "settingValue": data.ZeroTrustAccessAuthentication.ValueBool()})
	}
	if !data.ZeroTrustAccessAuthenticationTimeout.IsNull() {
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "sse.ztaAuthnTimeoutMinutes", "settingValue": data.ZeroTrustAccessAuthenticationTimeout.ValueInt64()})
	}
	if !data.EditingSessionDuration.IsNull() {
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "sse.policyEditSessionTimeoutMinutes", "settingValue": data.EditingSessionDuration.ValueInt64()})
	}
	if !data.EditingSessionExtension.IsNull() {
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "sse.policyEditSessionExtensionMinutes", "settingValue": data.EditingSessionExtension.ValueInt64()})
	}
	if !data.GlobalIpsProfileId.IsNull() {
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "umbrella.posture.ipsProfileId", "settingValue": data.GlobalIpsProfileId.ValueString()})
	}
	if !data.Microsoft365Compatibility.IsNull() {
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "umbrella.m365Compatibility", "settingValue": data.Microsoft365Compatibility.ValueBool()})
	}
	if !data.DecryptionGlobal.IsNull() {
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "sse.globalIPSEnabled", "settingValue": data.DecryptionGlobal.ValueBool()})
	}
	if !data.DecryptionDisableForSources.IsNull() {
		var values []string
		data.DecryptionDisableForSources.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "global.setting.disableDecryptionSource", "settingValue": values})
	}
	if !data.LogDecryptedInternetDestinations.IsNull() {
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "sse.decryption.logInternet", "settingValue": data.LogDecryptedInternetDestinations.ValueBool()})
	}
	if !data.LogDecryptedPrivateResources.IsNull() {
		body, _ = sjson.Set(body, "-1", map[string]any{"settingName": "sse.decryption.logPrivate", "settingValue": data.LogDecryptedPrivateResources.ValueBool()})
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *AccessPolicySettings) fromBody(ctx context.Context, res gjson.Result) {
	settings := helpers.GetSettings(res)
	if value, ok := settings["sse.ztaAuthnTimeoutEnabled"]; ok {
		data.ZeroTrustAccessAuthentication = types.BoolValue(value.Bool())
	} else {
		data.ZeroTrustAccessAuthentication = types.BoolNull()
	}
	if value, ok := settings["sse.ztaAuthnTimeoutMinutes"]; ok {
		data.ZeroTrustAccessAuthenticationTimeout = types.Int64Value(value.Int())
	} else {
		data.ZeroTrustAccessAuthenticationTimeout = types.Int64Null()
	}
	if value, ok := settings["sse.policyEditSessionTimeoutMinutes"]; ok {
		data.EditingSessionDuration = types.Int64Value(value.Int())
	} else {
		data.EditingSessionDuration = types.Int64Null()
	}
	if value, ok := settings["sse.policyEditSessionExtensionMinutes"]; ok {
		data.EditingSessionExtension = types.Int64Value(value.Int())
	} else {
		data.EditingSessionExtension = types.Int64Null()
	}
	if value, ok := settings["umbrella.posture.ipsProfileId"]; ok {
		data.GlobalIpsProfileId = types.StringValue(value.String())
	} else {
		data.GlobalIpsProfileId = types.StringNull()
	}
	if value, ok := settings["umbrella.m365Compatibility"]; ok {
		data.Microsoft365Compatibility = types.BoolValue(value.Bool())
	} else {
		data.Microsoft365Compatibility = types.BoolNull()
	}
	if value, ok := settings["sse.globalIPSEnabled"]; ok {
		data.DecryptionGlobal = types.BoolValue(value.Bool())
	} else {
		data.DecryptionGlobal = types.BoolNull()
	}
	if value, ok := settings["global.setting.disableDecryptionSource"]; ok {
		data.DecryptionDisableForSources = helpers.GetStringList(value.Array())
	} else {
		data.DecryptionDisableForSources = types.ListNull(types.StringType)
	}
	if value, ok := settings["sse.decryption.logInternet"]; ok {
		data.LogDecryptedInternetDestinations = types.BoolValue(value.Bool())
	} else {
		data.LogDecryptedInternetDestinations = types.BoolNull()
	}
	if value, ok := settings["sse.decryption.logPrivate"]; ok {
		data.LogDecryptedPrivateResources = types.BoolValue(value.Bool())
	} else {
		data.LogDecryptedPrivateResources = types.BoolNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *AccessPolicySettings) fromBodyPartial(ctx context.Context, res gjson.Result) {
	settings := helpers.GetSettings(res)
	if value, ok := settings["sse.ztaAuthnTimeoutEnabled"]; ok && !data.ZeroTrustAccessAuthentication.IsNull() {
		data.ZeroTrustAccessAuthentication = types.BoolValue(value.Bool())
	} else {
		data.ZeroTrustAccessAuthentication = types.BoolNull()
	}
	if value, ok := settings["sse.ztaAuthnTimeoutMinutes"]; ok && !data.ZeroTrustAccessAuthenticationTimeout.IsNull() {
		data.ZeroTrustAccessAuthenticationTimeout = types.Int64Value(value.Int())
	} else {
		data.ZeroTrustAccessAuthenticationTimeout = types.Int64Null()
	}
	if value, ok := settings["sse.policyEditSessionTimeoutMinutes"]; ok && !data.EditingSessionDuration.IsNull() {
		data.EditingSessionDuration = types.Int64Value(value.Int())
	} else {
		data.EditingSessionDuration = types.Int64Null()
	}
	if value, ok := settings["sse.policyEditSessionExtensionMinutes"]; ok && !data.EditingSessionExtension.IsNull() {
		data.EditingSessionExtension = types.Int64Value(value.Int())
	} else {
		data.EditingSessionExtension = types.Int64Null()
	}
	if value, ok := settings["umbrella.posture.ipsProfileId"]; ok && !data.GlobalIpsProfileId.IsNull() {
		data.GlobalIpsProfileId = types.StringValue(value.String())
	} else {
		data.GlobalIpsProfileId = types.StringNull()
	}
	if value, ok := settings["umbrella.m365Compatibility"]; ok && !data.Microsoft365Compatibility.IsNull() {
		data.Microsoft365Compatibility = types.BoolValue(value.Bool())
	} else {
		data.Microsoft365Compatibility = types.BoolNull()
	}
	if value, ok := settings["sse.globalIPSEnabled"]; ok && !data.DecryptionGlobal.IsNull() {
		data.DecryptionGlobal = types.BoolValue(value.Bool())
	} else {
		data.DecryptionGlobal = types.BoolNull()
	}
	if value, ok := settings["global.setting.disableDecryptionSource"]; ok && !data.DecryptionDisableForSources.IsNull() {
		data.DecryptionDisableForSources = helpers.GetStringList(value.Array())
	} else {
		data.DecryptionDisableForSources = types.ListNull(types.StringType)
	}
	if value, ok := settings["sse.decryption.logInternet"]; ok && !data.LogDecryptedInternetDestinations.IsNull() {
		data.LogDecryptedInternetDestinations = types.BoolValue(value.Bool())
	} else {
		data.LogDecryptedInternetDestinations = types.BoolNull()
	}
	if value, ok := settings["sse.decryption.logPrivate"]; ok && !data.LogDecryptedPrivateResources.IsNull() {
		data.LogDecryptedPrivateResources = types.BoolValue(value.Bool())
	} else {
		data.LogDecryptedPrivateResources = types.BoolNull()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *AccessPolicySettings) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
