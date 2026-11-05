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
	"net/netip"
	"net/url"
	"regexp"
	"slices"
	"strings"

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

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data DoNotDecryptList) toBody(ctx context.Context, state DoNotDecryptList) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	if len(data.Destinations) > 0 {
		var destinationsBody strings.Builder
		destinationsBody.WriteString("[")
		for _, item := range data.Destinations {
			itemBody := ""
			if !item.Destination.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "destination", item.Destination.ValueString())
			}
			if !item.Comment.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "comment", item.Comment.ValueString())
			}
			if itemBody != "" {
				if destinationsBody.Len() > 1 {
					destinationsBody.WriteString(",")
				}
				destinationsBody.WriteString(itemBody)
			}
		}
		destinationsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "destinations", destinationsBody.String())
	}
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

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *DoNotDecryptList) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
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

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *DoNotDecryptList) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
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

// End of section. //template:end fromBodyPartial

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

// adjustBody makes sure that `destinations` is part of the request, the API requires it also if there are none.
func (data DoNotDecryptList) adjustBody(ctx context.Context, req string) string {
	if !gjson.Get(req, "destinations").Exists() {
		req, _ = sjson.SetRaw(req, "destinations", "[]")
	}
	return req
}

// doNotDecryptListMaxItemsPerRequest is the maximum number of destinations the API accepts in a single request,
// for create as well as for the add and remove lists of an update.
const doNotDecryptListMaxItemsPerRequest = 500

// doNotDecryptListListRemovalChunk is the number of destinations that are removed from a destination list (see
// doNotDecryptListDestinationListPath) with a single request, the same limit that applies to the destination lists
// themselves.
const doNotDecryptListListRemovalChunk = 499

// doNotDecryptListDestinationListPath is the path of the destinations of a destination list. This is the destination
// lists API (/policies/v2/destinationlists), not the Do Not Decrypt Lists API, which has no /destinations path.
// The id to pass is `exceptiondomainlistId` of a Do Not Decrypt List, the destination list that holds its domains, IP
// addresses and CIDR blocks. It is never the id of the Do Not Decrypt List itself.
func doNotDecryptListDestinationListPath(listId string) string {
	return "/policies/v2/destinationlists/" + url.QueryEscape(listId) + "/destinations"
}

// doNotDecryptListDestinationsBody builds a JSON array of destinations. Comments are left out if withComment is false,
// the API does not need them to remove a destination.
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

// doNotDecryptListFirstDuplicate returns the first destination that is listed more than once. The `string_patterns` reject
// uppercase letters, so an exact comparison is enough.
func doNotDecryptListFirstDuplicate(destinations []string) (string, bool) {
	seen := make(map[string]struct{}, len(destinations))
	for _, destination := range destinations {
		if _, ok := seen[destination]; ok {
			return destination, true
		}
		seen[destination] = struct{}{}
	}
	return "", false
}

// doNotDecryptListIsIP reports whether a destination is an IP address or a CIDR block, as opposed to a domain name.
func doNotDecryptListIsIP(destination string) bool {
	if _, err := netip.ParsePrefix(destination); err == nil {
		return true
	}
	_, err := netip.ParseAddr(destination)
	return err == nil
}

// doNotDecryptListDestinationChanges compares the destinations of the state and of the plan:
//   - upsert: new destinations, and destinations with a new comment (adding an existing destination updates its comment)
//   - gone: destinations that are not in the plan anymore
//   - cleared: destinations whose comment is removed. They are removed and added again, as adding a destination without a
//     comment leaves the existing comment in place
func doNotDecryptListDestinationChanges(state, plan []DoNotDecryptListDestinations) (upsert, gone, cleared []DoNotDecryptListDestinations) {
	stateComments := make(map[string]types.String, len(state))
	for _, item := range state {
		stateComments[item.Destination.ValueString()] = item.Comment
	}
	planned := make(map[string]struct{}, len(plan))
	for _, item := range plan {
		destination := item.Destination.ValueString()
		planned[destination] = struct{}{}
		old, exists := stateComments[destination]
		switch {
		case !exists:
			upsert = append(upsert, item)
		case item.Comment.Equal(old):
		case item.Comment.IsNull():
			cleared = append(cleared, item)
		default:
			upsert = append(upsert, item)
		}
	}
	for _, item := range state {
		if _, ok := planned[item.Destination.ValueString()]; !ok {
			gone = append(gone, item)
		}
	}
	return upsert, gone, cleared
}

