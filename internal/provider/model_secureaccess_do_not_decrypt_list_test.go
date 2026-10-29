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

// Manual unit tests (no API needed): the update logic of the Do Not Decrypt List builds the add/remove requests the API expects

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func dndDestinations(n int) []DoNotDecryptListDestinations {
	items := make([]DoNotDecryptListDestinations, n)
	for i := range items {
		items[i] = DoNotDecryptListDestinations{Destination: types.StringValue(fmt.Sprintf("host%d.example.com", i)), Comment: types.StringNull()}
	}
	return items
}

func dndDestination(destination, comment string) DoNotDecryptListDestinations {
	item := DoNotDecryptListDestinations{Destination: types.StringValue(destination), Comment: types.StringNull()}
	if comment != "" {
		item.Comment = types.StringValue(comment)
	}
	return item
}

func dndInt64Set(values ...int64) types.Set {
	if values == nil {
		return types.SetNull(types.Int64Type)
	}
	elements := make([]attr.Value, len(values))
	for i, v := range values {
		elements[i] = types.Int64Value(v)
	}
	return types.SetValueMust(types.Int64Type, elements)
}

func TestDoNotDecryptListChangeStepsChunking(t *testing.T) {
	for _, tc := range []struct {
		destinations, applications int
		wantSteps                  []int // number of destinations per request
	}{
		{0, 0, nil},
		{0, 2, []int{0}},
		{1, 0, []int{1}},
		{500, 0, []int{500}},
		{501, 0, []int{500, 1}},
		{1000, 0, []int{500, 500}},
		{1001, 2, []int{500, 500, 1}},
	} {
		for _, add := range []bool{true, false} {
			var applications []int64
			for i := 0; i < tc.applications; i++ {
				applications = append(applications, int64(100+i))
			}
			steps := doNotDecryptListChangeSteps("name", add, dndDestinations(tc.destinations), applications, nil)
			if len(steps) != len(tc.wantSteps) {
				t.Fatalf("destinations=%d applications=%d add=%v: got %d requests, want %d", tc.destinations, tc.applications, add, len(steps), len(tc.wantSteps))
			}
			verb := map[bool]string{true: "add", false: "remove"}[add]
			for i, step := range steps {
				if got := len(gjson.Get(step.body, verb+"Destinations").Array()); got != tc.wantSteps[i] {
					t.Errorf("destinations=%d add=%v request %d: got %d destinations, want %d", tc.destinations, add, i, got, tc.wantSteps[i])
				}
				if name := gjson.Get(step.body, "name").String(); name != "name" {
					t.Errorf("request %d: name = %q, every request must carry the name", i, name)
				}
				if got, want := gjson.Get(step.body, verb+"Applications").Exists(), i == 0 && tc.applications > 0; got != want {
					t.Errorf("destinations=%d add=%v request %d: applications present = %v, want %v (only in the first request)", tc.destinations, add, i, got, want)
				}
			}
		}
	}
}

