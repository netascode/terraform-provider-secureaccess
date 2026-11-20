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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/netascode/go-secureaccess"
)

// This file is not generated. It tests the data sources application_exception_list and application_exception_lists, and the paging that they share,
// against a fake of the Secure Access API that behaves as the exception application list routes did on the int environment: the answer is
// {status, meta{page, limit, total}, data}, without `limit` a page holds 100 items, `limit` accepts up to 1000 (more is a 400), `limit` 0 and `page` 0 are
// a 400, a page past the end is an empty `data`, the settings of a list are in `data` of the answer of GET by id, and the answer has no `meta`.
// Only the ids, the names and the dates are made up. Lists of more than 125 applications were not observed, the large ones here are synthetic.

// excApp is an application of the list as /applications returns it.
func excApp(i int) map[string]any {
	return map[string]any{
		"applicationId": 1000 + i*3, "applicationName": "App " + strconv.Itoa(i), "action": "no_decrypt", "hasUrls": i%7 == 0,
		"groupId": 40 + i%5, "groupName": "Group " + strconv.Itoa(i%5),
		"selectedAction": map[string]any{"name": "no_decrypt", "displayName": "Allow", "actionId": 1},
	}
}

func excApps(n int) []any {
	out := make([]any, n)
	for i := range out {
		out[i] = excApp(i)
	}
	return out
}

// excSettings is GET /policies/v2/exceptionApplicationLists/{id} as the API answers it, with the counts of n applications.
func excSettings(n int) map[string]any {
	return map[string]any{
		"status": map[string]any{"code": 200, "text": "OK"},
		"data": map[string]any{
			"id": 5001, "organizationId": 2222222, "isDefault": nil, "name": "my_do_not_decrypt_list", "type": "application",
			"createdAt": "2026-10-09 16:16:12", "modifiedAt": "2026-10-16 17:19:09", "markedForDeletion": false, "bundleTypeId": 1, "isSwgDefault": nil,
			"blockedCount": 0, "allowedCount": n, "warnedCount": 0, "applicationsCategories": []any{},
		},
	}
}

func excWrite(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.Marshal(body)
	_, _ = w.Write(b)
}

// excPaged answers a collection route the way the API does, from the given items.
func excPaged(w http.ResponseWriter, r *http.Request, items []any) {
	q := r.URL.Query()
	page, limit := 1, 100
	var err error
	if v := q.Get("page"); v != "" {
		if page, err = strconv.Atoi(v); err != nil || page < 1 {
			excWrite(w, http.StatusBadRequest, map[string]any{"status": map[string]any{"code": 400, "text": "Bad Request"}, "data": map[string]any{"error": "invalid_request", "message": "Page must be a number greater than 0"}})
			return
		}
	}
	if v := q.Get("limit"); v != "" {
		if limit, err = strconv.Atoi(v); err != nil || limit < 1 || limit > 1000 {
			excWrite(w, http.StatusBadRequest, map[string]any{"status": map[string]any{"code": 400, "text": "Bad Request"}, "data": map[string]any{"error": "invalid_request", "message": "Limit must be a number between 1 and 1000"}})
			return
		}
	}
	start := (page - 1) * limit
	end := start + limit
	if end > len(items) {
		end = len(items)
	}
	data := []any{}
	if start < len(items) {
		data = items[start:end]
	}
	excWrite(w, http.StatusOK, map[string]any{
		"status": map[string]any{"code": 200, "text": "OK"},
		"meta":   map[string]any{"page": page, "limit": limit, "total": len(items)},
		"data":   data,
	})
}

// excServer is a fake Secure Access API: it answers the token request and everything else with the handler, and records the method, the escaped path
// and the query of every other request.
type excServer struct {
	*httptest.Server
	mu    sync.Mutex
	paths []string
}

func newExcServer(t *testing.T, handle func(w http.ResponseWriter, r *http.Request)) *excServer {
	t.Helper()
	s := &excServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/v2/token" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"token","token_type":"bearer","expires_in":3600}`))
			return
		}
		s.mu.Lock()
		s.paths = append(s.paths, r.Method+" "+r.URL.EscapedPath()+"?"+r.URL.RawQuery)
		s.mu.Unlock()
		handle(w, r)
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *excServer) requests() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.paths...)
}

// excListServer serves one list with the given applications under the id 5001, everything else is a 404 as the API answers an unknown id.
func excListServer(t *testing.T, apps []any) *excServer {
	t.Helper()
	return newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists/5001":
			excWrite(w, http.StatusOK, excSettings(len(apps)))
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists/5001/applications":
			excPaged(w, r, apps)
		default:
			excWrite(w, http.StatusNotFound, map[string]any{"status": map[string]any{"code": 404, "text": "Not Found"}, "data": map[string]any{"error": "object_not_found", "message": "DND exception application list not found"}})
		}
	})
}

