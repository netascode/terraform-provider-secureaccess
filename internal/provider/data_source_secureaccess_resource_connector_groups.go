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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ResourceConnectorGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &ResourceConnectorGroupsDataSource{}
)

func NewResourceConnectorGroupsDataSource() datasource.DataSource {
	return &ResourceConnectorGroupsDataSource{}
}

type ResourceConnectorGroupsDataSource struct {
	client *secureaccess.Client
}

func (d *ResourceConnectorGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_connector_groups"
}

func (d *ResourceConnectorGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource reads Resource Connector Groups.\n If `items` is empty, all Resource Connector Groups are returned, otherwise only the listed ones.\n If `filter` is provided, only Resource Connector Groups matching the filter are returned.").String,

		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				MarkdownDescription: "Case-insensitive substring used to filter resource connector groups by name.",
				Optional:            true,
				Computed:            true,
			},
			"items": schema.MapNestedAttribute{
				MarkdownDescription: "Map of Resource Connector Groups. The key of the map is the name of the individual Resource Connector Group.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Id of the Resource Connector Group.",
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
				},
			},
		},
	}
}

func (d *ResourceConnectorGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (d *ResourceConnectorGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ResourceConnectorGroups

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", "Resource Connector Groups"))

	urlPath := config.getPath()
	// Client handles pagination and returns all items merged into a single list
	res, err := d.client.GetAll(ctx, urlPath, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}
	if !res.IsArray() {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, expected a list, got: %s", res.Type))
		return
	}

	// Read all items if user did not provide any specific item names in the config
	if len(config.Items) == 0 {
		config.Items = map[string]ResourceConnectorGroupsItems{}
		res.ForEach(func(_, v gjson.Result) bool {
			if name := v.Get("name").String(); name != "" {
				config.Items[name] = ResourceConnectorGroupsItems{}
			}
			return true
		})
	}

	// Keep only items whose name contains the filter (case-insensitive)
	if filter := strings.ToLower(config.Filter.ValueString()); filter != "" {
		for name := range config.Items {
			if !strings.Contains(strings.ToLower(name), filter) {
				tflog.Debug(ctx, fmt.Sprintf("item does not match filter, removing: name=%v", name))
				delete(config.Items, name)
			}
		}
	}

	config.fromBody(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", "Resource Connector Groups"))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read
