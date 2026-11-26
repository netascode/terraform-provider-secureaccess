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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
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
			input := dndDestinations(tc.destinations)
			steps := doNotDecryptListChangeSteps("name", add, input, applications, applications) // the same ids as categories
			// applications are only added, the API cannot remove them: a removal without destinations is no request at all
			wantSteps := tc.wantSteps
			if !add && tc.destinations == 0 {
				wantSteps = nil
			}
			if len(steps) != len(wantSteps) {
				t.Fatalf("destinations=%d applications=%d add=%v: got %d requests, want %d", tc.destinations, tc.applications, add, len(steps), len(wantSteps))
			}
			verb := map[bool]string{true: "add", false: "remove"}[add]
			var sent []string
			for _, step := range steps {
				for _, destination := range gjson.Get(step.body, verb+"Destinations.#.destination").Array() {
					sent = append(sent, destination.String())
				}
			}
			var want []string
			for _, item := range input {
				want = append(want, item.Destination.ValueString())
			}
			if !slices.Equal(sent, want) {
				t.Errorf("destinations=%d add=%v: the requests do not carry every destination exactly once, in order", tc.destinations, add)
			}
			for i, step := range steps {
				if got := len(gjson.Get(step.body, verb+"Destinations").Array()); got != wantSteps[i] {
					t.Errorf("destinations=%d add=%v request %d: got %d destinations, want %d", tc.destinations, add, i, got, wantSteps[i])
				}
				if name := gjson.Get(step.body, "name").String(); name != "name" {
					t.Errorf("request %d: name = %q, every request must carry the name", i, name)
				}
				for _, field := range []string{verb + "Applications", verb + "DecryptExceptionCategories"} {
					if got, want := gjson.Get(step.body, field).Exists(), add && i == 0 && tc.applications > 0; got != want {
						t.Errorf("destinations=%d add=%v request %d: %s present = %v, want %v (only added, and only in the first request)", tc.destinations, add, i, field, got, want)
					}
				}
			}
		}
	}
}