func newExcClient(t *testing.T, server *excServer) *secureaccess.Client {
	t.Helper()
	client, err := secureaccess.NewClient(server.URL, "key", "secret", secureaccess.MaxRetries(0))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

// readExceptionApplicationList runs the Read() of the data source application_exception_list with the given id in the config.
func readExceptionApplicationList(t *testing.T, server *excServer, id string) (ApplicationExceptionList, datasource.ReadResponse) {
	t.Helper()
	ctx := context.Background()
	ds := &ApplicationExceptionListDataSource{client: newExcClient(t, server)}

	var schemaResp datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	objType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	config := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		config[name] = tftypes.NewValue(typ, nil)
	}
	config["id"] = tftypes.NewValue(tftypes.String, id)

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, config)}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, nil)}}
	ds.Read(ctx, req, &resp)

	var state ApplicationExceptionList
	if !resp.Diagnostics.HasError() {
		if diags := resp.State.Get(ctx, &state); diags.HasError() {
			t.Fatalf("reading the state: %v", diags)
		}
	}
	return state, resp
}

// readExceptionApplicationLists runs the Read() of the data source application_exception_lists.
func readExceptionApplicationLists(t *testing.T, server *excServer) (ApplicationExceptionLists, datasource.ReadResponse) {
	t.Helper()
	ctx := context.Background()
	ds := &ApplicationExceptionListsDataSource{client: newExcClient(t, server)}

	var schemaResp datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	objType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	config := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		config[name] = tftypes.NewValue(typ, nil)
	}

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, config)}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, nil)}}
	ds.Read(ctx, req, &resp)

	var state ApplicationExceptionLists
	if !resp.Diagnostics.HasError() {
		if diags := resp.State.Get(ctx, &state); diags.HasError() {
			t.Fatalf("reading the state: %v", diags)
		}
	}
	return state, resp
}

func excErrorDetail(t *testing.T, resp datasource.ReadResponse) string {
	t.Helper()
	if !resp.Diagnostics.HasError() {
		t.Fatalf("want an error, got none")
	}
	d := resp.Diagnostics.Errors()[0]
	if d.Summary() != "Client Error" {
		t.Errorf("summary = %q, want Client Error", d.Summary())
	}
	// the state must not be touched by a read that failed, not even partly (the helpers do not read it after an error)
	if !resp.State.Raw.IsNull() {
		t.Errorf("state must stay empty after an error, got %s", resp.State.Raw)
	}
	return d.Detail()
}

func TestApplicationExceptionListSchema(t *testing.T) {
	ctx := context.Background()

	for _, tc := range []struct {
		name     string
		ds       datasource.DataSource
		typeName string
		want     string
		idNeeded bool
	}{
		{"single list", &ApplicationExceptionListDataSource{}, "secureaccess_application_exception_list", "allowed_count,applications,blocked_count,id,name,warned_count", true},
		{"ids", &ApplicationExceptionListsDataSource{}, "secureaccess_application_exception_lists", "ids", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var meta datasource.MetadataResponse
			tc.ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "secureaccess"}, &meta)
			if meta.TypeName != tc.typeName {
				t.Errorf("type name = %q, want %q", meta.TypeName, tc.typeName)
			}
			var resp datasource.SchemaResponse
			tc.ds.Schema(ctx, datasource.SchemaRequest{}, &resp)
			var names []string
			for name, a := range resp.Schema.Attributes {
				names = append(names, name)
				if name == "id" {
					if !a.IsRequired() || a.IsComputed() {
						t.Errorf("id must be required (lookup by id only), required=%t computed=%t", a.IsRequired(), a.IsComputed())
					}
					continue
				}
				if !a.IsComputed() || a.IsRequired() || a.IsOptional() {
					t.Errorf("%s must be read-only, required=%t optional=%t computed=%t", name, a.IsRequired(), a.IsOptional(), a.IsComputed())
				}
			}
			sort.Strings(names)
			if got := strings.Join(names, ","); got != tc.want {
				t.Errorf("attributes = %s\nwant       %s", got, tc.want)
			}
			// a lookup by name would need a read of every list, there is none
			if _, ok := tc.ds.(datasource.DataSourceWithConfigValidators); ok {
				t.Errorf("there must be no config validators")
			}
			var registered bool
			for _, factory := range New("test")().DataSources(ctx) {
				var m datasource.MetadataResponse
				factory().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "secureaccess"}, &m)
				if m.TypeName == tc.typeName {
					registered = true
				}
			}
			if !registered {
				t.Error("the data source is not registered in the provider")
			}
		})
	}

	t.Run("the nested attributes of an application", func(t *testing.T) {
		var resp datasource.SchemaResponse
		(&ApplicationExceptionListDataSource{}).Schema(ctx, datasource.SchemaRequest{}, &resp)
		objType := resp.Schema.Type().TerraformType(ctx).(tftypes.Object)
		list, ok := objType.AttributeTypes["applications"].(tftypes.List)
		if !ok {
			t.Fatalf("applications is a %T, want a list (the order of the API is kept)", objType.AttributeTypes["applications"])
		}
		var names []string
		for name := range list.ElementType.(tftypes.Object).AttributeTypes {
			names = append(names, name)
		}
		sort.Strings(names)
		if got, want := strings.Join(names, ","), "action,application_id,application_name,group_id,group_name,has_urls"; got != want {
			t.Errorf("attributes of an application = %s, want %s (selectedAction is not exposed)", got, want)
		}
		nested, ok := resp.Schema.Attributes["applications"].(schema.ListNestedAttribute)
		if !ok {
			t.Fatalf("applications is a %T", resp.Schema.Attributes["applications"])
		}
		for name, a := range nested.NestedObject.Attributes {
			if !a.IsComputed() || a.IsRequired() || a.IsOptional() {
				t.Errorf("applications.%s must be read-only, required=%t optional=%t computed=%t", name, a.IsRequired(), a.IsOptional(), a.IsComputed())
			}
		}
	})

	t.Run("Configure keeps the client of the provider", func(t *testing.T) {
		client := &secureaccess.Client{}
		single := &ApplicationExceptionListDataSource{}
		single.Configure(ctx, datasource.ConfigureRequest{ProviderData: &SecureAccessProviderData{Client: client}}, &datasource.ConfigureResponse{})
		plural := &ApplicationExceptionListsDataSource{}
		plural.Configure(ctx, datasource.ConfigureRequest{ProviderData: &SecureAccessProviderData{Client: client}}, &datasource.ConfigureResponse{})
		if single.client != client || plural.client != client {
			t.Errorf("the data sources did not keep the client: %p %p, want %p", single.client, plural.client, client)
		}
	})
}

