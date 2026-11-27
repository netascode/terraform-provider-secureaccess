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
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/netascode/go-secureaccess"
	"github.com/tidwall/gjson"
)

// This file is not generated. It tests our own code of the data source do_not_decrypt_list_application_exception_list by calling it directly: the schema, the mapping
// of a body to the model (the generated fromBody()), what Read() does with the answers of the API (the check of the settings, the check and the merge
// of the applications), the path of the request and the text of an error.
// There is no test with a fake of the API here. Such a test only shows that Read() understands what we think the API answers, so it cannot notice that the API changed.
// Paging and the answers of the API belong to the client (go-secureaccess). Whether Read() works against a real organization is for an acceptance test,
// which needs a way to create an exception application list (the API has only GET routes for them).

// excAppJSON is an application of a list as the API returns it (made up). Only the keys that fromBody() reads matter; selectedAction is not exposed.
const excAppJSON = `{"applicationId":1000,"applicationName":"App 0","action":"no_decrypt","hasUrls":true,"groupId":40,"groupName":"Group 0","selectedAction":{"name":"no_decrypt","displayName":"Allow","actionId":1}}`

// excSettingsJSON is the `data` of the answer of GET by id (made up), without the applications.
const excSettingsJSON = `{"id":5001,"organizationId":2222222,"isDefault":null,"name":"my_do_not_decrypt_list","type":"application","markedForDeletion":false,"blockedCount":0,"allowedCount":2,"warnedCount":0,"applicationsCategories":[]}`

