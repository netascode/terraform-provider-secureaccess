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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/planmodifiers"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworkTunnelGroupResource{}
	_ resource.ResourceWithImportState = &NetworkTunnelGroupResource{}
)

func NewNetworkTunnelGroupResource() resource.Resource {
	return &NetworkTunnelGroupResource{}
}

type NetworkTunnelGroupResource struct {
	client *secureaccess.Client
}

func (r *NetworkTunnelGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_tunnel_group"
}

func (r *NetworkTunnelGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages a Network Tunnel Group.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the Network Tunnel Group. It cannot exceed 50 characters, and can only contain alphanumeric characters, spaces, and hyphens.").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 50),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9 -]+$`), ""),
				},
			},
			"device_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The type of device that establishes the network tunnel.").AddStringEnumDescription("ASA", "AWS S2S VPN", "AZURE S2S VPN", "FTD", "ISR", "Meraki MX", "Viptela cEdge", "Viptela vEdge", "other").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ASA", "AWS S2S VPN", "AZURE S2S VPN", "FTD", "ISR", "Meraki MX", "Viptela cEdge", "Viptela vEdge", "other"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The name of the region that the system uses to obtain the primary and secondary data centers for the Hubs.").String,
				Required:            true,
			},
			"tunnel_id_format_email": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Prefix for the authentication ID in email format (<tunnel_id_format_email>@<org><hub>.sse.cisco.com). The string should not have any special characters besides the period(.), underscore(_), and dash(-) characters.").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 100),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9._-]+$`), ""),
				},
			},
			"tunnel_id_format_ip_addresses": schema.ListAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("List of IP addresses for the authentication ID in IP address format. Please include exactly two IP addresses.").String,
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(2),
					listvalidator.SizeAtMost(2),
				},
			},
			"preshared_key": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Passphrase for the network tunnel. It must include at least one upper case letter, one lower case letter, one number, and cannot include any special characters.").String,
				Required:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(16, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`[a-z]`), ""),
					stringvalidator.RegexMatches(regexp.MustCompile(`[A-Z]`), ""),
					stringvalidator.RegexMatches(regexp.MustCompile(`[0-9]`), ""),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9]+$`), ""),
				},
			},
			"routing_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Type of routing for the network tunnel.").AddStringEnumDescription("static").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("static"),
				},
			},
			"routing_network_cidrs": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"hubs": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Remote connection endpoints for connecting network tunnel group").String,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					planmodifiers.ConditionalUseStateForUnknownList("region"),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the hub.").String,
							Computed:            true,
						},
						"is_primary": schema.BoolAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Indicates if the hub is the primary hub.").String,
							Computed:            true,
						},
						"tunnel_id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Tunnel ID").String,
							Computed:            true,
						},
						"datacenter_name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of the datacenter.").String,
							Computed:            true,
						},
						"datacenter_ip": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("IP address of the datacenter.").String,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (r *NetworkTunnelGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create

func (r *NetworkTunnelGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkTunnelGroup

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, NetworkTunnelGroup{})
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

func (r *NetworkTunnelGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkTunnelGroup

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

func (r *NetworkTunnelGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkTunnelGroup

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

	body := plan.toBodyPatch(ctx, state)
	if body != "" {
		res, err := r.client.Patch(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString()), body, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PATCH), got error: %s, %s", err, res.String()))
			return
		}

		// Hubs are recomputed on region change, retrieve them if not included in the PATCH response
		if plan.Hubs.IsUnknown() && !res.Get("hubs").Exists() {
			res, err = r.client.Get(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString()), reqMods...)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
				return
			}
		}
		plan.fromBodyUnknowns(ctx, res.Result)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Section below is generated&owned by "gen/generator.go". //template:begin delete

func (r *NetworkTunnelGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkTunnelGroup

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
func (r *NetworkTunnelGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
