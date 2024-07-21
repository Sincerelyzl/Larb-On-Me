package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(ctx context.Context, instanceId, serviceName, hostPort string) error
	Unregister(ctx context.Context, instanceId, serviceName string) error
	Discover(ctx context.Context, instanceId string) ([]string, error)
	HealthCheck(instanceId, serviceName string) error
}

func GenerateInstaceId(serviceName string) string {
	unique := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	instanceId := fmt.Sprintf("%s-%d", serviceName, unique)
	return instanceId
}
