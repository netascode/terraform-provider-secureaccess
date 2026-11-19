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
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/netascode/go-secureaccess"
	"github.com/tidwall/gjson"
)

// This file is not generated. It tests the data source security_profile_policy_setting with the payload that the Secure
// Access API returns for the Default Policy (only the ids and the dates are made up), and with the other states of the
// settings that the Security Profiles API guide describes. The cases for ai_semantic_inspection are synthetic: the guide does
// not describe that setting and it was never seen enabled.

// A Policy Setting as returned by GET /policies/v2/policysettings/{id}, every setting is off.
const policySettingDefaultPayload = `{
  "id": 11111111,
  "organizationId": 2222222,
  "isDefault": null,
  "name": "Default Policy",
  "intelligentProxy": false,
  "intelligentProxyHttps": false,
  "upIntelligentProxy": false,
  "upIntelligentProxyHttps": false,
  "ipFiltering": false,
  "fileInspectionId": null,
  "reporting": "enabled",
  "whitelistOnly": false,
  "safeSearch": false,
  "createdAt": "2025-01-09 14:38:57",
  "modifiedAt": "2025-01-09 14:38:57",
  "markedForDeletion": false,
  "decryptExceptionCategories": "0",
  "exceptiondomainlistId": 0,
  "saml": 0,
  "swgDisplayBlockPage": 0,
  "aiSupplyChainCategories": "0",
  "mcpSemanticInspection": false,
  "httpHeaderConstraints": false,
  "httpHeaderContentLength": 0,
  "cdrDownload": false,
  "categoryIds": [],
  "categoryNames": []
}`

// The mask of the AI supply chain categories that the guide gives for all three categories (54 characters).
const policySettingAllAiCategories = "E00000000000000000000000000000000000000000000000000000"

// policySettingBody returns a body with the keys that the data source reads, everything off, and the given keys replaced.
func policySettingBody(t *testing.T, over map[string]any) string {
	t.Helper()
	body := map[string]any{
		"name": "p", "intelligentProxyHttps": false, "intelligentProxy": false, "swgDisplayBlockPage": 0,
		"saml": 0, "safeSearch": false, "aiSupplyChainCategories": "0", "mcpSemanticInspection": false, "reporting": "enabled",
	}
	for k, v := range over {
		body[k] = v
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b)
}

