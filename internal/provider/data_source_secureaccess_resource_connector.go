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
	_ datasource.DataSource              = &ResourceConnectorDataSource{}
	_ datasource.DataSourceWithConfigure = &ResourceConnectorDataSource{}
)

func NewResourceConnectorDataSource() datasource.DataSource {
	return &ResourceConnectorDataSource{}
}

type ResourceConnectorDataSource struct {
	client *secureaccess.Client
}

func (d *ResourceConnectorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_connector"
}

func (d *ResourceConnectorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the Resource Connector.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Optional:            true,
				Computed:            true,
			},
			"instance_id": schema.StringAttribute{
				MarkdownDescription: "The instance ID of the Resource Connector.",
				Optional:            true,
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the Connector can receive traffic.",
				Computed:            true,
			},
			"confirmed": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the Connector exists. Once confirmed, it cannot be changed.",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname of the Resource Connector.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the Resource Connector.",
				Computed:            true,
			},
			"resource_connector_group_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Resource Connector Group to which this connector belongs.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The version of the Resource Connector.",
				Computed:            true,
			},
			"delete_on_destroy": schema.BoolAttribute{
				MarkdownDescription: "When `true`, the Resource Connector is removed from Secure Access on destroy. When `false`, it is only removed from the Terraform state.",
				Computed:            true,
			},
		},
	}
}
func (d *ResourceConnectorDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("instance_id"),
		),
	}
}

func (d *ResourceConnectorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (d *ResourceConnectorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ResourceConnector

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.String()))

	var res gjson.Result
	if config.Id.IsNull() && !config.InstanceId.IsNull() {
		obj, err := getResourceConnectorByInstanceId(ctx, d.client, config.InstanceId.ValueString(), reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
		config.Id = types.StringValue(obj.Get("id").String())
		tflog.Debug(ctx, fmt.Sprintf("%s: Found object with instance_id '%v'", config.Id.ValueString(), config.InstanceId.ValueString()))
		res = obj
	} else {
		urlPath := config.getPath() + "/" + url.QueryEscape(config.Id.ValueString())
		obj, err := d.client.Get(ctx, urlPath, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
			return
		}
		res = obj
	}

	config.fromBody(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
