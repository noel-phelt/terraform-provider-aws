// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configservice

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func expandAccountAggregationSources(configured []interface{}) []*configservice.AccountAggregationSource {
	var results []*configservice.AccountAggregationSource
	for _, item := range configured {
		detail := item.(map[string]interface{})
		source := configservice.AccountAggregationSource{
			AllAwsRegions: aws.Bool(detail["all_regions"].(bool)),
		}

		if v, ok := detail["account_ids"]; ok {
			accountIDs := v.([]interface{})
			if len(accountIDs) > 0 {
				source.AccountIds = flex.ExpandStringList(accountIDs)
			}
		}

		if v, ok := detail["regions"]; ok {
			regions := v.([]interface{})
			if len(regions) > 0 {
				source.AwsRegions = flex.ExpandStringList(regions)
			}
		}

		results = append(results, &source)
	}
	return results
}

func expandOrganizationAggregationSource(configured map[string]interface{}) *configservice.OrganizationAggregationSource {
	source := configservice.OrganizationAggregationSource{
		AllAwsRegions: aws.Bool(configured["all_regions"].(bool)),
		RoleArn:       aws.String(configured["role_arn"].(string)),
	}

	if v, ok := configured["regions"]; ok {
		regions := v.([]interface{})
		if len(regions) > 0 {
			source.AwsRegions = flex.ExpandStringList(regions)
		}
	}

	return &source
}

func expandRecordingGroup(group map[string]interface{}) *configservice.RecordingGroup {
	recordingGroup := configservice.RecordingGroup{}

	if v, ok := group["all_supported"]; ok {
		recordingGroup.AllSupported = aws.Bool(v.(bool))
	}

	if v, ok := group["exclusion_by_resource_types"]; ok {
		if len(v.([]interface{})) > 0 {
			recordingGroup.ExclusionByResourceTypes = expandRecordingGroupExclusionByResourceTypes(v.([]interface{}))
		}
	}

	if v, ok := group["include_global_resource_types"]; ok {
		recordingGroup.IncludeGlobalResourceTypes = aws.Bool(v.(bool))
	}

	if v, ok := group["recording_strategy"]; ok {
		if len(v.([]interface{})) > 0 {
			recordingGroup.RecordingStrategy = expandRecordingGroupRecordingStrategy(v.([]interface{}))
		}
	}

	if v, ok := group["resource_types"]; ok {
		recordingGroup.ResourceTypes = flex.ExpandStringSet(v.(*schema.Set))
	}
	return &recordingGroup
}

func expandRecordingGroupExclusionByResourceTypes(configured []interface{}) *configservice.ExclusionByResourceTypes {
	exclusionByResourceTypes := configservice.ExclusionByResourceTypes{}
	exclusion := configured[0].(map[string]interface{})
	if v, ok := exclusion["resource_types"]; ok {
		exclusionByResourceTypes.ResourceTypes = flex.ExpandStringSet(v.(*schema.Set))
	}
	return &exclusionByResourceTypes
}

func expandRecordingGroupRecordingStrategy(configured []interface{}) *configservice.RecordingStrategy {
	recordingStrategy := configservice.RecordingStrategy{}
	strategy := configured[0].(map[string]interface{})
	if v, ok := strategy["use_only"].(string); ok {
		recordingStrategy.UseOnly = aws.String(v)
	}
	return &recordingStrategy
}

func expandRecordingMode(mode map[string]interface{}) *configservice.RecordingMode {
	recordingMode := configservice.RecordingMode{}

	if v, ok := mode["recording_frequency"].(string); ok {
		recordingMode.RecordingFrequency = aws.String(v)
	}

	if v, ok := mode["recording_mode_override"]; ok {
		recordingMode.RecordingModeOverrides = expandRecordingModeRecordingModeOverrides(v.([]interface{}))
	}

	return &recordingMode
}

func expandRecordingModeRecordingModeOverrides(configured []interface{}) []*configservice.RecordingModeOverride {
	var out []*configservice.RecordingModeOverride
	for _, val := range configured {
		m, ok := val.(map[string]interface{})
		if !ok {
			continue
		}

		e := &configservice.RecordingModeOverride{}

		if v, ok := m["description"].(string); ok && v != "" {
			e.Description = aws.String(v)
		}

		if v, ok := m["resource_types"]; ok {
			e.ResourceTypes = flex.ExpandStringSet(v.(*schema.Set))
		}

		if v, ok := m["recording_frequency"].(string); ok && v != "" {
			e.RecordingFrequency = aws.String(v)
		}

		out = append(out, e)
	}
	return out
}

func expandRulesEvaluationModes(in []interface{}) []*configservice.EvaluationModeConfiguration {
	if len(in) == 0 {
		return nil
	}

	var out []*configservice.EvaluationModeConfiguration
	for _, val := range in {
		m, ok := val.(map[string]interface{})
		if !ok {
			continue
		}

		e := &configservice.EvaluationModeConfiguration{}
		if v, ok := m["mode"].(string); ok && v != "" {
			e.Mode = aws.String(v)
		}

		out = append(out, e)
	}

	return out
}