func TestDoNotDecryptListApplicationExceptionListSchema(t *testing.T) {
	ctx := context.Background()

	t.Run("the attributes of a list", func(t *testing.T) {
		ds := &DoNotDecryptListApplicationExceptionListDataSource{}
		const typeName = "secureaccess_do_not_decrypt_list_application_exception_list"
		var meta datasource.MetadataResponse
		ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "secureaccess"}, &meta)
		if meta.TypeName != typeName {
			t.Errorf("type name = %q, want %q", meta.TypeName, typeName)
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
		if got, want := strings.Join(names, ","), "allowed_count,applications,blocked_count,id,name,warned_count"; got != want {
			t.Errorf("attributes = %s\nwant       %s", got, want)
		}
		// a lookup by name would need a read of every list, there is none
		if _, ok := datasource.DataSource(ds).(datasource.DataSourceWithConfigValidators); ok {
			t.Errorf("there must be no config validators")
		}
		var registered bool
		for _, factory := range New("test")().DataSources(ctx) {
			var m datasource.MetadataResponse
			factory().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "secureaccess"}, &m)
			if m.TypeName == typeName {
				registered = true
			}
		}
		if !registered {
			t.Error("the data source is not registered in the provider")
		}
	})

	t.Run("the nested attributes of an application", func(t *testing.T) {
		var resp datasource.SchemaResponse
		(&DoNotDecryptListApplicationExceptionListDataSource{}).Schema(ctx, datasource.SchemaRequest{}, &resp)
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
		if got, want := strings.Join(names, ","), "action,group_id,group_name,has_urls,id,name"; got != want {
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
}

func TestDoNotDecryptListApplicationExceptionListFromBody(t *testing.T) {
	ctx := context.Background()

	t.Run("every count and every field of an application is read from its own key", func(t *testing.T) {
		// every value is different, so that a field read from the key of another one is noticed
		var got DoNotDecryptListApplicationExceptionList
		got.fromBody(ctx, gjson.Parse(`{"name":"list","allowedCount":4,"blockedCount":2,"warnedCount":3,"applications":[
			{"applicationId":11,"applicationName":"name","groupId":22,"groupName":"group","hasUrls":true,"action":"action_x",
			 "selectedAction":{"name":"name_y","displayName":"display_z","actionId":33}}]}`))
		if got.Name.ValueString() != "list" || got.AllowedCount.ValueInt64() != 4 || got.BlockedCount.ValueInt64() != 2 || got.WarnedCount.ValueInt64() != 3 {
			t.Errorf("name=%s allowed=%s blocked=%s warned=%s, want list 4 2 3", got.Name, got.AllowedCount, got.BlockedCount, got.WarnedCount)
		}
		if len(got.Applications) != 1 {
			t.Fatalf("applications = %d, want 1", len(got.Applications))
		}
		a := got.Applications[0]
		if a.Id.ValueInt64() != 11 || a.Name.ValueString() != "name" || a.GroupId.ValueInt64() != 22 || a.GroupName.ValueString() != "group" || !a.HasUrls.ValueBool() || a.Action.ValueString() != "action_x" {
			t.Errorf("application = %+v, want every field from its own key (action is not selectedAction.name)", a)
		}
	})

	t.Run("a field that the API leaves out is null, not zero", func(t *testing.T) {
		var got DoNotDecryptListApplicationExceptionList
		got.fromBody(ctx, gjson.Parse(`{"id":5001,"applications":[{"applicationId":11}]}`))
		if !got.Name.IsNull() || !got.AllowedCount.IsNull() || !got.BlockedCount.IsNull() || !got.WarnedCount.IsNull() {
			t.Errorf("settings = %s %s %s %s, want null", got.Name, got.AllowedCount, got.BlockedCount, got.WarnedCount)
		}
		if len(got.Applications) != 1 {
			t.Fatalf("applications = %d, want 1", len(got.Applications))
		}
		a := got.Applications[0]
		if a.Id.ValueInt64() != 11 || !a.Name.IsNull() || !a.GroupId.IsNull() || !a.GroupName.IsNull() || !a.HasUrls.IsNull() || !a.Action.IsNull() {
			t.Errorf("application = %+v, want only the id set", a)
		}
	})

	t.Run("a null in a field is read as the zero value", func(t *testing.T) {
		// said in the YAML: the generated fromBody() reads a JSON null as 0, "" or false. Not seen on the int environment.
		var got DoNotDecryptListApplicationExceptionList
		got.fromBody(ctx, gjson.Parse(`{"applications":[{"applicationId":11,"groupId":null,"groupName":null,"hasUrls":null}]}`))
		a := got.Applications[0]
		if a.GroupId.IsNull() || a.GroupId.ValueInt64() != 0 || a.GroupName.IsNull() || a.GroupName.ValueString() != "" || a.HasUrls.IsNull() || a.HasUrls.ValueBool() {
			t.Errorf("application = %+v, want 0, an empty string and false", a)
		}
	})

	t.Run("a list without applications is an empty list, not a null", func(t *testing.T) {
		var got DoNotDecryptListApplicationExceptionList
		got.fromBody(ctx, gjson.Parse(`{"name":"list","applications":[]}`))
		if got.Applications == nil || len(got.Applications) != 0 {
			t.Errorf("applications = %#v, want an empty list", got.Applications)
		}
	})

	t.Run("the applications are in the order of the API", func(t *testing.T) {
		var got DoNotDecryptListApplicationExceptionList
		got.fromBody(ctx, gjson.Parse(`{"applications":[{"applicationId":30},{"applicationId":10},{"applicationId":20}]}`))
		var ids []int64
		for _, a := range got.Applications {
			ids = append(ids, a.Id.ValueInt64())
		}
		if fmt.Sprint(ids) != "[30 10 20]" {
			t.Errorf("ids = %v, want [30 10 20]", ids)
		}
	})

	t.Run("the id is the id of the config, not the id in the answer", func(t *testing.T) {
		got := DoNotDecryptListApplicationExceptionList{Id: types.StringValue("5001")}
		got.fromBody(ctx, gjson.Parse(`{"id":7777,"name":"list"}`))
		if got.Id.ValueString() != "5001" {
			t.Errorf("id = %q, want the id of the config", got.Id.ValueString())
		}
	})
}

func TestExceptionApplicationListPath(t *testing.T) {
	const prefix = "/policies/v2/exceptionApplicationLists/"
	for _, tc := range []struct{ name, id, want string }{
		{"a number", "5001", prefix + "5001"},
		{"slashes cannot reach another path", "1/../../threatCategories", prefix + "1%2F..%2F..%2FthreatCategories"},
		{"a ? adds no query", "5001?x=1", prefix + "5001%3Fx%3D1"},
		{"a # adds no fragment", "5001#x", prefix + "5001%23x"},
		{"a % is escaped, an escaped slash stays escaped", "5001%2F..%2FthreatCategories", prefix + "5001%252F..%252FthreatCategories"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := exceptionApplicationListPath(tc.id)
			if got != tc.want {
				t.Errorf("path = %q, want %q", got, tc.want)
			}
			// whatever the id is, it is one path segment: no slash, no query, no fragment after the collection
			if rest := strings.TrimPrefix(got, prefix); strings.ContainsAny(rest, "/?#") {
				t.Errorf("the id is not one path segment: %q", rest)
			}
		})
	}
}

func TestExceptionApplicationListSettings(t *testing.T) {
	t.Run("the settings are the object in data", func(t *testing.T) {
		got, err := exceptionApplicationListSettings(gjson.Parse(`{"status":{"code":200},"data":` + excSettingsJSON + `}`))
		if err != nil {
			t.Fatal(err)
		}
		if got.Get("name").String() != "my_do_not_decrypt_list" || got.Get("allowedCount").Int() != 2 {
			t.Errorf("settings = %s", got.Raw)
		}
	})

	for _, tc := range []struct{ name, body, got string }{
		{"no data", `{"status":{"code":200}}`, "Null"},
		{"data is a list", `{"data":[]}`, "JSON"},
		{"data is a string", `{"data":"x"}`, "String"},
		{"data is a number", `{"data":5}`, "Number"},
		{"data is null", `{"data":null}`, "Null"},
	} {
		t.Run(tc.name+" is an error", func(t *testing.T) {
			_, err := exceptionApplicationListSettings(gjson.Parse(tc.body))
			if want := "expected an object in `data`, got: " + tc.got; err == nil || err.Error() != want {
				t.Errorf("error = %v, want %q", err, want)
			}
		})
	}
}

func TestExceptionApplicationListBody(t *testing.T) {
	ctx := context.Background()
	settings := gjson.Parse(excSettingsJSON)
	second := `{"applicationId":1003,"applicationName":"App 1","action":"no_decrypt","hasUrls":false,"groupId":41,"groupName":"Group 1"}`

	t.Run("the applications are put into the settings as applications", func(t *testing.T) {
		body, err := exceptionApplicationListBody(settings, gjson.Parse("["+excAppJSON+","+second+"]"))
		if err != nil {
			t.Fatal(err)
		}
		if got := gjson.Get(body, "applications.#").Int(); got != 2 {
			t.Errorf("applications = %d, want 2: %s", got, body)
		}
		if got := gjson.Get(body, "applications.1.applicationId").Int(); got != 1003 {
			t.Errorf("the second application has the id %d, want 1003 (the order is kept)", got)
		}
		// the settings stay as they are
		if got := gjson.Get(body, "name").String(); got != "my_do_not_decrypt_list" {
			t.Errorf("name = %q, want the name of the settings", got)
		}
		if !gjson.Get(body, "applicationsCategories").IsArray() {
			t.Errorf("the other settings are lost: %s", body)
		}
	})

	t.Run("no applications are an empty array, not a missing key", func(t *testing.T) {
		body, err := exceptionApplicationListBody(settings, gjson.Parse("[]"))
		if err != nil {
			t.Fatal(err)
		}
		if got := gjson.Get(body, "applications"); !got.IsArray() || len(got.Array()) != 0 {
			t.Errorf("applications = %s, want []", got.Raw)
		}
	})

	t.Run("the body is what fromBody reads", func(t *testing.T) {
		body, err := exceptionApplicationListBody(settings, gjson.Parse("["+excAppJSON+","+second+"]"))
		if err != nil {
			t.Fatal(err)
		}
		var got DoNotDecryptListApplicationExceptionList
		got.fromBody(ctx, gjson.Parse(body))
		if got.Name.ValueString() != "my_do_not_decrypt_list" || got.AllowedCount.ValueInt64() != 2 || got.BlockedCount.ValueInt64() != 0 || got.WarnedCount.ValueInt64() != 0 {
			t.Errorf("name=%s allowed=%s blocked=%s warned=%s", got.Name, got.AllowedCount, got.BlockedCount, got.WarnedCount)
		}
		if len(got.Applications) != 2 {
			t.Fatalf("applications = %d, want 2", len(got.Applications))
		}
		first, last := got.Applications[0], got.Applications[1]
		if first.Id.ValueInt64() != 1000 || first.Name.ValueString() != "App 0" || first.GroupId.ValueInt64() != 40 || first.GroupName.ValueString() != "Group 0" || !first.HasUrls.ValueBool() || first.Action.ValueString() != "no_decrypt" {
			t.Errorf("first application = %+v", first)
		}
		if last.Id.ValueInt64() != 1003 || last.Name.ValueString() != "App 1" || last.HasUrls.ValueBool() {
			t.Errorf("last application = %+v", last)
		}
	})

	for _, tc := range []struct{ name, bad, got string }{
		{"a string", `"x"`, "String"}, {"null", `null`, "Null"}, {"a number", `5`, "Number"}, {"a bool", `true`, "True"}, {"a list", `[1]`, "JSON"},
	} {
		t.Run("an application that is "+tc.name+" is an error", func(t *testing.T) {
			// the bad one is the first of two and the last of two: the index in the error is its position
			for i, applications := range []string{"[" + tc.bad + "," + excAppJSON + "]", "[" + excAppJSON + "," + tc.bad + "]"} {
				_, err := exceptionApplicationListBody(settings, gjson.Parse(applications))
				if want := fmt.Sprintf("expected an object as application %d, got: %s", i, tc.got); err == nil || err.Error() != want {
					t.Errorf("%s: error = %v, want %q", applications, err, want)
				}
			}
		})
	}

	t.Run("an application late in a long list is checked too", func(t *testing.T) {
		items := make([]string, 125)
		for i := range items {
			items[i] = excAppJSON
		}
		items[120] = `"x"`
		_, err := exceptionApplicationListBody(settings, gjson.Parse("["+strings.Join(items, ",")+"]"))
		if want := "expected an object as application 120, got: String"; err == nil || err.Error() != want {
			t.Errorf("error = %v, want %q", err, want)
		}
	})
}

func TestExceptionApplicationListError(t *testing.T) {
	err := errors.New("StatusCode 500")

	t.Run("the answer of the API is added to the error of the client", func(t *testing.T) {
		got := exceptionApplicationListError(err, secureaccess.Res{Result: gjson.Parse(`{"message":"internal"}`)})
		if want := `StatusCode 500, answer of the API: {"message":"internal"}`; got != want {
			t.Errorf("error = %q, want %q", got, want)
		}
	})

	t.Run("without an answer there is only the error of the client", func(t *testing.T) {
		got := exceptionApplicationListError(err, secureaccess.Res{})
		if want := "StatusCode 500"; got != want {
			t.Errorf("error = %q, want %q", got, want)
		}
	})
}
