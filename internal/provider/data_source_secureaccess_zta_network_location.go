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
	_ datasource.DataSource              = &ZTANetworkLocationDataSource{}
	_ datasource.DataSourceWithConfigure = &ZTANetworkLocationDataSource{}
)

func NewZTANetworkLocationDataSource() datasource.DataSource {
	return &ZTANetworkLocationDataSource{}
}

type ZTANetworkLocationDataSource struct {
	client *secureaccess.Client
}

func (d *ZTANetworkLocationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zta_network_location"
}

func (d *ZTANetworkLocationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the ZTA Network Location.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the Network Location. Must be unique within the organization. Must be alphanumeric or space.").String,
				Optional:            true,
				Computed:            true,
			},
			"egress_criteria": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Egress criteria evaluated against the source IP observed by Secure Access cloud.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source_type": schema.StringAttribute{
							MarkdownDescription: "Type of the criterion.",
							Computed:            true,
						},
						"manual_address": schema.StringAttribute{
							MarkdownDescription: "Single IPv4 or IPv6 address, or IPv4 or IPv6 CIDR. Address ranges are not accepted. Required when `source_type` is `manual`.",
							Computed:            true,
						},
						"network_object_id": schema.StringAttribute{
							MarkdownDescription: "Network Object ID. Required when `source_type` is `networkObject`.",
							Computed:            true,
						},
						"network_object_group_id": schema.StringAttribute{
							MarkdownDescription: "Network Object Group ID. Required when `source_type` is `networkObjectGroup`.",
							Computed:            true,
						},
					},
				},
			},
			"egress_exclusion_criteria": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Egress criteria that exclude matching source IPs from the Network Location.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source_type": schema.StringAttribute{
							MarkdownDescription: "Type of the criterion.",
							Computed:            true,
						},
						"manual_address": schema.StringAttribute{
							MarkdownDescription: "Single IPv4 or IPv6 address, or IPv4 or IPv6 CIDR. Address ranges are not accepted. Required when `source_type` is `manual`.",
							Computed:            true,
						},
						"network_object_id": schema.StringAttribute{
							MarkdownDescription: "Network Object ID. Required when `source_type` is `networkObject`.",
							Computed:            true,
						},
						"network_object_group_id": schema.StringAttribute{
							MarkdownDescription: "Network Object Group ID. Required when `source_type` is `networkObjectGroup`.",
							Computed:            true,
						},
					},
				},
			},
			"private_subnet_criteria": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Private or internal subnet criteria.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source_type": schema.StringAttribute{
							MarkdownDescription: "Type of the criterion.",
							Computed:            true,
						},
						"manual_address": schema.StringAttribute{
							MarkdownDescription: "Single IPv4 or IPv6 address, or IPv4 or IPv6 CIDR. Address ranges are not accepted. Required when `source_type` is `manual`.",
							Computed:            true,
						},
						"network_object_id": schema.StringAttribute{
							MarkdownDescription: "Network Object ID. Required when `source_type` is `networkObject`.",
							Computed:            true,
						},
						"network_object_group_id": schema.StringAttribute{
							MarkdownDescription: "Network Object Group ID. Required when `source_type` is `networkObjectGroup`.",
							Computed:            true,
						},
					},
				},
			},
			"revision": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Revision of the configuration.").String,
				Computed:            true,
			},
		},
	}
}
func (d *ZTANetworkLocationDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *ZTANetworkLocationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (d *ZTANetworkLocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ZTANetworkLocation

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.String()))

	// Object found by a query attribute, it is not retrieved again by ID
	var res secureaccess.Res
	if config.Id.IsNull() && !config.Name.IsNull() {
		// Client handles pagination and returns all items merged into a single list
		items, err := d.client.GetAll(ctx, config.getPath(), reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, got error: %s", err))
			return
		}
		if !items.IsArray() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, expected a list, got: %s", items.Type))
			return
		}
		items.ForEach(func(k, v gjson.Result) bool {
			if config.Name.ValueString() == v.Get("networkLocationName").String() {
				config.Id = types.StringValue(v.Get("networkLocationId").String())
				res.Result = v
				tflog.Debug(ctx, fmt.Sprintf("%s: Found object with name '%v', id: %v", config.Id.ValueString(), config.Name.ValueString(), config.Id.ValueString()))
				return false
			}
			return true
		})

		if config.Id.IsNull() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to find object with name: %v", config.Name.ValueString()))
			return
		}
	}

	if !res.Exists() {
		urlPath := config.getPath() + "/" + url.QueryEscape(config.Id.ValueString())
		var err error
		res, err = d.client.Get(ctx, urlPath, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
			return
		}
	}

	config.fromBody(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read
