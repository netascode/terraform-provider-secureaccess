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
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type NetworkObjectGroup struct {
	Id                  types.String                            `tfsdk:"id"`
	Name                types.String                            `tfsdk:"name"`
	Description         types.String                            `tfsdk:"description"`
	NetworkObjectGroups []NetworkObjectGroupNetworkObjectGroups `tfsdk:"network_object_groups"`
	NetworkObjects      []NetworkObjectGroupNetworkObjects      `tfsdk:"network_objects"`
	Literals            []NetworkObjectGroupLiterals            `tfsdk:"literals"`
}

type NetworkObjectGroupNetworkObjectGroups struct {
	Id types.Int64 `tfsdk:"id"`
}

type NetworkObjectGroupNetworkObjects struct {
	Id types.Int64 `tfsdk:"id"`
}

type NetworkObjectGroupLiterals struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data NetworkObjectGroup) getPath() string {
	return "/policies/v2/objects/networkObjectGroups"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data NetworkObjectGroup) toBody(ctx context.Context, state NetworkObjectGroup) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	if !data.Description.IsNull() {
		body, _ = sjson.Set(body, "description", data.Description.ValueString())
	}
	if len(data.NetworkObjectGroups) > 0 {
		var networkObjectGroupsBody strings.Builder
		networkObjectGroupsBody.WriteString("[")
		for _, item := range data.NetworkObjectGroups {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if itemBody != "" {
				if networkObjectGroupsBody.Len() > 1 {
					networkObjectGroupsBody.WriteString(",")
				}
				networkObjectGroupsBody.WriteString(itemBody)
			}
		}
		networkObjectGroupsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "groups", networkObjectGroupsBody.String())
	}
	if len(data.NetworkObjects) > 0 {
		var networkObjectsBody strings.Builder
		networkObjectsBody.WriteString("[")
		for _, item := range data.NetworkObjects {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if itemBody != "" {
				if networkObjectsBody.Len() > 1 {
					networkObjectsBody.WriteString(",")
				}
				networkObjectsBody.WriteString(itemBody)
			}
		}
		networkObjectsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "objects", networkObjectsBody.String())
	}
	if len(data.Literals) > 0 {
		var literalsBody strings.Builder
		literalsBody.WriteString("[")
		for _, item := range data.Literals {
			itemBody := ""
			if !item.Type.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "type", item.Type.ValueString())
			}
			if !item.Value.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "addresses.0", item.Value.ValueString())
			}
			if itemBody != "" {
				if literalsBody.Len() > 1 {
					literalsBody.WriteString(",")
				}
				literalsBody.WriteString(itemBody)
			}
		}
		literalsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "values", literalsBody.String())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *NetworkObjectGroup) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("description"); value.Exists() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("groups"); value.Exists() {
		data.NetworkObjectGroups = make([]NetworkObjectGroupNetworkObjectGroups, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := NetworkObjectGroupNetworkObjectGroups{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			(*parent).NetworkObjectGroups = append((*parent).NetworkObjectGroups, data)
			return true
		})
	}
	if value := res.Get("objects"); value.Exists() {
		data.NetworkObjects = make([]NetworkObjectGroupNetworkObjects, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := NetworkObjectGroupNetworkObjects{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			(*parent).NetworkObjects = append((*parent).NetworkObjects, data)
			return true
		})
	}
	if value := res.Get("values"); value.Exists() {
		data.Literals = make([]NetworkObjectGroupLiterals, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := NetworkObjectGroupLiterals{}
			if value := res.Get("type"); value.Exists() {
				data.Type = types.StringValue(value.String())
			} else {
				data.Type = types.StringNull()
			}
			if value := res.Get("addresses.0"); value.Exists() {
				data.Value = types.StringValue(value.String())
			} else {
				data.Value = types.StringNull()
			}
			(*parent).Literals = append((*parent).Literals, data)
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
func (data *NetworkObjectGroup) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("description"); value.Exists() && !data.Description.IsNull() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	networkObjectGroupsArray := res.Get("groups")
	for i := 0; i < len(data.NetworkObjectGroups); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.NetworkObjectGroups[i].Id.ValueInt64(), 10)}

		parent := &data
		data := (*parent).NetworkObjectGroups[i]
		var res gjson.Result

		networkObjectGroupsArray.ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if v.Get(keys[ik]).String() != keyValues[ik] {
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
			tflog.Debug(ctx, fmt.Sprintf("removing NetworkObjectGroups[%d] = %+v",
				i,
				(*parent).NetworkObjectGroups[i],
			))
			(*parent).NetworkObjectGroups = slices.Delete((*parent).NetworkObjectGroups, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		(*parent).NetworkObjectGroups[i] = data
	}
	networkObjectsArray := res.Get("objects")
	for i := 0; i < len(data.NetworkObjects); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.NetworkObjects[i].Id.ValueInt64(), 10)}

		parent := &data
		data := (*parent).NetworkObjects[i]
		var res gjson.Result

		networkObjectsArray.ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if v.Get(keys[ik]).String() != keyValues[ik] {
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
			tflog.Debug(ctx, fmt.Sprintf("removing NetworkObjects[%d] = %+v",
				i,
				(*parent).NetworkObjects[i],
			))
			(*parent).NetworkObjects = slices.Delete((*parent).NetworkObjects, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		(*parent).NetworkObjects[i] = data
	}
	literalsArray := res.Get("values")
	for i := 0; i < len(data.Literals); i++ {
		keys := [...]string{"type", "addresses.0"}
		keyValues := [...]string{data.Literals[i].Type.ValueString(), data.Literals[i].Value.ValueString()}

		parent := &data
		data := (*parent).Literals[i]
		var res gjson.Result

		literalsArray.ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if v.Get(keys[ik]).String() != keyValues[ik] {
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
			tflog.Debug(ctx, fmt.Sprintf("removing Literals[%d] = %+v",
				i,
				(*parent).Literals[i],
			))
			(*parent).Literals = slices.Delete((*parent).Literals, i, i+1)
			i--

			continue
		}
		if value := res.Get("type"); value.Exists() && !data.Type.IsNull() {
			data.Type = types.StringValue(value.String())
		} else {
			data.Type = types.StringNull()
		}
		if value := res.Get("addresses.0"); value.Exists() && !data.Value.IsNull() {
			data.Value = types.StringValue(value.String())
		} else {
			data.Value = types.StringNull()
		}
		(*parent).Literals[i] = data
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *NetworkObjectGroup) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin Clone

// End of section. //template:end Clone

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyNonBulk

// End of section. //template:end toBodyNonBulk

// Section below is generated&owned by "gen/generator.go". //template:begin findObjectsToBeReplaced

// End of section. //template:end findObjectsToBeReplaced

// Section below is generated&owned by "gen/generator.go". //template:begin clearItemIds

// End of section. //template:end clearItemIds

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

func (data NetworkObjectGroup) adjustBody(_ context.Context, req string) string {

	req, _ = sjson.Delete(req, "objects")
	var objectIds []int64
	for _, object := range data.NetworkObjects {
		objectIds = append(objectIds, object.Id.ValueInt64())
	}
	req, _ = sjson.Set(req, "objectIds", objectIds)

	req, _ = sjson.Set(req, "groups", objectIds)
	var groupIds []int64
	for _, group := range data.NetworkObjectGroups {
		groupIds = append(groupIds, group.Id.ValueInt64())
	}
	req, _ = sjson.Set(req, "groupIds", groupIds)

	return req
}

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBodyBulk

// End of section. //template:end adjustBodyBulk
