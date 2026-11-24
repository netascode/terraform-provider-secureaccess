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
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &DLPRealTimeRuleResource{}
	_ resource.ResourceWithImportState = &DLPRealTimeRuleResource{}
)

func NewDLPRealTimeRuleResource() resource.Resource {
	return &DLPRealTimeRuleResource{}
}

type DLPRealTimeRuleResource struct {
	client *secureaccess.Client
}

func (r *DLPRealTimeRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dlp_real_time_rule"
}

func (r *DLPRealTimeRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages a DLP Real Time Rule.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the DLP Real Time Rule.").String,
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Description of the rule.").String,
				Optional:            true,
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Severity of the rule.").AddStringEnumDescription("CRITICAL", "ALERT", "WARNING", "INFO").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("CRITICAL", "ALERT", "WARNING", "INFO"),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the rule is enabled.").String,
				Required:            true,
			},
			"data_classification_scannable_contexts": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Select where to search for the selected data classifications.").AddStringEnumDescription("CONTENT", "FILENAME").String,
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("CONTENT", "FILENAME"),
					),
				},
			},
			"data_classifications": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Classifications to scan using OR boolean logic.").String,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the classification.").String,
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of the classification.").String,
							Optional:            true,
						},
					},
				},
			},
			"file_labels": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Up to 10 labels that this rule will search for when inspecting document properties. Applicable only when `unclassified_files_scan` is not set to `true`.").String,
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(10),
				},
			},
			"file_size_configuration_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Configuration type for file size control. With `RANGE`, all `file_size_*` attributes are required. With `ALL`, none of them can be set.").AddStringEnumDescription("ALL", "RANGE").AddDefaultValueDescription("ALL").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ALL", "RANGE"),
				},
				Default: stringdefault.StaticString("ALL"),
			},
			"file_size_from": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Minimum file size for file control. Required with `file_size_configuration_type` `RANGE`.").String,
				Optional:            true,
			},
			"file_size_from_unit": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Unit of the minimum file size for file control. Required with `file_size_configuration_type` `RANGE`.").AddStringEnumDescription("KB", "MB").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("KB", "MB"),
				},
			},
			"file_size_to": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Maximum file size for file control. Required with `file_size_configuration_type` `RANGE`.").String,
				Optional:            true,
			},
			"file_size_to_unit": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Unit of the maximum file size for file control. Required with `file_size_configuration_type` `RANGE`.").AddStringEnumDescription("KB", "MB").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("KB", "MB"),
				},
			},
			"file_types_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("File types included for file control.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the file type.").String,
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of the file type.").String,
							Optional:            true,
						},
					},
				},
			},
			"file_type_groups_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("File type groups included for file control.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the file type group.").String,
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of the file type group.").String,
							Optional:            true,
						},
					},
				},
			},
			"file_types_excluded": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("File types excluded from file control.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the file type.").String,
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of the file type.").String,
							Optional:            true,
						},
					},
				},
			},
			"file_type_groups_excluded": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("File type groups excluded from file control.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the file type group.").String,
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of the file type group.").String,
							Optional:            true,
						},
					},
				},
			},
			"unclassified_files_scan": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Scan and log files with no labels. Applicable only when `file_labels` is unset.").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"identities_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Identities included for scanning. Use `type_id` to specify whole types of identities (i.a. Network Tunnel Groups). Add `id` to specify individual identities within `type_id`.").String,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the identity.").String,
							Optional:            true,
						},
						"type_id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Type ID of the identity.").String,
							Required:            true,
						},
					},
				},
			},
			"identities_excluded": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Identities excluded from scanning. Use `type_id` to specify whole types of identities (i.a. Network Tunnel Groups). Add `id` to specify individual identities within `type_id`.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the identity.").String,
							Optional:            true,
						},
						"type_id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Type ID of the identity.").String,
							Required:            true,
						},
					},
				},
			},
			"all_destinations_scope": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Scan traffic for all destinations specified by the scope. Use `NONE` to specify destinations manually.").AddStringEnumDescription("NONE", "FILES_AND_FORMS", "FILES_AND_VETTED_APPS_FORMS").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("NONE", "FILES_AND_FORMS", "FILES_AND_VETTED_APPS_FORMS"),
				},
			},
			"destination_list_included_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Destination list IDs included for scanning.").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"destination_list_excluded_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Destination list IDs excluded from scanning.").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"applications_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Applications included for scanning.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the application.").String,
							Required:            true,
						},
						"traffic_direction": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Traffic direction of the application.").AddStringEnumDescription("REQUEST", "RESPONSE", "BOTH").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("REQUEST", "RESPONSE", "BOTH"),
							},
						},
					},
				},
			},
			"applications_excluded": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Applications excluded from scanning.").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"application_categories_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Application categories included for scanning.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the application category.").String,
							Required:            true,
						},
						"traffic_direction": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Traffic direction of the application.").AddStringEnumDescription("REQUEST", "RESPONSE", "BOTH").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("REQUEST", "RESPONSE", "BOTH"),
							},
						},
					},
				},
			},
			"application_categories_excluded": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Application categories excluded from scanning.").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"private_resources_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Private resources included for scanning.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the private resource.").String,
							Required:            true,
						},
						"traffic_direction": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Traffic direction of the private resource.").AddStringEnumDescription("REQUEST", "RESPONSE", "BOTH").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("REQUEST", "RESPONSE", "BOTH"),
							},
						},
					},
				},
			},
			"private_resources_excluded": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Private resources excluded from scanning.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the private resource.").String,
							Required:            true,
						},
						"traffic_direction": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Traffic direction of the private resource.").AddStringEnumDescription("REQUEST", "RESPONSE", "BOTH").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("REQUEST", "RESPONSE", "BOTH"),
							},
						},
					},
				},
			},
			"private_resource_groups_included": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Private resource groups included for scanning.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the private resource group.").String,
							Required:            true,
						},
						"traffic_direction": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Traffic direction of the private resource.").AddStringEnumDescription("REQUEST", "RESPONSE", "BOTH").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("REQUEST", "RESPONSE", "BOTH"),
							},
						},
					},
				},
			},
			"private_resource_groups_excluded": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Private resource groups excluded from scanning.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the private resource group.").String,
							Required:            true,
						},
						"traffic_direction": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Traffic direction of the private resource.").AddStringEnumDescription("REQUEST", "RESPONSE", "BOTH").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("REQUEST", "RESPONSE", "BOTH"),
							},
						},
					},
				},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Action associated with the rule.").AddStringEnumDescription("BLOCK", "MONITOR").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("BLOCK", "MONITOR"),
				},
			},
			"override_block": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Allow override with business justification. Requires `notification_popup` to be set to `true`.").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"notification_popup": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Sends endpoint notifications to the user's machine.").String,
				Optional:            true,
			},
			"notification_template_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("ID of the notification template for custom notifications.").String,
				Optional:            true,
			},
			"email_notify_actor": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Sends email notifications to the actor.").String,
				Optional:            true,
			},
			"email_notify_actors_manager": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Sends email notifications to the actors' manager.").String,
				Optional:            true,
			},
			"email_notify_custom_recipients": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Sends email notifications to custom email addresses.").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"email_notification_template_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("UUID of the notification template.").String,
				Optional:            true,
			},
		},
	}
}

