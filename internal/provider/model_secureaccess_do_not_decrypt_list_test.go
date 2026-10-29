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

func TestDoNotDecryptListDestinationChanges(t *testing.T) {
	state := []DoNotDecryptListDestinations{dndDestination("same.com", "x"), dndDestination("changed.com", "x"), dndDestination("cleared.com", "x"), dndDestination("gets-comment.com", ""), dndDestination("gone.com", "x")}
	plan := []DoNotDecryptListDestinations{dndDestination("same.com", "x"), dndDestination("changed.com", "y"), dndDestination("cleared.com", ""), dndDestination("gets-comment.com", "z"), dndDestination("new.com", "")}

	upsert, remove, readd := doNotDecryptListDestinationChanges(state, plan)
	names := func(items []DoNotDecryptListDestinations) []string {
		var out []string
		for _, item := range items {
			out = append(out, item.Destination.ValueString())
		}
		return out
	}
	if got, want := names(upsert), []string{"changed.com", "gets-comment.com", "new.com"}; !slices.Equal(got, want) {
		t.Errorf("upsert = %v, want %v (a new comment is set by adding the destination again, without removing it)", got, want)
	}
	if got, want := names(remove), []string{"cleared.com", "gone.com"}; !slices.Equal(got, want) {
		t.Errorf("remove = %v, want %v", got, want)
	}
	if got, want := names(readd), []string{"cleared.com"}; !slices.Equal(got, want) {
		t.Errorf("readd = %v, want %v (adding without a comment does not clear an existing one)", got, want)
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
		if len(steps) != 1 || steps[0].body != `{"name":"new"}` || steps[0].viaList {
			t.Fatalf("got %+v", steps)
		}
	})

	t.Run("additions come first and update comments, removals carry no comment", func(t *testing.T) {
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
		add, remove := steps[0].body, steps[1].body
		if got := gjson.Get(add, "addDestinations.#.destination").Raw; got != `["a.com","c.com"]` {
			t.Errorf("added destinations = %s", got)
		}
		if got := gjson.Get(add, "addDestinations.0.comment").String(); got != "changed" {
			t.Errorf("added destination comment = %q", got)
		}
		if got := gjson.Get(add, "addApplications").Raw; got != `[3]` {
			t.Errorf("added applications = %s", got)
		}
		if gjson.Get(add, "removeDestinations").Exists() || gjson.Get(remove, "addDestinations").Exists() {
			t.Errorf("additions and removals must be in separate requests: %s / %s", add, remove)
		}
		if got := gjson.Get(remove, "removeDestinations.#.destination").Raw; got != `["b.com"]` {
			t.Errorf("removed destinations = %s (a changed comment must not remove the destination)", got)
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
	})

	t.Run("IP addresses and CIDR blocks are removed through the destination list, after the PUTs", func(t *testing.T) {
		state := state
		state.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", ""), dndDestination("10.10.10.10", ""), dndDestination("192.0.2.0/24", "x"), dndDestination("keep.com", "")}
		plan := state
		plan.Destinations = []DoNotDecryptListDestinations{dndDestination("keep.com", ""), dndDestination("new.com", "")}
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 3 {
			t.Fatalf("got %d requests, want 3: %+v", len(steps), steps)
		}
		if steps[0].viaList || gjson.Get(steps[0].body, "addDestinations.#.destination").Raw != `["new.com"]` {
			t.Errorf("request 0 must add new.com: %+v", steps[0])
		}
		if steps[1].viaList || gjson.Get(steps[1].body, "removeDestinations.#.destination").Raw != `["a.com"]` {
			t.Errorf("request 1 must remove the domain with the update request: %+v", steps[1])
		}
		if !steps[2].viaList || steps[2].body != "" || len(steps[2].removeDestinations) != 2 ||
			steps[2].removeDestinations[0].Destination.ValueString() != "10.10.10.10" || steps[2].removeDestinations[1].Destination.ValueString() != "192.0.2.0/24" {
			t.Errorf("request 2 must remove both IP destinations through the destination list: %+v", steps[2])
		}
	})

	t.Run("a removed comment is removed and added again last", func(t *testing.T) {
		state := state
		state.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", "x"), dndDestination("10.10.10.10", "y")}
		plan := state
		plan.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", ""), dndDestination("10.10.10.10", "")}
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 3 {
			t.Fatalf("got %d requests, want 3: %+v", len(steps), steps)
		}
		if steps[0].viaList || gjson.Get(steps[0].body, "removeDestinations.#.destination").Raw != `["a.com"]` {
			t.Errorf("request 0 must remove a.com: %+v", steps[0])
		}
		if !steps[1].viaList || len(steps[1].removeDestinations) != 1 || steps[1].removeDestinations[0].Destination.ValueString() != "10.10.10.10" {
			t.Errorf("request 1 must remove 10.10.10.10 through the destination list: %+v", steps[1])
		}
		if steps[2].viaList || gjson.Get(steps[2].body, "addDestinations.#.destination").Raw != `["a.com","10.10.10.10"]` || gjson.Get(steps[2].body, "addDestinations.0.comment").Exists() {
			t.Errorf("request 2 must add both again without a comment: %+v", steps[2])
		}
	})

	t.Run("a new name is sent even if the only other change is an IP removal", func(t *testing.T) {
		state := state
		state.Destinations = []DoNotDecryptListDestinations{dndDestination("10.10.10.10", "")}
		plan := state
		plan.Name = types.StringValue("new")
		plan.Destinations = nil
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 2 || steps[0].viaList || steps[0].body != `{"name":"new"}` || !steps[1].viaList {
			t.Fatalf("got %+v", steps)
		}
	})

	t.Run("null destinations in the plan remove all", func(t *testing.T) {
		state := state
		state.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", ""), dndDestination("b.com", ""), dndDestination("10.10.10.10", "")}
		plan := state
		plan.Destinations = nil
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 2 || len(gjson.Get(steps[0].body, "removeDestinations").Array()) != 2 || !steps[1].viaList || len(steps[1].removeDestinations) != 1 {
			t.Fatalf("got %+v", steps)
		}
	})

	t.Run("removals are split into chunks", func(t *testing.T) {
		state := state
		state.Destinations = dndDestinations(501)
		for i := 0; i < 500; i++ {
			state.Destinations = append(state.Destinations, dndDestination(fmt.Sprintf("10.1.%d.%d", i/200, i%200), ""))
		}
		plan := state
		plan.Destinations = nil
		steps := plan.toUpdateSteps(ctx, state)
		var put, list []int
		for _, step := range steps {
			if step.viaList {
				list = append(list, len(step.removeDestinations))
			} else {
				put = append(put, len(gjson.Get(step.body, "removeDestinations").Array()))
			}
		}
		if !slices.Equal(put, []int{500, 1}) || !slices.Equal(list, []int{doNotDecryptListListRemovalChunk, 500 - doNotDecryptListListRemovalChunk}) {
			t.Errorf("PUT removals = %v, destination list removals = %v", put, list)
		}
	})
}

