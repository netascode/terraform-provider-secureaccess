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

// This file is not generated. It tests the parts of the read of the applications and of the exception application list id that the
// tests of model_secureaccess_do_not_decrypt_list_test.go leave out: the guards of the id, the fallback of the data source for an object found
// by its name, a state without applications, the answers to create and update, and a failing read by id.

package provider

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
)

// The id of the exception application list is a number of 1 or more, a string or a bool is not an id (it would be put into the path).
func TestDoNotDecryptListApplicationListId(t *testing.T) {
	for object, want := range map[string]string{
		`{"exceptionapplicationlistId":888}`: "888",
		`{"exceptionapplicationlistId":1}`:   "1",
	} {
		if got, err := doNotDecryptListApplicationListId(gjson.Parse(object)); err != nil || got != want {
			t.Errorf("%s: got %q, %v, want %q", object, got, err, want)
		}
	}
	for _, object := range []string{
		`{"exceptionapplicationlistId":"888","totalApplications":2}`, // a string would be sent as "%22888%22"
		`{"exceptionapplicationlistId":true,"totalApplications":2}`,  // Int() of true is 1
		`{"exceptionapplicationlistId":0,"totalApplications":2}`,
	} {
		if got, err := doNotDecryptListApplicationListId(gjson.Parse(object)); err == nil {
			t.Errorf("%s: got %q, want an error", object, got)
		}
	}
}

