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
	"github.com/tidwall/gjson"
)

// This file is not generated. It adds to model_secureaccess_do_not_decrypt_list_ids_test.go the cases that need a fake API
// that keeps state (a PUT changes the destination list that the read-back sees): the order of the read-back relative to the
// other requests, requests that fail, and what is read from the answers. The ids, names and totals are made up.

// dndGapsAPI is a fake API for one Do Not Decrypt List (id 5) whose destination list (id 777) is changed by POST and PUT.
// The id of a destination is made up from its name (dndGapsID). Every request except the token request is recorded.
type dndGapsAPI struct {
	mu         sync.Mutex
	items      []dndGapsItem
	reqs       []string
	posts      []string // the bodies of the POST requests
	failPut    int      // the n-th PUT (1 is the first) is answered with a 400, 0 never
	puts       int
	objectCode int  // the status of GET of the object, 0 is a 200
	listCode   int  // the status of GET of the destinations, 0 is a 200
	postTotals bool // the answer of POST has the totals (the answer of PUT always has them, the one of GET too)
	url        string
}

type dndGapsItem struct {
	destination, comment string
}

// dndGapsID is the id the fake API gives to a destination: host<N>.example.com has 1000+N.
func dndGapsID(destination string) int64 {
	var n int64
	if _, err := fmt.Sscanf(destination, "host%d.example.com", &n); err == nil {
		return 1000 + n
	}
	return map[string]int64{"a.com": 11, "b.com": 12, "c.com": 13}[destination]
}

func (a *dndGapsAPI) object(withList, withTotals bool) string {
	out := `{"id":5,"name":"n","isDefault":false`
	if withList {
		out += `,"exceptiondomainlistId":777`
	}
	if withTotals {
		out += fmt.Sprintf(`,"totalDomains":%d,"totalIps":0,"totalApplications":0`, len(a.items))
	}
	return out + "}"
}

func (a *dndGapsAPI) add(list gjson.Result) {
	list.ForEach(func(_, v gjson.Result) bool {
		a.items = append(a.items, dndGapsItem{destination: v.Get("destination").String(), comment: v.Get("comment").String()})
		return true
	})
}

func newDndGapsAPI(t *testing.T) *dndGapsAPI {
	t.Helper()
	a := &dndGapsAPI{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		a.mu.Lock()
		defer a.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/auth/v2/token" {
			a.reqs = append(a.reqs, r.Method+" "+r.URL.Path)
		}
		switch {
		case r.URL.Path == "/auth/v2/token":
			_, _ = io.WriteString(w, `{"access_token":"t","token_type":"bearer","expires_in":3600}`)
		case r.Method == http.MethodPost && r.URL.Path == "/policies/v2/doNotDecryptLists":
			a.posts = append(a.posts, string(body))
			a.add(gjson.GetBytes(body, "destinations"))
			_, _ = io.WriteString(w, a.object(true, a.postTotals))
		case r.Method == http.MethodPut && r.URL.Path == "/policies/v2/doNotDecryptLists/5":
			a.puts++
			if a.failPut == a.puts {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, `{"message":"no"}`)
				return
			}
			a.add(gjson.GetBytes(body, "addDestinations"))
			gjson.GetBytes(body, "removeDestinations").ForEach(func(_, v gjson.Result) bool {
				a.items = slices.DeleteFunc(a.items, func(i dndGapsItem) bool { return i.destination == v.Get("destination").String() })
				return true
			})
			_, _ = io.WriteString(w, a.object(false, true)) // no exceptiondomainlistId: the answer of PUT does not always have it
		case r.Method == http.MethodGet && r.URL.Path == "/policies/v2/doNotDecryptLists/5":
			if a.objectCode != 0 {
				w.WriteHeader(a.objectCode)
				_, _ = io.WriteString(w, `{"message":"object gone"}`)
				return
			}
			_, _ = io.WriteString(w, a.object(true, true))
		case r.Method == http.MethodGet && r.URL.Path == "/policies/v2/destinationlists/777/destinations":
			if a.listCode != 0 {
				w.WriteHeader(a.listCode)
				_, _ = io.WriteString(w, `{"message":"no"}`)
				return
			}
			items := make([]string, 0, len(a.items))
			for _, i := range a.items {
				item := fmt.Sprintf(`{"id":%d,"destination":%q,"type":"domain"`, dndGapsID(i.destination), i.destination)
				if i.comment != "" {
					item += fmt.Sprintf(`,"comment":%q`, i.comment)
				}
				items = append(items, item+"}")
			}
			_, _ = fmt.Fprintf(w, `{"data":[%s],"meta":{"total":%d}}`, strings.Join(items, ","), len(items))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	a.url = server.URL
	return a
}

func (a *dndGapsAPI) requests() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return slices.Clone(a.reqs)
}

