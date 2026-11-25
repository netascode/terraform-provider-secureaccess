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
	"math/big"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/netascode/go-secureaccess"
	"github.com/tidwall/gjson"
)

// This file is not generated. It tests the data source threat_category with items in the shape that the Secure Access API
// returns (the ids and the dates are made up). Values of categoryBits that were seen on a live organization:
//   - "0" for Default Web Settings;
//   - d000000000000000 for Default Settings, whose categories 83, 87 and 89 are the bits 60, 62 and 63 (default categories
//     without a checkbox);
//   - 4000000000000000000000000000 for an item with DNS Tunneling VPN only, whose category 176 is the bit 110;
//   - 5000000000002000000000000000 for an item with exactly Newly Seen Domains, Dynamic DNS and DNS Tunneling VPN ticked, as the
//     Web UI received it (the bits 108, 61 and 110; seen in the answer of the Web UI);
//   - 40000000003000000000002000000000000000 for an item with Potentially Harmful Domains, Newly Seen Domains, Dynamic DNS and
//     Cryptomining ticked, whose categories 174, 172, 85 and 403 are the bits 109, 108, 61 and 150.
//
// Each of the five bits that the data source reads from categoryBits (108, 109, 61, 110 and 150) was seen set together with its
// category in GET by id (a read-only probe of the API). They are the bits of the Web UI source (constants.jsx).
// categoryBits is a JSON string in the list and in GET by id of /policies/v2/threatCategories.

// The bit that the Web UI uses for each option that is read from categoryBits.
const (
	bitNewlySeenDomains       = 108
	bitPotentiallyHarmful     = 109
	bitDynamicDns             = 61
	bitDnsTunnelingVpn        = 110
	bitCryptomining           = 150
	bitDriveByDownloads       = 60 // the Web UI maps these three to default categories without a checkbox
	bitMobileThreats          = 62
	bitHighRiskSites          = 63
	threatCategoryDefaultName = "Default Settings"
)

// threatCategoryHex builds the hexadecimal number with exactly the given bits set, without leading zeros.
func threatCategoryHex(bits ...int) string {
	n := new(big.Int)
	for _, b := range bits {
		n.SetBit(n, b, 1)
	}
	return n.Text(16)
}

// The items of GET /policies/v2/threatCategories. The data source reads the name, the three flags and categoryBits, every
// other key of an item is ignored.
const threatCategoryDefaultItem = `{
  "id": 11111111, "organizationId": 2222222, "isDefault": true, "name": "Default Settings",
  "botnetProtection": true, "malwareProtection": true, "superBotnet": true, "superMalware": true, "phishingProtection": true,
  "suspiciousResponseFiltering": false, "malwareUrlProxy": false, "urlProxyHttps": false, "ipFiltering": false,
  "fileInspection": false, "ampScan": false, "createdAt": "2025-01-09 14:38:57", "modifiedAt": "2025-01-09 14:38:57",
  "categoryBits": "d000000000000000", "markedForDeletion": false, "bundleTypeId": 1, "isSwgDefault": null
}`

const threatCategoryDefaultWebItem = `{
  "id": 11111112, "organizationId": 2222222, "isDefault": false, "name": "Default Web Settings",
  "botnetProtection": true, "malwareProtection": true, "superBotnet": true, "superMalware": true, "phishingProtection": true,
  "suspiciousResponseFiltering": false, "malwareUrlProxy": false, "urlProxyHttps": false, "ipFiltering": false,
  "fileInspection": true, "ampScan": false, "createdAt": "2025-01-09 14:38:57", "modifiedAt": "2025-01-09 14:38:57",
  "categoryBits": "0", "markedForDeletion": false, "bundleTypeId": 2, "isSwgDefault": true
}`

// An item with exactly Newly Seen Domains, Dynamic DNS and DNS Tunneling VPN ticked and no flag, as the Web UI received it.
const threatCategoryCustomItem = `{
  "id": 11111113, "organizationId": 2222222, "isDefault": false, "name": "My Threat Category",
  "botnetProtection": false, "malwareProtection": false, "superBotnet": false, "superMalware": false, "phishingProtection": false,
  "suspiciousResponseFiltering": false, "malwareUrlProxy": false, "urlProxyHttps": false, "ipFiltering": false,
  "fileInspection": false, "ampScan": false, "createdAt": "2025-02-01 10:00:00", "modifiedAt": "2025-02-01 10:00:00",
  "categoryBits": "5000000000002000000000000000", "markedForDeletion": false, "bundleTypeId": 2, "isSwgDefault": null
}`

