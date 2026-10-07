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
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type InternalNetwork struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Prefix       types.String `tfsdk:"prefix"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	NetworkId    types.Int64  `tfsdk:"network_id"`
	SiteId       types.Int64  `tfsdk:"site_id"`
	TunnelId     types.Int64  `tfsdk:"tunnel_id"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data InternalNetwork) getPath() string {
	return "/deployments/v2/internalnetworks"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data InternalNetwork) toBody(ctx context.Context, state InternalNetwork) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	if !data.Prefix.IsNull() {
		body, _ = sjson.Set(body, "ipAddress", data.Prefix.ValueString())
	}
	if !data.PrefixLength.IsNull() {
		body, _ = sjson.Set(body, "prefixLength", data.PrefixLength.ValueInt64())
	}
	if !data.NetworkId.IsNull() {
		body, _ = sjson.Set(body, "networkId", data.NetworkId.ValueInt64())
	}
	if !data.SiteId.IsNull() {
		body, _ = sjson.Set(body, "siteId", data.SiteId.ValueInt64())
	}
	if !data.TunnelId.IsNull() {
		body, _ = sjson.Set(body, "tunnelId", data.TunnelId.ValueInt64())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *InternalNetwork) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("ipAddress"); value.Exists() {
		data.Prefix = types.StringValue(value.String())
	} else {
		data.Prefix = types.StringNull()
	}
	if value := res.Get("prefixLength"); value.Exists() {
		data.PrefixLength = types.Int64Value(value.Int())
	} else {
		data.PrefixLength = types.Int64Null()
	}
	if value := res.Get("networkId"); value.Exists() {
		data.NetworkId = types.Int64Value(value.Int())
	} else {
		data.NetworkId = types.Int64Null()
	}
	if value := res.Get("siteId"); value.Exists() {
		data.SiteId = types.Int64Value(value.Int())
	} else {
		data.SiteId = types.Int64Null()
	}
	if value := res.Get("tunnelId"); value.Exists() {
		data.TunnelId = types.Int64Value(value.Int())
	} else {
		data.TunnelId = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *InternalNetwork) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("ipAddress"); value.Exists() && !data.Prefix.IsNull() {
		data.Prefix = types.StringValue(value.String())
	} else {
		data.Prefix = types.StringNull()
	}
	if value := res.Get("prefixLength"); value.Exists() && !data.PrefixLength.IsNull() {
		data.PrefixLength = types.Int64Value(value.Int())
	} else {
		data.PrefixLength = types.Int64Null()
	}
	if value := res.Get("networkId"); value.Exists() && !data.NetworkId.IsNull() {
		data.NetworkId = types.Int64Value(value.Int())
	} else {
		data.NetworkId = types.Int64Null()
	}
	if value := res.Get("siteId"); value.Exists() && !data.SiteId.IsNull() {
		data.SiteId = types.Int64Value(value.Int())
	} else {
		data.SiteId = types.Int64Null()
	}
	if value := res.Get("tunnelId"); value.Exists() && !data.TunnelId.IsNull() {
		data.TunnelId = types.Int64Value(value.Int())
	} else {
		data.TunnelId = types.Int64Null()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *InternalNetwork) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
