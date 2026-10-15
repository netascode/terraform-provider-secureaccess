//go:build ignore
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

// Section below is generated&owned by "gen/generator.go". //template:begin imports
import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &{{camelCase .Name}}DataSource{}
	_ datasource.DataSourceWithConfigure = &{{camelCase .Name}}DataSource{}
)

func New{{camelCase .Name}}DataSource() datasource.DataSource {
	return &{{camelCase .Name}}DataSource{}
}

type {{camelCase .Name}}DataSource struct {
	client *secureaccess.Client
}

func (d *{{camelCase .Name}}DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_{{snakeCase .Name}}"
}

func (d *{{camelCase .Name}}DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("{{.DsDescription}}")
		{{- if .DeprecationMessage -}}
		.AddAttributeDescription("{{.DeprecationMessage}}")
		{{- end -}}
		.String,
		{{- if .DeprecationMessage }}
		DeprecationMessage:  helpers.NewAttributeDescription("{{.DeprecationMessage}}").String,
		{{- end}}

		Attributes: map[string]schema.Attribute{
			{{- if not .NoId}}
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the object",
				{{- if .SettingsList}}
				Optional:            true,
				Computed:            true,
				{{- else if not (hasDataSourceQuery .Attributes) }}
				Required:            true,
				{{- else}}
				Optional:            true,
				Computed:            true,
				{{- end}}
			},
			{{- end}}
			{{- range .Attributes}}
			{{- if not .Value}}
			"{{.TfName}}": schema.{{if isNestedListMapSet .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
				MarkdownDescription: "{{if .DsDescription}}{{.DsDescription}}{{else}}{{.Description}}{{end}}",
				{{- if isListSet .}}
				ElementType:         types.{{.ElementType}}Type,
				{{- end}}
				{{- if .Reference}}
				Required:            true,
				{{- else}}
				{{- if or .DataSourceQuery .DataSourceOptionalParameter }}
				Optional:            true,
				{{- end}}
				{{- if isNestedMap .}}
				Optional:            true,
				{{- end}}
				Computed:            true,
				{{- end}}
				{{- if .Sensitive}}
				Sensitive:           true,
				{{- end}}
				{{- if isNestedListMapSet .}}
				{{- $parentNestedMap := isNestedMap .}}
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						{{- range .Attributes}}
						{{- if not .Value}}
						"{{.TfName}}": schema.{{if isNestedListMapSet .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
							MarkdownDescription: "{{if .DsDescription}}{{.DsDescription}}{{else}}{{.Description}}{{end}}",
							{{- if isListSet .}}
							ElementType:         types.{{.ElementType}}Type,
							{{- end}}
							{{- if and .ResourceId $parentNestedMap (not $.IsBulk) }}
							Required:            true,
							{{- else}}
							Computed:            true,
							{{- end}}
							{{- if .Sensitive}}
							Sensitive:           true,
							{{- end}}
							{{- if isNestedListMapSet .}}
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									{{- range .Attributes}}
									{{- if not .Value}}
									"{{.TfName}}": schema.{{if isNestedListMapSet .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
										MarkdownDescription: "{{if .DsDescription}}{{.DsDescription}}{{else}}{{.Description}}{{end}}",
										{{- if isListSet .}}
										ElementType:         types.{{.ElementType}}Type,
										{{- end}}
										Computed:            true,
										{{- if .Sensitive}}
										Sensitive:           true,
										{{- end}}
										{{- if isNestedListMapSet .}}
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												{{- range .Attributes}}
												{{- if not .Value}}
												"{{.TfName}}": schema.{{if isNestedListMapSet .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
													MarkdownDescription: "{{if .DsDescription}}{{.DsDescription}}{{else}}{{.Description}}{{end}}",
													{{- if isListSet .}}
													ElementType:         types.{{.ElementType}}Type,
													{{- end}}
													Computed:            true,
													{{- if .Sensitive}}
													Sensitive:           true,
													{{- end}}
												},
												{{- end}}
												{{- end}}
											},
										},
										{{- end}}
									},
									{{- end}}
									{{- end}}
								},
							},
							{{- end}}
						},
						{{- end}}
						{{- end}}
					},
				},
				{{- end}}
			},
			{{- end}}
			{{- end}}
		},
	}
}
{{- $dataSourceAttributes := getDataSourceQueryAttributes .}}
{{- if hasDataSourceQuery .Attributes }}
func (d *{{camelCase .Name}}DataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
    return []datasource.ConfigValidator{
        datasourcevalidator.ExactlyOneOf(
            path.MatchRoot("id"),
			{{- range $index, $attr := $dataSourceAttributes }}
			path.MatchRoot("{{$attr.TfName}}"),
			{{- end }}
        ),
    }
}
{{- end}}

