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

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource = &RoamingComputerWebSecurityOverrideResource{}
)

func NewRoamingComputerWebSecurityOverrideResource() resource.Resource {
	return &RoamingComputerWebSecurityOverrideResource{}
}

type RoamingComputerWebSecurityOverrideResource struct {
	client *secureaccess.Client
}

func (r *RoamingComputerWebSecurityOverrideResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roaming_computer_web_security_override"
}

func (r *RoamingComputerWebSecurityOverrideResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages Secure Web Gateway override settings of Roaming Computers.\n Roaming Computers are referenced by their origin ID. Only the listed Roaming Computers are managed, on destroy their override setting is removed.").String,

		Attributes: map[string]schema.Attribute{
			"enabled_origin_ids": schema.ListAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The list of Roaming Computer origin IDs with Secure Web Gateway always enabled.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
			"disabled_origin_ids": schema.ListAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The list of Roaming Computer origin IDs with Secure Web Gateway always disabled.").String,
				ElementType:         types.Int64Type,
				Optional:            true,
			},
		},
	}
}

func (r *RoamingComputerWebSecurityOverrideResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

var (
	_ resource.ResourceWithConfigValidators = &RoamingComputerWebSecurityOverrideResource{}
	_ resource.ResourceWithValidateConfig   = &RoamingComputerWebSecurityOverrideResource{}
)

func (r *RoamingComputerWebSecurityOverrideResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("enabled_origin_ids"),
			path.MatchRoot("disabled_origin_ids"),
		),
	}
}

// ValidateConfig ensures an origin ID is listed only once across both lists.
func (r *RoamingComputerWebSecurityOverrideResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config RoamingComputerWebSecurityOverride

	diags := req.Config.Get(ctx, &config)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Lists might contain unknown values (e.g. references to resources not yet created)
	seen := map[int64]string{}
	for _, list := range []struct {
		name  string
		value types.List
	}{
		{"enabled_origin_ids", config.EnabledOriginIds},
		{"disabled_origin_ids", config.DisabledOriginIds},
	} {
		for _, v := range list.value.Elements() {
			id, ok := v.(types.Int64)
			if !ok || id.IsUnknown() || id.IsNull() {
				continue
			}
			if other, found := seen[id.ValueInt64()]; found {
				detail := fmt.Sprintf("Origin ID %d is listed more than once in %s.", id.ValueInt64(), list.name)
				if other != list.name {
					detail = fmt.Sprintf("Origin ID %d cannot be part of both %s and %s.", id.ValueInt64(), other, list.name)
				}
				resp.Diagnostics.AddAttributeError(path.Root(list.name), "Invalid Attribute Combination", detail)
				continue
			}
			seen[id.ValueInt64()] = list.name
		}
	}
}

func (r *RoamingComputerWebSecurityOverrideResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoamingComputerWebSecurityOverride

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, "Roaming Computer Web Security Override: Beginning Create")

	if err := r.set(ctx, plan, roamingComputerWebSecurityOverrideEnabled, plan.enabledIds(ctx), reqMods...); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if err := r.set(ctx, plan, roamingComputerWebSecurityOverrideDisabled, plan.disabledIds(ctx), reqMods...); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Debug(ctx, "Roaming Computer Web Security Override: Create finished successfully")

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *RoamingComputerWebSecurityOverrideResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoamingComputerWebSecurityOverride

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, "Roaming Computer Web Security Override: Beginning Read")

	// Only the managed origin IDs are read, the API requires them to be listed explicitly
	var items []gjson.Result
	for _, batch := range roamingComputerWebSecurityOverrideBatches(state.allIds(ctx)) {
		res, err := r.client.Post(ctx, state.getPath()+"/list", state.toBodyOriginIds(batch), reqMods...)
		if err != nil && strings.Contains(err.Error(), "StatusCode 404") {
			continue
		} else if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects (POST), got error: %s, %s", err, res.String()))
			return
		}
		items = append(items, res.Array()...)
	}

	state.fromBody(ctx, items)

	tflog.Debug(ctx, "Roaming Computer Web Security Override: Read finished successfully")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *RoamingComputerWebSecurityOverrideResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state RoamingComputerWebSecurityOverride

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

	tflog.Debug(ctx, "Roaming Computer Web Security Override: Beginning Update")

	// Remove origin IDs no longer managed
	removed := roamingComputerWebSecurityOverrideDifference(state.allIds(ctx), plan.allIds(ctx))
	if err := r.remove(ctx, plan, removed, reqMods...); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	// Set only origin IDs, which are new or moved between the lists
	enabled := roamingComputerWebSecurityOverrideDifference(plan.enabledIds(ctx), state.enabledIds(ctx))
	if err := r.set(ctx, plan, roamingComputerWebSecurityOverrideEnabled, enabled, reqMods...); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	disabled := roamingComputerWebSecurityOverrideDifference(plan.disabledIds(ctx), state.disabledIds(ctx))
	if err := r.set(ctx, plan, roamingComputerWebSecurityOverrideDisabled, disabled, reqMods...); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Debug(ctx, "Roaming Computer Web Security Override: Update finished successfully")

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *RoamingComputerWebSecurityOverrideResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoamingComputerWebSecurityOverride

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, "Roaming Computer Web Security Override: Beginning Delete")

	if err := r.remove(ctx, state, state.allIds(ctx), reqMods...); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Debug(ctx, "Roaming Computer Web Security Override: Delete finished successfully")

	resp.State.RemoveResource(ctx)
}

// set applies the override value to the origin IDs. The API reports failures per origin ID in a successful response.
func (r *RoamingComputerWebSecurityOverrideResource) set(ctx context.Context, data RoamingComputerWebSecurityOverride, value string, ids []int64, reqMods ...func(*secureaccess.Req)) error {
	for _, batch := range roamingComputerWebSecurityOverrideBatches(ids) {
		res, err := r.client.Post(ctx, data.getPath()+"/set", data.toBodySet(value, batch), reqMods...)
		if err != nil {
			return fmt.Errorf("failed to configure objects (POST), got error: %s, %s", err, res.String())
		}
		if res.Get("failCount").Int() > 0 {
			var failures []string
			for _, item := range res.Get("items").Array() {
				if code := item.Get("code").Int(); code < 200 || code > 299 {
					failures = append(failures, fmt.Sprintf("origin ID %d: %d %s", item.Get("originId").Int(), code, item.Get("message").String()))
				}
			}
			return fmt.Errorf("failed to configure %d of %d objects (POST): %s", res.Get("failCount").Int(), res.Get("totalCount").Int(), strings.Join(failures, "; "))
		}
	}
	return nil
}

// remove deletes the override setting of the origin IDs.
func (r *RoamingComputerWebSecurityOverrideResource) remove(ctx context.Context, data RoamingComputerWebSecurityOverride, ids []int64, reqMods ...func(*secureaccess.Req)) error {
	for _, batch := range roamingComputerWebSecurityOverrideBatches(ids) {
		res, err := r.client.Post(ctx, data.getPath()+"/remove", data.toBodyOriginIds(batch), reqMods...)
		if err != nil && !strings.Contains(err.Error(), "StatusCode 404") {
			return fmt.Errorf("failed to delete objects (POST), got error: %s, %s", err, res.String())
		}
	}
	return nil
}