func TestSecurityProfilePolicySettingFromBody(t *testing.T) {
	type want struct {
		name, ai, reporting                       string
		decryption, proxy, blockPage, sso, search bool
		semantic                                  bool
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "default policy, everything off",
			body: policySettingDefaultPayload,
			want: want{name: "Default Policy", ai: "0", reporting: "enabled"},
		},
		// Every boolean is turned on alone, so a value read from the wrong key shows up in at least one case.
		{"only decryption", policySettingBody(t, map[string]any{"intelligentProxyHttps": true}), want{name: "p", ai: "0", reporting: "enabled", decryption: true}},
		{"only intelligent proxy", policySettingBody(t, map[string]any{"intelligentProxy": true}), want{name: "p", ai: "0", reporting: "enabled", proxy: true}},
		{"only block page mode", policySettingBody(t, map[string]any{"swgDisplayBlockPage": 1}), want{name: "p", ai: "0", reporting: "enabled", blockPage: true}},
		{"only single sign-on", policySettingBody(t, map[string]any{"saml": 1}), want{name: "p", ai: "0", reporting: "enabled", sso: true}},
		{"only safe search", policySettingBody(t, map[string]any{"safeSearch": true}), want{name: "p", ai: "0", reporting: "enabled", search: true}},
		{"only ai semantic inspection", policySettingBody(t, map[string]any{"mcpSemanticInspection": true}), want{name: "p", ai: "0", reporting: "enabled", semantic: true}},
		{
			// the AI supply chain mask is a different setting from AI Semantic Inspection, one must not turn on the other
			name: "only an AI supply chain category",
			body: policySettingBody(t, map[string]any{"aiSupplyChainCategories": policySettingAllAiCategories}),
			want: want{name: "p", ai: policySettingAllAiCategories, reporting: "enabled"},
		},
		{
			// the guide: Decryption sets intelligentProxy and intelligentProxyHttps to true and swgDisplayBlockPage to 0
			name: "the state of Decryption in the guide",
			body: policySettingBody(t, map[string]any{"intelligentProxy": true, "intelligentProxyHttps": true}),
			want: want{name: "p", ai: "0", reporting: "enabled", decryption: true, proxy: true},
		},
		{
			// the guide: Block page mode sets intelligentProxy to true, intelligentProxyHttps to false and swgDisplayBlockPage to 1
			name: "the state of Block page mode in the guide",
			body: policySettingBody(t, map[string]any{"intelligentProxy": true, "swgDisplayBlockPage": 1}),
			want: want{name: "p", ai: "0", reporting: "enabled", proxy: true, blockPage: true},
		},
		{
			name: "everything on, other strings",
			body: policySettingBody(t, map[string]any{"name": "all", "intelligentProxyHttps": true, "intelligentProxy": true, "swgDisplayBlockPage": 1, "saml": 1, "safeSearch": true, "aiSupplyChainCategories": policySettingAllAiCategories, "mcpSemanticInspection": true, "reporting": "disabled"}),
			want: want{name: "all", ai: policySettingAllAiCategories, reporting: "disabled", decryption: true, proxy: true, blockPage: true, sso: true, search: true, semantic: true},
		},
		{"a number other than 1 is true, as the data source documents", policySettingBody(t, map[string]any{"saml": 2, "swgDisplayBlockPage": 2}), want{name: "p", ai: "0", reporting: "enabled", sso: true, blockPage: true}},
		{
			// documented in the definition: the generated fromBody() reads a JSON null as an empty string and as false
			name: "a JSON null is read as an empty string and as false",
			body: policySettingBody(t, map[string]any{"name": nil, "reporting": nil, "aiSupplyChainCategories": nil, "intelligentProxy": nil, "saml": nil, "mcpSemanticInspection": nil}),
			want: want{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got SecurityProfilePolicySetting
			got.fromBody(context.Background(), gjson.Parse(tt.body))
			if !got.Id.IsNull() {
				t.Errorf("id = %s, fromBody must not set it, it comes from the config", got.Id)
			}
			for _, c := range []struct {
				attr      string
				got, want attr.Value
			}{
				{"name", got.Name, types.StringValue(tt.want.name)},
				{"ai_supply_chain_categories", got.AiSupplyChainCategories, types.StringValue(tt.want.ai)},
				{"reporting", got.Reporting, types.StringValue(tt.want.reporting)},
				{"decryption_enabled", got.DecryptionEnabled, types.BoolValue(tt.want.decryption)},
				{"intelligent_proxy", got.IntelligentProxy, types.BoolValue(tt.want.proxy)},
				{"end_user_notification", got.EndUserNotification, types.BoolValue(tt.want.blockPage)},
				{"single_sign_on", got.SingleSignOn, types.BoolValue(tt.want.sso)},
				{"safe_search", got.SafeSearch, types.BoolValue(tt.want.search)},
				{"ai_semantic_inspection", got.AiSemanticInspection, types.BoolValue(tt.want.semantic)},
			} {
				if !c.got.Equal(c.want) {
					t.Errorf("%s = %s, want %s", c.attr, c.got, c.want)
				}
			}
		})
	}

	t.Run("keys that are missing are null, not false and not empty", func(t *testing.T) {
		var got SecurityProfilePolicySetting
		got.fromBody(context.Background(), gjson.Parse(`{}`))
		for attr, null := range map[string]bool{
			"id": got.Id.IsNull(), "name": got.Name.IsNull(), "ai_supply_chain_categories": got.AiSupplyChainCategories.IsNull(),
			"reporting": got.Reporting.IsNull(), "decryption_enabled": got.DecryptionEnabled.IsNull(),
			"intelligent_proxy": got.IntelligentProxy.IsNull(), "end_user_notification": got.EndUserNotification.IsNull(),
			"single_sign_on": got.SingleSignOn.IsNull(), "safe_search": got.SafeSearch.IsNull(), "ai_semantic_inspection": got.AiSemanticInspection.IsNull(),
		} {
			if !null {
				t.Errorf("%s is not null", attr)
			}
		}
	})
}

