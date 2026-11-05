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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type NetworkTunnelGroup struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	DeviceType                types.String `tfsdk:"device_type"`
	Region                    types.String `tfsdk:"region"`
	BackupRegion              types.String `tfsdk:"backup_region"`
	TunnelIdFormatEmail       types.String `tfsdk:"tunnel_id_format_email"`
	TunnelIdFormatIpAddresses types.List   `tfsdk:"tunnel_id_format_ip_addresses"`
	PresharedKey              types.String `tfsdk:"preshared_key"`
	RoutingType               types.String `tfsdk:"routing_type"`
	StaticNetworkCidrs        types.Set    `tfsdk:"static_network_cidrs"`
	BgpAsNumber               types.String `tfsdk:"bgp_as_number"`
	BgpMultihopCount          types.Int64  `tfsdk:"bgp_multihop_count"`
	BgpMultihopIpRanges       types.Set    `tfsdk:"bgp_multihop_ip_ranges"`
	BgpPeeringTypes           types.Set    `tfsdk:"bgp_peering_types"`
	BgpOverrideSummarization  types.Bool   `tfsdk:"bgp_override_summarization"`
	Ipv6Routing               types.Bool   `tfsdk:"ipv6_routing"`
	Hubs                      types.List   `tfsdk:"hubs"`
}

type NetworkTunnelGroupHubs struct {
	Id             types.Int64  `tfsdk:"id"`
	IsPrimary      types.Bool   `tfsdk:"is_primary"`
	TunnelId       types.String `tfsdk:"tunnel_id"`
	DatacenterName types.String `tfsdk:"datacenter_name"`
	DatacenterIp   types.String `tfsdk:"datacenter_ip"`
}

