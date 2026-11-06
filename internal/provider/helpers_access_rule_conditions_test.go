// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
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
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netascode/go-secureaccess"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	tN  = "umbrella.destination.networkObjectIds"
	tNG = "umbrella.destination.networkObjectGroupIds"
	tS  = "umbrella.destination.serviceObjectIds"
	tSG = "umbrella.destination.serviceObjectGroupIds"
	tC  = "umbrella.destination.category_ids"
	tCL = "umbrella.destination.category_list_ids"
	tA  = "umbrella.destination.application_ids"
	tAC = "umbrella.destination.application_category_ids"
	tAL = "umbrella.destination.application_list_ids"
	tPR = "umbrella.destination.private_resource_ids"
	tID = "umbrella.source.identity_ids"
	tIT = "umbrella.source.identity_type_ids"
	tGE = "umbrella.source.geolocations_shared"
	tSN = "umbrella.source.networkObjectIds"
	tSC = "umbrella.source.composite_inline_ip"
)

// testAccessRuleBody returns a request body with the given condition names
func testAccessRuleBody(names ...string) string {
	body := `{"ruleName":"test"}`
	for _, n := range names {
		body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": n, "attributeOperator": "IN", "attributeValue": []int64{1}})
	}
	return body
}

func testAccessRuleCondition(body, name string) gjson.Result {
	return gjson.Get(body, `ruleConditions.#(attributeName=="`+name+`")`)
}

// Expressions accepted by the API (see gen/definitions/access_rule.yaml)
func TestAdjustAccessRuleConditionsDestination(t *testing.T) {
	cases := []struct {
		name       string
		conditions []string
		merge      string
		expected   string
	}{
		{"network and service", []string{tN, tS}, "AND", tN + " && " + tS},
		{"network and service with groups", []string{tN, tNG, tS, tSG}, "AND", "(" + tN + " || " + tNG + ") && (" + tS + " || " + tSG + ")"},
		{"category and objects merged", []string{tN, tS, tC}, "OR", tC + " && (" + tN + " || " + tS + ")"},
		{"category, objects merged and application", []string{tC, tN, tS, tA}, "OR", tA + " && " + tC + " && (" + tN + " || " + tS + ")"},
		{"category, network and service", []string{tC, tN, tS}, "AND", tC + " && " + tN + " && " + tS},
		{"applications and network", []string{tA, tAC, tN}, "AND", "(" + tA + " || " + tAC + ") && " + tN},
		{"category and application list", []string{tC, tAL}, "AND", tAL + " && " + tC},
		{"category lists and service groups", []string{tSG, tCL}, "AND", tCL + " && " + tSG},
		{"merge without category", []string{tN, tS}, "OR", "(" + tN + " || " + tS + ")"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body := adjustAccessRuleConditions(testAccessRuleBody(append([]string{tID}, c.conditions...)...), "OR", "AND", c.merge)
			got := testAccessRuleCondition(body, "umbrella.destination.logical_operator")
			if got.Get("attributeOperator").String() != "AND" || got.Get("attributeValue").String() != c.expected {
				t.Errorf("got %s, expected %q", got.Raw, c.expected)
			}
			// The built expression is read back with the configured operator
			if op := accessRuleNetworkServiceOperator(c.expected); op != c.merge && len(c.conditions) > 2 {
				t.Errorf("network/service operator of %q: got %s, expected %s", c.expected, op, c.merge)
			}
		})
	}
}

func TestAdjustAccessRuleConditionsSource(t *testing.T) {
	body := adjustAccessRuleConditions(testAccessRuleBody(tID, tIT, tGE), "AND", "OR", "AND")
	if got := testAccessRuleCondition(body, "umbrella.source.logical_operator").Get("attributeValue").String(); got != tID+" && "+tIT+" && "+tGE {
		t.Errorf("unexpected source expression %q", got)
	}
	if testAccessRuleCondition(body, "umbrella.destination.logical_operator").Exists() {
		t.Errorf("unexpected destination logical operator")
	}
	if !testAccessRuleCondition(body, "umbrella.destination.all").Get("attributeValue").Bool() {
		t.Errorf("expected destination all, got %s", body)
	}
}

