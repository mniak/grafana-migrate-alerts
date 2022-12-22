package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/samber/lo"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"golang.org/x/exp/slices"
)

type TerraformGenerator struct {
	IncludeAPIGroups    bool
	RuleGroupNameSuffix string
	RuleTitleSuffix     string
}

func (tg TerraformGenerator) GenerateFiles(jfolders Folders) (map[string]*hclwrite.File, error) {
	result := make(map[string]*hclwrite.File)

	tfFoldersFile := hclwrite.NewEmptyFile()
	result["folders"] = tfFoldersFile

	for folderName, jgroups := range jfolders {
		tfFile, err := tg.generateFolder(folderName, jgroups, tfFoldersFile)
		if err != nil {
			return result, err
		}
		if tfFile != nil {
			result[folderName] = tfFile
		}
	}
	return result, nil
}

func (tg TerraformGenerator) generateFolder(folderName string, jgroups []RuleGroup, tfFoldersFile *hclwrite.File) (*hclwrite.File, error) {
	result := hclwrite.NewEmptyFile()
	b := result.Body()

	slices.SortStableFunc(jgroups, func(a, b RuleGroup) bool {
		return strings.Compare(a.Name, b.Name) > 0
	})

	var groupCount int
	for _, jgroup := range jgroups {
		tfBlock, err := tg.generateRuleGroup(folderName, jgroup)
		if err != nil {
			return result, err
		}
		if tfBlock != nil {
			b.AppendBlock(tfBlock)
			groupCount++
		}
	}

	if groupCount == 0 {
		return nil, nil
	}
	folderBlock := tfFoldersFile.Body().AppendNewBlock("data", []string{"grafana_folder", folderName})
	folderBlock.Body().SetAttributeValue("title", cty.StringVal(folderName))
	return result, nil
}

func (tg TerraformGenerator) generateRuleGroup(folderName string, jgroup RuleGroup) (*hclwrite.Block, error) {
	jgroup.Rules = lo.Filter[Rule](jgroup.Rules, func(item Rule, index int) bool {
		return tg.IncludeAPIGroups || item.GrafanaAlert.Provenance != "api"
	})

	if len(jgroup.Rules) == 0 {
		return nil, nil
	}

	result := hclwrite.NewBlock("resource", []string{"grafana_rule_group", jgroup.Name})
	tfgroup := result.Body()

	tfgroup.SetAttributeValue("org_id", cty.NumberIntVal(1))
	tfgroup.SetAttributeTraversal("folder_uid", hcl.Traversal{
		hcl.TraverseRoot{Name: "data"},
		hcl.TraverseAttr{Name: "grafana_folder"},
		hcl.TraverseAttr{Name: folderName},
		hcl.TraverseAttr{Name: "uid"},
	})
	tfgroup.SetAttributeValue("name", cty.StringVal(fmt.Sprintf("%s%s", jgroup.Name, tg.RuleGroupNameSuffix)))
	jgroupIntervalDuration, err := time.ParseDuration(jgroup.Interval)
	if err != nil {
		return result, err
	}
	tfgroup.SetAttributeValue("interval_seconds", cty.NumberIntVal(int64(jgroupIntervalDuration.Seconds())))

	slices.SortStableFunc(jgroup.Rules, func(a, b Rule) bool {
		return strings.Compare(a.GrafanaAlert.Title, b.GrafanaAlert.Title) > 0
	})

	for _, jrule := range jgroup.Rules {

		tfBlock, err := tg.generateRule(jrule)
		if err != nil {
			return result, err
		}
		tfgroup.AppendBlock(tfBlock)
	}
	return result, nil
}

func (tg TerraformGenerator) generateRule(jrule Rule) (*hclwrite.Block, error) {
	result := hclwrite.NewBlock("rule", nil)
	tfrule := result.Body()
	tfrule.SetAttributeValue("name", cty.StringVal(fmt.Sprintf("%s%s", jrule.GrafanaAlert.Title, tg.RuleTitleSuffix)))
	tfrule.SetAttributeValue("condition", cty.StringVal(jrule.GrafanaAlert.Condition))
	tfrule.SetAttributeValue("for", cty.StringVal(jrule.For))
	tfrule.SetAttributeValue("no_data_state", cty.StringVal(jrule.GrafanaAlert.NoDataState))
	tfrule.SetAttributeValue("exec_err_state", cty.StringVal(jrule.GrafanaAlert.ExecErrState))

	if len(jrule.Annotations) > 0 {
		annotationsMap := lo.MapValues(jrule.Annotations, func(value, key string) cty.Value {
			return cty.StringVal(value)
		})
		tfrule.SetAttributeValue("annotations", cty.MapVal(annotationsMap))
	}

	if len(jrule.Labels) > 0 {
		labelsMap := lo.MapValues(jrule.Labels, func(value, key string) cty.Value {
			return cty.StringVal(value)
		})
		tfrule.SetAttributeValue("labels", cty.MapVal(labelsMap))
	}

	for _, jdata := range jrule.GrafanaAlert.Data {
		tfBlock, err := tg.generateDataBlock(jdata)
		if err != nil {
			return result, err
		}
		tfrule.AppendBlock(tfBlock)
	}
	return result, nil
}

func (tg TerraformGenerator) generateDataBlock(d Datum) (*hclwrite.Block, error) {
	result := hclwrite.NewBlock("data", nil)
	b := result.Body()

	b.SetAttributeValue("ref_id", cty.StringVal(d.RefID))
	b.SetAttributeValue("query_type", cty.StringVal(d.QueryType))
	b.SetAttributeValue("datasource_uid", cty.StringVal(d.DatasourceUid))

	tfrelativetimerange := b.AppendNewBlock("relative_time_range", nil).Body()
	tfrelativetimerange.SetAttributeValue("from", cty.NumberIntVal(d.RelativeTimeRange.From))
	tfrelativetimerange.SetAttributeValue("to", cty.NumberIntVal(d.RelativeTimeRange.To))

	tfmodel, err := parseCty(d.Model)
	if err != nil {
		return result, err
	}
	b.SetAttributeRaw("model", hclwrite.TokensForFunctionCall("jsonencode", hclwrite.TokensForValue(tfmodel)))
	return result, nil
}

func parseCty(obj any) (cty.Value, error) {
	switch value := obj.(type) {
	case map[string]any:
		if len(value) == 0 {
			return cty.MapValEmpty(cty.String), nil
		}
		newmap := make(map[string]cty.Value)
		for k, v := range value {
			parsed, err := parseCty(v)
			if err != nil {
				return cty.Value{}, err
			}
			newmap[k] = parsed
		}
		return cty.ObjectVal(newmap), nil
	case []any:
		if len(value) == 0 {
			return cty.ListValEmpty(cty.String), nil
		}
		newlist := make([]cty.Value, len(value))
		for i, v := range value {
			parsed, err := parseCty(v)
			if err != nil {
				return cty.Value{}, err
			}
			newlist[i] = parsed
		}
		return cty.ListVal(newlist), nil
	}

	t, err := gocty.ImpliedType(obj)
	if err != nil {
		return cty.Value{}, err
	}

	return gocty.ToCtyValue(obj, t)
}
