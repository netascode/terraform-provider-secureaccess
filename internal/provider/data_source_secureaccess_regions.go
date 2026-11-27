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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &RegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &RegionsDataSource{}
)

func NewRegionsDataSource() datasource.DataSource {
	return &RegionsDataSource{}
}

type RegionsDataSource struct {
	client *secureaccess.Client
}

func (d *RegionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regions"
}

func (d *RegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the regions that Secure Access offers for Network Tunnel Groups. It requires an API key with the `deployments.regions:read` scope.").String,

		Attributes: map[string]schema.Attribute{
			"regions": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The regions, in the order in which the API returns them.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Display name of the region. To select a region for a Network Tunnel Group use `region`, not `name`.",
							Computed:            true,
						},
						"region": schema.StringAttribute{
							MarkdownDescription: "Identifier of the region. This is the value to use for the `region` of a Network Tunnel Group.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Country and location of the region.",
							Computed:            true,
						},
						"continent": schema.StringAttribute{
							MarkdownDescription: "Continent of the region.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RegionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (d *RegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config Regions

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", "Regions"))

	urlPath := config.getPath()
	res, err := d.client.Get(ctx, urlPath, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}

	config.fromBody(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", "Regions"))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read