// doNotDecryptListStep is one request of an update, together with the changes it makes. The changes are kept to be able
// to tell what is applied if a later request of the same update fails.
type doNotDecryptListStep struct {
	// body is the PUT payload. It is empty if viaList is set.
	body string
	// viaList marks a step that removes IP addresses and CIDR blocks through the destination list behind the object.
	// The update request of the API answered "destination not found" for a private IP address and for a CIDR block, so
	// all of them are removed by id from the list, which worked for both.
	viaList                             bool
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

// doNotDecryptListListRemovalSteps builds the steps that remove IP addresses and CIDR blocks through the destination
// list, in chunks of doNotDecryptListListRemovalChunk.
func doNotDecryptListListRemovalSteps(destinations []DoNotDecryptListDestinations) []doNotDecryptListStep {
	var steps []doNotDecryptListStep
	for i := 0; i < len(destinations); i += doNotDecryptListListRemovalChunk {
		end := min(i+doNotDecryptListListRemovalChunk, len(destinations))
		steps = append(steps, doNotDecryptListStep{viaList: true, removeDestinations: destinations[i:end]})
	}
	return steps
}

// doNotDecryptListNotFound matches the error of the API for a destination that cannot be removed because it does not exist.
var doNotDecryptListNotFound = regexp.MustCompile(`destination '([^']+)' not found`)

// doNotDecryptListMissingDestination returns the destination that an error message of the API reports as not found.
func doNotDecryptListMissingDestination(message string) string {
	if match := doNotDecryptListNotFound.FindStringSubmatch(message); match != nil {
		return match[1]
	}
	return ""
}

// removes reports whether the step removes the destination.
func (step doNotDecryptListStep) removes(destination string) bool {
	return slices.ContainsFunc(step.removeDestinations, func(item DoNotDecryptListDestinations) bool {
		return item.Destination.ValueString() == destination
	})
}

// bodyWithout returns the PUT body of a removal step without the given destinations.
func (step doNotDecryptListStep) bodyWithout(dropped map[string]struct{}) string {
	var remaining []DoNotDecryptListDestinations
	for _, item := range step.removeDestinations {
		if _, ok := dropped[item.Destination.ValueString()]; !ok {
			remaining = append(remaining, item)
		}
	}
	if len(remaining) == 0 {
		body, _ := sjson.Delete(step.body, "removeDestinations")
		return body
	}
	body, _ := sjson.SetRaw(step.body, "removeDestinations", doNotDecryptListDestinationsBody(remaining, false))
	return body
}

// doNotDecryptListSplitIPs splits destinations into domain names and IP addresses or CIDR blocks.
func doNotDecryptListSplitIPs(destinations []DoNotDecryptListDestinations) (domains, ips []DoNotDecryptListDestinations) {
	for _, item := range destinations {
		if doNotDecryptListIsIP(item.Destination.ValueString()) {
			ips = append(ips, item)
		} else {
			domains = append(domains, item)
		}
	}
	return domains, ips
}

// toUpdateSteps builds the requests for an update. The API does not take full lists on update, but lists of
// destinations, applications and categories to add and to remove. Those are calculated from the state and the plan.
// The requests are ordered so that a failure leaves as much as possible in place:
//  1. additions, which are idempotent (an existing destination only gets its comment updated)
//  2. removal of the destinations that are gone: domain names with the update request, IP addresses and CIDR blocks
//     through the destination list (see doNotDecryptListStep)
//  3. removal of the destinations whose comment is removed, in the same way
//  4. adding those destinations again, without a comment
//
// If nothing but the name changes, a single request with the name is returned.
func (data DoNotDecryptList) toUpdateSteps(ctx context.Context, state DoNotDecryptList) []doNotDecryptListStep {
	name := data.Name.ValueString()
	planApplications := helpers.Int64SetElements(ctx, data.Applications)
	stateApplications := helpers.Int64SetElements(ctx, state.Applications)
	planCategories := helpers.Int64SetElements(ctx, data.DecryptExceptionCategories)
	stateCategories := helpers.Int64SetElements(ctx, state.DecryptExceptionCategories)

	upsert, gone, cleared := doNotDecryptListDestinationChanges(state.Destinations, data.Destinations)
	goneDomains, goneIPs := doNotDecryptListSplitIPs(gone)
	clearedDomains, clearedIPs := doNotDecryptListSplitIPs(cleared)

	steps := doNotDecryptListChangeSteps(name, true, upsert,
		helpers.DifferenceInt64(planApplications, stateApplications),
		helpers.DifferenceInt64(planCategories, stateCategories))
	steps = append(steps, doNotDecryptListChangeSteps(name, false, goneDomains,
		helpers.DifferenceInt64(stateApplications, planApplications),
		helpers.DifferenceInt64(stateCategories, planCategories))...)
	steps = append(steps, doNotDecryptListListRemovalSteps(goneIPs)...)
	steps = append(steps, doNotDecryptListChangeSteps(name, false, clearedDomains, nil, nil)...)
	steps = append(steps, doNotDecryptListListRemovalSteps(clearedIPs)...)
	steps = append(steps, doNotDecryptListChangeSteps(name, true, cleared, nil, nil)...)

	// The name is part of every PUT. If the name changed, it is sent before anything is removed through the destination
	// list, so that a rejected name changes nothing.
	if len(steps) == 0 || (steps[0].viaList && !data.Name.Equal(state.Name)) {
		body, _ := sjson.Set("", "name", name)
		steps = slices.Insert(steps, 0, doNotDecryptListStep{body: body})
	}
	return steps
}

// doNotDecryptListApplyInt64 removes and adds values to an Int64 set, the result is null if nothing is left.
func doNotDecryptListApplyInt64(ctx context.Context, set types.Set, remove, add []int64) types.Set {
	values := helpers.Int64SetElements(ctx, set)
	var result []attr.Value
	for _, v := range values {
		if !slices.Contains(remove, v) {
			result = append(result, types.Int64Value(v))
		}
	}
	for _, v := range add {
		if !slices.Contains(values, v) {
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
	if !step.viaList {
		data.Name = name
	}
	touched := make(map[string]struct{}, len(step.removeDestinations)+len(step.addDestinations))
	for _, item := range step.removeDestinations {
		touched[item.Destination.ValueString()] = struct{}{}
	}
	for _, item := range step.addDestinations {
		touched[item.Destination.ValueString()] = struct{}{}
	}
	var destinations []DoNotDecryptListDestinations
	for _, item := range data.Destinations {
		if _, ok := touched[item.Destination.ValueString()]; !ok {
			destinations = append(destinations, item)
		}
	}
	data.Destinations = append(destinations, step.addDestinations...)
	data.Applications = doNotDecryptListApplyInt64(ctx, data.Applications, step.removeApplications, step.addApplications)
	data.DecryptExceptionCategories = doNotDecryptListApplyInt64(ctx, data.DecryptExceptionCategories, step.removeCategories, step.addCategories)
}
