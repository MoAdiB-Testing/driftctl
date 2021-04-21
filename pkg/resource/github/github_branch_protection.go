// GENERATED, DO NOT EDIT THIS FILE
package github

import (
	"github.com/cloudskiff/driftctl/pkg/resource"
	rescty "github.com/cloudskiff/driftctl/pkg/resource/cty"
	"github.com/zclconf/go-cty/cty"
)

const GithubBranchProtectionResourceType = "github_branch_protection"

type GithubBranchProtection struct {
	AllowsDeletions            *bool     `cty:"allows_deletions"`
	AllowsForcePushes          *bool     `cty:"allows_force_pushes"`
	EnforceAdmins              *bool     `cty:"enforce_admins"`
	Id                         string    `cty:"id" computed:"true"`
	Pattern                    *string   `cty:"pattern"`
	PushRestrictions           *[]string `cty:"push_restrictions"`
	RepositoryId               *string   `cty:"repository_id" diff:"-"` // Terraform provider is always returning nil
	RequireSignedCommits       *bool     `cty:"require_signed_commits"`
	RequiredPullRequestReviews *[]struct {
		DismissStaleReviews          *bool     `cty:"dismiss_stale_reviews"`
		DismissalRestrictions        *[]string `cty:"dismissal_restrictions"`
		RequireCodeOwnerReviews      *bool     `cty:"require_code_owner_reviews"`
		RequiredApprovingReviewCount *int      `cty:"required_approving_review_count"`
	} `cty:"required_pull_request_reviews"`
	RequiredStatusChecks *[]struct {
		Contexts *[]string `cty:"contexts"`
		Strict   *bool     `cty:"strict"`
	} `cty:"required_status_checks"`
	CtyVal *cty.Value `diff:"-"`
}

func (r *GithubBranchProtection) TerraformId() string {
	return r.Id
}

func (r *GithubBranchProtection) TerraformType() string {
	return GithubBranchProtectionResourceType
}

func (r *GithubBranchProtection) CtyValue() *cty.Value {
	return r.CtyVal
}

func initGithubBranchProtectionMetadata(resourceSchemaRepository *resource.SchemaRepository) {
	resourceSchemaRepository.SetNormalizeFunc(GithubBranchProtectionResourceType, func(val *rescty.CtyAttributes) {
		val.SafeDelete([]string{"repository_id"})
	})
}
