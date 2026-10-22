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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &NetworkTunnelGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkTunnelGroupDataSource{}
)

func NewNetworkTunnelGroupDataSource() datasource.DataSource {
	return &NetworkTunnelGroupDataSource{}
}

type NetworkTunnelGroupDataSource struct {
	client *secureaccess.Client
}

func (d *NetworkTunnelGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_tunnel_group"
}

func (d *NetworkTunnelGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the Network Tunnel Group.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Network Tunnel Group. It cannot exceed 50 characters, and can only contain alphanumeric characters, spaces, and hyphens.",
				Computed:            true,
			},
			"device_type": schema.StringAttribute{
				MarkdownDescription: "The type of device that establishes the network tunnel.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The name of the region that the system uses to obtain the primary and secondary data centers for the Hubs.",
				Computed:            true,
			},
			"tunnel_id_format_email": schema.StringAttribute{
				MarkdownDescription: "Prefix for the authentication ID in email format (<tunnel_id_format_email>@<org><hub>.sse.cisco.com). The string should not have any special characters besides the period(.), underscore(_), and dash(-) characters.",
				Computed:            true,
			},
			"tunnel_id_format_ip_addresses": schema.ListAttribute{
				MarkdownDescription: "List of IP addresses for the authentication ID in IP address format. Please include exactly two IP addresses.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"preshared_key": schema.StringAttribute{
				MarkdownDescription: "Passphrase for the network tunnel. It must include at least one upper case letter, one lower case letter, one number, and cannot include any special characters.",
				Computed:            true,
				Sensitive:           true,
			},
			"routing_type": schema.StringAttribute{
				MarkdownDescription: "Type of routing for the network tunnel.",
				Computed:            true,
			},
			"routing_network_cidrs": schema.SetAttribute{
				MarkdownDescription: "",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"hubs": schema.ListNestedAttribute{
				MarkdownDescription: "Remote connection endpoints for connecting network tunnel group",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "ID of the hub.",
							Computed:            true,
						},
						"is_primary": schema.BoolAttribute{
							MarkdownDescription: "Indicates if the hub is the primary hub.",
							Computed:            true,
						},
						"tunnel_id": schema.StringAttribute{
							MarkdownDescription: "Tunnel ID",
							Computed:            true,
						},
						"datacenter_name": schema.StringAttribute{
							MarkdownDescription: "Name of the datacenter.",
							Computed:            true,
						},
						"datacenter_ip": schema.StringAttribute{
							MarkdownDescription: "IP address of the datacenter.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *NetworkTunnelGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (d *NetworkTunnelGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NetworkTunnelGroup

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