func TestApplicationExceptionListDataSourceRead(t *testing.T) {
	t.Run("reads the settings and every application with limit 1000, in two requests", func(t *testing.T) {
		server := excListServer(t, excApps(125))
		state, resp := readExceptionApplicationList(t, server, "5001")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		want := []string{
			"GET /policies/v2/exceptionApplicationLists/5001?",
			"GET /policies/v2/exceptionApplicationLists/5001/applications?page=1&limit=1000",
		}
		if got := server.requests(); strings.Join(got, "|") != strings.Join(want, "|") {
			t.Errorf("requests = %q, want %q (one page of limit 1000: a call without limit is cut at 100)", got, want)
		}
		if state.Id.ValueString() != "5001" || state.Name.ValueString() != "my_do_not_decrypt_list" {
			t.Errorf("id = %q, name = %q", state.Id.ValueString(), state.Name.ValueString())
		}
		if state.AllowedCount.ValueInt64() != 125 || state.BlockedCount.IsNull() || state.BlockedCount.ValueInt64() != 0 || state.WarnedCount.IsNull() || state.WarnedCount.ValueInt64() != 0 {
			t.Errorf("counts = %s %s %s, want 125 0 0", state.AllowedCount, state.BlockedCount, state.WarnedCount)
		}
		if len(state.Applications) != 125 {
			t.Fatalf("applications = %d, want 125", len(state.Applications))
		}
		first, last := state.Applications[0], state.Applications[124]
		if first.ApplicationId.ValueInt64() != 1000 || first.ApplicationName.ValueString() != "App 0" || first.GroupId.ValueInt64() != 40 || first.GroupName.ValueString() != "Group 0" || !first.HasUrls.ValueBool() || first.Action.ValueString() != "no_decrypt" {
			t.Errorf("first application = %+v", first)
		}
		if last.ApplicationId.ValueInt64() != 1000+124*3 || last.ApplicationName.ValueString() != "App 124" || last.HasUrls.ValueBool() {
			t.Errorf("last application = %+v", last)
		}
	})

	t.Run("reads a list of more than 1000 applications page by page, in the order of the API", func(t *testing.T) {
		server := excListServer(t, excApps(2345))
		state, resp := readExceptionApplicationList(t, server, "5001")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if got := server.requests(); len(got) != 4 || got[1] != "GET /policies/v2/exceptionApplicationLists/5001/applications?page=1&limit=1000" || got[3] != "GET /policies/v2/exceptionApplicationLists/5001/applications?page=3&limit=1000" {
			t.Errorf("requests = %q, want the settings and three pages of 1000", got)
		}
		if len(state.Applications) != 2345 {
			t.Fatalf("applications = %d, want 2345", len(state.Applications))
		}
		for i, a := range state.Applications {
			if a.ApplicationId.ValueInt64() != int64(1000+i*3) {
				t.Fatalf("application %d has the id %d, want %d (a page was repeated, skipped or reordered)", i, a.ApplicationId.ValueInt64(), 1000+i*3)
			}
		}
	})

	t.Run("a list of exactly 1000 applications is one page, a list of 1001 is two", func(t *testing.T) {
		for n, wantRequests := range map[int]int{1000: 2, 1001: 3} {
			server := excListServer(t, excApps(n))
			state, resp := readExceptionApplicationList(t, server, "5001")
			if resp.Diagnostics.HasError() {
				t.Fatalf("%d applications: unexpected error: %v", n, resp.Diagnostics)
			}
			if len(state.Applications) != n || len(server.requests()) != wantRequests {
				t.Errorf("%d applications: read %d with %d requests, want %d requests", n, len(state.Applications), len(server.requests()), wantRequests)
			}
		}
	})

	t.Run("a list without applications has an empty list, not a null", func(t *testing.T) {
		server := excListServer(t, excApps(0))
		state, resp := readExceptionApplicationList(t, server, "5001")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if state.Applications == nil || len(state.Applications) != 0 {
			t.Errorf("applications = %#v, want an empty list", state.Applications)
		}
		if state.Name.ValueString() != "my_do_not_decrypt_list" {
			t.Errorf("name = %q", state.Name)
		}
	})

	t.Run("a null in a list item is read as the zero value, not as an error", func(t *testing.T) {
		// the generated fromBody() reads a JSON null as 0, "" or false (said in the YAML). Not seen on the int environment.
		apps := excApps(2)
		apps[1].(map[string]any)["groupId"] = nil
		server := excListServer(t, apps)
		state, resp := readExceptionApplicationList(t, server, "5001")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if len(state.Applications) != 2 || state.Applications[1].GroupId.ValueInt64() != 0 {
			t.Errorf("applications = %+v", state.Applications)
		}
	})

	t.Run("the pages of the applications are checked against meta.total", func(t *testing.T) {
		cases := []struct {
			name, want string
			apps       func(w http.ResponseWriter, r *http.Request)
		}{
			{"total grows while the pages are read", "`meta.total` changed from 1500 to 1501", func(w http.ResponseWriter, r *http.Request) {
				page := r.URL.Query().Get("page")
				total := 1500
				if page == "2" {
					total = 1501
				}
				items := make([]any, 0)
				for i := 0; i < 1000; i++ {
					items = append(items, excApp(i))
				}
				excWrite(w, 200, map[string]any{"meta": map[string]any{"page": 1, "limit": 1000, "total": total}, "data": items})
			}},
			{"the API ignores page, so the first page comes again", "page 2 is identical to page 1, the API ignores `page`", func(w http.ResponseWriter, r *http.Request) {
				items := make([]any, 0)
				for i := 0; i < 1000; i++ {
					items = append(items, excApp(i))
				}
				// 2000 is two full pages: the counts add up to the total, only the repeated page shows the fault
				excWrite(w, 200, map[string]any{"meta": map[string]any{"page": 1, "limit": 1000, "total": 2000}, "data": items})
			}},
			{"the API answers more items than meta.total", "read 2000 items, but `meta.total` is 1500", func(w http.ResponseWriter, r *http.Request) {
				items := make([]any, 0)
				for i := 0; i < 1000; i++ {
					n := i
					if r.URL.Query().Get("page") == "2" {
						n += 1000
					}
					items = append(items, excApp(n))
				}
				excWrite(w, 200, map[string]any{"meta": map[string]any{"page": 1, "limit": 1000, "total": 1500}, "data": items})
			}},
			{"a page is empty before the total is reached", "page 2 is empty, but only 1000 of 1500 items were read", func(w http.ResponseWriter, r *http.Request) {
				items := make([]any, 0)
				if r.URL.Query().Get("page") == "1" {
					for i := 0; i < 1000; i++ {
						items = append(items, excApp(i))
					}
				}
				excWrite(w, 200, map[string]any{"meta": map[string]any{"page": 1, "limit": 1000, "total": 1500}, "data": items})
			}},
			{"the API cuts the page at 100 although it was asked for 1000", "page 2 is empty, but only 100 of 125 items were read", func(w http.ResponseWriter, r *http.Request) {
				items := make([]any, 0)
				if r.URL.Query().Get("page") == "1" {
					for i := 0; i < 100; i++ {
						items = append(items, excApp(i))
					}
				}
				excWrite(w, 200, map[string]any{"meta": map[string]any{"page": 1, "limit": 100, "total": 125}, "data": items})
			}},
			{"meta.total is missing", "expected a number in `meta.total`, got: Null", func(w http.ResponseWriter, r *http.Request) {
				excWrite(w, 200, map[string]any{"data": []any{excApp(0)}})
			}},
			{"meta.total is a string", "expected a number in `meta.total`, got: String", func(w http.ResponseWriter, r *http.Request) {
				excWrite(w, 200, map[string]any{"meta": map[string]any{"total": "1"}, "data": []any{excApp(0)}})
			}},
			{"data is not a list", "expected a list in `data`, got: JSON", func(w http.ResponseWriter, r *http.Request) {
				excWrite(w, 200, map[string]any{"meta": map[string]any{"total": 1}, "data": map[string]any{"applications": []any{}}})
			}},
			{"a page holds more than the limit", "1001 items, more than the limit of 1000", func(w http.ResponseWriter, r *http.Request) {
				items := make([]any, 0)
				for i := 0; i < 1001; i++ {
					items = append(items, excApp(i))
				}
				excWrite(w, 200, map[string]any{"meta": map[string]any{"total": 1001}, "data": items})
			}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				// a read that does not stop at the first fault would otherwise run into the limit of 1000 pages, a thousand requests of 0.2 s
				defer func(old int) { exceptionApplicationListMaxPages = old }(exceptionApplicationListMaxPages)
				exceptionApplicationListMaxPages = 10
				server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
					if r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists/5001" {
						excWrite(w, http.StatusOK, excSettings(1))
						return
					}
					tc.apps(w, r)
				})
				state, resp := readExceptionApplicationList(t, server, "5001")
				detail := excErrorDetail(t, resp)
				if !strings.Contains(detail, "Failed to retrieve the applications of the object") || !strings.Contains(detail, tc.want) {
					t.Errorf("detail = %q, want the failed read of the applications and %q", detail, tc.want)
				}
				if !state.Id.IsNull() || state.Applications != nil {
					t.Errorf("state must stay empty after an error, got %+v", state)
				}
			})
		}
	})

	t.Run("an error on a later page is an error of the read, not a partial list", func(t *testing.T) {
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists/5001":
				excWrite(w, http.StatusOK, excSettings(2345))
			case r.URL.Query().Get("page") == "3":
				excWrite(w, http.StatusInternalServerError, map[string]any{"message": "internal"})
			default:
				excPaged(w, r, excApps(2345))
			}
		})
		state, resp := readExceptionApplicationList(t, server, "5001")
		detail := excErrorDetail(t, resp)
		if !strings.Contains(detail, "page 3") || !strings.Contains(detail, "StatusCode 500") {
			t.Errorf("detail = %q, want page 3 and the status 500", detail)
		}
		if !state.Id.IsNull() || state.Applications != nil {
			t.Errorf("state must stay empty after an error, got %+v", state)
		}
	})

	t.Run("the settings without an object in data are an error", func(t *testing.T) {
		for name, body := range map[string]map[string]any{
			"no data":          {"status": map[string]any{"code": 200}},
			"data is a list":   {"data": []any{}},
			"data is a string": {"data": "x"},
		} {
			server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) { excWrite(w, http.StatusOK, body) })
			_, resp := readExceptionApplicationList(t, server, "5001")
			if detail := excErrorDetail(t, resp); !strings.Contains(detail, "expected an object in `data`") {
				t.Errorf("%s: detail = %q", name, detail)
			}
			if got := server.requests(); len(got) != 1 {
				t.Errorf("%s: requests = %q, want no read of the applications", name, got)
			}
		}
	})

	errorCases := []struct {
		name, id string
		status   int
		body     string
	}{
		{"an id that does not exist is an error, not an empty data source", "999999999", http.StatusNotFound, `{"status":{"code":404,"text":"Not Found"},"data":{"error":"object_not_found","message":"DND exception application list not found"}}`},
		{"the id 0 is an error", "0", http.StatusBadRequest, `{"status":{"code":400,"text":"Bad Request"},"data":{"error":"invalid_request","message":{"exceptionapplicationlists_id":"invalid"}}}`},
		{"an id that is not a number is an error of the gateway", "abc", http.StatusNotFound, `{"message":"no Route matched with those values","request_id":"x"}`},
		{"an error of the API is an error", "5002", http.StatusInternalServerError, `{"message":"internal"}`},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
				excWrite(w, tc.status, json.RawMessage(tc.body))
			})
			state, resp := readExceptionApplicationList(t, server, tc.id)
			detail := excErrorDetail(t, resp)
			if !strings.Contains(detail, "Failed to retrieve object") || !strings.Contains(detail, "StatusCode "+strconv.Itoa(tc.status)) {
				t.Errorf("detail = %q, want the failed read and the status %d", detail, tc.status)
			}
			if got := server.requests(); len(got) != 1 || !strings.HasPrefix(got[0], "GET /policies/v2/exceptionApplicationLists/"+tc.id+"?") {
				t.Errorf("requests = %q, want one GET of the list, no read of the applications", got)
			}
			if !state.Id.IsNull() || !state.Name.IsNull() {
				t.Errorf("state must stay empty after an error, got %+v", state)
			}
		})
	}

	t.Run("the id is escaped, it cannot reach another path", func(t *testing.T) {
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			excWrite(w, http.StatusNotFound, map[string]any{"message": "not found"})
		})
		readExceptionApplicationList(t, server, "1/../../threatCategories")
		got := server.requests()
		if len(got) != 1 || !strings.HasPrefix(got[0], "GET /policies/v2/exceptionApplicationLists/1%2F") || strings.Contains(got[0], "Lists/1/") {
			t.Errorf("requests = %q, want one request below /policies/v2/exceptionApplicationLists/ with the slashes of the id escaped", got)
		}
	})
}