// The item above in the shape of GET by id, which has the categories in addition. It is made up as a whole: only its categoryBits
// was seen, in the answer of the Web UI.
const threatCategoryCustomItemById = `{
  "id": 11111113, "organizationId": 2222222, "isDefault": false, "name": "My Threat Category",
  "botnetProtection": false, "malwareProtection": false, "superBotnet": false, "superMalware": false, "phishingProtection": false,
  "suspiciousResponseFiltering": false, "malwareUrlProxy": false, "urlProxyHttps": false, "ipFiltering": false,
  "fileInspection": false, "ampScan": false, "createdAt": "2025-02-01 10:00:00", "modifiedAt": "2025-02-01 10:00:00",
  "categoryBits": "5000000000002000000000000000", "markedForDeletion": false, "bundleTypeId": 2, "isSwgDefault": null,
  "categories": [{"categoryId": 172, "name": "Newly Seen Domains"}, {"categoryId": 85, "name": "Dynamic DNS"}, {"categoryId": 176, "name": "DNS Tunneling VPN"}]
}`

// An item with flags and options together (synthetic): Malware, Phishing Attacks, Potentially Harmful Domains and Cryptomining.
var threatCategoryMixedItem = `{
  "id": 11111114, "organizationId": 2222222, "isDefault": false, "name": "Mixed Category",
  "botnetProtection": false, "malwareProtection": true, "superBotnet": false, "superMalware": true, "phishingProtection": true,
  "suspiciousResponseFiltering": false, "malwareUrlProxy": false, "urlProxyHttps": false, "ipFiltering": false,
  "fileInspection": false, "ampScan": false, "createdAt": "2025-02-02 10:00:00", "modifiedAt": "2025-02-02 10:00:00",
  "categoryBits": "` + threatCategoryHex(bitPotentiallyHarmful, bitCryptomining) + `", "markedForDeletion": false, "bundleTypeId": 2, "isSwgDefault": null
}`

// An item with Potentially Harmful Domains, Newly Seen Domains, Dynamic DNS and Cryptomining ticked and no flags. Its categoryBits is the
// one that the API returned for such an item (its categories were 174, 172, 85 and 403); the ids and the dates are made up.
const threatCategoryFourOptionsItem = `{
  "id": 11111115, "organizationId": 2222222, "isDefault": false, "name": "Four Options",
  "botnetProtection": false, "malwareProtection": false, "superBotnet": false, "superMalware": false, "phishingProtection": false,
  "suspiciousResponseFiltering": false, "malwareUrlProxy": false, "urlProxyHttps": false, "ipFiltering": false,
  "fileInspection": false, "ampScan": false, "createdAt": "2025-02-03 10:00:00", "modifiedAt": "2025-02-03 10:00:00",
  "categoryBits": "40000000003000000000002000000000000000", "markedForDeletion": false, "bundleTypeId": 2, "isSwgDefault": null
}`

// The item above in the shape of GET by id: the categories are the ones that were seen with this categoryBits (the name of the
// category 174 is "Potentially Harmful").
var threatCategoryFourOptionsItemById = strings.TrimSuffix(threatCategoryFourOptionsItem, "\n}") + `,
  "categories": [{"categoryId": 85, "name": "Dynamic DNS"}, {"categoryId": 172, "name": "Newly Seen Domains"}, {"categoryId": 174, "name": "Potentially Harmful"}, {"categoryId": 403, "name": "Cryptomining"}]
}`

// threatCategoryBool is one Bool attribute of the data source and how to turn it on alone in an item body.
type threatCategoryBool struct {
	attr string
	get  func(ThreatCategory) types.Bool
	// on returns the keys of an item body that turn this attribute on and everything else off
	on map[string]any
}

// All Bool attributes of the data source. Each one is turned on alone below, so that a value read from the wrong key or the
// wrong bit shows up in at least one case.
var threatCategoryBools = []threatCategoryBool{
	{"malware", func(c ThreatCategory) types.Bool { return c.Malware }, map[string]any{"malwareProtection": true}},
	{"command_and_control_callbacks", func(c ThreatCategory) types.Bool { return c.CommandAndControlCallbacks }, map[string]any{"botnetProtection": true}},
	{"phishing_attacks", func(c ThreatCategory) types.Bool { return c.PhishingAttacks }, map[string]any{"phishingProtection": true}},
	{"potentially_harmful_domains", func(c ThreatCategory) types.Bool { return c.PotentiallyHarmfulDomains }, map[string]any{"categoryBits": threatCategoryHex(bitPotentiallyHarmful)}},
	{"newly_seen_domains", func(c ThreatCategory) types.Bool { return c.NewlySeenDomains }, map[string]any{"categoryBits": threatCategoryHex(bitNewlySeenDomains)}},
	{"dynamic_dns", func(c ThreatCategory) types.Bool { return c.DynamicDns }, map[string]any{"categoryBits": threatCategoryHex(bitDynamicDns)}},
	{"dns_tunneling_vpn", func(c ThreatCategory) types.Bool { return c.DnsTunnelingVpn }, map[string]any{"categoryBits": threatCategoryHex(bitDnsTunnelingVpn)}},
	{"cryptomining", func(c ThreatCategory) types.Bool { return c.Cryptomining }, map[string]any{"categoryBits": threatCategoryHex(bitCryptomining)}},
}

