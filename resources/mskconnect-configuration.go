package resources

import (
	"context"

	"github.com/aws/aws-sdk-go/service/kafkaconnect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MSKConnectConfigurationResource = "MSKConnectConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     MSKConnectConfigurationResource,
		Scope:    nuke.Account,
		Resource: &MSKConnectConfiguration{},
		Lister:   &MSKConnectConfigurationLister{},
	})
}

type MSKConnectConfigurationLister struct{}

func (l *MSKConnectConfigurationLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	var resources []resource.Resource
	svc := kafkaconnect.New(opts.Session)

	err := svc.ListWorkerConfigurationsPages(
		&kafkaconnect.ListWorkerConfigurationsInput{},
		func(page *kafkaconnect.ListWorkerConfigurationsOutput, lastPage bool) bool {
			for _, configuration := range page.WorkerConfigurations {
				resources = append(resources, &MSKConnectConfiguration{
					svc:  svc,
					Arn:  configuration.WorkerConfigurationArn,
					Name: configuration.Name,
				})
			}
			return !lastPage
		})
	if err != nil {
		return nil, err
	}

	return resources, nil
}

type MSKConnectConfiguration struct {
	svc  *kafkaconnect.KafkaConnect
	Arn  *string `description:"The ARN of the worker configuration"`
	Name *string `description:"The name of the worker configuration"`
}

func (m *MSKConnectConfiguration) Remove(_ context.Context) error {
	_, err := m.svc.DeleteWorkerConfiguration(&kafkaconnect.DeleteWorkerConfigurationInput{
		WorkerConfigurationArn: m.Arn,
	})
	return err
}

func (m *MSKConnectConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(m)
}