func TestApplicationExceptionListsDataSourceRead(t *testing.T) {
	idsServer := func(t *testing.T, ids []any) *excServer {
		t.Helper()
		return newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists" {
				excPaged(w, r, ids)
				return
			}
			excWrite(w, http.StatusNotFound, map[string]any{"message": "unexpected request"})
		})
	}
	idsOf := func(t *testing.T, state ApplicationExceptionLists) []int64 {
		t.Helper()
		var out []int64
		for _, e := range state.Ids.Elements() {
			var n int64
			if _, err := fmt.Sscan(strings.Trim(e.String(), `"`), &n); err != nil {
				t.Fatalf("id %s is not a number", e)
			}
			out = append(out, n)
		}
		sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
		return out
	}

	t.Run("reads the ids with limit 1000 in one request", func(t *testing.T) {
		server := idsServer(t, []any{16708649, 16708650, 17000001})
		state, resp := readExceptionApplicationLists(t, server)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if got := server.requests(); len(got) != 1 || got[0] != "GET /policies/v2/exceptionApplicationLists?page=1&limit=1000" {
			t.Errorf("requests = %q, want one page of limit 1000", got)
		}
		if got := idsOf(t, state); fmt.Sprint(got) != "[16708649 16708650 17000001]" {
			t.Errorf("ids = %v", got)
		}
	})

	t.Run("reads more than 1000 ids page by page", func(t *testing.T) {
		ids := make([]any, 1500)
		for i := range ids {
			ids[i] = 100000 + i
		}
		server := idsServer(t, ids)
		state, resp := readExceptionApplicationLists(t, server)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		got := idsOf(t, state)
		if len(got) != 1500 || got[0] != 100000 || got[1499] != 101499 || len(server.requests()) != 2 {
			t.Errorf("read %d ids with %d requests", len(got), len(server.requests()))
		}
	})

	t.Run("no list is an empty set, not a null", func(t *testing.T) {
		server := idsServer(t, []any{})
		state, resp := readExceptionApplicationLists(t, server)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if state.Ids.IsNull() || len(state.Ids.Elements()) != 0 {
			t.Errorf("ids = %s, want an empty set", state.Ids)
		}
	})

	t.Run("an id that is not a number is an error, it is not read as 0", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			bad  any
			want string
		}{
			{"a string", "x", "String"}, {"a bool", true, "True"}, {"null", nil, "Null"}, {"an object", map[string]any{"id": 1}, "JSON"}, {"a list", []any{1}, "JSON"},
		} {
			server := idsServer(t, []any{16708649, tc.bad})
			state, resp := readExceptionApplicationLists(t, server)
			if detail := excErrorDetail(t, resp); !strings.Contains(detail, "expected a number as id, got: "+tc.want) {
				t.Errorf("%s: detail = %q", tc.name, detail)
			}
			if !state.Ids.IsNull() {
				t.Errorf("%s: state must stay empty after an error, got %s", tc.name, state.Ids)
			}
		}
	})

	t.Run("a total that is not reached is an error", func(t *testing.T) {
		defer func(old int) { exceptionApplicationListMaxPages = old }(exceptionApplicationListMaxPages)
		exceptionApplicationListMaxPages = 10
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			data := []any{}
			if r.URL.Query().Get("page") == "1" {
				data = []any{16708649}
			}
			excWrite(w, http.StatusOK, map[string]any{"meta": map[string]any{"page": 1, "limit": 1000, "total": 3}, "data": data})
		})
		_, resp := readExceptionApplicationLists(t, server)
		if detail := excErrorDetail(t, resp); !strings.Contains(detail, "page 2 is empty, but only 1 of 3 items were read") || !strings.Contains(detail, "Failed to retrieve objects") {
			t.Errorf("detail = %q", detail)
		}
	})

	t.Run("an API that ignores page is an error, not a set with the same id three times", func(t *testing.T) {
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			excWrite(w, http.StatusOK, map[string]any{"meta": map[string]any{"page": 1, "limit": 1000, "total": 3}, "data": []any{16708649}})
		})
		state, resp := readExceptionApplicationLists(t, server)
		if detail := excErrorDetail(t, resp); !strings.Contains(detail, "page 2 is identical to page 1") {
			t.Errorf("detail = %q", detail)
		}
		if !state.Ids.IsNull() {
			t.Errorf("state must stay empty after an error, got %s", state.Ids)
		}
	})

	t.Run("an id that the API lists twice is an error", func(t *testing.T) {
		server := idsServer(t, []any{16708649, 16708650, 16708649})
		state, resp := readExceptionApplicationLists(t, server)
		if detail := excErrorDetail(t, resp); !strings.Contains(detail, "the id 16708649 is listed more than once") {
			t.Errorf("detail = %q", detail)
		}
		if !state.Ids.IsNull() {
			t.Errorf("state must stay empty after an error, got %s", state.Ids)
		}
	})

	t.Run("an error of the API is an error", func(t *testing.T) {
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			excWrite(w, http.StatusForbidden, map[string]any{"message": "Access denied"})
		})
		state, resp := readExceptionApplicationLists(t, server)
		if detail := excErrorDetail(t, resp); !strings.Contains(detail, "StatusCode 403") || !strings.Contains(detail, "page 1") {
			t.Errorf("detail = %q", detail)
		}
		if !state.Ids.IsNull() {
			t.Errorf("state must stay empty after an error, got %s", state.Ids)
		}
	})
}