var NetworkTunnelGroupHubsAttrTypes = map[string]attr.Type{
	"id":              types.Int64Type,
	"is_primary":      types.BoolType,
	"tunnel_id":       types.StringType,
	"datacenter_name": types.StringType,
	"datacenter_ip":   types.StringType,
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data NetworkTunnelGroup) getPath() string {
	return "/deployments/v2/networktunnelgroups"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data NetworkTunnelGroup) toBody(ctx context.Context, state NetworkTunnelGroup) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	if !data.DeviceType.IsNull() {
		body, _ = sjson.Set(body, "deviceType", data.DeviceType.ValueString())
	}
	if !data.Region.IsNull() {
		body, _ = sjson.Set(body, "region", data.Region.ValueString())
	}
	if !data.BackupRegion.IsNull() {
		body, _ = sjson.Set(body, "backupRegion", data.BackupRegion.ValueString())
	}
	if !data.TunnelIdFormatEmail.IsNull() {
		body, _ = sjson.Set(body, "authIdPrefix", data.TunnelIdFormatEmail.ValueString())
	}
	if !data.TunnelIdFormatIpAddresses.IsNull() {
		var values []string
		data.TunnelIdFormatIpAddresses.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "authIdPrefix", values)
	}
	if !data.PresharedKey.IsNull() {
		body, _ = sjson.Set(body, "passphrase", data.PresharedKey.ValueString())
	}
	if !data.RoutingType.IsNull() {
		body, _ = sjson.Set(body, "routing.type", data.RoutingType.ValueString())
	}
	if !data.StaticNetworkCidrs.IsNull() {
		var values []string
		data.StaticNetworkCidrs.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "routing.data.networkCIDRs", values)
	}
	if !data.BgpAsNumber.IsNull() {
		body, _ = sjson.Set(body, "routing.data.asNumber", data.BgpAsNumber.ValueString())
	}
	if !data.BgpMultihopCount.IsNull() {
		body, _ = sjson.Set(body, "routing.data.bgpHopCount", data.BgpMultihopCount.ValueInt64())
	}
	if !data.BgpMultihopIpRanges.IsNull() {
		var values []string
		data.BgpMultihopIpRanges.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "routing.data.bgpNeighborCIDRs", values)
	}
	if !data.BgpPeeringTypes.IsNull() {
		var values []string
		data.BgpPeeringTypes.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "routing.data.bgpPeeringTypes", values)
	}
	if !data.BgpOverrideSummarization.IsNull() {
		body, _ = sjson.Set(body, "resourceAttributes.isOverrideBGPSummarization", data.BgpOverrideSummarization.ValueBool())
	}
	if !data.Ipv6Routing.IsNull() {
		body, _ = sjson.Set(body, "routing.data.enableIPv6", data.Ipv6Routing.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// networkTunnelGroupPatchPaths lists top level body keys that can be updated with PATCH
var networkTunnelGroupPatchPaths = []string{"name", "authIdPrefix", "passphrase", "region", "backupRegion", "routing", "resourceAttributes"}

// toBodyPatch builds JSON Patch "replace" operations for the paths which differ between data and state.
// Returns an empty string if there is nothing to update.
func (data NetworkTunnelGroup) toBodyPatch(ctx context.Context, state NetworkTunnelGroup) string {
	planBody := data.toBody(ctx, NetworkTunnelGroup{})
	stateBody := state.toBody(ctx, NetworkTunnelGroup{})
	ops := "[]"
	for _, p := range networkTunnelGroupPatchPaths {
		planValue := gjson.Get(planBody, p)
		stateValue := gjson.Get(stateBody, p)
		if !planValue.Exists() {
			if stateValue.Exists() {
				tflog.Warn(ctx, fmt.Sprintf("%s: %s removed from configuration, it cannot be unset with PATCH", data.Id.ValueString(), p))
			}
			continue
		}
		if stateValue.Exists() && planValue.Raw == stateValue.Raw {
			continue
		}
		op := ""
		op, _ = sjson.Set(op, "op", "replace")
		op, _ = sjson.Set(op, "path", "/"+p)
		op, _ = sjson.SetRaw(op, "value", planValue.Raw)
		ops, _ = sjson.SetRaw(ops, "-1", op)
	}
	if !gjson.Get(ops, "0").Exists() {
		return ""
	}
	return ops
}

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *NetworkTunnelGroup) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("deviceType"); value.Exists() {
		data.DeviceType = types.StringValue(value.String())
	} else {
		data.DeviceType = types.StringNull()
	}
	if value := res.Get("region"); value.Exists() {
		data.Region = types.StringValue(value.String())
	} else {
		data.Region = types.StringNull()
	}
	if value := res.Get("backupRegion"); value.Exists() {
		data.BackupRegion = types.StringValue(value.String())
	} else {
		data.BackupRegion = types.StringNull()
	}
	if value := res.Get("routing.type"); value.Exists() {
		data.RoutingType = types.StringValue(value.String())
	} else {
		data.RoutingType = types.StringNull()
	}
	if value := res.Get("routing.data.networkCIDRs"); value.Exists() {
		data.StaticNetworkCidrs = helpers.GetStringSet(value.Array())
	} else {
		data.StaticNetworkCidrs = types.SetNull(types.StringType)
	}
	if value := res.Get("routing.data.asNumber"); value.Exists() {
		data.BgpAsNumber = types.StringValue(value.String())
	} else {
		data.BgpAsNumber = types.StringNull()
	}
	if value := res.Get("routing.data.bgpHopCount"); value.Exists() {
		data.BgpMultihopCount = types.Int64Value(value.Int())
	} else {
		data.BgpMultihopCount = types.Int64Null()
	}
	if value := res.Get("routing.data.bgpNeighborCIDRs"); value.Exists() {
		data.BgpMultihopIpRanges = helpers.GetStringSet(value.Array())
	} else {
		data.BgpMultihopIpRanges = types.SetNull(types.StringType)
	}
	if value := res.Get("routing.data.bgpPeeringTypes"); value.Exists() {
		data.BgpPeeringTypes = helpers.GetStringSet(value.Array())
	} else {
		data.BgpPeeringTypes = types.SetNull(types.StringType)
	}
	if value := res.Get("resourceAttributes.isOverrideBGPSummarization"); value.Exists() {
		data.BgpOverrideSummarization = types.BoolValue(value.Bool())
	} else {
		data.BgpOverrideSummarization = types.BoolNull()
	}
	if value := res.Get("routing.data.enableIPv6"); value.Exists() {
		data.Ipv6Routing = types.BoolValue(value.Bool())
	} else {
		data.Ipv6Routing = types.BoolNull()
	}
	data.fromBodyHubs(ctx, res)
}

