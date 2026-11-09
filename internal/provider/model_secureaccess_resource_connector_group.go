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

type ResourceConnectorGroup struct {
	Id                          types.String                       `tfsdk:"id"`
	Name                        types.String                       `tfsdk:"name"`
	Region                      types.String                       `tfsdk:"region"`
	Environment                 types.String                       `tfsdk:"environment"`
	DnsServers                  []ResourceConnectorGroupDnsServers `tfsdk:"dns_servers"`
	PrivateResourceIds          types.List                         `tfsdk:"private_resource_ids"`
	ProvisioningKeyExpiresAt    types.String                       `tfsdk:"provisioning_key_expires_at"`
	Status                      types.String                       `tfsdk:"status"`
	ConnectorsCount             types.Int64                        `tfsdk:"connectors_count"`
	ConnectedConnectorsCount    types.Int64                        `tfsdk:"connected_connectors_count"`
	DisconnectedConnectorsCount types.Int64                        `tfsdk:"disconnected_connectors_count"`
	UnableToSyncConnectorsCount types.Int64                        `tfsdk:"unable_to_sync_connectors_count"`
	NeedUpdateConnectorsCount   types.Int64                        `tfsdk:"need_update_connectors_count"`
	FailedUpdateAgentsAttempts  types.Int64                        `tfsdk:"failed_update_agents_attempts"`
}