// threatCategoryBody returns an item body with the three flags off and categoryBits "0", and the given keys replaced.
func threatCategoryBody(t *testing.T, over map[string]any) string {
	t.Helper()
	body := map[string]any{
		"name": "c", "categoryBits": "0", "malwareProtection": false, "botnetProtection": false, "phishingProtection": false,
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

// threatCategoryBoolsOf returns the value of every Bool attribute of a data source.
func threatCategoryBoolsOf(c ThreatCategory) map[string]types.Bool {
	got := map[string]types.Bool{}
	for _, b := range threatCategoryBools {
		got[b.attr] = b.get(c)
	}
	return got
}

func TestThreatCategoryFromBody(t *testing.T) {
	ctx := context.Background()
	flags := []string{"malware", "command_and_control_callbacks", "phishing_attacks"}
	options := []string{"potentially_harmful_domains", "newly_seen_domains", "dynamic_dns", "dns_tunneling_vpn", "cryptomining"}

	t.Run("each attribute is read from its own key or bit", func(t *testing.T) {
		for _, on := range threatCategoryBools {
			var got ThreatCategory
			got.fromBody(ctx, gjson.Parse(threatCategoryBody(t, on.on)))
			for attr, v := range threatCategoryBoolsOf(got) {
				if want := attr == on.attr; v.IsNull() || v.ValueBool() != want {
					t.Errorf("with only %s on, %s = %s, want %t", on.attr, attr, v, want)
				}
			}
		}
	})

	// the number needs a letter, or the upper case spelling would be the same as the lower case one: the bits 109, 110 and 111 are
	// the digit e (and 111 is not mapped to an option)
	upperCaseHex := strings.ToUpper(threatCategoryHex(bitPotentiallyHarmful, bitDnsTunnelingVpn, 111, bitCryptomining))
	if upperCaseHex == strings.ToLower(upperCaseHex) {
		t.Fatalf("test fixture: %s has no letter, the case would not be tested", upperCaseHex)
	}

	bitCases := []struct {
		name string
		bits any // the JSON value of categoryBits
		on   []string
	}{
		{"the value of an item with three options", "5000000000002000000000000000", []string{"newly_seen_domains", "dynamic_dns", "dns_tunneling_vpn"}},
		{"the example of one option", "4000000000000000000000000000", []string{"dns_tunneling_vpn"}},
		{"nothing set", "0", nil},
		{"an empty string is nothing set", "", nil},
		{"leading zeros", "00000000000000000000000000004000000000000000000000000000", []string{"dns_tunneling_vpn"}},
		{"upper case digits", upperCaseHex, []string{"potentially_harmful_domains", "dns_tunneling_vpn", "cryptomining"}},
		{"all five options", threatCategoryHex(bitNewlySeenDomains, bitPotentiallyHarmful, bitDynamicDns, bitDnsTunnelingVpn, bitCryptomining), options},
		{"the default categories have no option", threatCategoryHex(bitDriveByDownloads, bitMobileThreats, bitHighRiskSites), nil},
		{"other bits are ignored next to the options", threatCategoryHex(bitDriveByDownloads, bitNewlySeenDomains, bitHighRiskSites, 200), []string{"newly_seen_domains"}},
		{"a bit that is not mapped", threatCategoryHex(154), nil},
		{"bits far above the options are ignored", threatCategoryHex(4096, bitNewlySeenDomains, bitCryptomining), []string{"newly_seen_domains", "cryptomining"}},
		{"the bits next to a mapped bit", threatCategoryHex(107, 111, 60, 62, 149, 151), nil},
	}
	for _, tt := range bitCases {
		t.Run(tt.name, func(t *testing.T) {
			var got ThreatCategory
			got.fromBody(ctx, gjson.Parse(threatCategoryBody(t, map[string]any{"categoryBits": tt.bits})))
			on := map[string]bool{}
			for _, a := range tt.on {
				on[a] = true
			}
			for _, a := range options {
				v := threatCategoryBoolsOf(got)[a]
				if v.IsNull() || v.ValueBool() != on[a] {
					t.Errorf("%s = %s, want %t", a, v, on[a])
				}
			}
			// the bits do not change a flag
			for _, a := range flags {
				if v := threatCategoryBoolsOf(got)[a]; v.IsNull() || v.ValueBool() {
					t.Errorf("%s = %s, want false", a, v)
				}
			}
		})
	}

	t.Run("flags that are on do not turn an option on", func(t *testing.T) {
		var got ThreatCategory
		got.fromBody(ctx, gjson.Parse(threatCategoryBody(t, map[string]any{"malwareProtection": true, "botnetProtection": true, "phishingProtection": true})))
		for _, a := range options {
			if v := threatCategoryBoolsOf(got)[a]; v.IsNull() || v.ValueBool() {
				t.Errorf("%s = %s, want false", a, v)
			}
		}
	})

	// A categoryBits that cannot be read as a hexadecimal number leaves the options null, they are neither false nor true.
	// The flags of the item are still read.
	for name, bits := range map[string]any{"null": nil, "a number": 5, "not hexadecimal": "zz", "a sign": "-1", "a 0x prefix": "0x4000", "an object": map[string]any{}} {
		t.Run("categoryBits is "+name, func(t *testing.T) {
			var got ThreatCategory
			got.fromBody(ctx, gjson.Parse(threatCategoryBody(t, map[string]any{"categoryBits": bits, "malwareProtection": true})))
			for _, a := range options {
				if v := threatCategoryBoolsOf(got)[a]; !v.IsNull() {
					t.Errorf("%s = %s, want null", a, v)
				}
			}
			if !got.Malware.Equal(types.BoolValue(true)) || !got.PhishingAttacks.Equal(types.BoolValue(false)) {
				t.Errorf("the flags are not read: malware = %s, phishing_attacks = %s", got.Malware, got.PhishingAttacks)
			}
		})
	}

	t.Run("keys that are missing are null, not false", func(t *testing.T) {
		var got ThreatCategory
		got.fromBody(ctx, gjson.Parse(`{}`))
		if !got.Id.IsNull() || !got.Name.IsNull() {
			t.Errorf("id = %s, name = %s, want null", got.Id, got.Name)
		}
		for attr, v := range threatCategoryBoolsOf(got) {
			if !v.IsNull() {
				t.Errorf("%s = %s, want null", attr, v)
			}
		}
	})

	t.Run("the name is read, the id is not: it comes from the config or the lookup", func(t *testing.T) {
		var got ThreatCategory
		got.fromBody(ctx, gjson.Parse(threatCategoryBody(t, map[string]any{"name": "Some Name", "id": 11111113})))
		if got.Name.ValueString() != "Some Name" {
			t.Errorf("name = %s", got.Name)
		}
		if !got.Id.IsNull() {
			t.Errorf("id = %s, fromBody must not set it", got.Id)
		}
	})

	t.Run("the items as the API returns them", func(t *testing.T) {
		// the number of the item with four options is written out above, as it came from the API: it must be the four bits
		if want := threatCategoryHex(bitPotentiallyHarmful, bitNewlySeenDomains, bitDynamicDns, bitCryptomining); !strings.Contains(threatCategoryFourOptionsItem, `"categoryBits": "`+want+`"`) {
			t.Fatalf("test fixture: the categoryBits of the item with four options is not %s", want)
		}
		for _, tt := range []struct {
			item string
			on   []string
		}{
			{threatCategoryDefaultItem, flags},
			{threatCategoryDefaultWebItem, flags},
			{threatCategoryCustomItem, []string{"newly_seen_domains", "dynamic_dns", "dns_tunneling_vpn"}},
			{threatCategoryMixedItem, []string{"malware", "phishing_attacks", "potentially_harmful_domains", "cryptomining"}},
			{threatCategoryFourOptionsItem, []string{"potentially_harmful_domains", "newly_seen_domains", "dynamic_dns", "cryptomining"}},
		} {
			var got ThreatCategory
			got.fromBody(ctx, gjson.Parse(tt.item))
			on := map[string]bool{}
			for _, a := range tt.on {
				on[a] = true
			}
			for attr, v := range threatCategoryBoolsOf(got) {
				if v.IsNull() || v.ValueBool() != on[attr] {
					t.Errorf("%s: %s = %s, want %t", got.Name, attr, v, on[attr])
				}
			}
		}
	})
}

// The data source reads the attributes that the definition lists and nothing else: a new attribute would change the docs
// and the interface of the provider silently.
func TestThreatCategorySchema(t *testing.T) {
	ctx := context.Background()
	ds := &ThreatCategoryDataSource{}

	var meta datasource.MetadataResponse
	ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "secureaccess"}, &meta)
	if meta.TypeName != "secureaccess_threat_category" {
		t.Errorf("type name = %q", meta.TypeName)
	}

	var resp datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &resp)
	var names []string
	for name, a := range resp.Schema.Attributes {
		names = append(names, name)
		switch name {
		case "id", "name":
			if !a.IsOptional() || !a.IsComputed() || a.IsRequired() {
				t.Errorf("%s is a query attribute and must be optional and computed, required=%t optional=%t computed=%t", name, a.IsRequired(), a.IsOptional(), a.IsComputed())
			}
		default:
			if !a.IsComputed() || a.IsRequired() || a.IsOptional() {
				t.Errorf("%s must be read-only, required=%t optional=%t computed=%t", name, a.IsRequired(), a.IsOptional(), a.IsComputed())
			}
		}
	}
	sort.Strings(names)
	want := "command_and_control_callbacks,cryptomining,dns_tunneling_vpn,dynamic_dns,id,malware,name,newly_seen_domains,phishing_attacks,potentially_harmful_domains"
	if got := strings.Join(names, ","); got != want {
		t.Errorf("attributes = %s\nwant       %s", got, want)
	}

	// a Threat Category is found by its id or by its name, one of them and not both
	validators := ds.ConfigValidators(ctx)
	if len(validators) != 1 {
		t.Fatalf("config validators = %d, want 1 (id or name)", len(validators))
	}
	objType := resp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	for _, c := range []struct {
		name      string
		id, name2 any
		wantErr   bool
	}{
		{"id only", "1", nil, false},
		{"name only", nil, "n", false},
		{"id and name", "1", "n", true},
		{"neither", nil, nil, true},
	} {
		values := map[string]tftypes.Value{}
		for name, typ := range objType.AttributeTypes {
			values[name] = tftypes.NewValue(typ, nil)
		}
		values["id"] = tftypes.NewValue(tftypes.String, c.id)
		values["name"] = tftypes.NewValue(tftypes.String, c.name2)
		var vresp datasource.ValidateConfigResponse
		validators[0].ValidateDataSource(ctx, datasource.ValidateConfigRequest{Config: tfsdk.Config{Schema: resp.Schema, Raw: tftypes.NewValue(objType, values)}}, &vresp)
		if vresp.Diagnostics.HasError() != c.wantErr {
			t.Errorf("%s: error = %t, want %t (%v)", c.name, vresp.Diagnostics.HasError(), c.wantErr, vresp.Diagnostics)
		}
	}

	var registered bool
	for _, factory := range New("test")().DataSources(ctx) {
		var m datasource.MetadataResponse
		factory().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "secureaccess"}, &m)
		if m.TypeName == "secureaccess_threat_category" {
			registered = true
		}
	}
	if !registered {
		t.Error("the data source is not registered in the provider")
	}
}