func TestReadExceptionApplicationListPagesStops(t *testing.T) {
	// an API that answers one item per page and a total that is never reached must not keep the read busy for ever
	defer func(old int) { exceptionApplicationListMaxPages = old }(exceptionApplicationListMaxPages)
	exceptionApplicationListMaxPages = 7
	var mu sync.Mutex
	requests := 0
	server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests++
		mu.Unlock()
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		excWrite(w, http.StatusOK, map[string]any{"meta": map[string]any{"total": 5000}, "data": []any{page}})
	})
	_, err := readExceptionApplicationListPages(context.Background(), newExcClient(t, server), "/policies/v2/exceptionApplicationLists")
	if err == nil || !strings.Contains(err.Error(), "the total of 5000 was not reached") {
		t.Fatalf("err = %v, want the stop at the maximum number of pages", err)
	}
	if requests != exceptionApplicationListMaxPages {
		t.Errorf("requests = %d, want %d", requests, exceptionApplicationListMaxPages)
	}
}

func TestReadExceptionApplicationListPagesRejectsAnImpossibleTotal(t *testing.T) {
	// a total that more pages than the maximum would be needed for is an error on the first page, not after the maximum number of requests
	defer func(old int) { exceptionApplicationListMaxPages = old }(exceptionApplicationListMaxPages)
	exceptionApplicationListMaxPages = 7
	for name, total := range map[string]int64{"7001 items": 7001, "five billion": 5000000000, "the largest int64": 9223372036854775807} {
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			excWrite(w, http.StatusOK, map[string]any{"meta": map[string]any{"total": total}, "data": []any{1}})
		})
		_, err := readExceptionApplicationListPages(context.Background(), newExcClient(t, server), "/policies/v2/exceptionApplicationLists")
		if err == nil || !strings.Contains(err.Error(), "page 1: `meta.total` is "+strconv.FormatInt(total, 10)+", more than 7 pages of 1000 items") {
			t.Errorf("%s: err = %v", name, err)
		}
		if got := server.requests(); len(got) != 1 {
			t.Errorf("%s: requests = %q, want the first page only", name, got)
		}
	}
	// the largest total that fits is read
	server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
		excPaged(w, r, excApps(7000))
	})
	out, err := readExceptionApplicationListPages(context.Background(), newExcClient(t, server), "/policies/v2/exceptionApplicationLists")
	if err != nil || len(server.requests()) != 7 || strings.Count(out, `"applicationId"`) != 7000 {
		t.Errorf("7000 items in 7 pages: err = %v, requests = %d", err, len(server.requests()))
	}
}

