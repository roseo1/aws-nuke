package resources

import (
	"context"

	"github.com/aws/aws-sdk-go/service/kafkaconnect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MSKConnectConnectorResource = "MSKConnectConnector"

func init() {
	registry.Register(&registry.Registration{
		Name:     MSKConnectConnectorResource,
		Scope:    nuke.Account,
		Resource: &MSKConnectConnector{},
		Lister:   &MSKConnectConnectorLister{},
	})
}

type MSKConnectConnectorLister struct{}

func (l *MSKConnectConnectorLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	var resources []resource.Resource
	svc := kafkaconnect.New(opts.Session)

	err := svc.ListConnectorsPages(
		&kafkaconnect.ListConnectorsInput{},
		func(page *kafkaconnect.ListConnectorsOutput, lastPage bool) bool {
			for _, configuration := range page.Connectors {
				resources = append(resources, &MSKConnectConnector{
					svc:  svc,
					Arn:  configuration.ConnectorArn,
					Name: configuration.ConnectorName,
				})
			}
			return !lastPage
		})
	if err != nil {
		return nil, err
	}

	return resources, nil
}

type MSKConnectConnector struct {
	svc  *kafkaconnect.KafkaConnect
	Arn  *string `description:"The ARN of the connector"`
	Name *string `description:"The name of the connector"`
}

func (m *MSKConnectConnector) Remove(_ context.Context) error {
	_, err := m.svc.DeleteConnector(&kafkaconnect.DeleteConnectorInput{
		ConnectorArn: m.Arn,
	})
	return err
}

func (m *MSKConnectConnector) Properties() types.Properties {
	return types.NewPropertiesFromStruct(m)
}
