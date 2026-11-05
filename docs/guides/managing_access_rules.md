---
subcategory: "Guides"
page_title: "Managing Access Rules"
description: |-
    How to configure the sources and destinations of Access Rules, and how to combine them with source_operator and destination_operator.
---

# Managing Access Rules

An Access Rule (`secureaccess_access_rule`) decides what happens (`action`) to traffic coming **from** some sources
and going **to** some destinations. This guide explains how the source and destination attributes work together,
with a focus on `source_operator`, `destination_operator` and `destination_operator_network_service`.

## The short version

- Every source attribute (`source_*`) and destination attribute (`destination_*`) is a list of IDs or values.
- Within one attribute, the values are always alternatives: `destination_network_object_ids = [A, B]` means *A or B*.
- Between attributes, the default is also **OR**: the rule matches if the traffic matches *any* of the configured sources
  and *any* of the configured destinations.
- `source_operator = "AND"` and `destination_operator = "AND"` make the traffic match **every** configured type instead.
- A side without any attribute matches **everything** (any source, or any destination).

The source side and the destination side are always combined with AND: traffic must match the sources **and** the destinations.

## Rule types: `traffic_type`

`traffic_type` is mandatory and selects the kind of rule. It decides which sources, destinations and actions are available.

| | `PRIVATE_NETWORK` (private access) | `PUBLIC_INTERNET` (internet access) |
|---|---|---|
| Actions | `allow`, `block` | `allow`, `block`, `warn`, `isolate` |
| `source_identity_ids`, `source_identity_type_ids` | yes | yes |
| `source_geolocations` | yes | yes |
| `source_network_object_ids`, `source_network_object_group_ids` | yes | yes |
| `source_private_literals` (IPs) | yes | - |
| `source_internet_literals` (IPs with ports) | - | yes |
| `destination_network_object_ids`, `destination_network_object_group_ids` | yes | yes |
| `destination_service_object_ids`, `destination_service_object_group_ids` | yes | yes |
| `destination_literals` (IPs with ports and protocol) | yes | yes |
| `destination_private_resource_ids`, `destination_private_resource_group_ids` | yes | - |
| `destination_security_group_tag_ids` | yes | - |
| `destination_content_category_ids`, `destination_content_category_list_ids` | - | yes |
| `destination_application_ids`, `destination_application_category_ids`, `destination_application_list_ids` | - | yes |
| `destination_list_ids` | - | yes |
| `destination_geolocations` | - | yes, `block` only |

`traffic_type` can be changed on an existing rule, the provider updates it in place.

The provider checks these combinations when planning, so an unsupported combination fails at `terraform plan`, not at `apply`.

## Sources and `source_operator`

### `source_operator = "OR"` (default)

The rule matches traffic from **any** of the configured sources.

```terraform
resource "secureaccess_access_rule" "block_contractors_or_lab" {
  name         = "Block contractors or lab network"
  action       = "block"
  log_level    = "LOG_ALL"
  traffic_type = "PUBLIC_INTERNET"

  source_identity_ids       = [data.secureaccess_identities.contractors.items["directory_group"]["Contractors"].id]
  source_network_object_ids = [secureaccess_network_object.lab.id]
}
```

Effect: traffic from a member of the *Contractors* group, **or** traffic from the lab network, to any destination, is blocked.

### `source_operator = "AND"`

The rule matches traffic only if it matches **every** configured source attribute. Use it to narrow the sources down,
for example to *a group of users, but only from certain countries*.

```terraform
resource "secureaccess_access_rule" "finance_from_poland" {
  name         = "Finance users from Poland"
  action       = "allow"
  log_level    = "LOG_ALL"
  traffic_type = "PUBLIC_INTERNET"

  source_identity_ids = [data.secureaccess_identities.finance.items["directory_group"]["Finance"].id]
  source_geolocations = ["PL"]
  source_operator     = "AND"
}
```

