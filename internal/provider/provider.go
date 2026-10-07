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

// Section below is generated&owned by "gen/generator.go". //template:begin provider
import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-secureaccess"
)

// SecureAccessProvider defines the provider implementation.
type SecureAccessProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SecureAccessProviderModel describes the provider data model.
type SecureAccessProviderModel struct {
	ApiKey       types.String `tfsdk:"api_key"`
	ApiKeySecret types.String `tfsdk:"api_key_secret"`
	URL          types.String `tfsdk:"url"`
	Insecure     types.Bool   `tfsdk:"insecure"`
	ReqTimeout   types.String `tfsdk:"req_timeout"`
	Retries      types.Int64  `tfsdk:"retries"`
}

// SecureAccessProviderData describes the data maintained by the provider.
type SecureAccessProviderData struct {
	Client *secureaccess.Client
}

// Metadata returns the provider type name.
func (p *SecureAccessProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "secureaccess"
	resp.Version = p.version
}

func (p *SecureAccessProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for the SecureAccess instance. This can also be set as the SECUREACCESS_API_KEY environment variable.",
				Optional:            true,
			},
			"api_key_secret": schema.StringAttribute{
				MarkdownDescription: "API Key Secret for the SecureAccess instance. This can also be set as the SECUREACCESS_API_KEY_SECRET environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of the Cisco Secure Access instance (https://api.sse.cisco.com). This can also be set as the SECUREACCESS_URL environment variable.",
				Optional:            true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Allow insecure HTTPS client. This can also be set as the SECUREACCESS_INSECURE environment variable. Defaults to `false`.",
				Optional:            true,
			},
			"req_timeout": schema.StringAttribute{
				MarkdownDescription: "Timeout for a single HTTPS request made to REST API before it is retried. This can also be set as the SECUREACCESS_REQTIMEOUT environment variable. A string like `\"1s\"` means one second. Defaults to unlimited.",
				Optional:            true,
			},
			"retries": schema.Int64Attribute{
				MarkdownDescription: "Number of retries for REST API calls. This can also be set as the SECUREACCESS_RETRIES environment variable. Defaults to `3`.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 9),
				},
			},
		},
	}
}

func (p *SecureAccessProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config SecureAccessProviderModel
	var err error

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get API Key
	var apiKey string
	if config.ApiKey.IsNull() || config.ApiKey.IsUnknown() {
		apiKey = os.Getenv("SECUREACCESS_API_KEY")
	} else {
		apiKey = config.ApiKey.ValueString()
	}

	// Get API Key Secret
	var apiKeySecret string
	if config.ApiKeySecret.IsNull() || config.ApiKeySecret.IsUnknown() {
		apiKeySecret = os.Getenv("SECUREACCESS_API_KEY_SECRET")
	} else {
		apiKeySecret = config.ApiKeySecret.ValueString()
	}

	// Fail if the user didn't provider any credentials
	if apiKey == "" || apiKeySecret == "" {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Please provide credentials for Secure Access (API key and API key secret)",
		)
		return
	}

	var url string
	if config.URL.IsNull() {
		url = os.Getenv("SECUREACCESS_URL")
	} else {
		url = config.URL.ValueString()
	}

	if url == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find url",
			"URL cannot be an empty string",
		)
		return
	}

	var insecure bool
	if config.Insecure.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as insecure",
		)
		return
	}

	if config.Insecure.IsNull() {
		insecureStr := os.Getenv("SECUREACCESS_INSECURE")
		if insecureStr == "" {
			insecure = false
		} else {
			insecure, _ = strconv.ParseBool(insecureStr)
		}
	} else {
		insecure = config.Insecure.ValueBool()
	}

	var reqTimeout time.Duration
	if config.ReqTimeout.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as req_timeout",
		)
		return
	}

	var reqTimeoutStr string
	if config.ReqTimeout.IsNull() {
		reqTimeoutStr = os.Getenv("SECUREACCESS_REQTIMEOUT")
		if reqTimeoutStr == "" {
			reqTimeoutStr = "0s"
		}
	} else {
		reqTimeoutStr = config.ReqTimeout.ValueString()
	}
	reqTimeout, err = time.ParseDuration(reqTimeoutStr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			fmt.Sprintf("Cannot parse the req_timeout string: %v", err),
		)
		return
	}

	var retries int64
	if config.Retries.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as retries",
		)
		return
	}

	if config.Retries.IsNull() {
		retriesStr := os.Getenv("SECUREACCESS_RETRIES")
		if retriesStr == "" {
			retries = 3
		} else {
			retries, _ = strconv.ParseInt(retriesStr, 0, 64)
		}
	} else {
		retries = config.Retries.ValueInt64()
	}

	tflog.Debug(ctx, fmt.Sprint("Creating a new Secure Access client",
		"  url=", url,
		"  insecure=", insecure,
		"  req_timeout=", reqTimeout,
		"  retries=", retries,
	))

	// Create a new Secure Access client and set it to the provider client
	var c secureaccess.Client

	c, err = secureaccess.NewClient(url, apiKey, apiKeySecret, secureaccess.Insecure(insecure), secureaccess.MaxRetries(int(retries)), secureaccess.RequestTimeout(reqTimeout))

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create Secure Access client:\n\n"+err.Error(),
		)
		return
	}

	data := SecureAccessProviderData{Client: &c}
	resp.DataSourceData = &data
	resp.ResourceData = &data
}

func (p *SecureAccessProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDestinationListResource,
		NewInternalDomainResource,
		NewInternalNetworkResource,
		NewNetworkResource,
		NewNetworkObjectResource,
		NewNetworkObjectGroupResource,
		NewSiteResource,
	}
}

func (p *SecureAccessProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDestinationListDataSource,
		NewInternalDomainDataSource,
		NewInternalNetworkDataSource,
		NewNetworkDataSource,
		NewNetworkObjectDataSource,
		NewNetworkObjectGroupDataSource,
		NewSiteDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SecureAccessProvider{
			version: version,
		}
	}
}

// End of section. //template:end provider