func TestAdjustAccessRuleConditionsSourceInternetLiterals(t *testing.T) {
	body := adjustAccessRuleConditions(testAccessRuleBody(tSC, tIT), "AND", "OR", "AND")
	if got := testAccessRuleCondition(body, "umbrella.source.logical_operator").Get("attributeValue").String(); got != tIT+" && "+tSC {
		t.Errorf("unexpected source expression %q", got)
	}
}

func TestAdjustAccessRuleConditionsAll(t *testing.T) {
	body := adjustAccessRuleConditions(testAccessRuleBody(), "OR", "OR", "AND")
	for _, side := range []string{accessRuleSource, accessRuleDestination} {
		c := testAccessRuleCondition(body, accessRuleAllName(side))
		if c.Get("attributeOperator").String() != "=" || !c.Get("attributeValue").Bool() {
			t.Errorf("expected %s all, got %s", side, body)
		}
	}
	// OR does not add a logical operator nor `all` on a side with conditions
	body = adjustAccessRuleConditions(testAccessRuleBody(tSN, tN, tS), "OR", "OR", "AND")
	if n := len(gjson.Get(body, "ruleConditions").Array()); n != 3 {
		t.Errorf("expected 3 conditions, got %s", body)
	}
	// Unsupported conditions are kept as a separate term, rejected by the API and ValidateConfig
	body = adjustAccessRuleConditions(testAccessRuleBody(tPR, tN), "OR", "AND", "AND")
	if got := testAccessRuleCondition(body, "umbrella.destination.logical_operator").Get("attributeValue").String(); got != tN+" && "+tPR {
		t.Errorf("unexpected expression %q", got)
	}
}

func TestAccessRuleNetworkServiceOperator(t *testing.T) {
	cases := map[string]string{
		// Written by the portal
		"(" + tN + "&&" + tS + ")": "AND",
		// From the API documentation
		"((" + tN + " || " + tNG + ") && (" + tS + " || " + tSG + "))":       "AND",
		tC + " && (" + tS + " || " + tN + ")":                                "OR",
		tAL + " && (" + tS + " || " + tN + " || " + tSG + " || " + tNG + ")": "OR",
		"(" + tAC + " || " + tA + ") && (" + tS + " ||  " + tSG + ")":        "AND",
		"(" + tAC + " || " + tA + ") && " + tN:                               "AND",
		tC + " && " + tN + " && " + tS:                                       "AND",
		tC + " && ((" + tN + " || " + tNG + ") || " + tS + ")":               "OR",
	}
	for expression, expected := range cases {
		if got := accessRuleNetworkServiceOperator(expression); got != expected {
			t.Errorf("%q: got %s, expected %s", expression, got, expected)
		}
	}
}

func TestSplitAccessRuleExpression(t *testing.T) {
	cases := map[string][]string{
		"(a&&b)":             {"a", "b"},
		"(a) && (b)":         {"a", "b"},
		"((a || b) && c)":    {"a || b", "c"},
		"a && (b || (c))":    {"a", "b || (c)"},
		"  ((a || b))  ":     {"a || b"},
		"(a || b) && (c||d)": {"a || b", "c||d"},
	}
	for expression, expected := range cases {
		if got := splitAccessRuleExpression(expression); !reflect.DeepEqual(got, expected) {
			t.Errorf("%q: got %q, expected %q", expression, got, expected)
		}
	}
}

func TestRetryAccessRuleInvalidIds(t *testing.T) {
	defer func(d time.Duration) { accessRuleInvalidIdsDelay = d }(accessRuleInvalidIdsDelay)
	accessRuleInvalidIdsDelay = time.Millisecond
	response := func(status int, msg string) secureaccess.Res {
		return secureaccess.Res{Result: gjson.Parse(`{"msg":"` + msg + `"}`), StatusCode: status}
	}
	invalidIds := response(400, "invalid data passed. the ID's provided for umbrella.source.identity_ids are not valid for org 1. received [1]")
	cases := []struct {
		name     string
		failures []secureaccess.Res
		calls    int
		fails    bool
	}{
		{"success", nil, 1, false},
		{"invalid ids then success", []secureaccess.Res{invalidIds, invalidIds}, 3, false},
		{"invalid ids exhausted", []secureaccess.Res{invalidIds, invalidIds, invalidIds, invalidIds, invalidIds}, 4, true},
		{"other 400 not retried", []secureaccess.Res{response(400, "invalid data passed. Protocol must be one of TCP, UDP, ICMP, ANY")}, 1, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			calls := 0
			_, err := retryAccessRuleInvalidIds(context.Background(), func() (secureaccess.Res, error) {
				calls++
				if calls <= len(c.failures) {
					return c.failures[calls-1], errors.New("HTTP Request failed")
				}
				return response(200, ""), nil
			})
			if calls != c.calls || (err != nil) != c.fails {
				t.Errorf("got %d calls, error %v, expected %d calls, error %v", calls, err, c.calls, c.fails)
			}
		})
	}
}

