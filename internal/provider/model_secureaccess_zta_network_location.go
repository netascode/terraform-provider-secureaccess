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

type ZTANetworkLocation struct {
	Id                      types.String                                `tfsdk:"id"`
	Name                    types.String                                `tfsdk:"name"`
	EgressCriteria          []ZTANetworkLocationEgressCriteria          `tfsdk:"egress_criteria"`
	EgressExclusionCriteria []ZTANetworkLocationEgressExclusionCriteria `tfsdk:"egress_exclusion_criteria"`
	PrivateSubnetCriteria   []ZTANetworkLocationPrivateSubnetCriteria   `tfsdk:"private_subnet_criteria"`
	Revision                types.Int64                                 `tfsdk:"revision"`
}

type ZTANetworkLocationEgressCriteria struct {
	SourceType           types.String `tfsdk:"source_type"`
	ManualAddress        types.String `tfsdk:"manual_address"`
	NetworkObjectId      types.String `tfsdk:"network_object_id"`
	NetworkObjectGroupId types.String `tfsdk:"network_object_group_id"`
}

type ZTANetworkLocationEgressExclusionCriteria struct {
	SourceType           types.String `tfsdk:"source_type"`
	ManualAddress        types.String `tfsdk:"manual_address"`
	NetworkObjectId      types.String `tfsdk:"network_object_id"`
	NetworkObjectGroupId types.String `tfsdk:"network_object_group_id"`
}

