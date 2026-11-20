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
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type ApplicationExceptionList struct {
	Id           types.String                           `tfsdk:"id"`
	Name         types.String                           `tfsdk:"name"`
	AllowedCount types.Int64                            `tfsdk:"allowed_count"`
	BlockedCount types.Int64                            `tfsdk:"blocked_count"`
	WarnedCount  types.Int64                            `tfsdk:"warned_count"`
	Applications []ApplicationExceptionListApplications `tfsdk:"applications"`
}

type ApplicationExceptionListApplications struct {
	ApplicationId   types.Int64  `tfsdk:"application_id"`
	ApplicationName types.String `tfsdk:"application_name"`
	GroupId         types.Int64  `tfsdk:"group_id"`
	GroupName       types.String `tfsdk:"group_name"`
	HasUrls         types.Bool   `tfsdk:"has_urls"`
	Action          types.String `tfsdk:"action"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ApplicationExceptionList) getPath() string {
	return "/policies/v2/exceptionApplicationLists"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ApplicationExceptionList) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("allowedCount"); value.Exists() {
		data.AllowedCount = types.Int64Value(value.Int())
	} else {
		data.AllowedCount = types.Int64Null()
	}
	if value := res.Get("blockedCount"); value.Exists() {
		data.BlockedCount = types.Int64Value(value.Int())
	} else {
		data.BlockedCount = types.Int64Null()
	}
	if value := res.Get("warnedCount"); value.Exists() {
		data.WarnedCount = types.Int64Value(value.Int())
	} else {
		data.WarnedCount = types.Int64Null()
	}
	if value := res.Get("applications"); value.Exists() {
		data.Applications = make([]ApplicationExceptionListApplications, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ApplicationExceptionListApplications{}
			if value := res.Get("applicationId"); value.Exists() {
				data.ApplicationId = types.Int64Value(value.Int())
			} else {
				data.ApplicationId = types.Int64Null()
			}
			if value := res.Get("applicationName"); value.Exists() {
				data.ApplicationName = types.StringValue(value.String())
			} else {
				data.ApplicationName = types.StringNull()
			}
			if value := res.Get("groupId"); value.Exists() {
				data.GroupId = types.Int64Value(value.Int())
			} else {
				data.GroupId = types.Int64Null()
			}
			if value := res.Get("groupName"); value.Exists() {
				data.GroupName = types.StringValue(value.String())
			} else {
				data.GroupName = types.StringNull()
			}
			if value := res.Get("hasUrls"); value.Exists() {
				data.HasUrls = types.BoolValue(value.Bool())
			} else {
				data.HasUrls = types.BoolNull()
			}
			if value := res.Get("action"); value.Exists() {
				data.Action = types.StringValue(value.String())
			} else {
				data.Action = types.StringNull()
			}
			(*parent).Applications = append((*parent).Applications, data)
			return true
		})
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

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
