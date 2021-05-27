package aws

import (
	"github.com/cloudskiff/driftctl/pkg/alerter"
	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/cloudskiff/driftctl/pkg/remote/aws/client"
	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/aws"
	"github.com/cloudskiff/driftctl/pkg/terraform"
)

const RemoteAWSTerraform = "aws+tf"

/**
 * Initialize remote (configure credentials, launch tf providers and start gRPC clients)
 * Required to use Scanner
 */
func Init(alerter *alerter.Alerter, providerLibrary *terraform.ProviderLibrary, supplierLibrary *resource.SupplierLibrary, progress output.Progress, resourceSchemaRepository *resource.SchemaRepository) error {
	provider, err := NewAWSTerraformProvider(progress)
	if err != nil {
		return err
	}
	err = provider.Init()
	if err != nil {
		return err
	}

	repositoryCache := cache.New(100)

	s3Repository := repository.NewS3Repository(client.NewAWSClientFactory(provider.session), repositoryCache)
	ec2repository := repository.NewEC2Repository(provider.session, repositoryCache)
	route53repository := repository.NewRoute53Repository(provider.session)
	lambdaRepository := repository.NewLambdaRepository(provider.session)
	rdsRepository := repository.NewRDSRepository(provider.session)
	sqsRepository := repository.NewSQSClient(provider.session)
	snsRepository := repository.NewSNSClient(provider.session)
	dynamoDBRepository := repository.NewDynamoDBRepository(provider.session)
	cloudfrontRepository := repository.NewCloudfrontClient(provider.session)
	kmsRepository := repository.NewKMSRepository(provider.session)
	ecrRepository := repository.NewECRRepository(provider.session)

	providerLibrary.AddProvider(terraform.AWS, provider)

	supplierLibrary.AddSupplier(NewS3BucketSupplier(provider, s3Repository))
	supplierLibrary.AddSupplier(NewS3BucketAnalyticSupplier(provider, s3Repository))
	supplierLibrary.AddSupplier(NewS3BucketInventorySupplier(provider, s3Repository))
	supplierLibrary.AddSupplier(NewS3BucketMetricSupplier(provider, s3Repository))
	supplierLibrary.AddSupplier(NewS3BucketNotificationSupplier(provider, s3Repository))
	supplierLibrary.AddSupplier(NewS3BucketPolicySupplier(provider, s3Repository))
	supplierLibrary.AddSupplier(NewEC2EipSupplier(provider, ec2repository))
	supplierLibrary.AddSupplier(NewEC2EipAssociationSupplier(provider, ec2repository))
	supplierLibrary.AddSupplier(NewEC2EbsVolumeSupplier(provider, ec2repository))
	supplierLibrary.AddSupplier(NewEC2EbsSnapshotSupplier(provider, ec2repository))
	supplierLibrary.AddSupplier(NewRoute53ZoneSupplier(provider, route53repository))
	supplierLibrary.AddSupplier(NewRoute53RecordSupplier(provider, route53repository))
	supplierLibrary.AddSupplier(NewEC2InstanceSupplier(provider, ec2repository))
	supplierLibrary.AddSupplier(NewEC2AmiSupplier(provider, ec2repository))
	supplierLibrary.AddSupplier(NewEC2KeyPairSupplier(provider, ec2repository))
	supplierLibrary.AddSupplier(NewLambdaFunctionSupplier(provider, lambdaRepository))
	supplierLibrary.AddSupplier(NewDBSubnetGroupSupplier(provider, rdsRepository))
	supplierLibrary.AddSupplier(NewDBInstanceSupplier(provider, rdsRepository))
	supplierLibrary.AddSupplier(NewVPCSecurityGroupSupplier(provider))
	supplierLibrary.AddSupplier(NewIamUserSupplier(provider))
	supplierLibrary.AddSupplier(NewIamUserPolicySupplier(provider))
	supplierLibrary.AddSupplier(NewIamUserPolicyAttachmentSupplier(provider))
	supplierLibrary.AddSupplier(NewIamAccessKeySupplier(provider))
	supplierLibrary.AddSupplier(NewIamRoleSupplier(provider))
	supplierLibrary.AddSupplier(NewIamPolicySupplier(provider))
	supplierLibrary.AddSupplier(NewIamRolePolicySupplier(provider))
	supplierLibrary.AddSupplier(NewIamRolePolicyAttachmentSupplier(provider))
	supplierLibrary.AddSupplier(NewVPCSecurityGroupRuleSupplier(provider))
	supplierLibrary.AddSupplier(NewVPCSupplier(provider))
	supplierLibrary.AddSupplier(NewSubnetSupplier(provider))
	supplierLibrary.AddSupplier(NewRouteTableSupplier(provider))
	supplierLibrary.AddSupplier(NewRouteSupplier(provider))
	supplierLibrary.AddSupplier(NewRouteTableAssociationSupplier(provider))
	supplierLibrary.AddSupplier(NewNatGatewaySupplier(provider))
	supplierLibrary.AddSupplier(NewInternetGatewaySupplier(provider))
	supplierLibrary.AddSupplier(NewSqsQueueSupplier(provider, sqsRepository))
	supplierLibrary.AddSupplier(NewSqsQueuePolicySupplier(provider, sqsRepository))
	supplierLibrary.AddSupplier(NewSNSTopicSupplier(provider, snsRepository))
	supplierLibrary.AddSupplier(NewSNSTopicPolicySupplier(provider, snsRepository))
	supplierLibrary.AddSupplier(NewSNSTopicSubscriptionSupplier(provider, alerter, snsRepository))
	supplierLibrary.AddSupplier(NewDynamoDBTableSupplier(provider, dynamoDBRepository))
	supplierLibrary.AddSupplier(NewRoute53HealthCheckSupplier(provider, route53repository))
	supplierLibrary.AddSupplier(NewCloudfrontDistributionSupplier(provider, cloudfrontRepository))
	supplierLibrary.AddSupplier(NewECRRepositorySupplier(provider, ecrRepository))
	supplierLibrary.AddSupplier(NewKMSKeySupplier(provider, kmsRepository))
	supplierLibrary.AddSupplier(NewKMSAliasSupplier(provider, kmsRepository))
	supplierLibrary.AddSupplier(NewLambdaEventSourceMappingSupplier(provider, lambdaRepository))

	resourceSchemaRepository.Init(provider.Schema())
	aws.InitResourcesMetadata(resourceSchemaRepository)

	return nil
}
