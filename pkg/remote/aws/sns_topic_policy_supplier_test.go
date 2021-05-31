package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"

	"github.com/aws/aws-sdk-go/service/sns"

	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"

	resourceaws "github.com/cloudskiff/driftctl/pkg/resource/aws"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/cloudskiff/driftctl/pkg/parallel"

	awsdeserializer "github.com/cloudskiff/driftctl/pkg/resource/aws/deserializer"

	"github.com/cloudskiff/driftctl/test/goldenfile"
	mocks2 "github.com/cloudskiff/driftctl/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cloudskiff/driftctl/mocks"

	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/cloudskiff/driftctl/test"
)

func TestSNSTopicPolicySupplier_Resources(t *testing.T) {
	cases := []struct {
		test    string
		dirName string
		mocks   func(client *mocks.SNSRepository)
		err     error
	}{
		{
			test:    "no SNS Topic policy",
			dirName: "sns_topic_policy_empty",
			mocks: func(client *mocks.SNSRepository) {
				client.On("ListAllTopics").Return([]*sns.Topic{}, nil)
			},
			err: nil,
		},
		{
			test:    "Multiple SNSTopicPolicy",
			dirName: "sns_topic_policy_multiple",
			mocks: func(client *mocks.SNSRepository) {
				client.On("ListAllTopics").Return([]*sns.Topic{
					{TopicArn: aws.String("arn:aws:sns:us-east-1:526954929923:my-topic-with-policy")},
					{TopicArn: aws.String("arn:aws:sns:us-east-1:526954929923:my-topic-with-policy2")},
				}, nil)
			},
			err: nil,
		},
		{
			test:    "cannot list SNSTopic",
			dirName: "sns_topic_policy_topic_list",
			mocks: func(client *mocks.SNSRepository) {
				client.On("ListAllTopics").Return(nil, awserr.NewRequestFailure(nil, 403, ""))
			},
			err: remoteerror.NewResourceEnumerationErrorWithType(awserr.NewRequestFailure(nil, 403, ""), resourceaws.AwsSnsTopicPolicyResourceType, resourceaws.AwsSnsTopicResourceType),
		},
	}
	for _, c := range cases {
		shouldUpdate := c.dirName == *goldenfile.Update

		providerLibrary := terraform.NewProviderLibrary()
		supplierLibrary := resource.NewSupplierLibrary()

		if shouldUpdate {
			provider, err := InitTestAwsProvider(providerLibrary)
			if err != nil {
				t.Fatal(err)
			}
			supplierLibrary.AddSupplier(NewSNSTopicSupplier(provider, repository.NewSNSClient(provider.session)))
		}

		t.Run(c.test, func(tt *testing.T) {
			fakeClient := mocks.SNSRepository{}
			c.mocks(&fakeClient)
			provider := mocks2.NewMockedGoldenTFProvider(c.dirName, providerLibrary.Provider(terraform.AWS), shouldUpdate)
			SNSTopicPolicyDeserializer := awsdeserializer.NewSNSTopicPolicyDeserializer()
			s := &SNSTopicPolicySupplier{
				provider,
				SNSTopicPolicyDeserializer,
				&fakeClient,
				terraform.NewParallelResourceReader(parallel.NewParallelRunner(context.TODO(), 10)),
			}
			got, err := s.Resources()
			assert.Equal(tt, c.err, err)

			mock.AssertExpectationsForObjects(tt)
			test.CtyTestDiff(got, c.dirName, provider, SNSTopicPolicyDeserializer, shouldUpdate, tt)
		})
	}
}
