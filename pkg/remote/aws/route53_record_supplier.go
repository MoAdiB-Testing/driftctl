package aws

import (
	"strings"

	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"

	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"
	awsdeserializer "github.com/cloudskiff/driftctl/pkg/resource/aws/deserializer"

	"github.com/cloudskiff/driftctl/pkg/remote/deserializer"

	"github.com/cloudskiff/driftctl/pkg/resource"
	resourceaws "github.com/cloudskiff/driftctl/pkg/resource/aws"
	"github.com/cloudskiff/driftctl/pkg/terraform"

	"github.com/zclconf/go-cty/cty"
)

type Route53RecordSupplier struct {
	reader       terraform.ResourceReader
	deserializer deserializer.CTYDeserializer
	client       repository.Route53Repository
	runner       *terraform.ParallelResourceReader
}

func NewRoute53RecordSupplier(provider *AWSTerraformProvider, repo repository.Route53Repository) *Route53RecordSupplier {
	return &Route53RecordSupplier{
		provider,
		awsdeserializer.NewRoute53RecordDeserializer(),
		repo,
		terraform.NewParallelResourceReader(provider.Runner().SubRunner())}
}

func (s *Route53RecordSupplier) Resources() ([]resource.Resource, error) {

	zones, err := s.listZones()
	if err != nil {
		return nil, remoteerror.NewResourceEnumerationErrorWithType(err, resourceaws.AwsRoute53RecordResourceType, resourceaws.AwsRoute53ZoneResourceType)
	}

	for _, zone := range zones {
		if err := s.listRecordsForZone(zone[0], zone[1]); err != nil {
			return nil, remoteerror.NewResourceEnumerationError(err, resourceaws.AwsRoute53RecordResourceType)
		}
	}

	results, err := s.runner.Wait()
	if err != nil {
		return nil, err
	}
	return s.deserializer.Deserialize(results)
}

func (s *Route53RecordSupplier) listZones() ([][2]string, error) {
	results := make([][2]string, 0)
	zones, err := s.client.ListAllZones()
	if err != nil {
		return nil, err
	}

	for _, hostedZone := range zones {
		results = append(results, [2]string{strings.TrimPrefix(*hostedZone.Id, "/hostedzone/"), *hostedZone.Name})
	}

	return results, nil
}

func (s *Route53RecordSupplier) listRecordsForZone(zoneId string, zoneName string) error {

	records, err := s.client.ListRecordsForZone(zoneId)

	if err != nil {
		return err
	}

	for _, raw := range records {
		rawType := *raw.Type
		rawName := *raw.Name
		rawSetIdentifier := raw.SetIdentifier
		s.runner.Run(func() (cty.Value, error) {
			vars := []string{
				zoneId,
				strings.ToLower(strings.TrimSuffix(rawName, ".")),
				rawType,
			}
			if rawSetIdentifier != nil {
				vars = append(vars, *rawSetIdentifier)
			}

			record, err := s.reader.ReadResource(
				terraform.ReadResourceArgs{
					Ty: resourceaws.AwsRoute53RecordResourceType,
					ID: strings.Join(vars, "_"),
				},
			)
			if err != nil {
				return cty.NilVal, err
			}

			return *record, nil
		})

	}
	return nil
}
