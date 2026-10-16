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
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tidwall/gjson"
)

type Identities struct {
	Filter         types.String `tfsdk:"filter"`
	IncludeDeleted types.Bool   `tfsdk:"include_deleted"`
	IdentityTypes  types.List   `tfsdk:"identity_types"`
	// Items are map[type]map[name]IdentitiesItems
	Items types.Map `tfsdk:"items"`
}

type IdentitiesItems struct {
	Id      types.String `tfsdk:"id"`
	Deleted types.Bool   `tfsdk:"deleted"`
}

var identitiesItemType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":      types.StringType,
	"deleted": types.BoolType,
}}

var identitiesItemsElemType = types.MapType{ElemType: identitiesItemType}

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data Identities) getPath() string {
	return "/reports/v2/identities"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

// End of section. //template:end toBody

func (data *Identities) fromBody(ctx context.Context, res gjson.Result) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.IncludeDeleted.IsNull() {
		data.IncludeDeleted = types.BoolValue(false)
	}
	includeDeleted := data.IncludeDeleted.ValueBool()

	// Types and names requested in config, names are matched case-insensitively
	var requested map[string]map[string]IdentitiesItems
	if !data.Items.IsNull() && !data.Items.IsUnknown() {
		diags.Append(data.Items.ElementsAs(ctx, &requested, false)...)
		if diags.HasError() {
			return diags
		}
	}
	requestedNames := make(map[string]map[string]bool, len(requested))
	for t, names := range requested {
		requestedNames[t] = make(map[string]bool, len(names))
		for name := range names {
			requestedNames[t][strings.ToLower(name)] = true
		}
	}

	items := make(map[string]map[string]IdentitiesItems)
	for t := range requested {
		items[t] = map[string]IdentitiesItems{}
	}
	res.ForEach(func(_, v gjson.Result) bool {
		t := v.Get("type.type").String()
		name := v.Get("label").String()
		if t == "" || name == "" {
			return true
		}
		// API has no parameter for deleted identities, they are skipped here.
		// Some deleted identities are only marked by `(deleted)` in the name.
		if !includeDeleted && (v.Get("deleted").Bool() || strings.Contains(name, "(deleted)")) {
			return true
		}
		if requested != nil {
			names, ok := requestedNames[t]
			if !ok || (len(names) > 0 && !names[strings.ToLower(name)]) {
				return true
			}
		}
		if items[t] == nil {
			items[t] = map[string]IdentitiesItems{}
		}
		if existing, ok := items[t][name]; ok {
			diags.AddWarning("Duplicate identity name",
				fmt.Sprintf("Multiple identities of type %q are named %q, keeping id %s and skipping id %s.", t, name, existing.Id.ValueString(), v.Get("id").String()))
			return true
		}
		items[t][name] = IdentitiesItems{
			Id:      types.StringValue(v.Get("id").String()),
			Deleted: types.BoolValue(v.Get("deleted").Bool()),
		}
		return true
	})
	for t, names := range requested {
		for name := range names {
			found := false
			for n := range items[t] {
				if strings.EqualFold(n, name) {
					found = true
					break
				}
			}
			if !found {
				tflog.Debug(ctx, fmt.Sprintf("subresource not found, removing: type=%v name=%v", t, name))
			}
		}
	}

	value, d := types.MapValueFrom(ctx, identitiesItemsElemType, items)
	diags.Append(d...)
	data.Items = value
	return diags
}

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
