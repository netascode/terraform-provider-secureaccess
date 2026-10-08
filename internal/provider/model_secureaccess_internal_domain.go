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

type InternalDomain struct {
	Id                          types.String `tfsdk:"id"`
	Domain                      types.String `tfsdk:"domain"`
	Description                 types.String `tfsdk:"description"`
	IncludeAllMobileDevices     types.Bool   `tfsdk:"include_all_mobile_devices"`
	IncludeAllVirtualAppliances types.Bool   `tfsdk:"include_all_virtual_appliances"`
	SiteIds                     types.List   `tfsdk:"site_ids"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data InternalDomain) getPath() string {
	return "/deployments/v2/internaldomains"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data InternalDomain) toBody(ctx context.Context, state InternalDomain) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if !data.Domain.IsNull() {
		body, _ = sjson.Set(body, "domain", data.Domain.ValueString())
	}
	if !data.Description.IsNull() {
		body, _ = sjson.Set(body, "description", data.Description.ValueString())
	}
	if !data.IncludeAllMobileDevices.IsNull() {
		body, _ = sjson.Set(body, "includeAllMobileDevices", data.IncludeAllMobileDevices.ValueBool())
	}
	if !data.IncludeAllVirtualAppliances.IsNull() {
		body, _ = sjson.Set(body, "includeAllVAs", data.IncludeAllVirtualAppliances.ValueBool())
	}
	if !data.SiteIds.IsNull() {
		var values []int64
		data.SiteIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "siteIds", values)
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *InternalDomain) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("domain"); value.Exists() {
		data.Domain = types.StringValue(value.String())
	} else {
		data.Domain = types.StringNull()
	}
	if value := res.Get("description"); value.Exists() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("includeAllMobileDevices"); value.Exists() {
		data.IncludeAllMobileDevices = types.BoolValue(value.Bool())
	} else {
		data.IncludeAllMobileDevices = types.BoolNull()
	}
	if value := res.Get("includeAllVAs"); value.Exists() {
		data.IncludeAllVirtualAppliances = types.BoolValue(value.Bool())
	} else {
		data.IncludeAllVirtualAppliances = types.BoolNull()
	}
	if value := res.Get("siteIds"); value.Exists() {
		data.SiteIds = helpers.GetInt64List(value.Array())
	} else {
		data.SiteIds = types.ListNull(types.Int64Type)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *InternalDomain) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("domain"); value.Exists() && !data.Domain.IsNull() {
		data.Domain = types.StringValue(value.String())
	} else {
		data.Domain = types.StringNull()
	}
	if value := res.Get("description"); value.Exists() && !data.Description.IsNull() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("includeAllMobileDevices"); value.Exists() && !data.IncludeAllMobileDevices.IsNull() {
		data.IncludeAllMobileDevices = types.BoolValue(value.Bool())
	} else {
		data.IncludeAllMobileDevices = types.BoolNull()
	}
	if value := res.Get("includeAllVAs"); value.Exists() && !data.IncludeAllVirtualAppliances.IsNull() {
		data.IncludeAllVirtualAppliances = types.BoolValue(value.Bool())
	} else {
		data.IncludeAllVirtualAppliances = types.BoolNull()
	}
	if value := res.Get("siteIds"); value.Exists() && !data.SiteIds.IsNull() {
		data.SiteIds = helpers.GetInt64List(value.Array())
	} else {
		data.SiteIds = types.ListNull(types.Int64Type)
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *InternalDomain) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

func (data InternalDomain) adjustBody(ctx context.Context, req string) string {
	req, _ = sjson.Delete(req, "id")
	return req
}
