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

type PrivateResource struct {
	Id                        types.String                               `tfsdk:"id"`
	Name                      types.String                               `tfsdk:"name"`
	Description               types.String                               `tfsdk:"description"`
	CertificateId             types.Int64                                `tfsdk:"certificate_id"`
	DnsServerId               types.Int64                                `tfsdk:"dns_server_id"`
	EndpointConnectionMethods []PrivateResourceEndpointConnectionMethods `tfsdk:"endpoint_connection_methods"`
	Addresses                 []PrivateResourceAddresses                 `tfsdk:"addresses"`
}

type PrivateResourceEndpointConnectionMethods struct {
	Type                      types.String `tfsdk:"type"`
	ClientReachableAddresses  types.List   `tfsdk:"client_reachable_addresses"`
	BrowserExternalFqdnPrefix types.String `tfsdk:"browser_external_fqdn_prefix"`
	BrowserProtocol           types.String `tfsdk:"browser_protocol"`
	BrowserSni                types.String `tfsdk:"browser_sni"`
	BrowserSslVerification    types.Bool   `tfsdk:"browser_ssl_verification"`
}

type PrivateResourceAddresses struct {
	Addresses         types.List                                  `tfsdk:"addresses"`
	ProtocolsAndPorts []PrivateResourceAddressesProtocolsAndPorts `tfsdk:"protocols_and_ports"`
}