func TestApplicationExceptionListDataSourceReadDetails(t *testing.T) {
	t.Run("the id of the state is the id of the config, not the id in the answer", func(t *testing.T) {
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists/5001" {
				settings := excSettings(1)
				settings["data"].(map[string]any)["id"] = 7777
				excWrite(w, http.StatusOK, settings)
				return
			}
			excPaged(w, r, excApps(1))
		})
		state, resp := readExceptionApplicationList(t, server, "5001")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if state.Id.ValueString() != "5001" {
			t.Errorf("id = %q, want the id of the config", state.Id.ValueString())
		}
	})

	t.Run("every count and every field of an application is read from its own key", func(t *testing.T) {
		apps := []any{map[string]any{
			"applicationId": 11, "applicationName": "name", "groupId": 22, "groupName": "group", "hasUrls": true, "action": "action_x",
			"selectedAction": map[string]any{"name": "name_y", "displayName": "display_z", "actionId": 33},
		}}
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists/5001" {
				settings := excSettings(1)
				d := settings["data"].(map[string]any)
				d["allowedCount"], d["blockedCount"], d["warnedCount"] = 4, 2, 3
				excWrite(w, http.StatusOK, settings)
				return
			}
			excPaged(w, r, apps)
		})
		state, resp := readExceptionApplicationList(t, server, "5001")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if state.AllowedCount.ValueInt64() != 4 || state.BlockedCount.ValueInt64() != 2 || state.WarnedCount.ValueInt64() != 3 {
			t.Errorf("counts allowed=%s blocked=%s warned=%s, want 4 2 3", state.AllowedCount, state.BlockedCount, state.WarnedCount)
		}
		a := state.Applications[0]
		if a.ApplicationId.ValueInt64() != 11 || a.ApplicationName.ValueString() != "name" || a.GroupId.ValueInt64() != 22 || a.GroupName.ValueString() != "group" || !a.HasUrls.ValueBool() || a.Action.ValueString() != "action_x" {
			t.Errorf("application = %+v, want every field from its own key (action is not selectedAction.name)", a)
		}
	})

	t.Run("a field that the API leaves out is null, not zero", func(t *testing.T) {
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists/5001" {
				excWrite(w, http.StatusOK, map[string]any{"data": map[string]any{"id": 5001}})
				return
			}
			excPaged(w, r, []any{map[string]any{"applicationId": 11}})
		})
		state, resp := readExceptionApplicationList(t, server, "5001")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if !state.Name.IsNull() || !state.AllowedCount.IsNull() || !state.BlockedCount.IsNull() || !state.WarnedCount.IsNull() {
			t.Errorf("settings = %s %s %s %s, want null", state.Name, state.AllowedCount, state.BlockedCount, state.WarnedCount)
		}
		a := state.Applications[0]
		if a.ApplicationId.ValueInt64() != 11 || !a.ApplicationName.IsNull() || !a.GroupId.IsNull() || !a.GroupName.IsNull() || !a.HasUrls.IsNull() || !a.Action.IsNull() {
			t.Errorf("application = %+v, want only the id set", a)
		}
	})
}

