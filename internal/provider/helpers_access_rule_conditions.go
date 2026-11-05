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
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Helpers for the `ruleConditions` of the access rule, see the API behaviour documented in gen/definitions/access_rule.yaml.

const (
	accessRuleSource      = "source"
	accessRuleDestination = "destination"

	accessRuleOperatorAnd = "AND"
	accessRuleOperatorOr  = "OR"

	accessRulePrivateNetwork = "PRIVATE_NETWORK"
	accessRulePublicInternet = "PUBLIC_INTERNET"
)

// accessRuleGroup is the type of a condition. In a logical operator expression, conditions of the same group are
// ORed within one term, different groups are ANDed. `groupOther` conditions are not supported in expressions.
type accessRuleGroup string

const (
	groupIdentity      accessRuleGroup = "identities"
	groupIdentityType  accessRuleGroup = "identity types"
	groupGeolocation   accessRuleGroup = "geolocations"
	groupSourceLiteral accessRuleGroup = "source literals"
	groupNetwork       accessRuleGroup = "network objects"
	groupService       accessRuleGroup = "service objects"
	groupCategory      accessRuleGroup = "content categories"
	groupApplication   accessRuleGroup = "applications"
	groupOther         accessRuleGroup = ""
)

// accessRuleDestinationGroupRank is the order of the destination terms in a built expression
var accessRuleDestinationGroupRank = map[accessRuleGroup]int{groupApplication: 0, groupCategory: 1, groupNetwork: 2, groupService: 3}

// accessRuleDestinationTermRank returns the rank of a term key, a group or the name of a `groupOther` condition
func accessRuleDestinationTermRank(key string) int {
	if rank, ok := accessRuleDestinationGroupRank[accessRuleGroup(key)]; ok {
		return rank
	}
	return len(accessRuleDestinationGroupRank)
}

type accessRuleCondition struct {
	tfName      string
	name        string // `attributeName` sent by the provider
	altName     string // alternative `attributeName` the API may return for the same condition (e.g. written by the portal)
	side        string
	group       accessRuleGroup
	trafficType string // only allowed for this `umbrella.default.traffic` value, empty if allowed for both
}

// accessRuleConditions lists all conditions, in the order of the resource attributes
var accessRuleConditions = []accessRuleCondition{
	{"source_identity_ids", "umbrella.source.identity_ids", "umbrella.source.identity_ids_shared", accessRuleSource, groupIdentity, ""},
	{"source_identity_type_ids", "umbrella.source.identity_type_ids", "umbrella.source.identity_type_ids_shared", accessRuleSource, groupIdentityType, ""},
	{"source_geolocations", "umbrella.source.geolocations_shared", "umbrella.source.geolocations", accessRuleSource, groupGeolocation, ""},
	{"source_network_object_ids", "umbrella.source.networkObjectIds", "", accessRuleSource, groupOther, ""},
	{"source_network_object_group_ids", "umbrella.source.networkObjectGroupIds", "", accessRuleSource, groupOther, ""},
	{"source_private_literals", "umbrella.source.ip_address", "", accessRuleSource, groupOther, accessRulePrivateNetwork},
	{"source_internet_literals", "umbrella.source.composite_inline_ip", "", accessRuleSource, groupSourceLiteral, accessRulePublicInternet},
	{"destination_network_object_ids", "umbrella.destination.networkObjectIds", "", accessRuleDestination, groupNetwork, ""},
	{"destination_network_object_group_ids", "umbrella.destination.networkObjectGroupIds", "", accessRuleDestination, groupNetwork, ""},
	{"destination_service_object_ids", "umbrella.destination.serviceObjectIds", "", accessRuleDestination, groupService, ""},
	{"destination_service_object_group_ids", "umbrella.destination.serviceObjectGroupIds", "", accessRuleDestination, groupService, ""},
	{"destination_private_resource_ids", "umbrella.destination.private_resource_ids", "", accessRuleDestination, groupOther, accessRulePrivateNetwork},
	{"destination_private_resource_group_ids", "umbrella.destination.private_resource_group_ids", "", accessRuleDestination, groupOther, accessRulePrivateNetwork},
	{"destination_security_group_tag_ids", "umbrella.destination.security_group_tag_ids", "", accessRuleDestination, groupOther, accessRulePrivateNetwork},
	{"destination_content_category_ids", "umbrella.destination.category_ids", "", accessRuleDestination, groupCategory, accessRulePublicInternet},
	{"destination_content_category_list_ids", "umbrella.destination.category_list_ids", "", accessRuleDestination, groupCategory, accessRulePublicInternet},
	{"destination_application_ids", "umbrella.destination.application_ids", "", accessRuleDestination, groupApplication, accessRulePublicInternet},
	{"destination_application_category_ids", "umbrella.destination.application_category_ids", "", accessRuleDestination, groupApplication, accessRulePublicInternet},
	{"destination_application_list_ids", "umbrella.destination.application_list_ids", "", accessRuleDestination, groupApplication, accessRulePublicInternet},
	{"destination_list_ids", "umbrella.destination.destination_list_ids", "", accessRuleDestination, groupOther, accessRulePublicInternet},
	{"destination_literals", "umbrella.destination.composite_inline_ip", "", accessRuleDestination, groupOther, ""},
	{"destination_geolocations", "umbrella.destination.geolocations", "", accessRuleDestination, groupOther, accessRulePublicInternet},
}