func TestDoNotDecryptListApplyStep(t *testing.T) {
	ctx := context.Background()
	state := DoNotDecryptList{
		Name:                       types.StringValue("old"),
		Destinations:               []DoNotDecryptListDestinations{dndDestination("a.com", "x"), dndDestination("b.com", "y"), dndDestination("10.10.10.10", "")},
		Applications:               dndInt64Set(1, 2),
		DecryptExceptionCategories: dndInt64Set(10),
	}
	plan := state
	plan.Name = types.StringValue("new")
	plan.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", "changed"), dndDestination("c.com", "")}
	plan.Applications = dndInt64Set(2, 3)
	plan.DecryptExceptionCategories = types.SetNull(types.Int64Type)

	steps := plan.toUpdateSteps(ctx, state) // add, remove (b.com), remove through the list (10.10.10.10)
	if len(steps) != 3 {
		t.Fatalf("got %d requests, want 3: %+v", len(steps), steps)
	}
	original := slices.Clone(state.Destinations)
	interim := state
	interim.applyStep(ctx, plan.Name, steps[0]) // only the addition request was accepted

	if interim.Name.ValueString() != "new" {
		t.Errorf("name = %q", interim.Name.ValueString())
	}
	var got []string
	for _, item := range interim.Destinations {
		got = append(got, item.Destination.ValueString()+"="+item.Comment.ValueString())
	}
	if want := []string{"b.com=y", "10.10.10.10=", "a.com=changed", "c.com="}; !slices.Equal(got, want) {
		t.Errorf("after the addition request destinations = %v, want %v (an existing destination is updated, not duplicated)", got, want)
	}
	if got := doNotDecryptListInt64Values(ctx, interim.Applications); !slices.Equal(got, []int64{1, 2, 3}) {
		t.Errorf("applications = %v, want [1 2 3]", got)
	}
	if !slices.EqualFunc(original, state.Destinations, func(a, b DoNotDecryptListDestinations) bool { return a == b }) {
		t.Errorf("the state the interim copy was made from must not be modified")
	}

	interim.applyStep(ctx, plan.Name, steps[1]) // the removal request
	if got := doNotDecryptListInt64Values(ctx, interim.Applications); !slices.Equal(got, []int64{2, 3}) {
		t.Errorf("applications = %v, want [2 3]", got)
	}
	if !interim.DecryptExceptionCategories.IsNull() {
		t.Errorf("categories = %v, want null", interim.DecryptExceptionCategories)
	}

	interim.Name = types.StringValue("not touched by the destination list")
	interim.applyStep(ctx, plan.Name, steps[2]) // removal through the destination list does not carry the name
	if interim.Name.ValueString() != "not touched by the destination list" {
		t.Errorf("a removal through the destination list must not change the name, got %q", interim.Name.ValueString())
	}
	got = nil
	for _, item := range interim.Destinations {
		got = append(got, item.Destination.ValueString())
	}
	if want := []string{"a.com", "c.com"}; !slices.Equal(got, want) {
		t.Errorf("destinations = %v, want %v", got, want)
	}
}

