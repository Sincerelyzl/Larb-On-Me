package discovery

import "context"

type Registry interface {
	Register(ctx context.Context, instanceId, serviceName, hostPort string) error
	Unregister(ctx context.Context, instanceId, serviceName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceId, serviceName string) error
}
