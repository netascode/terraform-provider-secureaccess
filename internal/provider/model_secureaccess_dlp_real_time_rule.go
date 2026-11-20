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
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/terraform-provider-secureaccess/internal/provider/helpers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types

type DLPRealTimeRule struct {
	Id                                  types.String                                   `tfsdk:"id"`
	Name                                types.String                                   `tfsdk:"name"`
	Description                         types.String                                   `tfsdk:"description"`
	Severity                            types.String                                   `tfsdk:"severity"`
	Enabled                             types.Bool                                     `tfsdk:"enabled"`
	DataClassificationScannableContexts types.Set                                      `tfsdk:"data_classification_scannable_contexts"`
	DataClassifications                 []DLPRealTimeRuleDataClassifications           `tfsdk:"data_classifications"`
	FileLabels                          types.Set                                      `tfsdk:"file_labels"`
	FileSizeConfigurationType           types.String                                   `tfsdk:"file_size_configuration_type"`
	FileSizeFrom                        types.Int64                                    `tfsdk:"file_size_from"`
	FileSizeFromUnit                    types.String                                   `tfsdk:"file_size_from_unit"`
	FileSizeTo                          types.Int64                                    `tfsdk:"file_size_to"`
	FileSizeToUnit                      types.String                                   `tfsdk:"file_size_to_unit"`
	FileTypesIncluded                   []DLPRealTimeRuleFileTypesIncluded             `tfsdk:"file_types_included"`
	FileTypeGroupsIncluded              []DLPRealTimeRuleFileTypeGroupsIncluded        `tfsdk:"file_type_groups_included"`
	FileTypesExcluded                   []DLPRealTimeRuleFileTypesExcluded             `tfsdk:"file_types_excluded"`
	FileTypeGroupsExcluded              []DLPRealTimeRuleFileTypeGroupsExcluded        `tfsdk:"file_type_groups_excluded"`
	UnclassifiedFilesScan               types.Bool                                     `tfsdk:"unclassified_files_scan"`
	IdentitiesIncluded                  []DLPRealTimeRuleIdentitiesIncluded            `tfsdk:"identities_included"`
	IdentitiesExcluded                  []DLPRealTimeRuleIdentitiesExcluded            `tfsdk:"identities_excluded"`
	AllDestinationsScope                types.String                                   `tfsdk:"all_destinations_scope"`
	DestinationListIncludedIds          types.Set                                      `tfsdk:"destination_list_included_ids"`
	DestinationListExcludedIds          types.Set                                      `tfsdk:"destination_list_excluded_ids"`
	ApplicationsIncluded                []DLPRealTimeRuleApplicationsIncluded          `tfsdk:"applications_included"`
	ApplicationsExcluded                types.Set                                      `tfsdk:"applications_excluded"`
	ApplicationCategoriesIncluded       []DLPRealTimeRuleApplicationCategoriesIncluded `tfsdk:"application_categories_included"`
	ApplicationCategoriesExcluded       types.Set                                      `tfsdk:"application_categories_excluded"`
	PrivateResourcesIncluded            []DLPRealTimeRulePrivateResourcesIncluded      `tfsdk:"private_resources_included"`
	PrivateResourcesExcluded            []DLPRealTimeRulePrivateResourcesExcluded      `tfsdk:"private_resources_excluded"`
	PrivateResourceGroupsIncluded       []DLPRealTimeRulePrivateResourceGroupsIncluded `tfsdk:"private_resource_groups_included"`
	PrivateResourceGroupsExcluded       []DLPRealTimeRulePrivateResourceGroupsExcluded `tfsdk:"private_resource_groups_excluded"`
	Action                              types.String                                   `tfsdk:"action"`
	OverrideBlock                       types.Bool                                     `tfsdk:"override_block"`
	NotificationPopup                   types.Bool                                     `tfsdk:"notification_popup"`
	NotificationTemplateId              types.String                                   `tfsdk:"notification_template_id"`
	EmailNotifyActor                    types.Bool                                     `tfsdk:"email_notify_actor"`
	EmailNotifyActorsManager            types.Bool                                     `tfsdk:"email_notify_actors_manager"`
	EmailNotifyCustomRecipients         types.Set                                      `tfsdk:"email_notify_custom_recipients"`
	EmailNotificationTemplateId         types.String                                   `tfsdk:"email_notification_template_id"`
}

type DLPRealTimeRuleDataClassifications struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type DLPRealTimeRuleFileTypesIncluded struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type DLPRealTimeRuleFileTypeGroupsIncluded struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type DLPRealTimeRuleFileTypesExcluded struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type DLPRealTimeRuleFileTypeGroupsExcluded struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type DLPRealTimeRuleIdentitiesIncluded struct {
	Id     types.Int64 `tfsdk:"id"`
	TypeId types.Int64 `tfsdk:"type_id"`
}

type DLPRealTimeRuleIdentitiesExcluded struct {
	Id     types.Int64 `tfsdk:"id"`
	TypeId types.Int64 `tfsdk:"type_id"`
}

type DLPRealTimeRuleApplicationsIncluded struct {
	Id               types.Int64  `tfsdk:"id"`
	TrafficDirection types.String `tfsdk:"traffic_direction"`
}

type DLPRealTimeRuleApplicationCategoriesIncluded struct {
	Id               types.Int64  `tfsdk:"id"`
	TrafficDirection types.String `tfsdk:"traffic_direction"`
}

type DLPRealTimeRulePrivateResourcesIncluded struct {
	Id               types.Int64  `tfsdk:"id"`
	TrafficDirection types.String `tfsdk:"traffic_direction"`
}

type DLPRealTimeRulePrivateResourcesExcluded struct {
	Id               types.Int64  `tfsdk:"id"`
	TrafficDirection types.String `tfsdk:"traffic_direction"`
}

type DLPRealTimeRulePrivateResourceGroupsIncluded struct {
	Id               types.Int64  `tfsdk:"id"`
	TrafficDirection types.String `tfsdk:"traffic_direction"`
}