func TestDoNotDecryptListDestinationProblem(t *testing.T) {
	for destination, wantProblem := range map[string]bool{
		"example.com":          false,
		"sub.example.com":      false,
		"10.10.10.10":          false,
		"192.0.2.0/24":         false,
		"https://foo.bar/path": true,
		"foo.bar/path":         true,
		"ftp://10.0.0.1":       true,
		"Example.com":          true,
		"EXAMPLE.COM":          true,
		"":                     true,
		" example.com":         true,
		"example.com ":         true,
		"example.com.":         true,
	} {
		if got := doNotDecryptListDestinationProblem(destination) != ""; got != wantProblem {
			t.Errorf("%q: problem = %q, want a problem: %v", destination, doNotDecryptListDestinationProblem(destination), wantProblem)
		}
	}
}

func TestDoNotDecryptListIsIP(t *testing.T) {
	for destination, want := range map[string]bool{
		"10.10.10.10":     true,
		"203.0.113.5":     true,
		"192.0.2.0/24":    true,
		"example.com":     false,
		"1.2.3":           false,
		"a.b.c.d":         false,
		"sub.example.com": false,
	} {
		if got := doNotDecryptListIsIP(destination); got != want {
			t.Errorf("%q: IP = %v, want %v", destination, got, want)
		}
	}
}

func TestDoNotDecryptListToUpdateStepsEdgeCases(t *testing.T) {
	ctx := context.Background()
	state := DoNotDecryptList{
		Name:                       types.StringValue("old"),
		Destinations:               []DoNotDecryptListDestinations{dndDestination("a.com", "x"), dndDestination("10.10.10.10", "y")},
		Applications:               dndInt64Set(1),
		DecryptExceptionCategories: dndInt64Set(10),
	}

	t.Run("no change sends a single request with the name", func(t *testing.T) {
		steps := state.toUpdateSteps(ctx, state)
		if len(steps) != 1 || steps[0].viaList || steps[0].body != `{"name":"old"}` {
			t.Fatalf("got %+v", steps)
		}
	})

	t.Run("a category and an application are added in the first request", func(t *testing.T) {
		plan := state
		plan.DecryptExceptionCategories = dndInt64Set(10, 11)
		plan.Applications = dndInt64Set(1, 2)
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 1 || gjson.Get(steps[0].body, "addDecryptExceptionCategories").Raw != `[11]` || gjson.Get(steps[0].body, "addApplications").Raw != `[2]` {
			t.Fatalf("got %+v", steps)
		}
	})

	t.Run("removing only an IP address with an unchanged name is a single request through the destination list", func(t *testing.T) {
		plan := state
		plan.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", "x")}
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 1 || !steps[0].viaList || steps[0].body != "" {
			t.Fatalf("got %+v", steps)
		}
	})

	t.Run("a changed name is sent before the first removal through the destination list", func(t *testing.T) {
		plan := state
		plan.Name = types.StringValue("new")
		plan.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", "x"), dndDestination("10.10.10.10", "")} // comment removed: list removal, then add
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 3 || steps[0].viaList || steps[0].body != `{"name":"new"}` || !steps[1].viaList || steps[2].viaList {
			t.Fatalf("got %+v", steps)
		}
	})
}

func TestDoNotDecryptListMissingDestination(t *testing.T) {
	for message, want := range map[string]string{
		`HTTP Request failed: StatusCode 400, {"statusCode":400,"error":"Bad Request","message":"destination 'gone.example.com' not found","txId":"x"}`: "gone.example.com",
		`destination '10.10.10.10' not found`: "10.10.10.10",
		`HTTP Request failed: StatusCode 400, {"message":"addDestinations cannot contain more than 500 destinations."}`: "",
		"": "",
	} {
		if got := doNotDecryptListMissingDestination(message); got != want {
			t.Errorf("%q: got %q, want %q", message, got, want)
		}
	}

	ctx := context.Background()
	state := DoNotDecryptList{Name: types.StringValue("n"), Destinations: []DoNotDecryptListDestinations{dndDestination("a.com", ""), dndDestination("b.com", ""), dndDestination("c.com", "")}}
	plan := DoNotDecryptList{Name: types.StringValue("n")}
	steps := plan.toUpdateSteps(ctx, state)
	if len(steps) != 1 || !steps[0].removes("b.com") || steps[0].removes("d.com") {
		t.Fatalf("got %+v", steps)
	}
	if got := gjson.Get(steps[0].bodyWithout(map[string]struct{}{"b.com": {}}), "removeDestinations.#.destination").Raw; got != `["a.com","c.com"]` {
		t.Errorf("without b.com: %s", got)
	}
	body := steps[0].bodyWithout(map[string]struct{}{"a.com": {}, "b.com": {}, "c.com": {}})
	if body != `{"name":"n"}` {
		t.Errorf("without all destinations the request only carries the name, got %s", body)
	}
}
