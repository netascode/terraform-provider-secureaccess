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
	"encoding/json"
	"fmt"
	"slices"
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
	_ datasource.DataSource              = &IdentityTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &IdentityTypesDataSource{}
)

func NewIdentityTypesDataSource() datasource.DataSource {
	return &IdentityTypesDataSource{}
}

type IdentityTypesDataSource struct {
	client *secureaccess.Client
}

func (d *IdentityTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_types"
}

func (d *IdentityTypesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource reads Identity Types.\n Available Identity Types include, but are not limited to: `Registered Networks`, `AD Groups`, `AD Devices`, `AD Users`, `Apple Devices`, `Internal Networks`, `Sites`, `Network Devices`, `ChromeOS Devices`, `Network Tunnel Groups`, `Google Workspace Users`, `Google Workspace Organizational Units`, `Catalyst SD-WAN Service VPN IDs`, `Security Group Tags`, `iOS/Android Devices` and `macOS/Windows Devices`.\n").String,

		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Case-insensitive substring used to filter identity types by name.").String,
				Optional:            true,
				Computed:            true,
			},
			"items": schema.MapNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Map of Identity Types. The key of the map is the name of the individual Identity Type.").String,
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Id of the Identity Type.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *IdentityTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (d *IdentityTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config IdentityTypes

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", "Identity Types"))

	body := identityTypesStaticBody()
	res := secureaccess.Res{
		Result: gjson.ParseBytes(body),
	}

	// Read all items if user did not provide any specific item names in the config
	if len(config.Items) == 0 {
		config.Items = map[string]IdentityTypesItems{}
		res.ForEach(func(_, v gjson.Result) bool {
			if name := v.Get("label").String(); name != "" {
				config.Items[name] = IdentityTypesItems{}
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

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", "Identity Types"))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// identityTypesStaticBody returns identityTypesStatic and identityTypesCombined as a JSON list, mimicking the API response
func identityTypesStaticBody() []byte {
	items := slices.Clip(identityTypesStatic) // append copies
	for _, c := range identityTypesCombined {
		items = append(items, identityTypeStatic{Id: c.Id, Label: c.Label})
	}
	body, _ := json.Marshal(items)
	return body
}