type DLPRealTimeRulePrivateResourceGroupsExcluded struct {
	Id               types.Int64  `tfsdk:"id"`
	TrafficDirection types.String `tfsdk:"traffic_direction"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath

func (data DLPRealTimeRule) getPath() string {
	return "/policies/v2/dlp/realTime/rules"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin toBody

func (data DLPRealTimeRule) toBody(ctx context.Context, state DLPRealTimeRule) string {
	body := ""
	if data.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", data.Id.ValueString())
	}
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	if !data.Description.IsNull() {
		body, _ = sjson.Set(body, "description", data.Description.ValueString())
	}
	if !data.Severity.IsNull() {
		body, _ = sjson.Set(body, "severity", data.Severity.ValueString())
	}
	if !data.Enabled.IsNull() {
		body, _ = sjson.Set(body, "enabled", data.Enabled.ValueBool())
	}
	if !data.DataClassificationScannableContexts.IsNull() {
		var values []string
		data.DataClassificationScannableContexts.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "scannableContexts", values)
	}
	if len(data.DataClassifications) > 0 {
		var dataClassificationsBody strings.Builder
		dataClassificationsBody.WriteString("[")
		for _, item := range data.DataClassifications {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueString())
			}
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if itemBody != "" {
				if dataClassificationsBody.Len() > 1 {
					dataClassificationsBody.WriteString(",")
				}
				dataClassificationsBody.WriteString(itemBody)
			}
		}
		dataClassificationsBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "classifications", dataClassificationsBody.String())
	}
	if !data.FileLabels.IsNull() {
		var values []string
		data.FileLabels.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "mipTags", values)
	}
	if !data.FileSizeConfigurationType.IsNull() {
		body, _ = sjson.Set(body, "fileSizeConfigType", data.FileSizeConfigurationType.ValueString())
	}
	if !data.FileSizeFrom.IsNull() {
		body, _ = sjson.Set(body, "fileSizeFrom", data.FileSizeFrom.ValueInt64())
	}
	if !data.FileSizeFromUnit.IsNull() {
		body, _ = sjson.Set(body, "fileSizeFromUnit", data.FileSizeFromUnit.ValueString())
	}
	if !data.FileSizeTo.IsNull() {
		body, _ = sjson.Set(body, "fileSizeTo", data.FileSizeTo.ValueInt64())
	}
	if !data.FileSizeToUnit.IsNull() {
		body, _ = sjson.Set(body, "fileSizeToUnit", data.FileSizeToUnit.ValueString())
	}
	if len(data.FileTypesIncluded) > 0 {
		var fileTypesIncludedBody strings.Builder
		fileTypesIncludedBody.WriteString("[")
		for _, item := range data.FileTypesIncluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if itemBody != "" {
				if fileTypesIncludedBody.Len() > 1 {
					fileTypesIncludedBody.WriteString(",")
				}
				fileTypesIncludedBody.WriteString(itemBody)
			}
		}
		fileTypesIncludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "fileTypesIncluded", fileTypesIncludedBody.String())
	}
	if len(data.FileTypeGroupsIncluded) > 0 {
		var fileTypeGroupsIncludedBody strings.Builder
		fileTypeGroupsIncludedBody.WriteString("[")
		for _, item := range data.FileTypeGroupsIncluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if itemBody != "" {
				if fileTypeGroupsIncludedBody.Len() > 1 {
					fileTypeGroupsIncludedBody.WriteString(",")
				}
				fileTypeGroupsIncludedBody.WriteString(itemBody)
			}
		}
		fileTypeGroupsIncludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "fileTypeGroupsIncluded", fileTypeGroupsIncludedBody.String())
	}
	if len(data.FileTypesExcluded) > 0 {
		var fileTypesExcludedBody strings.Builder
		fileTypesExcludedBody.WriteString("[")
		for _, item := range data.FileTypesExcluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if itemBody != "" {
				if fileTypesExcludedBody.Len() > 1 {
					fileTypesExcludedBody.WriteString(",")
				}
				fileTypesExcludedBody.WriteString(itemBody)
			}
		}
		fileTypesExcludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "fileTypesExcluded", fileTypesExcludedBody.String())
	}
	if len(data.FileTypeGroupsExcluded) > 0 {
		var fileTypeGroupsExcludedBody strings.Builder
		fileTypeGroupsExcludedBody.WriteString("[")
		for _, item := range data.FileTypeGroupsExcluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if itemBody != "" {
				if fileTypeGroupsExcludedBody.Len() > 1 {
					fileTypeGroupsExcludedBody.WriteString(",")
				}
				fileTypeGroupsExcludedBody.WriteString(itemBody)
			}
		}
		fileTypeGroupsExcludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "fileTypeGroupsExcluded", fileTypeGroupsExcludedBody.String())
	}
	if !data.UnclassifiedFilesScan.IsNull() {
		body, _ = sjson.Set(body, "detectUnclassifiedFiles", data.UnclassifiedFilesScan.ValueBool())
	}
	if len(data.IdentitiesIncluded) > 0 {
		var identitiesIncludedBody strings.Builder
		identitiesIncludedBody.WriteString("[")
		for _, item := range data.IdentitiesIncluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "originId", item.Id.ValueInt64())
			}
			if !item.TypeId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "originTypeId", item.TypeId.ValueInt64())
			}
			if itemBody != "" {
				if identitiesIncludedBody.Len() > 1 {
					identitiesIncludedBody.WriteString(",")
				}
				identitiesIncludedBody.WriteString(itemBody)
			}
		}
		identitiesIncludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "identities", identitiesIncludedBody.String())
	}
	if len(data.IdentitiesExcluded) > 0 {
		var identitiesExcludedBody strings.Builder
		identitiesExcludedBody.WriteString("[")
		for _, item := range data.IdentitiesExcluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "originId", item.Id.ValueInt64())
			}
			if !item.TypeId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "originTypeId", item.TypeId.ValueInt64())
			}
			if itemBody != "" {
				if identitiesExcludedBody.Len() > 1 {
					identitiesExcludedBody.WriteString(",")
				}
				identitiesExcludedBody.WriteString(itemBody)
			}
		}
		identitiesExcludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "excludedIdentities", identitiesExcludedBody.String())
	}
	if !data.AllDestinationsScope.IsNull() {
		body, _ = sjson.Set(body, "allDestinationsScope", data.AllDestinationsScope.ValueString())
	}
	if !data.DestinationListIncludedIds.IsNull() {
		var values []string
		data.DestinationListIncludedIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "destinationListIds", values)
	}
	if !data.DestinationListExcludedIds.IsNull() {
		var values []string
		data.DestinationListExcludedIds.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "excludedDestinationListIds", values)
	}
	if len(data.ApplicationsIncluded) > 0 {
		var applicationsIncludedBody strings.Builder
		applicationsIncludedBody.WriteString("[")
		for _, item := range data.ApplicationsIncluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.TrafficDirection.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "trafficDirection", item.TrafficDirection.ValueString())
			}
			if itemBody != "" {
				if applicationsIncludedBody.Len() > 1 {
					applicationsIncludedBody.WriteString(",")
				}
				applicationsIncludedBody.WriteString(itemBody)
			}
		}
		applicationsIncludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "applications", applicationsIncludedBody.String())
	}
	if !data.ApplicationsExcluded.IsNull() {
		var values []string
		data.ApplicationsExcluded.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "excludedApplicationIds", values)
	}
	if len(data.ApplicationCategoriesIncluded) > 0 {
		var applicationCategoriesIncludedBody strings.Builder
		applicationCategoriesIncludedBody.WriteString("[")
		for _, item := range data.ApplicationCategoriesIncluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.TrafficDirection.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "trafficDirection", item.TrafficDirection.ValueString())
			}
			if itemBody != "" {
				if applicationCategoriesIncludedBody.Len() > 1 {
					applicationCategoriesIncludedBody.WriteString(",")
				}
				applicationCategoriesIncludedBody.WriteString(itemBody)
			}
		}
		applicationCategoriesIncludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "applicationCategories", applicationCategoriesIncludedBody.String())
	}
	if !data.ApplicationCategoriesExcluded.IsNull() {
		var values []string
		data.ApplicationCategoriesExcluded.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "excludedApplicationCategoryIds", values)
	}
	if len(data.PrivateResourcesIncluded) > 0 {
		var privateResourcesIncludedBody strings.Builder
		privateResourcesIncludedBody.WriteString("[")
		for _, item := range data.PrivateResourcesIncluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.TrafficDirection.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "trafficDirection", item.TrafficDirection.ValueString())
			}
			if itemBody != "" {
				if privateResourcesIncludedBody.Len() > 1 {
					privateResourcesIncludedBody.WriteString(",")
				}
				privateResourcesIncludedBody.WriteString(itemBody)
			}
		}
		privateResourcesIncludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "privateResources", privateResourcesIncludedBody.String())
	}
	if len(data.PrivateResourcesExcluded) > 0 {
		var privateResourcesExcludedBody strings.Builder
		privateResourcesExcludedBody.WriteString("[")
		for _, item := range data.PrivateResourcesExcluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.TrafficDirection.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "trafficDirection", item.TrafficDirection.ValueString())
			}
			if itemBody != "" {
				if privateResourcesExcludedBody.Len() > 1 {
					privateResourcesExcludedBody.WriteString(",")
				}
				privateResourcesExcludedBody.WriteString(itemBody)
			}
		}
		privateResourcesExcludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "excludedPrivateResources", privateResourcesExcludedBody.String())
	}
	if len(data.PrivateResourceGroupsIncluded) > 0 {
		var privateResourceGroupsIncludedBody strings.Builder
		privateResourceGroupsIncludedBody.WriteString("[")
		for _, item := range data.PrivateResourceGroupsIncluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.TrafficDirection.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "trafficDirection", item.TrafficDirection.ValueString())
			}
			if itemBody != "" {
				if privateResourceGroupsIncludedBody.Len() > 1 {
					privateResourceGroupsIncludedBody.WriteString(",")
				}
				privateResourceGroupsIncludedBody.WriteString(itemBody)
			}
		}
		privateResourceGroupsIncludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "privateResourceGroups", privateResourceGroupsIncludedBody.String())
	}
	if len(data.PrivateResourceGroupsExcluded) > 0 {
		var privateResourceGroupsExcludedBody strings.Builder
		privateResourceGroupsExcludedBody.WriteString("[")
		for _, item := range data.PrivateResourceGroupsExcluded {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueInt64())
			}
			if !item.TrafficDirection.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "trafficDirection", item.TrafficDirection.ValueString())
			}
			if itemBody != "" {
				if privateResourceGroupsExcludedBody.Len() > 1 {
					privateResourceGroupsExcludedBody.WriteString(",")
				}
				privateResourceGroupsExcludedBody.WriteString(itemBody)
			}
		}
		privateResourceGroupsExcludedBody.WriteString("]")
		body, _ = sjson.SetRaw(body, "excludedPrivateResourceGroups", privateResourceGroupsExcludedBody.String())
	}
	if !data.Action.IsNull() {
		body, _ = sjson.Set(body, "action", data.Action.ValueString())
	}
	if !data.OverrideBlock.IsNull() {
		body, _ = sjson.Set(body, "overrideBlock", data.OverrideBlock.ValueBool())
	}
	if !data.NotificationPopup.IsNull() {
		body, _ = sjson.Set(body, "popupNotification", data.NotificationPopup.ValueBool())
	}
	if !data.NotificationTemplateId.IsNull() {
		body, _ = sjson.Set(body, "edlpNotificationTemplateId", data.NotificationTemplateId.ValueString())
	}
	if !data.EmailNotifyActor.IsNull() {
		body, _ = sjson.Set(body, "notifyActor", data.EmailNotifyActor.ValueBool())
	}
	if !data.EmailNotifyActorsManager.IsNull() {
		body, _ = sjson.Set(body, "notifyManager", data.EmailNotifyActorsManager.ValueBool())
	}
	if !data.EmailNotifyCustomRecipients.IsNull() {
		var values []string
		data.EmailNotifyCustomRecipients.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "customEmailsNotification", values)
	}
	if !data.EmailNotificationTemplateId.IsNull() {
		body, _ = sjson.Set(body, "notificationTemplateUuid", data.EmailNotificationTemplateId.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody

func (data *DLPRealTimeRule) fromBody(ctx context.Context, res gjson.Result) {
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
	if value := res.Get("severity"); value.Exists() {
		data.Severity = types.StringValue(value.String())
	} else {
		data.Severity = types.StringNull()
	}
	if value := res.Get("enabled"); value.Exists() {
		data.Enabled = types.BoolValue(value.Bool())
	} else {
		data.Enabled = types.BoolNull()
	}
	if value := res.Get("scannableContexts"); value.Exists() {
		data.DataClassificationScannableContexts = helpers.GetStringSet(value.Array())
	} else {
		data.DataClassificationScannableContexts = types.SetNull(types.StringType)
	}
	if value := res.Get("classifications"); value.Exists() {
		data.DataClassifications = make([]DLPRealTimeRuleDataClassifications, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRuleDataClassifications{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.StringValue(value.String())
			} else {
				data.Id = types.StringNull()
			}
			if value := res.Get("name"); value.Exists() {
				data.Name = types.StringValue(value.String())
			} else {
				data.Name = types.StringNull()
			}
			(*parent).DataClassifications = append((*parent).DataClassifications, data)
			return true
		})
	}
	if value := res.Get("mipTags"); value.Exists() {
		data.FileLabels = helpers.GetStringSet(value.Array())
	} else {
		data.FileLabels = types.SetNull(types.StringType)
	}
	if value := res.Get("fileSizeConfigType"); value.Exists() {
		data.FileSizeConfigurationType = types.StringValue(value.String())
	} else {
		data.FileSizeConfigurationType = types.StringNull()
	}
	if value := res.Get("fileSizeFrom"); value.Exists() {
		data.FileSizeFrom = types.Int64Value(value.Int())
	} else {
		data.FileSizeFrom = types.Int64Null()
	}
	if value := res.Get("fileSizeFromUnit"); value.Exists() {
		data.FileSizeFromUnit = types.StringValue(value.String())
	} else {
		data.FileSizeFromUnit = types.StringNull()
	}
	if value := res.Get("fileSizeTo"); value.Exists() {
		data.FileSizeTo = types.Int64Value(value.Int())
	} else {
		data.FileSizeTo = types.Int64Null()
	}
	if value := res.Get("fileSizeToUnit"); value.Exists() {
		data.FileSizeToUnit = types.StringValue(value.String())
	} else {
		data.FileSizeToUnit = types.StringNull()
	}
	if value := res.Get("fileTypesIncluded"); value.Exists() {
		data.FileTypesIncluded = make([]DLPRealTimeRuleFileTypesIncluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRuleFileTypesIncluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("name"); value.Exists() {
				data.Name = types.StringValue(value.String())
			} else {
				data.Name = types.StringNull()
			}
			(*parent).FileTypesIncluded = append((*parent).FileTypesIncluded, data)
			return true
		})
	}
	if value := res.Get("fileTypeGroupsIncluded"); value.Exists() {
		data.FileTypeGroupsIncluded = make([]DLPRealTimeRuleFileTypeGroupsIncluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRuleFileTypeGroupsIncluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("name"); value.Exists() {
				data.Name = types.StringValue(value.String())
			} else {
				data.Name = types.StringNull()
			}
			(*parent).FileTypeGroupsIncluded = append((*parent).FileTypeGroupsIncluded, data)
			return true
		})
	}
	if value := res.Get("fileTypesExcluded"); value.Exists() {
		data.FileTypesExcluded = make([]DLPRealTimeRuleFileTypesExcluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRuleFileTypesExcluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("name"); value.Exists() {
				data.Name = types.StringValue(value.String())
			} else {
				data.Name = types.StringNull()
			}
			(*parent).FileTypesExcluded = append((*parent).FileTypesExcluded, data)
			return true
		})
	}
	if value := res.Get("fileTypeGroupsExcluded"); value.Exists() {
		data.FileTypeGroupsExcluded = make([]DLPRealTimeRuleFileTypeGroupsExcluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRuleFileTypeGroupsExcluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("name"); value.Exists() {
				data.Name = types.StringValue(value.String())
			} else {
				data.Name = types.StringNull()
			}
			(*parent).FileTypeGroupsExcluded = append((*parent).FileTypeGroupsExcluded, data)
			return true
		})
	}
	if value := res.Get("detectUnclassifiedFiles"); value.Exists() {
		data.UnclassifiedFilesScan = types.BoolValue(value.Bool())
	} else {
		data.UnclassifiedFilesScan = types.BoolNull()
	}
	if value := res.Get("identities"); value.Exists() {
		data.IdentitiesIncluded = make([]DLPRealTimeRuleIdentitiesIncluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRuleIdentitiesIncluded{}
			if value := res.Get("originId"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("originTypeId"); value.Exists() {
				data.TypeId = types.Int64Value(value.Int())
			} else {
				data.TypeId = types.Int64Null()
			}
			(*parent).IdentitiesIncluded = append((*parent).IdentitiesIncluded, data)
			return true
		})
	}
	if value := res.Get("excludedIdentities"); value.Exists() {
		data.IdentitiesExcluded = make([]DLPRealTimeRuleIdentitiesExcluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRuleIdentitiesExcluded{}
			if value := res.Get("originId"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("originTypeId"); value.Exists() {
				data.TypeId = types.Int64Value(value.Int())
			} else {
				data.TypeId = types.Int64Null()
			}
			(*parent).IdentitiesExcluded = append((*parent).IdentitiesExcluded, data)
			return true
		})
	}
	if value := res.Get("allDestinationsScope"); value.Exists() {
		data.AllDestinationsScope = types.StringValue(value.String())
	} else {
		data.AllDestinationsScope = types.StringNull()
	}
	if value := res.Get("destinationListIds"); value.Exists() {
		data.DestinationListIncludedIds = helpers.GetStringSet(value.Array())
	} else {
		data.DestinationListIncludedIds = types.SetNull(types.StringType)
	}
	if value := res.Get("excludedDestinationListIds"); value.Exists() {
		data.DestinationListExcludedIds = helpers.GetStringSet(value.Array())
	} else {
		data.DestinationListExcludedIds = types.SetNull(types.StringType)
	}
	if value := res.Get("applications"); value.Exists() {
		data.ApplicationsIncluded = make([]DLPRealTimeRuleApplicationsIncluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRuleApplicationsIncluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("trafficDirection"); value.Exists() {
				data.TrafficDirection = types.StringValue(value.String())
			} else {
				data.TrafficDirection = types.StringNull()
			}
			(*parent).ApplicationsIncluded = append((*parent).ApplicationsIncluded, data)
			return true
		})
	}
	if value := res.Get("excludedApplicationIds"); value.Exists() {
		data.ApplicationsExcluded = helpers.GetStringSet(value.Array())
	} else {
		data.ApplicationsExcluded = types.SetNull(types.StringType)
	}
	if value := res.Get("applicationCategories"); value.Exists() {
		data.ApplicationCategoriesIncluded = make([]DLPRealTimeRuleApplicationCategoriesIncluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRuleApplicationCategoriesIncluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("trafficDirection"); value.Exists() {
				data.TrafficDirection = types.StringValue(value.String())
			} else {
				data.TrafficDirection = types.StringNull()
			}
			(*parent).ApplicationCategoriesIncluded = append((*parent).ApplicationCategoriesIncluded, data)
			return true
		})
	}
	if value := res.Get("excludedApplicationCategoryIds"); value.Exists() {
		data.ApplicationCategoriesExcluded = helpers.GetStringSet(value.Array())
	} else {
		data.ApplicationCategoriesExcluded = types.SetNull(types.StringType)
	}
	if value := res.Get("privateResources"); value.Exists() {
		data.PrivateResourcesIncluded = make([]DLPRealTimeRulePrivateResourcesIncluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRulePrivateResourcesIncluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("trafficDirection"); value.Exists() {
				data.TrafficDirection = types.StringValue(value.String())
			} else {
				data.TrafficDirection = types.StringNull()
			}
			(*parent).PrivateResourcesIncluded = append((*parent).PrivateResourcesIncluded, data)
			return true
		})
	}
	if value := res.Get("excludedPrivateResources"); value.Exists() {
		data.PrivateResourcesExcluded = make([]DLPRealTimeRulePrivateResourcesExcluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRulePrivateResourcesExcluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("trafficDirection"); value.Exists() {
				data.TrafficDirection = types.StringValue(value.String())
			} else {
				data.TrafficDirection = types.StringNull()
			}
			(*parent).PrivateResourcesExcluded = append((*parent).PrivateResourcesExcluded, data)
			return true
		})
	}
	if value := res.Get("privateResourceGroups"); value.Exists() {
		data.PrivateResourceGroupsIncluded = make([]DLPRealTimeRulePrivateResourceGroupsIncluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRulePrivateResourceGroupsIncluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("trafficDirection"); value.Exists() {
				data.TrafficDirection = types.StringValue(value.String())
			} else {
				data.TrafficDirection = types.StringNull()
			}
			(*parent).PrivateResourceGroupsIncluded = append((*parent).PrivateResourceGroupsIncluded, data)
			return true
		})
	}
	if value := res.Get("excludedPrivateResourceGroups"); value.Exists() {
		data.PrivateResourceGroupsExcluded = make([]DLPRealTimeRulePrivateResourceGroupsExcluded, 0, int(value.Get("#").Int()))
		value.ForEach(func(k, res gjson.Result) bool {
			parent := &data
			data := DLPRealTimeRulePrivateResourceGroupsExcluded{}
			if value := res.Get("id"); value.Exists() {
				data.Id = types.Int64Value(value.Int())
			} else {
				data.Id = types.Int64Null()
			}
			if value := res.Get("trafficDirection"); value.Exists() {
				data.TrafficDirection = types.StringValue(value.String())
			} else {
				data.TrafficDirection = types.StringNull()
			}
			(*parent).PrivateResourceGroupsExcluded = append((*parent).PrivateResourceGroupsExcluded, data)
			return true
		})
	}
	if value := res.Get("action"); value.Exists() {
		data.Action = types.StringValue(value.String())
	} else {
		data.Action = types.StringNull()
	}
	if value := res.Get("overrideBlock"); value.Exists() {
		data.OverrideBlock = types.BoolValue(value.Bool())
	} else {
		data.OverrideBlock = types.BoolNull()
	}
	if value := res.Get("popupNotification"); value.Exists() {
		data.NotificationPopup = types.BoolValue(value.Bool())
	} else {
		data.NotificationPopup = types.BoolNull()
	}
	if value := res.Get("edlpNotificationTemplateId"); value.Exists() {
		data.NotificationTemplateId = types.StringValue(value.String())
	} else {
		data.NotificationTemplateId = types.StringNull()
	}
	if value := res.Get("notifyActor"); value.Exists() {
		data.EmailNotifyActor = types.BoolValue(value.Bool())
	} else {
		data.EmailNotifyActor = types.BoolNull()
	}
	if value := res.Get("notifyManager"); value.Exists() {
		data.EmailNotifyActorsManager = types.BoolValue(value.Bool())
	} else {
		data.EmailNotifyActorsManager = types.BoolNull()
	}
	if value := res.Get("customEmailsNotification"); value.Exists() {
		data.EmailNotifyCustomRecipients = helpers.GetStringSet(value.Array())
	} else {
		data.EmailNotifyCustomRecipients = types.SetNull(types.StringType)
	}
	if value := res.Get("notificationTemplateUuid"); value.Exists() {
		data.EmailNotificationTemplateId = types.StringValue(value.String())
	} else {
		data.EmailNotificationTemplateId = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyPartial

// fromBodyPartial reads values from a gjson.Result into a tfstate model. It ignores null attributes in order to
// uncouple the provider from the exact values that the backend API might summon to replace nulls. (Such behavior might
// easily change across versions of the backend API.) For List/Set/Map attributes, the func only updates the
// "managed" elements, instead of all elements.
func (data *DLPRealTimeRule) fromBodyPartial(ctx context.Context, res gjson.Result) {
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
	if value := res.Get("severity"); value.Exists() && !data.Severity.IsNull() {
		data.Severity = types.StringValue(value.String())
	} else {
		data.Severity = types.StringNull()
	}
	if value := res.Get("enabled"); value.Exists() && !data.Enabled.IsNull() {
		data.Enabled = types.BoolValue(value.Bool())
	} else {
		data.Enabled = types.BoolNull()
	}
	if value := res.Get("scannableContexts"); value.Exists() && !data.DataClassificationScannableContexts.IsNull() {
		data.DataClassificationScannableContexts = helpers.GetStringSet(value.Array())
	} else {
		data.DataClassificationScannableContexts = types.SetNull(types.StringType)
	}
	dataClassificationsArray := res.Get("classifications")
	for i := 0; i < len(data.DataClassifications); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{data.DataClassifications[i].Id.ValueString()}

		parent := &data
		data := (*parent).DataClassifications[i]
		var res gjson.Result

		dataClassificationsArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing DataClassifications[%d] = %+v",
				i,
				(*parent).DataClassifications[i],
			))
			(*parent).DataClassifications = slices.Delete((*parent).DataClassifications, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.StringValue(value.String())
		} else {
			data.Id = types.StringNull()
		}
		if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
			data.Name = types.StringValue(value.String())
		} else {
			data.Name = types.StringNull()
		}
		(*parent).DataClassifications[i] = data
	}
	if value := res.Get("mipTags"); value.Exists() && !data.FileLabels.IsNull() {
		data.FileLabels = helpers.GetStringSet(value.Array())
	} else {
		data.FileLabels = types.SetNull(types.StringType)
	}
	if value := res.Get("fileSizeConfigType"); value.Exists() && !data.FileSizeConfigurationType.IsNull() {
		data.FileSizeConfigurationType = types.StringValue(value.String())
	} else {
		data.FileSizeConfigurationType = types.StringNull()
	}
	if value := res.Get("fileSizeFrom"); value.Exists() && !data.FileSizeFrom.IsNull() {
		data.FileSizeFrom = types.Int64Value(value.Int())
	} else {
		data.FileSizeFrom = types.Int64Null()
	}
	if value := res.Get("fileSizeFromUnit"); value.Exists() && !data.FileSizeFromUnit.IsNull() {
		data.FileSizeFromUnit = types.StringValue(value.String())
	} else {
		data.FileSizeFromUnit = types.StringNull()
	}
	if value := res.Get("fileSizeTo"); value.Exists() && !data.FileSizeTo.IsNull() {
		data.FileSizeTo = types.Int64Value(value.Int())
	} else {
		data.FileSizeTo = types.Int64Null()
	}
	if value := res.Get("fileSizeToUnit"); value.Exists() && !data.FileSizeToUnit.IsNull() {
		data.FileSizeToUnit = types.StringValue(value.String())
	} else {
		data.FileSizeToUnit = types.StringNull()
	}
	fileTypesIncludedArray := res.Get("fileTypesIncluded")
	for i := 0; i < len(data.FileTypesIncluded); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.FileTypesIncluded[i].Id.ValueInt64(), 10)}

		parent := &data
		data := (*parent).FileTypesIncluded[i]
		var res gjson.Result

		fileTypesIncludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing FileTypesIncluded[%d] = %+v",
				i,
				(*parent).FileTypesIncluded[i],
			))
			(*parent).FileTypesIncluded = slices.Delete((*parent).FileTypesIncluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
			data.Name = types.StringValue(value.String())
		} else {
			data.Name = types.StringNull()
		}
		(*parent).FileTypesIncluded[i] = data
	}
	fileTypeGroupsIncludedArray := res.Get("fileTypeGroupsIncluded")
	for i := 0; i < len(data.FileTypeGroupsIncluded); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.FileTypeGroupsIncluded[i].Id.ValueInt64(), 10)}

		parent := &data
		data := (*parent).FileTypeGroupsIncluded[i]
		var res gjson.Result

		fileTypeGroupsIncludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing FileTypeGroupsIncluded[%d] = %+v",
				i,
				(*parent).FileTypeGroupsIncluded[i],
			))
			(*parent).FileTypeGroupsIncluded = slices.Delete((*parent).FileTypeGroupsIncluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
			data.Name = types.StringValue(value.String())
		} else {
			data.Name = types.StringNull()
		}
		(*parent).FileTypeGroupsIncluded[i] = data
	}
	fileTypesExcludedArray := res.Get("fileTypesExcluded")
	for i := 0; i < len(data.FileTypesExcluded); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.FileTypesExcluded[i].Id.ValueInt64(), 10)}

		parent := &data
		data := (*parent).FileTypesExcluded[i]
		var res gjson.Result

		fileTypesExcludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing FileTypesExcluded[%d] = %+v",
				i,
				(*parent).FileTypesExcluded[i],
			))
			(*parent).FileTypesExcluded = slices.Delete((*parent).FileTypesExcluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
			data.Name = types.StringValue(value.String())
		} else {
			data.Name = types.StringNull()
		}
		(*parent).FileTypesExcluded[i] = data
	}
	fileTypeGroupsExcludedArray := res.Get("fileTypeGroupsExcluded")
	for i := 0; i < len(data.FileTypeGroupsExcluded); i++ {
		keys := [...]string{"id"}
		keyValues := [...]string{strconv.FormatInt(data.FileTypeGroupsExcluded[i].Id.ValueInt64(), 10)}

		parent := &data
		data := (*parent).FileTypeGroupsExcluded[i]
		var res gjson.Result

		fileTypeGroupsExcludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing FileTypeGroupsExcluded[%d] = %+v",
				i,
				(*parent).FileTypeGroupsExcluded[i],
			))
			(*parent).FileTypeGroupsExcluded = slices.Delete((*parent).FileTypeGroupsExcluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
			data.Name = types.StringValue(value.String())
		} else {
			data.Name = types.StringNull()
		}
		(*parent).FileTypeGroupsExcluded[i] = data
	}
	if value := res.Get("detectUnclassifiedFiles"); value.Exists() && !data.UnclassifiedFilesScan.IsNull() {
		data.UnclassifiedFilesScan = types.BoolValue(value.Bool())
	} else {
		data.UnclassifiedFilesScan = types.BoolNull()
	}
	identitiesIncludedArray := res.Get("identities")
	for i := 0; i < len(data.IdentitiesIncluded); i++ {
		keys := [...]string{"originTypeId"}
		keyValues := [...]string{strconv.FormatInt(data.IdentitiesIncluded[i].TypeId.ValueInt64(), 10)}

		parent := &data
		data := (*parent).IdentitiesIncluded[i]
		var res gjson.Result

		identitiesIncludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing IdentitiesIncluded[%d] = %+v",
				i,
				(*parent).IdentitiesIncluded[i],
			))
			(*parent).IdentitiesIncluded = slices.Delete((*parent).IdentitiesIncluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("originId"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("originTypeId"); value.Exists() && !data.TypeId.IsNull() {
			data.TypeId = types.Int64Value(value.Int())
		} else {
			data.TypeId = types.Int64Null()
		}
		(*parent).IdentitiesIncluded[i] = data
	}
	identitiesExcludedArray := res.Get("excludedIdentities")
	for i := 0; i < len(data.IdentitiesExcluded); i++ {
		keys := [...]string{"originTypeId"}
		keyValues := [...]string{strconv.FormatInt(data.IdentitiesExcluded[i].TypeId.ValueInt64(), 10)}

		parent := &data
		data := (*parent).IdentitiesExcluded[i]
		var res gjson.Result

		identitiesExcludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing IdentitiesExcluded[%d] = %+v",
				i,
				(*parent).IdentitiesExcluded[i],
			))
			(*parent).IdentitiesExcluded = slices.Delete((*parent).IdentitiesExcluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("originId"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("originTypeId"); value.Exists() && !data.TypeId.IsNull() {
			data.TypeId = types.Int64Value(value.Int())
		} else {
			data.TypeId = types.Int64Null()
		}
		(*parent).IdentitiesExcluded[i] = data
	}
	if value := res.Get("allDestinationsScope"); value.Exists() && !data.AllDestinationsScope.IsNull() {
		data.AllDestinationsScope = types.StringValue(value.String())
	} else {
		data.AllDestinationsScope = types.StringNull()
	}
	if value := res.Get("destinationListIds"); value.Exists() && !data.DestinationListIncludedIds.IsNull() {
		data.DestinationListIncludedIds = helpers.GetStringSet(value.Array())
	} else {
		data.DestinationListIncludedIds = types.SetNull(types.StringType)
	}
	if value := res.Get("excludedDestinationListIds"); value.Exists() && !data.DestinationListExcludedIds.IsNull() {
		data.DestinationListExcludedIds = helpers.GetStringSet(value.Array())
	} else {
		data.DestinationListExcludedIds = types.SetNull(types.StringType)
	}
	applicationsIncludedArray := res.Get("applications")
	for i := 0; i < len(data.ApplicationsIncluded); i++ {
		keys := [...]string{"id", "trafficDirection"}
		keyValues := [...]string{strconv.FormatInt(data.ApplicationsIncluded[i].Id.ValueInt64(), 10), data.ApplicationsIncluded[i].TrafficDirection.ValueString()}

		parent := &data
		data := (*parent).ApplicationsIncluded[i]
		var res gjson.Result

		applicationsIncludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing ApplicationsIncluded[%d] = %+v",
				i,
				(*parent).ApplicationsIncluded[i],
			))
			(*parent).ApplicationsIncluded = slices.Delete((*parent).ApplicationsIncluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("trafficDirection"); value.Exists() && !data.TrafficDirection.IsNull() {
			data.TrafficDirection = types.StringValue(value.String())
		} else {
			data.TrafficDirection = types.StringNull()
		}
		(*parent).ApplicationsIncluded[i] = data
	}
	if value := res.Get("excludedApplicationIds"); value.Exists() && !data.ApplicationsExcluded.IsNull() {
		data.ApplicationsExcluded = helpers.GetStringSet(value.Array())
	} else {
		data.ApplicationsExcluded = types.SetNull(types.StringType)
	}
	applicationCategoriesIncludedArray := res.Get("applicationCategories")
	for i := 0; i < len(data.ApplicationCategoriesIncluded); i++ {
		keys := [...]string{"id", "trafficDirection"}
		keyValues := [...]string{strconv.FormatInt(data.ApplicationCategoriesIncluded[i].Id.ValueInt64(), 10), data.ApplicationCategoriesIncluded[i].TrafficDirection.ValueString()}

		parent := &data
		data := (*parent).ApplicationCategoriesIncluded[i]
		var res gjson.Result

		applicationCategoriesIncludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing ApplicationCategoriesIncluded[%d] = %+v",
				i,
				(*parent).ApplicationCategoriesIncluded[i],
			))
			(*parent).ApplicationCategoriesIncluded = slices.Delete((*parent).ApplicationCategoriesIncluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("trafficDirection"); value.Exists() && !data.TrafficDirection.IsNull() {
			data.TrafficDirection = types.StringValue(value.String())
		} else {
			data.TrafficDirection = types.StringNull()
		}
		(*parent).ApplicationCategoriesIncluded[i] = data
	}
	if value := res.Get("excludedApplicationCategoryIds"); value.Exists() && !data.ApplicationCategoriesExcluded.IsNull() {
		data.ApplicationCategoriesExcluded = helpers.GetStringSet(value.Array())
	} else {
		data.ApplicationCategoriesExcluded = types.SetNull(types.StringType)
	}
	privateResourcesIncludedArray := res.Get("privateResources")
	for i := 0; i < len(data.PrivateResourcesIncluded); i++ {
		keys := [...]string{"id", "trafficDirection"}
		keyValues := [...]string{strconv.FormatInt(data.PrivateResourcesIncluded[i].Id.ValueInt64(), 10), data.PrivateResourcesIncluded[i].TrafficDirection.ValueString()}

		parent := &data
		data := (*parent).PrivateResourcesIncluded[i]
		var res gjson.Result

		privateResourcesIncludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing PrivateResourcesIncluded[%d] = %+v",
				i,
				(*parent).PrivateResourcesIncluded[i],
			))
			(*parent).PrivateResourcesIncluded = slices.Delete((*parent).PrivateResourcesIncluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("trafficDirection"); value.Exists() && !data.TrafficDirection.IsNull() {
			data.TrafficDirection = types.StringValue(value.String())
		} else {
			data.TrafficDirection = types.StringNull()
		}
		(*parent).PrivateResourcesIncluded[i] = data
	}
	privateResourcesExcludedArray := res.Get("excludedPrivateResources")
	for i := 0; i < len(data.PrivateResourcesExcluded); i++ {
		keys := [...]string{"id", "trafficDirection"}
		keyValues := [...]string{strconv.FormatInt(data.PrivateResourcesExcluded[i].Id.ValueInt64(), 10), data.PrivateResourcesExcluded[i].TrafficDirection.ValueString()}

		parent := &data
		data := (*parent).PrivateResourcesExcluded[i]
		var res gjson.Result

		privateResourcesExcludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing PrivateResourcesExcluded[%d] = %+v",
				i,
				(*parent).PrivateResourcesExcluded[i],
			))
			(*parent).PrivateResourcesExcluded = slices.Delete((*parent).PrivateResourcesExcluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("trafficDirection"); value.Exists() && !data.TrafficDirection.IsNull() {
			data.TrafficDirection = types.StringValue(value.String())
		} else {
			data.TrafficDirection = types.StringNull()
		}
		(*parent).PrivateResourcesExcluded[i] = data
	}
	privateResourceGroupsIncludedArray := res.Get("privateResourceGroups")
	for i := 0; i < len(data.PrivateResourceGroupsIncluded); i++ {
		keys := [...]string{"id", "trafficDirection"}
		keyValues := [...]string{strconv.FormatInt(data.PrivateResourceGroupsIncluded[i].Id.ValueInt64(), 10), data.PrivateResourceGroupsIncluded[i].TrafficDirection.ValueString()}

		parent := &data
		data := (*parent).PrivateResourceGroupsIncluded[i]
		var res gjson.Result

		privateResourceGroupsIncludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing PrivateResourceGroupsIncluded[%d] = %+v",
				i,
				(*parent).PrivateResourceGroupsIncluded[i],
			))
			(*parent).PrivateResourceGroupsIncluded = slices.Delete((*parent).PrivateResourceGroupsIncluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("trafficDirection"); value.Exists() && !data.TrafficDirection.IsNull() {
			data.TrafficDirection = types.StringValue(value.String())
		} else {
			data.TrafficDirection = types.StringNull()
		}
		(*parent).PrivateResourceGroupsIncluded[i] = data
	}
	privateResourceGroupsExcludedArray := res.Get("excludedPrivateResourceGroups")
	for i := 0; i < len(data.PrivateResourceGroupsExcluded); i++ {
		keys := [...]string{"id", "trafficDirection"}
		keyValues := [...]string{strconv.FormatInt(data.PrivateResourceGroupsExcluded[i].Id.ValueInt64(), 10), data.PrivateResourceGroupsExcluded[i].TrafficDirection.ValueString()}

		parent := &data
		data := (*parent).PrivateResourceGroupsExcluded[i]
		var res gjson.Result

		privateResourceGroupsExcludedArray.ForEach(
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
			tflog.Debug(ctx, fmt.Sprintf("removing PrivateResourceGroupsExcluded[%d] = %+v",
				i,
				(*parent).PrivateResourceGroupsExcluded[i],
			))
			(*parent).PrivateResourceGroupsExcluded = slices.Delete((*parent).PrivateResourceGroupsExcluded, i, i+1)
			i--

			continue
		}
		if value := res.Get("id"); value.Exists() && !data.Id.IsNull() {
			data.Id = types.Int64Value(value.Int())
		} else {
			data.Id = types.Int64Null()
		}
		if value := res.Get("trafficDirection"); value.Exists() && !data.TrafficDirection.IsNull() {
			data.TrafficDirection = types.StringValue(value.String())
		} else {
			data.TrafficDirection = types.StringNull()
		}
		(*parent).PrivateResourceGroupsExcluded[i] = data
	}
	if value := res.Get("action"); value.Exists() && !data.Action.IsNull() {
		data.Action = types.StringValue(value.String())
	} else {
		data.Action = types.StringNull()
	}
	if value := res.Get("overrideBlock"); value.Exists() && !data.OverrideBlock.IsNull() {
		data.OverrideBlock = types.BoolValue(value.Bool())
	} else {
		data.OverrideBlock = types.BoolNull()
	}
	if value := res.Get("popupNotification"); value.Exists() && !data.NotificationPopup.IsNull() {
		data.NotificationPopup = types.BoolValue(value.Bool())
	} else {
		data.NotificationPopup = types.BoolNull()
	}
	if value := res.Get("edlpNotificationTemplateId"); value.Exists() && !data.NotificationTemplateId.IsNull() {
		data.NotificationTemplateId = types.StringValue(value.String())
	} else {
		data.NotificationTemplateId = types.StringNull()
	}
	if value := res.Get("notifyActor"); value.Exists() && !data.EmailNotifyActor.IsNull() {
		data.EmailNotifyActor = types.BoolValue(value.Bool())
	} else {
		data.EmailNotifyActor = types.BoolNull()
	}
	if value := res.Get("notifyManager"); value.Exists() && !data.EmailNotifyActorsManager.IsNull() {
		data.EmailNotifyActorsManager = types.BoolValue(value.Bool())
	} else {
		data.EmailNotifyActorsManager = types.BoolNull()
	}
	if value := res.Get("customEmailsNotification"); value.Exists() && !data.EmailNotifyCustomRecipients.IsNull() {
		data.EmailNotifyCustomRecipients = helpers.GetStringSet(value.Array())
	} else {
		data.EmailNotifyCustomRecipients = types.SetNull(types.StringType)
	}
	if value := res.Get("notificationTemplateUuid"); value.Exists() && !data.EmailNotificationTemplateId.IsNull() {
		data.EmailNotificationTemplateId = types.StringValue(value.String())
	} else {
		data.EmailNotificationTemplateId = types.StringNull()
	}
}

// End of section. //template:end fromBodyPartial

// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// fromBodyUnknowns updates the Unknown Computed tfstate values from a JSON.
// Known values are not changed (usual for Computed attributes with UseStateForUnknown or with Default).
func (data *DLPRealTimeRule) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
}

// End of section. //template:end fromBodyUnknowns

// Section below is generated&owned by "gen/generator.go". //template:begin toBodyPutDelete

// End of section. //template:end toBodyPutDelete

func (data DLPRealTimeRule) adjustBody(ctx context.Context, req string) string {

	// API Needs this to be cleared for unclassified files to work
	if len(data.FileLabels.Elements()) == 0 {
		req, _ = sjson.SetRaw(req, "mipTags", "[]")
	}

	if data.EmailNotifyActor.ValueBool() && data.EmailNotificationTemplateId.IsNull() {
		req, _ = sjson.Set(req, "notificationTemplateUuid", "00000000-0000-0000-0000-000000000000")
	}

	req, _ = sjson.Delete(req, "classifications")
	var classificationIds []string
	for _, classification := range data.DataClassifications {
		classificationIds = append(classificationIds, classification.Id.ValueString())
	}
	req, _ = sjson.Set(req, "classifications", classificationIds)

	req, _ = sjson.Delete(req, "fileTypesIncluded")
	var fileTypeIncludedIds []int64
	for _, fileType := range data.FileTypesIncluded {
		fileTypeIncludedIds = append(fileTypeIncludedIds, fileType.Id.ValueInt64())
	}
	req, _ = sjson.Set(req, "fileTypesIncluded", fileTypeIncludedIds)

	req, _ = sjson.Delete(req, "fileTypeGroupsIncluded")
	var fileTypeGroupIncludedIds []int64
	for _, fileTypeGroup := range data.FileTypeGroupsIncluded {
		fileTypeGroupIncludedIds = append(fileTypeGroupIncludedIds, fileTypeGroup.Id.ValueInt64())
	}
	req, _ = sjson.Set(req, "fileTypeGroupsIncluded", fileTypeGroupIncludedIds)

	req, _ = sjson.Delete(req, "fileTypesExcluded")
	var fileTypeExcludedIds []int64
	for _, fileType := range data.FileTypesExcluded {
		fileTypeExcludedIds = append(fileTypeExcludedIds, fileType.Id.ValueInt64())
	}
	req, _ = sjson.Set(req, "fileTypesExcluded", fileTypeExcludedIds)

	req, _ = sjson.Delete(req, "fileTypeGroupsExcluded")
	var fileTypeGroupExcludedIds []int64
	for _, fileTypeGroup := range data.FileTypeGroupsExcluded {
		fileTypeGroupExcludedIds = append(fileTypeGroupExcludedIds, fileTypeGroup.Id.ValueInt64())
	}
	req, _ = sjson.Set(req, "fileTypeGroupsExcluded", fileTypeGroupExcludedIds)

	return req
}

// Section below is generated&owned by "gen/generator.go". //template:begin adjustRead

// End of section. //template:end adjustRead
