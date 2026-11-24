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
	"maps"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &ZTAProfileResource{}
	_ resource.ResourceWithImportState = &ZTAProfileResource{}
)

func NewZTAProfileResource() resource.Resource {
	return &ZTAProfileResource{}
}

type ZTAProfileResource struct {
	client *secureaccess.Client
}

func (r *ZTAProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zta_profile"
}

func (r *ZTAProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource manages ZTA Profile. The order of created profiles is determined by order they are created by Terraform and `priority`, hence they may not be created in the order `priority` suggests. Subsequent update will fix ordering.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the ZTA Profile. Must be unique within the organization. Must be alphanumeric or space. The name `Default ZTA Profile` manages the existing default profile, which is updated instead of created, and reset instead of deleted.").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9 ]+$`), ""),
				},
			},
			"revision": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Revision of the configuration.").String,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Priority of the ZTA Profile. Profiles should have contiguous priority values starting from 1. The default profile, named `Default ZTA Profile`, always has priority `0`.").String,
				Required:            true,
			},
			"operating_system_macos": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the ZTA Profile is enabled for macOS clients.").AddDefaultValueDescription("true").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"operating_system_windows": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the ZTA Profile is enabled for Windows clients.").AddDefaultValueDescription("true").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"operating_system_linux": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the ZTA Profile is enabled for Linux clients.").AddDefaultValueDescription("true").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"operating_system_ios": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the ZTA Profile is enabled for iOS clients.").AddDefaultValueDescription("true").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"operating_system_android": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the ZTA Profile is enabled for generic Android clients.").AddDefaultValueDescription("true").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"operating_system_android_knox": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the ZTA Profile is enabled for Samsung Knox Android clients.").AddDefaultValueDescription("true").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"all_users": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the ZTA Profile is enabled for all users of the organization.").String,
				Optional:            true,
			},
			"all_groups": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the ZTA Profile is enabled for all groups of the organization.").String,
				Optional:            true,
			},
			"network_location_ids": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Network Locations associated with the ZTA Profile.").String,
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(50),
				},
			},
			"private_access_dns_steering_destinations": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("DNS steering destinations for secure private access.").String,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(20),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("ID of the DNS steering destination.").String,
							Required:            true,
						},
					},
				},
			},
			"private_access_use_trusted_networks": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether Trusted Network Detection is enabled for secure private access.").String,
				Optional:            true,
			},
			"private_access_end_user_pause": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether ZTA enforcement pause is enabled for secure private access.").String,
				Optional:            true,
			},
			"private_access_end_user_pause_duration": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Enforcement pause duration in minutes for secure private access. When not set while the pause is enabled, the pause is indefinite.").AddIntegerRangeDescription(1, 1440).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 1440),
				},
			},
			"internet_access_steering_mode": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Internet steering mode: `0` - no internet destinations, `1` - all internet destinations, `2` - selected internet destinations.").AddIntegerRangeDescription(0, 2).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 2),
				},
			},
			"internet_access_use_trusted_networks": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether Trusted Network Detection is enabled for secure internet access.").String,
				Optional:            true,
			},
			"internet_access_end_user_pause": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether ZTA enforcement pause is enabled for secure internet access.").String,
				Optional:            true,
			},
			"internet_access_end_user_pause_duration": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Enforcement pause duration in minutes for secure internet access. When not set while the pause is enabled, the pause is indefinite.").AddIntegerRangeDescription(1, 1440).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 1440),
				},
			},
		},
	}
}

func (r *ZTAProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

var _ resource.ResourceWithValidateConfig = &ZTAProfileResource{}

// ValidateConfig rejects changes of the operating systems and network locations of the default profile, which cannot be changed.
func (r *ZTAProfileResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config ZTAProfile

	diags := req.Config.Get(ctx, &config)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if config.Name.ValueString() != ztaDefaultProfileName {
		return
	}

	operatingSystems := map[string]types.Bool{
		"operating_system_macos":        config.OperatingSystemMacos,
		"operating_system_windows":      config.OperatingSystemWindows,
		"operating_system_linux":        config.OperatingSystemLinux,
		"operating_system_ios":          config.OperatingSystemIos,
		"operating_system_android":      config.OperatingSystemAndroid,
		"operating_system_android_knox": config.OperatingSystemAndroidKnox,
	}
	for _, name := range slices.Sorted(maps.Keys(operatingSystems)) {
		if !operatingSystems[name].IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root(name), "Invalid Attribute Value",
				fmt.Sprintf("Operating systems of the %s cannot be changed, `%s` must not be set.", ztaDefaultProfileName, name))
		}
	}

	if !config.NetworkLocationIds.IsNull() && (config.NetworkLocationIds.IsUnknown() || len(config.NetworkLocationIds.Elements()) > 0) {
		resp.Diagnostics.AddAttributeError(path.Root("network_location_ids"), "Invalid Attribute Value",
			fmt.Sprintf("Network locations of the %s cannot be changed, `network_location_ids` must be empty or not set.", ztaDefaultProfileName))
	}
}

func (r *ZTAProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ztaProfileMu.Lock()
	defer ztaProfileMu.Unlock()

	var plan ZTAProfile

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	body := plan.toBody(ctx, ZTAProfile{})
	body = plan.adjustBody(ctx, body)
	var res secureaccess.Res
	var err error
	if plan.Name.ValueString() == ztaDefaultProfileName {
		// The default profile always exists with a fixed ID and cannot be created, it is read to get its revision, and updated
		plan.Id = types.StringValue(ztaDefaultProfileId)
		res, err = r.client.Get(ctx, plan.getPath()+"/"+ztaDefaultProfileId, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve default profile (GET), got error: %s, %s", err, res.String()))
			return
		}
		body, _ = sjson.Set(body, "rev", res.Get("rev").Int())
		res, err = r.client.Put(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString()), body, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s, %s", err, res.String()))
			return
		}
	} else {
		// Create object
		res, err = r.client.Post(ctx, plan.getPath(), body, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST/PUT), got error: %s, %s", err, res.String()))
			return
		}
		plan.Id = types.StringValue(res.Get("profileId").String())
	}
	plan.fromBodyUnknowns(ctx, res.Result)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (r *ZTAProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ZTAProfile

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

func (r *ZTAProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ztaProfileMu.Lock()
	defer ztaProfileMu.Unlock()

	var plan, state ZTAProfile

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
	res, err := r.client.Put(ctx, plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString()), body, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ZTAProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ztaProfileMu.Lock()
	defer ztaProfileMu.Unlock()

	var state ZTAProfile

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))
	if state.Id.ValueString() == ztaDefaultProfileId {
		// The default profile cannot be deleted, it is reset instead
		body := state.toBodyPutDelete(ctx)
		res, err := r.client.Put(ctx, state.getPath()+"/"+url.QueryEscape(state.Id.ValueString()), body, reqMods...)
		if err != nil && !strings.Contains(err.Error(), "StatusCode 404") {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (PUT), got error: %s, %s", err, res.String()))
			return
		}
	} else {
		res, err := r.client.Delete(ctx, state.getPath()+"/"+url.QueryEscape(state.Id.ValueString()), reqMods...)
		if err != nil && !strings.Contains(err.Error(), "StatusCode 404") {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *ZTAProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
