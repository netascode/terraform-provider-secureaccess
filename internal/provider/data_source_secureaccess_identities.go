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

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
		MarkdownDescription: helpers.NewAttributeDescription("This is Early Access and may be subject to changes. This resource reads Identities of a single identity type.\n If `items` is empty, all Identities of `identity_type` are returned, otherwise only the listed ones.\n If `search` is provided, only Identities matching the search string are returned.\n").String,

		Attributes: map[string]schema.Attribute{
			"identity_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Identity type to read.").AddStringEnumDescription("ad_connector", "android", "anyconnect", "app_connector_agent", "bdr", "branch", "chromebook_user", "connected_devices", "device", "directory_computer", "directory_domain", "directory_group", "directory_user", "domain_controller", "edlp_client", "gsuite_domain", "gsuite_group", "gsuite_user", "internal_network", "meraki_auto_vpn", "mobile_device", "network", "network_device", "network_dns_forwarder", "network_tunnel_hub", "on_network", "policy", "roaming", "saml_group", "saml_user", "secure_web_appliance", "sgt_id", "site", "tunnel_device", "vpn", "vpn_id", "ztna_client").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ad_connector", "android", "anyconnect", "app_connector_agent", "bdr", "branch", "chromebook_user", "connected_devices", "device", "directory_computer", "directory_domain", "directory_group", "directory_user", "domain_controller", "edlp_client", "gsuite_domain", "gsuite_group", "gsuite_user", "internal_network", "meraki_auto_vpn", "mobile_device", "network", "network_device", "network_dns_forwarder", "network_tunnel_hub", "on_network", "policy", "roaming", "saml_group", "saml_user", "secure_web_appliance", "sgt_id", "site", "tunnel_device", "vpn", "vpn_id", "ztna_client"),
				},
			},
			"search": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Search string, only Identities where it appears in the Identity data are returned.").String,
				Optional:            true,
				Computed:            true,
			},
			"include_deleted": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Specifies whether to include deleted identities, i.e. marked as deleted or with `(deleted)` in the name.").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
			},
			"items": schema.MapNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Map of Identities. The key of the map is the name of the individual Identity.").String,
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Id of the Identity.",
							Computed:            true,
						},
						"type_id": schema.Int64Attribute{
							MarkdownDescription: "Type ID of the Identity.",
							Computed:            true,
						},
						"deleted": schema.BoolAttribute{
							MarkdownDescription: "Deletion status of the Identity.",
							Computed:            true,
						},
					},
				},
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

// End of section. //template:end model

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

	// Default applies only when not configured, fromBody() always sets it
	includeDeleted := config.IncludeDeleted
	if includeDeleted.IsNull() || includeDeleted.IsUnknown() {
		includeDeleted = types.BoolValue(false)
	}

	// Search is passed to the API as query parameter, identity type is part of getPath()
	query := url.Values{}
	if search := config.Search.ValueString(); search != "" {
		query.Set("search", search)
	}

	// Endpoint uses its own offset/limit pagination without total count, where
	// limit must always be set. A page shorter than the limit is the last one.
	// Deleted identities are skipped here, as the API has no parameter for them.
	// Some deleted identities are only marked by `(deleted)` in the name.
	const maxItems = 5000
	seen := map[string]int64{}
	var sb strings.Builder
	sb.WriteByte('[')
	for offset := 0; ; offset += maxItems {
		query.Set("limit", strconv.Itoa(maxItems))
		query.Set("offset", strconv.Itoa(offset))
		urlPath := config.getPath() + "&" + query.Encode()
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
			count++
			name := v.Get("label").String()
			if name == "" {
				return true
			}
			if !includeDeleted.ValueBool() && (v.Get("deleted").Bool() || strings.Contains(name, "(deleted)")) {
				return true
			}
			// Names are matched case-insensitively, the first identity is kept
			if id, ok := seen[strings.ToLower(name)]; ok {
				resp.Diagnostics.AddWarning("Duplicate identity name",
					fmt.Sprintf("Multiple identities of type %q are named %q, keeping id %d and skipping id %d.", config.IdentityType.ValueString(), name, id, v.Get("id").Int()))
				return true
			}
			seen[strings.ToLower(name)] = v.Get("id").Int()
			if len(seen) > 1 {
				sb.WriteByte(',')
			}
			sb.WriteString(v.Raw)
			return true
		})
		if count < maxItems {
			break
		}
	}
	sb.WriteByte(']')
	res := gjson.Parse(sb.String())

	// Item names requested in config, used to warn about the ones not found
	requestedNames := make([]string, 0, len(config.Items))
	for name := range config.Items {
		requestedNames = append(requestedNames, name)
	}
	sort.Strings(requestedNames)

	// Read all items if user did not provide any specific item names in the config
	if len(config.Items) == 0 {
		config.Items = map[string]IdentitiesItems{}
		res.ForEach(func(_, v gjson.Result) bool {
			if name := v.Get("label").String(); name != "" {
				config.Items[name] = IdentitiesItems{}
			}
			return true
		})
	}

	config.fromBody(ctx, res)
	for _, name := range requestedNames {
		if _, ok := config.Items[name]; !ok {
			resp.Diagnostics.AddWarning("Item not found", fmt.Sprintf("%s: %q not found, removed from `items`.", "Identities", name))
		}
	}
	config.IncludeDeleted = includeDeleted

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", "Identities"))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