func newDndGapsDriver(t *testing.T, api *dndGapsAPI) *dndIdsDriver {
	t.Helper()
	s := newDndReadServer(t, api.url)
	object := s.resType.(tftypes.Object)
	return &dndIdsDriver{t: t, s: s, object: object, elemType: object.AttributeTypes["destinations"].(tftypes.Set).ElementType}
}

// attribute returns the value of a top level attribute of a state as text: "null", "unknown", or the number or string.
func (d *dndIdsDriver) attribute(state *tfprotov6.DynamicValue, name string) string {
	d.t.Helper()
	v, err := state.Unmarshal(d.object)
	if err != nil {
		d.t.Fatal(err)
	}
	var object map[string]tftypes.Value
	if err := v.As(&object); err != nil {
		d.t.Fatal(err)
	}
	switch a := object[name]; {
	case !a.IsKnown():
		return "unknown"
	case a.IsNull():
		return "null"
	default:
		var f big.Float
		if err := a.As(&f); err != nil {
			d.t.Fatalf("%s: %v", name, err)
		}
		return f.Text('f', 0)
	}
}

// prior is a state of the list with id 5 that holds the given destinations.
func (d *dndIdsDriver) prior(totalDomains int64, elements ...tftypes.Value) *tfprotov6.DynamicValue {
	return d.value(map[string]tftypes.Value{
		"id":                 tftypes.NewValue(tftypes.String, "5"),
		"name":               tftypes.NewValue(tftypes.String, "n"),
		"is_default":         tftypes.NewValue(tftypes.Bool, false),
		"total_domains":      tftypes.NewValue(tftypes.Number, totalDomains),
		"total_ips":          tftypes.NewValue(tftypes.Number, 0),
		"total_applications": tftypes.NewValue(tftypes.Number, 0),
		"destinations":       tftypes.NewValue(tftypes.Set{ElementType: d.elemType}, elements),
	})
}

// hosts returns n new destinations (unknown ids), host0.example.com and so on.
func (d *dndIdsDriver) hosts(n int) []tftypes.Value {
	out := make([]tftypes.Value, n)
	for i := range out {
		out[i] = d.destination(nil, fmt.Sprintf("host%d.example.com", i), "")
	}
	return out
}

// checkHostIds fails unless the state holds exactly the destinations a.com (if withA) and host0 .. host<n-1>, with the ids of
// the fake API (want: every id filled), or with a null id for the hosts (!want).
func (d *dndIdsDriver) checkHostIds(state *tfprotov6.DynamicValue, n int, withA, want bool) {
	d.t.Helper()
	got := d.destinationIds(state)
	size := n
	if withA {
		size++
		if got["a.com"] != "11" {
			d.t.Errorf("a.com has the id %q, want 11 (it is carried from the state)", got["a.com"])
		}
	}
	if len(got) != size {
		d.t.Errorf("the state has %d destinations, want %d", len(got), size)
	}
	known := 0
	for destination, id := range got {
		if destination == "a.com" {
			continue
		}
		switch {
		case want && id != fmt.Sprint(dndGapsID(destination)):
			d.t.Errorf("%s has the id %q, want %d", destination, id, dndGapsID(destination))
		case !want && id != "null":
			d.t.Errorf("%s has the id %q, want null", destination, id)
		default:
			known++
		}
	}
	if known != n {
		d.t.Errorf("%d of the %d new destinations are as wanted", known, n)
	}
}

