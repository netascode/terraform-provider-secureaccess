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

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// This file is not generated. It tests the `id` of the destinations of a Do Not Decrypt List through the framework server
// (planning and applying like Terraform does): the API does not return the destinations with the object, so Create() and
// Update() read the destination list behind the object once more to fill the ids that are unknown after the apply. The ids,
// the names and the totals in the fake API are made up.

// dndIdsAPI is a fake Secure Access API for one Do Not Decrypt List (id 5) whose destination list (id 777) holds the given
// items. Every request except the token request is recorded.
type dndIdsAPI struct {
	mu       sync.Mutex
	items    []string // the destinations of the destination list, JSON objects
	listCode int      // the status of the request for the destinations, 0 is a 200
	reqs     []string
	url      string
}

func newDndIdsAPI(t *testing.T) *dndIdsAPI {
	t.Helper()
	a := &dndIdsAPI{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		a.mu.Lock()
		defer a.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/auth/v2/token" {
			a.reqs = append(a.reqs, r.Method+" "+r.URL.Path)
		}
		switch {
		case r.URL.Path == "/auth/v2/token":
			_, _ = io.WriteString(w, `{"access_token":"t","token_type":"bearer","expires_in":3600}`)
		case (r.Method == http.MethodPost && r.URL.Path == "/policies/v2/doNotDecryptLists") || r.URL.Path == "/policies/v2/doNotDecryptLists/5":
			_, _ = fmt.Fprintf(w, `{"id":5,"name":"n","isDefault":false,"exceptiondomainlistId":777,"totalDomains":%d,"totalIps":0,"totalApplications":0}`, len(a.items))
		case r.Method == http.MethodGet && r.URL.Path == "/policies/v2/destinationlists/777/destinations":
			if a.listCode != 0 {
				w.WriteHeader(a.listCode)
				_, _ = io.WriteString(w, `{"message":"no"}`)
				return
			}
			_, _ = fmt.Fprintf(w, `{"data":[%s],"meta":{"total":%d}}`, strings.Join(a.items, ","), len(a.items))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	a.url = server.URL
	return a
}

func (a *dndIdsAPI) set(code int, items ...string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.items, a.listCode, a.reqs = items, code, nil
}

func (a *dndIdsAPI) requests() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return slices.Clone(a.reqs)
}

// dndIdsDriver builds states and plans for the resource and runs them through the framework server.
type dndIdsDriver struct {
	t        *testing.T
	s        *dndReadServer
	object   tftypes.Object
	elemType tftypes.Type
}

func newDndIdsDriver(t *testing.T, api *dndIdsAPI) *dndIdsDriver {
	t.Helper()
	s := newDndReadServer(t, api.url)
	object := s.resType.(tftypes.Object)
	return &dndIdsDriver{t: t, s: s, object: object, elemType: object.AttributeTypes["destinations"].(tftypes.Set).ElementType}
}

// value is an object of the resource: the given attributes, every other one is null.
func (d *dndIdsDriver) value(set map[string]tftypes.Value) *tfprotov6.DynamicValue {
	d.t.Helper()
	values := map[string]tftypes.Value{}
	for name, typ := range d.object.AttributeTypes {
		values[name] = tftypes.NewValue(typ, nil)
	}
	for name, v := range set {
		values[name] = v
	}
	dv, err := tfprotov6.NewDynamicValue(d.object, tftypes.NewValue(d.object, values))
	if err != nil {
		d.t.Fatal(err)
	}
	return &dv
}

// destination is an element of `destinations`, with the given id (nil is null) and comment ("" is null).
func (d *dndIdsDriver) destination(id *int64, destination, comment string) tftypes.Value {
	idValue := tftypes.NewValue(tftypes.Number, nil)
	if id != nil {
		idValue = tftypes.NewValue(tftypes.Number, *id)
	}
	commentValue := tftypes.NewValue(tftypes.String, nil)
	if comment != "" {
		commentValue = tftypes.NewValue(tftypes.String, comment)
	}
	return tftypes.NewValue(d.elemType, map[string]tftypes.Value{"id": idValue, "destination": tftypes.NewValue(tftypes.String, destination), "comment": commentValue})
}

func (d *dndIdsDriver) config(elements ...tftypes.Value) *tfprotov6.DynamicValue {
	return d.value(map[string]tftypes.Value{
		"name":         tftypes.NewValue(tftypes.String, "n"),
		"destinations": tftypes.NewValue(tftypes.Set{ElementType: d.elemType}, elements),
	})
}

// plan runs PlanResourceChange, the framework marks the computed values that are null as unknown.
func (d *dndIdsDriver) plan(prior, proposed *tfprotov6.DynamicValue) *tfprotov6.PlanResourceChangeResponse {
	d.t.Helper()
	resp, err := d.s.srv.PlanResourceChange(context.Background(), &tfprotov6.PlanResourceChangeRequest{
		TypeName: "secureaccess_do_not_decrypt_list", PriorState: prior, ProposedNewState: proposed, Config: proposed,
	})
	if err != nil {
		d.t.Fatal(err)
	}
	for _, diag := range resp.Diagnostics {
		if diag.Severity == tfprotov6.DiagnosticSeverityError {
			d.t.Fatalf("plan: %s: %s", diag.Summary, diag.Detail)
		}
	}
	return resp
}

