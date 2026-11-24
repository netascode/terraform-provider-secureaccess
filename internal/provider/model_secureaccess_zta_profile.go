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
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type ZTAProfile struct {
	Id                                   types.String                                     `tfsdk:"id"`
	Name                                 types.String                                     `tfsdk:"name"`
	Revision                             types.Int64                                      `tfsdk:"revision"`
	Priority                             types.Int64                                      `tfsdk:"priority"`
	OperatingSystemMacos                 types.Bool                                       `tfsdk:"operating_system_macos"`
	OperatingSystemWindows               types.Bool                                       `tfsdk:"operating_system_windows"`
	OperatingSystemLinux                 types.Bool                                       `tfsdk:"operating_system_linux"`
	OperatingSystemIos                   types.Bool                                       `tfsdk:"operating_system_ios"`
	OperatingSystemAndroid               types.Bool                                       `tfsdk:"operating_system_android"`
	OperatingSystemAndroidKnox           types.Bool                                       `tfsdk:"operating_system_android_knox"`
	AllUsers                             types.Bool                                       `tfsdk:"all_users"`
	AllGroups                            types.Bool                                       `tfsdk:"all_groups"`
	NetworkLocationIds                   types.Set                                        `tfsdk:"network_location_ids"`
	PrivateAccessDnsSteeringDestinations []ZTAProfilePrivateAccessDnsSteeringDestinations `tfsdk:"private_access_dns_steering_destinations"`
	PrivateAccessUseTrustedNetworks      types.Bool                                       `tfsdk:"private_access_use_trusted_networks"`
	PrivateAccessEndUserPause            types.Bool                                       `tfsdk:"private_access_end_user_pause"`
	PrivateAccessEndUserPauseDuration    types.Int64                                      `tfsdk:"private_access_end_user_pause_duration"`
	InternetAccessSteeringMode           types.Int64                                      `tfsdk:"internet_access_steering_mode"`
	InternetAccessUseTrustedNetworks     types.Bool                                       `tfsdk:"internet_access_use_trusted_networks"`
	InternetAccessEndUserPause           types.Bool                                       `tfsdk:"internet_access_end_user_pause"`
	InternetAccessEndUserPauseDuration   types.Int64                                      `tfsdk:"internet_access_end_user_pause_duration"`
}

