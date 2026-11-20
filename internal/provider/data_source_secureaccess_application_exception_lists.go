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
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	_ datasource.DataSource              = &ApplicationExceptionListsDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationExceptionListsDataSource{}
)

func NewApplicationExceptionListsDataSource() datasource.DataSource {
	return &ApplicationExceptionListsDataSource{}
}

type ApplicationExceptionListsDataSource struct {
	client *secureaccess.Client
}

func (d *ApplicationExceptionListsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_exception_lists"
}

func (d *ApplicationExceptionListsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the ids of the exception application lists that the Do Not Decrypt lists of the organization use. A Do Not Decrypt list that has applications uses one, the API returns its id as `exceptionapplicationlistId` of that Do Not Decrypt list. Use the `secureaccess_application_exception_list` data source to read a list.").String,

		Attributes: map[string]schema.Attribute{
			"ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Ids of the exception application lists that Do Not Decrypt lists use.").String,
				ElementType:         types.Int64Type,
				Computed:            true,
			},
		},
	}
}

func (d *ApplicationExceptionListsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (d *ApplicationExceptionListsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ApplicationExceptionLists

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", "Application Exception Lists"))

	// The ids are read page by page, all of them or an error
	ids, err := readExceptionApplicationListPages(ctx, d.client, config.getPath(), reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, got error: %s", err))
		return
	}
	// The ids are numbers: the generated fromBody() would read anything else as 0
	var invalid gjson.Result
	var duplicate string
	seen := map[string]bool{}
	gjson.Parse(ids).ForEach(func(_, id gjson.Result) bool {
		if id.Type != gjson.Number {
			invalid = id
			return false
		}
		if seen[id.Raw] {
			duplicate = id.Raw
			return false
		}
		seen[id.Raw] = true
		return true
	})
	if invalid.Exists() {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, expected a number as id, got: %s", invalid.Type))
		return
	}
	if duplicate != "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, the id %s is listed more than once", duplicate))
		return
	}
	body, err := sjson.SetRaw("{}", "ids", ids)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to build the object, got error: %s", err))
		return
	}

	config.fromBody(ctx, gjson.Parse(body))

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", "Application Exception Lists"))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
