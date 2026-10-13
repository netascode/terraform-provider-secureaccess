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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type ResourceConnector struct {
	Id                       types.String `tfsdk:"id"`
	InstanceId               types.String `tfsdk:"instance_id"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	Confirmed                types.Bool   `tfsdk:"confirmed"`
	Hostname                 types.String `tfsdk:"hostname"`
	Status                   types.String `tfsdk:"status"`
	ResourceConnectorGroupId types.Int64  `tfsdk:"resource_connector_group_id"`
	Version                  types.String `tfsdk:"version"`
	DeleteOnDestroy          types.Bool   `tfsdk:"delete_on_destroy"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ResourceConnector) getPath() string {
	return "/deployments/v2/connectorAgents"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ResourceConnector) toBody(ctx context.Context, state ResourceConnector) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if !data.InstanceId.IsNull() {
		body, _ = sjson.Set(body, "instanceId", data.InstanceId.ValueString())
	}
	if !data.Enabled.IsNull() {
		body, _ = sjson.Set(body, "enabled", data.Enabled.ValueBool())
	}
	if !data.Confirmed.IsNull() {
		body, _ = sjson.Set(body, "confirmed", data.Confirmed.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// resourceConnectorPatchPaths lists top level body keys that can be updated with PATCH
var resourceConnectorPatchPaths = []string{"enabled", "confirmed"}

// toBodyPatch builds JSON Patch "replace" operations for the paths which differ between data and current.
// current is either the prior state (Update) or the object as retrieved from the backend (Create).
// Returns an empty string if there is nothing to update.
func (data ResourceConnector) toBodyPatch(ctx context.Context, current ResourceConnector) string {
	planBody := data.toBody(ctx, ResourceConnector{})
	currentBody := current.toBody(ctx, ResourceConnector{})
	ops := "[]"
	for _, p := range resourceConnectorPatchPaths {
		planValue := gjson.Get(planBody, p)
		currentValue := gjson.Get(currentBody, p)
		if !planValue.Exists() {
			if currentValue.Exists() {
				tflog.Debug(ctx, fmt.Sprintf("%s: %s not configured, leaving it unchanged", data.Id.ValueString(), p))
			}
			continue
		}
		if currentValue.Exists() && planValue.Raw == currentValue.Raw {
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

func (data *ResourceConnector) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("instanceId"); value.Exists() {
		data.InstanceId = types.StringValue(value.String())
	} else {
		data.InstanceId = types.StringNull()
	}
	if value := res.Get("enabled"); value.Exists() {
		data.Enabled = types.BoolValue(value.Bool())
	} else {
		data.Enabled = types.BoolNull()
	}
	if value := res.Get("confirmed"); value.Exists() {
		data.Confirmed = types.BoolValue(value.Bool())
	} else {
		data.Confirmed = types.BoolNull()
	}
	if value := res.Get("hostname"); value.Exists() {
		data.Hostname = types.StringValue(value.String())
	} else {
		data.Hostname = types.StringNull()
	}
	if value := res.Get("status"); value.Exists() {
		data.Status = types.StringValue(value.String())
	} else {
		data.Status = types.StringNull()
	}
	if value := res.Get("groupId"); value.Exists() {
		data.ResourceConnectorGroupId = types.Int64Value(value.Int())
	} else {
		data.ResourceConnectorGroupId = types.Int64Null()
	}
	if value := res.Get("version"); value.Exists() {
		data.Version = types.StringValue(value.String())
	} else {
		data.Version = types.StringNull()
	}
	data.DeleteOnDestroy = types.BoolValue(false)
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *ResourceConnector) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("instanceId"); value.Exists() && !data.InstanceId.IsNull() {
		data.InstanceId = types.StringValue(value.String())
	} else {
		data.InstanceId = types.StringNull()
	}
	if value := res.Get("enabled"); value.Exists() && !data.Enabled.IsNull() {
		data.Enabled = types.BoolValue(value.Bool())
	} else {
		data.Enabled = types.BoolNull()
	}
	if value := res.Get("confirmed"); value.Exists() && !data.Confirmed.IsNull() {
		data.Confirmed = types.BoolValue(value.Bool())
	} else {
		data.Confirmed = types.BoolNull()
	}
	if value := res.Get("hostname"); value.Exists() && !data.Hostname.IsNull() {
		data.Hostname = types.StringValue(value.String())
	} else {
		data.Hostname = types.StringNull()
	}
	if value := res.Get("status"); value.Exists() && !data.Status.IsNull() {
		data.Status = types.StringValue(value.String())
	} else {
		data.Status = types.StringNull()
	}
	if value := res.Get("groupId"); value.Exists() && !data.ResourceConnectorGroupId.IsNull() {
		data.ResourceConnectorGroupId = types.Int64Value(value.Int())
	} else {
		data.ResourceConnectorGroupId = types.Int64Null()
	}
	if value := res.Get("version"); value.Exists() {
		data.Version = types.StringValue(value.String())
	} else {
		data.Version = types.StringNull()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ResourceConnector) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	if data.Hostname.IsUnknown() {
		if value := res.Get("hostname"); value.Exists() {
			data.Hostname = types.StringValue(value.String())
		} else {
			data.Hostname = types.StringNull()
		}
	}
	if data.Status.IsUnknown() {
		if value := res.Get("status"); value.Exists() {
			data.Status = types.StringValue(value.String())
		} else {
			data.Status = types.StringNull()
		}
	}
	if data.ResourceConnectorGroupId.IsUnknown() {
		if value := res.Get("groupId"); value.Exists() {
			data.ResourceConnectorGroupId = types.Int64Value(value.Int())
		} else {
			data.ResourceConnectorGroupId = types.Int64Null()
		}
	}
	if value := res.Get("version"); value.Exists() {
		data.Version = types.StringValue(value.String())
	} else {
		data.Version = types.StringNull()
	}
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