func TestReadAccessRuleAltConditions(t *testing.T) {
	ctx := context.Background()
	res := gjson.Parse(`{"ruleConditions":[
		{"attributeName":"umbrella.source.identity_ids","attributeValue":[1,2]},
		{"attributeName":"umbrella.source.identity_ids_shared","attributeValue":[2,3]},
		{"attributeName":"umbrella.source.identity_type_ids_shared","attributeValue":[7]}]}`)
	ids := func(v ...int64) types.Set {
		s, _ := types.SetValueFrom(ctx, types.Int64Type, v)
		return s
	}
	// Full read (import): both names are merged, alternative names alone are read
	data := AccessRule{SourceIdentityIds: ids(1, 2), SourceIdentityTypeIds: types.SetNull(types.Int64Type), SourceGeolocations: types.SetNull(types.StringType)}
	readAccessRuleAltConditions(ctx, &data, nil, res)
	if !data.SourceIdentityIds.Equal(ids(1, 2, 3)) || !data.SourceIdentityTypeIds.Equal(ids(7)) || !data.SourceGeolocations.IsNull() {
		t.Errorf("full read: got %v, %v, %v", data.SourceIdentityIds, data.SourceIdentityTypeIds, data.SourceGeolocations)
	}
	// Partial read (refresh): unmanaged attributes are left null
	prior := AccessRule{SourceIdentityIds: ids(1), SourceIdentityTypeIds: types.SetNull(types.Int64Type), SourceGeolocations: types.SetNull(types.StringType)}
	data = prior
	readAccessRuleAltConditions(ctx, &data, &prior, res)
	if !data.SourceIdentityIds.Equal(ids(1, 2, 3)) || !data.SourceIdentityTypeIds.IsNull() {
		t.Errorf("partial read: got %v, %v", data.SourceIdentityIds, data.SourceIdentityTypeIds)
	}
}

func TestExpandAccessRuleIdentityTypes(t *testing.T) {
	body := `{"ruleConditions":[{"attributeName":"umbrella.source.identity_type_ids","attributeOperator":"INTERSECT","attributeValue":[7,-1,-2,36]}]}`
	got := testAccessRuleCondition(expandAccessRuleIdentityTypes(body), tIT).Get("attributeValue").Raw
	if got != "[7,36,50,34,9]" {
		t.Errorf("unexpected identity types %s", got)
	}
}

func TestCollapseAccessRuleIdentityTypes(t *testing.T) {
	ctx := context.Background()
	ids := func(v ...int64) types.Set {
		s, _ := types.SetValueFrom(ctx, types.Int64Type, v)
		return s
	}
	for _, tc := range []struct {
		name  string
		prior *AccessRule
		want  types.Set
	}{
		{"import", nil, ids(7, -1, -2)},
		{"combined", &AccessRule{SourceIdentityTypeIds: ids(7, -1, -2)}, ids(7, -1, -2)},
		{"individual", &AccessRule{SourceIdentityTypeIds: ids(7, 36, 50, 34, 9)}, ids(7, 36, 50, 34, 9)},
		{"mixed", &AccessRule{SourceIdentityTypeIds: ids(7, 36, 50, -2)}, ids(7, 36, 50, -2)},
	} {
		data := AccessRule{SourceIdentityTypeIds: ids(7, 36, 50, 34, 9)}
		collapseAccessRuleIdentityTypes(ctx, &data, tc.prior)
		if !data.SourceIdentityTypeIds.Equal(tc.want) {
			t.Errorf("%s: got %v, want %v", tc.name, data.SourceIdentityTypeIds, tc.want)
		}
	}
	// Incomplete pairs are left as returned
	data := AccessRule{SourceIdentityTypeIds: ids(7, 36)}
	collapseAccessRuleIdentityTypes(ctx, &data, nil)
	if !data.SourceIdentityTypeIds.Equal(ids(7, 36)) {
		t.Errorf("incomplete: got %v", data.SourceIdentityTypeIds)
	}
}

