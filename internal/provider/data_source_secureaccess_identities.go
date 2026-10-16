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
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &IdentitiesDataSource{}
	_ datasource.DataSourceWithConfigure = &IdentitiesDataSource{}
)

func NewIdentitiesDataSource() datasource.DataSource {
	return &IdentitiesDataSource{}
}

type IdentitiesDataSource struct {
	client *secureaccess.Client
}

func (d *IdentitiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identities"
}

func (d *IdentitiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource reads Identities, grouped by identity type.\n If `items` is empty, all Identities are returned, otherwise only the listed types (and names).\n If `filter` is provided, only Identities matching the filter are returned.").String,

		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				MarkdownDescription: "Case-insensitive substring used to filter identities by name.",
				Optional:            true,
				Computed:            true,
			},
			"include_deleted": schema.BoolAttribute{
				MarkdownDescription: "Specifies whether to include deleted identities, i.e. marked as deleted or with `(deleted)` in the name.",
				Optional:            true,
				Computed:            true,
			},
			"identity_types": schema.ListAttribute{
				MarkdownDescription: "List of identity types to include.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"items": schema.MapAttribute{
				MarkdownDescription: "Map of Identities grouped by type. The key of the outer map is the identity type (e.g. `site`), the key of the inner map is the name of the individual Identity, each with `id` and `deleted` attributes. If set, only the listed types are returned, and only the listed names of a type whose inner map is not empty.",
				ElementType:         identitiesItemsElemType,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (d *IdentitiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

func (d *IdentitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config Identities

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", "Identities"))

	// Filter and identity types are passed to the API as query parameters
	query := url.Values{}
	if filter := config.Filter.ValueString(); filter != "" {
		query.Set("search", filter)
	}
	if !config.IdentityTypes.IsNull() && !config.IdentityTypes.IsUnknown() {
		var identityTypes []string
		resp.Diagnostics.Append(config.IdentityTypes.ElementsAs(ctx, &identityTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(identityTypes) > 0 {
			query.Set("identitytypes", strings.Join(identityTypes, ","))
		}
	}
	// Without explicit identity types, only the types listed in items are requested
	if query.Get("identitytypes") == "" && !config.Items.IsNull() && !config.Items.IsUnknown() && len(config.Items.Elements()) > 0 {
		itemTypes := make([]string, 0, len(config.Items.Elements()))
		for t := range config.Items.Elements() {
			itemTypes = append(itemTypes, t)
		}
		sort.Strings(itemTypes)
		query.Set("identitytypes", strings.Join(itemTypes, ","))
	}

	// Endpoint uses its own offset/limit pagination without total count, where
	// limit must always be set. A page shorter than the limit is the last one.
	const maxItems = 5000
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for offset := 0; ; offset += maxItems {
		query.Set("limit", strconv.Itoa(maxItems))
		query.Set("offset", strconv.Itoa(offset))
		urlPath := config.getPath() + "?" + query.Encode()
		page, err := d.client.Get(ctx, urlPath, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
			return
		}
		// Items are wrapped in {"data": [...]}
		items := page.Get("data")
		if !items.IsArray() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, expected a list, got: %s", items.Type))
			return
		}
		count := 0
		items.ForEach(func(_, v gjson.Result) bool {
			if !first {
				sb.WriteByte(',')
			}
			sb.WriteString(v.Raw)
			first = false
			count++
			return true
		})
		if count < maxItems {
			break
		}
	}
	sb.WriteByte(']')
	res := gjson.Parse(sb.String())

	resp.Diagnostics.Append(config.fromBody(ctx, res)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", "Identities"))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