type ZTANetworkLocationPrivateSubnetCriteria struct {
	SourceType           types.String `tfsdk:"source_type"`
	ManualAddress        types.String `tfsdk:"manual_address"`
	NetworkObjectId      types.String `tfsdk:"network_object_id"`
	NetworkObjectGroupId types.String `tfsdk:"network_object_group_id"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ZTANetworkLocation) getPath() string {
	return "/deployments/v2/ztna/networkLocations"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ZTANetworkLocation) toBody(ctx context.Context, state ZTANetworkLocation) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "networkLocationId", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "networkLocationName", data.Name.ValueString())
	}
	if len(data.EgressCriteria) > 0 {
		var egressCriteriaBody strings.Builder
		egressCriteriaBody.WriteString("[")
		for _, item := range data.EgressCriteria {
			itemBody := ""
			if !item.SourceType.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "sourceType", item.SourceType.ValueString())
			}
			if !item.ManualAddress.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "address", item.ManualAddress.ValueString())
			}
			if !item.NetworkObjectId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "networkObjectId", item.NetworkObjectId.ValueString())
			}
			if !item.NetworkObjectGroupId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "networkObjectGroupId", item.NetworkObjectGroupId.ValueString())
			}
			if itemBody != "" {
				if egressCriteriaBody.Len() > 1 {
					egressCriteriaBody.WriteString(",")
				}
				egressCriteriaBody.WriteString(itemBody)
			}
		}
		egressCriteriaBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "egressCriteria", egressCriteriaBody.String())
	}
	if len(data.EgressExclusionCriteria) > 0 {
		var egressExclusionCriteriaBody strings.Builder
		egressExclusionCriteriaBody.WriteString("[")
		for _, item := range data.EgressExclusionCriteria {
			itemBody := ""
			if !item.SourceType.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "sourceType", item.SourceType.ValueString())
			}
			if !item.ManualAddress.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "address", item.ManualAddress.ValueString())
			}
			if !item.NetworkObjectId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "networkObjectId", item.NetworkObjectId.ValueString())
			}
			if !item.NetworkObjectGroupId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "networkObjectGroupId", item.NetworkObjectGroupId.ValueString())
			}
			if itemBody != "" {
				if egressExclusionCriteriaBody.Len() > 1 {
					egressExclusionCriteriaBody.WriteString(",")
				}
				egressExclusionCriteriaBody.WriteString(itemBody)
			}
		}
		egressExclusionCriteriaBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "egressExclusionCriteria", egressExclusionCriteriaBody.String())
	}
	if len(data.PrivateSubnetCriteria) > 0 {
		var privateSubnetCriteriaBody strings.Builder
		privateSubnetCriteriaBody.WriteString("[")
		for _, item := range data.PrivateSubnetCriteria {
			itemBody := ""
			if !item.SourceType.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "sourceType", item.SourceType.ValueString())
			}
			if !item.ManualAddress.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "address", item.ManualAddress.ValueString())
			}
			if !item.NetworkObjectId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "networkObjectId", item.NetworkObjectId.ValueString())
			}
			if !item.NetworkObjectGroupId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "networkObjectGroupId", item.NetworkObjectGroupId.ValueString())
			}
			if itemBody != "" {
				if privateSubnetCriteriaBody.Len() > 1 {
					privateSubnetCriteriaBody.WriteString(",")
				}
				privateSubnetCriteriaBody.WriteString(itemBody)
			}
		}
		privateSubnetCriteriaBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "privateSubnetCriteria", privateSubnetCriteriaBody.String())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ZTANetworkLocation) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("networkLocationName"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("egressCriteria"); value.Exists() {
		data.EgressCriteria = make([]ZTANetworkLocationEgressCriteria, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTANetworkLocationEgressCriteria{}
			if value := res.Get("sourceType"); value.Exists() {
				data.SourceType = types.StringValue(value.String())
			} else {
				data.SourceType = types.StringNull()
			}
			if value := res.Get("address"); value.Exists() {
				data.ManualAddress = types.StringValue(value.String())
			} else {
				data.ManualAddress = types.StringNull()
			}
			if value := res.Get("networkObjectId"); value.Exists() {
				data.NetworkObjectId = types.StringValue(value.String())
			} else {
				data.NetworkObjectId = types.StringNull()
			}
			if value := res.Get("networkObjectGroupId"); value.Exists() {
				data.NetworkObjectGroupId = types.StringValue(value.String())
			} else {
				data.NetworkObjectGroupId = types.StringNull()
			}
			(*parent).EgressCriteria = append((*parent).EgressCriteria, data)
			return true
		})
	}
	if value := res.Get("egressExclusionCriteria"); value.Exists() {
		data.EgressExclusionCriteria = make([]ZTANetworkLocationEgressExclusionCriteria, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTANetworkLocationEgressExclusionCriteria{}
			if value := res.Get("sourceType"); value.Exists() {
				data.SourceType = types.StringValue(value.String())
			} else {
				data.SourceType = types.StringNull()
			}
			if value := res.Get("address"); value.Exists() {
				data.ManualAddress = types.StringValue(value.String())
			} else {
				data.ManualAddress = types.StringNull()
			}
			if value := res.Get("networkObjectId"); value.Exists() {
				data.NetworkObjectId = types.StringValue(value.String())
			} else {
				data.NetworkObjectId = types.StringNull()
			}
			if value := res.Get("networkObjectGroupId"); value.Exists() {
				data.NetworkObjectGroupId = types.StringValue(value.String())
			} else {
				data.NetworkObjectGroupId = types.StringNull()
			}
			(*parent).EgressExclusionCriteria = append((*parent).EgressExclusionCriteria, data)
			return true
		})
	}
	if value := res.Get("privateSubnetCriteria"); value.Exists() {
		data.PrivateSubnetCriteria = make([]ZTANetworkLocationPrivateSubnetCriteria, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTANetworkLocationPrivateSubnetCriteria{}
			if value := res.Get("sourceType"); value.Exists() {
				data.SourceType = types.StringValue(value.String())
			} else {
				data.SourceType = types.StringNull()
			}
			if value := res.Get("address"); value.Exists() {
				data.ManualAddress = types.StringValue(value.String())
			} else {
				data.ManualAddress = types.StringNull()
			}
			if value := res.Get("networkObjectId"); value.Exists() {
				data.NetworkObjectId = types.StringValue(value.String())
			} else {
				data.NetworkObjectId = types.StringNull()
			}
			if value := res.Get("networkObjectGroupId"); value.Exists() {
				data.NetworkObjectGroupId = types.StringValue(value.String())
			} else {
				data.NetworkObjectGroupId = types.StringNull()
			}
			(*parent).PrivateSubnetCriteria = append((*parent).PrivateSubnetCriteria, data)
			return true
		})
	}
	if value := res.Get("rev"); value.Exists() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *ZTANetworkLocation) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("networkLocationName"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	egressCriteriaArray := res.Get("egressCriteria")
	for i := 0; i < len(data.EgressCriteria); i++ {
		keys := [...]string{"sourceType", "address", "networkObjectId", "networkObjectGroupId"}
		keyValues := [...]string{data.EgressCriteria[i].SourceType.ValueString(), data.EgressCriteria[i].ManualAddress.ValueString(), data.EgressCriteria[i].NetworkObjectId.ValueString(), data.EgressCriteria[i].NetworkObjectGroupId.ValueString()}

		parent := &data
		data := (*parent).EgressCriteria[i]
		var res gjson.Result

		egressCriteriaArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing EgressCriteria[%d] = %+v",
				i,
				(*parent).EgressCriteria[i],
			))
			(*parent).EgressCriteria = slices.Delete((*parent).EgressCriteria, i, i+1)
			i--

			continue
		}
		if value := res.Get("sourceType"); value.Exists() && !data.SourceType.IsNull() {
			data.SourceType = types.StringValue(value.String())
		} else {
			data.SourceType = types.StringNull()
		}
		if value := res.Get("address"); value.Exists() && !data.ManualAddress.IsNull() {
			data.ManualAddress = types.StringValue(value.String())
		} else {
			data.ManualAddress = types.StringNull()
		}
		if value := res.Get("networkObjectId"); value.Exists() && !data.NetworkObjectId.IsNull() {
			data.NetworkObjectId = types.StringValue(value.String())
		} else {
			data.NetworkObjectId = types.StringNull()
		}
		if value := res.Get("networkObjectGroupId"); value.Exists() && !data.NetworkObjectGroupId.IsNull() {
			data.NetworkObjectGroupId = types.StringValue(value.String())
		} else {
			data.NetworkObjectGroupId = types.StringNull()
		}
		(*parent).EgressCriteria[i] = data
	}
	egressExclusionCriteriaArray := res.Get("egressExclusionCriteria")
	for i := 0; i < len(data.EgressExclusionCriteria); i++ {
		keys := [...]string{"sourceType", "address", "networkObjectId", "networkObjectGroupId"}
		keyValues := [...]string{data.EgressExclusionCriteria[i].SourceType.ValueString(), data.EgressExclusionCriteria[i].ManualAddress.ValueString(), data.EgressExclusionCriteria[i].NetworkObjectId.ValueString(), data.EgressExclusionCriteria[i].NetworkObjectGroupId.ValueString()}

		parent := &data
		data := (*parent).EgressExclusionCriteria[i]
		var res gjson.Result

		egressExclusionCriteriaArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing EgressExclusionCriteria[%d] = %+v",
				i,
				(*parent).EgressExclusionCriteria[i],
			))
			(*parent).EgressExclusionCriteria = slices.Delete((*parent).EgressExclusionCriteria, i, i+1)
			i--

			continue
		}
		if value := res.Get("sourceType"); value.Exists() && !data.SourceType.IsNull() {
			data.SourceType = types.StringValue(value.String())
		} else {
			data.SourceType = types.StringNull()
		}
		if value := res.Get("address"); value.Exists() && !data.ManualAddress.IsNull() {
			data.ManualAddress = types.StringValue(value.String())
		} else {
			data.ManualAddress = types.StringNull()
		}
		if value := res.Get("networkObjectId"); value.Exists() && !data.NetworkObjectId.IsNull() {
			data.NetworkObjectId = types.StringValue(value.String())
		} else {
			data.NetworkObjectId = types.StringNull()
		}
		if value := res.Get("networkObjectGroupId"); value.Exists() && !data.NetworkObjectGroupId.IsNull() {
			data.NetworkObjectGroupId = types.StringValue(value.String())
		} else {
			data.NetworkObjectGroupId = types.StringNull()
		}
		(*parent).EgressExclusionCriteria[i] = data
	}
	privateSubnetCriteriaArray := res.Get("privateSubnetCriteria")
	for i := 0; i < len(data.PrivateSubnetCriteria); i++ {
		keys := [...]string{"sourceType", "address", "networkObjectId", "networkObjectGroupId"}
		keyValues := [...]string{data.PrivateSubnetCriteria[i].SourceType.ValueString(), data.PrivateSubnetCriteria[i].ManualAddress.ValueString(), data.PrivateSubnetCriteria[i].NetworkObjectId.ValueString(), data.PrivateSubnetCriteria[i].NetworkObjectGroupId.ValueString()}

		parent := &data
		data := (*parent).PrivateSubnetCriteria[i]
		var res gjson.Result

		privateSubnetCriteriaArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing PrivateSubnetCriteria[%d] = %+v",
				i,
				(*parent).PrivateSubnetCriteria[i],
			))
			(*parent).PrivateSubnetCriteria = slices.Delete((*parent).PrivateSubnetCriteria, i, i+1)
			i--

			continue
		}
		if value := res.Get("sourceType"); value.Exists() && !data.SourceType.IsNull() {
			data.SourceType = types.StringValue(value.String())
		} else {
			data.SourceType = types.StringNull()
		}
		if value := res.Get("address"); value.Exists() && !data.ManualAddress.IsNull() {
			data.ManualAddress = types.StringValue(value.String())
		} else {
			data.ManualAddress = types.StringNull()
		}
		if value := res.Get("networkObjectId"); value.Exists() && !data.NetworkObjectId.IsNull() {
			data.NetworkObjectId = types.StringValue(value.String())
		} else {
			data.NetworkObjectId = types.StringNull()
		}
		if value := res.Get("networkObjectGroupId"); value.Exists() && !data.NetworkObjectGroupId.IsNull() {
			data.NetworkObjectGroupId = types.StringValue(value.String())
		} else {
			data.NetworkObjectGroupId = types.StringNull()
		}
		(*parent).PrivateSubnetCriteria[i] = data
	}
	if value := res.Get("rev"); value.Exists() && !data.Revision.IsNull() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ZTANetworkLocation) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
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

func (data ZTANetworkLocation) adjustBody(ctx context.Context, req string) string {

	req, _ = sjson.Delete(req, "networkLocationId")
	// Revision is kept from state in the plan (unknown on create), the API requires the current one on update
	if !data.Revision.IsUnknown() && !data.Revision.IsNull() {
		req, _ = sjson.Set(req, "rev", data.Revision.ValueInt64())
	}

	return req
}

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
