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
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type ZTAProfileInternetAccessExceptions struct {
	Id           types.String                                   `tfsdk:"id"`
	ZtaProfileId types.String                                   `tfsdk:"zta_profile_id"`
	Revision     types.Int64                                    `tfsdk:"revision"`
	Exceptions   []ZTAProfileInternetAccessExceptionsExceptions `tfsdk:"exceptions"`
}

type ZTAProfileInternetAccessExceptionsExceptions struct {
	Destination types.String `tfsdk:"destination"`
	Description types.String `tfsdk:"description"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ZTAProfileInternetAccessExceptions) getPath() string {
	return fmt.Sprintf("/deployments/v2/ztna/profiles/%v/internetSteering/exclusions", url.QueryEscape(data.ZtaProfileId.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ZTAProfileInternetAccessExceptions) toBody(ctx context.Context, state ZTAProfileInternetAccessExceptions) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if len(data.Exceptions) > 0 {
		var exceptionsBody strings.Builder
		exceptionsBody.WriteString("[")
		for _, item := range data.Exceptions {
			itemBody := ""
			if !item.Destination.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "exclusion", item.Destination.ValueString())
			}
			if !item.Description.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "description", item.Description.ValueString())
			}
			if itemBody != "" {
				if exceptionsBody.Len() > 1 {
					exceptionsBody.WriteString(",")
				}
				exceptionsBody.WriteString(itemBody)
			}
		}
		exceptionsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "items", exceptionsBody.String())
	}
	return body
}

// End of section. //template:end toBody

// toBodyPatch builds the PATCH body with the add/remove operations which change the current exceptions into the
// desired ones, the replace operations for the descriptions of existing exceptions with a different description, and
// the revision the current exceptions are based on. Current exceptions are given as destination to API item, as removal
// and replacement need the API assigned ID, desired ones as destination to description. Returns an empty string if
// there is nothing to change.
func (data ZTAProfileInternetAccessExceptions) toBodyPatch(current map[string]gjson.Result, desired map[string]string, rev int64) string {
	body := ""
	for _, destination := range slices.Sorted(maps.Keys(desired)) {
		item, exists := current[destination]
		if !exists {
			value := map[string]any{"exclusion": destination}
			if desired[destination] != "" {
				value["description"] = desired[destination]
			}
			body, _ = sjson.Set(body, "operations.-1", map[string]any{"op": "add", "path": "/items", "value": value})
		} else if item.Get("description").String() != desired[destination] {
			// Replace sets an empty description when not given
			value := map[string]any{"exclusion": destination, "description": desired[destination]}
			body, _ = sjson.Set(body, "operations.-1", map[string]any{"op": "replace", "path": "/items/" + item.Get("id").String(), "value": value})
		}
	}
	for _, destination := range slices.Sorted(maps.Keys(current)) {
		if _, ok := desired[destination]; !ok {
			body, _ = sjson.Set(body, "operations.-1", map[string]any{"op": "remove", "path": "/items", "value": map[string]any{"id": current[destination].Get("id").String()}})
		}
	}
	if body == "" {
		return ""
	}
	body, _ = sjson.Set(body, "rev", rev)
	return body
}

// exceptions returns the exceptions, as destination to description.
func (data ZTAProfileInternetAccessExceptions) exceptions() map[string]string {
	exceptions := make(map[string]string, len(data.Exceptions))
	for _, item := range data.Exceptions {
		exceptions[item.Destination.ValueString()] = item.Description.ValueString()
	}
	return exceptions
}

// currentExceptions returns the exceptions of an API response, as destination to API item.
func (data ZTAProfileInternetAccessExceptions) currentExceptions(res gjson.Result) map[string]gjson.Result {
	exceptions := make(map[string]gjson.Result)
	for _, item := range res.Get("items").Array() {
		exceptions[item.Get("exclusion").String()] = item
	}
	return exceptions
}

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ZTAProfileInternetAccessExceptions) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("rev"); value.Exists() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
	if value := res.Get("items"); value.Exists() {
		data.Exceptions = make([]ZTAProfileInternetAccessExceptionsExceptions, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTAProfileInternetAccessExceptionsExceptions{}
			if value := res.Get("exclusion"); value.Exists() {
				data.Destination = types.StringValue(value.String())
			} else {
				data.Destination = types.StringNull()
			}
			if value := res.Get("description"); value.Exists() {
				data.Description = types.StringValue(value.String())
			} else {
				data.Description = types.StringNull()
			}
			(*parent).Exceptions = append((*parent).Exceptions, data)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *ZTAProfileInternetAccessExceptions) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("rev"); value.Exists() && !data.Revision.IsNull() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
	exceptionsArray := res.Get("items")
	for i := 0; i < len(data.Exceptions); i++ {
		keys := [...]string{"exclusion"}
		keyValues := [...]string{data.Exceptions[i].Destination.ValueString()}

		parent := &data
		data := (*parent).Exceptions[i]
		var res gjson.Result

		exceptionsArray.ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if helpers.GetKeyValue(v.Get(keys[ik])) != keyValues[ik] {
						found = false
						break
					}
					found = true
				}
				if found {
					res = v
					return false
				}
				return true
			},
		)
		if !res.Exists() {
			tflog.Debug(ctx, fmt.Sprintf("removing Exceptions[%d] = %+v",
				i,
				(*parent).Exceptions[i],
			))
			(*parent).Exceptions = slices.Delete((*parent).Exceptions, i, i+1)
			i--

			continue
		}
		if value := res.Get("exclusion"); value.Exists() && !data.Destination.IsNull() {
			data.Destination = types.StringValue(value.String())
		} else {
			data.Destination = types.StringNull()
		}
		if value := res.Get("description"); value.Exists() && !data.Description.IsNull() {
			data.Description = types.StringValue(value.String())
		} else {
			data.Description = types.StringNull()
		}
		(*parent).Exceptions[i] = data
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ZTAProfileInternetAccessExceptions) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	if data.Revision.IsUnknown() {
		if value := res.Get("rev"); value.Exists() {
			data.Revision = types.Int64Value(value.Int())
		} else {
			data.Revision = types.Int64Null()
		}
	}
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
