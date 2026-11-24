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

type ZTAProfilePrivateAccessSteeringPrivateResourceMappings struct {
	Id               types.String                                                             `tfsdk:"id"`
	ZtaProfileId     types.String                                                             `tfsdk:"zta_profile_id"`
	Revision         types.Int64                                                              `tfsdk:"revision"`
	PrivateResources []ZTAProfilePrivateAccessSteeringPrivateResourceMappingsPrivateResources `tfsdk:"private_resources"`
}

type ZTAProfilePrivateAccessSteeringPrivateResourceMappingsPrivateResources struct {
	Id types.String `tfsdk:"id"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ZTAProfilePrivateAccessSteeringPrivateResourceMappings) getPath() string {
	return fmt.Sprintf("/deployments/v2/ztna/profiles/%v/privateSteering/resourceMappings", url.QueryEscape(data.ZtaProfileId.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ZTAProfilePrivateAccessSteeringPrivateResourceMappings) toBody(ctx context.Context, state ZTAProfilePrivateAccessSteeringPrivateResourceMappings) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if len(data.PrivateResources) > 0 {
		var privateResourcesBody strings.Builder
		privateResourcesBody.WriteString("[")
		for _, item := range data.PrivateResources {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueString())
			}
			if itemBody != "" {
				if privateResourcesBody.Len() > 1 {
					privateResourcesBody.WriteString(",")
				}
				privateResourcesBody.WriteString(itemBody)
			}
		}
		privateResourcesBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "items", privateResourcesBody.String())
	}
	return body
}

// End of section. //template:end toBody

// toBodyPatch builds the PATCH body with the add/remove operations which change the current mappings into the
// desired ones, and the revision the current mappings are based on. Returns an empty string if there is nothing to change.
func (data ZTAProfilePrivateAccessSteeringPrivateResourceMappings) toBodyPatch(current, desired []string, rev int64) string {
	body := ""
	for _, id := range desired {
		if !slices.Contains(current, id) {
			body, _ = sjson.Set(body, "operations.-1", map[string]any{"op": "add", "path": "/items", "value": map[string]any{"id": id}})
		}
	}
	for _, id := range current {
		if !slices.Contains(desired, id) {
			body, _ = sjson.Set(body, "operations.-1", map[string]any{"op": "remove", "path": "/items", "value": map[string]any{"id": id}})
		}
	}
	if body == "" {
		return ""
	}
	body, _ = sjson.Set(body, "rev", rev)
	return body
}

// privateResourceIds returns the IDs of the mapped private resources.
func (data ZTAProfilePrivateAccessSteeringPrivateResourceMappings) privateResourceIds() []string {
	var ids []string
	for _, item := range data.PrivateResources {
		ids = append(ids, item.Id.ValueString())
	}
	return ids
}

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ZTAProfilePrivateAccessSteeringPrivateResourceMappings) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("rev"); value.Exists() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
	if value := res.Get("items"); value.Exists() {
		data.PrivateResources = make([]ZTAProfilePrivateAccessSteeringPrivateResourceMappingsPrivateResources, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTAProfilePrivateAccessSteeringPrivateResourceMappingsPrivateResources{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.StringValue(value.String())
			} else {
				data.Id = types.StringNull()
			}
			(*parent).PrivateResources = append((*parent).PrivateResources, data)
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
func (data *ZTAProfilePrivateAccessSteeringPrivateResourceMappings) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("rev"); value.Exists() && !data.Revision.IsNull() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
	privateResourcesArray := res.Get("items")
	for i := 0; i < len(data.PrivateResources); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{data.PrivateResources[i].Id.ValueString()}

		parent := &data
		data := (*parent).PrivateResources[i]
		var res gjson.Result

		privateResourcesArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing PrivateResources[%d] = %+v",
				i,
				(*parent).PrivateResources[i],
			))
			(*parent).PrivateResources = slices.Delete((*parent).PrivateResources, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.StringValue(value.String())
		} else {
			data.Id = types.StringNull()
		}
		(*parent).PrivateResources[i] = data
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ZTAProfilePrivateAccessSteeringPrivateResourceMappings) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
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
