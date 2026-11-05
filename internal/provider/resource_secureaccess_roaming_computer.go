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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource = &RoamingComputerResource{}
)

func NewRoamingComputerResource() resource.Resource {
	return &RoamingComputerResource{}
}

type RoamingComputerResource struct {
	client *secureaccess.Client
}

func (r *RoamingComputerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roaming_computer"
}

func (r *RoamingComputerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages a Roaming Computer, which is previously registered with Cisco Secure Access outside of Terraform.\n Implicit import of the resource will be performed based on the `name` attribute.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The name of the Roaming Computer. Used to find the Roaming Computer on create, a later change renames it.").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 50),
				},
			},
			"device_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The hex ID of the Roaming Computer.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"anyconnect_device_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The ID of the device that has the Cisco Secure Client deployed with the Internet Security module.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"origin_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("The origin ID of the Roaming Computer.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"applied_bundle": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("The ID of the policy applied to the Roaming Computer.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The type of the Roaming Computer.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"os_version": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The OS version of the Roaming Computer.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"os_version_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The OS version name of the Roaming Computer.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"has_ip_blocking": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Specifies whether the Roaming Computer has IP blocking.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The status of the Roaming Computer with DNS-layer security.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secure_web_gateway_status": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The status of the Roaming Computer with Internet security (Secure Web Gateway).").String,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secure_client_version": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The version of the Roaming Computer client.").String,
				Computed:            true,
			},
			"delete_on_destroy": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("When `true`, the Roaming Computer is removed from Secure Access on destroy. When `false`, it is only removed from the Terraform state.").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *RoamingComputerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (r *RoamingComputerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoamingComputer

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Implicit import of the object
	res, err := getRoamingComputerByName(ctx, r.client, plan.Name.ValueString(), reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	plan.Id = types.StringValue(res.Get("deviceId").String())

	plan.fromBodyUnknowns(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (r *RoamingComputerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoamingComputer

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

func (r *RoamingComputerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state RoamingComputer

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

	// Only the name can be updated, other changes (e.g. delete_on_destroy) are state only
	urlPath := plan.getPath() + "/" + url.QueryEscape(plan.Id.ValueString())
	var res secureaccess.Res
	var err error
	if body := plan.toBodyPut(ctx, state); body != "" {
		res, err = r.client.Put(ctx, urlPath, body, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s, %s", err, res.String()))
			return
		}
	} else {
		res, err = r.client.Get(ctx, urlPath, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
			return
		}
	}
	plan.fromBodyUnknowns(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// getRoamingComputerByName retrieves the single Roaming Computer matching the name.
// Returns an error if the lookup fails, or if no or multiple roaming computers match.
func getRoamingComputerByName(ctx context.Context, client *secureaccess.Client, name string, reqMods ...func(*secureaccess.Req)) (secureaccess.Res, error) {
	// The name query parameter filters the list server side, a single page is sufficient
	query := url.Values{}
	query.Set("name", name)
	query.Set("limit", "100")

	res, err := client.Get(ctx, RoamingComputer{}.getPath()+"?"+query.Encode(), reqMods...)
	if err != nil {
		return res, fmt.Errorf("failed to retrieve objects (GET), got error: %s, %s", err, res.String())
	}

	// The filter might not be an exact match, compare the names
	var matches []gjson.Result
	for _, v := range res.Array() {
		if v.Get("name").String() == name {
			matches = append(matches, v)
		}
	}
	if len(matches) == 0 {
		return secureaccess.Res{}, fmt.Errorf("Roaming Computer with name %s not found", name)
	}
	if len(matches) > 1 {
		return secureaccess.Res{}, fmt.Errorf("multiple Roaming Computers with name %s found", name)
	}
	return secureaccess.Res{Result: matches[0], StatusCode: res.StatusCode}, nil
}

func (r *RoamingComputerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoamingComputer

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	// The roaming computer is adopted on Create, by default it is only removed from the state
	if state.DeleteOnDestroy.ValueBool() {
		res, err := r.client.Delete(ctx, state.getPath()+"/"+url.QueryEscape(state.Id.ValueString()), reqMods...)
		if err != nil && !strings.Contains(err.Error(), "StatusCode 404") {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
			return
		}
	} else {
		tflog.Debug(ctx, fmt.Sprintf("%s: delete_on_destroy is false, removing from state only", state.Id.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// Section below is generated&owned by "gen/generator.go". //template:begin import
// End of section. //template:end import