func expandRuleScope(l []interface{}) *configservice.Scope {
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	configured := l[0].(map[string]interface{})
	scope := &configservice.Scope{}

	if v, ok := configured["compliance_resource_id"].(string); ok && v != "" {
		scope.ComplianceResourceId = aws.String(v)
	}
	if v, ok := configured["compliance_resource_types"]; ok {
		l := v.(*schema.Set)
		if l.Len() > 0 {
			scope.ComplianceResourceTypes = flex.ExpandStringSet(l)
		}
	}
	if v, ok := configured["tag_key"].(string); ok && v != "" {
		scope.TagKey = aws.String(v)
	}
	if v, ok := configured["tag_value"].(string); ok && v != "" {
		scope.TagValue = aws.String(v)
	}

	return scope
}

func expandRuleSource(configured []interface{}) *configservice.Source {
	cfg := configured[0].(map[string]interface{})
	source := configservice.Source{
		Owner: aws.String(cfg["owner"].(string)),
	}

	if v, ok := cfg["source_identifier"].(string); ok && v != "" {
		source.SourceIdentifier = aws.String(v)
	}

	if details, ok := cfg["source_detail"]; ok {
		source.SourceDetails = expandRuleSourceDetails(details.(*schema.Set))
	}

	if v, ok := cfg["custom_policy_details"].([]interface{}); ok && len(v) > 0 {
		source.CustomPolicyDetails = expandRuleSourceCustomPolicyDetails(v)
	}

	return &source
}

func expandRuleSourceDetails(configured *schema.Set) []*configservice.SourceDetail {
	var results []*configservice.SourceDetail

	for _, item := range configured.List() {
		detail := item.(map[string]interface{})
		src := configservice.SourceDetail{}

		if msgType, ok := detail["message_type"].(string); ok && msgType != "" {
			src.MessageType = aws.String(msgType)
		}
		if eventSource, ok := detail["event_source"].(string); ok && eventSource != "" {
			src.EventSource = aws.String(eventSource)
		}
		if maxExecFreq, ok := detail["maximum_execution_frequency"].(string); ok && maxExecFreq != "" {
			src.MaximumExecutionFrequency = aws.String(maxExecFreq)
		}

		results = append(results, &src)
	}

	return results
}

func expandRuleSourceCustomPolicyDetails(configured []interface{}) *configservice.CustomPolicyDetails {
	cfg := configured[0].(map[string]interface{})
	source := configservice.CustomPolicyDetails{
		PolicyRuntime:          aws.String(cfg["policy_runtime"].(string)),
		PolicyText:             aws.String(cfg["policy_text"].(string)),
		EnableDebugLogDelivery: aws.Bool(cfg["enable_debug_log_delivery"].(bool)),
	}

	return &source
}

func flattenAccountAggregationSources(sources []*configservice.AccountAggregationSource) []interface{} {
	var result []interface{}

	if len(sources) == 0 {
		return result
	}

	source := sources[0]
	m := make(map[string]interface{})
	m["account_ids"] = flex.FlattenStringList(source.AccountIds)
	m["all_regions"] = aws.BoolValue(source.AllAwsRegions)
	m["regions"] = flex.FlattenStringList(source.AwsRegions)
	result = append(result, m)
	return result
}

func flattenOrganizationAggregationSource(source *configservice.OrganizationAggregationSource) []interface{} {
	var result []interface{}

	if source == nil {
		return result
	}

	m := make(map[string]interface{})
	m["all_regions"] = aws.BoolValue(source.AllAwsRegions)
	m["regions"] = flex.FlattenStringList(source.AwsRegions)
	m["role_arn"] = aws.StringValue(source.RoleArn)
	result = append(result, m)
	return result
}

func flattenRecordingGroup(g *configservice.RecordingGroup) []map[string]interface{} {
	m := make(map[string]interface{}, 1)

	if g.AllSupported != nil {
		m["all_supported"] = aws.BoolValue(g.AllSupported)
	}

	if g.ExclusionByResourceTypes != nil {
		m["exclusion_by_resource_types"] = flattenExclusionByResourceTypes(g.ExclusionByResourceTypes)
	}

	if g.IncludeGlobalResourceTypes != nil {
		m["include_global_resource_types"] = aws.BoolValue(g.IncludeGlobalResourceTypes)
	}

	if g.RecordingStrategy != nil {
		m["recording_strategy"] = flattenRecordingGroupRecordingStrategy(g.RecordingStrategy)
	}

	if g.ResourceTypes != nil && len(g.ResourceTypes) > 0 {
		m["resource_types"] = flex.FlattenStringSet(g.ResourceTypes)
	}

	return []map[string]interface{}{m}
}

