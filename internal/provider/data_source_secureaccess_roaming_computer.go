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

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &RoamingComputerDataSource{}
	_ datasource.DataSourceWithConfigure = &RoamingComputerDataSource{}
)

func NewRoamingComputerDataSource() datasource.DataSource {
	return &RoamingComputerDataSource{}
}

type RoamingComputerDataSource struct {
	client *secureaccess.Client
}

func (d *RoamingComputerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roaming_computer"
}

func (d *RoamingComputerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the Roaming Computer.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Roaming Computer. Used to find the Roaming Computer on create, a later change renames it.",
				Optional:            true,
				Computed:            true,
			},
			"device_id": schema.StringAttribute{
				MarkdownDescription: "The hex ID of the Roaming Computer.",
				Computed:            true,
			},
			"anyconnect_device_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the device that has the Cisco Secure Client deployed with the Internet Security module.",
				Computed:            true,
			},
			"origin_id": schema.Int64Attribute{
				MarkdownDescription: "The origin ID of the Roaming Computer.",
				Computed:            true,
			},
			"applied_bundle": schema.Int64Attribute{
				MarkdownDescription: "The ID of the policy applied to the Roaming Computer.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the Roaming Computer.",
				Computed:            true,
			},
			"os_version": schema.StringAttribute{
				MarkdownDescription: "The OS version of the Roaming Computer.",
				Computed:            true,
			},
			"os_version_name": schema.StringAttribute{
				MarkdownDescription: "The OS version name of the Roaming Computer.",
				Computed:            true,
			},
			"has_ip_blocking": schema.BoolAttribute{
				MarkdownDescription: "Specifies whether the Roaming Computer has IP blocking.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the Roaming Computer with DNS-layer security.",
				Computed:            true,
			},
			"secure_web_gateway_status": schema.StringAttribute{
				MarkdownDescription: "The status of the Roaming Computer with Internet security (Secure Web Gateway).",
				Computed:            true,
			},
			"secure_client_version": schema.StringAttribute{
				MarkdownDescription: "The version of the Roaming Computer client.",
				Computed:            true,
			},
			"delete_on_destroy": schema.BoolAttribute{
				MarkdownDescription: "When `true`, the Roaming Computer is removed from Secure Access on destroy. When `false`, it is only removed from the Terraform state.",
				Computed:            true,
			},
		},
	}
}
func (d *RoamingComputerDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *RoamingComputerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (d *RoamingComputerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config RoamingComputer

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.String()))

	var res gjson.Result
	if config.Id.IsNull() && !config.Name.IsNull() {
		obj, err := getRoamingComputerByName(ctx, d.client, config.Name.ValueString(), reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
		config.Id = types.StringValue(obj.Get("deviceId").String())
		tflog.Debug(ctx, fmt.Sprintf("%s: Found object with name '%v'", config.Id.ValueString(), config.Name.ValueString()))
		res = obj
	} else {
		urlPath := config.getPath() + "/" + url.QueryEscape(config.Id.ValueString())
		obj, err := d.client.Get(ctx, urlPath, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
			return
		}
		res = obj
	}

	config.fromBody(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
