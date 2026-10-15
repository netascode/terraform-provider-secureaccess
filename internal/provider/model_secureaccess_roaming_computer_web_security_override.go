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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type RoamingComputerWebSecurityOverride struct {
	EnabledOriginIds  types.List `tfsdk:"enabled_origin_ids"`
	DisabledOriginIds types.List `tfsdk:"disabled_origin_ids"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data RoamingComputerWebSecurityOverride) getPath() string {
	return "/deployments/v2/deviceSettings/SWGEnabled"
}

// End of section. //template:end getPath

const (
	// The API accepts 1-100 origin IDs per request
	roamingComputerWebSecurityOverrideBatchSize = 100
	roamingComputerWebSecurityOverrideEnabled   = "1"
	roamingComputerWebSecurityOverrideDisabled  = "0"
)

// enabledIds returns the origin IDs with Secure Web Gateway always enabled.
func (data RoamingComputerWebSecurityOverride) enabledIds(ctx context.Context) []int64 {
	var ids []int64
	data.EnabledOriginIds.ElementsAs(ctx, &ids, false)
	return ids
}

// disabledIds returns the origin IDs with Secure Web Gateway always disabled.
func (data RoamingComputerWebSecurityOverride) disabledIds(ctx context.Context) []int64 {
	var ids []int64
	data.DisabledOriginIds.ElementsAs(ctx, &ids, false)
	return ids
}

// allIds returns all managed origin IDs, enabled first.
func (data RoamingComputerWebSecurityOverride) allIds(ctx context.Context) []int64 {
	return append(data.enabledIds(ctx), data.disabledIds(ctx)...)
}

// toBodyOriginIds builds the body of the /list and /remove requests.
func (data RoamingComputerWebSecurityOverride) toBodyOriginIds(ids []int64) string {
	body, _ := sjson.Set("", "originIds", ids)
	return body
}

// toBodySet builds the body of the /set request.
func (data RoamingComputerWebSecurityOverride) toBodySet(value string, ids []int64) string {
	body, _ := sjson.Set("", "value", value)
	body, _ = sjson.Set(body, "originIds", ids)
	return body
}

// fromBody updates the managed origin IDs based on the /list response items.
// Origin IDs keep their order, an origin ID whose value changed outside of Terraform is moved to the other list,
// and an origin ID without an override setting is removed.
func (data *RoamingComputerWebSecurityOverride) fromBody(ctx context.Context, items []gjson.Result) {
	values := make(map[int64]string, len(items))
	for _, v := range items {
		values[v.Get("originId").Int()] = v.Get("value").String()
	}

	var enabled, disabled []attr.Value
	for _, id := range data.allIds(ctx) {
		switch values[id] {
		case roamingComputerWebSecurityOverrideEnabled:
			enabled = append(enabled, types.Int64Value(id))
		case roamingComputerWebSecurityOverrideDisabled:
			disabled = append(disabled, types.Int64Value(id))
		}
	}

	data.EnabledOriginIds = roamingComputerWebSecurityOverrideList(data.EnabledOriginIds, enabled)
	data.DisabledOriginIds = roamingComputerWebSecurityOverrideList(data.DisabledOriginIds, disabled)
}

// roamingComputerWebSecurityOverrideList returns a list of the values, a null list is kept if there are no values.
func roamingComputerWebSecurityOverrideList(current types.List, values []attr.Value) types.List {
	if len(values) == 0 && current.IsNull() {
		return current
	}
	if values == nil {
		values = []attr.Value{}
	}
	return types.ListValueMust(types.Int64Type, values)
}

// roamingComputerWebSecurityOverrideBatches splits the origin IDs into chunks accepted by the API.
func roamingComputerWebSecurityOverrideBatches(ids []int64) [][]int64 {
	var batches [][]int64
	for len(ids) > roamingComputerWebSecurityOverrideBatchSize {
		batches = append(batches, ids[:roamingComputerWebSecurityOverrideBatchSize])
		ids = ids[roamingComputerWebSecurityOverrideBatchSize:]
	}
	if len(ids) > 0 {
		batches = append(batches, ids)
	}
	return batches
}

// roamingComputerWebSecurityOverrideDifference returns the origin IDs of a, which are not part of b.
func roamingComputerWebSecurityOverrideDifference(a, b []int64) []int64 {
	exclude := make(map[int64]struct{}, len(b))
	for _, id := range b {
		exclude[id] = struct{}{}
	}
	var result []int64
	for _, id := range a {
		if _, ok := exclude[id]; !ok {
			result = append(result, id)
		}
	}
	return result
}
