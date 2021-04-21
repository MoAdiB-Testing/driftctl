package analyser

import (
	"sort"
	"strings"

	resourceaws "github.com/cloudskiff/driftctl/pkg/resource/aws"
	"github.com/cloudskiff/driftctl/pkg/resource/cty"
	"github.com/r3labs/diff/v2"

	"github.com/cloudskiff/driftctl/pkg/alerter"
	"github.com/cloudskiff/driftctl/pkg/resource"
)

type UnmanagedSecurityGroupRulesAlert struct{}

func newUnmanagedSecurityGroupRulesAlert() *UnmanagedSecurityGroupRulesAlert {
	return &UnmanagedSecurityGroupRulesAlert{}
}

func (u *UnmanagedSecurityGroupRulesAlert) Message() string {
	return "You have unmanaged security group rules that could be false positives, find out more at https://github.com/cloudskiff/driftctl/blob/main/doc/LIMITATIONS.md#terraform-resources"
}

func (u *UnmanagedSecurityGroupRulesAlert) ShouldIgnoreResource() bool {
	return false
}

type ComputedDiffAlert struct{}

func NewComputedDiffAlert() *ComputedDiffAlert {
	return &ComputedDiffAlert{}
}

func (c *ComputedDiffAlert) Message() string {
	return "You have diffs on computed fields, check the documentation for potential false positive drifts"
}

func (c *ComputedDiffAlert) ShouldIgnoreResource() bool {
	return false
}

type Analyzer struct {
	alerter                  *alerter.Alerter
	resourceSchemaRepository *resource.SchemaRepository
}

type Filter interface {
	IsResourceIgnored(res resource.Resource) bool
	IsFieldIgnored(res resource.Resource, path []string) bool
}

func NewAnalyzer(alerter *alerter.Alerter, resourceSchemaRepository *resource.SchemaRepository) Analyzer {
	return Analyzer{alerter, resourceSchemaRepository}
}

func (a Analyzer) Analyze(remoteResources, resourcesFromState []resource.Resource, filter Filter) (Analysis, error) {
	analysis := Analysis{}

	// Iterate on remote resources and filter ignored resources
	filteredRemoteResource := make([]resource.Resource, 0, len(remoteResources))
	for _, remoteRes := range remoteResources {
		if filter.IsResourceIgnored(remoteRes) || a.alerter.IsResourceIgnored(remoteRes) {
			continue
		}
		filteredRemoteResource = append(filteredRemoteResource, remoteRes)
	}

	haveComputedDiff := false
	for _, stateRes := range resourcesFromState {
		i, remoteRes, found := findCorrespondingRes(filteredRemoteResource, stateRes)

		if filter.IsResourceIgnored(stateRes) || a.alerter.IsResourceIgnored(stateRes) {
			continue
		}

		if !found {
			analysis.AddDeleted(stateRes)
			continue
		}

		// Remove managed resources, so it will remain only unmanaged ones
		filteredRemoteResource = removeResourceByIndex(i, filteredRemoteResource)
		analysis.AddManaged(stateRes)

		var stateAttrs map[string]interface{}
		if res, ok := stateRes.(*resource.AbstractResource); ok {
			stateAttrs = res.Attrs
		} else {
			stateAttrs = cty.ToCtyAttributes(stateRes.CtyValue()).Attrs
		}

		var remoteAttrs map[string]interface{}
		if res, ok := remoteRes.(*resource.AbstractResource); ok {
			remoteAttrs = res.Attrs
		} else {
			remoteAttrs = cty.ToCtyAttributes(stateRes.CtyValue()).Attrs
		}

		delta, _ := diff.Diff(stateAttrs, remoteAttrs)
		if len(delta) > 0 {
			sort.Slice(delta, func(i, j int) bool {
				return strings.Join(delta[i].Path, ".") < strings.Join(delta[j].Path, ".") || delta[i].Type < delta[j].Type
			})
			changelog := make([]Change, 0, len(delta))
			for _, change := range delta {
				if filter.IsFieldIgnored(stateRes, change.Path) {
					continue
				}
				c := Change{Change: change}
				resSchema, exist := a.resourceSchemaRepository.GetSchema(stateRes.TerraformType())
				if exist {
					c.Computed = resSchema.IsComputedField(c.Path)
					c.JsonString = resSchema.IsJsonStringField(c.Path)
				}
				if c.Computed {
					haveComputedDiff = true
				}
				changelog = append(changelog, c)
			}
			if len(changelog) > 0 {
				analysis.AddDifference(Difference{
					Res:       stateRes,
					Changelog: changelog,
				})
			}
		}
	}

	if a.hasUnmanagedSecurityGroupRules(filteredRemoteResource) {
		a.alerter.SendAlert("", newUnmanagedSecurityGroupRulesAlert())
	}

	if haveComputedDiff {
		a.alerter.SendAlert("", NewComputedDiffAlert())
	}

	// Add remaining unmanaged resources
	analysis.AddUnmanaged(filteredRemoteResource...)

	// Sort resources by Terraform Id
	// The purpose is to have a predictable output
	analysis.SortResources()

	analysis.SetAlerts(a.alerter.Retrieve())

	return analysis, nil
}

func findCorrespondingRes(resources []resource.Resource, res resource.Resource) (int, resource.Resource, bool) {
	for i, r := range resources {
		if resource.IsSameResource(res, r) {
			return i, r, true
		}
	}
	return -1, nil, false
}

func removeResourceByIndex(i int, resources []resource.Resource) []resource.Resource {
	if i == len(resources)-1 {
		return resources[:len(resources)-1]
	}
	return append(resources[:i], resources[i+1:]...)
}

// hasUnmanagedSecurityGroupRules returns true if we find at least one unmanaged
// security group rule
func (a Analyzer) hasUnmanagedSecurityGroupRules(unmanagedResources []resource.Resource) bool {
	for _, res := range unmanagedResources {
		if res.TerraformType() == resourceaws.AwsSecurityGroupRuleResourceType {
			return true
		}
	}
	return false
}
