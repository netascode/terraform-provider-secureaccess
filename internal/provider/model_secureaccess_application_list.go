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

type ApplicationList struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ApplicationIds         types.Set    `tfsdk:"application_ids"`
	ApplicationCategoryIds types.Set    `tfsdk:"application_category_ids"`
	IsDefault              types.Bool   `tfsdk:"is_default"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ApplicationList) getPath() string {
	return "/policies/v2/applicationLists"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ApplicationList) toBody(ctx context.Context, state ApplicationList) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "applicationListId", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "applicationListName", data.Name.ValueString())
	}
	if !data.ApplicationIds.IsNull() {
		var values []int64
		data.ApplicationIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "applicationIds", values)
	}
	if !data.ApplicationCategoryIds.IsNull() {
		var values []int64
		data.ApplicationCategoryIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "applicationCategoryIds", values)
	}
	if !data.IsDefault.IsNull() {
		body, _ = sjson.Set(body, "isDefault", data.IsDefault.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ApplicationList) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("applicationListName"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("applicationIds"); value.Exists() {
		data.ApplicationIds = helpers.GetInt64Set(value.Array())
	} else {
		data.ApplicationIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get("applicationCategoryIds"); value.Exists() {
		data.ApplicationCategoryIds = helpers.GetInt64Set(value.Array())
	} else {
		data.ApplicationCategoryIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get("isDefault"); value.Exists() {
		data.IsDefault = types.BoolValue(value.Bool())
	} else {
		data.IsDefault = types.BoolValue(false)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *ApplicationList) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("applicationListName"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("applicationIds"); value.Exists() && !data.ApplicationIds.IsNull() {
		data.ApplicationIds = helpers.GetInt64Set(value.Array())
	} else {
		data.ApplicationIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get("applicationCategoryIds"); value.Exists() && !data.ApplicationCategoryIds.IsNull() {
		data.ApplicationCategoryIds = helpers.GetInt64Set(value.Array())
	} else {
		data.ApplicationCategoryIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get("isDefault"); value.Exists() && !data.IsDefault.IsNull() {
		data.IsDefault = types.BoolValue(value.Bool())
	} else if data.IsDefault.ValueBool() != false {
		data.IsDefault = types.BoolNull()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ApplicationList) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

func (data ApplicationList) adjustBody(ctx context.Context, req string) string {
	// API requires `applicationIds` to be present, even if empty
	if !gjson.Get(req, "applicationIds").Exists() {
		req, _ = sjson.SetRaw(req, "applicationIds", "[]")
	}
	return req
}
