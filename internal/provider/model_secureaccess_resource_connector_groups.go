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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type ResourceConnectorGroups struct {
	Filter types.String                            `tfsdk:"filter"`
	Items  map[string]ResourceConnectorGroupsItems `tfsdk:"items"`
}

type ResourceConnectorGroupsItems struct {
	Id                          types.String `tfsdk:"id"`
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

func (data ResourceConnectorGroups) getPath() string {
	return "/deployments/v2/connectorGroups"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ResourceConnectorGroups) fromBody(ctx context.Context, res gjson.Result) {
	// Build lookup map for O(1) access
	itemsByName := make(map[string]gjson.Result)
	res.ForEach(func(_, v gjson.Result) bool {
		if name := v.Get("name").String(); name != "" {
			// Item names are matched case-insensitively
			itemsByName[strings.ToLower(name)] = v
		}
		return true
	})
	for k := range data.Items {
		parent := &data
		data := (*parent).Items[k]
		res, found := itemsByName[strings.ToLower(k)]
		if !found {
			tflog.Debug(ctx, fmt.Sprintf("subresource not found, removing: name=%v", k))
			delete((*parent).Items, k)
			continue
		}
		if value := res.Get("id"); value.Exists() {
			data.Id = types.StringValue(value.String())
		} else {
			data.Id = types.StringNull()
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
		(*parent).Items[k] = data
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