func accessRuleAllName(side string) string {
	return "umbrella." + side + ".all"
}

func accessRuleLogicalOperatorName(side string) string {
	return "umbrella." + side + ".logical_operator"
}

// buildAccessRuleExpression returns the logical operator expression ANDing the groups of `conditions`.
// With `mergeNetworkService`, network and service objects are ORed within a single term.
func buildAccessRuleExpression(side string, conditions []accessRuleCondition, mergeNetworkService bool) string {
	termOf := func(c accessRuleCondition) string {
		if c.group == groupOther {
			// Not supported in expressions, kept as a separate term to let the API report it
			return c.name
		}
		if mergeNetworkService && c.group == groupService {
			return string(groupNetwork)
		}
		return string(c.group)
	}

	var keys []string
	names := map[string][]string{}
	for _, c := range conditions {
		k := termOf(c)
		if _, ok := names[k]; !ok {
			keys = append(keys, k)
		}
		names[k] = append(names[k], c.name)
	}
	if side == accessRuleDestination {
		sort.SliceStable(keys, func(i, j int) bool {
			return accessRuleDestinationTermRank(keys[i]) < accessRuleDestinationTermRank(keys[j])
		})
	}

	terms := make([]string, 0, len(keys))
	for _, k := range keys {
		if len(names[k]) == 1 {
			terms = append(terms, names[k][0])
		} else {
			terms = append(terms, "("+strings.Join(names[k], " || ")+")")
		}
	}
	return strings.Join(terms, " && ")
}

// adjustAccessRuleConditions completes the `ruleConditions` of a request body: a side without any condition matches
// all (`umbrella.<side>.all`), a side with operator AND gets the logical operator expression.
func adjustAccessRuleConditions(body, sourceOperator, destinationOperator, networkServiceOperator string) string {
	present := map[string]bool{}
	for _, c := range gjson.Get(body, "ruleConditions").Array() {
		present[c.Get("attributeName").String()] = true
	}
	operators := map[string]string{accessRuleSource: sourceOperator, accessRuleDestination: destinationOperator}
	for _, side := range []string{accessRuleSource, accessRuleDestination} {
		var conditions []accessRuleCondition
		for _, c := range accessRuleConditions {
			if c.side == side && present[c.name] {
				conditions = append(conditions, c)
			}
		}
		if len(conditions) == 0 {
			body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": accessRuleAllName(side), "attributeOperator": "=", "attributeValue": true})
		} else if operators[side] == accessRuleOperatorAnd {
			merge := side == accessRuleDestination && networkServiceOperator == accessRuleOperatorOr
			expression := buildAccessRuleExpression(side, conditions, merge)
			body, _ = sjson.Set(body, "ruleConditions.-1", map[string]any{"attributeName": accessRuleLogicalOperatorName(side), "attributeOperator": accessRuleOperatorAnd, "attributeValue": expression})
		}
	}
	return body
}

// accessRuleLogicalOperator returns the logical operator expression of a side from a response, empty if there is none
func accessRuleLogicalOperator(res gjson.Result, side string) string {
	return res.Get(`ruleConditions.#(attributeName=="` + accessRuleLogicalOperatorName(side) + `").attributeValue`).String()
}

// splitAccessRuleExpression returns the top-level AND terms of an expression, without enclosing parentheses
func splitAccessRuleExpression(expression string) []string {
	var terms []string
	depth, start := 0, 0
	expression = trimAccessRuleParentheses(expression)
	for i := 0; i < len(expression); i++ {
		switch expression[i] {
		case '(':
			depth++
		case ')':
			depth--
		case '&':
			if depth == 0 && i+1 < len(expression) && expression[i+1] == '&' {
				terms = append(terms, trimAccessRuleParentheses(expression[start:i]))
				start = i + 2
				i++
			}
		}
	}
	return append(terms, trimAccessRuleParentheses(expression[start:]))
}