func (d *{{camelCase .Name}}DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*SecureAccessProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (d *{{camelCase .Name}}DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config {{camelCase .Name}}

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set request modifiers
	reqMods := [](func(*secureaccess.Req)){}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", {{if .NoId}}"{{.Name}}"{{else}}config.Id.String(){{end}}))

	{{- if hasDataSourceQuery .Attributes }}

	// Object found by a query attribute, it is not retrieved again by ID
	var res gjson.Result
	{{- range $index, $dataSourceAttribute := $dataSourceAttributes }}
	if config.Id.IsNull() && !config.{{toGoName $dataSourceAttribute.TfName}}.IsNull() {
		// Client handles pagination and returns all items merged into a single list
		items, err := d.client.GetAll(ctx, config.getPath(), reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, got error: %s", err))
			return
		}
		if !items.IsArray() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, expected a list, got: %s", items.Type))
			return
		}
		items.ForEach(func(k, v gjson.Result) bool {
			if config.{{toGoName $dataSourceAttribute.TfName}}.
				{{- if eq $dataSourceAttribute.Type "Int64" -}}ValueInt64()
				{{- else -}}ValueString(){{- end -}} == v.Get("{{range $dataSourceAttribute.DataPath}}{{.}}.{{end}}{{$dataSourceAttribute.ModelName}}").
				{{- if eq $dataSourceAttribute.Type "Int64" -}}Int()
				{{- else -}}String(){{- end -}} {
				config.Id = types.StringValue(v.Get("id").String())
				res = v
				tflog.Debug(ctx, fmt.Sprintf("%s: Found object with {{$dataSourceAttribute.TfName}} '%v', id: %v", config.Id.ValueString(), config.{{toGoName $dataSourceAttribute.TfName}}.{{if eq $dataSourceAttribute.Type "Int64"}}ValueInt64(){{else}}ValueString(){{end}}, config.Id.ValueString()))
				return false
			}
			return true
		})

		if config.Id.IsNull() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to find object with {{$dataSourceAttribute.TfName}}: %v", config.{{toGoName $dataSourceAttribute.TfName}}.{{if eq $dataSourceAttribute.Type "Int64"}}ValueInt64(){{else}}ValueString(){{end}}))
			return
		}
	}
	{{- end}}
	{{- end}}

	{{if hasDataSourceQuery .Attributes}}
	if !res.Exists() {
		urlPath := config.getPath()+"/"+url.QueryEscape(config.Id.ValueString())
		var err error
		res, err = d.client.Get(ctx, urlPath, reqMods...)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
			return
		}
	}
	{{- else}}
	{{- if or .NoId .SettingsList}}
	urlPath := config.getPath()
	{{- else}}
	urlPath := config.getPath()+"/"+url.QueryEscape(config.Id.ValueString())
	{{- end}}
	{{- if .IsBulk}}
	// Client handles pagination and returns all items merged into a single list
	res, err := d.client.GetAll(ctx, urlPath, reqMods...)
	{{- else}}
	res, err := d.client.Get(ctx, urlPath, reqMods...)
	{{- end}}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}
	{{- if .IsBulk}}
	if !res.IsArray() {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, expected a list, got: %s", res.Type))
		return
	}
	{{- end}}
	{{- end}}

	{{- if .IsBulk}}
	{{- $items := getAttributeByTfName .Attributes "items"}}
	{{- $filter := getAttributeByTfName .Attributes "filter"}}

	// Read all items if user did not provide any specific item names in the config
	if len(config.{{toGoName $items.TfName}}) == 0 {
		config.{{toGoName $items.TfName}} = map[string]{{$items.GoTypeName}}{}
		res.{{if $items.ModelName}}Get("{{range $items.DataPath}}{{.}}.{{end}}{{$items.ModelName}}").{{end}}ForEach(func(_, v gjson.Result) bool {
			if name := v.Get("name").String(); name != "" {
				config.{{toGoName $items.TfName}}[name] = {{$items.GoTypeName}}{}
			}
			return true
		})
	}
	{{- if $filter.TfName}}

	// Keep only items whose name contains the filter (case-insensitive)
	if filter := strings.ToLower(config.{{toGoName $filter.TfName}}.ValueString()); filter != "" {
		for name := range config.{{toGoName $items.TfName}} {
			if !strings.Contains(strings.ToLower(name), filter) {
				tflog.Debug(ctx, fmt.Sprintf("item does not match filter, removing: name=%v", name))
				delete(config.{{toGoName $items.TfName}}, name)
			}
		}
	}
	{{- end}}
	{{- end}}

	config.fromBody(ctx, res)
	{{- if .SettingsList}}
	config.Id = types.StringValue("{{snakeCase .Name}}")
	{{- end}}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", {{if .NoId}}"{{.Name}}"{{else}}config.Id.ValueString(){{end}}))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read
