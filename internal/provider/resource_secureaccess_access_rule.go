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

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &AccessRuleResource{}
	_ resource.ResourceWithImportState = &AccessRuleResource{}
)

func NewAccessRuleResource() resource.Resource {
	return &AccessRuleResource{}
}

type AccessRuleResource struct {
	client *secureaccess.Client
}

func (r *AccessRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_rule"
}

func (r *AccessRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages an Access Rule.\n ~> **Note:** This resource is under active development, its schema and behaviour may change in future releases.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the Access Control Rule.").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 50),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Description of the Access Control Rule.").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Action of the Access Control Rule.").AddStringEnumDescription("allow", "block", "warn", "isolate").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("allow", "block", "warn", "isolate"),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Priority of the Access Control Rule. This value is respected on creation, but not tracked afterwards. All rules with lower priorities must be created before this rule.").String,
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Indicates whether the Access Control Rule is enabled.").String,
				Optional:            true,
			},
			"source_identity_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the source identities (users, groups, devices). For Sites and Roaming Computers, use their `origin_id`, not their `id`.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"source_identity_type_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the source identity types.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"source_geolocations": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Country codes (ISO 3166-1 alpha-2, uppercase) of the source geolocations. Combined with identities, it requires `source_operator` `AND` and AD users or groups. Cannot be combined with `destination_geolocations`.").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"source_network_object_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the source Network Objects. Not supported with `source_operator` `AND`.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"source_network_object_group_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the source Network Object Groups. Not supported with `source_operator` `AND`.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"source_private_literals": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Source IP addresses or networks in CIDR notation. Only for `PRIVATE_NETWORK` rules. Not supported with `source_operator` `AND`.").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"source_internet_literals": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Source IP addresses with ports. Only for `PUBLIC_INTERNET` rules. Cannot be combined with destination service objects or service object groups, more than one `destination_literals` item or `destination_literals` with protocol `ICMP`.").String,
				Optional:            true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ips": schema.SetAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("IP addresses or networks in CIDR notation.").String,
							ElementType:         types.StringType,
							Required:            true,
						},
						"ports": schema.SetAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Ports or port ranges, e.g. `443` or `8000-8010`.").String,
							ElementType:         types.StringType,
							Required:            true,
						},
					},
				},
			},
			"source_operator": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Logical operator between the source conditions. With `AND`, a source must match every configured attribute of `source_identity_ids`, `source_identity_type_ids`, `source_geolocations` and `source_internet_literals`, at least two of them must be configured. Source network objects and `source_private_literals` are not supported with `AND`. See the [Managing Access Rules guide](../guides/managing_access_rules) for examples.").AddStringEnumDescription("OR", "AND").AddDefaultValueDescription("OR").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("OR", "AND"),
				},
				Default: stringdefault.StaticString("OR"),
			},
			"destination_network_object_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Network Objects.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_network_object_group_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Network Object Groups.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_service_object_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Service Objects.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_service_object_group_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Service Object Groups.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_private_resource_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Private Resources. Only for `PRIVATE_NETWORK` rules.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_private_resource_group_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Private Resource Groups. Only for `PRIVATE_NETWORK` rules.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_security_group_tag_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Security Group Tags, i.e. identities of type `sgt_id`. Only for `PRIVATE_NETWORK` rules.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_content_category_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Content Categories. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_content_category_list_ids`.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_content_category_list_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Content Category Lists. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_content_category_ids`.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_application_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Applications. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_application_list_ids`.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_application_category_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Application Categories. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_application_list_ids`.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_application_list_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the destination Application Lists. Only for `PUBLIC_INTERNET` rules. Cannot be combined with `destination_application_ids` or `destination_application_category_ids`.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_list_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the Destination Lists. Only for `PUBLIC_INTERNET` rules.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"destination_geolocations": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Country codes (ISO 3166-1 alpha-2, uppercase) of the destination geolocations. Only for `PUBLIC_INTERNET` rules with action `block`. Cannot be combined with other destinations or with `source_geolocations`.").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"destination_literals": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Destination IP addresses with ports and protocol. Not supported with `destination_operator` `AND`.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ips": schema.SetAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("IP addresses or networks in CIDR notation.").String,
							ElementType:         types.StringType,
							Required:            true,
						},
						"ports": schema.SetAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Ports or port ranges, e.g. `443` or `8000-8010`. Use `0-65535` for any port, also with `ICMP`.").String,
							ElementType:         types.StringType,
							Required:            true,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Protocol.").AddStringEnumDescription("TCP", "UDP", "ICMP", "ANY").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("TCP", "UDP", "ICMP", "ANY"),
							},
						},
					},
				},
			},
			"destination_operator": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Logical operator between the destination conditions. With `AND`, the conditions are grouped by type (network objects, service objects, content categories, applications), a destination must match every group and any condition within a group. At least two groups must be configured. Destination Lists, Private Resources, Private Resource Groups, `destination_security_group_tag_ids`, `destination_literals` and `destination_geolocations` are not supported with `AND`. See the [Managing Access Rules guide](../guides/managing_access_rules) for examples.").AddStringEnumDescription("OR", "AND").AddDefaultValueDescription("OR").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("OR", "AND"),
				},
				Default: stringdefault.StaticString("OR"),
			},
			"destination_operator_network_service": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Logical operator between the network object and service object groups, only relevant with `destination_operator` `AND`. With `OR`, a destination must match any network or service object, and at least one other group must be configured (e.g. `destination_content_category_ids`).").AddStringEnumDescription("AND", "OR").AddDefaultValueDescription("AND").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("AND", "OR"),
				},
				Default: stringdefault.StaticString("AND"),
			},
			"log_level": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Logging level for traffic matching the Access Control Rule.").AddStringEnumDescription("LOG_ALL", "LOG_SECURITY", "LOG_NONE").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("LOG_ALL", "LOG_SECURITY", "LOG_NONE"),
				},
			},
			"traffic_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Type of the Access Control Rule, private access (`PRIVATE_NETWORK`) or internet access (`PUBLIC_INTERNET`). The available sources, destinations and actions depend on it.").AddStringEnumDescription("PRIVATE_NETWORK", "PUBLIC_INTERNET").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("PRIVATE_NETWORK", "PUBLIC_INTERNET"),
				},
			},
			"public_security_profile_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Set the public security profile ID.").String,
				Optional:            true,
			},
			"private_security_profile_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Set the private security profile ID.").String,
				Optional:            true,
			},
			"tenant_control_profile_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Set the tenant control profile ID.").String,
				Optional:            true,
			},
			"ips_profile_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Set the IPS posture profile ID.").String,
				Optional:            true,
			},
		},
	}
}