// fromBodyHubs replaces the computed only hubs list with the values from a JSON.
func (data *NetworkTunnelGroup) fromBodyHubs(ctx context.Context, res gjson.Result) {
	elemType := types.ObjectType{AttrTypes: NetworkTunnelGroupHubsAttrTypes}
	value := res.Get("hubs")
	if !value.Exists() {
		data.Hubs = types.ListNull(elemType)
		return
	}
	items := make([]NetworkTunnelGroupHubs, 0, int(value.Get("#").Int()))
	value.ForEach(func(k, res gjson.Result) bool {
		data := NetworkTunnelGroupHubs{}
		if value := res.Get("id"); value.Exists() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("isPrimary"); value.Exists() {
			data.IsPrimary = types.BoolValue(value.Bool())
		} else {
			data.IsPrimary = types.BoolNull()
		}
		if value := res.Get("authId"); value.Exists() {
			data.TunnelId = types.StringValue(value.String())
		} else {
			data.TunnelId = types.StringNull()
		}
		if value := res.Get("datacenter.name"); value.Exists() {
			data.DatacenterName = types.StringValue(value.String())
		} else {
			data.DatacenterName = types.StringNull()
		}
		if value := res.Get("datacenter.ip"); value.Exists() {
			data.DatacenterIp = types.StringValue(value.String())
		} else {
			data.DatacenterIp = types.StringNull()
		}
		items = append(items, data)
		return true
	})
	list, diags := types.ListValueFrom(ctx, elemType, items)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("failed to convert hubs: %v", diags))
	}
	data.Hubs = list
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *NetworkTunnelGroup) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("deviceType"); value.Exists() && !data.DeviceType.IsNull() {
		data.DeviceType = types.StringValue(value.String())
	} else {
		data.DeviceType = types.StringNull()
	}
	if value := res.Get("region"); value.Exists() && !data.Region.IsNull() {
		data.Region = types.StringValue(value.String())
	} else {
		data.Region = types.StringNull()
	}
	if value := res.Get("backupRegion"); value.Exists() && !data.BackupRegion.IsNull() {
		data.BackupRegion = types.StringValue(value.String())
	} else {
		data.BackupRegion = types.StringNull()
	}
	if value := res.Get("routing.type"); value.Exists() && !data.RoutingType.IsNull() {
		data.RoutingType = types.StringValue(value.String())
	} else {
		data.RoutingType = types.StringNull()
	}
	if value := res.Get("routing.data.networkCIDRs"); value.Exists() && !data.StaticNetworkCidrs.IsNull() {
		data.StaticNetworkCidrs = helpers.GetStringSet(value.Array())
	} else {
		data.StaticNetworkCidrs = types.SetNull(types.StringType)
	}
	if value := res.Get("routing.data.asNumber"); value.Exists() && !data.BgpAsNumber.IsNull() {
		data.BgpAsNumber = types.StringValue(value.String())
	} else {
		data.BgpAsNumber = types.StringNull()
	}
	if value := res.Get("routing.data.bgpHopCount"); value.Exists() && !data.BgpMultihopCount.IsNull() {
		data.BgpMultihopCount = types.Int64Value(value.Int())
	} else {
		data.BgpMultihopCount = types.Int64Null()
	}
	if value := res.Get("routing.data.bgpNeighborCIDRs"); value.Exists() && !data.BgpMultihopIpRanges.IsNull() {
		data.BgpMultihopIpRanges = helpers.GetStringSet(value.Array())
	} else {
		data.BgpMultihopIpRanges = types.SetNull(types.StringType)
	}
	if value := res.Get("routing.data.bgpPeeringTypes"); value.Exists() && !data.BgpPeeringTypes.IsNull() {
		data.BgpPeeringTypes = helpers.GetStringSet(value.Array())
	} else {
		data.BgpPeeringTypes = types.SetNull(types.StringType)
	}
	if value := res.Get("resourceAttributes.isOverrideBGPSummarization"); value.Exists() && !data.BgpOverrideSummarization.IsNull() {
		data.BgpOverrideSummarization = types.BoolValue(value.Bool())
	} else {
		data.BgpOverrideSummarization = types.BoolNull()
	}
	if value := res.Get("routing.data.enableIPv6"); value.Exists() && !data.Ipv6Routing.IsNull() {
		data.Ipv6Routing = types.BoolValue(value.Bool())
	} else {
		data.Ipv6Routing = types.BoolNull()
	}
	data.fromBodyHubs(ctx, res)
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *NetworkTunnelGroup) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	if data.Hubs.IsUnknown() {
		data.fromBodyHubs(ctx, res)
	}
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
