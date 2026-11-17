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
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type DLPCustomClassification struct {
	Id                      types.String                                     `tfsdk:"id"`
	Name                    types.String                                     `tfsdk:"name"`
	Description             types.String                                     `tfsdk:"description"`
	EvaluationStrategy      types.String                                     `tfsdk:"evaluation_strategy"`
	IncludedDataIdentifiers []DLPCustomClassificationIncludedDataIdentifiers `tfsdk:"included_data_identifiers"`
	ExcludedDataIdentifiers []DLPCustomClassificationExcludedDataIdentifiers `tfsdk:"excluded_data_identifiers"`
}

type DLPCustomClassificationIncludedDataIdentifiers struct {
	Id types.String `tfsdk:"id"`
}

type DLPCustomClassificationExcludedDataIdentifiers struct {
	Id types.String `tfsdk:"id"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data DLPCustomClassification) getPath() string {
	return "/policies/v2/dlp/custom/classifications"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data DLPCustomClassification) toBody(ctx context.Context, state DLPCustomClassification) string {
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
	if !data.EvaluationStrategy.IsNull() {
		body, _ = sjson.Set(body, "classifierEvaluationStrategy", data.EvaluationStrategy.ValueString())
	}
	if len(data.IncludedDataIdentifiers) > 0 {
		var includedDataIdentifiersBody strings.Builder
		includedDataIdentifiersBody.WriteString("[")
		for _, item := range data.IncludedDataIdentifiers {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueString())
			}
			if itemBody != "" {
				if includedDataIdentifiersBody.Len() > 1 {
					includedDataIdentifiersBody.WriteString(",")
				}
				includedDataIdentifiersBody.WriteString(itemBody)
			}
		}
		includedDataIdentifiersBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "dataIdentifiers", includedDataIdentifiersBody.String())
	}
	if len(data.ExcludedDataIdentifiers) > 0 {
		var excludedDataIdentifiersBody strings.Builder
		excludedDataIdentifiersBody.WriteString("[")
		for _, item := range data.ExcludedDataIdentifiers {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueString())
			}
			if itemBody != "" {
				if excludedDataIdentifiersBody.Len() > 1 {
					excludedDataIdentifiersBody.WriteString(",")
				}
				excludedDataIdentifiersBody.WriteString(itemBody)
			}
		}
		excludedDataIdentifiersBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "excludedDataIdentifiers", excludedDataIdentifiersBody.String())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *DLPCustomClassification) fromBody(ctx context.Context, res gjson.Result) {
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
	if value := res.Get("classifierEvaluationStrategy"); value.Exists() {
		data.EvaluationStrategy = types.StringValue(value.String())
	} else {
		data.EvaluationStrategy = types.StringNull()
	}
	if value := res.Get("dataIdentifiers"); value.Exists() {
		data.IncludedDataIdentifiers = make([]DLPCustomClassificationIncludedDataIdentifiers, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPCustomClassificationIncludedDataIdentifiers{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.StringValue(value.String())
			} else {
				data.Id = types.StringNull()
			}
			(*parent).IncludedDataIdentifiers = append((*parent).IncludedDataIdentifiers, data)
			return true
		})
	}
	if value := res.Get("excludedDataIdentifiers"); value.Exists() {
		data.ExcludedDataIdentifiers = make([]DLPCustomClassificationExcludedDataIdentifiers, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPCustomClassificationExcludedDataIdentifiers{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.StringValue(value.String())
			} else {
				data.Id = types.StringNull()
			}
			(*parent).ExcludedDataIdentifiers = append((*parent).ExcludedDataIdentifiers, data)
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
func (data *DLPCustomClassification) fromBodyPartial(ctx context.Context, res gjson.Result) {
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
	if value := res.Get("classifierEvaluationStrategy"); value.Exists() && !data.EvaluationStrategy.IsNull() {
		data.EvaluationStrategy = types.StringValue(value.String())
	} else {
		data.EvaluationStrategy = types.StringNull()
	}
	includedDataIdentifiersArray := res.Get("dataIdentifiers")
	for i := 0; i < len(data.IncludedDataIdentifiers); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{data.IncludedDataIdentifiers[i].Id.ValueString()}

		parent := &data
		data := (*parent).IncludedDataIdentifiers[i]
		var res gjson.Result

		includedDataIdentifiersArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing IncludedDataIdentifiers[%d] = %+v",
				i,
				(*parent).IncludedDataIdentifiers[i],
			))
			(*parent).IncludedDataIdentifiers = slices.Delete((*parent).IncludedDataIdentifiers, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.StringValue(value.String())
		} else {
			data.Id = types.StringNull()
		}
		(*parent).IncludedDataIdentifiers[i] = data
	}
	excludedDataIdentifiersArray := res.Get("excludedDataIdentifiers")
	for i := 0; i < len(data.ExcludedDataIdentifiers); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{data.ExcludedDataIdentifiers[i].Id.ValueString()}

		parent := &data
		data := (*parent).ExcludedDataIdentifiers[i]
		var res gjson.Result

		excludedDataIdentifiersArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing ExcludedDataIdentifiers[%d] = %+v",
				i,
				(*parent).ExcludedDataIdentifiers[i],
			))
			(*parent).ExcludedDataIdentifiers = slices.Delete((*parent).ExcludedDataIdentifiers, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.StringValue(value.String())
		} else {
			data.Id = types.StringNull()
		}
		(*parent).ExcludedDataIdentifiers[i] = data
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *DLPCustomClassification) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

func (data DLPCustomClassification) adjustBody(ctx context.Context, req string) string {

	req, _ = sjson.Delete(req, "dataIdentifiers")
	var includedDataIdentifierIds []string
	for _, includedDataIdentifier := range data.IncludedDataIdentifiers {
		includedDataIdentifierIds = append(includedDataIdentifierIds, includedDataIdentifier.Id.ValueString())
	}
	req, _ = sjson.Set(req, "dataIdentifiers", includedDataIdentifierIds)

	req, _ = sjson.Delete(req, "excludedDataIdentifiers")
	var excludedDataIdentifierIds []string
	for _, excludedDataIdentifier := range data.ExcludedDataIdentifiers {
		excludedDataIdentifierIds = append(excludedDataIdentifierIds, excludedDataIdentifier.Id.ValueString())
	}
	req, _ = sjson.Set(req, "excludedDataIdentifiers", excludedDataIdentifierIds)

	return req
}

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