Effect: traffic from a member of the *Finance* group **and** located in Poland. Finance users in other countries,
and other users in Poland, do not match.

Rules for `source_operator = "AND"`:

- At least **two** of these attributes must be set: `source_identity_ids`, `source_identity_type_ids`,
  `source_geolocations`, `source_internet_literals`.
- `source_network_object_ids`, `source_network_object_group_ids` and `source_private_literals` cannot be used with `AND`.

### Geolocations need `AND` when combined with identities

`source_geolocations` alone works with `OR` (traffic from these countries). Combined with `source_identity_ids` or
`source_identity_type_ids`, it requires `source_operator = "AND"`, i.e. "these identities from these countries".
Geolocations can only be combined with Active Directory users or groups.

## Destinations and `destination_operator`

### `destination_operator = "OR"` (default)

The rule matches traffic to **any** of the configured destinations. Any destination attribute can be used.

```terraform
resource "secureaccess_access_rule" "block_social_and_gambling" {
  name         = "Block social networks and gambling"
  action       = "block"
  log_level    = "LOG_ALL"
  traffic_type = "PUBLIC_INTERNET"

  destination_content_category_ids = [
    data.secureaccess_content_categories.example.items["Social Networking"].id,
    data.secureaccess_content_categories.example.items["Gambling"].id,
  ]
  destination_list_ids = [secureaccess_destination_list.blocked.id]
}
```

Effect: everybody (no source is configured) is blocked from social networks, **or** gambling sites,
**or** any destination of the destination list.

### `destination_operator = "AND"`

With `AND`, the destination attributes are grouped by type. Traffic must match **every group**, and **any value within a group**:

| Group | Attributes |
|---|---|
| Network objects | `destination_network_object_ids`, `destination_network_object_group_ids` |
| Service objects | `destination_service_object_ids`, `destination_service_object_group_ids` |
| Content categories | `destination_content_category_ids`, `destination_content_category_list_ids` |
| Applications | `destination_application_ids`, `destination_application_category_ids`, `destination_application_list_ids` |

Use it to narrow the destinations down, for example *these servers, but only on these ports*.

```terraform
resource "secureaccess_access_rule" "web_servers_https" {
  name         = "Web servers over HTTPS"
  action       = "allow"
  log_level    = "LOG_ALL"
  traffic_type = "PRIVATE_NETWORK"

  destination_network_object_ids       = [secureaccess_network_object.web_servers.id]
  destination_network_object_group_ids = [secureaccess_network_object_group.dmz.id]
  destination_service_object_ids       = [secureaccess_service_object.https.id]
  destination_operator                 = "AND"
}
```

Effect: traffic to the web servers **or** the DMZ (network objects group) **and** on HTTPS (service objects group).
Traffic to the web servers on another port does not match, nor does HTTPS traffic to other servers.

Rules for `destination_operator = "AND"`:

- At least **two groups** must be configured. Two attributes of the same group (e.g. network objects and network
  object groups only) are not enough, as they are ORed within the group anyway.
- `destination_list_ids`, `destination_private_resource_ids`, `destination_private_resource_group_ids`,
  `destination_security_group_tag_ids`, `destination_literals` and `destination_geolocations` cannot be used with `AND`.

More examples of `AND`, for internet rules:

```terraform
# Social networks (category) AND only on the corporate proxy port (service)
destination_content_category_ids = [data.secureaccess_content_categories.example.items["Social Networking"].id]
destination_service_object_ids   = [secureaccess_service_object.proxy_port.id]
destination_operator             = "AND"

# File sharing applications (application group) AND only to the partner networks (network group)
destination_application_category_ids = [tonumber(data.secureaccess_application_categories.example.items["File Sharing"].id)]
destination_network_object_ids       = [secureaccess_network_object.partners.id]
destination_operator                 = "AND"
```

### `destination_operator_network_service`

Only relevant with `destination_operator = "AND"`. It decides how the **network objects** and **service objects**
groups are combined with each other:

- `"AND"` (default): network objects and service objects are separate groups, traffic must match both, as in the
  *web servers over HTTPS* example above.
- `"OR"`: network objects and service objects are merged into **one** group, traffic must match either of them.
  Because a group alone is not enough for `AND`, at least one other group (content categories or applications) must be configured.

```terraform
resource "secureaccess_access_rule" "gambling_on_partner_networks_or_ports" {
  name         = "Gambling via partner networks or proxy port"
  action       = "block"
  log_level    = "LOG_ALL"
  traffic_type = "PUBLIC_INTERNET"

  destination_content_category_ids     = [data.secureaccess_content_categories.example.items["Gambling"].id]
  destination_network_object_ids       = [secureaccess_network_object.partners.id]
  destination_service_object_ids       = [secureaccess_service_object.proxy_port.id]
  destination_operator                 = "AND"
  destination_operator_network_service = "OR"
}
```

Effect: gambling sites (content categories) **and** (the partner networks **or** the proxy port).

With `destination_operator_network_service = "AND"`, the same attributes would mean gambling sites **and** the
partner networks **and** the proxy port.

## Summary of the combinations

Using **N** for network objects, **S** for service objects, **C** for content categories and **A** for applications:

| Configured | `destination_operator` | `destination_operator_network_service` | Matches |
|---|---|---|---|
| N, S | `OR` | - | N or S |
| N, S | `AND` | `AND` | N and S |
| N, S | `AND` | `OR` | not allowed, a single group |
| C, N, S | `OR` | - | C or N or S |
| C, N, S | `AND` | `AND` | C and N and S |
| C, N, S | `AND` | `OR` | C and (N or S) |
| A, C | `AND` | - | A and C |
| A, C, N | `AND` | - | A and C and N |
| N only (any number of objects and groups) | `AND` | - | not allowed, a single group |
| N, destination list | `AND` | - | not allowed, destination lists cannot be used with `AND` |

## Literal IP addresses

Instead of network and service objects, IP addresses can be configured directly in the rule.

| Attribute | Rule type | Value | With `AND` |
|---|---|---|---|
| `source_private_literals` | `PRIVATE_NETWORK` | IPs or CIDRs | no |
| `source_internet_literals` | `PUBLIC_INTERNET` | **one** item of IPs or CIDRs with ports | yes |
| `destination_literals` | both | items of IPs or CIDRs with ports and protocol | no |

```terraform
resource "secureaccess_access_rule" "branch_to_dns" {
  name         = "Branch subnets to internal DNS"
  action       = "allow"
  log_level    = "LOG_ALL"
  traffic_type = "PRIVATE_NETWORK"

  source_private_literals = ["10.10.0.0/16", "10.20.0.0/16"]
  destination_literals = [
    { ips = ["10.0.0.53", "10.0.1.53"], ports = ["53"], protocol = "UDP" },
    { ips = ["10.0.0.53", "10.0.1.53"], ports = ["53"], protocol = "TCP" },
  ]
}
```

Notes:

- Ranges like `10.0.0.1-10.0.0.9` and FQDNs are not supported, use CIDRs.
- `ports` accepts single ports and ranges (`"8000-8010"`). `destination_literals` always requires ports, use `"0-65535"`
  for any port, including with `protocol = "ICMP"`.
- With `source_internet_literals`, destinations cannot include service objects or service object groups, more than one
  `destination_literals` item, or a `destination_literals` item with protocol `ICMP`.

## Geolocations

Countries are configured with uppercase ISO 3166-1 alpha-2 codes, e.g. `["PL", "DE"]`.

