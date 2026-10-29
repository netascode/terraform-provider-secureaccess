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
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type ServiceObjectGroup struct {
	Id                  types.String                            `tfsdk:"id"`
	Name                types.String                            `tfsdk:"name"`
	Description         types.String                            `tfsdk:"description"`
	ServiceObjectGroups []ServiceObjectGroupServiceObjectGroups `tfsdk:"service_object_groups"`
	ServiceObjects      []ServiceObjectGroupServiceObjects      `tfsdk:"service_objects"`
	Literals            []ServiceObjectGroupLiterals            `tfsdk:"literals"`
}

type ServiceObjectGroupServiceObjectGroups struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type ServiceObjectGroupServiceObjects struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type ServiceObjectGroupLiterals struct {
	Protocol types.String `tfsdk:"protocol"`
	Ports    types.List   `tfsdk:"ports"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ServiceObjectGroup) getPath() string {
	return "/policies/v2/objects/serviceObjectGroups"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ServiceObjectGroup) toBody(ctx context.Context, state ServiceObjectGroup) string {
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
	if len(data.ServiceObjectGroups) > 0 {
		var serviceObjectGroupsBody strings.Builder
		serviceObjectGroupsBody.WriteString("[")
		for _, item := range data.ServiceObjectGroups {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if itemBody != "" {
				if serviceObjectGroupsBody.Len() > 1 {
					serviceObjectGroupsBody.WriteString(",")
				}
				serviceObjectGroupsBody.WriteString(itemBody)
			}
		}
		serviceObjectGroupsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "groups", serviceObjectGroupsBody.String())
	}
	if len(data.ServiceObjects) > 0 {
		var serviceObjectsBody strings.Builder
		serviceObjectsBody.WriteString("[")
		for _, item := range data.ServiceObjects {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if itemBody != "" {
				if serviceObjectsBody.Len() > 1 {
					serviceObjectsBody.WriteString(",")
				}
				serviceObjectsBody.WriteString(itemBody)
			}
		}
		serviceObjectsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "objects", serviceObjectsBody.String())
	}
	if len(data.Literals) > 0 {
		var literalsBody strings.Builder
		literalsBody.WriteString("[")
		for _, item := range data.Literals {
			itemBody := ""
			if !item.Protocol.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "protocol", item.Protocol.ValueString())
			}
			if !item.Ports.IsNull() {
				var values []string
				item.Ports.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "ports", values)
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

func (data *ServiceObjectGroup) fromBody(ctx context.Context, res gjson.Result) {
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
		data.ServiceObjectGroups = make([]ServiceObjectGroupServiceObjectGroups, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ServiceObjectGroupServiceObjectGroups{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("name"); value.Exists() {
				data.Name = types.StringValue(value.String())
			} else {
				data.Name = types.StringNull()
			}
			(*parent).ServiceObjectGroups = append((*parent).ServiceObjectGroups, data)
			return true
		})
	}
	if value := res.Get("objects"); value.Exists() {
		data.ServiceObjects = make([]ServiceObjectGroupServiceObjects, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ServiceObjectGroupServiceObjects{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("name"); value.Exists() {
				data.Name = types.StringValue(value.String())
			} else {
				data.Name = types.StringNull()
			}
			(*parent).ServiceObjects = append((*parent).ServiceObjects, data)
			return true
		})
	}
	if value := res.Get("values"); value.Exists() {
		data.Literals = make([]ServiceObjectGroupLiterals, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ServiceObjectGroupLiterals{}
			if value := res.Get("protocol|@case:lower"); value.Exists() {
				data.Protocol = types.StringValue(value.String())
			} else {
				data.Protocol = types.StringNull()
			}
			if value := res.Get("ports"); value.Exists() {
				data.Ports = helpers.GetStringList(value.Array())
			} else {
				data.Ports = types.ListNull(types.StringType)
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
func (data *ServiceObjectGroup) fromBodyPartial(ctx context.Context, res gjson.Result) {
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
	serviceObjectGroupsArray := res.Get("groups")
	for i := 0; i < len(data.ServiceObjectGroups); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.ServiceObjectGroups[i].Id.ValueInt64(), 10)}

		parent := &data
		data := (*parent).ServiceObjectGroups[i]
		var res gjson.Result

		serviceObjectGroupsArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing ServiceObjectGroups[%d] = %+v",
				i,
				(*parent).ServiceObjectGroups[i],
			))
			(*parent).ServiceObjectGroups = slices.Delete((*parent).ServiceObjectGroups, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("name"); value.Exists() {
			data.Name = types.StringValue(value.String())
		} else {
			data.Name = types.StringNull()
		}
		(*parent).ServiceObjectGroups[i] = data
	}
	serviceObjectsArray := res.Get("objects")
	for i := 0; i < len(data.ServiceObjects); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.ServiceObjects[i].Id.ValueInt64(), 10)}

		parent := &data
		data := (*parent).ServiceObjects[i]
		var res gjson.Result

		serviceObjectsArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing ServiceObjects[%d] = %+v",
				i,
				(*parent).ServiceObjects[i],
			))
			(*parent).ServiceObjects = slices.Delete((*parent).ServiceObjects, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("name"); value.Exists() {
			data.Name = types.StringValue(value.String())
		} else {
			data.Name = types.StringNull()
		}
		(*parent).ServiceObjects[i] = data
	}
	literalsArray := res.Get("values")
	for i := 0; i < len(data.Literals); i++ {
		keys := [...]string{"protocol|@case:lower"}
		keyValues := [...]string{data.Literals[i].Protocol.ValueString()}

		parent := &data
		data := (*parent).Literals[i]
		var res gjson.Result

		literalsArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing Literals[%d] = %+v",
				i,
				(*parent).Literals[i],
			))
			(*parent).Literals = slices.Delete((*parent).Literals, i, i+1)
			i--

			continue
		}
		if value := res.Get("protocol|@case:lower"); value.Exists() && !data.Protocol.IsNull() {
			data.Protocol = types.StringValue(value.String())
		} else {
			data.Protocol = types.StringNull()
		}
		if value := res.Get("ports"); value.Exists() && !data.Ports.IsNull() {
			data.Ports = helpers.GetStringList(value.Array())
		} else {
			data.Ports = types.ListNull(types.StringType)
		}
		(*parent).Literals[i] = data
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ServiceObjectGroup) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	serviceObjectGroupsArray := res.Get("groups")
	for i := range data.ServiceObjectGroups {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.ServiceObjectGroups[i].Id.ValueInt64(), 10)}

		var r gjson.Result
		serviceObjectGroupsArray.ForEach(
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
					r = v
					return false
				}
				return true
			},
		)
		if v := data.ServiceObjectGroups[i]; v.Name.IsUnknown() {
			if value := r.Get("name"); value.Exists() {
				v.Name = types.StringValue(value.String())
			} else {
				v.Name = types.StringNull()
			}
			data.ServiceObjectGroups[i] = v
		}
	}
	serviceObjectsArray := res.Get("objects")
	for i := range data.ServiceObjects {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.ServiceObjects[i].Id.ValueInt64(), 10)}

		var r gjson.Result
		serviceObjectsArray.ForEach(
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
					r = v
					return false
				}
				return true
			},
		)
		if v := data.ServiceObjects[i]; v.Name.IsUnknown() {
			if value := r.Get("name"); value.Exists() {
				v.Name = types.StringValue(value.String())
			} else {
				v.Name = types.StringNull()
			}
			data.ServiceObjects[i] = v
		}
	}
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

func (data ServiceObjectGroup) adjustBody(_ context.Context, req string) string {

	req, _ = sjson.Delete(req, "objects")
	var objectIds []int64
	for _, object := range data.ServiceObjects {
		objectIds = append(objectIds, object.Id.ValueInt64())
	}
	req, _ = sjson.Set(req, "objectIds", objectIds)

	req, _ = sjson.Delete(req, "groups")
	var groupIds []int64
	for _, group := range data.ServiceObjectGroups {
		groupIds = append(groupIds, group.Id.ValueInt64())
	}
	req, _ = sjson.Set(req, "groupIds", groupIds)

	return req
}
