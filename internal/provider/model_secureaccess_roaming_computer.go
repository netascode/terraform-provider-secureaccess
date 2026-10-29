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

type RoamingComputer struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	DeviceId               types.String `tfsdk:"device_id"`
	AnyconnectDeviceId     types.String `tfsdk:"anyconnect_device_id"`
	OriginId               types.Int64  `tfsdk:"origin_id"`
	AppliedBundle          types.Int64  `tfsdk:"applied_bundle"`
	Type                   types.String `tfsdk:"type"`
	OsVersion              types.String `tfsdk:"os_version"`
	OsVersionName          types.String `tfsdk:"os_version_name"`
	HasIpBlocking          types.Bool   `tfsdk:"has_ip_blocking"`
	Status                 types.String `tfsdk:"status"`
	SecureWebGatewayStatus types.String `tfsdk:"secure_web_gateway_status"`
	SecureClientVersion    types.String `tfsdk:"secure_client_version"`
	DeleteOnDestroy        types.Bool   `tfsdk:"delete_on_destroy"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data RoamingComputer) getPath() string {
	return "/deployments/v2/roamingcomputers"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data RoamingComputer) toBody(ctx context.Context, state RoamingComputer) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "deviceId", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// toBodyPut builds the PUT body with the name, if it differs from the state.
// Returns an empty string if there is nothing to update.
func (data RoamingComputer) toBodyPut(ctx context.Context, state RoamingComputer) string {
	if data.Name.Equal(state.Name) {
		return ""
	}
	body, _ := sjson.Set("", "name", data.Name.ValueString())
	return body
}

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *RoamingComputer) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("deviceId"); value.Exists() {
		data.DeviceId = types.StringValue(value.String())
	} else {
		data.DeviceId = types.StringNull()
	}
	if value := res.Get("anyconnectDeviceId"); value.Exists() {
		data.AnyconnectDeviceId = types.StringValue(value.String())
	} else {
		data.AnyconnectDeviceId = types.StringNull()
	}
	if value := res.Get("originId"); value.Exists() {
		data.OriginId = types.Int64Value(value.Int())
	} else {
		data.OriginId = types.Int64Null()
	}
	if value := res.Get("appliedBundle"); value.Exists() {
		data.AppliedBundle = types.Int64Value(value.Int())
	} else {
		data.AppliedBundle = types.Int64Null()
	}
	if value := res.Get("type"); value.Exists() {
		data.Type = types.StringValue(value.String())
	} else {
		data.Type = types.StringNull()
	}
	if value := res.Get("osVersion"); value.Exists() {
		data.OsVersion = types.StringValue(value.String())
	} else {
		data.OsVersion = types.StringNull()
	}
	if value := res.Get("osVersionName"); value.Exists() {
		data.OsVersionName = types.StringValue(value.String())
	} else {
		data.OsVersionName = types.StringNull()
	}
	if value := res.Get("hasIpBlocking"); value.Exists() {
		data.HasIpBlocking = types.BoolValue(value.Bool())
	} else {
		data.HasIpBlocking = types.BoolNull()
	}
	if value := res.Get("status"); value.Exists() {
		data.Status = types.StringValue(value.String())
	} else {
		data.Status = types.StringNull()
	}
	if value := res.Get("swgStatus"); value.Exists() {
		data.SecureWebGatewayStatus = types.StringValue(value.String())
	} else {
		data.SecureWebGatewayStatus = types.StringNull()
	}
	if value := res.Get("version"); value.Exists() {
		data.SecureClientVersion = types.StringValue(value.String())
	} else {
		data.SecureClientVersion = types.StringNull()
	}
	data.DeleteOnDestroy = types.BoolValue(false)
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *RoamingComputer) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("deviceId"); value.Exists() && !data.DeviceId.IsNull() {
		data.DeviceId = types.StringValue(value.String())
	} else {
		data.DeviceId = types.StringNull()
	}
	if value := res.Get("anyconnectDeviceId"); value.Exists() && !data.AnyconnectDeviceId.IsNull() {
		data.AnyconnectDeviceId = types.StringValue(value.String())
	} else {
		data.AnyconnectDeviceId = types.StringNull()
	}
	if value := res.Get("originId"); value.Exists() && !data.OriginId.IsNull() {
		data.OriginId = types.Int64Value(value.Int())
	} else {
		data.OriginId = types.Int64Null()
	}
	if value := res.Get("appliedBundle"); value.Exists() && !data.AppliedBundle.IsNull() {
		data.AppliedBundle = types.Int64Value(value.Int())
	} else {
		data.AppliedBundle = types.Int64Null()
	}
	if value := res.Get("type"); value.Exists() && !data.Type.IsNull() {
		data.Type = types.StringValue(value.String())
	} else {
		data.Type = types.StringNull()
	}
	if value := res.Get("osVersion"); value.Exists() && !data.OsVersion.IsNull() {
		data.OsVersion = types.StringValue(value.String())
	} else {
		data.OsVersion = types.StringNull()
	}
	if value := res.Get("osVersionName"); value.Exists() && !data.OsVersionName.IsNull() {
		data.OsVersionName = types.StringValue(value.String())
	} else {
		data.OsVersionName = types.StringNull()
	}
	if value := res.Get("hasIpBlocking"); value.Exists() && !data.HasIpBlocking.IsNull() {
		data.HasIpBlocking = types.BoolValue(value.Bool())
	} else {
		data.HasIpBlocking = types.BoolNull()
	}
	if value := res.Get("status"); value.Exists() && !data.Status.IsNull() {
		data.Status = types.StringValue(value.String())
	} else {
		data.Status = types.StringNull()
	}
	if value := res.Get("swgStatus"); value.Exists() && !data.SecureWebGatewayStatus.IsNull() {
		data.SecureWebGatewayStatus = types.StringValue(value.String())
	} else {
		data.SecureWebGatewayStatus = types.StringNull()
	}
	if value := res.Get("version"); value.Exists() {
		data.SecureClientVersion = types.StringValue(value.String())
	} else {
		data.SecureClientVersion = types.StringNull()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *RoamingComputer) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	if data.DeviceId.IsUnknown() {
		if value := res.Get("deviceId"); value.Exists() {
			data.DeviceId = types.StringValue(value.String())
		} else {
			data.DeviceId = types.StringNull()
		}
	}
	if data.AnyconnectDeviceId.IsUnknown() {
		if value := res.Get("anyconnectDeviceId"); value.Exists() {
			data.AnyconnectDeviceId = types.StringValue(value.String())
		} else {
			data.AnyconnectDeviceId = types.StringNull()
		}
	}
	if data.OriginId.IsUnknown() {
		if value := res.Get("originId"); value.Exists() {
			data.OriginId = types.Int64Value(value.Int())
		} else {
			data.OriginId = types.Int64Null()
		}
	}
	if data.AppliedBundle.IsUnknown() {
		if value := res.Get("appliedBundle"); value.Exists() {
			data.AppliedBundle = types.Int64Value(value.Int())
		} else {
			data.AppliedBundle = types.Int64Null()
		}
	}
	if data.Type.IsUnknown() {
		if value := res.Get("type"); value.Exists() {
			data.Type = types.StringValue(value.String())
		} else {
			data.Type = types.StringNull()
		}
	}
	if data.OsVersion.IsUnknown() {
		if value := res.Get("osVersion"); value.Exists() {
			data.OsVersion = types.StringValue(value.String())
		} else {
			data.OsVersion = types.StringNull()
		}
	}
	if data.OsVersionName.IsUnknown() {
		if value := res.Get("osVersionName"); value.Exists() {
			data.OsVersionName = types.StringValue(value.String())
		} else {
			data.OsVersionName = types.StringNull()
		}
	}
	if data.HasIpBlocking.IsUnknown() {
		if value := res.Get("hasIpBlocking"); value.Exists() {
			data.HasIpBlocking = types.BoolValue(value.Bool())
		} else {
			data.HasIpBlocking = types.BoolNull()
		}
	}
	if data.Status.IsUnknown() {
		if value := res.Get("status"); value.Exists() {
			data.Status = types.StringValue(value.String())
		} else {
			data.Status = types.StringNull()
		}
	}
	if data.SecureWebGatewayStatus.IsUnknown() {
		if value := res.Get("swgStatus"); value.Exists() {
			data.SecureWebGatewayStatus = types.StringValue(value.String())
		} else {
			data.SecureWebGatewayStatus = types.StringNull()
		}
	}
	if value := res.Get("version"); value.Exists() {
		data.SecureClientVersion = types.StringValue(value.String())
	} else {
		data.SecureClientVersion = types.StringNull()
	}
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
