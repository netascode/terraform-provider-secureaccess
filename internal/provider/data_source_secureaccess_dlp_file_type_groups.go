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
	"sort"

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
	_ datasource.DataSource              = &DLPFileTypeGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &DLPFileTypeGroupsDataSource{}
)

func NewDLPFileTypeGroupsDataSource() datasource.DataSource {
	return &DLPFileTypeGroupsDataSource{}
}

type DLPFileTypeGroupsDataSource struct {
	client *secureaccess.Client
}

func (d *DLPFileTypeGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dlp_file_type_groups"
}

func (d *DLPFileTypeGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource reads DLP File Type Groups.\n If `items` is empty, all DLP File Type Groups are returned, otherwise only the listed ones.\n").String,

		Attributes: map[string]schema.Attribute{
			"items": schema.MapNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Map of DLP File Type Groups. The key of the map is the name of the individual DLP File Type Group.").String,
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Id of the DLP File Type Group.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DLPFileTypeGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (d *DLPFileTypeGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DLPFileTypeGroups

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", "DLP File Type Groups"))

	urlPath := config.getPath()
	// Endpoint is not paginated, a single request returns all items
	res, err := d.client.Get(ctx, urlPath, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}
	// Items are wrapped in {"fileTypes": [...], "fileTypeGroups": [...]}
	resData := res.Get("fileTypeGroups")
	if !resData.IsArray() {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, expected a list, got: %s", resData.Type))
		return
	}

	// Item names requested in config, used to warn about the ones not found
	requestedNames := make([]string, 0, len(config.Items))
	for name := range config.Items {
		requestedNames = append(requestedNames, name)
	}
	sort.Strings(requestedNames)

	// Read all items if user did not provide any specific item names in the config
	if len(config.Items) == 0 {
		config.Items = map[string]DLPFileTypeGroupsItems{}
		resData.ForEach(func(_, v gjson.Result) bool {
			if name := v.Get("name").String(); name != "" {
				config.Items[name] = DLPFileTypeGroupsItems{}
			}
			return true
		})
	}

	config.fromBody(ctx, resData)
	for _, name := range requestedNames {
		if _, ok := config.Items[name]; !ok {
			resp.Diagnostics.AddWarning("Item not found", fmt.Sprintf("%s: %q not found, removed from `items`.", "DLP File Type Groups", name))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", "DLP File Type Groups"))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