// trimAccessRuleParentheses removes whitespace and parentheses enclosing the whole expression, e.g. `((a && b))`
func trimAccessRuleParentheses(expression string) string {
	for {
		expression = strings.TrimSpace(expression)
		if len(expression) < 2 || expression[0] != '(' || expression[len(expression)-1] != ')' {
			return expression
		}
		// Make sure the first parenthesis closes at the very end, i.e. not `(a) && (b)`
		depth := 0
		for i := 0; i < len(expression)-1; i++ {
			if expression[i] == '(' {
				depth++
			} else if expression[i] == ')' {
				depth--
			}
			if depth == 0 {
				return expression
			}
		}
		expression = expression[1 : len(expression)-1]
	}
}

// accessRuleSetField returns the Set attribute `tfName` of an access rule model, nil if there is none
func accessRuleSetField(data *AccessRule, tfName string) *types.Set {
	v := reflect.ValueOf(data).Elem()
	for i := range v.NumField() {
		if v.Type().Field(i).Tag.Get("tfsdk") == tfName {
			value, _ := v.Field(i).Addr().Interface().(*types.Set)
			return value
		}
	}
	return nil
}

// readAccessRuleAltConditions merges the values of the conditions returned under their name and alternative name
// into their attribute. In a partial read (`prior` not nil), attributes that were null before the read are left null.
func readAccessRuleAltConditions(ctx context.Context, data, prior *AccessRule, res gjson.Result) {
	for _, c := range accessRuleConditions {
		if c.altName == "" {
			continue
		}
		value := accessRuleSetField(data, c.tfName)
		if value == nil || prior != nil && accessRuleSetField(prior, c.tfName).IsNull() {
			continue
		}
		seen := map[string]bool{}
		var elements []gjson.Result
		for _, name := range []string{c.name, c.altName} {
			for _, e := range res.Get(`ruleConditions.#(attributeName=="` + name + `").attributeValue`).Array() {
				if !seen[e.Raw] {
					seen[e.Raw] = true
					elements = append(elements, e)
				}
			}
		}
		elementType := value.ElementType(ctx)
		switch {
		case len(elements) == 0:
			*value = types.SetNull(elementType)
		case elementType == types.StringType:
			*value = helpers.GetStringSet(elements)
		default:
			*value = helpers.GetInt64Set(elements)
		}
	}
}

// accessRuleNetworkServiceOperator returns OR if a destination expression has a term matching both network and
// service objects (e.g. `category_ids && (networkObjectIds || serviceObjectIds)`), AND otherwise
func accessRuleNetworkServiceOperator(expression string) string {
	groups := map[string]accessRuleGroup{}
	for _, c := range accessRuleConditions {
		if c.side == accessRuleDestination {
			groups[c.name] = c.group
		}
	}
	for _, term := range splitAccessRuleExpression(expression) {
		var network, service bool
		for operand := range strings.SplitSeq(term, "||") {
			// Parentheses within a term only group operands of the same OR
			switch groups[strings.Trim(operand, "() \t\n")] {
			case groupNetwork:
				network = true
			case groupService:
				service = true
			}
		}
		if network && service {
			return accessRuleOperatorOr
		}
	}
	return accessRuleOperatorAnd
}

