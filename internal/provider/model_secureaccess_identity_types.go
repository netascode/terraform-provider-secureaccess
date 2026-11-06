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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type IdentityTypes struct {
	Filter types.String                  `tfsdk:"filter"`
	Items  map[string]IdentityTypesItems `tfsdk:"items"`
}

type IdentityTypesItems struct {
	Id types.Int64 `tfsdk:"id"`
}

// End of section. //template:end types

// identityTypesStatic holds hardcoded Identity Types until a relevant API endpoint is available
type identityTypeStatic struct {
	Id    int64  `json:"id"`
	Label string `json:"label"`
}

var identityTypesStatic = []identityTypeStatic{
	{Id: 1, Label: "Registered Networks"},
	{Id: 3, Label: "AD Groups"},
	{Id: 5, Label: "AD Devices"},
	{Id: 7, Label: "AD Users"},
	// API bundles 9 & 11 into macOS/Windows Devices
	// 9 & 34 and 36 & 50 must be used together, see identityTypesCombined
	// {Id: 9, Label: "Windows Devices"},
	{Id: 11, Label: "Apple Devices"},
	// {Id: 13, Label: "Virtual Appliances"},
	// {Id: 15, Label: "AD Domains"},
	// {Id: 17, Label: "On Network"},
	{Id: 19, Label: "Internal Networks"},
	{Id: 21, Label: "Sites"},
	// {Id: 22, Label: "Organization"},
	// {Id: 24, Label: "Legacy Device"},
	// {Id: 26, Label: "Legacy Policy"},
	// {Id: 28, Label: "AD Connector"},
	// {Id: 30, Label: "Domain Controller"},
	{Id: 32, Label: "Network Devices"},
	// {Id: 34, Label: "macOS/Windows Devices"},
	// {Id: 36, Label: "iOS Mobile Devices"},
	{Id: 38, Label: "ChromeOS Devices"},
	{Id: 40, Label: "Network Tunnel Groups"},
	// {Id: 41, Label: "GSuite Domains"},
	{Id: 43, Label: "Google Workspace Users"},
	{Id: 45, Label: "Google Workspace Organizational Units"},
	// {Id: 46, Label: "Web Groups"},
	// {Id: 48, Label: "Web Users"},
	// {Id: 50, Label: "Android Mobile Devices"},
	// {Id: 51, Label: "Windows Mobile devices"},
	{Id: 52, Label: "Catalyst SD-WAN Service VPN IDs"},
	{Id: 54, Label: "Security Group Tags"},
	// {Id: 57, Label: "ZTNA Client"},
	// {Id: 62, Label: "Store Branch"},
}

// identityTypesCombined are provider-only Identity Types (negative IDs) standing for IDs the API requires together.
// Access rules expand them on write and collapse them on read (helpers_access_rule_conditions.go).
var identityTypesCombined = []struct {
	Id    int64
	Label string
	Ids   []int64
}{
	{-1, "iOS/Android Devices", []int64{36, 50}},
	{-2, "macOS/Windows Devices", []int64{34, 9}},
}

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data IdentityTypes) getPath() string {
	return "NOT_AVAILABLE"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *IdentityTypes) fromBody(ctx context.Context, res gjson.Result) {
	// Build lookup map for O(1) access
	itemsByName := make(map[string]gjson.Result)
	res.ForEach(func(_, v gjson.Result) bool {
		if name := v.Get("label").String(); name != "" {
			// Item names are matched case-insensitively
			itemsByName[strings.ToLower(name)] = v
		}
		return true
	})
	for k := range data.Items {
		parent := &data
		data := (*parent).Items[k]
		res, found := itemsByName[strings.ToLower(k)]
		if !found {
			tflog.Debug(ctx, fmt.Sprintf("subresource not found, removing: name=%v", k))
			delete((*parent).Items, k)
			continue
		}
		if value := res.Get("id"); value.Exists() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		(*parent).Items[k] = data
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
