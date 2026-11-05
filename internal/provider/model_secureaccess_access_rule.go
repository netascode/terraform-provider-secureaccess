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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type AccessRule struct {
	Id                                 types.String                       `tfsdk:"id"`
	Name                               types.String                       `tfsdk:"name"`
	Description                        types.String                       `tfsdk:"description"`
	Action                             types.String                       `tfsdk:"action"`
	Priority                           types.Int64                        `tfsdk:"priority"`
	Enabled                            types.Bool                         `tfsdk:"enabled"`
	SourceIdentityIds                  types.Set                          `tfsdk:"source_identity_ids"`
	SourceIdentityTypeIds              types.Set                          `tfsdk:"source_identity_type_ids"`
	SourceGeolocations                 types.Set                          `tfsdk:"source_geolocations"`
	SourceNetworkObjectIds             types.Set                          `tfsdk:"source_network_object_ids"`
	SourceNetworkObjectGroupIds        types.Set                          `tfsdk:"source_network_object_group_ids"`
	SourcePrivateLiterals              types.Set                          `tfsdk:"source_private_literals"`
	SourceInternetLiterals             []AccessRuleSourceInternetLiterals `tfsdk:"source_internet_literals"`
	SourceOperator                     types.String                       `tfsdk:"source_operator"`
	DestinationNetworkObjectIds        types.Set                          `tfsdk:"destination_network_object_ids"`
	DestinationNetworkObjectGroupIds   types.Set                          `tfsdk:"destination_network_object_group_ids"`
	DestinationServiceObjectIds        types.Set                          `tfsdk:"destination_service_object_ids"`
	DestinationServiceObjectGroupIds   types.Set                          `tfsdk:"destination_service_object_group_ids"`
	DestinationPrivateResourceIds      types.Set                          `tfsdk:"destination_private_resource_ids"`
	DestinationPrivateResourceGroupIds types.Set                          `tfsdk:"destination_private_resource_group_ids"`
	DestinationSecurityGroupTagIds     types.Set                          `tfsdk:"destination_security_group_tag_ids"`
	DestinationContentCategoryIds      types.Set                          `tfsdk:"destination_content_category_ids"`
	DestinationContentCategoryListIds  types.Set                          `tfsdk:"destination_content_category_list_ids"`
	DestinationApplicationIds          types.Set                          `tfsdk:"destination_application_ids"`
	DestinationApplicationCategoryIds  types.Set                          `tfsdk:"destination_application_category_ids"`
	DestinationApplicationListIds      types.Set                          `tfsdk:"destination_application_list_ids"`
	DestinationListIds                 types.Set                          `tfsdk:"destination_list_ids"`
	DestinationGeolocations            types.Set                          `tfsdk:"destination_geolocations"`
	DestinationLiterals                []AccessRuleDestinationLiterals    `tfsdk:"destination_literals"`
	DestinationOperator                types.String                       `tfsdk:"destination_operator"`
	DestinationOperatorNetworkService  types.String                       `tfsdk:"destination_operator_network_service"`
	LogLevel                           types.String                       `tfsdk:"log_level"`
	TrafficType                        types.String                       `tfsdk:"traffic_type"`
	PublicSecurityProfileId            types.Int64                        `tfsdk:"public_security_profile_id"`
	PrivateSecurityProfileId           types.Int64                        `tfsdk:"private_security_profile_id"`
	TenantControlProfileId             types.Int64                        `tfsdk:"tenant_control_profile_id"`
	IpsProfileId                       types.Int64                        `tfsdk:"ips_profile_id"`
}

type AccessRuleSourceInternetLiterals struct {
	Ips   types.Set `tfsdk:"ips"`
	Ports types.Set `tfsdk:"ports"`
}