// validateAccessRuleConditions checks the conditions of an access rule configuration against the rule access type
// and the logical operators, to report at plan time what the API would reject. Unknown values are skipped.
func validateAccessRuleConditions(ctx context.Context, config tfsdk.Config, diags *diag.Diagnostics) {
	getString := func(name string) (string, bool) {
		var v types.String
		diags.Append(config.GetAttribute(ctx, path.Root(name), &v)...)
		return v.ValueString(), !v.IsUnknown()
	}
	trafficType, trafficTypeKnown := getString("traffic_type")
	action, actionKnown := getString("action")
	sourceOperator, sourceOperatorKnown := getString("source_operator")
	destinationOperator, destinationOperatorKnown := getString("destination_operator")
	networkServiceOperator, networkServiceOperatorKnown := getString("destination_operator_network_service")
	if diags.HasError() {
		return
	}
	// Optional attributes with a default are null in the configuration if not set
	if sourceOperator == "" {
		sourceOperator = accessRuleOperatorOr
	}
	if destinationOperator == "" {
		destinationOperator = accessRuleOperatorOr
	}
	if networkServiceOperator == "" {
		networkServiceOperator = accessRuleOperatorAnd
	}

	// Configured conditions per side and group
	configured := map[string][]accessRuleCondition{}
	groups := map[string]map[accessRuleGroup]bool{accessRuleSource: {}, accessRuleDestination: {}}
	for _, c := range accessRuleConditions {
		// Conditions are Sets, or Lists for conditions with a single item
		var v attr.Value
		diags.Append(config.GetAttribute(ctx, path.Root(c.tfName), &v)...)
		if v == nil || v.IsNull() {
			continue
		}
		if elements, ok := v.(interface{ Elements() []attr.Value }); ok && !v.IsUnknown() && len(elements.Elements()) == 0 {
			diags.AddAttributeError(path.Root(c.tfName), "Invalid Attribute Value", fmt.Sprintf("%s must not be empty, omit it instead.", c.tfName))
		}
		if trafficTypeKnown && c.trafficType != "" && trafficType != c.trafficType {
			diags.AddAttributeError(path.Root(c.tfName), "Invalid Attribute Combination", fmt.Sprintf("%s is only supported with traffic_type %q.", c.tfName, c.trafficType))
		}
		configured[c.side] = append(configured[c.side], c)
		groups[c.side][c.group] = true
	}

	// Lists are exclusive with their individual entries (portal restriction, not enforced by the API)
	for _, e := range [][2]string{
		{"destination_application_list_ids", "destination_application_ids"},
		{"destination_application_list_ids", "destination_application_category_ids"},
		{"destination_content_category_list_ids", "destination_content_category_ids"},
	} {
		if accessRuleConfigured(configured[accessRuleDestination], e[0]) && accessRuleConfigured(configured[accessRuleDestination], e[1]) {
			diags.AddAttributeError(path.Root(e[1]), "Invalid Attribute Combination", fmt.Sprintf("%s cannot be combined with %s.", e[1], e[0]))
		}
	}

	// Destinations restricted with source internet literals
	if groups[accessRuleSource][groupSourceLiteral] {
		if groups[accessRuleDestination][groupService] {
			diags.AddAttributeError(path.Root("source_internet_literals"), "Invalid Attribute Combination", "source_internet_literals cannot be combined with destination service objects or service object groups.")
		}
		var literals types.Set
		diags.Append(config.GetAttribute(ctx, path.Root("destination_literals"), &literals)...)
		if !literals.IsNull() && !literals.IsUnknown() {
			if len(literals.Elements()) > 1 {
				diags.AddAttributeError(path.Root("destination_literals"), "Invalid Attribute Combination", "destination_literals cannot have more than one item with source_internet_literals.")
			}
			for _, e := range literals.Elements() {
				if o, ok := e.(types.Object); ok {
					if protocol, ok := o.Attributes()["protocol"].(types.String); ok && protocol.ValueString() == "ICMP" {
						diags.AddAttributeError(path.Root("destination_literals"), "Invalid Attribute Combination", "destination_literals protocol \"ICMP\" cannot be combined with source_internet_literals.")
					}
				}
			}
		}
	}

	// Destination geolocations
	if accessRuleConfigured(configured[accessRuleDestination], "destination_geolocations") {
		if actionKnown && action != "block" {
			diags.AddAttributeError(path.Root("destination_geolocations"), "Invalid Attribute Combination", "destination_geolocations is only supported with action \"block\".")
		}
		if len(configured[accessRuleDestination]) > 1 {
			diags.AddAttributeError(path.Root("destination_geolocations"), "Invalid Attribute Combination", "destination_geolocations cannot be combined with other destinations.")
		}
		if groups[accessRuleSource][groupGeolocation] {
			diags.AddAttributeError(path.Root("destination_geolocations"), "Invalid Attribute Combination", "destination_geolocations cannot be combined with source_geolocations.")
		}
	}

	if trafficTypeKnown && actionKnown && trafficType == accessRulePrivateNetwork && (action == "warn" || action == "isolate") {
		diags.AddAttributeError(path.Root("action"), "Invalid Attribute Combination", fmt.Sprintf("action %q is only supported with traffic_type %q.", action, accessRulePublicInternet))
	}

	// Source logical operator
	if sourceOperatorKnown && sourceOperator == accessRuleOperatorAnd {
		if groups[accessRuleSource][groupOther] {
			diags.AddAttributeError(path.Root("source_operator"), "Invalid Attribute Combination", "source_operator \"AND\" does not support source_private_literals, source network objects or network object groups.")
		} else if len(configured[accessRuleSource]) < 2 {
			diags.AddAttributeError(path.Root("source_operator"), "Invalid Attribute Combination", "source_operator \"AND\" requires at least two of source_identity_ids, source_identity_type_ids, source_geolocations and source_internet_literals.")
		}
	}
	if sourceOperatorKnown && sourceOperator == accessRuleOperatorOr && groups[accessRuleSource][groupGeolocation] && (groups[accessRuleSource][groupIdentity] || groups[accessRuleSource][groupIdentityType]) {
		diags.AddAttributeError(path.Root("source_geolocations"), "Invalid Attribute Combination", "source_geolocations combined with source identities requires source_operator \"AND\".")
	}

	// Destination logical operator
	if destinationOperatorKnown && networkServiceOperatorKnown && networkServiceOperator == accessRuleOperatorOr {
		if destinationOperator != accessRuleOperatorAnd {
			diags.AddAttributeError(path.Root("destination_operator_network_service"), "Invalid Attribute Combination", "destination_operator_network_service \"OR\" is only supported with destination_operator \"AND\".")
		} else if !groups[accessRuleDestination][groupNetwork] || !groups[accessRuleDestination][groupService] || !groups[accessRuleDestination][groupCategory] && !groups[accessRuleDestination][groupApplication] {
			diags.AddAttributeError(path.Root("destination_operator_network_service"), "Invalid Attribute Combination", "destination_operator_network_service \"OR\" requires network objects, service objects and content categories or applications.")
		}
	}
	if destinationOperatorKnown && destinationOperator == accessRuleOperatorAnd {
		terms := 0
		for _, g := range []accessRuleGroup{groupNetwork, groupService, groupCategory, groupApplication} {
			if groups[accessRuleDestination][g] {
				terms++
			}
		}
		if groups[accessRuleDestination][groupNetwork] && groups[accessRuleDestination][groupService] && networkServiceOperator == accessRuleOperatorOr {
			terms--
		}
		if groups[accessRuleDestination][groupOther] {
			diags.AddAttributeError(path.Root("destination_operator"), "Invalid Attribute Combination", "destination_operator \"AND\" does not support destination_literals, destination_geolocations, Destination Lists, Private Resources, Private Resource Groups and Security Group Tags.")
		} else if terms < 2 {
			diags.AddAttributeError(path.Root("destination_operator"), "Invalid Attribute Combination", fmt.Sprintf("destination_operator \"AND\" requires conditions of at least two of: %s, %s, %s, %s.", groupNetwork, groupService, groupCategory, groupApplication))
		}
	}
}

