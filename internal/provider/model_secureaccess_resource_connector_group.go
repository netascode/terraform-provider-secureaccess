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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type ResourceConnectorGroup struct {
	Id                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Region                      types.String `tfsdk:"region"`
	Environment                 types.String `tfsdk:"environment"`
	ProvisioningKeyExpiresAt    types.String `tfsdk:"provisioning_key_expires_at"`
	Status                      types.String `tfsdk:"status"`
	ConnectorsCount             types.Int64  `tfsdk:"connectors_count"`
	ConnectedConnectorsCount    types.Int64  `tfsdk:"connected_connectors_count"`
	DisconnectedConnectorsCount types.Int64  `tfsdk:"disconnected_connectors_count"`
	UnableToSyncConnectorsCount types.Int64  `tfsdk:"unable_to_sync_connectors_count"`
	NeedUpdateConnectorsCount   types.Int64  `tfsdk:"need_update_connectors_count"`
	FailedUpdateAgentsAttempts  types.Int64  `tfsdk:"failed_update_agents_attempts"`
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
	return body
}

// End of section. //template:end toBody

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
