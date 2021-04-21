// GENERATED, DO NOT EDIT THIS FILE
package aws

import (
	"github.com/cloudskiff/driftctl/pkg/resource"
	rescty "github.com/cloudskiff/driftctl/pkg/resource/cty"
	"github.com/zclconf/go-cty/cty"
)

const AwsLambdaFunctionResourceType = "aws_lambda_function"

type AwsLambdaFunction struct {
	Arn                          *string           `cty:"arn" computed:"true"`
	CodeSigningConfigArn         *string           `cty:"code_signing_config_arn"`
	Description                  *string           `cty:"description"`
	Filename                     *string           `cty:"filename" diff:"-"`
	FunctionName                 *string           `cty:"function_name"`
	Handler                      *string           `cty:"handler"`
	Id                           string            `cty:"id" computed:"true"`
	ImageUri                     *string           `cty:"image_uri"`
	InvokeArn                    *string           `cty:"invoke_arn" computed:"true"`
	KmsKeyArn                    *string           `cty:"kms_key_arn"`
	LastModified                 *string           `cty:"last_modified" computed:"true"`
	Layers                       []string          `cty:"layers"`
	MemorySize                   *int              `cty:"memory_size"`
	PackageType                  *string           `cty:"package_type"`
	Publish                      *bool             `cty:"publish" diff:"-"`
	QualifiedArn                 *string           `cty:"qualified_arn" computed:"true"`
	ReservedConcurrentExecutions *int              `cty:"reserved_concurrent_executions"`
	Role                         *string           `cty:"role"`
	Runtime                      *string           `cty:"runtime"`
	S3Bucket                     *string           `cty:"s3_bucket"`
	S3Key                        *string           `cty:"s3_key"`
	S3ObjectVersion              *string           `cty:"s3_object_version"`
	SigningJobArn                *string           `cty:"signing_job_arn" computed:"true"`
	SigningProfileVersionArn     *string           `cty:"signing_profile_version_arn" computed:"true"`
	SourceCodeHash               *string           `cty:"source_code_hash" computed:"true"`
	SourceCodeSize               *int              `cty:"source_code_size" computed:"true"`
	Tags                         map[string]string `cty:"tags"`
	Timeout                      *int              `cty:"timeout"`
	Version                      *string           `cty:"version" computed:"true"`
	DeadLetterConfig             *[]struct {
		TargetArn *string `cty:"target_arn"`
	} `cty:"dead_letter_config"`
	Environment *[]struct {
		Variables map[string]string `cty:"variables"`
	} `cty:"environment"`
	FileSystemConfig *[]struct {
		Arn            *string `cty:"arn"`
		LocalMountPath *string `cty:"local_mount_path"`
	} `cty:"file_system_config"`
	ImageConfig *[]struct {
		Command          []string `cty:"command"`
		EntryPoint       []string `cty:"entry_point"`
		WorkingDirectory *string  `cty:"working_directory"`
	} `cty:"image_config"`
	Timeouts *struct {
		Create *string `cty:"create"`
	} `cty:"timeouts" diff:"-"`
	TracingConfig *[]struct {
		Mode *string `cty:"mode"`
	} `cty:"tracing_config"`
	VpcConfig *[]struct {
		SecurityGroupIds []string `cty:"security_group_ids"`
		SubnetIds        []string `cty:"subnet_ids"`
		VpcId            *string  `cty:"vpc_id" computed:"true"`
	} `cty:"vpc_config"`
	CtyVal *cty.Value `diff:"-"`
}

func (r *AwsLambdaFunction) TerraformId() string {
	return r.Id
}

func (r *AwsLambdaFunction) TerraformType() string {
	return AwsLambdaFunctionResourceType
}

func (r *AwsLambdaFunction) CtyValue() *cty.Value {
	return r.CtyVal
}

func initAwsLambdaFunctionMetaData(resourceSchemaRepository *resource.SchemaRepository) {
	resourceSchemaRepository.SetNormalizeFunc(AwsLambdaFunctionResourceType, func(val *rescty.CtyAttributes) {
		val.SafeDelete([]string{"filename"})
		val.SafeDelete([]string{"publish"})
		val.SafeDelete([]string{"timeouts"})
		val.SafeDelete([]string{"last_modified"})
	})
}