// threatCategoryServer is a fake Secure Access API: it answers the token request and everything else with the handler,
// and records the method and the escaped path of every other request.
type threatCategoryServer struct {
	*httptest.Server
	mu    sync.Mutex
	paths []string
}

func newThreatCategoryServer(t *testing.T, handle func(w http.ResponseWriter, r *http.Request)) *threatCategoryServer {
	t.Helper()
	s := &threatCategoryServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/v2/token" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"token","token_type":"bearer","expires_in":3600}`))
			return
		}
		s.mu.Lock()
		s.paths = append(s.paths, r.Method+" "+r.URL.EscapedPath())
		s.mu.Unlock()
		handle(w, r)
	}))
	t.Cleanup(s.Close)
	return s
}

// requests returns the method and the escaped path of every request, without the query: the paginator of the client
// decides which query a list request has, it is not what is tested here.
func (s *threatCategoryServer) requests() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.paths...)
}

// threatCategoryAPI answers with the shapes that were seen for this route: the list is a bare array that ignores page and limit,
// GET by id returns one object that has the categories in addition. The 404 for an unknown id has the body that the API answered
// for it (the txId is made up).
func threatCategoryAPI() func(w http.ResponseWriter, r *http.Request) {
	byId := map[string]string{
		"/policies/v2/threatCategories/11111111": threatCategoryDefaultItem,
		"/policies/v2/threatCategories/11111112": threatCategoryDefaultWebItem,
		"/policies/v2/threatCategories/11111113": threatCategoryCustomItemById,
		"/policies/v2/threatCategories/11111114": threatCategoryMixedItem,
		"/policies/v2/threatCategories/11111115": threatCategoryFourOptionsItemById,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/policies/v2/threatCategories":
			_, _ = w.Write([]byte("[" + threatCategoryDefaultItem + "," + threatCategoryDefaultWebItem + "," + threatCategoryCustomItem + "," + threatCategoryMixedItem + "," + threatCategoryFourOptionsItem + "]"))
		case r.Method == http.MethodGet && byId[r.URL.EscapedPath()] != "":
			_, _ = w.Write([]byte(byId[r.URL.EscapedPath()]))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"statusCode":404,"error":"Not Found","message":"object_not_found","txId":"00000000-0000-0000-0000-000000000000"}`))
		}
	}
}

// readThreatCategory runs the Read() of the data source with the given id or name in the config (an empty one is not set).
func readThreatCategory(t *testing.T, server *threatCategoryServer, id, name string) (ThreatCategory, datasource.ReadResponse) {
	t.Helper()
	ctx := context.Background()
	client, err := secureaccess.NewClient(server.URL, "key", "secret", secureaccess.MaxRetries(0))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	ds := &ThreatCategoryDataSource{client: client}

	var schemaResp datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	objType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	config := map[string]tftypes.Value{}
	for attrName, typ := range objType.AttributeTypes {
		config[attrName] = tftypes.NewValue(typ, nil)
	}
	if id != "" {
		config["id"] = tftypes.NewValue(tftypes.String, id)
	}
	if name != "" {
		config["name"] = tftypes.NewValue(tftypes.String, name)
	}

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, config)}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, nil)}}
	ds.Read(ctx, req, &resp)

	var state ThreatCategory
	if !resp.Diagnostics.HasError() {
		if diags := resp.State.Get(ctx, &state); diags.HasError() {
			t.Fatalf("reading the state: %v", diags)
		}
	}
	return state, resp
}

// checkThreatCategoryState compares every Bool attribute of the state with the attributes that are on.
func checkThreatCategoryState(t *testing.T, state ThreatCategory, on ...string) {
	t.Helper()
	want := map[string]bool{}
	for _, a := range on {
		want[a] = true
	}
	for attr, v := range threatCategoryBoolsOf(state) {
		if v.IsNull() || v.ValueBool() != want[attr] {
			t.Errorf("%s = %s, want %t", attr, v, want[attr])
		}
	}
}

// Two items of the list with the same name: the first one of the list is found. That is how the generated Read() of every data
// source with a lookup by name behaves. The test makes a change of it visible, it does not say that duplicate names are expected.
func TestThreatCategoryDataSourceDuplicateNames(t *testing.T) {
	duplicate := strings.Replace(threatCategoryMixedItem, `"id": 11111114`, `"id": 11111199`, 1)
	duplicate = strings.Replace(duplicate, `"name": "Mixed Category"`, `"name": "My Threat Category"`, 1)
	if !strings.Contains(duplicate, `"id": 11111199`) || !strings.Contains(duplicate, `"name": "My Threat Category"`) {
		t.Fatal("test helper: the duplicate item was not built")
	}
	for _, tt := range []struct {
		name  string
		list  string
		id    string
		onAre []string
	}{
		{"the first item of the list", "[" + threatCategoryCustomItem + "," + duplicate + "]", "11111113", []string{"newly_seen_domains", "dynamic_dns", "dns_tunneling_vpn"}},
		{"the other item when it is first", "[" + duplicate + "," + threatCategoryCustomItem + "]", "11111199", []string{"malware", "phishing_attacks", "potentially_harmful_domains", "cryptomining"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := newThreatCategoryServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.list))
			})
			state, resp := readThreatCategory(t, server, "", "My Threat Category")
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			if state.Id.ValueString() != tt.id {
				t.Errorf("id = %q, want %q", state.Id.ValueString(), tt.id)
			}
			checkThreatCategoryState(t, state, tt.onAre...)
		})
	}
}

// Configure takes the client of the provider. Without provider data (the provider is not configured yet) it does nothing.
func TestThreatCategoryDataSourceConfigure(t *testing.T) {
	ctx := context.Background()

	t.Run("without provider data the client stays empty and there is no error", func(t *testing.T) {
		ds := &ThreatCategoryDataSource{}
		var resp datasource.ConfigureResponse
		ds.Configure(ctx, datasource.ConfigureRequest{ProviderData: nil}, &resp)
		if resp.Diagnostics.HasError() || ds.client != nil {
			t.Errorf("client = %v, diagnostics = %v, want no client and no error", ds.client, resp.Diagnostics)
		}
	})

	t.Run("the client of the provider is used", func(t *testing.T) {
		server := newThreatCategoryServer(t, threatCategoryAPI())
		client, err := secureaccess.NewClient(server.URL, "key", "secret", secureaccess.MaxRetries(0))
		if err != nil {
			t.Fatalf("NewClient: %v", err)
		}
		ds := &ThreatCategoryDataSource{}
		var resp datasource.ConfigureResponse
		ds.Configure(ctx, datasource.ConfigureRequest{ProviderData: &SecureAccessProviderData{Client: client}}, &resp)
		if resp.Diagnostics.HasError() || ds.client != client {
			t.Errorf("client = %p, want %p, diagnostics = %v", ds.client, client, resp.Diagnostics)
		}
	})
}

func TestThreatCategoryDataSourceRead(t *testing.T) {
	customOn := []string{"newly_seen_domains", "dynamic_dns", "dns_tunneling_vpn"}

	t.Run("by id: one request, the options are read from categoryBits", func(t *testing.T) {
		server := newThreatCategoryServer(t, threatCategoryAPI())
		state, resp := readThreatCategory(t, server, "11111113", "")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if got := server.requests(); len(got) != 1 || got[0] != "GET /policies/v2/threatCategories/11111113" {
			t.Errorf("requests = %q, want exactly one GET by id", got)
		}
		if state.Id.ValueString() != "11111113" || state.Name.ValueString() != "My Threat Category" {
			t.Errorf("id = %q, name = %q", state.Id.ValueString(), state.Name.ValueString())
		}
		checkThreatCategoryState(t, state, customOn...)
	})

	t.Run("by name: one request for the list, no second request by id, the same values as by id", func(t *testing.T) {
		server := newThreatCategoryServer(t, threatCategoryAPI())
		byName, resp := readThreatCategory(t, server, "", "My Threat Category")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if got := server.requests(); len(got) != 1 || got[0] != "GET /policies/v2/threatCategories" {
			t.Errorf("requests = %q, want exactly one GET of the list", got)
		}
		if byName.Id.ValueString() != "11111113" || byName.Name.ValueString() != "My Threat Category" {
			t.Errorf("id = %q, name = %q, want the item that has the name", byName.Id.ValueString(), byName.Name.ValueString())
		}
		checkThreatCategoryState(t, byName, customOn...)

		byId, resp := readThreatCategory(t, newThreatCategoryServer(t, threatCategoryAPI()), "11111113", "")
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if !byName.Id.Equal(byId.Id) || !byName.Name.Equal(byId.Name) {
			t.Errorf("by name: %+v, by id: %+v", byName, byId)
		}
		for attr, v := range threatCategoryBoolsOf(byName) {
			if !v.Equal(threatCategoryBoolsOf(byId)[attr]) {
				t.Errorf("%s: by name %s, by id %s", attr, v, threatCategoryBoolsOf(byId)[attr])
			}
		}
	})

	for _, tt := range []struct {
		name, id string
		on       []string
	}{
		{threatCategoryDefaultName, "11111111", []string{"malware", "command_and_control_callbacks", "phishing_attacks"}},
		{"Default Web Settings", "11111112", []string{"malware", "command_and_control_callbacks", "phishing_attacks"}},
		{"Mixed Category", "11111114", []string{"malware", "phishing_attacks", "potentially_harmful_domains", "cryptomining"}},
		{"Four Options", "11111115", []string{"potentially_harmful_domains", "newly_seen_domains", "dynamic_dns", "cryptomining"}},
	} {
		t.Run("the item "+tt.name+" is found by its name", func(t *testing.T) {
			server := newThreatCategoryServer(t, threatCategoryAPI())
			state, resp := readThreatCategory(t, server, "", tt.name)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			if state.Id.ValueString() != tt.id {
				t.Errorf("id = %q, want %q", state.Id.ValueString(), tt.id)
			}
			checkThreatCategoryState(t, state, tt.on...)
		})
		// the same item by its id: the same values from the shape of GET by id
		t.Run("the item "+tt.name+" is read by its id", func(t *testing.T) {
			server := newThreatCategoryServer(t, threatCategoryAPI())
			state, resp := readThreatCategory(t, server, tt.id, "")
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			if state.Id.ValueString() != tt.id || state.Name.ValueString() != tt.name {
				t.Errorf("id = %q, name = %q, want %q and %q", state.Id.ValueString(), state.Name.ValueString(), tt.id, tt.name)
			}
			checkThreatCategoryState(t, state, tt.on...)
			if got := server.requests(); len(got) != 1 || got[0] != "GET /policies/v2/threatCategories/"+tt.id {
				t.Errorf("requests = %q, want exactly one GET by id", got)
			}
		})
	}

	t.Run("a name that does not exist is an error", func(t *testing.T) {
		server := newThreatCategoryServer(t, threatCategoryAPI())
		_, resp := readThreatCategory(t, server, "", "No Such Category")
		if n := len(resp.Diagnostics.Errors()); n != 1 {
			t.Fatalf("errors = %d, want exactly 1: %v", n, resp.Diagnostics)
		}
		d := resp.Diagnostics.Errors()[0]
		if d.Summary() != "Client Error" || !strings.Contains(d.Detail(), "Failed to find object with name: No Such Category") {
			t.Errorf("diagnostic = %q / %q", d.Summary(), d.Detail())
		}
		if got := server.requests(); len(got) != 1 {
			t.Errorf("requests = %q, want one request for the list and nothing after it", got)
		}
		if !resp.State.Raw.IsNull() {
			t.Errorf("the state must stay empty after an error, got %s", resp.State.Raw)
		}
	})

	t.Run("the names are compared exactly", func(t *testing.T) {
		server := newThreatCategoryServer(t, threatCategoryAPI())
		for _, name := range []string{"my threat category", "My Threat Category ", "My Threat"} {
			if _, resp := readThreatCategory(t, server, "", name); !resp.Diagnostics.HasError() {
				t.Errorf("the name %q must not match", name)
			}
		}
	})

	// The client rejects this answer itself (its paginator expects a bare array), so the data source gets an error from the client:
	// the check of the data source for a list that is not an array is not reached with this client.
	t.Run("an answer that wraps the list in an object is an error", func(t *testing.T) {
		server := newThreatCategoryServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":{"code":200},"meta":{"total":1},"data":[` + threatCategoryCustomItem + `]}`))
		})
		_, resp := readThreatCategory(t, server, "", "My Threat Category")
		if n := len(resp.Diagnostics.Errors()); n != 1 {
			t.Fatalf("errors = %d, want exactly 1: %v", n, resp.Diagnostics)
		}
		if d := resp.Diagnostics.Errors()[0]; !strings.Contains(d.Detail(), "Failed to retrieve objects") {
			t.Errorf("diagnostic = %q / %q", d.Summary(), d.Detail())
		}
		if !resp.State.Raw.IsNull() {
			t.Errorf("the state must stay empty after an error, got %s", resp.State.Raw)
		}
	})

	// The bodies of the first two cases are the ones that the API answered for these ids (a read-only probe, the txId is made up).
	// The body of the error of the list was not seen, it is made up.
	errorCases := []struct {
		name, id, itemName string
		status             int
		body, wantDetail   string
	}{
		{"an id that does not exist is an error, not an empty data source", "999999999", "", http.StatusNotFound, `{"statusCode":404,"error":"Not Found","message":"object_not_found","txId":"00000000-0000-0000-0000-000000000000"}`, "Failed to retrieve object,"},
		{"an id that is not a number is an error", "abc", "", http.StatusBadRequest, `{"statusCode":400,"error":"Bad Request","message":{"securitysettings_id":{"code":"int_invalid","code_text":"not a valid integer"}},"txId":"00000000-0000-0000-0000-000000000000"}`, "Failed to retrieve object,"},
		{"an error of the list is an error", "", "My Threat Category", http.StatusBadRequest, `{"message":"bad request"}`, "Failed to retrieve objects,"},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newThreatCategoryServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			})
			_, resp := readThreatCategory(t, server, tc.id, tc.itemName)
			if n := len(resp.Diagnostics.Errors()); n != 1 {
				t.Fatalf("errors = %d, want exactly 1: %v", n, resp.Diagnostics)
			}
			d := resp.Diagnostics.Errors()[0]
			if d.Summary() != "Client Error" || !strings.Contains(d.Detail(), tc.wantDetail) || !strings.Contains(d.Detail(), "StatusCode "+strconv.Itoa(tc.status)) {
				t.Errorf("diagnostic = %q / %q, want a Client Error that names %q and the status %d", d.Summary(), d.Detail(), tc.wantDetail, tc.status)
			}
			if got := server.requests(); len(got) != 1 {
				t.Errorf("requests = %q, want one request and no fallback to another", got)
			}
			if !resp.State.Raw.IsNull() {
				t.Errorf("the state must stay empty after an error, got %s", resp.State.Raw)
			}
		})
	}

	t.Run("the id is escaped, it cannot reach another path", func(t *testing.T) {
		server := newThreatCategoryServer(t, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"message":"not found"}`, http.StatusNotFound)
		})
		readThreatCategory(t, server, "1/../../securityProfiles", "")
		got := server.requests()
		if len(got) != 1 || !strings.HasPrefix(got[0], "GET /policies/v2/threatCategories/1%2F") || strings.Contains(got[0], "threatCategories/1/") {
			t.Errorf("requests = %q, want one request below /policies/v2/threatCategories/ with the slashes of the id escaped", got)
		}
	})
}

// The state is checked by attribute name, not through the Go struct: the Read tests above get the struct back from the same state,
// so two swapped `tfsdk` tags would pass them. Each option is turned on alone in a list item and every attribute of the state is
// compared by its name.
func TestThreatCategoryStateAttributeNames(t *testing.T) {
	for _, on := range threatCategoryBools {
		t.Run(on.attr, func(t *testing.T) {
			over := map[string]any{"id": 11111120, "name": "Only One"}
			for k, v := range on.on {
				over[k] = v
			}
			item := threatCategoryBody(t, over)
			server := newThreatCategoryServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte("[" + item + "]"))
			})
			_, resp := readThreatCategory(t, server, "", "Only One")
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			var raw map[string]tftypes.Value
			if err := resp.State.Raw.As(&raw); err != nil {
				t.Fatalf("state: %v", err)
			}
			for _, other := range threatCategoryBools {
				var got *bool
				if err := raw[other.attr].As(&got); err != nil {
					t.Fatalf("%s: %v", other.attr, err)
				}
				if got == nil {
					t.Errorf("with only %s on, %s is null in the state", on.attr, other.attr)
				} else if want := other.attr == on.attr; *got != want {
					t.Errorf("with only %s on, %s = %t in the state, want %t", on.attr, other.attr, *got, want)
				}
			}
		})
	}
}

// Each description names the option of its own attribute. The descriptions come from the YAML and regenerate cleanly if two of them
// are swapped, so only a test notices.
func TestThreatCategoryDescriptionsNameTheirOption(t *testing.T) {
	ctx := context.Background()
	var resp datasource.SchemaResponse
	(&ThreatCategoryDataSource{}).Schema(ctx, datasource.SchemaRequest{}, &resp)

	for attr, option := range map[string]string{
		"malware":                       "Malware",
		"command_and_control_callbacks": "Command and Control Callbacks",
		"phishing_attacks":              "Phishing Attacks",
		"potentially_harmful_domains":   "Potentially Harmful Domains",
		"newly_seen_domains":            "Newly Seen Domains",
		"dynamic_dns":                   "Dynamic DNS",
		"dns_tunneling_vpn":             "DNS Tunneling VPN",
		"cryptomining":                  "Cryptomining",
	} {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Errorf("attribute %s is missing", attr)
			continue
		}
		description := a.GetMarkdownDescription()
		if !strings.HasPrefix(description, "Specifies whether "+option+" ") || !strings.HasSuffix(description, "(Web UI Threat Categories, "+option+").") ||
			!(strings.Contains(description, " is enabled ") || strings.Contains(description, " are enabled ")) {
			t.Errorf("%s: description %q does not say \"Specifies whether %s is|are enabled (Web UI Threat Categories, %s).\"", attr, description, option, option)
		}
	}
}