type ResourceConnectorGroupDnsServers struct {
	Domains     types.String `tfsdk:"domains"`
	DnsServerId types.Int64  `tfsdk:"dns_server_id"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ResourceConnectorGroup) getPath() string {
	return "/deployments/v2/connectorGroups"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ResourceConnectorGroup) toBody(ctx context.Context, state ResourceConnectorGroup) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	if !data.Region.IsNull() {
		body, _ = sjson.Set(body, "location", data.Region.ValueString())
	}
	if !data.Environment.IsNull() {
		body, _ = sjson.Set(body, "environment", data.Environment.ValueString())
	}
	if len(data.DnsServers) > 0 {
		var dnsServersBody strings.Builder
		dnsServersBody.WriteString("[")
		for _, item := range data.DnsServers {
			itemBody := ""
			if !item.Domains.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "domains", item.Domains.ValueString())
			}
			if !item.DnsServerId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dnsResourceId", item.DnsServerId.ValueInt64())
			}
			if itemBody != "" {
				if dnsServersBody.Len() > 1 {
					dnsServersBody.WriteString(",")
				}
				dnsServersBody.WriteString(itemBody)
			}
		}
		dnsServersBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "forwardDNS", dnsServersBody.String())
	}
	if !data.PrivateResourceIds.IsNull() {
		var values []int64
		data.PrivateResourceIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "resourceIds", values)
	}
	return body
}

// End of section. //template:end toBody

// resourceConnectorGroupPatchPaths lists top level body keys that can be updated with PATCH
var resourceConnectorGroupPatchPaths = []string{"name", "forwardDNS", "resourceIds"}

// toBodyPatch builds JSON Patch "replace" operations for the paths which differ between data and state.
// Returns an empty string if there is nothing to update.
func (data ResourceConnectorGroup) toBodyPatch(ctx context.Context, state ResourceConnectorGroup) string {
	planBody := data.toBody(ctx, ResourceConnectorGroup{})
	stateBody := state.toBody(ctx, ResourceConnectorGroup{})
	ops := "[]"
	for _, p := range resourceConnectorGroupPatchPaths {
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

func (data *ResourceConnectorGroup) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("location"); value.Exists() {
		data.Region = types.StringValue(value.String())
	} else {
		data.Region = types.StringNull()
	}
	if value := res.Get("environment"); value.Exists() {
		data.Environment = types.StringValue(value.String())
	} else {
		data.Environment = types.StringNull()
	}
	if value := res.Get("forwardDNS"); value.Exists() {
		data.DnsServers = make([]ResourceConnectorGroupDnsServers, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ResourceConnectorGroupDnsServers{}
			if value := res.Get("domains"); value.Exists() {
				data.Domains = types.StringValue(value.String())
			} else {
				data.Domains = types.StringNull()
			}
			if value := res.Get("dnsResourceId"); value.Exists() {
				data.DnsServerId = types.Int64Value(value.Int())
			} else {
				data.DnsServerId = types.Int64Null()
			}
			(*parent).DnsServers = append((*parent).DnsServers, data)
			return true
		})
	}
	if value := res.Get("resourceIds"); value.Exists() {
		data.PrivateResourceIds = helpers.GetInt64List(value.Array())
	} else {
		data.PrivateResourceIds = types.ListNull(types.Int64Type)
	}
	if value := res.Get("provisioningKeyExpiresAt"); value.Exists() {
		data.ProvisioningKeyExpiresAt = types.StringValue(value.String())
	} else {
		data.ProvisioningKeyExpiresAt = types.StringNull()
	}
	if value := res.Get("status"); value.Exists() {
		data.Status = types.StringValue(value.String())
	} else {
		data.Status = types.StringNull()
	}
	if value := res.Get("connectorsCount"); value.Exists() {
		data.ConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.ConnectorsCount = types.Int64Null()
	}
	if value := res.Get("connectedConnectorsCount"); value.Exists() {
		data.ConnectedConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.ConnectedConnectorsCount = types.Int64Null()
	}
	if value := res.Get("disconnectedConnectorsCount"); value.Exists() {
		data.DisconnectedConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.DisconnectedConnectorsCount = types.Int64Null()
	}
	if value := res.Get("unableToSyncConnectorsCount"); value.Exists() {
		data.UnableToSyncConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.UnableToSyncConnectorsCount = types.Int64Null()
	}
	if value := res.Get("needUpdateConnectorsCount"); value.Exists() {
		data.NeedUpdateConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.NeedUpdateConnectorsCount = types.Int64Null()
	}
	if value := res.Get("failedUpdateAgentsAttempts"); value.Exists() {
		data.FailedUpdateAgentsAttempts = types.Int64Value(value.Int())
	} else {
		data.FailedUpdateAgentsAttempts = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *ResourceConnectorGroup) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("location"); value.Exists() && !data.Region.IsNull() {
		data.Region = types.StringValue(value.String())
	} else {
		data.Region = types.StringNull()
	}
	if value := res.Get("environment"); value.Exists() && !data.Environment.IsNull() {
		data.Environment = types.StringValue(value.String())
	} else {
		data.Environment = types.StringNull()
	}
	dnsServersArray := res.Get("forwardDNS")
	for i := 0; i < len(data.DnsServers); i++ {
		keys := [...]string{"domains", "dnsResourceId"}
		keyValues := [...]string{data.DnsServers[i].Domains.ValueString(), strconv.FormatInt(data.DnsServers[i].DnsServerId.ValueInt64(), 10)}

		parent := &data
		data := (*parent).DnsServers[i]
		var res gjson.Result

		dnsServersArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing DnsServers[%d] = %+v",
				i,
				(*parent).DnsServers[i],
			))
			(*parent).DnsServers = slices.Delete((*parent).DnsServers, i, i+1)
			i--

			continue
		}
		if value := res.Get("domains"); value.Exists() && !data.Domains.IsNull() {
			data.Domains = types.StringValue(value.String())
		} else {
			data.Domains = types.StringNull()
		}
		if value := res.Get("dnsResourceId"); value.Exists() && !data.DnsServerId.IsNull() {
			data.DnsServerId = types.Int64Value(value.Int())
		} else {
			data.DnsServerId = types.Int64Null()
		}
		(*parent).DnsServers[i] = data
	}
	if value := res.Get("resourceIds"); value.Exists() && !data.PrivateResourceIds.IsNull() {
		data.PrivateResourceIds = helpers.GetInt64List(value.Array())
	} else {
		data.PrivateResourceIds = types.ListNull(types.Int64Type)
	}
	if value := res.Get("provisioningKeyExpiresAt"); value.Exists() {
		data.ProvisioningKeyExpiresAt = types.StringValue(value.String())
	} else {
		data.ProvisioningKeyExpiresAt = types.StringNull()
	}
	if value := res.Get("status"); value.Exists() {
		data.Status = types.StringValue(value.String())
	} else {
		data.Status = types.StringNull()
	}
	if value := res.Get("connectorsCount"); value.Exists() {
		data.ConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.ConnectorsCount = types.Int64Null()
	}
	if value := res.Get("connectedConnectorsCount"); value.Exists() {
		data.ConnectedConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.ConnectedConnectorsCount = types.Int64Null()
	}
	if value := res.Get("disconnectedConnectorsCount"); value.Exists() {
		data.DisconnectedConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.DisconnectedConnectorsCount = types.Int64Null()
	}
	if value := res.Get("unableToSyncConnectorsCount"); value.Exists() {
		data.UnableToSyncConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.UnableToSyncConnectorsCount = types.Int64Null()
	}
	if value := res.Get("needUpdateConnectorsCount"); value.Exists() {
		data.NeedUpdateConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.NeedUpdateConnectorsCount = types.Int64Null()
	}
	if value := res.Get("failedUpdateAgentsAttempts"); value.Exists() {
		data.FailedUpdateAgentsAttempts = types.Int64Value(value.Int())
	} else {
		data.FailedUpdateAgentsAttempts = types.Int64Null()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ResourceConnectorGroup) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	if value := res.Get("provisioningKeyExpiresAt"); value.Exists() {
		data.ProvisioningKeyExpiresAt = types.StringValue(value.String())
	} else {
		data.ProvisioningKeyExpiresAt = types.StringNull()
	}
	if value := res.Get("status"); value.Exists() {
		data.Status = types.StringValue(value.String())
	} else {
		data.Status = types.StringNull()
	}
	if value := res.Get("connectorsCount"); value.Exists() {
		data.ConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.ConnectorsCount = types.Int64Null()
	}
	if value := res.Get("connectedConnectorsCount"); value.Exists() {
		data.ConnectedConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.ConnectedConnectorsCount = types.Int64Null()
	}
	if value := res.Get("disconnectedConnectorsCount"); value.Exists() {
		data.DisconnectedConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.DisconnectedConnectorsCount = types.Int64Null()
	}
	if value := res.Get("unableToSyncConnectorsCount"); value.Exists() {
		data.UnableToSyncConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.UnableToSyncConnectorsCount = types.Int64Null()
	}
	if value := res.Get("needUpdateConnectorsCount"); value.Exists() {
		data.NeedUpdateConnectorsCount = types.Int64Value(value.Int())
	} else {
		data.NeedUpdateConnectorsCount = types.Int64Null()
	}
	if value := res.Get("failedUpdateAgentsAttempts"); value.Exists() {
		data.FailedUpdateAgentsAttempts = types.Int64Value(value.Int())
	} else {
		data.FailedUpdateAgentsAttempts = types.Int64Null()
	}
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