// The definition says that a Policy Setting is found by its id only and that the list is never read (the names are not
// unique), and the attributes are the ones documented. Both would change silently with a lookup by name or a new attribute.
func TestSecurityProfilePolicySettingSchema(t *testing.T) {
	ctx := context.Background()
	ds := &SecurityProfilePolicySettingDataSource{}

	var meta datasource.MetadataResponse
	ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "secureaccess"}, &meta)
	if meta.TypeName != "secureaccess_security_profile_policy_setting" {
		t.Errorf("type name = %q", meta.TypeName)
	}

	var resp datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &resp)
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
	want := "ai_semantic_inspection,ai_supply_chain_categories,decryption_enabled,end_user_notification,id,intelligent_proxy,name,reporting,safe_search,single_sign_on"
	if got := strings.Join(names, ","); got != want {
		t.Errorf("attributes = %s\nwant       %s", got, want)
	}
	// the generator adds config validators (id or name) with a lookup by name, there is no second way to find a Policy Setting
	if _, ok := any(ds).(datasource.DataSourceWithConfigValidators); ok {
		t.Errorf("there must be no config validators")
	}

	var registered bool
	for _, factory := range New("test")().DataSources(ctx) {
		var m datasource.MetadataResponse
		factory().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "secureaccess"}, &m)
		if m.TypeName == "secureaccess_security_profile_policy_setting" {
			registered = true
		}
	}
	if !registered {
		t.Error("the data source is not registered in the provider")
	}
}

// policySettingServer is a fake Secure Access API: it answers the token request and everything else with the handler,
// and records the method, the escaped path and the query of every other request.
type policySettingServer struct {
	*httptest.Server
	mu    sync.Mutex
	paths []string
}

func newPolicySettingServer(t *testing.T, handle func(w http.ResponseWriter, r *http.Request)) *policySettingServer {
	t.Helper()
	s := &policySettingServer{}
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

func (s *policySettingServer) requests() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.paths...)
}

// readPolicySetting runs the Read() of the data source with the given id in the config.
func readPolicySetting(t *testing.T, server *policySettingServer, id string) (SecurityProfilePolicySetting, datasource.ReadResponse) {
	t.Helper()
	ctx := context.Background()
	client, err := secureaccess.NewClient(server.URL, "key", "secret", secureaccess.MaxRetries(0))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	ds := &SecurityProfilePolicySettingDataSource{client: client}

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

	var state SecurityProfilePolicySetting
	if !resp.Diagnostics.HasError() {
		if diags := resp.State.Get(ctx, &state); diags.HasError() {
			t.Fatalf("reading the state: %v", diags)
		}
	}
	return state, resp
}

