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
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &SecurityProfilePolicySettingDataSource{}
	_ datasource.DataSourceWithConfigure = &SecurityProfilePolicySettingDataSource{}
)

func NewSecurityProfilePolicySettingDataSource() datasource.DataSource {
	return &SecurityProfilePolicySettingDataSource{}
}

type SecurityProfilePolicySettingDataSource struct {
	client *secureaccess.Client
}

func (d *SecurityProfilePolicySettingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_profile_policy_setting"
}

func (d *SecurityProfilePolicySettingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the Policy Setting of a Security Profile, found by its `id`. The `id` of the Policy Setting of a Security Profile is the `policySettingGroupId` of that profile.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the Policy Setting.").String,
				Computed:            true,
			},
			"decryption_enabled": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Specifies whether Decryption is enabled (Web UI Decryption, Enabled). According to the Security Profiles API guide, enabling Decryption sets it to true.").String,
				Computed:            true,
			},
			"intelligent_proxy": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Specifies whether the intelligent proxy is enabled. According to the Security Profiles API guide, the API also sets it to true when Decryption or Block page mode is enabled, but it can be true without either of them.").String,
				Computed:            true,
			},
			"end_user_notification": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Specifies whether the Block page mode of the End-User Notifications is enabled (Web UI End-User Notifications, Block page mode). The API returns 0 or 1, which is read as `false` or `true`.").String,
				Computed:            true,
			},
			"single_sign_on": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Specifies whether Single Sign-On (SSO) is enabled. The API returns 0 or 1, which is read as `false` or `true`.").String,
				Computed:            true,
			},
			"safe_search": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Specifies whether SafeSearch is enabled.").String,
				Computed:            true,
			},
			"ai_supply_chain_categories": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Mask of the blocked AI supply chain categories (Web UI AI Supply Chain Blocking). The API returns `0` if the blocking is off, otherwise a mask string, according to the Security Profiles API guide its first digit adds up 2 for Prohibited Suppliers, 4 for Code Execution and 8 for Copyleft License (as a hexadecimal digit) and the 53 digits after it are 0.").String,
				Computed:            true,
			},
			"reporting": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Reporting setting of the Policy Setting.").String,
				Computed:            true,
			},
		},
	}
}

func (d *SecurityProfilePolicySettingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (d *SecurityProfilePolicySettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config SecurityProfilePolicySetting

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.String()))

	urlPath := config.getPath() + "/" + url.QueryEscape(config.Id.ValueString())
	res, err := d.client.Get(ctx, urlPath, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}

	config.fromBody(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read