func flattenRecordingMode(g *configservice.RecordingMode) []map[string]interface{} {
	m := make(map[string]interface{}, 1)

	if g.RecordingFrequency != nil {
		m["recording_frequency"] = aws.StringValue(g.RecordingFrequency)
	}

	if g.RecordingModeOverrides != nil && len(g.RecordingModeOverrides) > 0 {
		m["recording_mode_override"] = flattenRecordingModeRecordingModeOverrides(g.RecordingModeOverrides)
	}

	return []map[string]interface{}{m}
}

func flattenRecordingModeRecordingModeOverrides(in []*configservice.RecordingModeOverride) []interface{} {
	var out []interface{}
	for _, v := range in {
		m := map[string]interface{}{
			"description":         aws.StringValue(v.Description),
			"recording_frequency": aws.StringValue(v.RecordingFrequency),
		}

		if v.ResourceTypes != nil {
			m["resource_types"] = flex.FlattenStringSet(v.ResourceTypes)
		}

		out = append(out, m)
	}

	return out
}

func flattenExclusionByResourceTypes(exclusionByResourceTypes *configservice.ExclusionByResourceTypes) []interface{} {
	if exclusionByResourceTypes == nil {
		return nil
	}
	m := make(map[string]interface{})
	if exclusionByResourceTypes.ResourceTypes != nil {
		m["resource_types"] = flex.FlattenStringSet(exclusionByResourceTypes.ResourceTypes)
	}

	return []interface{}{m}
}

func flattenRuleEvaluationMode(in []*configservice.EvaluationModeConfiguration) []interface{} {
	if len(in) == 0 {
		return nil
	}

	var out []interface{}
	for _, v := range in {
		m := map[string]interface{}{
			"mode": aws.StringValue(v.Mode),
		}

		out = append(out, m)
	}

	return out
}

func flattenRecordingGroupRecordingStrategy(recordingStrategy *configservice.RecordingStrategy) []interface{} {
	if recordingStrategy == nil {
		return nil
	}
	m := make(map[string]interface{})
	if recordingStrategy.UseOnly != nil {
		m["use_only"] = aws.StringValue(recordingStrategy.UseOnly)
	}

	return []interface{}{m}
}

func flattenRuleScope(scope *configservice.Scope) []interface{} {
	var items []interface{}

	m := make(map[string]interface{})
	if scope.ComplianceResourceId != nil {
		m["compliance_resource_id"] = aws.StringValue(scope.ComplianceResourceId)
	}
	if scope.ComplianceResourceTypes != nil {
		m["compliance_resource_types"] = flex.FlattenStringSet(scope.ComplianceResourceTypes)
	}
	if scope.TagKey != nil {
		m["tag_key"] = aws.StringValue(scope.TagKey)
	}
	if scope.TagValue != nil {
		m["tag_value"] = aws.StringValue(scope.TagValue)
	}

	items = append(items, m)
	return items
}

func flattenRuleSource(source *configservice.Source) []interface{} {
	var result []interface{}
	m := make(map[string]interface{})
	m["owner"] = aws.StringValue(source.Owner)
	m["source_identifier"] = aws.StringValue(source.SourceIdentifier)

	if source.CustomPolicyDetails != nil {
		m["custom_policy_details"] = flattenRuleSourceCustomPolicyDetails(source.CustomPolicyDetails)
	}

	if len(source.SourceDetails) > 0 {
		m["source_detail"] = flattenRuleSourceDetails(source.SourceDetails)
	}

	result = append(result, m)
	return result
}

func flattenRuleSourceCustomPolicyDetails(source *configservice.CustomPolicyDetails) []interface{} {
	var result []interface{}
	m := make(map[string]interface{})
	m["policy_runtime"] = aws.StringValue(source.PolicyRuntime)
	m["policy_text"] = aws.StringValue(source.PolicyText)
	m["enable_debug_log_delivery"] = aws.BoolValue(source.EnableDebugLogDelivery)

	result = append(result, m)
	return result
}

func flattenRuleSourceDetails(details []*configservice.SourceDetail) []interface{} {
	var items []interface{}
	for _, d := range details {
		m := make(map[string]interface{})
		if d.MessageType != nil {
			m["message_type"] = aws.StringValue(d.MessageType)
		}
		if d.EventSource != nil {
			m["event_source"] = aws.StringValue(d.EventSource)
		}
		if d.MaximumExecutionFrequency != nil {
			m["maximum_execution_frequency"] = aws.StringValue(d.MaximumExecutionFrequency)
		}

		items = append(items, m)
	}

	return items
}

func flattenSnapshotDeliveryProperties(p *configservice.ConfigSnapshotDeliveryProperties) []map[string]interface{} {
	m := make(map[string]interface{})

	if p.DeliveryFrequency != nil {
		m["delivery_frequency"] = aws.StringValue(p.DeliveryFrequency)
	}

	return []map[string]interface{}{m}
}
