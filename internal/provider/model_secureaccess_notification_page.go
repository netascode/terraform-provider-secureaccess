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

type NotificationPage struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	IsDefault          types.Bool   `tfsdk:"is_default"`
	IsSwgDefault       types.Bool   `tfsdk:"is_swg_default"`
	ContactForm        types.Bool   `tfsdk:"contact_form"`
	ShowDiagnosticInfo types.Bool   `tfsdk:"show_diagnostic_info"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data NotificationPage) getPath() string {
	return "/policies/v2/landerSettings"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *NotificationPage) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("isDefault"); value.Exists() {
		data.IsDefault = types.BoolValue(value.Bool())
	} else {
		data.IsDefault = types.BoolNull()
	}
	if value := res.Get("isSwgDefault"); value.Exists() {
		data.IsSwgDefault = types.BoolValue(value.Bool())
	} else {
		data.IsSwgDefault = types.BoolNull()
	}
	if value := res.Get("contactForm"); value.Exists() {
		data.ContactForm = types.BoolValue(value.Bool())
	} else {
		data.ContactForm = types.BoolNull()
	}
	if value := res.Get("showDiagnosticInfo"); value.Exists() {
		data.ShowDiagnosticInfo = types.BoolValue(value.Bool())
	} else {
		data.ShowDiagnosticInfo = types.BoolNull()
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
