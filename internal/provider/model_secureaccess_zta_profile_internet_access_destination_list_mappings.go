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

type ZTAProfileInternetAccessDestinationListMappings struct {
	Id               types.String                                                      `tfsdk:"id"`
	ZtaProfileId     types.String                                                      `tfsdk:"zta_profile_id"`
	Revision         types.Int64                                                       `tfsdk:"revision"`
	DestinationLists []ZTAProfileInternetAccessDestinationListMappingsDestinationLists `tfsdk:"destination_lists"`
}

type ZTAProfileInternetAccessDestinationListMappingsDestinationLists struct {
	Id         types.String `tfsdk:"id"`
	Exclusions types.Set    `tfsdk:"exclusions"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ZTAProfileInternetAccessDestinationListMappings) getPath() string {
	return fmt.Sprintf("/deployments/v2/ztna/profiles/%v/internetSteering/destinationListMappings", url.QueryEscape(data.ZtaProfileId.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ZTAProfileInternetAccessDestinationListMappings) toBody(ctx context.Context, state ZTAProfileInternetAccessDestinationListMappings) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if len(data.DestinationLists) > 0 {
		var destinationListsBody strings.Builder
		destinationListsBody.WriteString("[")
		for _, item := range data.DestinationLists {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueString())
			}
			if !item.Exclusions.IsNull() {
				var values []string
				item.Exclusions.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "exclusions", values)
			}
			if itemBody != "" {
				if destinationListsBody.Len() > 1 {
					destinationListsBody.WriteString(",")
				}
				destinationListsBody.WriteString(itemBody)
			}
		}
		destinationListsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "items", destinationListsBody.String())
	}
	return body
}

// End of section. //template:end toBody

// toBodyPatch builds the PATCH body with the add/remove operations which change the current mappings into the
// desired ones, the replace operations for the exclusions of destination lists mapped in both with different
// exclusions, and the revision the current mappings are based on. Mappings are given as destination list ID to
// exclusions. Returns an empty string if there is nothing to change.
func (data ZTAProfileInternetAccessDestinationListMappings) toBodyPatch(current, desired map[string][]string, rev int64) string {
	body := ""
	for _, id := range slices.Sorted(maps.Keys(desired)) {
		exclusions, mapped := current[id]
		if !mapped {
			value := map[string]any{"id": id}
			if len(desired[id]) > 0 {
				value["exclusions"] = desired[id]
			}
			body, _ = sjson.Set(body, "operations.-1", map[string]any{"op": "add", "path": "/items", "value": value})
		} else if !slices.Equal(slices.Sorted(slices.Values(exclusions)), slices.Sorted(slices.Values(desired[id]))) {
			// Exclusions are replaced as a whole, an empty list removes them
			value := map[string]any{"exclusions": append([]string{}, desired[id]...)}
			body, _ = sjson.Set(body, "operations.-1", map[string]any{"op": "replace", "path": "/items/" + id, "value": value})
		}
	}
	for _, id := range slices.Sorted(maps.Keys(current)) {
		if _, ok := desired[id]; !ok {
			body, _ = sjson.Set(body, "operations.-1", map[string]any{"op": "remove", "path": "/items", "value": map[string]any{"id": id}})
		}
	}
	if body == "" {
		return ""
	}
	body, _ = sjson.Set(body, "rev", rev)
	return body
}

// destinationLists returns the mapped destination lists, as destination list ID to exclusions.
func (data ZTAProfileInternetAccessDestinationListMappings) destinationLists(ctx context.Context) map[string][]string {
	mappings := make(map[string][]string, len(data.DestinationLists))
	for _, item := range data.DestinationLists {
		var exclusions []string
		item.Exclusions.ElementsAs(ctx, &exclusions, false)
		mappings[item.Id.ValueString()] = exclusions
	}
	return mappings
}

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ZTAProfileInternetAccessDestinationListMappings) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("rev"); value.Exists() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
	if value := res.Get("items"); value.Exists() {
		data.DestinationLists = make([]ZTAProfileInternetAccessDestinationListMappingsDestinationLists, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTAProfileInternetAccessDestinationListMappingsDestinationLists{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.StringValue(value.String())
			} else {
				data.Id = types.StringNull()
			}
			if value := res.Get("exclusions"); value.Exists() {
				data.Exclusions = helpers.GetStringSet(value.Array())
			} else {
				data.Exclusions = types.SetNull(types.StringType)
			}
			(*parent).DestinationLists = append((*parent).DestinationLists, data)
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
func (data *ZTAProfileInternetAccessDestinationListMappings) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("rev"); value.Exists() && !data.Revision.IsNull() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
	destinationListsArray := res.Get("items")
	for i := 0; i < len(data.DestinationLists); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{data.DestinationLists[i].Id.ValueString()}

		parent := &data
		data := (*parent).DestinationLists[i]
		var res gjson.Result

		destinationListsArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing DestinationLists[%d] = %+v",
				i,
				(*parent).DestinationLists[i],
			))
			(*parent).DestinationLists = slices.Delete((*parent).DestinationLists, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.StringValue(value.String())
		} else {
			data.Id = types.StringNull()
		}
		if value := res.Get("exclusions"); value.Exists() && !data.Exclusions.IsNull() {
			data.Exclusions = helpers.GetStringSet(value.Array())
		} else {
			data.Exclusions = types.SetNull(types.StringType)
		}
		(*parent).DestinationLists[i] = data
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ZTAProfileInternetAccessDestinationListMappings) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
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
