package aws

import "github.com/cloudskiff/driftctl/pkg/resource"

func InitResourcesMetadata(resourceSchemaRepository *resource.SchemaRepository) {
	initAwsAmiMetaData(resourceSchemaRepository)
	initAwsCloudfrontDistributionMetaData(resourceSchemaRepository)
	initAwsDbInstanceMetaData(resourceSchemaRepository)
	initAwsDbSubnetGroupMetaData(resourceSchemaRepository)
	initAwsDefaultRouteTableMetaData(resourceSchemaRepository)
	initAwsDefaultSecurityGroupMetaData(resourceSchemaRepository)
	initAwsDefaultSubnetMetaData(resourceSchemaRepository)
	initAwsDefaultVpcMetaData(resourceSchemaRepository)
	initAwsDynamodbTableMetaData(resourceSchemaRepository)
	initAwsEbsSnapshotMetaData(resourceSchemaRepository)
	initAwsEbsVolumeMetaData(resourceSchemaRepository)
	initAwsEcrRepositoryMetaData(resourceSchemaRepository)
	initAwsEipMetaData(resourceSchemaRepository)
	initAwsEipAssociationMetaData(resourceSchemaRepository)
	initAwsIamAccessKeyMetaData(resourceSchemaRepository)
	initAwsIamPolicyMetaData(resourceSchemaRepository)
	initAwsIamPolicyAttachmentMetaData(resourceSchemaRepository)
	initAwsIamRoleMetaData(resourceSchemaRepository)
	initAwsIamRolePolicyMetaData(resourceSchemaRepository)
	initAwsIamRolePolicyAttachmentMetaData(resourceSchemaRepository)
	initAwsIamUserMetaData(resourceSchemaRepository)
	initAwsIamUserPolicyMetaData(resourceSchemaRepository)
	initAwsIamUserPolicyAttachmentMetaData(resourceSchemaRepository)
	initAwsInstanceMetaData(resourceSchemaRepository)
	initAwsInternetGatewayMetaData(resourceSchemaRepository)
	initAwsKeyPairMetaData(resourceSchemaRepository)
	initAwsKmsAliasMetaData(resourceSchemaRepository)
	initAwsKmsKeyMetaData(resourceSchemaRepository)
	initAwsLambdaEventSourceMappingMetaData(resourceSchemaRepository)
	initAwsLambdaFunctionMetaData(resourceSchemaRepository)
	initAwsRouteMetaData(resourceSchemaRepository)
	initAwsRoute53RecordMetaData(resourceSchemaRepository)
	initAwsRoute53ZoneMetaData(resourceSchemaRepository)
	initAwsRouteTableMetaData(resourceSchemaRepository)
	initAwsS3BucketMetaData(resourceSchemaRepository)
	initAwsS3BucketNotificationMetaData(resourceSchemaRepository)
	initAwsS3BucketPolicyMetaData(resourceSchemaRepository)
	initAwsSecurityGroupMetaData(resourceSchemaRepository)
	initAwsSecurityGroupRuleMetaData(resourceSchemaRepository)
	initAwsSnsTopicMetaData(resourceSchemaRepository)
	initAwsSnsTopicPolicyMetaData(resourceSchemaRepository)
	initAwsSnsTopicSubscriptionMetaData(resourceSchemaRepository)
	initAwsSqsQueuePolicyMetaData(resourceSchemaRepository)
	initAwsSubnetMetaData(resourceSchemaRepository)
}
