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
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type DoNotDecryptList struct {
	Id                         types.String                   `tfsdk:"id"`
	Name                       types.String                   `tfsdk:"name"`
	Destinations               []DoNotDecryptListDestinations `tfsdk:"destinations"`
	Applications               types.Set                      `tfsdk:"applications"`
	DecryptExceptionCategories types.Set                      `tfsdk:"decrypt_exception_categories"`
	IsDefault                  types.Bool                     `tfsdk:"is_default"`
	TotalDomains               types.Int64                    `tfsdk:"total_domains"`
	TotalIps                   types.Int64                    `tfsdk:"total_ips"`
	TotalApplications          types.Int64                    `tfsdk:"total_applications"`
}

type DoNotDecryptListDestinations struct {
	Destination types.String `tfsdk:"destination"`
	Comment     types.String `tfsdk:"comment"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data DoNotDecryptList) getPath() string {
	return "/policies/v2/doNotDecryptLists"
}

// End of section. //template:end getPath

// toBody builds the POST payload. The API requires `destinations` to be present, also if there are none.
func (data DoNotDecryptList) toBody(ctx context.Context, state DoNotDecryptList) string {
	body := ""
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	body, _ = sjson.SetRaw(body, "destinations", doNotDecryptListDestinationsBody(data.Destinations, true))
	if !data.Applications.IsNull() {
		var values []int64
		data.Applications.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "applications", values)
	}
	if !data.DecryptExceptionCategories.IsNull() {
		var values []int64
		data.DecryptExceptionCategories.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "decryptExceptionCategories", values)
	}
	return body
}

// fromBody reads a full response, used after `terraform import`. The API does not return destinations and
// applications (only their totals), so those stay null. Categories are returned as `{categoryId, name}` objects.
func (data *DoNotDecryptList) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("decryptExceptionCategories.#.categoryId").Array(); len(value) > 0 {
		data.DecryptExceptionCategories = helpers.GetInt64Set(value)
	} else {
		data.DecryptExceptionCategories = types.SetNull(types.Int64Type)
	}
	if value := res.Get("isDefault"); value.Exists() {
		data.IsDefault = types.BoolValue(value.Bool())
	} else {
		data.IsDefault = types.BoolNull()
	}
	if value := res.Get("totalDomains"); value.Exists() {
		data.TotalDomains = types.Int64Value(value.Int())
	} else {
		data.TotalDomains = types.Int64Null()
	}
	if value := res.Get("totalIps"); value.Exists() {
		data.TotalIps = types.Int64Value(value.Int())
	} else {
		data.TotalIps = types.Int64Null()
	}
	if value := res.Get("totalApplications"); value.Exists() {
		data.TotalApplications = types.Int64Value(value.Int())
	} else {
		data.TotalApplications = types.Int64Null()
	}
}

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls.
// Destinations and applications are never returned by the API, so the values from the state are kept as they are.
func (data *DoNotDecryptList) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if !data.DecryptExceptionCategories.IsNull() {
		data.DecryptExceptionCategories = helpers.GetInt64Set(res.Get("decryptExceptionCategories.#.categoryId").Array())
	}
	if value := res.Get("isDefault"); value.Exists() && !data.IsDefault.IsNull() {
		data.IsDefault = types.BoolValue(value.Bool())
	} else {
		data.IsDefault = types.BoolNull()
	}
	if value := res.Get("totalDomains"); value.Exists() {
		data.TotalDomains = types.Int64Value(value.Int())
	} else {
		data.TotalDomains = types.Int64Null()
	}
	if value := res.Get("totalIps"); value.Exists() {
		data.TotalIps = types.Int64Value(value.Int())
	} else {
		data.TotalIps = types.Int64Null()
	}
	if value := res.Get("totalApplications"); value.Exists() {
		data.TotalApplications = types.Int64Value(value.Int())
	} else {
		data.TotalApplications = types.Int64Null()
	}
}

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *DoNotDecryptList) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	if data.IsDefault.IsUnknown() {
		if value := res.Get("isDefault"); value.Exists() {
			data.IsDefault = types.BoolValue(value.Bool())
		} else {
			data.IsDefault = types.BoolNull()
		}
	}
	if value := res.Get("totalDomains"); value.Exists() {
		data.TotalDomains = types.Int64Value(value.Int())
	} else {
		data.TotalDomains = types.Int64Null()
	}
	if value := res.Get("totalIps"); value.Exists() {
		data.TotalIps = types.Int64Value(value.Int())
	} else {
		data.TotalIps = types.Int64Null()
	}
	if value := res.Get("totalApplications"); value.Exists() {
		data.TotalApplications = types.Int64Value(value.Int())
	} else {
		data.TotalApplications = types.Int64Null()
	}
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody

// doNotDecryptListMaxItemsPerRequest is the maximum number of destinations sent in a single request (create allows 500)
const doNotDecryptListMaxItemsPerRequest = 500

// doNotDecryptListDestinationsBody builds a JSON array of destinations. Comments are left out if withComment is false,
// the API only requires the `destination` value to remove a destination.
func doNotDecryptListDestinationsBody(items []DoNotDecryptListDestinations, withComment bool) string {
	body := "[]"
	for _, item := range items {
		itemBody := ""
		itemBody, _ = sjson.Set(itemBody, "destination", item.Destination.ValueString())
		if withComment && !item.Comment.IsNull() && !item.Comment.IsUnknown() {
			itemBody, _ = sjson.Set(itemBody, "comment", item.Comment.ValueString())
		}
		body, _ = sjson.SetRaw(body, "-1", itemBody)
	}
	return body
}

// doNotDecryptListDestinationKey identifies a destination including its comment, a destination with a changed
// comment is treated as removed and added again.
func doNotDecryptListDestinationKey(item DoNotDecryptListDestinations) string {
	return item.Destination.ValueString() + "\x00" + item.Comment.ValueString()
}

// doNotDecryptListDestinationsDiff returns the destinations that are in `a` but not in `b`.
func doNotDecryptListDestinationsDiff(a, b []DoNotDecryptListDestinations) []DoNotDecryptListDestinations {
	inB := make(map[string]struct{}, len(b))
	for _, item := range b {
		inB[doNotDecryptListDestinationKey(item)] = struct{}{}
	}
	var diff []DoNotDecryptListDestinations
	for _, item := range a {
		if _, ok := inB[doNotDecryptListDestinationKey(item)]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

// doNotDecryptListInt64Values returns the elements of an Int64 set, or nil if the set is null or unknown.
func doNotDecryptListInt64Values(ctx context.Context, set types.Set) []int64 {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	var values []int64
	set.ElementsAs(ctx, &values, false)
	return values
}

// doNotDecryptListInt64Diff returns the values that are in `a` but not in `b`.
func doNotDecryptListInt64Diff(a, b []int64) []int64 {
	var diff []int64
	for _, v := range a {
		if !slices.Contains(b, v) {
			diff = append(diff, v)
		}
	}
	return diff
}

// doNotDecryptListStep is one PUT request of an update, together with the changes it makes. The changes are kept to
// be able to tell what is applied if a later request of the same update fails.
type doNotDecryptListStep struct {
	body                                string
	removeDestinations, addDestinations []DoNotDecryptListDestinations
	removeApplications, addApplications []int64
	removeCategories, addCategories     []int64
}

// doNotDecryptListChangeSteps builds the PUT requests that add (`add` is true) or remove destinations, applications and
// categories. Destinations are split into chunks, applications and categories are sent with the first chunk.
// Returns nothing if there is nothing to change.
func doNotDecryptListChangeSteps(name string, add bool, destinations []DoNotDecryptListDestinations, applications, categories []int64) []doNotDecryptListStep {
	if len(destinations) == 0 && len(applications) == 0 && len(categories) == 0 {
		return nil
	}
	verb := "remove"
	if add {
		verb = "add"
	}
	var steps []doNotDecryptListStep
	for first, i := true, 0; first || i < len(destinations); first, i = false, i+doNotDecryptListMaxItemsPerRequest {
		var step doNotDecryptListStep
		step.body, _ = sjson.Set("", "name", name)
		if first {
			if len(applications) > 0 {
				step.body, _ = sjson.Set(step.body, verb+"Applications", applications)
			}
			if len(categories) > 0 {
				step.body, _ = sjson.Set(step.body, verb+"DecryptExceptionCategories", categories)
			}
		}
		var chunk []DoNotDecryptListDestinations
		if end := min(i+doNotDecryptListMaxItemsPerRequest, len(destinations)); end > i {
			chunk = destinations[i:end]
			step.body, _ = sjson.SetRaw(step.body, verb+"Destinations", doNotDecryptListDestinationsBody(chunk, add))
		}
		if add {
			step.addDestinations = chunk
			if first {
				step.addApplications, step.addCategories = applications, categories
			}
		} else {
			step.removeDestinations = chunk
			if first {
				step.removeApplications, step.removeCategories = applications, categories
			}
		}
		steps = append(steps, step)
	}
	return steps
}

// toUpdateSteps builds the PUT requests for an update. The API does not take full lists on update, but lists of
// destinations, applications and categories to add and to remove. Those are calculated from the state and the plan.
// Removals are sent first, so that a destination with a changed comment can be removed and added again.
// If nothing but the name changes, a single request with the name is returned.
func (data DoNotDecryptList) toUpdateSteps(ctx context.Context, state DoNotDecryptList) []doNotDecryptListStep {
	name := data.Name.ValueString()
	planApplications := doNotDecryptListInt64Values(ctx, data.Applications)
	stateApplications := doNotDecryptListInt64Values(ctx, state.Applications)
	planCategories := doNotDecryptListInt64Values(ctx, data.DecryptExceptionCategories)
	stateCategories := doNotDecryptListInt64Values(ctx, state.DecryptExceptionCategories)

	steps := doNotDecryptListChangeSteps(name, false,
		doNotDecryptListDestinationsDiff(state.Destinations, data.Destinations),
		doNotDecryptListInt64Diff(stateApplications, planApplications),
		doNotDecryptListInt64Diff(stateCategories, planCategories))
	steps = append(steps, doNotDecryptListChangeSteps(name, true,
		doNotDecryptListDestinationsDiff(data.Destinations, state.Destinations),
		doNotDecryptListInt64Diff(planApplications, stateApplications),
		doNotDecryptListInt64Diff(planCategories, stateCategories))...)
	if len(steps) == 0 {
		body, _ := sjson.Set("", "name", name)
		steps = append(steps, doNotDecryptListStep{body: body})
	}
	return steps
}

// doNotDecryptListApplyInt64 removes and adds values to an Int64 set, the result is null if nothing is left.
func doNotDecryptListApplyInt64(ctx context.Context, set types.Set, remove, add []int64) types.Set {
	values := doNotDecryptListInt64Values(ctx, set)
	var result []attr.Value
	for _, v := range values {
		if !slices.Contains(remove, v) {
			result = append(result, types.Int64Value(v))
		}
	}
	for _, v := range add {
		if !slices.Contains(values, v) || slices.Contains(remove, v) {
			result = append(result, types.Int64Value(v))
		}
	}
	if len(result) == 0 {
		return types.SetNull(types.Int64Type)
	}
	return types.SetValueMust(types.Int64Type, result)
}

// applyStep applies the changes of a request that the API accepted to the model. It is used to keep track of what the
// API holds if a later request of the same update fails: the API does not return destinations and applications, so
// the state could not be corrected by a refresh.
func (data *DoNotDecryptList) applyStep(ctx context.Context, name types.String, step doNotDecryptListStep) {
	data.Name = name
	var destinations []DoNotDecryptListDestinations
	for _, item := range data.Destinations {
		if !slices.ContainsFunc(step.removeDestinations, func(r DoNotDecryptListDestinations) bool {
			return r.Destination.ValueString() == item.Destination.ValueString()
		}) {
			destinations = append(destinations, item)
		}
	}
	data.Destinations = append(destinations, step.addDestinations...)
	data.Applications = doNotDecryptListApplyInt64(ctx, data.Applications, step.removeApplications, step.addApplications)
	data.DecryptExceptionCategories = doNotDecryptListApplyInt64(ctx, data.DecryptExceptionCategories, step.removeCategories, step.addCategories)
}
