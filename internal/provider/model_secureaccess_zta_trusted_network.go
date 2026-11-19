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
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type ZTATrustedNetwork struct {
	Id             types.String                      `tfsdk:"id"`
	Name           types.String                      `tfsdk:"name"`
	Inspect        types.Int64                       `tfsdk:"inspect"`
	DnsServers     []ZTATrustedNetworkDnsServers     `tfsdk:"dns_servers"`
	DnsDomains     []ZTATrustedNetworkDnsDomains     `tfsdk:"dns_domains"`
	TrustedServers []ZTATrustedNetworkTrustedServers `tfsdk:"trusted_servers"`
	Revision       types.Int64                       `tfsdk:"revision"`
}

type ZTATrustedNetworkDnsServers struct {
	Ip types.String `tfsdk:"ip"`
}

type ZTATrustedNetworkDnsDomains struct {
	Name types.String `tfsdk:"name"`
}

type ZTATrustedNetworkTrustedServers struct {
	Url             types.String `tfsdk:"url"`
	CertificateHash types.String `tfsdk:"certificate_hash"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data ZTATrustedNetwork) getPath() string {
	return "/deployments/v2/ztna/trustedNetworks"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data ZTATrustedNetwork) toBody(ctx context.Context, state ZTATrustedNetwork) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "networkId", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "networkName", data.Name.ValueString())
	}
	if !data.Inspect.IsNull() {
		body, _ = sjson.Set(body, "interfaceType", data.Inspect.ValueInt64())
	}
	if len(data.DnsServers) > 0 {
		var dnsServersBody strings.Builder
		dnsServersBody.WriteString("[")
		for _, item := range data.DnsServers {
			itemBody := ""
			if !item.Ip.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "serverIp", item.Ip.ValueString())
			}
			if itemBody != "" {
				if dnsServersBody.Len() > 1 {
					dnsServersBody.WriteString(",")
				}
				dnsServersBody.WriteString(itemBody)
			}
		}
		dnsServersBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "criteria.dnsServers", dnsServersBody.String())
	}
	if len(data.DnsDomains) > 0 {
		var dnsDomainsBody strings.Builder
		dnsDomainsBody.WriteString("[")
		for _, item := range data.DnsDomains {
			itemBody := ""
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if itemBody != "" {
				if dnsDomainsBody.Len() > 1 {
					dnsDomainsBody.WriteString(",")
				}
				dnsDomainsBody.WriteString(itemBody)
			}
		}
		dnsDomainsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "criteria.dnsDomains", dnsDomainsBody.String())
	}
	if len(data.TrustedServers) > 0 {
		var trustedServersBody strings.Builder
		trustedServersBody.WriteString("[")
		for _, item := range data.TrustedServers {
			itemBody := ""
			if !item.Url.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "url", item.Url.ValueString())
			}
			if !item.CertificateHash.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "certificateHash", item.CertificateHash.ValueString())
			}
			if itemBody != "" {
				if trustedServersBody.Len() > 1 {
					trustedServersBody.WriteString(",")
				}
				trustedServersBody.WriteString(itemBody)
			}
		}
		trustedServersBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "criteria.trustedServers", trustedServersBody.String())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *ZTATrustedNetwork) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("networkName"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("interfaceType"); value.Exists() {
		data.Inspect = types.Int64Value(value.Int())
	} else {
		data.Inspect = types.Int64Value(0)
	}
	if value := res.Get("criteria.dnsServers"); value.Exists() {
		data.DnsServers = make([]ZTATrustedNetworkDnsServers, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTATrustedNetworkDnsServers{}
			if value := res.Get("serverIp"); value.Exists() {
				data.Ip = types.StringValue(value.String())
			} else {
				data.Ip = types.StringNull()
			}
			(*parent).DnsServers = append((*parent).DnsServers, data)
			return true
		})
	}
	if value := res.Get("criteria.dnsDomains"); value.Exists() {
		data.DnsDomains = make([]ZTATrustedNetworkDnsDomains, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTATrustedNetworkDnsDomains{}
			if value := res.Get("name"); value.Exists() {
				data.Name = types.StringValue(value.String())
			} else {
				data.Name = types.StringNull()
			}
			(*parent).DnsDomains = append((*parent).DnsDomains, data)
			return true
		})
	}
	if value := res.Get("criteria.trustedServers"); value.Exists() {
		data.TrustedServers = make([]ZTATrustedNetworkTrustedServers, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := ZTATrustedNetworkTrustedServers{}
			if value := res.Get("url"); value.Exists() {
				data.Url = types.StringValue(value.String())
			} else {
				data.Url = types.StringNull()
			}
			if value := res.Get("certificateHash"); value.Exists() {
				data.CertificateHash = types.StringValue(value.String())
			} else {
				data.CertificateHash = types.StringNull()
			}
			(*parent).TrustedServers = append((*parent).TrustedServers, data)
			return true
		})
	}
	if value := res.Get("rev"); value.Exists() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *ZTATrustedNetwork) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("networkName"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("interfaceType"); value.Exists() && !data.Inspect.IsNull() {
		data.Inspect = types.Int64Value(value.Int())
	} else if data.Inspect.ValueInt64() != 0 {
		data.Inspect = types.Int64Null()
	}
	dnsServersArray := res.Get("criteria.dnsServers")
	for i := 0; i < len(data.DnsServers); i++ {
		keys := [...]string{"serverIp"}
		keyValues := [...]string{data.DnsServers[i].Ip.ValueString()}

		parent := &data
		data := (*parent).DnsServers[i]
		var res gjson.Result

		dnsServersArray.ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if helpers.GetKeyValue(v.Get(keys[ik])) != keyValues[ik] {
						found = false
						break
					}
					found = true
				}
				if found {
					res = v
					return false
				}
				return true
			},
		)
		if !res.Exists() {
			tflog.Debug(ctx, fmt.Sprintf("removing DnsServers[%d] = %+v",
				i,
				(*parent).DnsServers[i],
			))
			(*parent).DnsServers = slices.Delete((*parent).DnsServers, i, i+1)
			i--

			continue
		}
		if value := res.Get("serverIp"); value.Exists() && !data.Ip.IsNull() {
			data.Ip = types.StringValue(value.String())
		} else {
			data.Ip = types.StringNull()
		}
		(*parent).DnsServers[i] = data
	}
	dnsDomainsArray := res.Get("criteria.dnsDomains")
	for i := 0; i < len(data.DnsDomains); i++ {
		keys := [...]string{"name"}
		keyValues := [...]string{data.DnsDomains[i].Name.ValueString()}

		parent := &data
		data := (*parent).DnsDomains[i]
		var res gjson.Result

		dnsDomainsArray.ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if helpers.GetKeyValue(v.Get(keys[ik])) != keyValues[ik] {
						found = false
						break
					}
					found = true
				}
				if found {
					res = v
					return false
				}
				return true
			},
		)
		if !res.Exists() {
			tflog.Debug(ctx, fmt.Sprintf("removing DnsDomains[%d] = %+v",
				i,
				(*parent).DnsDomains[i],
			))
			(*parent).DnsDomains = slices.Delete((*parent).DnsDomains, i, i+1)
			i--

			continue
		}
		if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
			data.Name = types.StringValue(value.String())
		} else {
			data.Name = types.StringNull()
		}
		(*parent).DnsDomains[i] = data
	}
	trustedServersArray := res.Get("criteria.trustedServers")
	for i := 0; i < len(data.TrustedServers); i++ {
		keys := [...]string{"url"}
		keyValues := [...]string{data.TrustedServers[i].Url.ValueString()}

		parent := &data
		data := (*parent).TrustedServers[i]
		var res gjson.Result

		trustedServersArray.ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if helpers.GetKeyValue(v.Get(keys[ik])) != keyValues[ik] {
						found = false
						break
					}
					found = true
				}
				if found {
					res = v
					return false
				}
				return true
			},
		)
		if !res.Exists() {
			tflog.Debug(ctx, fmt.Sprintf("removing TrustedServers[%d] = %+v",
				i,
				(*parent).TrustedServers[i],
			))
			(*parent).TrustedServers = slices.Delete((*parent).TrustedServers, i, i+1)
			i--

			continue
		}
		if value := res.Get("url"); value.Exists() && !data.Url.IsNull() {
			data.Url = types.StringValue(value.String())
		} else {
			data.Url = types.StringNull()
		}
		if value := res.Get("certificateHash"); value.Exists() && !data.CertificateHash.IsNull() {
			data.CertificateHash = types.StringValue(value.String())
		} else {
			data.CertificateHash = types.StringNull()
		}
		(*parent).TrustedServers[i] = data
	}
	if value := res.Get("rev"); value.Exists() && !data.Revision.IsNull() {
		data.Revision = types.Int64Value(value.Int())
	} else {
		data.Revision = types.Int64Null()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *ZTATrustedNetwork) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
	if data.Revision.IsUnknown() {
		if value := res.Get("rev"); value.Exists() {
			data.Revision = types.Int64Value(value.Int())
		} else {
			data.Revision = types.Int64Null()
		}
	}
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

func (data ZTATrustedNetwork) adjustBody(ctx context.Context, req string) string {

	req, _ = sjson.Delete(req, "networkId")
	// Revision is kept from state in the plan (unknown on create), the API requires the current one on update
	if !data.Revision.IsUnknown() && !data.Revision.IsNull() {
		req, _ = sjson.Set(req, "rev", data.Revision.ValueInt64())
	}

	return req
}

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
