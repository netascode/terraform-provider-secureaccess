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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	_ resource.Resource                = &DoNotDecryptListResource{}
	_ resource.ResourceWithImportState = &DoNotDecryptListResource{}
)

func NewDoNotDecryptListResource() resource.Resource {
	return &DoNotDecryptListResource{}
}

type DoNotDecryptListResource struct {
	client *secureaccess.Client
}

func (r *DoNotDecryptListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_do_not_decrypt_list"
}

func (r *DoNotDecryptListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages a Do Not Decrypt List. The Secure Access API does not return the destinations and applications of a list, therefore changes made to them outside of Terraform are not detected and they are not populated by `terraform import`.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the Do Not Decrypt List.").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1024),
				},
			},
			"destinations": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Destinations that are excluded from decryption, domain names (in lowercase), IPv4 addresses or CIDR blocks. URLs are not supported. Large lists are split into multiple requests automatically.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"destination": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("The destination value (domain name in lowercase, IPv4 address or CIDR block).").String,
							Required:            true,
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Comment for the destination.").String,
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 255),
							},
						},
					},
				},
			},
			"applications": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the applications that are excluded from decryption.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"decrypt_exception_categories": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IDs of the content categories that are excluded from decryption.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"is_default": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Specifies whether this is the default Do Not Decrypt List.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"total_domains": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of domains in the Do Not Decrypt List.").String,
				Computed:            true,
			},
			"total_ips": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of IP addresses in the Do Not Decrypt List.").String,
				Computed:            true,
			},
			"total_applications": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of applications in the Do Not Decrypt List.").String,
				Computed:            true,
			},
		},
	}
}

