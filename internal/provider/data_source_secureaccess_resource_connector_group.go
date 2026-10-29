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
	_ datasource.DataSource              = &ResourceConnectorGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &ResourceConnectorGroupDataSource{}
)

func NewResourceConnectorGroupDataSource() datasource.DataSource {
	return &ResourceConnectorGroupDataSource{}
}

type ResourceConnectorGroupDataSource struct {
	client *secureaccess.Client
}

func (d *ResourceConnectorGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_connector_group"
}

func (d *ResourceConnectorGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the Resource Connector Group.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Resource Connector Group.",
				Optional:            true,
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region where to deploy Resource Connector Group on Cisco Secure Access.",
				Computed:            true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Environment in which agents will be deployed.",
				Computed:            true,
			},
			"provisioning_key_expires_at": schema.StringAttribute{
				MarkdownDescription: "The expiration date and time of the Resource Connector Group's provisioning key.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The label that describes the status of the Resource Connector Group.",
				Computed:            true,
			},
			"connectors_count": schema.Int64Attribute{
				MarkdownDescription: "The number of connectors in the Resource Connector Group.",
				Computed:            true,
			},
			"connected_connectors_count": schema.Int64Attribute{
				MarkdownDescription: "The number of Connectors in the Resource Connector Group that are connected.",
				Computed:            true,
			},
			"disconnected_connectors_count": schema.Int64Attribute{
				MarkdownDescription: "The number of Connectors in the Resource Connector Group that are disconnected.",
				Computed:            true,
			},
			"unable_to_sync_connectors_count": schema.Int64Attribute{
				MarkdownDescription: "The number of Connectors in the Resource Connector Group that are unable to sync.",
				Computed:            true,
			},
			"need_update_connectors_count": schema.Int64Attribute{
				MarkdownDescription: "The number of Connectors in the Resource Connector Group that need an update.",
				Computed:            true,
			},
			"failed_update_agents_attempts": schema.Int64Attribute{
				MarkdownDescription: "The number of Connectors in the Resource Connector Group that failed to update their agents.",
				Computed:            true,
			},
		},
	}
}
func (d *ResourceConnectorGroupDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *ResourceConnectorGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (d *ResourceConnectorGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ResourceConnectorGroup

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
	var res gjson.Result
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
			if config.Name.ValueString() == v.Get("name").String() {
				config.Id = types.StringValue(v.Get("id").String())
				res = v
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

	config.fromBody(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read
