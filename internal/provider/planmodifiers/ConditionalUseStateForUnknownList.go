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

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.List = conditionalUseStateForUnknownList{}

// ConditionalUseStateForUnknownList keeps the prior state value of a computed list,
// unless any of the dependent (string) attributes is changing, in which case the list is left as `known after apply`.
func ConditionalUseStateForUnknownList(dependentAttributes ...string) planmodifier.List {
	return conditionalUseStateForUnknownList{dependentAttributes: dependentAttributes}
}

type conditionalUseStateForUnknownList struct {
	dependentAttributes []string
}

func (c conditionalUseStateForUnknownList) Description(ctx context.Context) string {
	return "Conditionally set to `known after apply` based on change of `dependentAttributes`"
}

func (c conditionalUseStateForUnknownList) MarkdownDescription(ctx context.Context) string {
	return c.Description(ctx)
}

func (c conditionalUseStateForUnknownList) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// Do nothing on resource creation, or when the value is configured
	if req.StateValue.IsNull() || !req.PlanValue.IsUnknown() {
		return
	}

	for _, dependentAttribute := range c.dependentAttributes {
		// Get path to the dependent attribute
		dependentAttributePath := req.Path.ParentPath().AtName(dependentAttribute)

		// Get state value of dependent attribute
		var stateValue types.String
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, dependentAttributePath, &stateValue)...)

		// Get plan value of dependent attribute
		var planValue types.String
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, dependentAttributePath, &planValue)...)

		// If state and plan values don't match, do nothing (sets value to `known after apply`)
		if resp.Diagnostics.HasError() || !stateValue.Equal(planValue) {
			return
		}
	}

	// If all dependent attributes match, set plan to state value
	resp.PlanValue = req.StateValue
}