// dndPartialListServer serves the Do Not Decrypt Lists API with a list item that carries only some of the fields
// that point to the lists behind it, and the lists behind it.
func dndPartialListServer(t *testing.T, listItem string) (*secureaccess.Client, *dndReadAPI) {
	t.Helper()
	api := &dndReadAPI{}
	const full = `{"id":5,"name":"n","isDefault":false,"exceptiondomainlistId":777,"exceptionapplicationlistId":888,"totalDomains":2,"totalIps":0,"totalApplications":3}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/v2/token" {
			api.mu.Lock()
			api.paths = append(api.paths, r.URL.Path)
			api.mu.Unlock()
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/auth/v2/token":
			_, _ = io.WriteString(w, `{"access_token":"t","token_type":"bearer","expires_in":3600}`)
		case "/policies/v2/doNotDecryptLists":
			_, _ = io.WriteString(w, "["+listItem+"]")
		case "/policies/v2/doNotDecryptLists/5":
			_, _ = io.WriteString(w, full)
		case "/policies/v2/destinationlists/777/destinations":
			_, _ = io.WriteString(w, `{"data":[{"id":1,"destination":"a.com","type":"domain"},{"id":2,"destination":"c.com","type":"domain"}],"meta":{"total":2}}`)
		case "/policies/v2/exceptionApplicationLists/888/applications":
			excPaged(w, r, excApps(3))
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

// An object found by name is read again by its id if the list item lacks the destination list or the field of the exception application
// list, and is not if it carries both.
func TestDoNotDecryptListDataSourceReadByNameFallback(t *testing.T) {
	ctx := context.Background()
	const lists = "/policies/v2/doNotDecryptLists"
	const byId = "/policies/v2/doNotDecryptLists/5"
	const dest = "/policies/v2/destinationlists/777/destinations"
	const apps = "/policies/v2/exceptionApplicationLists/888/applications"
	for name, tc := range map[string]struct {
		item string
		want []string
	}{
		"the item lacks exceptionapplicationlistId": {
			`{"id":5,"name":"n","isDefault":false,"exceptiondomainlistId":777,"totalDomains":2,"totalIps":0,"totalApplications":3}`,
			[]string{lists, byId, dest, apps}},
		"the item lacks exceptiondomainlistId": {
			`{"id":5,"name":"n","isDefault":false,"exceptionapplicationlistId":888,"totalDomains":2,"totalIps":0,"totalApplications":3}`,
			[]string{lists, byId, dest, apps}},
		"the item carries both": {
			`{"id":5,"name":"n","isDefault":false,"exceptiondomainlistId":777,"exceptionapplicationlistId":888,"totalDomains":2,"totalIps":0,"totalApplications":3}`,
			[]string{lists, dest, apps}},
	} {
		t.Run(name, func(t *testing.T) {
			client, api := dndPartialListServer(t, tc.item)
			d := &DoNotDecryptListDataSource{client: client}
			var schemaResp datasource.SchemaResponse
			d.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
			config := DoNotDecryptList{Id: types.StringNull(), Name: types.StringValue("n"), Applications: types.SetNull(types.Int64Type), DecryptExceptionCategories: types.SetNull(types.Int64Type)}
			nullRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)
			configState := tfsdk.State{Schema: schemaResp.Schema, Raw: nullRaw}
			if diags := configState.Set(ctx, &config); diags.HasError() {
				t.Fatal(diags)
			}
			resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: nullRaw}}
			d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: configState.Raw}}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatal(resp.Diagnostics)
			}
			var got DoNotDecryptList
			if diags := resp.State.Get(ctx, &got); diags.HasError() {
				t.Fatal(diags)
			}
			if ids := helpers.Int64SetElements(ctx, got.Applications); !slices.Equal(ids, []int64{1000, 1003, 1006}) || got.ApplicationExceptionListId.ValueInt64() != 888 {
				t.Errorf("applications = %v, application_exception_list_id = %v", ids, got.ApplicationExceptionListId)
			}
			if reqs := api.requests(); !slices.Equal(reqs, tc.want) {
				t.Errorf("requests = %v, want %v", reqs, tc.want)
			}
		})
	}
}

// A refresh of a state whose `applications` is null (the configuration omits them) leaves them null instead of reporting the
// applications of the list as a change.
func TestDoNotDecryptListResourceReadApplicationsNullInState(t *testing.T) {
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/auth/v2/token":
			_, _ = io.WriteString(w, `{"access_token":"t","token_type":"bearer","expires_in":3600}`)
		case "/policies/v2/doNotDecryptLists/5":
			_, _ = io.WriteString(w, `{"id":5,"name":"n","isDefault":false,"exceptiondomainlistId":777,"exceptionapplicationlistId":888,"totalDomains":1,"totalIps":0,"totalApplications":3}`)
		case "/policies/v2/destinationlists/777/destinations":
			_, _ = io.WriteString(w, `{"data":[{"id":1,"destination":"a.com","type":"domain"}],"meta":{"total":1}}`)
		case "/policies/v2/exceptionApplicationLists/888/applications":
			excPaged(w, r, excApps(3))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	s := newDndReadServer(t, server.URL)

	sch := dndTestSchema(t)
	model := dndTestModel(types.SetNull(types.Int64Type), dndInt64Set(2))
	model.Id = types.StringValue("5")
	model.Destinations = []DoNotDecryptListDestinations{dndDestination("a.com", "")}
	st := dndTestState(t, sch, &model)
	dv, err := tfprotov6.NewDynamicValue(sch.Type().TerraformType(context.Background()), st.Raw)
	if err != nil {
		t.Fatal(err)
	}
	rr := s.read(&dv, nil)
	if len(rr.Diagnostics) > 0 {
		t.Fatal(rr.Diagnostics)
	}
	if attr := s.attribute(rr.NewState, "applications"); !attr.IsNull() {
		t.Errorf("applications = %s after a refresh of a state without applications, want null", attr)
	}
	if got := s.applicationListId(rr.NewState); got != "888" {
		t.Errorf("application_exception_list_id = %s, want 888", got)
	}
}

// The answers to the create and update requests carry exceptionapplicationlistId null (no application yet) or 0 (only content
// categories); fromBodyUnknowns() reads both as 0, like fromBody() and fromBodyPartial(), so that `application_exception_list_id > 0`
// is usable right after the apply.
func TestDoNotDecryptListFromBodyUnknownsApplicationExceptionListId(t *testing.T) {
	ctx := context.Background()
	for body, want := range map[string]int64{
		`{"exceptionapplicationlistId":null}`: 0,
		`{"exceptionapplicationlistId":0}`:    0,
		`{"exceptionapplicationlistId":888}`:  888,
	} {
		data := DoNotDecryptList{ApplicationExceptionListId: types.Int64Unknown()}
		data.fromBodyUnknowns(ctx, gjson.Parse(body))
		if data.ApplicationExceptionListId.IsNull() || data.ApplicationExceptionListId.ValueInt64() != want {
			t.Errorf("%s: application_exception_list_id = %v, want %d", body, data.ApplicationExceptionListId, want)
		}
	}
}

// A failing GET by id of the data source is one error that says so, it does not go on and add the misleading error of an object
// without lists.
func TestDoNotDecryptListDataSourceReadByIdFailure(t *testing.T) {
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/auth/v2/token":
			_, _ = io.WriteString(w, `{"access_token":"t","token_type":"bearer","expires_in":3600}`)
		case "/policies/v2/doNotDecryptLists/5":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"message":"no"}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	client, err := secureaccess.NewClient(server.URL, "key", "secret")
	if err != nil {
		t.Fatal(err)
	}
	d := &DoNotDecryptListDataSource{client: client}
	var schemaResp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	config := DoNotDecryptList{Id: types.StringValue("5"), Applications: types.SetNull(types.Int64Type), DecryptExceptionCategories: types.SetNull(types.Int64Type)}
	nullRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)
	configState := tfsdk.State{Schema: schemaResp.Schema, Raw: nullRaw}
	if diags := configState.Set(ctx, &config); diags.HasError() {
		t.Fatal(diags)
	}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: nullRaw}}
	d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: configState.Raw}}, &resp)
	if resp.Diagnostics.ErrorsCount() != 1 || !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), "Failed to retrieve object") {
		t.Errorf("diagnostics = %v, want one error: Failed to retrieve object", resp.Diagnostics)
	}
}
