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
	_ datasource.DataSource              = &ApplicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationsDataSource{}
)

func NewApplicationsDataSource() datasource.DataSource {
	return &ApplicationsDataSource{}
}

type ApplicationsDataSource struct {
	client *secureaccess.Client
}

func (d *ApplicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

func (d *ApplicationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource reads Applications.\n If `items` is empty, all Applications are returned, otherwise only the listed ones.\n If `filter` is provided, only Applications matching the filter are returned.").String,

		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				MarkdownDescription: "Case-insensitive substring used to filter applications by name.",
				Optional:            true,
				Computed:            true,
			},
			"items": schema.MapNestedAttribute{
				MarkdownDescription: "Map of Applications. The key of the map is the name of the individual Application.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Id of the Application.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the Application",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							MarkdownDescription: "Category of the Application",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ApplicationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (d *ApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config Applications

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", "Applications"))

	urlPath := config.getPath()
	// Endpoint is not paginated, a single request returns all items
	res, err := d.client.Get(ctx, urlPath, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}
	// Items are wrapped in {"data": {"applications": [...]}}
	res = res.Get("data.applications")
	if !res.IsArray() {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, expected a list, got: %s", res.Type))
		return
	}

	// Read all items if user did not provide any specific item names in the config
	if len(config.Items) == 0 {
		config.Items = map[string]ApplicationsItems{}
		res.ForEach(func(_, v gjson.Result) bool {
			if name := v.Get("label").String(); name != "" {
				config.Items[name] = ApplicationsItems{}
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

	config.fromBody(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", "Applications"))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