func TestReadExceptionApplicationListPagesDetails(t *testing.T) {
	t.Run("the default limit of pages is large enough for the largest list and small enough to stop", func(t *testing.T) {
		if exceptionApplicationListMaxPages != 1000 || exceptionApplicationListPageSize != 1000 {
			t.Errorf("max pages = %d, page size = %d, want 1000 and 1000", exceptionApplicationListMaxPages, exceptionApplicationListPageSize)
		}
	})

	t.Run("request modifiers reach every request", func(t *testing.T) {
		var mu sync.Mutex
		seen := map[string]int{}
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			seen[r.Header.Get("X-Test-Mod")]++
			mu.Unlock()
			excPaged(w, r, excApps(1500))
		})
		mod := func(req *secureaccess.Req) { req.HttpReq.Header.Set("X-Test-Mod", "yes") }
		out, err := readExceptionApplicationListPages(context.Background(), newExcClient(t, server), "/policies/v2/exceptionApplicationLists/1/applications", mod)
		if err != nil || strings.Count(out, `"applicationId"`) != 1500 {
			t.Fatalf("err = %v", err)
		}
		if seen["yes"] != 2 || len(seen) != 1 {
			t.Errorf("requests by value of the header = %v, want both pages with the modifier", seen)
		}
	})

	t.Run("an error returns no items and keeps the error of the client", func(t *testing.T) {
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("page") == "2" {
				excWrite(w, http.StatusInternalServerError, map[string]any{"message": "internal"})
				return
			}
			excPaged(w, r, excApps(1500))
		})
		out, err := readExceptionApplicationListPages(context.Background(), newExcClient(t, server), "/policies/v2/exceptionApplicationLists/1/applications")
		if err == nil || out != "" {
			t.Fatalf("out = %q, err = %v, want no items and an error", out, err)
		}
		if errors.Unwrap(err) == nil || !strings.Contains(err.Error(), "page 2") {
			t.Errorf("err = %v, want the error of the client wrapped, with the page", err)
		}
	})

	t.Run("a total that shrinks while the pages are read is an error", func(t *testing.T) {
		server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
			total := 1500
			if r.URL.Query().Get("page") == "2" {
				total = 1400
			}
			items := make([]any, 0)
			for i := 0; i < 1000; i++ {
				items = append(items, excApp(i+1000*(map[bool]int{true: 1, false: 0}[r.URL.Query().Get("page") == "2"])))
			}
			excWrite(w, http.StatusOK, map[string]any{"meta": map[string]any{"total": total}, "data": items})
		})
		out, err := readExceptionApplicationListPages(context.Background(), newExcClient(t, server), "/policies/v2/exceptionApplicationLists/1/applications")
		if err == nil || !strings.Contains(err.Error(), "`meta.total` changed from 1500 to 1400") || out != "" {
			t.Errorf("out has %d bytes, err = %v", len(out), err)
		}
	})
}

