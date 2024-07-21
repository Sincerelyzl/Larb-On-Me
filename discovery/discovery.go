package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/middleware"
)

type Registry interface {
	Register(ctx context.Context, instanceId, serviceName, hostPort string) error
	Unregister(ctx context.Context, instanceId, serviceName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceId, serviceName string) error
}

func GenerateInstaceId(serviceName string) string {
	unique := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	instanceId := fmt.Sprintf("%s-%d", serviceName, unique)
	return instanceId
}

func CreateThreadHealthCheck(ctx context.Context, registry Registry, instanceId, serviceName string) {
	go func() {
		for {
			if err := registry.HealthCheck(instanceId, serviceName); err != nil {
				middleware.LogGlobal.Log.Error("health check", "error", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