type ZTAProfilePrivateAccessDnsSteeringDestinations struct {
	Id types.String `tfsdk:"id"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ZTAProfile) getPath() string {
	return "/deployments/v2/ztna/profiles"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ZTAProfile) toBody(ctx context.Context, state ZTAProfile) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "profileId", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "profileName", data.Name.ValueString())
	}
	if !data.Priority.IsNull() {
		body, _ = sjson.Set(body, "priority", data.Priority.ValueInt64())
	}
	if !data.OperatingSystemMacos.IsNull() {
		body, _ = sjson.Set(body, "operatingSystems.mac-intel.enabled", data.OperatingSystemMacos.ValueBool())
	}
	if !data.OperatingSystemWindows.IsNull() {
		body, _ = sjson.Set(body, "operatingSystems.win.enabled", data.OperatingSystemWindows.ValueBool())
	}
	if !data.OperatingSystemLinux.IsNull() {
		body, _ = sjson.Set(body, "operatingSystems.linux-64.enabled", data.OperatingSystemLinux.ValueBool())
	}
	if !data.OperatingSystemIos.IsNull() {
		body, _ = sjson.Set(body, "operatingSystems.apple-ios.enabled", data.OperatingSystemIos.ValueBool())
	}
	if !data.OperatingSystemAndroid.IsNull() {
		body, _ = sjson.Set(body, "operatingSystems.android.generic-android.enabled", data.OperatingSystemAndroid.ValueBool())
	}
	if !data.OperatingSystemAndroidKnox.IsNull() {
		body, _ = sjson.Set(body, "operatingSystems.android.knox-android.enabled", data.OperatingSystemAndroidKnox.ValueBool())
	}
	if !data.AllUsers.IsNull() {
		body, _ = sjson.Set(body, "usersData.allUsersEnabled", data.AllUsers.ValueBool())
	}
	if !data.AllGroups.IsNull() {
		body, _ = sjson.Set(body, "groupsData.allGroupsEnabled", data.AllGroups.ValueBool())
	}
	if !data.NetworkLocationIds.IsNull() {
		var values []string
		data.NetworkLocationIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "networkLocationIds", values)
	}
	if len(data.PrivateAccessDnsSteeringDestinations) > 0 {
		var privateAccessDnsSteeringDestinationsBody strings.Builder
		privateAccessDnsSteeringDestinationsBody.WriteString("[")
		for _, item := range data.PrivateAccessDnsSteeringDestinations {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueString())
			}
			if itemBody != "" {
				if privateAccessDnsSteeringDestinationsBody.Len() > 1 {
					privateAccessDnsSteeringDestinationsBody.WriteString(",")
				}
				privateAccessDnsSteeringDestinationsBody.WriteString(itemBody)
			}
		}
		privateAccessDnsSteeringDestinationsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "securePrivateAccess.dnsSteeringDestinations", privateAccessDnsSteeringDestinationsBody.String())
	}
	if !data.PrivateAccessUseTrustedNetworks.IsNull() {
		body, _ = sjson.Set(body, "securePrivateAccess.trustedNetworksEnabled", data.PrivateAccessUseTrustedNetworks.ValueBool())
	}
	if !data.PrivateAccessEndUserPause.IsNull() {
		body, _ = sjson.Set(body, "securePrivateAccess.enforcementPause.enabled", data.PrivateAccessEndUserPause.ValueBool())
	}
	if !data.PrivateAccessEndUserPauseDuration.IsNull() {
		body, _ = sjson.Set(body, "securePrivateAccess.enforcementPause.durationMinutes", data.PrivateAccessEndUserPauseDuration.ValueInt64())
	}
	if !data.InternetAccessSteeringMode.IsNull() {
		body, _ = sjson.Set(body, "secureInternetAccess.steeringMode", data.InternetAccessSteeringMode.ValueInt64())
	}
	if !data.InternetAccessUseTrustedNetworks.IsNull() {
		body, _ = sjson.Set(body, "secureInternetAccess.trustedNetworksEnabled", data.InternetAccessUseTrustedNetworks.ValueBool())
	}
	if !data.InternetAccessEndUserPause.IsNull() {
		body, _ = sjson.Set(body, "secureInternetAccess.enforcementPause.enabled", data.InternetAccessEndUserPause.ValueBool())
	}
	if !data.InternetAccessEndUserPauseDuration.IsNull() {
		body, _ = sjson.Set(body, "secureInternetAccess.enforcementPause.durationMinutes", data.InternetAccessEndUserPauseDuration.ValueInt64())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ZTAProfile) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("profileName"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("rev"); value.Exists() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
	if value := res.Get("priority"); value.Exists() {
		data.Priority = types.Int64Value(value.Int())
	} else {
		data.Priority = types.Int64Null()
	}
	if value := res.Get("operatingSystems.mac-intel.enabled"); value.Exists() {
		data.OperatingSystemMacos = types.BoolValue(value.Bool())
	} else {
		data.OperatingSystemMacos = types.BoolValue(true)
	}
	if value := res.Get("operatingSystems.win.enabled"); value.Exists() {
		data.OperatingSystemWindows = types.BoolValue(value.Bool())
	} else {
		data.OperatingSystemWindows = types.BoolValue(true)
	}
	if value := res.Get("operatingSystems.linux-64.enabled"); value.Exists() {
		data.OperatingSystemLinux = types.BoolValue(value.Bool())
	} else {
		data.OperatingSystemLinux = types.BoolValue(true)
	}
	if value := res.Get("operatingSystems.apple-ios.enabled"); value.Exists() {
		data.OperatingSystemIos = types.BoolValue(value.Bool())
	} else {
		data.OperatingSystemIos = types.BoolValue(true)
	}
	if value := res.Get("operatingSystems.android.generic-android.enabled"); value.Exists() {
		data.OperatingSystemAndroid = types.BoolValue(value.Bool())
	} else {
		data.OperatingSystemAndroid = types.BoolValue(true)
	}
	if value := res.Get("operatingSystems.android.knox-android.enabled"); value.Exists() {
		data.OperatingSystemAndroidKnox = types.BoolValue(value.Bool())
	} else {
		data.OperatingSystemAndroidKnox = types.BoolValue(true)
	}
	if value := res.Get("usersData.allUsersEnabled"); value.Exists() {
		data.AllUsers = types.BoolValue(value.Bool())
	} else {
		data.AllUsers = types.BoolNull()
	}
	if value := res.Get("groupsData.allGroupsEnabled"); value.Exists() {
		data.AllGroups = types.BoolValue(value.Bool())
	} else {
		data.AllGroups = types.BoolNull()
	}
	if value := res.Get("networkLocationIds"); value.Exists() {
		data.NetworkLocationIds = helpers.GetStringSet(value.Array())
	} else {
		data.NetworkLocationIds = types.SetNull(types.StringType)
	}
	if value := res.Get("securePrivateAccess.dnsSteeringDestinations"); value.Exists() {
		data.PrivateAccessDnsSteeringDestinations = make([]ZTAProfilePrivateAccessDnsSteeringDestinations, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTAProfilePrivateAccessDnsSteeringDestinations{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.StringValue(value.String())
			} else {
				data.Id = types.StringNull()
			}
			(*parent).PrivateAccessDnsSteeringDestinations = append((*parent).PrivateAccessDnsSteeringDestinations, data)
			return true
		})
	}
	if value := res.Get("securePrivateAccess.trustedNetworksEnabled"); value.Exists() {
		data.PrivateAccessUseTrustedNetworks = types.BoolValue(value.Bool())
	} else {
		data.PrivateAccessUseTrustedNetworks = types.BoolNull()
	}
	if value := res.Get("securePrivateAccess.enforcementPause.enabled"); value.Exists() {
		data.PrivateAccessEndUserPause = types.BoolValue(value.Bool())
	} else {
		data.PrivateAccessEndUserPause = types.BoolNull()
	}
	if value := res.Get("securePrivateAccess.enforcementPause.durationMinutes"); value.Exists() {
		data.PrivateAccessEndUserPauseDuration = types.Int64Value(value.Int())
	} else {
		data.PrivateAccessEndUserPauseDuration = types.Int64Null()
	}
	if value := res.Get("secureInternetAccess.steeringMode"); value.Exists() {
		data.InternetAccessSteeringMode = types.Int64Value(value.Int())
	} else {
		data.InternetAccessSteeringMode = types.Int64Null()
	}
	if value := res.Get("secureInternetAccess.trustedNetworksEnabled"); value.Exists() {
		data.InternetAccessUseTrustedNetworks = types.BoolValue(value.Bool())
	} else {
		data.InternetAccessUseTrustedNetworks = types.BoolNull()
	}
	if value := res.Get("secureInternetAccess.enforcementPause.enabled"); value.Exists() {
		data.InternetAccessEndUserPause = types.BoolValue(value.Bool())
	} else {
		data.InternetAccessEndUserPause = types.BoolNull()
	}
	if value := res.Get("secureInternetAccess.enforcementPause.durationMinutes"); value.Exists() {
		data.InternetAccessEndUserPauseDuration = types.Int64Value(value.Int())
	} else {
		data.InternetAccessEndUserPauseDuration = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *ZTAProfile) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("profileName"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("rev"); value.Exists() && !data.Revision.IsNull() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
	if value := res.Get("priority"); value.Exists() && !data.Priority.IsNull() {
		data.Priority = types.Int64Value(value.Int())
	} else {
		data.Priority = types.Int64Null()
	}
	if value := res.Get("operatingSystems.mac-intel.enabled"); value.Exists() && !data.OperatingSystemMacos.IsNull() {
		data.OperatingSystemMacos = types.BoolValue(value.Bool())
	} else if data.OperatingSystemMacos.ValueBool() != true {
		data.OperatingSystemMacos = types.BoolNull()
	}
	if value := res.Get("operatingSystems.win.enabled"); value.Exists() && !data.OperatingSystemWindows.IsNull() {
		data.OperatingSystemWindows = types.BoolValue(value.Bool())
	} else if data.OperatingSystemWindows.ValueBool() != true {
		data.OperatingSystemWindows = types.BoolNull()
	}
	if value := res.Get("operatingSystems.linux-64.enabled"); value.Exists() && !data.OperatingSystemLinux.IsNull() {
		data.OperatingSystemLinux = types.BoolValue(value.Bool())
	} else if data.OperatingSystemLinux.ValueBool() != true {
		data.OperatingSystemLinux = types.BoolNull()
	}
	if value := res.Get("operatingSystems.apple-ios.enabled"); value.Exists() && !data.OperatingSystemIos.IsNull() {
		data.OperatingSystemIos = types.BoolValue(value.Bool())
	} else if data.OperatingSystemIos.ValueBool() != true {
		data.OperatingSystemIos = types.BoolNull()
	}
	if value := res.Get("operatingSystems.android.generic-android.enabled"); value.Exists() && !data.OperatingSystemAndroid.IsNull() {
		data.OperatingSystemAndroid = types.BoolValue(value.Bool())
	} else if data.OperatingSystemAndroid.ValueBool() != true {
		data.OperatingSystemAndroid = types.BoolNull()
	}
	if value := res.Get("operatingSystems.android.knox-android.enabled"); value.Exists() && !data.OperatingSystemAndroidKnox.IsNull() {
		data.OperatingSystemAndroidKnox = types.BoolValue(value.Bool())
	} else if data.OperatingSystemAndroidKnox.ValueBool() != true {
		data.OperatingSystemAndroidKnox = types.BoolNull()
	}
	if value := res.Get("usersData.allUsersEnabled"); value.Exists() && !data.AllUsers.IsNull() {
		data.AllUsers = types.BoolValue(value.Bool())
	} else {
		data.AllUsers = types.BoolNull()
	}
	if value := res.Get("groupsData.allGroupsEnabled"); value.Exists() && !data.AllGroups.IsNull() {
		data.AllGroups = types.BoolValue(value.Bool())
	} else {
		data.AllGroups = types.BoolNull()
	}
	if value := res.Get("networkLocationIds"); value.Exists() && !data.NetworkLocationIds.IsNull() {
		data.NetworkLocationIds = helpers.GetStringSet(value.Array())
	} else {
		data.NetworkLocationIds = types.SetNull(types.StringType)
	}
	privateAccessDnsSteeringDestinationsArray := res.Get("securePrivateAccess.dnsSteeringDestinations")
	for i := 0; i < len(data.PrivateAccessDnsSteeringDestinations); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{data.PrivateAccessDnsSteeringDestinations[i].Id.ValueString()}

		parent := &data
		data := (*parent).PrivateAccessDnsSteeringDestinations[i]
		var res gjson.Result

		privateAccessDnsSteeringDestinationsArray.ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if helpers.GetKeyValue(v.Get(keys[ik])) != keyValues[ik] {
						found = false
						break
					}
					found = true
				}
				if found {
					res = v
					return false
				}
				return true
			},
		)
		if !res.Exists() {
			tflog.Debug(ctx, fmt.Sprintf("removing PrivateAccessDnsSteeringDestinations[%d] = %+v",
				i,
				(*parent).PrivateAccessDnsSteeringDestinations[i],
			))
			(*parent).PrivateAccessDnsSteeringDestinations = slices.Delete((*parent).PrivateAccessDnsSteeringDestinations, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.StringValue(value.String())
		} else {
			data.Id = types.StringNull()
		}
		(*parent).PrivateAccessDnsSteeringDestinations[i] = data
	}
	if value := res.Get("securePrivateAccess.trustedNetworksEnabled"); value.Exists() && !data.PrivateAccessUseTrustedNetworks.IsNull() {
		data.PrivateAccessUseTrustedNetworks = types.BoolValue(value.Bool())
	} else {
		data.PrivateAccessUseTrustedNetworks = types.BoolNull()
	}
	if value := res.Get("securePrivateAccess.enforcementPause.enabled"); value.Exists() && !data.PrivateAccessEndUserPause.IsNull() {
		data.PrivateAccessEndUserPause = types.BoolValue(value.Bool())
	} else {
		data.PrivateAccessEndUserPause = types.BoolNull()
	}
	if value := res.Get("securePrivateAccess.enforcementPause.durationMinutes"); value.Exists() && !data.PrivateAccessEndUserPauseDuration.IsNull() {
		data.PrivateAccessEndUserPauseDuration = types.Int64Value(value.Int())
	} else {
		data.PrivateAccessEndUserPauseDuration = types.Int64Null()
	}
	if value := res.Get("secureInternetAccess.steeringMode"); value.Exists() && !data.InternetAccessSteeringMode.IsNull() {
		data.InternetAccessSteeringMode = types.Int64Value(value.Int())
	} else {
		data.InternetAccessSteeringMode = types.Int64Null()
	}
	if value := res.Get("secureInternetAccess.trustedNetworksEnabled"); value.Exists() && !data.InternetAccessUseTrustedNetworks.IsNull() {
		data.InternetAccessUseTrustedNetworks = types.BoolValue(value.Bool())
	} else {
		data.InternetAccessUseTrustedNetworks = types.BoolNull()
	}
	if value := res.Get("secureInternetAccess.enforcementPause.enabled"); value.Exists() && !data.InternetAccessEndUserPause.IsNull() {
		data.InternetAccessEndUserPause = types.BoolValue(value.Bool())
	} else {
		data.InternetAccessEndUserPause = types.BoolNull()
	}
	if value := res.Get("secureInternetAccess.enforcementPause.durationMinutes"); value.Exists() && !data.InternetAccessEndUserPauseDuration.IsNull() {
		data.InternetAccessEndUserPauseDuration = types.Int64Value(value.Int())
	} else {
		data.InternetAccessEndUserPauseDuration = types.Int64Null()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ZTAProfile) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	if data.Revision.IsUnknown() {
		if value := res.Get("rev"); value.Exists() {
			data.Revision = types.Int64Value(value.Int())
		} else {
			data.Revision = types.Int64Null()
		}
	}
}

// End of section. //template:end fromBodyUnknowns

// toBodyPutDelete is used to create the body for PUT requests to reset the default profile, which cannot be deleted.
// Only the fields required by the API are set, with the revision from state.
func (data ZTAProfile) toBodyPutDelete(ctx context.Context) string {
	body := ""
	body, _ = sjson.Set(body, "profileName", ztaDefaultProfileName)
	body, _ = sjson.Set(body, "priority", 0)
	body, _ = sjson.Set(body, "rev", data.Revision.ValueInt64())
	return body
}

func (data ZTAProfile) adjustBody(ctx context.Context, req string) string {

	req, _ = sjson.Delete(req, "profileId")
	// Revision is kept from state in the plan (unknown on create), the API requires the current one on update
	if !data.Revision.IsUnknown() && !data.Revision.IsNull() {
		req, _ = sjson.Set(req, "rev", data.Revision.ValueInt64())
	}

	return req
}

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