func accessRuleConfigured(conditions []accessRuleCondition, tfName string) bool {
	for _, c := range conditions {
		if c.tfName == tfName {
			return true
		}
	}
	return false
}

// accessRuleInvalidIdsRetries is the number of retries of a rule POST/PUT rejected for referencing IDs unknown to the
// policy service, which happens for objects created shortly before (e.g. in the same apply)
const accessRuleInvalidIdsRetries = 3

// accessRuleInvalidIdsDelay is multiplied by the attempt number to get the delay before a retry
var accessRuleInvalidIdsDelay = 10 * time.Second

// accessRuleInvalidIdsMessages are parts of the `msg` of the 400 responses for IDs unknown to the policy service
var accessRuleInvalidIdsMessages = []string{
	"the ID's provided for",
	"contains invalid network object",
	"contains invalid service object",
	"were not found for organization",
}

func isAccessRuleInvalidIdsError(res secureaccess.Res) bool {
	if res.StatusCode != 400 {
		return false
	}
	msg := res.Get("msg").String()
	for _, m := range accessRuleInvalidIdsMessages {
		if strings.Contains(msg, m) {
			return true
		}
	}
	return false
}

// retryAccessRuleInvalidIds calls `write` (a rule POST or PUT) and retries it if the API rejects referenced IDs
func retryAccessRuleInvalidIds(ctx context.Context, write func() (secureaccess.Res, error)) (secureaccess.Res, error) {
	for attempt := 1; ; attempt++ {
		res, err := write()
		if err == nil || attempt > accessRuleInvalidIdsRetries || !isAccessRuleInvalidIdsError(res) {
			return res, err
		}
		delay := time.Duration(attempt) * accessRuleInvalidIdsDelay
		tflog.Warn(ctx, fmt.Sprintf("Access rule references IDs unknown to the API, retry %d/%d in %v: %s", attempt, accessRuleInvalidIdsRetries, delay, res.Get("msg").String()))
		select {
		case <-ctx.Done():
			return res, err
		case <-time.After(delay):
		}
	}
}