func TestApplicationExceptionListTotalChangesOnAnyPage(t *testing.T) {
	// the items are always the n applications of the list, only meta.total moves, on page 2 or on page 3
	for _, tc := range []struct {
		name          string
		n, newTotal   int
		changedOnPage string
	}{
		{"shrinks on page 2 of 2", 1500, 1499, "2"},
		{"grows on page 2 of 2", 1500, 1501, "2"},
		{"shrinks on page 3 of 3", 2345, 2344, "3"},
		{"grows on page 3 of 3", 2345, 2346, "3"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func(old int) { exceptionApplicationListMaxPages = old }(exceptionApplicationListMaxPages)
			exceptionApplicationListMaxPages = 10
			all := excApps(tc.n)
			server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists/5001" {
					excWrite(w, http.StatusOK, excSettings(tc.n))
					return
				}
				page := r.URL.Query().Get("page")
				total := tc.n
				if page == tc.changedOnPage {
					total = tc.newTotal
				}
				items := []any{}
				p, _ := strconv.Atoi(page)
				if lo := (p - 1) * 1000; lo < tc.n {
					hi := lo + 1000
					if hi > tc.n {
						hi = tc.n
					}
					items = all[lo:hi]
				}
				excWrite(w, http.StatusOK, map[string]any{"meta": map[string]any{"total": total}, "data": items})
			})
			_, resp := readExceptionApplicationList(t, server, "5001")
			detail := excErrorDetail(t, resp)
			want := fmt.Sprintf("page %s: `meta.total` changed from %d to %d", tc.changedOnPage, tc.n, tc.newTotal)
			if !strings.Contains(detail, want) {
				t.Errorf("detail = %q, want %q", detail, want)
			}
		})
	}
}

func TestApplicationExceptionListApplicationsMustBeObjects(t *testing.T) {
	for name, bad := range map[string]any{"a string": "x", "null": nil, "a number": 5, "a list": []any{1}} {
		t.Run(name, func(t *testing.T) {
			server := newExcServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.EscapedPath() == "/policies/v2/exceptionApplicationLists/5001" {
					excWrite(w, http.StatusOK, excSettings(2))
					return
				}
				excPaged(w, r, []any{excApp(0), bad})
			})
			_, resp := readExceptionApplicationList(t, server, "5001")
			if detail := excErrorDetail(t, resp); !strings.Contains(detail, "expected an object as application 1") {
				t.Errorf("detail = %q", detail)
			}
		})
	}
}
