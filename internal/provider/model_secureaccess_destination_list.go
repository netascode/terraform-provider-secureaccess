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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type DestinationList struct {
	Id           types.String                  `tfsdk:"id"`
	Name         types.String                  `tfsdk:"name"`
	Destinations []DestinationListDestinations `tfsdk:"destinations"`
}

type DestinationListDestinations struct {
	Id          types.Int64  `tfsdk:"id"`
	Destination types.String `tfsdk:"destination"`
	Comment     types.String `tfsdk:"comment"`
	Type        types.String `tfsdk:"type"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data DestinationList) getPath() string {
	return "/policies/v2/destinationlists"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data DestinationList) toBody(ctx context.Context, state DestinationList) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	body, _ = sjson.Set(body, "access", "none")
	body, _ = sjson.Set(body, "bundleTypeId", 2)
	body, _ = sjson.Set(body, "isGlobal", false)
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
			if !item.Type.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "type", item.Type.ValueString())
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
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *DestinationList) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("destinations"); value.Exists() {
		data.Destinations = make([]DestinationListDestinations, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DestinationListDestinations{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("destination"); value.Exists() {
				data.Destination = types.StringValue(value.String())
			} else {
				data.Destination = types.StringNull()
			}
			if value := res.Get("comment"); value.Exists() {
				data.Comment = types.StringValue(value.String())
			} else {
				data.Comment = types.StringNull()
			}
			if value := res.Get("type"); value.Exists() {
				data.Type = types.StringValue(value.String())
			} else {
				data.Type = types.StringNull()
			}
			(*parent).Destinations = append((*parent).Destinations, data)
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
func (data *DestinationList) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	destinationsArray := res.Get("destinations")
	for i := 0; i < len(data.Destinations); i++ {
		keys := [...]string{"destination"}
		keyValues := [...]string{data.Destinations[i].Destination.ValueString()}

		parent := &data
		data := (*parent).Destinations[i]
		var res gjson.Result

		destinationsArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing Destinations[%d] = %+v",
				i,
				(*parent).Destinations[i],
			))
			(*parent).Destinations = slices.Delete((*parent).Destinations, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("destination"); value.Exists() && !data.Destination.IsNull() {
			data.Destination = types.StringValue(value.String())
		} else {
			data.Destination = types.StringNull()
		}
		if value := res.Get("comment"); value.Exists() && !data.Comment.IsNull() {
			data.Comment = types.StringValue(value.String())
		} else {
			data.Comment = types.StringNull()
		}
		if value := res.Get("type"); value.Exists() && !data.Type.IsNull() {
			data.Type = types.StringValue(value.String())
		} else {
			data.Type = types.StringNull()
		}
		(*parent).Destinations[i] = data
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *DestinationList) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	destinationsArray := res.Get("destinations")
	for i := range data.Destinations {
		keys := [...]string{"destination"}
		keyValues := [...]string{data.Destinations[i].Destination.ValueString()}

		var r gjson.Result
		destinationsArray.ForEach(
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
					r = v
					return false
				}
				return true
			},
		)
		if v := data.Destinations[i]; v.Id.IsUnknown() {
			if value := r.Get("id"); value.Exists() {
				v.Id = types.Int64Value(value.Int())
			} else {
				v.Id = types.Int64Null()
			}
			data.Destinations[i] = v
		}
	}
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