func (r *DoNotDecryptListResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (r *DoNotDecryptListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DoNotDecryptList

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object, the API accepts at most 500 destinations here, the rest is added with updates afterwards
	createPlan := plan
	var remaining []DoNotDecryptListDestinations
	if len(plan.Destinations) > doNotDecryptListMaxItemsPerRequest {
		createPlan.Destinations = plan.Destinations[:doNotDecryptListMaxItemsPerRequest]
		remaining = plan.Destinations[doNotDecryptListMaxItemsPerRequest:]
	}
	body := createPlan.toBody(ctx, DoNotDecryptList{})
	res, err := r.client.Post(ctx, plan.getPath(), body, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST), got error: %s, %s", err, res.String()))
		return
	}
	if res.Get("id").String() == "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST), the response has no id: %s", res.String()))
		return
	}
	plan.Id = types.StringValue(res.Get("id").String())

	for i, step := range doNotDecryptListChangeSteps(plan.Name.ValueString(), true, remaining, nil, nil) {
		res, err = r.client.Put(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString()), step.body, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to add remaining destinations (PUT %d), got error: %s, %s", i+1, err, res.String()))
			// Persist the id so the object is tracked (and tainted) rather than orphaned
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.Id)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), plan.Name)...)
			return
		}
	}
	plan.fromBodyUnknowns(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (r *DoNotDecryptListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DoNotDecryptList

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
		state.fromBody(ctx, res)
	} else {
		state.fromBodyPartial(ctx, res)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// End of section. //template:end read

func (r *DoNotDecryptListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DoNotDecryptList

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

	// The API takes lists of items to add and to remove instead of the full lists, which can require multiple requests.
	// If one of them fails after others were accepted, the state is set to what the API holds by then. The API does not
	// return destinations and applications, so Terraform could not find out about the difference otherwise.
	steps := plan.toUpdateSteps(ctx, state)
	interim := state
	objectPath := plan.getPath() + "/" + url.QueryEscape(plan.Id.ValueString())
	var res gjson.Result
	for i, step := range steps {
		var err error
		if step.viaList {
			res, err = r.removeThroughDestinationList(ctx, objectPath, step.removeDestinations, reqMods...)
			if err != nil {
				err = fmt.Errorf("failed to remove IP addresses through the destination list: %w", err)
			}
		} else {
			res, err = r.client.Put(ctx, objectPath, step.body, reqMods...)
			if err != nil {
				err = fmt.Errorf("%w, %s", err, res.String())
			}
		}
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (request %d of %d), got error: %s", i+1, len(steps), err))
			if i > 0 {
				resp.Diagnostics.Append(resp.State.Set(ctx, &interim)...)
			}
			return
		}
		interim.applyStep(ctx, plan.Name, step)
		interim.fromBodyUnknowns(ctx, res)
	}
	plan.fromBodyUnknowns(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Section below is generated&owned by "gen/generator.go". //template:begin delete

func (r *DoNotDecryptListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DoNotDecryptList

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
func (r *DoNotDecryptListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

var _ resource.ResourceWithValidateConfig = &DoNotDecryptListResource{}

// removeThroughDestinationList removes IP addresses and CIDR blocks. The update request of the API cannot remove them
// ("destination not found"), but the destinations of the object are also available as a destination list
// (exceptiondomainlistId), from which they can be removed by id. A destination that is not in the list anymore is
// skipped, as it is gone already. Returns the object as it is afterwards.
func (r *DoNotDecryptListResource) removeThroughDestinationList(ctx context.Context, objectPath string, items []DoNotDecryptListDestinations, reqMods ...func(*secureaccess.Req)) (gjson.Result, error) {
	res, err := r.client.Get(ctx, objectPath, reqMods...)
	if err != nil {
		return res, fmt.Errorf("failed to retrieve object (GET): %w, %s", err, res.String())
	}
	listId := res.Get("exceptiondomainlistId").String()
	if listId == "" || listId == "0" {
		return res, fmt.Errorf("the object has no destination list (exceptiondomainlistId), got: %s", res.String())
	}

	// Find the ids of the destinations, the client reads all pages
	wanted := make(map[string]struct{}, len(items))
	for _, item := range items {
		wanted[item.Destination.ValueString()] = struct{}{}
	}
	destinations, err := r.client.GetAll(ctx, doNotDecryptListDestinationListPath(listId), reqMods...)
	if err != nil {
		return res, fmt.Errorf("failed to retrieve destinations (GET): %w, %s", err, destinations.String())
	}
	found := make(map[string]int64, len(items))
	for _, item := range destinations.Array() {
		destination := item.Get("destination").String()
		if _, ok := wanted[destination]; ok {
			found[destination] = item.Get("id").Int()
		}
	}
	var ids []int64
	for _, item := range items {
		if id, ok := found[item.Destination.ValueString()]; ok {
			ids = append(ids, id)
		} else {
			tflog.Warn(ctx, fmt.Sprintf("%s: %s is not in the destination list anymore, skipping its removal", objectPath, item.Destination.ValueString()))
		}
	}

	for i := 0; i < len(ids); i += doNotDecryptListListRemovalChunk {
		end := min(i+doNotDecryptListListRemovalChunk, len(ids))
		tflog.Debug(ctx, fmt.Sprintf("%s: Removing destinations %d-%d of %d from destination list %s", objectPath, i, end, len(ids), listId))
		body, _ := json.Marshal(ids[i:end])
		resDel, err := r.client.DeleteWithBody(ctx, doNotDecryptListDestinationListPath(listId)+"/remove", string(body), reqMods...)
		if err != nil {
			return res, fmt.Errorf("failed to remove destinations %d-%d (DELETE): %w, %s", i, end, err, resDel.String())
		}
	}
	return r.client.Get(ctx, objectPath, reqMods...)
}

// ValidateConfig rejects destinations that the API cannot store as they are (URLs, domain names with uppercase
// letters, which the API stores in lowercase), and a destination that is listed more than once. Destinations are
// removed by their value without the comment, so two entries with the same value and different comments could not be
// managed independently.
func (r *DoNotDecryptListResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var destinations types.Set
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("destinations"), &destinations)...)
	if resp.Diagnostics.HasError() || destinations.IsNull() || destinations.IsUnknown() {
		return
	}
	seen := map[string]struct{}{}
	for _, element := range destinations.Elements() {
		object, ok := element.(types.Object)
		if !ok || object.IsNull() || object.IsUnknown() {
			continue
		}
		destination, ok := object.Attributes()["destination"].(types.String)
		if !ok || destination.IsNull() || destination.IsUnknown() {
			continue
		}
		if problem := doNotDecryptListDestinationProblem(destination.ValueString()); problem != "" {
			resp.Diagnostics.AddAttributeError(path.Root("destinations"), "Invalid destination",
				fmt.Sprintf("The destination %q cannot be used. %s", destination.ValueString(), problem))
			return
		}
		key := strings.ToLower(destination.ValueString())
		if _, exists := seen[key]; exists {
			resp.Diagnostics.AddAttributeError(path.Root("destinations"), "Duplicate destination",
				fmt.Sprintf("The destination %q is listed more than once. A destination can only be listed once, with a single comment.", destination.ValueString()))
			return
		}
		seen[key] = struct{}{}
	}
}
