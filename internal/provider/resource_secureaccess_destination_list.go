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
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

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
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &DestinationListResource{}
	_ resource.ResourceWithImportState = &DestinationListResource{}
)

func NewDestinationListResource() resource.Resource {
	return &DestinationListResource{}
}

type DestinationListResource struct {
	client *secureaccess.Client
}

func (r *DestinationListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination_list"
}

func (r *DestinationListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages a Destination List.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the Destination List.").String,
				Required:            true,
			},
			"destinations": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("List of destinations in the Destination List").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Unique identifier for the destination.").String,
							Computed:            true,
						},
						"destination": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("The destination value (domain name, URL, or IPv4 address).").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 253),
							},
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Comment for the destination.").String,
							Optional:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Type of the destination.").AddStringEnumDescription("domain", "url", "ipv4").String,
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("domain", "url", "ipv4"),
							},
						},
					},
				},
			},
		},
	}
}

func (r *DestinationListResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (r *DestinationListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DestinationList

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, DestinationList{})
	res, err := r.client.Post(ctx, plan.getPath(), body, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST/PUT), got error: %s, %s", err, res.String()))
		return
	}

	res = res.Get("data")
	plan.Id = types.StringValue(res.Get("id").String())

	// get destinations // TODO error handling
	resDest, _ := r.client.Get(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString())+"/destinations", reqMods...)
	resDest = resDest.Get("data")
	s := res.String()
	s, _ = sjson.SetRaw(s, "destinations", resDest.String())
	res = gjson.Parse(s)

	plan.fromBodyUnknowns(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

func (r *DestinationListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DestinationList

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

	res = res.Get("data")

	// get destinations
	resDest, _ := r.client.Get(ctx, state.getPath()+"/"+url.QueryEscape(state.Id.ValueString())+"/destinations", reqMods...)
	resDest = resDest.Get("data")

	s := res.String()
	s, _ = sjson.SetRaw(s, "destinations", resDest.String())
	res = gjson.Parse(s)

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

func (r *DestinationListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DestinationList

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

	if state.Name != plan.Name {
		// this could probably be simplified as name is the only thing that could change here
		body := plan.toBody(ctx, state)
		body, _ = sjson.Delete(body, "destinations")
		res, err := r.client.Patch(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString()), body, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PATCH), got error: %s, %s", err, res.String()))
			return
		}
	}

	// Update destinations if they have changed
	var toDelete = []int64{}
	var toCreate = []DestinationListDestinations{}

	var stateDestinations = map[int64]DestinationListDestinations{}
	for _, dest := range state.Destinations {
		stateDestinations[dest.Id.ValueInt64()] = dest
	}

	var planDestinations = map[int64]DestinationListDestinations{}
	for _, dest := range plan.Destinations {
		if dest.Id.IsUnknown() {
			toCreate = append(toCreate, dest)
			continue
		}
		planDestinations[dest.Id.ValueInt64()] = dest
	}

	for id := range stateDestinations {
		if _, exists := planDestinations[id]; !exists {
			toDelete = append(toDelete, id)
		} else {
			// Technically only comment could change, defensive check for destination change as well
			if stateDestinations[id].Destination.ValueString() != planDestinations[id].Destination.ValueString() ||
				stateDestinations[id].Comment.ValueString() != planDestinations[id].Comment.ValueString() {
				toDelete = append(toDelete, id)
				toCreate = append(toCreate, planDestinations[id])
			}
		}
	}

	for id := range planDestinations {
		if _, exists := stateDestinations[id]; !exists {
			toCreate = append(toCreate, planDestinations[id])
		}
	}

	// Delete destinations
	if len(toDelete) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("%s: Deleting destinations: %v", plan.Id.ValueString(), toDelete))
		json, _ := json.Marshal(toDelete)
		body := string(json)
		res, err := r.client.DeleteWithBody(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString())+"/destinations/remove", body)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete destinations, got error: %s, %s", err, res.String()))
			return
		}
	}

	// Create destinations
	if len(toCreate) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("%s: Creating destinations: %v", plan.Id.ValueString(), toCreate))
		tmp := DestinationList{
			Destinations: toCreate,
		}
		body := tmp.toBody(ctx, DestinationList{})
		body = gjson.Get(body, "destinations").String()
		res, err := r.client.Post(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString())+"/destinations", body)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create destinations, got error: %s, %s", err, res.String()))
			return
		}

		// Retrieve ID of the newly created destinations
		resDest, _ := r.client.Get(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString())+"/destinations", reqMods...)
		resDest = resDest.Get("data")
		s := res.String()
		s, _ = sjson.SetRaw(s, "destinations", resDest.String())
		res = gjson.Parse(s)

		plan.fromBodyUnknowns(ctx, res)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Section below is generated&owned by "gen/generator.go". //template:begin delete

func (r *DestinationListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DestinationList

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
func (r *DestinationListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse import ID
	var inputPattern = regexp.MustCompile(`^(?P<id>[^\s,]+?)$`)
	match := inputPattern.FindStringSubmatch(req.ID)
	if match == nil {
		errMsg := "Failed to parse import parameters.\nPlease provide import string in the following format: ,<id>\n" + fmt.Sprintf("Got: %q", req.ID)
		resp.Diagnostics.AddError("Import error", errMsg)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), match[inputPattern.SubexpIndex("id")])...)

	helpers.SetFlagImporting(ctx, true, resp.Private, &resp.Diagnostics)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin createSubresources

// End of section. //template:end createSubresources

// Section below is generated&owned by "gen/generator.go". //template:begin deleteSubresources

// End of section. //template:end deleteSubresources

// Section below is generated&owned by "gen/generator.go". //template:begin updateSubresources

// End of section. //template:end updateSubresources