func TestDoNotDecryptListDestinationChanges(t *testing.T) {
	state := []DoNotDecryptListDestinations{dndDestination("same.com", "x"), dndDestination("changed.com", "x"), dndDestination("cleared.com", "x"), dndDestination("gets-comment.com", ""), dndDestination("gone.com", "x")}
	plan := []DoNotDecryptListDestinations{dndDestination("same.com", "x"), dndDestination("changed.com", "y"), dndDestination("cleared.com", ""), dndDestination("gets-comment.com", "z"), dndDestination("new.com", "")}

	upsert, gone, cleared := doNotDecryptListDestinationChanges(state, plan)
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
	if got, want := names(gone), []string{"gone.com"}; !slices.Equal(got, want) {
		t.Errorf("gone = %v, want %v", got, want)
	}
	if got, want := names(cleared), []string{"cleared.com"}; !slices.Equal(got, want) {
		t.Errorf("cleared = %v, want %v (adding without a comment does not clear an existing one)", got, want)
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

	t.Run("additions come first and update comments, removals carry no comment, applications and categories are never removed", func(t *testing.T) {
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
		// the plan drops application 1 and all categories, which the API cannot do: Update() refuses such a plan before it
		// builds any request, so a request must never carry the fields that the API rejects as "not an allowed field"
		for _, step := range steps {
			for _, field := range []string{"removeApplications", "removeDecryptExceptionCategories"} {
				if gjson.Get(step.body, field).Exists() {
					t.Errorf("a request carries %s, which the API rejects: %s", field, step.body)
				}
			}
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

	t.Run("destinations that are gone are removed before those whose comment is removed", func(t *testing.T) {
		state := state
		state.Destinations = []DoNotDecryptListDestinations{dndDestination("gone.com", ""), dndDestination("10.0.0.1", ""), dndDestination("cleared.com", "x"), dndDestination("10.0.0.2", "y")}
		plan := state
		plan.Destinations = []DoNotDecryptListDestinations{dndDestination("cleared.com", ""), dndDestination("10.0.0.2", "")}
		steps := plan.toUpdateSteps(ctx, state)
		if len(steps) != 5 {
			t.Fatalf("got %d requests, want 5: %+v", len(steps), steps)
		}
		removed := func(step doNotDecryptListStep) string {
			var out []string
			for _, item := range step.removeDestinations {
				out = append(out, item.Destination.ValueString())
			}
			return fmt.Sprintf("%v/list=%v", out, step.viaList)
		}
		for i, want := range []string{"[gone.com]/list=false", "[10.0.0.1]/list=true", "[cleared.com]/list=false", "[10.0.0.2]/list=true"} {
			if got := removed(steps[i]); got != want {
				t.Errorf("request %d removes %s, want %s", i, got, want)
			}
		}
		if got := gjson.Get(steps[4].body, "addDestinations.#.destination").Raw; got != `["cleared.com","10.0.0.2"]` || steps[4].viaList {
			t.Errorf("the last request must add the cleared destinations again: %+v", steps[4])
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
		var putSent, listSent []string
		for _, step := range steps {
			if step.viaList {
				for _, item := range step.removeDestinations {
					listSent = append(listSent, item.Destination.ValueString())
				}
			} else {
				for _, destination := range gjson.Get(step.body, "removeDestinations.#.destination").Array() {
					putSent = append(putSent, destination.String())
				}
			}
		}
		if len(putSent) != 501 || putSent[0] != "host0.example.com" || putSent[500] != "host500.example.com" || len(listSent) != 500 || listSent[0] != "10.1.0.0" || listSent[499] != "10.1.2.99" {
			t.Errorf("removals do not carry every destination exactly once, in order: %d domains (%v...%v), %d IPs (%v...%v)", len(putSent), putSent[:1], putSent[len(putSent)-1:], len(listSent), listSent[:1], listSent[len(listSent)-1:])
		}
		if !slices.Equal(put, []int{500, 1}) || !slices.Equal(list, []int{destinationListMaxDestinationsPerRequest, 500 - destinationListMaxDestinationsPerRequest}) {
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
	plan.Applications = dndInt64Set(1, 2, 3) // applications and categories are only added, never removed
	plan.DecryptExceptionCategories = dndInt64Set(10, 11)

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
	if got := helpers.Int64SetElements(ctx, interim.Applications); !slices.Equal(got, []int64{1, 2, 3}) {
		t.Errorf("applications = %v, want [1 2 3]", got)
	}
	if got := helpers.Int64SetElements(ctx, interim.DecryptExceptionCategories); !slices.Equal(got, []int64{10, 11}) {
		t.Errorf("categories = %v, want [10 11]", got)
	}
	if !slices.EqualFunc(original, state.Destinations, func(a, b DoNotDecryptListDestinations) bool { return a == b }) {
		t.Errorf("the state the interim copy was made from must not be modified")
	}

	interim.applyStep(ctx, plan.Name, steps[1]) // the removal request removes a destination, nothing else
	if got := helpers.Int64SetElements(ctx, interim.Applications); !slices.Equal(got, []int64{1, 2, 3}) {
		t.Errorf("applications = %v, want [1 2 3]", got)
	}
	if got := helpers.Int64SetElements(ctx, interim.DecryptExceptionCategories); !slices.Equal(got, []int64{10, 11}) {
		t.Errorf("categories = %v, want [10 11]", got)
	}

	// a step that carries no applications or categories leaves null sets null (the state after an import)
	imported := DoNotDecryptList{Applications: types.SetNull(types.Int64Type), DecryptExceptionCategories: types.SetNull(types.Int64Type)}
	imported.applyStep(ctx, plan.Name, steps[1])
	if !imported.Applications.IsNull() || !imported.DecryptExceptionCategories.IsNull() {
		t.Errorf("null sets must stay null, got %v and %v", imported.Applications, imported.DecryptExceptionCategories)
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

// The values of the destinations are checked by the validators that the generator creates from `string_patterns`.
func TestDoNotDecryptListDestinationValidators(t *testing.T) {
	ctx := context.Background()
	var schemaResp resource.SchemaResponse
	NewDoNotDecryptListResource().Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	destinations, ok := schemaResp.Schema.Attributes["destinations"].(schema.SetNestedAttribute)
	if !ok {
		t.Fatalf("destinations is not a nested set attribute")
	}
	destination, ok := destinations.NestedObject.Attributes["destination"].(schema.StringAttribute)
	if !ok || len(destination.Validators) == 0 {
		t.Fatalf("the destination attribute has no validators, string_patterns of the definition are missing")
	}

	for value, valid := range map[string]bool{
		"example.com":           true,
		"sub.example.com":       true,
		"10.10.10.10":           true,
		"203.0.113.5":           true,
		"192.0.2.0/24":          true,
		"10.0.0.0/8":            true,
		"1.2.3.4/32":            true,
		"xn--bcher-kva.example": true,
		"my_host.example.com":   true,
		"*.example.com":         true,
		"https://foo.bar/path":  false,
		"foo.bar/path":          false,
		"ftp://10.0.0.1":        false,
		"Example.com":           false,
		"EXAMPLE.COM":           false,
		"":                      false,
		" example.com":          false,
		"example.com ":          false,
		"example.com.":          false,
		"10.10.10.010":          false,
		"127.1":                 false,
		"256.1.1.1":             false,
		"1.2.3.4/33":            false,
		"1.2.3.4/032":           false,
		"1.2.3.249":             true,
		"1.2.3.256":             false,
		"01.2.3.4":              false,
		"1234567890":            false,
		"exaMple.com":           false,
		"example.com\n":         false,
		"example.com\t":         false,
		"example.com\u00a0":     false,
		"example\u2003.com":     false,
		"É.com":                 false,
		"ЯНДЕКС.рф":             false,
		"bücher.example":        true,
	} {
		var diags diag.Diagnostics
		for _, v := range destination.Validators {
			resp := validator.StringResponse{}
			v.ValidateString(ctx, validator.StringRequest{Path: path.Root("destination"), ConfigValue: types.StringValue(value)}, &resp)
			diags.Append(resp.Diagnostics...)
		}
		if got := !diags.HasError(); got != valid {
			t.Errorf("%q: accepted = %v, want %v", value, got, valid)
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

// The request body of create is generated, adjustBody adds what the API requires on top of it.
func TestDoNotDecryptListCreateBody(t *testing.T) {
	ctx := context.Background()
	body := func(plan DoNotDecryptList) string {
		return plan.adjustBody(ctx, plan.toBody(ctx, DoNotDecryptList{}))
	}

	t.Run("only a name: destinations are sent as an empty list", func(t *testing.T) {
		plan := DoNotDecryptList{
			Name:                       types.StringValue("n"),
			Applications:               types.SetNull(types.Int64Type),
			DecryptExceptionCategories: types.SetNull(types.Int64Type),
		}
		if got := body(plan); got != `{"name":"n","destinations":[]}` {
			t.Errorf("got %s", got)
		}
	})

	t.Run("everything set", func(t *testing.T) {
		plan := DoNotDecryptList{
			Name:                       types.StringValue("n"),
			Destinations:               []DoNotDecryptListDestinations{dndDestination("a.com", "x"), dndDestination("10.10.10.10", "")},
			Applications:               dndInt64Set(101),
			DecryptExceptionCategories: dndInt64Set(12),
		}
		got := body(plan)
		if gjson.Get(got, "destinations.#.destination").Raw != `["a.com","10.10.10.10"]` ||
			gjson.Get(got, "destinations.0.comment").String() != "x" || gjson.Get(got, "destinations.1.comment").Exists() ||
			gjson.Get(got, "applications").Raw != `[101]` || gjson.Get(got, "decryptExceptionCategories").Raw != `[12]` {
			t.Errorf("got %s", got)
		}
	})
}

// The destinations are read from the destination list behind the object (Read() puts them into the body as `destinations`),
// applications and categories are not returned in the form that is written, so they are write only: the generated
// fromBody and fromBodyPartial never touch them.
func TestDoNotDecryptListFromBodyWithDestinations(t *testing.T) {
	ctx := context.Background()
	// the body that Read() builds: the object, with the destinations of the destination list behind it as `destinations`
	response := gjson.Parse(`{"name":"n","isDefault":false,"decryptExceptionCategories":[{"categoryId":12,"name":"A"}],"totalDomains":2,"totalIps":0,"totalApplications":0,` +
		`"destinations":[{"id":1,"destination":"a.com","comment":"x","type":"domain"},{"id":2,"destination":"c.com","comment":null,"type":"domain"}]}`)
	destinationNames := func(items []DoNotDecryptListDestinations) []string {
		var out []string
		for _, item := range items {
			out = append(out, item.Destination.ValueString())
		}
		return out
	}

	t.Run("a refresh reads the destinations of the state, applications and categories stay as they are", func(t *testing.T) {
		state := DoNotDecryptList{
			Name:                       types.StringValue("n"),
			Destinations:               []DoNotDecryptListDestinations{dndDestination("a.com", "old"), dndDestination("b.com", "y"), dndDestination("c.com", "")},
			Applications:               dndInt64Set(101),
			DecryptExceptionCategories: dndInt64Set(12, 27),
			IsDefault:                  types.BoolValue(false),
		}
		state.fromBodyPartial(ctx, response)
		// a.com is read (the comment changed outside of Terraform), b.com is gone from the destination list and is dropped,
		// c.com has no comment in the state and keeps none
		if got, want := destinationNames(state.Destinations), []string{"a.com", "c.com"}; !slices.Equal(got, want) {
			t.Fatalf("destinations = %v, want %v (a destination that was removed outside of Terraform is dropped)", got, want)
		}
		if got := state.Destinations[0].Comment.ValueString(); got != "x" {
			t.Errorf("the comment of a.com = %q, want the one of the destination list", got)
		}
		if !state.Destinations[1].Comment.IsNull() {
			t.Errorf("a destination without a comment in the state must keep none, got %v", state.Destinations[1].Comment)
		}
		if !slices.Equal(helpers.Int64SetElements(ctx, state.Applications), []int64{101}) ||
			!slices.Equal(helpers.Int64SetElements(ctx, state.DecryptExceptionCategories), []int64{12, 27}) {
			t.Errorf("write only attributes changed by a refresh: %+v", state)
		}
		if state.TotalDomains.ValueInt64() != 2 {
			t.Errorf("the totals are read: %v", state.TotalDomains)
		}
	})

	t.Run("a destination that was added outside of Terraform is not read by a refresh", func(t *testing.T) {
		state := DoNotDecryptList{Name: types.StringValue("n"), Destinations: []DoNotDecryptListDestinations{dndDestination("a.com", "x")}}
		state.fromBodyPartial(ctx, response)
		if got, want := destinationNames(state.Destinations), []string{"a.com"}; !slices.Equal(got, want) {
			t.Errorf("destinations = %v, want %v (only the destinations of the state are managed)", got, want)
		}
	})

	t.Run("an import reads all destinations and leaves applications and categories empty", func(t *testing.T) {
		var imported DoNotDecryptList
		imported.fromBody(ctx, response)
		if got, want := destinationNames(imported.Destinations), []string{"a.com", "c.com"}; !slices.Equal(got, want) {
			t.Fatalf("destinations = %v, want %v", got, want)
		}
		if got := imported.Destinations[0].Comment.ValueString(); got != "x" {
			t.Errorf("the comment of a.com = %q", got)
		}
		if !imported.Applications.IsNull() || !imported.DecryptExceptionCategories.IsNull() {
			t.Errorf("applications and categories are not read: %+v", imported)
		}
		if imported.Name.ValueString() != "n" || imported.TotalDomains.ValueInt64() != 2 {
			t.Errorf("name and totals must be read: %+v", imported)
		}
	})
}

// dndReadAPI is a fake API for the tests of the read functions. It serves the Do Not Decrypt Lists API and the
// destinations of destination lists, and keeps the paths of the requests it got.
type dndReadAPI struct {
	mu       sync.Mutex
	paths    []string
	listHas  bool // the items of the list carry exceptiondomainlistId
	destCode int  // status of the answer for the destinations, 0 is 200
}

func (a *dndReadAPI) requests() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return slices.Clone(a.paths)
}

func newDndReadClient(t *testing.T, listHas bool, destCode int) (*secureaccess.Client, *dndReadAPI) {
	t.Helper()
	api := &dndReadAPI{listHas: listHas, destCode: destCode}
	object := func(withList bool) map[string]any {
		o := map[string]any{"id": 5, "name": "n", "isDefault": false, "totalDomains": 2, "totalIps": 0, "totalApplications": 0}
		if withList {
			o["exceptiondomainlistId"] = 777
		}
		return o
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/v2/token" {
			api.mu.Lock()
			api.paths = append(api.paths, r.URL.Path)
			api.mu.Unlock()
		}
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/auth/v2/token":
			_, _ = io.WriteString(w, `{"access_token":"t","token_type":"bearer","expires_in":3600}`)
		case r.URL.Path == "/policies/v2/doNotDecryptLists":
			b, _ := json.Marshal([]any{object(api.listHas)})
			_, _ = w.Write(b)
		case r.URL.Path == "/policies/v2/doNotDecryptLists/5":
			b, _ := json.Marshal(object(true))
			_, _ = w.Write(b)
		case r.URL.Path == "/policies/v2/destinationlists/777/destinations":
			if api.destCode != 0 {
				w.WriteHeader(api.destCode)
				_, _ = io.WriteString(w, `{"message":"no"}`)
				return
			}
			_, _ = io.WriteString(w, `{"data":[{"id":1,"destination":"a.com","comment":"x","type":"domain"},{"id":2,"destination":"c.com","type":"domain"}],"meta":{"total":2}}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	client, err := secureaccess.NewClient(server.URL, "key", "secret")
	if err != nil {
		t.Fatal(err)
	}
	return client, api
}

func TestDoNotDecryptListWithDestinations(t *testing.T) {
	ctx := context.Background()

	t.Run("the destinations of the destination list are added to the object", func(t *testing.T) {
		client, api := newDndReadClient(t, true, 0)
		object := gjson.Parse(`{"id":5,"name":"n","exceptiondomainlistId":777,"totalDomains":2}`)
		body, err := doNotDecryptListWithDestinations(ctx, client, object)
		if err != nil {
			t.Fatal(err)
		}
		if got := body.Get("destinations.#.destination").Raw; got != `["a.com","c.com"]` {
			t.Errorf("destinations = %s in %s", got, body.Raw)
		}
		if body.Get("name").String() != "n" || body.Get("totalDomains").Int() != 2 || body.Get("exceptiondomainlistId").Int() != 777 {
			t.Errorf("the rest of the object must stay: %s", body.Raw)
		}
		if got := api.requests(); !slices.Equal(got, []string{"/policies/v2/destinationlists/777/destinations"}) {
			t.Errorf("requests = %v", got)
		}
	})

	t.Run("an object without a destination list is an error and nothing is requested", func(t *testing.T) {
		client, api := newDndReadClient(t, true, 0)
		for _, object := range []string{`{"id":5,"name":"n"}`, `{"id":5,"exceptiondomainlistId":0}`} {
			if _, err := doNotDecryptListWithDestinations(ctx, client, gjson.Parse(object)); err == nil || !strings.Contains(err.Error(), "exceptiondomainlistId") {
				t.Errorf("%s: err = %v, want one that names exceptiondomainlistId", object, err)
			}
		}
		if got := api.requests(); len(got) != 0 {
			t.Errorf("requests = %v, want none", got)
		}
	})

	t.Run("a failing destination request is an error", func(t *testing.T) {
		client, _ := newDndReadClient(t, true, http.StatusBadRequest)
		if _, err := doNotDecryptListWithDestinations(ctx, client, gjson.Parse(`{"id":5,"exceptiondomainlistId":777}`)); err == nil || !strings.Contains(err.Error(), "failed to retrieve destinations") {
			t.Errorf("err = %v", err)
		}
	})
}

// The data source reads the destinations of the destination list behind the object that it found, however it found it.
func TestDoNotDecryptListDataSourceRead(t *testing.T) {
	ctx := context.Background()
	read := func(t *testing.T, client *secureaccess.Client, id, name types.String) (datasource.ReadResponse, DoNotDecryptList) {
		t.Helper()
		d := &DoNotDecryptListDataSource{client: client}
		var schemaResp datasource.SchemaResponse
		d.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
		config := DoNotDecryptList{
			Id: id, Name: name,
			Applications:               types.SetNull(types.Int64Type),
			DecryptExceptionCategories: types.SetNull(types.Int64Type),
		}
		nullRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)
		configState := tfsdk.State{Schema: schemaResp.Schema, Raw: nullRaw}
		if diags := configState.Set(ctx, &config); diags.HasError() {
			t.Fatal(diags)
		}
		resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: nullRaw}}
		d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: configState.Raw}}, &resp)
		var out DoNotDecryptList
		if !resp.Diagnostics.HasError() {
			if diags := resp.State.Get(ctx, &out); diags.HasError() {
				t.Fatal(diags)
			}
		}
		return resp, out
	}
	names := func(items []DoNotDecryptListDestinations) []string {
		var out []string
		for _, item := range items {
			out = append(out, item.Destination.ValueString())
		}
		return out
	}

	for name, tc := range map[string]struct {
		id, name types.String
		listHas  bool
		want     []string
	}{
		"by name, the list carries the destination list": {types.StringNull(), types.StringValue("n"), true,
			[]string{"/policies/v2/doNotDecryptLists", "/policies/v2/destinationlists/777/destinations"}},
		"by name, the list lacks the destination list": {types.StringNull(), types.StringValue("n"), false,
			[]string{"/policies/v2/doNotDecryptLists", "/policies/v2/doNotDecryptLists/5", "/policies/v2/destinationlists/777/destinations"}},
		"by id": {types.StringValue("5"), types.StringNull(), true,
			[]string{"/policies/v2/doNotDecryptLists/5", "/policies/v2/destinationlists/777/destinations"}},
	} {
		t.Run(name, func(t *testing.T) {
			client, api := newDndReadClient(t, tc.listHas, 0)
			resp, got := read(t, client, tc.id, tc.name)
			if resp.Diagnostics.HasError() {
				t.Fatal(resp.Diagnostics)
			}
			if got.Id.ValueString() != "5" || got.Name.ValueString() != "n" || got.TotalDomains.ValueInt64() != 2 {
				t.Errorf("id, name and totals: %+v", got)
			}
			if want := []string{"a.com", "c.com"}; !slices.Equal(names(got.Destinations), want) {
				t.Fatalf("destinations = %v, want %v", names(got.Destinations), want)
			}
			if got.Destinations[0].Id.ValueInt64() != 1 || got.Destinations[1].Id.ValueInt64() != 2 {
				t.Errorf("the ids of the destinations = %v, %v, want 1 and 2 (the ids of the destination list)", got.Destinations[0].Id, got.Destinations[1].Id)
			}
			if got.Destinations[0].Comment.ValueString() != "x" {
				t.Errorf("the comment of a.com = %q", got.Destinations[0].Comment.ValueString())
			}
			if !got.Applications.IsNull() || !got.DecryptExceptionCategories.IsNull() {
				t.Errorf("applications and categories are not read: %+v", got)
			}
			if reqs := api.requests(); !slices.Equal(reqs, tc.want) {
				t.Errorf("requests = %v, want %v", reqs, tc.want)
			}
		})
	}

	t.Run("a failing destination request is an error", func(t *testing.T) {
		client, _ := newDndReadClient(t, true, http.StatusBadRequest)
		resp, _ := read(t, client, types.StringValue("5"), types.StringNull())
		if !resp.Diagnostics.HasError() || !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), "failed to retrieve destinations") {
			t.Errorf("diagnostics = %v", resp.Diagnostics)
		}
	})
}

func TestDoNotDecryptListFirstDuplicate(t *testing.T) {
	for name, tc := range map[string]struct {
		in   []string
		want string
	}{
		"none":                        {[]string{"a.com", "b.com", "10.10.10.10"}, ""},
		"empty":                       {nil, ""},
		"the same one twice":          {[]string{"a.com", "b.com", "a.com"}, "a.com"},
		"an IP address twice":         {[]string{"10.10.10.10", "a.com", "10.10.10.10"}, "10.10.10.10"},
		"a different case is another": {[]string{"a.com", "A.com"}, ""},
	} {
		got, found := doNotDecryptListFirstDuplicate(tc.in)
		if got != tc.want || found != (tc.want != "") {
			t.Errorf("%s: got %q, %v, want %q", name, got, found, tc.want)
		}
	}
}

func TestDoNotDecryptListAddedAndRemoved(t *testing.T) {
	ctx := context.Background()
	state, plan := dndInt64Set(1, 2, 3), dndInt64Set(2, 3, 4)
	if got := doNotDecryptListRemoved(ctx, state, plan); !slices.Equal(got, []int64{1}) {
		t.Errorf("removed = %v, want [1]", got)
	}
	if got := doNotDecryptListAdded(ctx, state, plan); !slices.Equal(got, []int64{4}) {
		t.Errorf("added = %v, want [4]", got)
	}
	// a set that is dropped from the configuration is null in the plan: everything in the state is removed
	if got := doNotDecryptListRemoved(ctx, state, types.SetNull(types.Int64Type)); !slices.Equal(got, []int64{1, 2, 3}) {
		t.Errorf("removed with a null plan = %v, want [1 2 3]", got)
	}
	if got := doNotDecryptListRemoved(ctx, types.SetNull(types.Int64Type), plan); len(got) != 0 {
		t.Errorf("nothing is removed from a null state (create, import), got %v", got)
	}
	if got := doNotDecryptListRemoved(ctx, state, state); len(got) != 0 {
		t.Errorf("an unchanged set removes nothing, got %v", got)
	}
}

func TestDoNotDecryptListRemovalDiagnostics(t *testing.T) {
	if diags := doNotDecryptListRemovalDiagnostics(nil, nil, false); len(diags) != 0 {
		t.Fatalf("nothing is removed, got %v", diags)
	}

	for _, asWarning := range []bool{false, true} {
		diags := doNotDecryptListRemovalDiagnostics([]int64{1}, []int64{72, 415}, asWarning)
		if len(diags) != 2 {
			t.Fatalf("asWarning=%v: got %d diagnostics, want one per attribute: %v", asWarning, len(diags), diags)
		}
		if diags.HasError() == asWarning {
			t.Errorf("asWarning=%v: HasError = %v (the plan must only warn, an error would block -replace)", asWarning, diags.HasError())
		}
		for i, want := range []struct {
			attribute, summary, ids string
		}{
			{"applications", "Removing applications is not supported", "1"},
			{"decrypt_exception_categories", "Removing content categories is not supported", "72, 415"},
		} {
			withPath, ok := diags[i].(diag.DiagnosticWithPath)
			if !ok || !withPath.Path().Equal(path.Root(want.attribute)) {
				t.Errorf("asWarning=%v diagnostic %d: not attached to %s: %#v", asWarning, i, want.attribute, diags[i])
			}
			if diags[i].Summary() != want.summary {
				t.Errorf("asWarning=%v diagnostic %d: summary = %q, want %q", asWarning, i, diags[i].Summary(), want.summary)
			}
			for _, part := range []string{want.ids, "-replace", "new id"} {
				if !strings.Contains(diags[i].Detail(), part) {
					t.Errorf("asWarning=%v diagnostic %d: detail lacks %q: %s", asWarning, i, part, diags[i].Detail())
				}
			}
			// the warning is also shown on a plan that replaces the list, so it must not claim that the plan is an update
			// that fails; the error is only raised when nothing was sent
			if got := strings.Contains(diags[i].Detail(), "Unless the list is replaced, applying this plan fails."); got != asWarning {
				t.Errorf("asWarning=%v diagnostic %d: conditional failure sentence present = %v: %s", asWarning, i, got, diags[i].Detail())
			}
			if got := strings.Contains(diags[i].Detail(), "Nothing was changed."); got == asWarning {
				t.Errorf("asWarning=%v diagnostic %d: 'Nothing was changed.' present = %v: %s", asWarning, i, got, diags[i].Detail())
			}
		}
	}

	if diags := doNotDecryptListRemovalDiagnostics(nil, []int64{415}, false); len(diags) != 1 || diags[0].Summary() != "Removing content categories is not supported" {
		t.Errorf("only the categories are removed, got %v", diags)
	}
}

// dndTestSchema returns the schema of the resource, the plan and state of the tests below are built from it.
func dndTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	var schemaResp resource.SchemaResponse
	NewDoNotDecryptListResource().Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	return schemaResp.Schema
}

func dndTestModel(applications, categories types.Set) DoNotDecryptList {
	return DoNotDecryptList{
		Id:                         types.StringValue("1"),
		Name:                       types.StringValue("n"),
		Applications:               applications,
		DecryptExceptionCategories: categories,
		IsDefault:                  types.BoolValue(false),
		TotalDomains:               types.Int64Value(0),
		TotalIps:                   types.Int64Value(0),
		TotalApplications:          types.Int64Value(0),
	}
}

func dndTestState(t *testing.T, s schema.Schema, model *DoNotDecryptList) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(context.Background()), nil)}
	if model != nil {
		if diags := state.Set(context.Background(), model); diags.HasError() {
			t.Fatal(diags)
		}
	}
	return state
}

func dndTestPlan(t *testing.T, s schema.Schema, model *DoNotDecryptList) tfsdk.Plan {
	t.Helper()
	state := dndTestState(t, s, model)
	return tfsdk.Plan{Schema: s, Raw: state.Raw}
}

// dndTestAPI is a fake API for the tests that run Update(): it answers the token request, keeps the body of every PUT to
// the list with id 1, answers the read of the object that Update() makes at the end (the list has no destinations) and
// fails the test on any other request.
type dndTestAPI struct {
	mu   sync.Mutex
	puts []string
	gets int
}

func (a *dndTestAPI) bodies() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return slices.Clone(a.puts)
}

func dndTestResource(t *testing.T) (*DoNotDecryptListResource, *dndTestAPI) {
	t.Helper()
	api := &dndTestAPI{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/auth/v2/token":
			_, _ = io.WriteString(w, `{"access_token":"t","token_type":"bearer","expires_in":3600}`)
		case r.Method == http.MethodPut && r.URL.Path == "/policies/v2/doNotDecryptLists/1":
			body, _ := io.ReadAll(r.Body)
			api.mu.Lock()
			api.puts = append(api.puts, string(body))
			api.mu.Unlock()
			_, _ = io.WriteString(w, `{"id":1,"name":"n","isDefault":false,"totalDomains":0,"totalIps":0,"totalApplications":3}`)
		case r.Method == http.MethodGet && r.URL.Path == "/policies/v2/doNotDecryptLists/1":
			api.mu.Lock()
			api.gets++
			api.mu.Unlock()
			_, _ = io.WriteString(w, `{"id":1,"name":"n","isDefault":false,"totalDomains":0,"totalIps":0,"totalApplications":3}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "unexpected request", http.StatusInternalServerError)
		}
	}))
	t.Cleanup(server.Close)
	client, err := secureaccess.NewClient(server.URL, "key", "secret")
	if err != nil {
		t.Fatal(err)
	}
	return &DoNotDecryptListResource{client: client}, api
}

// Update() must refuse a plan that removes applications or categories before it sends anything: no request reaches the
// API and the state stays as it was. A plan that only adds goes through, which proves that the guard does not block
// every update.
func TestDoNotDecryptListUpdateGuard(t *testing.T) {
	ctx := context.Background()
	s := dndTestSchema(t)
	state := dndTestModel(dndInt64Set(1, 2), dndInt64Set(72, 415))
	null := types.SetNull(types.Int64Type)
	categories, applications := "Removing content categories is not supported", "Removing applications is not supported"

	update := func(t *testing.T, plan DoNotDecryptList) (resource.UpdateResponse, tfsdk.State, *dndTestAPI) {
		t.Helper()
		r, api := dndTestResource(t)
		stateModel := state
		prior := dndTestState(t, s, &stateModel)
		resp := resource.UpdateResponse{State: dndTestState(t, s, &stateModel)}
		r.Update(ctx, resource.UpdateRequest{Plan: dndTestPlan(t, s, &plan), State: prior}, &resp)
		return resp, prior, api
	}

	for name, tc := range map[string]struct {
		plan DoNotDecryptList
		want []string
	}{
		"a category is dropped":                 {dndTestModel(dndInt64Set(1, 2), dndInt64Set(72)), []string{categories}},
		"an application is dropped":             {dndTestModel(dndInt64Set(1), dndInt64Set(72, 415)), []string{applications}},
		"both sets are removed from the file":   {dndTestModel(null, null), []string{applications, categories}},
		"one is added while another is dropped": {dndTestModel(dndInt64Set(1, 2, 3), dndInt64Set(415)), []string{categories}},
	} {
		t.Run("refused: "+name, func(t *testing.T) {
			resp, prior, api := update(t, tc.plan)
			var got []string
			for _, d := range resp.Diagnostics.Errors() {
				got = append(got, d.Summary())
				if !strings.Contains(d.Detail(), "Nothing was changed.") {
					t.Errorf("the error must say that nothing was changed: %s", d.Detail())
				}
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("errors = %v, want %v", got, tc.want)
			}
			if puts := api.bodies(); len(puts) != 0 {
				t.Errorf("nothing may be sent, got %d request(s): %v", len(puts), puts)
			}
			if !resp.State.Raw.Equal(prior.Raw) {
				t.Errorf("the state must stay as it was, nothing was applied")
			}
		})
	}

	t.Run("applications and categories are added", func(t *testing.T) {
		resp, _, api := update(t, dndTestModel(dndInt64Set(1, 2, 3), dndInt64Set(72, 415, 9)))
		if resp.Diagnostics.HasError() {
			t.Fatalf("an update that only adds must go through: %v", resp.Diagnostics)
		}
		if warnings := resp.Diagnostics.Warnings(); len(warnings) > 0 {
			t.Errorf("the read after the update must succeed, got warnings: %v", warnings)
		}
		puts := api.bodies()
		if len(puts) != 1 {
			t.Fatalf("got %d requests, want 1: %v", len(puts), puts)
		}
		api.mu.Lock()
		gets := api.gets
		api.mu.Unlock()
		if gets != 1 {
			t.Errorf("got %d reads of the object after the update, want 1 (it gives the ids of the destinations)", gets)
		}
		if got := gjson.Get(puts[0], "addApplications").Raw; got != `[3]` {
			t.Errorf("addApplications = %s in %s", got, puts[0])
		}
		if got := gjson.Get(puts[0], "addDecryptExceptionCategories").Raw; got != `[9]` {
			t.Errorf("addDecryptExceptionCategories = %s in %s", got, puts[0])
		}
		for _, field := range []string{"removeApplications", "removeDecryptExceptionCategories"} {
			if gjson.Get(puts[0], field).Exists() {
				t.Errorf("the request carries %s, which the API rejects: %s", field, puts[0])
			}
		}
	})

	t.Run("a rename without other changes goes through", func(t *testing.T) {
		plan := dndTestModel(state.Applications, state.DecryptExceptionCategories)
		plan.Name = types.StringValue("renamed")
		resp, _, api := update(t, plan)
		if resp.Diagnostics.HasError() {
			t.Fatalf("a rename must go through: %v", resp.Diagnostics)
		}
		if puts := api.bodies(); len(puts) != 1 || puts[0] != `{"name":"renamed"}` {
			t.Errorf("requests = %v, want one with only the name", puts)
		}
	})
}

// ModifyPlan only warns, an error would also block `terraform apply -replace`, the way to get rid of categories.
func TestDoNotDecryptListModifyPlan(t *testing.T) {
	ctx := context.Background()
	s := dndTestSchema(t)
	stateModel := dndTestModel(dndInt64Set(1, 2), dndInt64Set(72, 415))
	state := dndTestState(t, s, &stateModel)
	plan := func(applications, categories types.Set) tfsdk.Plan {
		model := dndTestModel(applications, categories)
		return dndTestPlan(t, s, &model)
	}
	run := func(state tfsdk.State, plan tfsdk.Plan) diag.Diagnostics {
		resp := resource.ModifyPlanResponse{Plan: plan}
		NewDoNotDecryptListResource().(resource.ResourceWithModifyPlan).ModifyPlan(ctx, resource.ModifyPlanRequest{State: state, Plan: plan}, &resp)
		return resp.Diagnostics
	}
	summaries := func(diags diag.Diagnostics) []string {
		var out []string
		for _, d := range diags {
			out = append(out, d.Summary())
		}
		return out
	}
	null := types.SetNull(types.Int64Type)
	categories, applications := "Removing content categories is not supported", "Removing applications is not supported"

	t.Run("nothing is checked on create", func(t *testing.T) {
		if diags := run(dndTestState(t, s, nil), plan(dndInt64Set(1), dndInt64Set(72))); len(diags) != 0 {
			t.Errorf("got %v", diags)
		}
	})
	t.Run("nothing is checked on destroy", func(t *testing.T) {
		if diags := run(state, dndTestPlan(t, s, nil)); len(diags) != 0 {
			t.Errorf("got %v", diags)
		}
	})
	t.Run("unchanged and added values do not warn", func(t *testing.T) {
		if diags := run(state, plan(dndInt64Set(1, 2), dndInt64Set(72, 415))); len(diags) != 0 {
			t.Errorf("unchanged: got %v", diags)
		}
		if diags := run(state, plan(dndInt64Set(1, 2, 3), dndInt64Set(72, 415, 7))); len(diags) != 0 {
			t.Errorf("added: got %v", diags)
		}
	})
	t.Run("removals warn and never fail", func(t *testing.T) {
		diags := run(state, plan(dndInt64Set(1), dndInt64Set(72, 415)))
		if got := summaries(diags); !slices.Equal(got, []string{applications}) || diags.HasError() {
			t.Errorf("application dropped: got %v (HasError %v)", got, diags.HasError())
		}
		diags = run(state, plan(null, null))
		if got := summaries(diags); !slices.Equal(got, []string{applications, categories}) || diags.HasError() {
			t.Errorf("both dropped: got %v (HasError %v)", got, diags.HasError())
		}
	})
	t.Run("values that are not known yet are left to Update", func(t *testing.T) {
		unknownPlan := plan(dndInt64Set(1, 2), dndInt64Set(72, 415))
		if diags := unknownPlan.SetAttribute(ctx, path.Root("decrypt_exception_categories"), types.SetUnknown(types.Int64Type)); diags.HasError() {
			t.Fatal(diags)
		}
		if diags := run(state, unknownPlan); len(diags) != 0 {
			t.Errorf("got %v", diags)
		}
	})
	t.Run("a set with an unknown element is left to Update", func(t *testing.T) {
		// applications = [1, <id of another resource>]: the unknown element cannot be read, which must not look like
		// the removal of 1 and 2
		partlyUnknown := plan(dndInt64Set(1, 2), dndInt64Set(72, 415))
		set := types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1), types.Int64Unknown()})
		if diags := partlyUnknown.SetAttribute(ctx, path.Root("applications"), set); diags.HasError() {
			t.Fatal(diags)
		}
		if diags := run(state, partlyUnknown); len(diags) != 0 {
			t.Errorf("got %v", diags)
		}
	})
	t.Run("unknown destinations do not break the check", func(t *testing.T) {
		// the whole model cannot be read while destinations are unknown, which is why only the two sets are read
		unknownPlan := plan(dndInt64Set(1), dndInt64Set(72, 415))
		destinations := s.Attributes["destinations"].GetType().(types.SetType)
		if diags := unknownPlan.SetAttribute(ctx, path.Root("destinations"), types.SetUnknown(destinations.ElemType)); diags.HasError() {
			t.Fatal(diags)
		}
		diags := run(state, unknownPlan)
		if got := summaries(diags); !slices.Equal(got, []string{applications}) || diags.HasError() {
			t.Errorf("got %v (HasError %v)", got, diags.HasError())
		}
	})
}

func TestDoNotDecryptListWithoutEmptyComments(t *testing.T) {
	in := gjson.Parse(`[{"id":1,"destination":"a.com","comment":"x"},{"id":2,"destination":"b.com","comment":null},{"id":3,"destination":"c.com","comment":""},{"id":4,"destination":"d.com"}]`)
	got := gjson.Parse(doNotDecryptListWithoutEmptyComments(in))
	if n := len(got.Array()); n != 4 {
		t.Fatalf("got %d destinations, want 4: %s", n, got.Raw)
	}
	if got.Get("0.comment").String() != "x" {
		t.Errorf("a comment must stay: %s", got.Raw)
	}
	for _, i := range []string{"1", "2", "3"} {
		if got.Get(i + ".comment").Exists() {
			t.Errorf("destination %s: a null, empty or missing comment must not be there: %s", i, got.Raw)
		}
	}
	if got.Get("#.destination").Raw != `["a.com","b.com","c.com","d.com"]` || got.Get("#.id").Raw != `[1,2,3,4]` {
		t.Errorf("the other fields must stay: %s", got.Raw)
	}
	if empty := doNotDecryptListWithoutEmptyComments(gjson.Parse(`[]`)); empty != `[]` {
		t.Errorf("an empty list stays empty, got %s", empty)
	}
}

// A list that has no destination list is read as one without destinations if its totals say so.
func TestDoNotDecryptListWithDestinationsWithoutDestinationList(t *testing.T) {
	client, api := newDndReadClient(t, true, 0)
	body, err := doNotDecryptListWithDestinations(context.Background(), client, gjson.Parse(`{"id":5,"name":"n","totalDomains":0,"totalIps":0}`))
	if err != nil {
		t.Fatal(err)
	}
	if body.Get("destinations").Raw != `[]` || body.Get("name").String() != "n" {
		t.Errorf("body = %s", body.Raw)
	}
	for _, object := range []string{`{"id":5,"totalDomains":1,"totalIps":0}`, `{"id":5,"totalDomains":0}`, `{"id":5}`} {
		if _, err := doNotDecryptListWithDestinations(context.Background(), client, gjson.Parse(object)); err == nil {
			t.Errorf("%s: an object that may have destinations but no destination list must be an error", object)
		}
	}
	if got := api.requests(); len(got) != 0 {
		t.Errorf("requests = %v, want none", got)
	}
}

// dndReadServer drives the real resource through the framework server (tfprotov6): ImportResourceState sets the private
// `importing` flag that Read() needs, so the import and the refresh branch of Read() can both be run without a seam.
type dndReadServer struct {
	t       *testing.T
	srv     tfprotov6.ProviderServer
	resType tftypes.Type
}

func newDndReadServer(t *testing.T, apiURL string) *dndReadServer {
	t.Helper()
	ctx := context.Background()
	t.Setenv("SECUREACCESS_URL", apiURL)
	t.Setenv("SECUREACCESS_API_KEY", "key")
	t.Setenv("SECUREACCESS_API_KEY_SECRET", "secret")
	srv := providerserver.NewProtocol6(New("test")())()
	sch, err := srv.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	provType := sch.Provider.ValueType().(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for k, typ := range provType.AttributeTypes {
		vals[k] = tftypes.NewValue(typ, nil)
	}
	cfg, err := tfprotov6.NewDynamicValue(provType, tftypes.NewValue(provType, vals))
	if err != nil {
		t.Fatal(err)
	}
	if cr, err := srv.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{TerraformVersion: "1.9.0", Config: &cfg}); err != nil || len(cr.Diagnostics) > 0 {
		t.Fatal(err, cr.Diagnostics)
	}
	return &dndReadServer{t: t, srv: srv, resType: sch.ResourceSchemas["secureaccess_do_not_decrypt_list"].ValueType()}
}

// destinations returns the destinations of a state, sorted, with their comments (nil for a null comment).
func (s *dndReadServer) destinations(state *tfprotov6.DynamicValue) map[string]*string {
	s.t.Helper()
	v, err := state.Unmarshal(s.resType)
	if err != nil {
		s.t.Fatal(err)
	}
	var obj map[string]tftypes.Value
	if err := v.As(&obj); err != nil {
		s.t.Fatal(err)
	}
	out := map[string]*string{}
	if obj["destinations"].IsNull() {
		return out
	}
	var set []tftypes.Value
	if err := obj["destinations"].As(&set); err != nil {
		s.t.Fatal(err)
	}
	for _, e := range set {
		var o map[string]tftypes.Value
		if err := e.As(&o); err != nil {
			s.t.Fatal(err)
		}
		var d string
		if err := o["destination"].As(&d); err != nil {
			s.t.Fatal(err)
		}
		if o["comment"].IsNull() {
			out[d] = nil
			continue
		}
		var c string
		if err := o["comment"].As(&c); err != nil {
			s.t.Fatal(err)
		}
		out[d] = &c
	}
	return out
}

func (s *dndReadServer) read(state *tfprotov6.DynamicValue, private []byte) *tfprotov6.ReadResourceResponse {
	s.t.Helper()
	rr, err := s.srv.ReadResource(context.Background(), &tfprotov6.ReadResourceRequest{TypeName: "secureaccess_do_not_decrypt_list", CurrentState: state, Private: private})
	if err != nil {
		s.t.Fatal(err)
	}
	return rr
}

func (s *dndReadServer) importThenRead(id string) *tfprotov6.ReadResourceResponse {
	s.t.Helper()
	imp, err := s.srv.ImportResourceState(context.Background(), &tfprotov6.ImportResourceStateRequest{TypeName: "secureaccess_do_not_decrypt_list", ID: id})
	if err != nil || len(imp.Diagnostics) > 0 {
		s.t.Fatal(err, imp.Diagnostics)
	}
	return s.read(imp.ImportedResources[0].State, imp.ImportedResources[0].Private)
}

// stateWith is a prior state of the list with id 5 (a refresh, no import flag).
func (s *dndReadServer) stateWith(destinations ...DoNotDecryptListDestinations) *tfprotov6.DynamicValue {
	s.t.Helper()
	sch := dndTestSchema(s.t)
	model := dndTestModel(dndInt64Set(1), dndInt64Set(2))
	model.Id = types.StringValue("5")
	model.Destinations = destinations
	st := dndTestState(s.t, sch, &model)
	dv, err := tfprotov6.NewDynamicValue(sch.Type().TerraformType(context.Background()), st.Raw)
	if err != nil {
		s.t.Fatal(err)
	}
	return &dv
}

// Read() of the resource: after an import (private flag) all destinations of the destination list are read, a refresh
// only updates the destinations of the state, a null or empty comment of the API is no comment, and a destination list
// that cannot be read fails the refresh instead of emptying the state.
func TestDoNotDecryptListResourceRead(t *testing.T) {
	var mu sync.Mutex
	const withList = `{"id":5,"name":"n","isDefault":false,"exceptiondomainlistId":777,"totalDomains":2,"totalIps":0,"totalApplications":0}`
	object := withList
	destCode := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/auth/v2/token":
			_, _ = io.WriteString(w, `{"access_token":"t","token_type":"bearer","expires_in":3600}`)
		case "/policies/v2/doNotDecryptLists/5":
			_, _ = io.WriteString(w, object)
		case "/policies/v2/destinationlists/777/destinations":
			if destCode != 0 {
				w.WriteHeader(destCode)
				_, _ = io.WriteString(w, `{"message":"no"}`)
				return
			}
			_, _ = io.WriteString(w, `{"data":[{"id":1,"destination":"a.com","comment":"x","type":"domain"},{"id":2,"destination":"c.com","comment":null,"type":"domain"},{"id":3,"destination":"d.com","comment":"","type":"domain"}],"meta":{"total":3}}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	s := newDndReadServer(t, server.URL)
	set := func(o string, code int) { mu.Lock(); object, destCode = o, code; mu.Unlock() }
	comment := func(m map[string]*string, destination string) string {
		c, ok := m[destination]
		switch {
		case !ok:
			return "(not there)"
		case c == nil:
			return "(null)"
		}
		return *c
	}

	t.Run("after an import all destinations are read, a null or empty comment is no comment", func(t *testing.T) {
		rr := s.importThenRead("5")
		if len(rr.Diagnostics) > 0 {
			t.Fatal(rr.Diagnostics)
		}
		got := s.destinations(rr.NewState)
		if len(got) != 3 || comment(got, "a.com") != "x" || comment(got, "c.com") != "(null)" || comment(got, "d.com") != "(null)" {
			t.Errorf("destinations = a.com:%s c.com:%s d.com:%s (%d in total), want x, null, null (3)", comment(got, "a.com"), comment(got, "c.com"), comment(got, "d.com"), len(got))
		}
	})

	t.Run("a refresh updates the destinations of the state only", func(t *testing.T) {
		rr := s.read(s.stateWith(dndDestination("a.com", "old"), dndDestination("b.com", "")), nil)
		if len(rr.Diagnostics) > 0 {
			t.Fatal(rr.Diagnostics)
		}
		got := s.destinations(rr.NewState)
		if len(got) != 1 || comment(got, "a.com") != "x" {
			t.Errorf("destinations = %v, want only a.com with the comment of the destination list (b.com is gone, c.com and d.com are not managed)", got)
		}
	})

	t.Run("a destination list that cannot be read fails the refresh", func(t *testing.T) {
		set(withList, http.StatusBadRequest)
		defer set(withList, 0)
		if rr := s.read(s.stateWith(dndDestination("a.com", "x")), nil); len(rr.Diagnostics) == 0 {
			t.Errorf("a failing destination request must fail the refresh, got the state %v", s.destinations(rr.NewState))
		}
	})

	t.Run("an object without a destination list but with destinations fails the refresh", func(t *testing.T) {
		set(`{"id":5,"name":"n","isDefault":false,"totalDomains":2,"totalIps":0,"totalApplications":0}`, 0)
		defer set(withList, 0)
		if rr := s.read(s.stateWith(dndDestination("a.com", "x")), nil); len(rr.Diagnostics) == 0 {
			t.Errorf("an object without exceptiondomainlistId must fail the refresh")
		}
	})

	t.Run("an object without a destination list and without destinations is read as empty", func(t *testing.T) {
		set(`{"id":5,"name":"n","isDefault":false,"totalDomains":0,"totalIps":0,"totalApplications":0}`, 0)
		defer set(withList, 0)
		rr := s.read(s.stateWith(dndDestination("a.com", "x")), nil)
		if len(rr.Diagnostics) > 0 {
			t.Fatal(rr.Diagnostics)
		}
		if got := s.destinations(rr.NewState); len(got) != 0 {
			t.Errorf("destinations = %v, want none", got)
		}
	})
}

// Update() removes IP addresses through the destination list behind the object, and treats a destination that the API
// reports as not found as removed if it is not in the destination list anymore (a destination that disappeared between
// the refresh and the apply).
func TestDoNotDecryptListUpdateRemovals(t *testing.T) {
	ctx := context.Background()
	var mu sync.Mutex
	var reqs []string
	putFailures := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		defer mu.Unlock()
		if r.URL.Path != "/auth/v2/token" {
			reqs = append(reqs, r.Method+" "+r.URL.Path+" "+string(body))
		}
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/auth/v2/token":
			_, _ = io.WriteString(w, `{"access_token":"t","token_type":"bearer","expires_in":3600}`)
		case r.Method == http.MethodGet && r.URL.Path == "/policies/v2/doNotDecryptLists/1":
			_, _ = io.WriteString(w, `{"id":1,"name":"n","isDefault":false,"exceptiondomainlistId":777,"totalDomains":1,"totalIps":1,"totalApplications":0}`)
		case r.Method == http.MethodGet && r.URL.Path == "/policies/v2/destinationlists/777/destinations":
			_, _ = io.WriteString(w, `{"data":[{"id":11,"destination":"a.com","type":"domain"},{"id":12,"destination":"10.10.10.10","type":"ipv4"}],"meta":{"total":2}}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/policies/v2/destinationlists/777/destinations/remove":
			_, _ = io.WriteString(w, `{}`)
		case r.Method == http.MethodPut && r.URL.Path == "/policies/v2/doNotDecryptLists/1":
			if putFailures > 0 {
				putFailures--
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, `{"message":"destination 'b.com' not found"}`)
				return
			}
			_, _ = io.WriteString(w, `{"id":1,"name":"n","isDefault":false,"totalDomains":1,"totalIps":0,"totalApplications":0}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	client, err := secureaccess.NewClient(server.URL, "key", "secret")
	if err != nil {
		t.Fatal(err)
	}
	r := &DoNotDecryptListResource{client: client}
	s := dndTestSchema(t)
	null := types.SetNull(types.Int64Type)
	update := func(t *testing.T, before []DoNotDecryptListDestinations, after ...DoNotDecryptListDestinations) ([]string, resource.UpdateResponse) {
		t.Helper()
		mu.Lock()
		reqs = nil
		mu.Unlock()
		state, plan := dndTestModel(null, null), dndTestModel(null, null)
		state.Destinations, plan.Destinations = before, after
		prior := dndTestState(t, s, &state)
		resp := resource.UpdateResponse{State: dndTestState(t, s, &state)}
		r.Update(ctx, resource.UpdateRequest{Plan: dndTestPlan(t, s, &plan), State: prior}, &resp)
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(reqs), resp
	}

	t.Run("an IP address is removed through the destination list behind the object", func(t *testing.T) {
		got, resp := update(t, []DoNotDecryptListDestinations{dndDestination("a.com", ""), dndDestination("10.10.10.10", "")}, dndDestination("a.com", ""))
		if resp.Diagnostics.HasError() {
			t.Fatal(resp.Diagnostics)
		}
		want := []string{
			"GET /policies/v2/doNotDecryptLists/1 ",
			"GET /policies/v2/destinationlists/777/destinations ",
			"DELETE /policies/v2/destinationlists/777/destinations/remove [12]",
			"GET /policies/v2/doNotDecryptLists/1 ",
			// the object and the destinations once more, to get the ids of the destinations
			"GET /policies/v2/doNotDecryptLists/1 ",
			"GET /policies/v2/destinationlists/777/destinations ",
		}
		if !slices.Equal(got, want) {
			t.Errorf("requests = %v, want %v", got, want)
		}
	})

	t.Run("a domain that is gone is left out and the removal is repeated", func(t *testing.T) {
		putFailures = 1
		got, resp := update(t, []DoNotDecryptListDestinations{dndDestination("a.com", ""), dndDestination("b.com", "")}, dndDestination("a.com", ""))
		if resp.Diagnostics.HasError() {
			t.Fatal(resp.Diagnostics)
		}
		if len(got) != 6 || !strings.Contains(got[0], `"removeDestinations":[{"destination":"b.com"}]`) ||
			got[1] != "GET /policies/v2/doNotDecryptLists/1 " || got[2] != "GET /policies/v2/destinationlists/777/destinations " ||
			strings.Contains(got[3], "removeDestinations") ||
			got[4] != "GET /policies/v2/doNotDecryptLists/1 " || got[5] != "GET /policies/v2/destinationlists/777/destinations " {
			t.Errorf("requests = %v", got)
		}
	})
}
