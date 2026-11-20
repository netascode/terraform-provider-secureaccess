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
	_ datasource.DataSource              = &DLPAIGuardrailsRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &DLPAIGuardrailsRuleDataSource{}
)

func NewDLPAIGuardrailsRuleDataSource() datasource.DataSource {
	return &DLPAIGuardrailsRuleDataSource{}
}

type DLPAIGuardrailsRuleDataSource struct {
	client *secureaccess.Client
}

func (d *DLPAIGuardrailsRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dlp_ai_guardrails_rule"
}

func (d *DLPAIGuardrailsRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This data source reads the DLP AI Guardrails Rule.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the DLP AI Guardrails Rule.").String,
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Description of the rule.").String,
				Computed:            true,
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Severity of the rule.").String,
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the rule is enabled.").String,
				Computed:            true,
			},
			"data_classification_scannable_contexts": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Select where to search for the selected data classifications.").String,
				ElementType:         types.StringType,
				Computed:            true,
			},
			"data_classifications": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Classifications to scan using OR boolean logic.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "ID of the classification.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the classification.",
							Computed:            true,
						},
					},
				},
			},
			"file_labels": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Up to 10 labels that this rule will search for when inspecting document properties. Applicable only when `unclassified_files_scan` is not set to `true`.").String,
				ElementType:         types.StringType,
				Computed:            true,
			},
			"file_size_configuration_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Configuration type for file size control.").String,
				Computed:            true,
			},
			"file_size_from": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Minimum file size for file control.").String,
				Computed:            true,
			},
			"file_size_from_unit": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Unit of the minimum file size for file control.").String,
				Computed:            true,
			},
			"file_size_to": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Maximum file size for file control.").String,
				Computed:            true,
			},
			"file_size_to_unit": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Unit of the maximum file size for file control.").String,
				Computed:            true,
			},
			"file_types_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("File types included for file control.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "ID of the file type.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the file type.",
							Computed:            true,
						},
					},
				},
			},
			"file_type_groups_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("File type groups included for file control.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "ID of the file type group.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the file type group.",
							Computed:            true,
						},
					},
				},
			},
			"file_types_excluded": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("File types excluded from file control.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "ID of the file type.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the file type.",
							Computed:            true,
						},
					},
				},
			},
			"file_type_groups_excluded": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("File type groups excluded from file control.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "ID of the file type group.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the file type group.",
							Computed:            true,
						},
					},
				},
			},
			"unclassified_files_scan": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Scan and log files with no labels. Applicable only when `file_labels` is unset.").String,
				Computed:            true,
			},
			"identities_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Identities included for scanning. Use `type_id` to specify whole types of identities (i.a. Network Tunnel Groups). Add `id` to specify individual identities within `type_id`.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "ID of the identity.",
							Computed:            true,
						},
						"type_id": schema.Int64Attribute{
							MarkdownDescription: "Type ID of the identity.",
							Computed:            true,
						},
					},
				},
			},
			"identities_excluded": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Identities excluded from scanning. Use `type_id` to specify whole types of identities (i.a. Network Tunnel Groups). Add `id` to specify individual identities within `type_id`.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "ID of the identity.",
							Computed:            true,
						},
						"type_id": schema.Int64Attribute{
							MarkdownDescription: "Type ID of the identity.",
							Computed:            true,
						},
					},
				},
			},
			"applications_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Applications included for scanning.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "ID of the application.",
							Computed:            true,
						},
						"traffic_direction": schema.StringAttribute{
							MarkdownDescription: "Traffic direction of the application.",
							Computed:            true,
						},
					},
				},
			},
			"application_categories_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Application categories included for scanning.").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "ID of the application category.",
							Computed:            true,
						},
						"traffic_direction": schema.StringAttribute{
							MarkdownDescription: "Traffic direction of the application.",
							Computed:            true,
						},
					},
				},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Action associated with the rule.").String,
				Computed:            true,
			},
			"override_block": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Allow override with business justification. Requires `notification_popup` to be set to `true`.").String,
				Computed:            true,
			},
			"notification_popup": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Sends endpoint notifications to the user's machine.").String,
				Computed:            true,
			},
			"notification_template_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("ID of the notification template for custom notifications.").String,
				Computed:            true,
			},
			"email_notify_actor": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Sends email notifications to the actor.").String,
				Computed:            true,
			},
			"email_notify_actors_manager": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Sends email notifications to the actors' manager.").String,
				Computed:            true,
			},
			"email_notify_custom_recipients": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Sends email notifications to custom email addresses.").String,
				ElementType:         types.StringType,
				Computed:            true,
			},
			"email_notification_template_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("UUID of the notification template.").String,
				Computed:            true,
			},
		},
	}
}
func (d *DLPAIGuardrailsRuleDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *DLPAIGuardrailsRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (d *DLPAIGuardrailsRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DLPAIGuardrailsRule

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.String()))

	// Object found by a query attribute, it is not retrieved again by ID
	var res secureaccess.Res
	if config.Id.IsNull() && !config.Name.IsNull() {
		// Client handles pagination and returns all items merged into a single list
		items, err := d.client.GetAll(ctx, config.getPath(), reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, got error: %s", err))
			return
		}
		if !items.IsArray() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, expected a list, got: %s", items.Type))
			return
		}
		items.ForEach(func(k, v gjson.Result) bool {
			if config.Name.ValueString() == v.Get("name").String() {
				config.Id = types.StringValue(v.Get("id").String())
				res.Result = v
				tflog.Debug(ctx, fmt.Sprintf("%s: Found object with name '%v', id: %v", config.Id.ValueString(), config.Name.ValueString(), config.Id.ValueString()))
				return false
			}
			return true
		})

		if config.Id.IsNull() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to find object with name: %v", config.Name.ValueString()))
			return
		}
	}

	urlPath := config.getPath() + "/" + url.QueryEscape(config.Id.ValueString())
	var err error
	res, err = d.client.Get(ctx, urlPath, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}

	config.fromBody(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
