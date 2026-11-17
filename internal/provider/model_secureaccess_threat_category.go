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
	Id                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	BundleTypeId                types.Int64  `tfsdk:"bundle_type_id"`
	IsDefault                   types.Bool   `tfsdk:"is_default"`
	IsSwgDefault                types.Bool   `tfsdk:"is_swg_default"`
	AmpScan                     types.Bool   `tfsdk:"amp_scan"`
	BotnetProtection            types.Bool   `tfsdk:"botnet_protection"`
	FileInspection              types.Bool   `tfsdk:"file_inspection"`
	IpFiltering                 types.Bool   `tfsdk:"ip_filtering"`
	MalwareProtection           types.Bool   `tfsdk:"malware_protection"`
	MalwareUrlProxy             types.Bool   `tfsdk:"malware_url_proxy"`
	PhishingProtection          types.Bool   `tfsdk:"phishing_protection"`
	SuperBotnet                 types.Bool   `tfsdk:"super_botnet"`
	SuperMalware                types.Bool   `tfsdk:"super_malware"`
	SuspiciousResponseFiltering types.Bool   `tfsdk:"suspicious_response_filtering"`
	UrlProxyHttps               types.Bool   `tfsdk:"url_proxy_https"`
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
	if value := res.Get("bundleTypeId"); value.Exists() {
		data.BundleTypeId = types.Int64Value(value.Int())
	} else {
		data.BundleTypeId = types.Int64Null()
	}
	if value := res.Get("isDefault"); value.Exists() {
		data.IsDefault = types.BoolValue(value.Bool())
	} else {
		data.IsDefault = types.BoolNull()
	}
	if value := res.Get("isSwgDefault"); value.Exists() {
		data.IsSwgDefault = types.BoolValue(value.Bool())
	} else {
		data.IsSwgDefault = types.BoolNull()
	}
	if value := res.Get("ampScan"); value.Exists() {
		data.AmpScan = types.BoolValue(value.Bool())
	} else {
		data.AmpScan = types.BoolNull()
	}
	if value := res.Get("botnetProtection"); value.Exists() {
		data.BotnetProtection = types.BoolValue(value.Bool())
	} else {
		data.BotnetProtection = types.BoolNull()
	}
	if value := res.Get("fileInspection"); value.Exists() {
		data.FileInspection = types.BoolValue(value.Bool())
	} else {
		data.FileInspection = types.BoolNull()
	}
	if value := res.Get("ipFiltering"); value.Exists() {
		data.IpFiltering = types.BoolValue(value.Bool())
	} else {
		data.IpFiltering = types.BoolNull()
	}
	if value := res.Get("malwareProtection"); value.Exists() {
		data.MalwareProtection = types.BoolValue(value.Bool())
	} else {
		data.MalwareProtection = types.BoolNull()
	}
	if value := res.Get("malwareUrlProxy"); value.Exists() {
		data.MalwareUrlProxy = types.BoolValue(value.Bool())
	} else {
		data.MalwareUrlProxy = types.BoolNull()
	}
	if value := res.Get("phishingProtection"); value.Exists() {
		data.PhishingProtection = types.BoolValue(value.Bool())
	} else {
		data.PhishingProtection = types.BoolNull()
	}
	if value := res.Get("superBotnet"); value.Exists() {
		data.SuperBotnet = types.BoolValue(value.Bool())
	} else {
		data.SuperBotnet = types.BoolNull()
	}
	if value := res.Get("superMalware"); value.Exists() {
		data.SuperMalware = types.BoolValue(value.Bool())
	} else {
		data.SuperMalware = types.BoolNull()
	}
	if value := res.Get("suspiciousResponseFiltering"); value.Exists() {
		data.SuspiciousResponseFiltering = types.BoolValue(value.Bool())
	} else {
		data.SuspiciousResponseFiltering = types.BoolNull()
	}
	if value := res.Get("urlProxyHttps"); value.Exists() {
		data.UrlProxyHttps = types.BoolValue(value.Bool())
	} else {
		data.UrlProxyHttps = types.BoolNull()
	}
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

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
