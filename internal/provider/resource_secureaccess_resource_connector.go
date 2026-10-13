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

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	_ resource.Resource = &ResourceConnectorResource{}
)

func NewResourceConnectorResource() resource.Resource {
	return &ResourceConnectorResource{}
}

type ResourceConnectorResource struct {
	client *secureaccess.Client
}

func (r *ResourceConnectorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_connector"
}

func (r *ResourceConnectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages a Resource Connector, which is previously registered with Cisco Secure Access outside of Terraform.\n Implicit import of the resource will be performed based on the `instance_id` attribute.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The instance ID of the Resource Connector.").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Indicates whether the Connector can receive traffic.").String,
				Optional:            true,
			},
			"confirmed": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Indicates whether the Connector exists. Once confirmed, it cannot be changed.").String,
				Optional:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Hostname of the Resource Connector.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The status of the Resource Connector.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_connector_group_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("The ID of the Resource Connector Group to which this connector belongs.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The version of the Resource Connector.").String,
				Computed:            true,
			},
			"delete_on_destroy": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("When `true`, the Resource Connector is removed from Secure Access on destroy. When `false`, it is only removed from the Terraform state.").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *ResourceConnectorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

func (r *ResourceConnectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceConnector

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Implicit import of the object
	obj, err := getResourceConnectorByInstanceId(ctx, r.client, plan.InstanceId.ValueString(), reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	plan.Id = types.StringValue(obj.Get("id").String())

	// Align the existing connector with the configuration
	current := ResourceConnector{Id: plan.Id}
	current.fromBody(ctx, obj)
	res, err := r.patch(ctx, plan, current, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object, got error: %s, %s", err, res.String()))
		return
	}

	plan.fromBodyUnknowns(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (r *ResourceConnectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceConnector

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

func (r *ResourceConnectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceConnector

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

	res, err := r.patch(ctx, plan, state, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object, got error: %s, %s", err, res.String()))
		return
	}
	plan.fromBodyUnknowns(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// patch sends JSON Patch "replace" operations for enabled/confirmed values which differ between plan and current,
// where current is the prior state or the object as retrieved from the backend.
// Returns the full object, either from the PATCH response or, if no PATCH was needed, retrieved with GET.
func (r *ResourceConnectorResource) patch(ctx context.Context, plan, current ResourceConnector, reqMods ...func(*secureaccess.Req)) (gjson.Result, error) {
	urlPath := plan.getPath() + "/" + url.QueryEscape(plan.Id.ValueString())

	body := plan.toBodyPatch(ctx, current)
	if body == "" {
		res, err := r.client.Get(ctx, urlPath, reqMods...)
		if err != nil {
			return res, fmt.Errorf("failed to retrieve object (GET): %w", err)
		}
		return res, nil
	}

	res, err := r.client.Patch(ctx, urlPath, body, reqMods...)
	if err != nil {
		return res, fmt.Errorf("failed to patch object (PATCH): %w", err)
	}
	return res, nil
}

// getResourceConnectorByInstanceId retrieves the single Resource Connector matching the instance ID.
// Returns an error if the lookup fails, or if no or multiple connectors match.
func getResourceConnectorByInstanceId(ctx context.Context, client *secureaccess.Client, instanceId string, reqMods ...func(*secureaccess.Req)) (gjson.Result, error) {
	filters, _ := sjson.Set("", "InstanceId", instanceId)
	query := url.Values{}
	query.Set("filters", filters)

	res, err := client.Get(ctx, ResourceConnector{}.getPath()+"?"+query.Encode(), reqMods...)
	if err != nil {
		return res, fmt.Errorf("failed to retrieve objects (GET), got error: %s, %s", err, res.String())
	}

	resArray := res.Get("data").Array()
	if len(resArray) == 0 {
		return gjson.Result{}, fmt.Errorf("Resource Connector with instance_id %s not found", instanceId)
	}
	if len(resArray) > 1 {
		return gjson.Result{}, fmt.Errorf("multiple Resource Connectors with instance_id %s found", instanceId)
	}
	return resArray[0], nil
}

func (r *ResourceConnectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceConnector

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	// The connector is adopted on Create, by default it is only removed from the state
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