func (r *AccessRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

var _ resource.ResourceWithValidateConfig = &AccessRuleResource{}

// ValidateConfig checks the source and destination conditions, see validateAccessRuleConditions.
func (r *AccessRuleResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	validateAccessRuleConditions(ctx, req.Config, &resp.Diagnostics)
}

func (r *AccessRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AccessRule

	accessRuleMu.Lock()
	defer accessRuleMu.Unlock()

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, AccessRule{})
	body = plan.adjustBody(ctx, body)
	res, err := retryAccessRuleInvalidIds(ctx, func() (secureaccess.Res, error) {
		return r.client.Post(ctx, plan.getPath(), body, reqMods...)
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST/PUT), got error: %s, %s", err, res.String()))
		return
	}
	plan.Id = types.StringValue(res.Get("ruleId").String())
	plan.fromBodyUnknowns(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (r *AccessRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AccessRule

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

func (r *AccessRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AccessRule

	accessRuleMu.Lock()
	defer accessRuleMu.Unlock()

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

	body := plan.toBody(ctx, state)
	body = plan.adjustBody(ctx, body)
	// API requires `rulePriority`, keep the current one if `priority` is not configured
	if !gjson.Get(body, "rulePriority").Exists() {
		res, err := r.client.Get(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString()), reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
			return
		}
		body, _ = sjson.Set(body, "rulePriority", res.Get("rulePriority").Int())
	}
	res, err := retryAccessRuleInvalidIds(ctx, func() (secureaccess.Res, error) {
		return r.client.Put(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString()), body, reqMods...)
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *AccessRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AccessRule

	accessRuleMu.Lock()
	defer accessRuleMu.Unlock()

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

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *AccessRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