func (d *dndIdsDriver) apply(prior *tfprotov6.DynamicValue, plan *tfprotov6.PlanResourceChangeResponse, config *tfprotov6.DynamicValue) *tfprotov6.ApplyResourceChangeResponse {
	d.t.Helper()
	resp, err := d.s.srv.ApplyResourceChange(context.Background(), &tfprotov6.ApplyResourceChangeRequest{
		TypeName: "secureaccess_do_not_decrypt_list", PriorState: prior, PlannedState: plan.PlannedState, Config: config, PlannedPrivate: plan.PlannedPrivate,
	})
	if err != nil {
		d.t.Fatal(err)
	}
	return resp
}

// destinationIds returns "11" for a destination with the id 11, "null" and "unknown" for the ids that are not known.
func (d *dndIdsDriver) destinationIds(state *tfprotov6.DynamicValue) map[string]string {
	d.t.Helper()
	v, err := state.Unmarshal(d.object)
	if err != nil {
		d.t.Fatal(err)
	}
	var object map[string]tftypes.Value
	if err := v.As(&object); err != nil {
		d.t.Fatal(err)
	}
	out := map[string]string{}
	if object["destinations"].IsNull() {
		return out
	}
	var elements []tftypes.Value
	if err := object["destinations"].As(&elements); err != nil {
		d.t.Fatal(err)
	}
	for _, e := range elements {
		var o map[string]tftypes.Value
		if err := e.As(&o); err != nil {
			d.t.Fatal(err)
		}
		var destination string
		if err := o["destination"].As(&destination); err != nil {
			d.t.Fatal(err)
		}
		switch {
		case !o["id"].IsKnown():
			out[destination] = "unknown"
		case o["id"].IsNull():
			out[destination] = "null"
		default:
			var f big.Float
			if err := o["id"].As(&f); err != nil {
				d.t.Fatal(err)
			}
			out[destination] = f.Text('f', 0)
		}
	}
	return out
}

func (d *dndIdsDriver) fullyKnown(state *tfprotov6.DynamicValue) bool {
	d.t.Helper()
	v, err := state.Unmarshal(d.object)
	if err != nil {
		d.t.Fatal(err)
	}
	return v.IsFullyKnown()
}

func dndWarnings(diags []*tfprotov6.Diagnostic) (warnings, errors []string) {
	for _, diag := range diags {
		if diag.Severity == tfprotov6.DiagnosticSeverityWarning {
			warnings = append(warnings, diag.Summary+": "+diag.Detail)
		} else {
			errors = append(errors, diag.Summary+": "+diag.Detail)
		}
	}
	return warnings, errors
}

func dndInt(i int64) *int64 { return &i }

const (
	dndIdsItemA = `{"id":11,"destination":"a.com","comment":"x","type":"domain"}`
	dndIdsItemB = `{"id":12,"destination":"b.com","type":"domain"}`
	dndIdsItemC = `{"id":13,"destination":"c.com","comment":"y","type":"domain"}`
)