func TestDoNotDecryptListToUpdateSteps(t *testing.T) {
	ctx := context.Background()
	state := DoNotDecryptList{
		Name:                       types.StringValue("old"),
		Destinations:               []DoNotDecryptListDestinations{dndDestination("a.com", "x"), dndDestination("b.com", "y"), dndDestination("keep.com", "")},
		Applications:               dndInt64Set(1, 2),
		DecryptExceptionCategories: dndInt64Set(10),
	}

	t.Run("name only", func(t *testing.T) {
		plan := state
		plan.Name = types.StringValue("new")
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 1 || steps[0].body != `{"name":"new"}` {
			t.Fatalf("got %+v", steps)
		}
	})

	t.Run("removals are sent before additions, removals carry no comment", func(t *testing.T) {
		plan := state
		plan.Name = types.StringValue("new")
		// a.com: comment changed, b.com: removed, c.com: added, keep.com: unchanged
		plan.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", "changed"), dndDestination("keep.com", ""), dndDestination("c.com", "z")}
		plan.Applications = dndInt64Set(2, 3)
		plan.DecryptExceptionCategories = types.SetNull(types.Int64Type)
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 2 {
			t.Fatalf("got %d requests, want 2: %+v", len(steps), steps)
		}
		remove, add := steps[0].body, steps[1].body
		if got := gjson.Get(remove, "removeDestinations.#.destination").Raw; got != `["a.com","b.com"]` {
			t.Errorf("removed destinations = %s", got)
		}
		if gjson.Get(remove, "removeDestinations.0.comment").Exists() {
			t.Errorf("removal must not carry a comment: %s", remove)
		}
		if got := gjson.Get(remove, "removeApplications").Raw; got != `[1]` {
			t.Errorf("removed applications = %s", got)
		}
		if got := gjson.Get(remove, "removeDecryptExceptionCategories").Raw; got != `[10]` {
			t.Errorf("a null set in the plan must remove the categories, got %s", got)
		}
		if gjson.Get(remove, "addDestinations").Exists() || gjson.Get(add, "removeDestinations").Exists() {
			t.Errorf("additions and removals must be in separate requests: %s / %s", remove, add)
		}
		if got := gjson.Get(add, "addDestinations.#.destination").Raw; got != `["a.com","c.com"]` {
			t.Errorf("added destinations = %s", got)
		}
		if got := gjson.Get(add, "addDestinations.0.comment").String(); got != "changed" {
			t.Errorf("added destination comment = %q", got)
		}
		if got := gjson.Get(add, "addApplications").Raw; got != `[3]` {
			t.Errorf("added applications = %s", got)
		}
	})

	t.Run("null destinations in the plan remove all", func(t *testing.T) {
		plan := state
		plan.Destinations = nil
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 1 || len(gjson.Get(steps[0].body, "removeDestinations").Array()) != 3 {
			t.Fatalf("got %+v", steps)
		}
	})
}

func TestDoNotDecryptListApplyStep(t *testing.T) {
	ctx := context.Background()
	state := DoNotDecryptList{
		Name:                       types.StringValue("old"),
		Destinations:               []DoNotDecryptListDestinations{dndDestination("a.com", "x"), dndDestination("b.com", "y")},
		Applications:               dndInt64Set(1, 2),
		DecryptExceptionCategories: dndInt64Set(10),
	}
	plan := state
	plan.Name = types.StringValue("new")
	plan.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", "changed"), dndDestination("c.com", "")}
	plan.Applications = dndInt64Set(2, 3)
	plan.DecryptExceptionCategories = types.SetNull(types.Int64Type)

	steps := plan.toUpdateSteps(ctx, state)
	original := slices.Clone(state.Destinations)
	interim := state
	interim.applyStep(ctx, plan.Name, steps[0]) // only the removal request was accepted

	if interim.Name.ValueString() != "new" {
		t.Errorf("name = %q", interim.Name.ValueString())
	}
	if len(interim.Destinations) != 0 {
		t.Errorf("after the removal request no destination is left, got %+v", interim.Destinations)
	}
	if got := doNotDecryptListInt64Values(ctx, interim.Applications); !slices.Equal(got, []int64{2}) {
		t.Errorf("applications = %v, want [2]", got)
	}
	if !interim.DecryptExceptionCategories.IsNull() {
		t.Errorf("categories = %v, want null", interim.DecryptExceptionCategories)
	}
	if !slices.EqualFunc(original, state.Destinations, func(a, b DoNotDecryptListDestinations) bool { return a == b }) {
		t.Errorf("the state the interim copy was made from must not be modified")
	}

	interim.applyStep(ctx, plan.Name, steps[1]) // the addition request
	if len(interim.Destinations) != 2 || interim.Destinations[0].Destination.ValueString() != "a.com" || interim.Destinations[0].Comment.ValueString() != "changed" || interim.Destinations[1].Destination.ValueString() != "c.com" {
		t.Errorf("destinations = %+v", interim.Destinations)
	}
	if got := doNotDecryptListInt64Values(ctx, interim.Applications); !slices.Equal(got, []int64{2, 3}) {
		t.Errorf("applications = %v, want [2 3]", got)
	}
}
