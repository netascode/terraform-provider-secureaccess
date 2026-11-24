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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	_ resource.Resource                = &ZTAProfileUserMappingsResource{}
	_ resource.ResourceWithImportState = &ZTAProfileUserMappingsResource{}
)

func NewZTAProfileUserMappingsResource() resource.Resource {
	return &ZTAProfileUserMappingsResource{}
}

type ZTAProfileUserMappingsResource struct {
	client *secureaccess.Client
}

func (r *ZTAProfileUserMappingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zta_profile_user_mappings"
}

func (r *ZTAProfileUserMappingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages a ZTA Profile User Mappings.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zta_profile_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("ID of the ZTA Profile.").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"revision": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Revision of the configuration.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"users": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("List of users mapped to the ZTA Profile.").String,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the user.").String,
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *ZTAProfileUserMappingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (r *ZTAProfileUserMappingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ZTAProfileUserMappings

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	plan.Id = plan.ZtaProfileId

	ztaProfileMu.RLock()
	defer ztaProfileMu.RUnlock()
	defer ztaProfileLocks.Lock(plan.ZtaProfileId.ValueString())()

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	res, err := r.setMappings(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object, got error: %s", err))
		return
	}
	plan.fromBodyUnknowns(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

func (r *ZTAProfileUserMappingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ZTAProfileUserMappings

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Id.String()))

	res, err := r.getMappings(ctx, state)
	if err != nil && strings.Contains(err.Error(), "StatusCode 404") {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s", err))
		return
	}

	imp, diags := helpers.IsFlagImporting(ctx, req)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// After `terraform import` we switch to a full read.
	if imp {
		state.fromBody(ctx, res)
	} else {
		state.fromBodyPartial(ctx, res)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

func (r *ZTAProfileUserMappingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ZTAProfileUserMappings

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ztaProfileMu.RLock()
	defer ztaProfileMu.RUnlock()
	defer ztaProfileLocks.Lock(plan.ZtaProfileId.ValueString())()
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	if _, err := r.setMappings(ctx, plan); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ZTAProfileUserMappingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ZTAProfileUserMappings

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ztaProfileMu.RLock()
	defer ztaProfileMu.RUnlock()
	defer ztaProfileLocks.Lock(state.ZtaProfileId.ValueString())()
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	// Remove only the mappings managed by this resource
	body := state.toBodyPatch(state.userIds(), nil, state.Revision.ValueInt64())
	if body != "" {
		res, err := r.client.Patch(ctx, state.getPath(), body)
		if err != nil && !strings.Contains(err.Error(), "StatusCode 404") {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (PATCH), got error: %s, %s", err, res.String()))
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

func (r *ZTAProfileUserMappingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zta_profile_id"), req.ID)...)

	helpers.SetFlagImporting(ctx, true, resp.Private, &resp.Diagnostics)
}

// getMappings reads all mappings, and the revision, which is not part of the merged pages, from a single page.
// The result has the shape of a single page response: {"items": [...], "rev": ...}.
func (r *ZTAProfileUserMappingsResource) getMappings(ctx context.Context, data ZTAProfileUserMappings) (gjson.Result, error) {
	// Client handles pagination and returns all items merged into a single list
	items, err := r.client.GetAll(ctx, data.getPath())
	if err != nil {
		return gjson.Result{}, fmt.Errorf("%w, %s", err, items.String())
	}
	res, err := r.client.Get(ctx, data.getPath()+"?limit=1")
	if err != nil {
		return gjson.Result{}, fmt.Errorf("%w, %s", err, res.String())
	}
	body := ""
	body, _ = sjson.SetRaw(body, "items", items.Raw)
	body, _ = sjson.Set(body, "rev", res.Get("rev").Int())
	return gjson.Parse(body), nil
}

// setMappings changes the current mappings into the ones in data. The API rejects adding an already mapped
// and removing a not mapped user, so the operations are built from the current mappings, read with their
// revision right before the change. Returns the resulting mappings.
func (r *ZTAProfileUserMappingsResource) setMappings(ctx context.Context, data ZTAProfileUserMappings) (gjson.Result, error) {
	current, err := r.getMappings(ctx, data)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("failed to retrieve current mappings (GET): %w", err)
	}
	var currentIds []string
	for _, id := range current.Get("items.#.id").Array() {
		currentIds = append(currentIds, id.String())
	}

	body := data.toBodyPatch(currentIds, data.userIds(), current.Get("rev").Int())
	if body == "" {
		return current, nil
	}
	res, err := r.client.Patch(ctx, data.getPath(), body)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("failed to change mappings (PATCH): %w, %s", err, res.String())
	}
	return res.Result, nil
}