func TestIdentityTypesStaticBody(t *testing.T) {
	res := gjson.ParseBytes(identityTypesStaticBody())
	if got := res.Get(`#(label=="iOS/Android Devices").id`).String(); got != "-1" {
		t.Errorf("unexpected combined id %q", got)
	}
	if res.Get(`#(label=="iOS Mobile Devices")`).Exists() {
		t.Errorf("unexpected individual identity type")
	}
	if int64(len(identityTypesStatic)+len(identityTypesCombined)) != res.Get("#").Int() {
		t.Errorf("identityTypesStatic modified")
	}
}

func TestAccessRuleApplicationCategoryAll(t *testing.T) {
	// Added to an existing condition, or as a new one, the expression references the condition once
	body := addAccessRuleApplicationCategoryAll(`{"ruleConditions":[{"attributeName":"` + tAC + `","attributeOperator":"INTERSECT","attributeValue":[102]}]}`)
	if got := testAccessRuleCondition(body, tAC).Get("attributeValue").Raw; got != "[102,90]" {
		t.Errorf("unexpected application categories %s", got)
	}
	body = addAccessRuleApplicationCategoryAll(testAccessRuleBody(tN))
	body = adjustAccessRuleConditions(body, "OR", "AND", "AND")
	if got := testAccessRuleCondition(body, tAC).Get("attributeValue").Raw; got != "[90]" {
		t.Errorf("unexpected application categories %s", got)
	}
	if got := testAccessRuleCondition(body, "umbrella.destination.logical_operator").Get("attributeValue").String(); got != tAC+" && "+tN {
		t.Errorf("unexpected destination expression %q", got)
	}

	ctx := context.Background()
	ids := func(v ...int64) types.Set {
		s, _ := types.SetValueFrom(ctx, types.Int64Type, v)
		return s
	}
	res := gjson.Parse(`{"ruleConditions":[{"attributeName":"` + tAC + `","attributeValue":[102,90]}]}`)
	// Full read (import)
	data := AccessRule{DestinationApplicationCategoryIds: ids(102, 90)}
	readAccessRuleApplicationCategoryAll(ctx, &data, nil, res)
	if !data.DestinationApplicationCategoryAll.ValueBool() || !data.DestinationApplicationCategoryIds.Equal(ids(102)) {
		t.Errorf("full read: got %v, %v", data.DestinationApplicationCategoryAll, data.DestinationApplicationCategoryIds)
	}
	// Partial read (refresh) of an unmanaged attribute
	prior := AccessRule{DestinationApplicationCategoryIds: ids(102), DestinationApplicationCategoryAll: types.BoolNull()}
	data = AccessRule{DestinationApplicationCategoryIds: ids(102, 90)}
	readAccessRuleApplicationCategoryAll(ctx, &data, &prior, res)
	if !data.DestinationApplicationCategoryAll.IsNull() || !data.DestinationApplicationCategoryIds.Equal(ids(102)) {
		t.Errorf("partial read: got %v, %v", data.DestinationApplicationCategoryAll, data.DestinationApplicationCategoryIds)
	}
	// Only all categories
	prior = AccessRule{DestinationApplicationCategoryIds: types.SetNull(types.Int64Type), DestinationApplicationCategoryAll: types.BoolValue(true)}
	data = AccessRule{DestinationApplicationCategoryIds: types.SetNull(types.Int64Type)}
	readAccessRuleApplicationCategoryAll(ctx, &data, &prior, gjson.Parse(`{"ruleConditions":[{"attributeName":"`+tAC+`","attributeValue":[90]}]}`))
	if !data.DestinationApplicationCategoryAll.ValueBool() || !data.DestinationApplicationCategoryIds.IsNull() {
		t.Errorf("all only: got %v, %v", data.DestinationApplicationCategoryAll, data.DestinationApplicationCategoryIds)
	}
}