func TestDoNotDecryptListDestinationIds(t *testing.T) {
	api := newDndIdsAPI(t)
	d := newDndIdsDriver(t, api)
	null, err := tfprotov6.NewDynamicValue(d.object, tftypes.NewValue(d.object, nil))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("the plan of a new list has unknown ids", func(t *testing.T) {
		plan := d.plan(&null, d.config(d.destination(nil, "a.com", "x"), d.destination(nil, "b.com", "")))
		got := d.destinationIds(plan.PlannedState)
		if got["a.com"] != "unknown" || got["b.com"] != "unknown" || len(got) != 2 {
			t.Fatalf("planned ids = %v, want both unknown (the test needs unknown ids to prove that they are read back)", got)
		}
	})

	// Create() reads the object and the destination list again and fills the ids, which the API did not return with the object.
	var created *tfprotov6.DynamicValue
	t.Run("create fills the ids from the destination list", func(t *testing.T) {
		api.set(0, dndIdsItemA, dndIdsItemB)
		config := d.config(d.destination(nil, "a.com", "x"), d.destination(nil, "b.com", ""))
		plan := d.plan(&null, config)
		applied := d.apply(&null, plan, config)
		if warnings, errors := dndWarnings(applied.Diagnostics); len(warnings) > 0 || len(errors) > 0 {
			t.Fatalf("warnings %v, errors %v", warnings, errors)
		}
		got := d.destinationIds(applied.NewState)
		if got["a.com"] != "11" || got["b.com"] != "12" || len(got) != 2 {
			t.Errorf("ids = %v, want a.com 11 and b.com 12, the ids of the destination list", got)
		}
		if !d.fullyKnown(applied.NewState) {
			t.Errorf("the state has an unknown value after the apply")
		}
		want := []string{
			"POST /policies/v2/doNotDecryptLists",
			"GET /policies/v2/doNotDecryptLists/5",
			"GET /policies/v2/destinationlists/777/destinations",
		}
		if got := api.requests(); !slices.Equal(got, want) {
			t.Errorf("requests = %v, want %v (the create, then the object and the destinations once more)", got, want)
		}
		created = applied.NewState
	})

	t.Run("update keeps the ids of the destinations that stay and fills the id of a new one", func(t *testing.T) {
		if created == nil {
			t.Skip("the create step failed")
		}
		api.set(0, dndIdsItemA, dndIdsItemB, dndIdsItemC)
		// the ids of the destinations that stay are carried from the state, the new one is computed
		config := d.config(d.destination(nil, "a.com", "x"), d.destination(nil, "b.com", ""), d.destination(nil, "c.com", "y"))
		proposed := d.value(map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "5"),
			"name": tftypes.NewValue(tftypes.String, "n"),
			"destinations": tftypes.NewValue(tftypes.Set{ElementType: d.elemType}, []tftypes.Value{
				d.destination(dndInt(11), "a.com", "x"), d.destination(dndInt(12), "b.com", ""), d.destination(nil, "c.com", "y"),
			}),
		})
		plan := d.plan(created, proposed)
		planned := d.destinationIds(plan.PlannedState)
		if planned["a.com"] != "11" || planned["b.com"] != "12" || planned["c.com"] != "unknown" {
			t.Fatalf("planned ids = %v, want 11, 12 and unknown", planned)
		}
		applied := d.apply(created, plan, config)
		if warnings, errors := dndWarnings(applied.Diagnostics); len(warnings) > 0 || len(errors) > 0 {
			t.Fatalf("warnings %v, errors %v", warnings, errors)
		}
		got := d.destinationIds(applied.NewState)
		if got["a.com"] != "11" || got["b.com"] != "12" || got["c.com"] != "13" || len(got) != 3 {
			t.Errorf("ids = %v, want 11, 12 and 13", got)
		}
		if !d.fullyKnown(applied.NewState) {
			t.Errorf("the state has an unknown value after the apply")
		}
		requests := api.requests()
		if n := len(requests); n < 2 || requests[n-2] != "GET /policies/v2/doNotDecryptLists/5" || requests[n-1] != "GET /policies/v2/destinationlists/777/destinations" {
			t.Errorf("requests = %v, want the object and the destinations as the last two", requests)
		}
	})

	t.Run("a read that fails after create is a warning, the ids stay empty and the object is kept", func(t *testing.T) {
		api.set(http.StatusBadRequest, dndIdsItemA)
		config := d.config(d.destination(nil, "a.com", "x"))
		plan := d.plan(&null, config)
		applied := d.apply(&null, plan, config)
		warnings, errors := dndWarnings(applied.Diagnostics)
		if len(errors) > 0 {
			t.Fatalf("the object is created, a failing read must not be an error: %v", errors)
		}
		if len(warnings) != 1 || !strings.Contains(warnings[0], "Destinations not read back") || !strings.Contains(warnings[0], "StatusCode 400") {
			t.Errorf("warnings = %v, want one warning that names the failed read", warnings)
		}
		if applied.NewState == nil {
			t.Fatal("the state must be kept, the object exists")
		}
		if got := d.destinationIds(applied.NewState); got["a.com"] != "null" || len(got) != 1 {
			t.Errorf("ids = %v, want a.com null: no unknown value may be stored", got)
		}
		if !d.fullyKnown(applied.NewState) {
			t.Errorf("the state has an unknown value")
		}
	})

	t.Run("a read that fails after update is a warning, the ids of new destinations stay empty", func(t *testing.T) {
		if created == nil {
			t.Skip("the create step failed")
		}
		api.set(http.StatusBadRequest, dndIdsItemA, dndIdsItemB, dndIdsItemC)
		config := d.config(d.destination(nil, "a.com", "x"), d.destination(nil, "b.com", ""), d.destination(nil, "c.com", "y"))
		proposed := d.value(map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "5"),
			"name": tftypes.NewValue(tftypes.String, "n"),
			"destinations": tftypes.NewValue(tftypes.Set{ElementType: d.elemType}, []tftypes.Value{
				d.destination(dndInt(11), "a.com", "x"), d.destination(dndInt(12), "b.com", ""), d.destination(nil, "c.com", "y"),
			}),
		})
		plan := d.plan(created, proposed)
		applied := d.apply(created, plan, config)
		warnings, errors := dndWarnings(applied.Diagnostics)
		if len(errors) > 0 || len(warnings) != 1 || !strings.Contains(warnings[0], "Destinations not read back") {
			t.Fatalf("warnings %v, errors %v, want one warning and no error", warnings, errors)
		}
		got := d.destinationIds(applied.NewState)
		if got["a.com"] != "11" || got["b.com"] != "12" || got["c.com"] != "null" {
			t.Errorf("ids = %v, want the known ids kept and c.com null", got)
		}
		if !d.fullyKnown(applied.NewState) {
			t.Errorf("the state has an unknown value")
		}
	})
}
