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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ApplicationExceptionListDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationExceptionListDataSource{}
)

func NewApplicationExceptionListDataSource() datasource.DataSource {
	return &ApplicationExceptionListDataSource{}
}

type ApplicationExceptionListDataSource struct {
	client *secureaccess.Client
}

func (d *ApplicationExceptionListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_exception_list"
}

func (d *ApplicationExceptionListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads an exception application list, found by its `id`. A Do Not Decrypt list that has applications uses one, the API returns its `id` as `exceptionapplicationlistId` of that Do Not Decrypt list. Get the ids from the `secureaccess_application_exception_lists` data source. The data source reads all applications of the list.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the exception application list. On the lists observed it is the name of the Do Not Decrypt list that uses it.").String,
				Computed:            true,
			},
			"allowed_count": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of allowed applications in the full list. On every list observed it was the number of all applications.").String,
				Computed:            true,
			},
			"blocked_count": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of blocked applications in the full list. 0 on every list observed.").String,
				Computed:            true,
			},
			"warned_count": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of warned applications in the full list. 0 on every list observed.").String,
				Computed:            true,
			},
			"applications": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Applications of the list, in the order of the API (ascending by `application_id` on the lists observed).").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"application_id": schema.Int64Attribute{
							MarkdownDescription: "Id of the application. It is the id of an application of the type AVC of the `secureaccess_applications` data source.",
							Computed:            true,
						},
						"application_name": schema.StringAttribute{
							MarkdownDescription: "Name of the application.",
							Computed:            true,
						},
						"group_id": schema.Int64Attribute{
							MarkdownDescription: "Id of the group of the application, the id of an application category of the `secureaccess_application_categories` data source.",
							Computed:            true,
						},
						"group_name": schema.StringAttribute{
							MarkdownDescription: "Name of the group of the application, the name of an application category.",
							Computed:            true,
						},
						"has_urls": schema.BoolAttribute{
							MarkdownDescription: "The `hasUrls` of the application as the API returns it. The API draft lists it as a boolean without a description, it was false on every application observed.",
							Computed:            true,
						},
						"action": schema.StringAttribute{
							MarkdownDescription: "Action of the application, `no_decrypt` on the lists observed (displayed as Allow).",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ApplicationExceptionListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (d *ApplicationExceptionListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ApplicationExceptionList

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.String()))

	// The settings of the list are in `data`, the answer has no applications
	urlPath := config.getPath() + "/" + url.QueryEscape(config.Id.ValueString())
	res, err := d.client.Get(ctx, urlPath, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}
	settings := res.Get("data")
	if !settings.IsObject() {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, expected an object in `data`, got: %s", settings.Type))
		return
	}

	// The client reads the applications page by page, with the `page` and `limit` of the answer
	applications, err := d.client.GetAllWithPaginator(ctx, urlPath+"/applications", secureaccess.PageLimitKeyedPagination, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the applications of the object, got error: %s", exceptionApplicationListError(err, applications)))
		return
	}
	// the generated fromBody() would read an item that is not an object as an application without any value
	for i, application := range applications.Array() {
		if !application.IsObject() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the applications of the object, expected an object as application %d, got: %s", i, application.Type))
			return
		}
	}
	body, err := sjson.SetRaw(settings.Raw, "applications", applications.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to build the object, got error: %s", err))
		return
	}

	config.fromBody(ctx, gjson.Parse(body))

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// exceptionApplicationListError adds the answer of the API to the error of the client, which carries only the status code.
func exceptionApplicationListError(err error, res secureaccess.Res) string {
	if answer := res.String(); answer != "" {
		return fmt.Sprintf("%s, answer of the API: %s", err, answer)
	}
	return err.Error()
}
