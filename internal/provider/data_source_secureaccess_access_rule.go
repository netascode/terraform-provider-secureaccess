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
	_ datasource.DataSource              = &AccessRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &AccessRuleDataSource{}
)

func NewAccessRuleDataSource() datasource.DataSource {
	return &AccessRuleDataSource{}
}

type AccessRuleDataSource struct {
	client *secureaccess.Client
}

func (d *AccessRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_rule"
}

func (d *AccessRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the Access Rule.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Access Control Rule.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the Access Control Rule.",
				Computed:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Action of the Access Control Rule.",
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority of the Access Control Rule. This value is respected on creation, but not tracked afterwards. All rules with lower priorities must be created before this rule.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the Access Control Rule is enabled.",
				Computed:            true,
			},
			"source_identity_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the source identities (users, groups, devices). For Sites and Roaming Computers, use their `origin_id`, not their `id`.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"source_identity_type_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the source identity types.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"source_geolocations": schema.SetAttribute{
				MarkdownDescription: "Country codes (ISO 3166-1 alpha-2, uppercase) of the source geolocations. Combined with identities, it requires `source_operator` `AND` and AD users or groups. Cannot be combined with `destination_geolocations`.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"source_network_object_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the source Network Objects. Not supported with `source_operator` `AND`.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"source_network_object_group_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the source Network Object Groups. Not supported with `source_operator` `AND`.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"source_private_literals": schema.SetAttribute{
				MarkdownDescription: "Source IP addresses or networks in CIDR notation. Only for `PRIVATE_NETWORK` rules. Not supported with `source_operator` `AND`.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"source_internet_literals": schema.ListNestedAttribute{
				MarkdownDescription: "Source IP addresses with ports. Only for `PUBLIC_INTERNET` rules. Cannot be combined with destination service objects or service object groups, more than one `destination_literals` item or `destination_literals` with protocol `ICMP`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ips": schema.SetAttribute{
							MarkdownDescription: "IP addresses or networks in CIDR notation.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"ports": schema.SetAttribute{
							MarkdownDescription: "Ports or port ranges, e.g. `443` or `8000-8010`.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
			"source_operator": schema.StringAttribute{
				MarkdownDescription: "Logical operator between the source conditions. With `AND`, a source must match every configured attribute of `source_identity_ids`, `source_identity_type_ids`, `source_geolocations` and `source_internet_literals`, at least two of them must be configured. Source network objects and `source_private_literals` are not supported with `AND`. See the [Managing Access Rules guide](../guides/managing_access_rules) for examples.",
				Computed:            true,
			},
			"destination_network_object_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Network Objects.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_network_object_group_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Network Object Groups.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_service_object_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Service Objects.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_service_object_group_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Service Object Groups.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_private_resource_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Private Resources. Only for `PRIVATE_NETWORK` rules.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_private_resource_group_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Private Resource Groups. Only for `PRIVATE_NETWORK` rules.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_security_group_tag_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Security Group Tags, i.e. identities of type `sgt_id`. Only for `PRIVATE_NETWORK` rules.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_content_category_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Content Categories. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_content_category_list_ids`.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_content_category_list_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Content Category Lists. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_content_category_ids`.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_application_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Applications. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_application_list_ids`.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_application_category_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Application Categories. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_application_list_ids`.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_application_list_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the destination Application Lists. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_application_ids` or `destination_application_category_ids`.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_list_ids": schema.SetAttribute{
				MarkdownDescription: "IDs of the Destination Lists. Only for `PUBLIC_INTERNET` rules.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"destination_geolocations": schema.SetAttribute{
				MarkdownDescription: "Country codes (ISO 3166-1 alpha-2, uppercase) of the destination geolocations. Only for `PUBLIC_INTERNET` rules with action `block`. Cannot be combined with other destinations or with `source_geolocations`.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"destination_literals": schema.SetNestedAttribute{
				MarkdownDescription: "Destination IP addresses with ports and protocol. Not supported with `destination_operator` `AND`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ips": schema.SetAttribute{
							MarkdownDescription: "IP addresses or networks in CIDR notation.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"ports": schema.SetAttribute{
							MarkdownDescription: "Ports or port ranges, e.g. `443` or `8000-8010`. Use `0-65535` for any port, also with `ICMP`.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "Protocol.",
							Computed:            true,
						},
					},
				},
			},
			"destination_operator": schema.StringAttribute{
				MarkdownDescription: "Logical operator between the destination conditions. With `AND`, the conditions are grouped by type (network objects, service objects, content categories, applications), a destination must match every group and any condition within a group. At least two groups must be configured. Destination Lists, Private Resources, Private Resource Groups, `destination_security_group_tag_ids`, `destination_literals` and `destination_geolocations` are not supported with `AND`. See the [Managing Access Rules guide](../guides/managing_access_rules) for examples.",
				Computed:            true,
			},
			"destination_operator_network_service": schema.StringAttribute{
				MarkdownDescription: "Logical operator between the network object and service object groups, only relevant with `destination_operator` `AND`. With `OR`, a destination must match any network or service object, and at least one other group must be configured (e.g. `destination_content_category_ids`).",
				Computed:            true,
			},
			"log_level": schema.StringAttribute{
				MarkdownDescription: "Logging level for traffic matching the Access Control Rule.",
				Computed:            true,
			},
			"traffic_type": schema.StringAttribute{
				MarkdownDescription: "Type of the Access Control Rule, private access (`PRIVATE_NETWORK`) or internet access (`PUBLIC_INTERNET`). The available sources, destinations and actions depend on it.",
				Computed:            true,
			},
			"public_security_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Set the public security profile ID.",
				Computed:            true,
			},
			"private_security_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Set the private security profile ID.",
				Computed:            true,
			},
			"tenant_control_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Set the tenant control profile ID.",
				Computed:            true,
			},
			"ips_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Set the IPS posture profile ID.",
				Computed:            true,
			},
		},
	}
}

func (d *AccessRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (d *AccessRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AccessRule

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