func (r *DLPRealTimeRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

var _ resource.ResourceWithValidateConfig = &DLPRealTimeRuleResource{}

// ValidateConfig requires all file size attributes with `RANGE` and none of them with `ALL` (the default).
func (r *DLPRealTimeRuleResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config DLPRealTimeRule
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() || config.FileSizeConfigurationType.IsUnknown() {
		return
	}
	fileSizes := []struct {
		name   string
		isNull bool
	}{
		{"file_size_from", config.FileSizeFrom.IsNull()},
		{"file_size_from_unit", config.FileSizeFromUnit.IsNull()},
		{"file_size_to", config.FileSizeTo.IsNull()},
		{"file_size_to_unit", config.FileSizeToUnit.IsNull()},
	}
	isRange := config.FileSizeConfigurationType.ValueString() == "RANGE"
	for _, fileSize := range fileSizes {
		name := fileSize.name
		if isRange && fileSize.isNull {
			resp.Diagnostics.AddAttributeError(path.Root(name), "Missing Attribute", fmt.Sprintf("`%s` is required when `file_size_configuration_type` is `RANGE`.", name))
		} else if !isRange && !fileSize.isNull {
			resp.Diagnostics.AddAttributeError(path.Root(name), "Invalid Attribute Combination", fmt.Sprintf("`%s` cannot be set unless `file_size_configuration_type` is `RANGE`.", name))
		}
	}
}

// Section below is generated&owned by "gen/generator.go". //template:begin create

func (r *DLPRealTimeRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DLPRealTimeRule

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, DLPRealTimeRule{})
	body = plan.adjustBody(ctx, body)
	res, err := r.client.Post(ctx, plan.getPath(), body, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST/PUT), got error: %s, %s", err, res.String()))
		return
	}
	plan.Id = types.StringValue(res.Get("id").String())
	plan.fromBodyUnknowns(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (r *DLPRealTimeRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DLPRealTimeRule

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Id.String()))

	urlPath := state.getPath() + "/" + url.QueryEscape(state.Id.ValueString())
	res, err := r.client.Get(ctx, urlPath, reqMods...)
	if err != nil && strings.Contains(err.Error(), "StatusCode 404") {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}

	imp, diags := helpers.IsFlagImporting(ctx, req)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// After `terraform import` we switch to a full read.
	if imp {
		state.fromBody(ctx, res.Result)
	} else {
		state.fromBodyPartial(ctx, res.Result)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// End of section. //template:end read

// Section below is generated&owned by "gen/generator.go". //template:begin update

func (r *DLPRealTimeRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DLPRealTimeRule

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Read state
	diags = req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	body := plan.toBodyPatchMerge(ctx, state)
	if body != "" {
		res, err := r.client.Patch(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString()), body, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PATCH), got error: %s, %s", err, res.String()))
			return
		}
	} else {
		tflog.Debug(ctx, fmt.Sprintf("%s: No changes to PATCH", plan.Id.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end update

// Section below is generated&owned by "gen/generator.go". //template:begin delete

func (r *DLPRealTimeRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DLPRealTimeRule

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))
	res, err := r.client.Delete(ctx, state.getPath()+"/"+url.QueryEscape(state.Id.ValueString()), reqMods...)
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *DLPRealTimeRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse import ID
	var inputPattern = regexp.MustCompile(`^(?P<id>[^\s,]+?)$`)
	match := inputPattern.FindStringSubmatch(req.ID)
	if match == nil {
		errMsg := "Failed to parse import parameters.\nPlease provide import string in the following format: <id>\n" + fmt.Sprintf("Got: %q", req.ID)
		resp.Diagnostics.AddError("Import error", errMsg)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), match[inputPattern.SubexpIndex("id")])...)

	helpers.SetFlagImporting(ctx, true, resp.Private, &resp.Diagnostics)
}

// End of section. //template:end import
