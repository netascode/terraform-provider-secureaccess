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

type DLPCustomDataIdentifier struct {
	Id              types.String                             `tfsdk:"id"`
	Name            types.String                             `tfsdk:"name"`
	Description     types.String                             `tfsdk:"description"`
	UniqueThreshold types.Bool                               `tfsdk:"unique_threshold"`
	ThresholdRanges []DLPCustomDataIdentifierThresholdRanges `tfsdk:"threshold_ranges"`
	Terms           types.Set                                `tfsdk:"terms"`
	Patterns        types.Set                                `tfsdk:"patterns"`
	Proximity       types.Set                                `tfsdk:"proximity"`
}

type DLPCustomDataIdentifierThresholdRanges struct {
	Minimum  types.Int64  `tfsdk:"minimum"`
	Maximum  types.Int64  `tfsdk:"maximum"`
	Severity types.String `tfsdk:"severity"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data DLPCustomDataIdentifier) getPath() string {
	return "/policies/v2/dlp/custom/dataIdentifiers"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data DLPCustomDataIdentifier) toBody(ctx context.Context, state DLPCustomDataIdentifier) string {
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
	if !data.UniqueThreshold.IsNull() {
		body, _ = sjson.Set(body, "countUniqueMatchesOnly", data.UniqueThreshold.ValueBool())
	}
	if len(data.ThresholdRanges) > 0 {
		var thresholdRangesBody strings.Builder
		thresholdRangesBody.WriteString("[")
		for _, item := range data.ThresholdRanges {
			itemBody := ""
			if !item.Minimum.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "minimum", item.Minimum.ValueInt64())
			}
			if !item.Maximum.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "maximum", item.Maximum.ValueInt64())
			}
			if !item.Severity.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "severity", item.Severity.ValueString())
			}
			if itemBody != "" {
				if thresholdRangesBody.Len() > 1 {
					thresholdRangesBody.WriteString(",")
				}
				thresholdRangesBody.WriteString(itemBody)
			}
		}
		thresholdRangesBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "thresholdRanges", thresholdRangesBody.String())
	}
	if !data.Terms.IsNull() {
		var values []string
		data.Terms.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "terms", values)
	}
	if !data.Patterns.IsNull() {
		var values []string
		data.Patterns.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "patterns", values)
	}
	if !data.Proximity.IsNull() {
		var values []string
		data.Proximity.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "proximity", values)
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *DLPCustomDataIdentifier) fromBody(ctx context.Context, res gjson.Result) {
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
	if value := res.Get("countUniqueMatchesOnly"); value.Exists() {
		data.UniqueThreshold = types.BoolValue(value.Bool())
	} else {
		data.UniqueThreshold = types.BoolNull()
	}
	if value := res.Get("thresholdRanges"); value.Exists() {
		data.ThresholdRanges = make([]DLPCustomDataIdentifierThresholdRanges, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPCustomDataIdentifierThresholdRanges{}
			if value := res.Get("minimum"); value.Exists() {
				data.Minimum = types.Int64Value(value.Int())
			} else {
				data.Minimum = types.Int64Null()
			}
			if value := res.Get("maximum"); value.Exists() {
				data.Maximum = types.Int64Value(value.Int())
			} else {
				data.Maximum = types.Int64Null()
			}
			if value := res.Get("severity"); value.Exists() {
				data.Severity = types.StringValue(value.String())
			} else {
				data.Severity = types.StringNull()
			}
			(*parent).ThresholdRanges = append((*parent).ThresholdRanges, data)
			return true
		})
	}
	if value := res.Get("terms"); value.Exists() {
		data.Terms = helpers.GetStringSet(value.Array())
	} else {
		data.Terms = types.SetNull(types.StringType)
	}
	if value := res.Get("patterns"); value.Exists() {
		data.Patterns = helpers.GetStringSet(value.Array())
	} else {
		data.Patterns = types.SetNull(types.StringType)
	}
	if value := res.Get("proximity"); value.Exists() {
		data.Proximity = helpers.GetStringSet(value.Array())
	} else {
		data.Proximity = types.SetNull(types.StringType)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *DLPCustomDataIdentifier) fromBodyPartial(ctx context.Context, res gjson.Result) {
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
	if value := res.Get("countUniqueMatchesOnly"); value.Exists() && !data.UniqueThreshold.IsNull() {
		data.UniqueThreshold = types.BoolValue(value.Bool())
	} else {
		data.UniqueThreshold = types.BoolNull()
	}
	thresholdRangesArray := res.Get("thresholdRanges")
	for i := 0; i < len(data.ThresholdRanges); i++ {
		keys := [...]string{"minimum", "maximum", "severity"}
		keyValues := [...]string{strconv.FormatInt(data.ThresholdRanges[i].Minimum.ValueInt64(), 10), strconv.FormatInt(data.ThresholdRanges[i].Maximum.ValueInt64(), 10), data.ThresholdRanges[i].Severity.ValueString()}

		parent := &data
		data := (*parent).ThresholdRanges[i]
		var res gjson.Result

		thresholdRangesArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing ThresholdRanges[%d] = %+v",
				i,
				(*parent).ThresholdRanges[i],
			))
			(*parent).ThresholdRanges = slices.Delete((*parent).ThresholdRanges, i, i+1)
			i--

			continue
		}
		if value := res.Get("minimum"); value.Exists() && !data.Minimum.IsNull() {
			data.Minimum = types.Int64Value(value.Int())
		} else {
			data.Minimum = types.Int64Null()
		}
		if value := res.Get("maximum"); value.Exists() && !data.Maximum.IsNull() {
			data.Maximum = types.Int64Value(value.Int())
		} else {
			data.Maximum = types.Int64Null()
		}
		if value := res.Get("severity"); value.Exists() && !data.Severity.IsNull() {
			data.Severity = types.StringValue(value.String())
		} else {
			data.Severity = types.StringNull()
		}
		(*parent).ThresholdRanges[i] = data
	}
	if value := res.Get("terms"); value.Exists() && !data.Terms.IsNull() {
		data.Terms = helpers.GetStringSet(value.Array())
	} else {
		data.Terms = types.SetNull(types.StringType)
	}
	if value := res.Get("patterns"); value.Exists() && !data.Patterns.IsNull() {
		data.Patterns = helpers.GetStringSet(value.Array())
	} else {
		data.Patterns = types.SetNull(types.StringType)
	}
	if value := res.Get("proximity"); value.Exists() && !data.Proximity.IsNull() {
		data.Proximity = helpers.GetStringSet(value.Array())
	} else {
		data.Proximity = types.SetNull(types.StringType)
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *DLPCustomDataIdentifier) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