type PrivateResourceAddressesProtocolsAndPorts struct {
	Protocol types.String `tfsdk:"protocol"`
	Ports    types.String `tfsdk:"ports"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data PrivateResource) getPath() string {
	return "/policies/v2/privateResources"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data PrivateResource) toBody(ctx context.Context, state PrivateResource) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "resourceId", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	if !data.Description.IsNull() {
		body, _ = sjson.Set(body, "description", data.Description.ValueString())
	}
	if !data.CertificateId.IsNull() {
		body, _ = sjson.Set(body, "certificateId", data.CertificateId.ValueInt64())
	}
	if !data.DnsServerId.IsNull() {
		body, _ = sjson.Set(body, "dnsServerId", data.DnsServerId.ValueInt64())
	}
	if len(data.EndpointConnectionMethods) > 0 {
		var endpointConnectionMethodsBody strings.Builder
		endpointConnectionMethodsBody.WriteString("[")
		for _, item := range data.EndpointConnectionMethods {
			itemBody := ""
			if !item.Type.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "type", item.Type.ValueString())
			}
			if !item.ClientReachableAddresses.IsNull() {
				var values []string
				item.ClientReachableAddresses.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "reachableAddresses", values)
			}
			if !item.BrowserExternalFqdnPrefix.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "externalFQDNPrefix", item.BrowserExternalFqdnPrefix.ValueString())
			}
			if !item.BrowserProtocol.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "protocol", item.BrowserProtocol.ValueString())
			}
			if !item.BrowserSni.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "sni", item.BrowserSni.ValueString())
			}
			if !item.BrowserSslVerification.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "sslVerificationEnabled", item.BrowserSslVerification.ValueBool())
			}
			if itemBody != "" {
				if endpointConnectionMethodsBody.Len() > 1 {
					endpointConnectionMethodsBody.WriteString(",")
				}
				endpointConnectionMethodsBody.WriteString(itemBody)
			}
		}
		endpointConnectionMethodsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "accessTypes", endpointConnectionMethodsBody.String())
	}
	if len(data.Addresses) > 0 {
		var addressesBody strings.Builder
		addressesBody.WriteString("[")
		for _, item := range data.Addresses {
			itemBody := ""
			if !item.Addresses.IsNull() {
				var values []string
				item.Addresses.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "destinationAddr", values)
			}
			if len(item.ProtocolsAndPorts) > 0 {
				var protocolsAndPortsChildBody strings.Builder
				protocolsAndPortsChildBody.WriteString("[")
				for _, childItem := range item.ProtocolsAndPorts {
					itemChildBody := ""
					if !childItem.Protocol.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "protocol", childItem.Protocol.ValueString())
					}
					if !childItem.Ports.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ports", childItem.Ports.ValueString())
					}
					if itemChildBody != "" {
						if protocolsAndPortsChildBody.Len() > 1 {
							protocolsAndPortsChildBody.WriteString(",")
						}
						protocolsAndPortsChildBody.WriteString(itemChildBody)
					}
				}
				protocolsAndPortsChildBody.WriteString("]")
				itemBody, _ = sjson.SetRaw(itemBody, "protocolPorts", protocolsAndPortsChildBody.String())
			}
			if itemBody != "" {
				if addressesBody.Len() > 1 {
					addressesBody.WriteString(",")
				}
				addressesBody.WriteString(itemBody)
			}
		}
		addressesBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "resourceAddresses", addressesBody.String())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *PrivateResource) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("description"); value.Exists() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("certificateId"); value.Exists() {
		data.CertificateId = types.Int64Value(value.Int())
	} else {
		data.CertificateId = types.Int64Null()
	}
	if value := res.Get("dnsServerId"); value.Exists() {
		data.DnsServerId = types.Int64Value(value.Int())
	} else {
		data.DnsServerId = types.Int64Null()
	}
	if value := res.Get("accessTypes"); value.Exists() {
		data.EndpointConnectionMethods = make([]PrivateResourceEndpointConnectionMethods, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := PrivateResourceEndpointConnectionMethods{}
			if value := res.Get("type"); value.Exists() {
				data.Type = types.StringValue(value.String())
			} else {
				data.Type = types.StringNull()
			}
			if value := res.Get("reachableAddresses"); value.Exists() {
				data.ClientReachableAddresses = helpers.GetStringList(value.Array())
			} else {
				data.ClientReachableAddresses = types.ListNull(types.StringType)
			}
			if value := res.Get("externalFQDNPrefix"); value.Exists() {
				data.BrowserExternalFqdnPrefix = types.StringValue(value.String())
			} else {
				data.BrowserExternalFqdnPrefix = types.StringNull()
			}
			if value := res.Get("protocol"); value.Exists() {
				data.BrowserProtocol = types.StringValue(value.String())
			} else {
				data.BrowserProtocol = types.StringNull()
			}
			if value := res.Get("sni"); value.Exists() {
				data.BrowserSni = types.StringValue(value.String())
			} else {
				data.BrowserSni = types.StringNull()
			}
			if value := res.Get("sslVerificationEnabled"); value.Exists() {
				data.BrowserSslVerification = types.BoolValue(value.Bool())
			} else {
				data.BrowserSslVerification = types.BoolNull()
			}
			(*parent).EndpointConnectionMethods = append((*parent).EndpointConnectionMethods, data)
			return true
		})
	}
	if value := res.Get("resourceAddresses"); value.Exists() {
		data.Addresses = make([]PrivateResourceAddresses, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := PrivateResourceAddresses{}
			if value := res.Get("destinationAddr"); value.Exists() {
				data.Addresses = helpers.GetStringList(value.Array())
			} else {
				data.Addresses = types.ListNull(types.StringType)
			}
			if value := res.Get("protocolPorts"); value.Exists() {
				data.ProtocolsAndPorts = make([]PrivateResourceAddressesProtocolsAndPorts, 0, int(value.Get("#").Int()))
				value.ForEach(func(k, res gjson.Result) bool {
					parent := &data
					data := PrivateResourceAddressesProtocolsAndPorts{}
					if value := res.Get("protocol"); value.Exists() {
						data.Protocol = types.StringValue(value.String())
					} else {
						data.Protocol = types.StringNull()
					}
					if value := res.Get("ports"); value.Exists() {
						data.Ports = types.StringValue(value.String())
					} else {
						data.Ports = types.StringNull()
					}
					(*parent).ProtocolsAndPorts = append((*parent).ProtocolsAndPorts, data)
					return true
				})
			}
			(*parent).Addresses = append((*parent).Addresses, data)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *PrivateResource) fromBodyPartial(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("description"); value.Exists() && !data.Description.IsNull() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("certificateId"); value.Exists() && !data.CertificateId.IsNull() {
		data.CertificateId = types.Int64Value(value.Int())
	} else {
		data.CertificateId = types.Int64Null()
	}
	if value := res.Get("dnsServerId"); value.Exists() && !data.DnsServerId.IsNull() {
		data.DnsServerId = types.Int64Value(value.Int())
	} else {
		data.DnsServerId = types.Int64Null()
	}
	endpointConnectionMethodsArray := res.Get("accessTypes")
	for i := 0; i < len(data.EndpointConnectionMethods); i++ {
		keys := [...]string{"type"}
		keyValues := [...]string{data.EndpointConnectionMethods[i].Type.ValueString()}

		parent := &data
		data := (*parent).EndpointConnectionMethods[i]
		var res gjson.Result

		endpointConnectionMethodsArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing EndpointConnectionMethods[%d] = %+v",
				i,
				(*parent).EndpointConnectionMethods[i],
			))
			(*parent).EndpointConnectionMethods = slices.Delete((*parent).EndpointConnectionMethods, i, i+1)
			i--

			continue
		}
		if value := res.Get("type"); value.Exists() && !data.Type.IsNull() {
			data.Type = types.StringValue(value.String())
		} else {
			data.Type = types.StringNull()
		}
		if value := res.Get("reachableAddresses"); value.Exists() && !data.ClientReachableAddresses.IsNull() {
			data.ClientReachableAddresses = helpers.GetStringList(value.Array())
		} else {
			data.ClientReachableAddresses = types.ListNull(types.StringType)
		}
		if value := res.Get("externalFQDNPrefix"); value.Exists() && !data.BrowserExternalFqdnPrefix.IsNull() {
			data.BrowserExternalFqdnPrefix = types.StringValue(value.String())
		} else {
			data.BrowserExternalFqdnPrefix = types.StringNull()
		}
		if value := res.Get("protocol"); value.Exists() && !data.BrowserProtocol.IsNull() {
			data.BrowserProtocol = types.StringValue(value.String())
		} else {
			data.BrowserProtocol = types.StringNull()
		}
		if value := res.Get("sni"); value.Exists() && !data.BrowserSni.IsNull() {
			data.BrowserSni = types.StringValue(value.String())
		} else {
			data.BrowserSni = types.StringNull()
		}
		if value := res.Get("sslVerificationEnabled"); value.Exists() && !data.BrowserSslVerification.IsNull() {
			data.BrowserSslVerification = types.BoolValue(value.Bool())
		} else {
			data.BrowserSslVerification = types.BoolNull()
		}
		(*parent).EndpointConnectionMethods[i] = data
	}
	addressesArray := res.Get("resourceAddresses")
	for i := 0; i < len(data.Addresses); i++ {
		keys := [...]string{"destinationAddr"}
		keyValues := [...]string{helpers.ListSetKeyValue(data.Addresses[i].Addresses)}

		parent := &data
		data := (*parent).Addresses[i]
		var res gjson.Result

		addressesArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing Addresses[%d] = %+v",
				i,
				(*parent).Addresses[i],
			))
			(*parent).Addresses = slices.Delete((*parent).Addresses, i, i+1)
			i--

			continue
		}
		if value := res.Get("destinationAddr"); value.Exists() && !data.Addresses.IsNull() {
			data.Addresses = helpers.GetStringList(value.Array())
		} else {
			data.Addresses = types.ListNull(types.StringType)
		}
		protocolsAndPortsArray := res.Get("protocolPorts")
		for i := 0; i < len(data.ProtocolsAndPorts); i++ {
			keys := [...]string{"protocol", "ports"}
			keyValues := [...]string{data.ProtocolsAndPorts[i].Protocol.ValueString(), data.ProtocolsAndPorts[i].Ports.ValueString()}

			parent := &data
			data := (*parent).ProtocolsAndPorts[i]
			var res gjson.Result

			protocolsAndPortsArray.ForEach(
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
				tflog.Debug(ctx, fmt.Sprintf("removing ProtocolsAndPorts[%d] = %+v",
					i,
					(*parent).ProtocolsAndPorts[i],
				))
				(*parent).ProtocolsAndPorts = slices.Delete((*parent).ProtocolsAndPorts, i, i+1)
				i--

				continue
			}
			if value := res.Get("protocol"); value.Exists() && !data.Protocol.IsNull() {
				data.Protocol = types.StringValue(value.String())
			} else {
				data.Protocol = types.StringNull()
			}
			if value := res.Get("ports"); value.Exists() && !data.Ports.IsNull() {
				data.Ports = types.StringValue(value.String())
			} else {
				data.Ports = types.StringNull()
			}
			(*parent).ProtocolsAndPorts[i] = data
		}
		(*parent).Addresses[i] = data
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *PrivateResource) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

// Section below is generated&owned by "gen/generator.go". //template:begin adjustBody

// End of section. //template:end adjustBody