func TestSecurityProfilePolicySettingDataSourceRead(t *testing.T) {
	t.Run("reads the policy setting with one request by id, without reading the list", func(t *testing.T) {
		// the id of the config and the id in the payload are different, the id of the state is the one of the config
		server := newPolicySettingServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.EscapedPath() == "/policies/v2/policysettings/33333333" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(policySettingDefaultPayload))
				return
			}
			http.Error(w, `{"message":"unexpected request"}`, http.StatusNotFound)
		})
		state, resp := readPolicySetting(t, server, "33333333")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if got := server.requests(); len(got) != 1 || got[0] != "GET /policies/v2/policysettings/33333333?" {
			t.Errorf("requests = %q, want exactly one GET of the policy setting by id and no read of the list", got)
		}
		if state.Id.ValueString() != "33333333" {
			t.Errorf("id = %q, want the id of the config", state.Id.ValueString())
		}
		if state.Name.ValueString() != "Default Policy" || state.Reporting.ValueString() != "enabled" || state.AiSupplyChainCategories.ValueString() != "0" {
			t.Errorf("state = %+v", state)
		}
		for attr, v := range map[string]types.Bool{"decryption_enabled": state.DecryptionEnabled, "intelligent_proxy": state.IntelligentProxy, "end_user_notification": state.EndUserNotification, "single_sign_on": state.SingleSignOn, "safe_search": state.SafeSearch, "ai_semantic_inspection": state.AiSemanticInspection} {
			if v.IsNull() || v.ValueBool() {
				t.Errorf("%s = %s, want false", attr, v)
			}
		}
	})

	t.Run("reads the settings that are on", func(t *testing.T) {
		server := newPolicySettingServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(policySettingBody(t, map[string]any{"intelligentProxy": true, "safeSearch": true, "saml": 1, "mcpSemanticInspection": true})))
		})
		state, resp := readPolicySetting(t, server, "44444444")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		got := map[string]bool{"decryption_enabled": state.DecryptionEnabled.ValueBool(), "intelligent_proxy": state.IntelligentProxy.ValueBool(), "end_user_notification": state.EndUserNotification.ValueBool(), "single_sign_on": state.SingleSignOn.ValueBool(), "safe_search": state.SafeSearch.ValueBool(), "ai_semantic_inspection": state.AiSemanticInspection.ValueBool()}
		want := map[string]bool{"decryption_enabled": false, "intelligent_proxy": true, "end_user_notification": false, "single_sign_on": true, "safe_search": true, "ai_semantic_inspection": true}
		for attr := range want {
			if got[attr] != want[attr] {
				t.Errorf("%s = %t, want %t", attr, got[attr], want[attr])
			}
		}
	})

	t.Run("reads AI Semantic Inspection alone, without turning on another setting", func(t *testing.T) {
		server := newPolicySettingServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(policySettingBody(t, map[string]any{"mcpSemanticInspection": true})))
		})
		state, resp := readPolicySetting(t, server, "66666666")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		got := map[string]types.Bool{
			"decryption_enabled": state.DecryptionEnabled, "intelligent_proxy": state.IntelligentProxy, "end_user_notification": state.EndUserNotification,
			"single_sign_on": state.SingleSignOn, "safe_search": state.SafeSearch, "ai_semantic_inspection": state.AiSemanticInspection,
		}
		for attr, v := range got {
			if want := attr == "ai_semantic_inspection"; v.IsNull() || v.ValueBool() != want {
				t.Errorf("%s = %s, want %t", attr, v, want)
			}
		}
	})

	errorCases := []struct {
		name, id string
		status   int
		body     string
	}{
		{"an id that does not exist is an error, not an empty data source", "999999999", http.StatusNotFound, `{"statusCode":404,"error":"Not Found","message":"object_not_found"}`},
		{"an id that is not a number is an error", "abc", http.StatusBadRequest, `{"statusCode":400,"error":"Bad Request","message":{"policysettings_id":{"code":"int_invalid","code_text":"not a valid integer"}}}`},
		{"an error of the API is an error", "55555555", http.StatusInternalServerError, `{"message":"internal"}`},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newPolicySettingServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			})
			state, resp := readPolicySetting(t, server, tc.id)
			if !resp.Diagnostics.HasError() {
				t.Fatalf("want an error for the status %d", tc.status)
			}
			d := resp.Diagnostics.Errors()[0]
			if d.Summary() != "Client Error" || !strings.Contains(d.Detail(), "Failed to retrieve object") || !strings.Contains(d.Detail(), "StatusCode "+strconv.Itoa(tc.status)) {
				t.Errorf("diagnostic = %q / %q, want a Client Error that names the failed read and the status %d", d.Summary(), d.Detail(), tc.status)
			}
			if got := server.requests(); len(got) != 1 || !strings.HasPrefix(got[0], "GET /policies/v2/policysettings/"+tc.id) {
				t.Errorf("requests = %q, want one GET of the id and no fallback to another request", got)
			}
			if !state.Id.IsNull() || !state.Name.IsNull() {
				t.Errorf("state must stay empty after an error, got %+v", state)
			}
		})
	}

	t.Run("the id is escaped, it cannot reach another path", func(t *testing.T) {
		server := newPolicySettingServer(t, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"message":"not found"}`, http.StatusNotFound)
		})
		readPolicySetting(t, server, "1/../../securityProfiles")
		got := server.requests()
		if len(got) != 1 || !strings.HasPrefix(got[0], "GET /policies/v2/policysettings/1%2F") || strings.Contains(got[0], "settings/1/") {
			t.Errorf("requests = %q, want one request below /policies/v2/policysettings/ with the slashes of the id escaped", got)
		}
	})
}