- `source_geolocations` matches traffic from these countries. It is available for both rule types, see
  [Geolocations need `AND` when combined with identities](#geolocations-need-and-when-combined-with-identities).
- `destination_geolocations` matches traffic to these countries. It is only available for `PUBLIC_INTERNET` rules
  with `action = "block"`, and it must be the only destination attribute of the rule.
- A rule can have geolocations either as sources or as destinations, never both.

```terraform
resource "secureaccess_access_rule" "block_countries" {
  name         = "Block traffic to selected countries"
  action       = "block"
  log_level    = "LOG_ALL"
  traffic_type = "PUBLIC_INTERNET"

  destination_geolocations = ["KP", "IR"]
}
```

Effect: everybody is blocked from destinations located in North Korea or Iran.

## Security Group Tags

`destination_security_group_tag_ids` matches traffic to devices tagged with these Security Group Tags (SGTs). It is only
available for `PRIVATE_NETWORK` rules and cannot be used with `destination_operator = "AND"`. The IDs are those of the
identities of type `sgt_id`.

```terraform
data "secureaccess_identities" "sgt" {
  identity_types = ["sgt_id"]
}

resource "secureaccess_access_rule" "employees_to_developers" {
  name         = "Employees to developer devices"
  action       = "allow"
  log_level    = "LOG_ALL"
  traffic_type = "PRIVATE_NETWORK"

  source_identity_ids                = [data.secureaccess_identities.employees.items["directory_group"]["Employees"].id]
  destination_security_group_tag_ids = [data.secureaccess_identities.sgt.items["sgt_id"]["Developers"].id]
}
```

Effect: members of the *Employees* group are allowed to reach the devices tagged with the *Developers* SGT.

## Attributes that cannot be combined

- `destination_application_list_ids` cannot be combined with `destination_application_ids` or `destination_application_category_ids`.
- `destination_content_category_list_ids` cannot be combined with `destination_content_category_ids`.
- `destination_geolocations` cannot be combined with any other destination attribute, nor with `source_geolocations`.
- `warn` and `isolate` actions are only available for `PUBLIC_INTERNET` rules. `warn` also requires content categories or applications as destinations.

## Finding the IDs

Most attributes take IDs. They can come from resources managed in the same configuration (e.g.
`secureaccess_network_object.example.id`), or from data sources for objects that already exist:

| Attribute | Data source |
|---|---|
| `source_identity_ids` | `secureaccess_identities`, `secureaccess_roaming_computer` (`origin_id`), `secureaccess_site` (`origin_id`) |
| `destination_content_category_ids` | `secureaccess_content_categories` |
| `destination_content_category_list_ids` | `secureaccess_content_category_list`, `secureaccess_content_category_lists` |
| `destination_application_ids` | `secureaccess_applications` (use `tonumber()` on the `id`) |
| `destination_application_category_ids` | `secureaccess_application_categories` (use `tonumber()` on the `id`) |
| `destination_application_list_ids` | `secureaccess_application_list` |
| `destination_list_ids` | `secureaccess_destination_list` |
| `destination_private_resource_ids`, `destination_private_resource_group_ids` | `secureaccess_private_resource`, `secureaccess_private_resource_group` |
| `destination_security_group_tag_ids` | `secureaccess_identities` with `identity_types = ["sgt_id"]` |

~> **Note:** Sites and Roaming Computers are referenced in `source_identity_ids` by their `origin_id`, not by their `id`.
Using `secureaccess_site.example.id` or `secureaccess_roaming_computer.example.id` fails at apply, use
`secureaccess_site.example.origin_id` or `secureaccess_roaming_computer.example.origin_id` instead.

```terraform
data "secureaccess_content_categories" "example" {
  items = {
    "Social Networking" = {}
    "Gambling"          = {}
  }
}

data "secureaccess_identities" "finance" {
  identity_types = ["directory_group"]
  filter         = "Finance"
}
```

Objects created in the same `terraform apply` are sometimes not yet known to the rules API, which then rejects their IDs.
The provider retries such requests up to 3 times before failing.

## Rules created in the dashboard

Existing rules can be imported with `terraform import secureaccess_access_rule.example <rule ID>`. The provider reads
the operators back from the rule. Rules using attributes not supported by the provider (e.g. Private Applications) lose
those conditions on the next update from Terraform.
