package resources

import (
	"context"

	"github.com/aws/aws-sdk-go/service/kafkaconnect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MSKConnectPluginResource = "MSKConnectPlugin"

func init() {
	registry.Register(&registry.Registration{
		Name:     MSKConnectPluginResource,
		Scope:    nuke.Account,
		Resource: &MSKConnectPlugin{},
		Lister:   &MSKConnectPluginLister{},
	})
}

type MSKConnectPluginLister struct{}

func (l *MSKConnectPluginLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	var resources []resource.Resource
	svc := kafkaconnect.New(opts.Session)

	err := svc.ListCustomPluginsPages(
		&kafkaconnect.ListCustomPluginsInput{},
		func(page *kafkaconnect.ListCustomPluginsOutput, lastPage bool) bool {
			for _, configuration := range page.CustomPlugins {
				resources = append(resources, &MSKConnectPlugin{
					svc:  svc,
					Arn:  configuration.CustomPluginArn,
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

type MSKConnectPlugin struct {
	svc  *kafkaconnect.KafkaConnect
	Arn  *string `description:"The ARN of the plugin"`
	Name *string `description:"The name of the plugin"`
}

func (m *MSKConnectPlugin) Remove(_ context.Context) error {
	_, err := m.svc.DeleteCustomPlugin(&kafkaconnect.DeleteCustomPluginInput{
		CustomPluginArn: m.Arn,
	})
	return err
}

func (m *MSKConnectPlugin) Properties() types.Properties {
	return types.NewPropertiesFromStruct(m)
}
