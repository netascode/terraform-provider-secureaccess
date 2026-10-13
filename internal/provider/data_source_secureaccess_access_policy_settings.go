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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AccessPolicySettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &AccessPolicySettingsDataSource{}
)

func NewAccessPolicySettingsDataSource() datasource.DataSource {
	return &AccessPolicySettingsDataSource{}
}

type AccessPolicySettingsDataSource struct {
	client *secureaccess.Client
}

func (d *AccessPolicySettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_policy_settings"
}

func (d *AccessPolicySettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the Access Policy Settings.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Optional:            true,
				Computed:            true,
			},
			"zero_trust_access_authentication": schema.BoolAttribute{
				MarkdownDescription: "Enable ZTA authentication timeout.",
				Computed:            true,
			},
			"zero_trust_access_authentication_timeout": schema.Int64Attribute{
				MarkdownDescription: "Set the ZTA authentication timeout in minutes.",
				Computed:            true,
			},
			"editing_session_duration": schema.Int64Attribute{
				MarkdownDescription: "How long an editing session lasts before it automatically expires.",
				Computed:            true,
			},
			"editing_session_extension": schema.Int64Attribute{
				MarkdownDescription: "Allow admins to extend their editing session once it's close to expiring. (minutes per extension)",
				Computed:            true,
			},
			"global_ips_profile_id": schema.StringAttribute{
				MarkdownDescription: "Set the IPS posture profile ID.",
				Computed:            true,
			},
			"microsoft_365_compatibility": schema.BoolAttribute{
				MarkdownDescription: "Enable or disable Microsoft 365 compatibility.",
				Computed:            true,
			},
			"decryption_global": schema.BoolAttribute{
				MarkdownDescription: "Enable or disable global IPS.",
				Computed:            true,
			},
			"decryption_disable_for_sources": schema.ListAttribute{
				MarkdownDescription: "List of sources for which decryption is disabled. Accepted formats include IP addresses and CIDR blocks.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"log_decrypted_internet_destinations": schema.BoolAttribute{
				MarkdownDescription: "Log decrypted traffic to internet destinations.",
				Computed:            true,
			},
			"log_decrypted_private_resources": schema.BoolAttribute{
				MarkdownDescription: "Log decrypted traffic to private resources.",
				Computed:            true,
			},
		},
	}
}

func (d *AccessPolicySettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (d *AccessPolicySettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AccessPolicySettings

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.String()))

	urlPath := config.getPath()
	res, err := d.client.Get(ctx, urlPath, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}

	config.fromBody(ctx, res)
	config.Id = types.StringValue("access_policy_settings")

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read