type AccessRuleDestinationLiterals struct {
	Ips      types.Set    `tfsdk:"ips"`
	Ports    types.Set    `tfsdk:"ports"`
	Protocol types.String `tfsdk:"protocol"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data AccessRule) getPath() string {
	return "/policies/v2/rules"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data AccessRule) toBody(ctx context.Context, state AccessRule) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "ruleId", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "ruleName", data.Name.ValueString())
	}
	if !data.Description.IsNull() {
		body, _ = sjson.Set(body, "ruleDescription", data.Description.ValueString())
	}
	if !data.Action.IsNull() {
		body, _ = sjson.Set(body, "ruleAction", data.Action.ValueString())
	}
	if !data.Priority.IsNull() {
		body, _ = sjson.Set(body, "rulePriority", data.Priority.ValueInt64())
	}
	if !data.Enabled.IsNull() {
		body, _ = sjson.Set(body, "ruleIsEnabled", data.Enabled.ValueBool())
	}
	if !data.SourceIdentityIds.IsNull() {
		var values []int64
		data.SourceIdentityIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.source.identity_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.SourceIdentityTypeIds.IsNull() {
		var values []int64
		data.SourceIdentityTypeIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.source.identity_type_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.SourceGeolocations.IsNull() {
		var values []string
		data.SourceGeolocations.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.source.geolocations_shared", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.SourceNetworkObjectIds.IsNull() {
		var values []int64
		data.SourceNetworkObjectIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.source.networkObjectIds", "attributeOperator": "IN", "attributeValue": values})
	}
	if !data.SourceNetworkObjectGroupIds.IsNull() {
		var values []int64
		data.SourceNetworkObjectGroupIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.source.networkObjectGroupIds", "attributeOperator": "IN", "attributeValue": values})
	}
	if !data.SourcePrivateLiterals.IsNull() {
		var values []string
		data.SourcePrivateLiterals.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.source.ip_address", "attributeOperator": "IN", "attributeValue": values})
	}
	if len(data.SourceInternetLiterals) > 0 {
		var sourceInternetLiteralsBody strings.Builder
		sourceInternetLiteralsBody.WriteString("[")
		for _, item := range data.SourceInternetLiterals {
			itemBody := ""
			if !item.Ips.IsNull() {
				var values []string
				item.Ips.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "ip", values)
			}
			if !item.Ports.IsNull() {
				var values []string
				item.Ports.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "port", values)
			}
			if itemBody != "" {
				if sourceInternetLiteralsBody.Len() > 1 {
					sourceInternetLiteralsBody.WriteString(",")
				}
				sourceInternetLiteralsBody.WriteString(itemBody)
			}
		}
		sourceInternetLiteralsBody.WriteString("]")
		entry, _ := sjson.SetRaw(`{"attributeName":"umbrella.source.composite_inline_ip","attributeOperator":"INTERSECT"}`, "attributeValue", sourceInternetLiteralsBody.String())
		body, _ = sjson.SetRaw(body, "ruleConditions.-1", entry)
	}
	if !data.DestinationNetworkObjectIds.IsNull() {
		var values []int64
		data.DestinationNetworkObjectIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.networkObjectIds", "attributeOperator": "IN", "attributeValue": values})
	}
	if !data.DestinationNetworkObjectGroupIds.IsNull() {
		var values []int64
		data.DestinationNetworkObjectGroupIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.networkObjectGroupIds", "attributeOperator": "IN", "attributeValue": values})
	}
	if !data.DestinationServiceObjectIds.IsNull() {
		var values []int64
		data.DestinationServiceObjectIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.serviceObjectIds", "attributeOperator": "IN", "attributeValue": values})
	}
	if !data.DestinationServiceObjectGroupIds.IsNull() {
		var values []int64
		data.DestinationServiceObjectGroupIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.serviceObjectGroupIds", "attributeOperator": "IN", "attributeValue": values})
	}
	if !data.DestinationPrivateResourceIds.IsNull() {
		var values []int64
		data.DestinationPrivateResourceIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.private_resource_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.DestinationPrivateResourceGroupIds.IsNull() {
		var values []int64
		data.DestinationPrivateResourceGroupIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.private_resource_group_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.DestinationSecurityGroupTagIds.IsNull() {
		var values []int64
		data.DestinationSecurityGroupTagIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.security_group_tag_ids", "attributeOperator": "=", "attributeValue": values})
	}
	if !data.DestinationContentCategoryIds.IsNull() {
		var values []int64
		data.DestinationContentCategoryIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.category_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.DestinationContentCategoryListIds.IsNull() {
		var values []int64
		data.DestinationContentCategoryListIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.category_list_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.DestinationApplicationIds.IsNull() {
		var values []int64
		data.DestinationApplicationIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.application_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.DestinationApplicationCategoryIds.IsNull() {
		var values []int64
		data.DestinationApplicationCategoryIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.application_category_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.DestinationApplicationListIds.IsNull() {
		var values []int64
		data.DestinationApplicationListIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.application_list_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.DestinationListIds.IsNull() {
		var values []int64
		data.DestinationListIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.destination_list_ids", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if !data.DestinationGeolocations.IsNull() {
		var values []string
		data.DestinationGeolocations.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": "umbrella.destination.geolocations", "attributeOperator": "INTERSECT", "attributeValue": values})
	}
	if len(data.DestinationLiterals) > 0 {
		var destinationLiteralsBody strings.Builder
		destinationLiteralsBody.WriteString("[")
		for _, item := range data.DestinationLiterals {
			itemBody := ""
			if !item.Ips.IsNull() {
				var values []string
				item.Ips.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "ip", values)
			}
			if !item.Ports.IsNull() {
				var values []string
				item.Ports.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "port", values)
			}
			if !item.Protocol.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "protocol", item.Protocol.ValueString())
			}
			if itemBody != "" {
				if destinationLiteralsBody.Len() > 1 {
					destinationLiteralsBody.WriteString(",")
				}
				destinationLiteralsBody.WriteString(itemBody)
			}
		}
		destinationLiteralsBody.WriteString("]")
		entry, _ := sjson.SetRaw(`{"attributeName":"umbrella.destination.composite_inline_ip","attributeOperator":"IN"}`, "attributeValue", destinationLiteralsBody.String())
		body, _ = sjson.SetRaw(body, "ruleConditions.-1", entry)
	}
	if !data.LogLevel.IsNull() {
		body, _ = sjson.Set(body, "ruleSettings.-1", map[string]any{"settingName": "umbrella.logLevel", "settingValue": data.LogLevel.ValueString()})
	}
	if !data.TrafficType.IsNull() {
		body, _ = sjson.Set(body, "ruleSettings.-1", map[string]any{"settingName": "umbrella.default.traffic", "settingValue": data.TrafficType.ValueString()})
	}
	if !data.PublicSecurityProfileId.IsNull() {
		body, _ = sjson.Set(body, "ruleSettings.-1", map[string]any{"settingName": "umbrella.posture.webProfileId", "settingValue": data.PublicSecurityProfileId.ValueInt64()})
	}
	if !data.PrivateSecurityProfileId.IsNull() {
		body, _ = sjson.Set(body, "ruleSettings.-1", map[string]any{"settingName": "umbrella.posture.privateSecurityProfileId", "settingValue": data.PrivateSecurityProfileId.ValueInt64()})
	}
	if !data.TenantControlProfileId.IsNull() {
		body, _ = sjson.Set(body, "ruleSettings.-1", map[string]any{"settingName": "sse.tenantControlProfileId", "settingValue": data.TenantControlProfileId.ValueInt64()})
	}
	if !data.IpsProfileId.IsNull() {
		body, _ = sjson.Set(body, "ruleSettings.-1", map[string]any{"settingName": "umbrella.posture.ipsProfileId", "settingValue": data.IpsProfileId.ValueInt64()})
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *AccessRule) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("ruleName"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("ruleDescription"); value.Exists() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("ruleAction"); value.Exists() {
		data.Action = types.StringValue(value.String())
	} else {
		data.Action = types.StringNull()
	}
	if value := res.Get("ruleIsEnabled"); value.Exists() {
		data.Enabled = types.BoolValue(value.Bool())
	} else {
		data.Enabled = types.BoolNull()
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.identity_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.SourceIdentityIds = helpers.GetInt64Set(value.Array())
	} else {
		data.SourceIdentityIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.identity_type_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.SourceIdentityTypeIds = helpers.GetInt64Set(value.Array())
	} else {
		data.SourceIdentityTypeIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.geolocations_shared").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.SourceGeolocations = helpers.GetStringSet(value.Array())
	} else {
		data.SourceGeolocations = types.SetNull(types.StringType)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.networkObjectIds").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.SourceNetworkObjectIds = helpers.GetInt64Set(value.Array())
	} else {
		data.SourceNetworkObjectIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.networkObjectGroupIds").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.SourceNetworkObjectGroupIds = helpers.GetInt64Set(value.Array())
	} else {
		data.SourceNetworkObjectGroupIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.ip_address").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.SourcePrivateLiterals = helpers.GetStringSet(value.Array())
	} else {
		data.SourcePrivateLiterals = types.SetNull(types.StringType)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.composite_inline_ip").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.SourceInternetLiterals = make([]AccessRuleSourceInternetLiterals, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := AccessRuleSourceInternetLiterals{}
			if value := res.Get("ip"); value.Exists() {
				data.Ips = helpers.GetStringSet(value.Array())
			} else {
				data.Ips = types.SetNull(types.StringType)
			}
			if value := res.Get("port"); value.Exists() {
				data.Ports = helpers.GetStringSet(value.Array())
			} else {
				data.Ports = types.SetNull(types.StringType)
			}
			(*parent).SourceInternetLiterals = append((*parent).SourceInternetLiterals, data)
			return true
		})
	}
	data.SourceOperator = types.StringValue("OR")
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.networkObjectIds").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationNetworkObjectIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationNetworkObjectIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.networkObjectGroupIds").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationNetworkObjectGroupIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationNetworkObjectGroupIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.serviceObjectIds").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationServiceObjectIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationServiceObjectIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.serviceObjectGroupIds").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationServiceObjectGroupIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationServiceObjectGroupIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.private_resource_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationPrivateResourceIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationPrivateResourceIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.private_resource_group_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationPrivateResourceGroupIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationPrivateResourceGroupIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.security_group_tag_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationSecurityGroupTagIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationSecurityGroupTagIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.category_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationContentCategoryIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationContentCategoryIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.category_list_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationContentCategoryListIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationContentCategoryListIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.application_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationApplicationIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationApplicationIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.application_category_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationApplicationCategoryIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationApplicationCategoryIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.application_list_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationApplicationListIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationApplicationListIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.destination_list_ids").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationListIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationListIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.geolocations").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationGeolocations = helpers.GetStringSet(value.Array())
	} else {
		data.DestinationGeolocations = types.SetNull(types.StringType)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.composite_inline_ip").attributeValue`); value.Exists() && value.Type != gjson.Null {
		data.DestinationLiterals = make([]AccessRuleDestinationLiterals, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := AccessRuleDestinationLiterals{}
			if value := res.Get("ip"); value.Exists() {
				data.Ips = helpers.GetStringSet(value.Array())
			} else {
				data.Ips = types.SetNull(types.StringType)
			}
			if value := res.Get("port"); value.Exists() {
				data.Ports = helpers.GetStringSet(value.Array())
			} else {
				data.Ports = types.SetNull(types.StringType)
			}
			if value := res.Get("protocol"); value.Exists() {
				data.Protocol = types.StringValue(value.String())
			} else {
				data.Protocol = types.StringNull()
			}
			(*parent).DestinationLiterals = append((*parent).DestinationLiterals, data)
			return true
		})
	}
	data.DestinationOperator = types.StringValue("OR")
	data.DestinationOperatorNetworkService = types.StringValue("AND")
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.logLevel").settingValue`); value.Exists() && value.Type != gjson.Null {
		data.LogLevel = types.StringValue(value.String())
	} else {
		data.LogLevel = types.StringNull()
	}
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.default.traffic").settingValue`); value.Exists() && value.Type != gjson.Null {
		data.TrafficType = types.StringValue(value.String())
	} else {
		data.TrafficType = types.StringNull()
	}
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.posture.webProfileId").settingValue`); value.Exists() && value.Type != gjson.Null {
		data.PublicSecurityProfileId = types.Int64Value(value.Int())
	} else {
		data.PublicSecurityProfileId = types.Int64Null()
	}
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.posture.privateSecurityProfileId").settingValue`); value.Exists() && value.Type != gjson.Null {
		data.PrivateSecurityProfileId = types.Int64Value(value.Int())
	} else {
		data.PrivateSecurityProfileId = types.Int64Null()
	}
	if value := res.Get(`ruleSettings.#(settingName=="sse.tenantControlProfileId").settingValue`); value.Exists() && value.Type != gjson.Null {
		data.TenantControlProfileId = types.Int64Value(value.Int())
	} else {
		data.TenantControlProfileId = types.Int64Null()
	}
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.posture.ipsProfileId").settingValue`); value.Exists() && value.Type != gjson.Null {
		data.IpsProfileId = types.Int64Value(value.Int())
	} else {
		data.IpsProfileId = types.Int64Null()
	}
	data.adjustRead(ctx, res, nil)
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *AccessRule) fromBodyPartial(ctx context.Context, res gjson.Result) {
	prior := *data
	if value := res.Get("ruleName"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("ruleDescription"); value.Exists() && !data.Description.IsNull() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("ruleAction"); value.Exists() && !data.Action.IsNull() {
		data.Action = types.StringValue(value.String())
	} else {
		data.Action = types.StringNull()
	}
	if value := res.Get("ruleIsEnabled"); value.Exists() && !data.Enabled.IsNull() {
		data.Enabled = types.BoolValue(value.Bool())
	} else {
		data.Enabled = types.BoolNull()
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.identity_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.SourceIdentityIds.IsNull() {
		data.SourceIdentityIds = helpers.GetInt64Set(value.Array())
	} else {
		data.SourceIdentityIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.identity_type_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.SourceIdentityTypeIds.IsNull() {
		data.SourceIdentityTypeIds = helpers.GetInt64Set(value.Array())
	} else {
		data.SourceIdentityTypeIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.geolocations_shared").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.SourceGeolocations.IsNull() {
		data.SourceGeolocations = helpers.GetStringSet(value.Array())
	} else {
		data.SourceGeolocations = types.SetNull(types.StringType)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.networkObjectIds").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.SourceNetworkObjectIds.IsNull() {
		data.SourceNetworkObjectIds = helpers.GetInt64Set(value.Array())
	} else {
		data.SourceNetworkObjectIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.networkObjectGroupIds").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.SourceNetworkObjectGroupIds.IsNull() {
		data.SourceNetworkObjectGroupIds = helpers.GetInt64Set(value.Array())
	} else {
		data.SourceNetworkObjectGroupIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.ip_address").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.SourcePrivateLiterals.IsNull() {
		data.SourcePrivateLiterals = helpers.GetStringSet(value.Array())
	} else {
		data.SourcePrivateLiterals = types.SetNull(types.StringType)
	}
	// The whole list is the value of a single list entry, it is read back as a whole
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.source.composite_inline_ip").attributeValue`); value.Exists() && value.Type != gjson.Null && len(data.SourceInternetLiterals) > 0 {
		data.SourceInternetLiterals = make([]AccessRuleSourceInternetLiterals, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := AccessRuleSourceInternetLiterals{}
			if value := res.Get("ip"); value.Exists() {
				data.Ips = helpers.GetStringSet(value.Array())
			} else {
				data.Ips = types.SetNull(types.StringType)
			}
			if value := res.Get("port"); value.Exists() {
				data.Ports = helpers.GetStringSet(value.Array())
			} else {
				data.Ports = types.SetNull(types.StringType)
			}
			(*parent).SourceInternetLiterals = append((*parent).SourceInternetLiterals, data)
			return true
		})
	} else {
		data.SourceInternetLiterals = nil
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.networkObjectIds").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationNetworkObjectIds.IsNull() {
		data.DestinationNetworkObjectIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationNetworkObjectIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.networkObjectGroupIds").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationNetworkObjectGroupIds.IsNull() {
		data.DestinationNetworkObjectGroupIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationNetworkObjectGroupIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.serviceObjectIds").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationServiceObjectIds.IsNull() {
		data.DestinationServiceObjectIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationServiceObjectIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.serviceObjectGroupIds").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationServiceObjectGroupIds.IsNull() {
		data.DestinationServiceObjectGroupIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationServiceObjectGroupIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.private_resource_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationPrivateResourceIds.IsNull() {
		data.DestinationPrivateResourceIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationPrivateResourceIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.private_resource_group_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationPrivateResourceGroupIds.IsNull() {
		data.DestinationPrivateResourceGroupIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationPrivateResourceGroupIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.security_group_tag_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationSecurityGroupTagIds.IsNull() {
		data.DestinationSecurityGroupTagIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationSecurityGroupTagIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.category_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationContentCategoryIds.IsNull() {
		data.DestinationContentCategoryIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationContentCategoryIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.category_list_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationContentCategoryListIds.IsNull() {
		data.DestinationContentCategoryListIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationContentCategoryListIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.application_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationApplicationIds.IsNull() {
		data.DestinationApplicationIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationApplicationIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.application_category_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationApplicationCategoryIds.IsNull() {
		data.DestinationApplicationCategoryIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationApplicationCategoryIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.application_list_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationApplicationListIds.IsNull() {
		data.DestinationApplicationListIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationApplicationListIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.destination_list_ids").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationListIds.IsNull() {
		data.DestinationListIds = helpers.GetInt64Set(value.Array())
	} else {
		data.DestinationListIds = types.SetNull(types.Int64Type)
	}
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.geolocations").attributeValue`); value.Exists() && value.Type != gjson.Null && !data.DestinationGeolocations.IsNull() {
		data.DestinationGeolocations = helpers.GetStringSet(value.Array())
	} else {
		data.DestinationGeolocations = types.SetNull(types.StringType)
	}
	// The whole list is the value of a single list entry, it is read back as a whole
	if value := res.Get(`ruleConditions.#(attributeName=="umbrella.destination.composite_inline_ip").attributeValue`); value.Exists() && value.Type != gjson.Null && len(data.DestinationLiterals) > 0 {
		data.DestinationLiterals = make([]AccessRuleDestinationLiterals, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := AccessRuleDestinationLiterals{}
			if value := res.Get("ip"); value.Exists() {
				data.Ips = helpers.GetStringSet(value.Array())
			} else {
				data.Ips = types.SetNull(types.StringType)
			}
			if value := res.Get("port"); value.Exists() {
				data.Ports = helpers.GetStringSet(value.Array())
			} else {
				data.Ports = types.SetNull(types.StringType)
			}
			if value := res.Get("protocol"); value.Exists() {
				data.Protocol = types.StringValue(value.String())
			} else {
				data.Protocol = types.StringNull()
			}
			(*parent).DestinationLiterals = append((*parent).DestinationLiterals, data)
			return true
		})
	} else {
		data.DestinationLiterals = nil
	}
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.logLevel").settingValue`); value.Exists() && value.Type != gjson.Null && !data.LogLevel.IsNull() {
		data.LogLevel = types.StringValue(value.String())
	} else {
		data.LogLevel = types.StringNull()
	}
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.default.traffic").settingValue`); value.Exists() && value.Type != gjson.Null && !data.TrafficType.IsNull() {
		data.TrafficType = types.StringValue(value.String())
	} else {
		data.TrafficType = types.StringNull()
	}
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.posture.webProfileId").settingValue`); value.Exists() && value.Type != gjson.Null && !data.PublicSecurityProfileId.IsNull() {
		data.PublicSecurityProfileId = types.Int64Value(value.Int())
	} else {
		data.PublicSecurityProfileId = types.Int64Null()
	}
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.posture.privateSecurityProfileId").settingValue`); value.Exists() && value.Type != gjson.Null && !data.PrivateSecurityProfileId.IsNull() {
		data.PrivateSecurityProfileId = types.Int64Value(value.Int())
	} else {
		data.PrivateSecurityProfileId = types.Int64Null()
	}
	if value := res.Get(`ruleSettings.#(settingName=="sse.tenantControlProfileId").settingValue`); value.Exists() && value.Type != gjson.Null && !data.TenantControlProfileId.IsNull() {
		data.TenantControlProfileId = types.Int64Value(value.Int())
	} else {
		data.TenantControlProfileId = types.Int64Null()
	}
	if value := res.Get(`ruleSettings.#(settingName=="umbrella.posture.ipsProfileId").settingValue`); value.Exists() && value.Type != gjson.Null && !data.IpsProfileId.IsNull() {
		data.IpsProfileId = types.Int64Value(value.Int())
	} else {
		data.IpsProfileId = types.Int64Null()
	}
	data.adjustRead(ctx, res, &prior)
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *AccessRule) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// adjustBody adds the `all` conditions and the logical operator expressions to the `ruleConditions`
func (data AccessRule) adjustBody(ctx context.Context, req string) string {
	return adjustAccessRuleConditions(req, data.SourceOperator.ValueString(), data.DestinationOperator.ValueString(), data.DestinationOperatorNetworkService.ValueString())
}

// adjustRead merges the conditions returned under an alternative name (e.g. `*_shared` written by the portal)
// and derives the logical operators from the `logical_operator` conditions. The expressions are not compared,
// as the API accepts many equivalent ones.
func (data *AccessRule) adjustRead(ctx context.Context, res gjson.Result, prior *AccessRule) {
	readAccessRuleAltConditions(ctx, data, prior, res)

	data.SourceOperator = types.StringValue(accessRuleOperatorOr)
	if accessRuleLogicalOperator(res, accessRuleSource) != "" {
		data.SourceOperator = types.StringValue(accessRuleOperatorAnd)
	}
	data.DestinationOperator = types.StringValue(accessRuleOperatorOr)
	data.DestinationOperatorNetworkService = types.StringValue(accessRuleOperatorAnd)
	if expression := accessRuleLogicalOperator(res, accessRuleDestination); expression != "" {
		data.DestinationOperator = types.StringValue(accessRuleOperatorAnd)
		data.DestinationOperatorNetworkService = types.StringValue(accessRuleNetworkServiceOperator(expression))
	}
}