func TestDoNotDecryptListReadBackGaps(t *testing.T) {
	const (
		post   = "POST /policies/v2/doNotDecryptLists"
		put    = "PUT /policies/v2/doNotDecryptLists/5"
		object = "GET /policies/v2/doNotDecryptLists/5"
		list   = "GET /policies/v2/destinationlists/777/destinations"
	)
	newNull := func(t *testing.T, d *dndIdsDriver) *tfprotov6.DynamicValue {
		t.Helper()
		null, err := tfprotov6.NewDynamicValue(d.object, tftypes.NewValue(d.object, nil))
		if err != nil {
			t.Fatal(err)
		}
		return &null
	}
	assertRequests := func(t *testing.T, api *dndGapsAPI, want ...string) {
		t.Helper()
		if got := api.requests(); !slices.Equal(got, want) {
			t.Errorf("requests = %v, want %v", got, want)
		}
	}

	t.Run("create: 501 destinations, the rest is added by a PUT, the read-back comes last and its ids are used", func(t *testing.T) {
		api := newDndGapsAPI(t)
		d := newDndGapsDriver(t, api)
		null := newNull(t, d)
		config := d.config(d.hosts(501)...)
		applied := d.apply(null, d.plan(null, config), config)
		if warnings, errors := dndWarnings(applied.Diagnostics); len(warnings) > 0 || len(errors) > 0 {
			t.Fatalf("warnings %v, errors %v", warnings, errors)
		}
		// the API takes 500 destinations with the create, the 501st is added by a PUT: only a read after it has its id
		assertRequests(t, api, post, put, object, list)
		d.checkHostIds(applied.NewState, 501, false, true)
		if !d.fullyKnown(applied.NewState) {
			t.Errorf("the state has an unknown value after the apply")
		}
		if got := d.attribute(applied.NewState, "total_domains"); got != "501" {
			t.Errorf("total_domains = %s, want 501", got)
		}
		// the ids are computed, the request must not carry them (an unknown id would be sent as 0)
		if len(api.posts) != 1 || gjson.Get(api.posts[0], "destinations.#").Int() != 500 || len(gjson.Get(api.posts[0], "destinations.#.id").Array()) != 0 {
			t.Errorf("the create request must carry 500 destinations without an id: %.300s", strings.Join(api.posts, ", "))
		}
	})

	t.Run("create without destinations: the totals are read from the object, the answer of POST has none", func(t *testing.T) {
		api := newDndGapsAPI(t)
		d := newDndGapsDriver(t, api)
		null := newNull(t, d)
		config := d.config()
		applied := d.apply(null, d.plan(null, config), config)
		if warnings, errors := dndWarnings(applied.Diagnostics); len(warnings) > 0 || len(errors) > 0 {
			t.Fatalf("warnings %v, errors %v", warnings, errors)
		}
		assertRequests(t, api, post, object, list)
		for _, name := range []string{"total_domains", "total_ips", "total_applications"} {
			if got := d.attribute(applied.NewState, name); got != "0" {
				t.Errorf("%s = %s, want 0 (from the object)", name, got)
			}
		}
	})

	// A read-back that fails is a warning, the answer of the last request is the fallback. Its totals must be kept.
	for _, tc := range []struct {
		name       string
		set        func(*dndGapsAPI)
		wantDetail []string // texts of the warning
	}{
		{"the destinations cannot be read", func(a *dndGapsAPI) { a.listCode = http.StatusBadRequest }, []string{"StatusCode 400"}},
		{"the object cannot be read", func(a *dndGapsAPI) { a.objectCode = http.StatusBadRequest }, []string{"StatusCode 400", "failed to retrieve object", "object gone"}},
	} {
		t.Run("create: "+tc.name+", the object is kept with the totals of the POST answer", func(t *testing.T) {
			api := newDndGapsAPI(t)
			api.postTotals = true
			tc.set(api)
			d := newDndGapsDriver(t, api)
			null := newNull(t, d)
			config := d.config(d.destination(nil, "a.com", "x"))
			applied := d.apply(null, d.plan(null, config), config)
			warnings, errors := dndWarnings(applied.Diagnostics)
			if len(errors) > 0 || len(warnings) != 1 || !strings.HasPrefix(warnings[0], "Destinations not read back: The object is configured") {
				t.Fatalf("warnings %v, errors %v, want one warning that says that the object is configured", warnings, errors)
			}
			for _, want := range tc.wantDetail {
				if !strings.Contains(warnings[0], want) {
					t.Errorf("the warning lacks %q: %s", want, warnings[0])
				}
			}
			if got := d.destinationIds(applied.NewState); got["a.com"] != "null" || len(got) != 1 {
				t.Errorf("ids = %v, want a.com null", got)
			}
			if got := d.attribute(applied.NewState, "total_domains"); got != "1" {
				t.Errorf("total_domains = %s, want 1, the value of the POST answer (the answer of the last request is the fallback)", got)
			}
		})

		t.Run("update: "+tc.name+", the totals of the PUT answer are kept", func(t *testing.T) {
			api := newDndGapsAPI(t)
			api.items = []dndGapsItem{{"a.com", "x"}}
			d := newDndGapsDriver(t, api)
			prior := d.prior(1, d.destination(dndInt(11), "a.com", "x"))
			config := d.config(d.destination(nil, "a.com", "x"), d.destination(nil, "b.com", ""))
			proposed := d.value(map[string]tftypes.Value{
				"id": tftypes.NewValue(tftypes.String, "5"), "name": tftypes.NewValue(tftypes.String, "n"),
				"destinations": tftypes.NewValue(tftypes.Set{ElementType: d.elemType}, []tftypes.Value{d.destination(dndInt(11), "a.com", "x"), d.destination(nil, "b.com", "")}),
			})
			tc.set(api)
			applied := d.apply(prior, d.plan(prior, proposed), config)
			warnings, errors := dndWarnings(applied.Diagnostics)
			if len(errors) > 0 || len(warnings) != 1 || !strings.HasPrefix(warnings[0], "Destinations not read back: The object is configured") {
				t.Fatalf("warnings %v, errors %v, want one warning that says that the object is configured", warnings, errors)
			}
			for _, want := range tc.wantDetail {
				if !strings.Contains(warnings[0], want) {
					t.Errorf("the warning lacks %q: %s", want, warnings[0])
				}
			}
			if got := d.destinationIds(applied.NewState); got["a.com"] != "11" || got["b.com"] != "null" || len(got) != 2 {
				t.Errorf("ids = %v, want a.com 11 and b.com null", got)
			}
			if got := d.attribute(applied.NewState, "total_domains"); got != "2" {
				t.Errorf("total_domains = %s, want 2, the value of the PUT answer (the answer of the last request is the fallback)", got)
			}
		})
	}

	t.Run("update with two requests: one read-back after the last, and its ids are used", func(t *testing.T) {
		api := newDndGapsAPI(t)
		api.items = []dndGapsItem{{"a.com", "x"}}
		d := newDndGapsDriver(t, api)
		prior := d.prior(1, d.destination(dndInt(11), "a.com", "x"))
		config := d.config(append([]tftypes.Value{d.destination(nil, "a.com", "x")}, d.hosts(501)...)...)
		proposed := d.value(map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.String, "5"), "name": tftypes.NewValue(tftypes.String, "n"),
			"destinations": tftypes.NewValue(tftypes.Set{ElementType: d.elemType}, append([]tftypes.Value{d.destination(dndInt(11), "a.com", "x")}, d.hosts(501)...)),
		})
		applied := d.apply(prior, d.plan(prior, proposed), config)
		if warnings, errors := dndWarnings(applied.Diagnostics); len(warnings) > 0 || len(errors) > 0 {
			t.Fatalf("warnings %v, errors %v", warnings, errors)
		}
		// 501 new destinations are two requests (500 and 1), the object and the destinations are read once, after both
		assertRequests(t, api, put, put, object, list)
		d.checkHostIds(applied.NewState, 501, true, true)
		if !d.fullyKnown(applied.NewState) {
			t.Errorf("the state has an unknown value after the apply")
		}
	})

	t.Run("update with two requests, the second fails: the state holds no unknown id", func(t *testing.T) {
		api := newDndGapsAPI(t)
		api.items = []dndGapsItem{{"a.com", "x"}}
		api.failPut = 2
		d := newDndGapsDriver(t, api)
		prior := d.prior(1, d.destination(dndInt(11), "a.com", "x"))
		config := d.config(append([]tftypes.Value{d.destination(nil, "a.com", "x")}, d.hosts(501)...)...)
		proposed := d.value(map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.String, "5"), "name": tftypes.NewValue(tftypes.String, "n"),
			"destinations": tftypes.NewValue(tftypes.Set{ElementType: d.elemType}, append([]tftypes.Value{d.destination(dndInt(11), "a.com", "x")}, d.hosts(501)...)),
		})
		applied := d.apply(prior, d.plan(prior, proposed), config)
		_, errors := dndWarnings(applied.Diagnostics)
		// the error of the second request is the only one: an unknown value in the state would be a second error of the framework
		if len(errors) != 1 || !strings.Contains(errors[0], "request 2 of 2") {
			t.Fatalf("errors %v, want only the one of the second request", errors)
		}
		assertRequests(t, api, put, put)
		if applied.NewState == nil || !d.fullyKnown(applied.NewState) {
			t.Fatalf("the state of the first request must be kept and have no unknown value")
		}
		// the 500 destinations of the first request are in the state, without ids: the PUT answer has none
		d.checkHostIds(applied.NewState, 500, true, false)
	})

	t.Run("a refresh fills the ids that a failed read-back left empty, an import reads them", func(t *testing.T) {
		api := newDndIdsAPI(t)
		api.set(0, dndIdsItemA, dndIdsItemB)
		d := newDndIdsDriver(t, api)
		prior := d.prior(2, d.destination(nil, "a.com", "x"), d.destination(nil, "b.com", ""))
		rr := d.s.read(prior, nil)
		if len(rr.Diagnostics) > 0 {
			t.Fatal(rr.Diagnostics)
		}
		if got := d.destinationIds(rr.NewState); got["a.com"] != "11" || got["b.com"] != "12" || len(got) != 2 {
			t.Errorf("after a refresh: ids = %v, want a.com 11 and b.com 12", got)
		}
		imp := d.s.importThenRead("5")
		if len(imp.Diagnostics) > 0 {
			t.Fatal(imp.Diagnostics)
		}
		if got := d.destinationIds(imp.NewState); got["a.com"] != "11" || got["b.com"] != "12" || len(got) != 2 {
			t.Errorf("after an import: ids = %v, want a.com 11 and b.com 12", got)
		}
	})
}
