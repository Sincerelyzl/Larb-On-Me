package consul

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Sincerelyzl/larb-on-me/common/constants"
	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	consul "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(addr, serviceName string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}

func (r *Registry) Register(ctx context.Context, instanceId, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return constants.ErrBadHostPortRegisterFormat
	}
	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return constants.ErrBadPortType
	}

	agentRegistration := &consul.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Address: host,
		Port:    port,
		Check: &consul.AgentServiceCheck{
			CheckID:                        instanceId,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	return r.client.Agent().ServiceRegister(agentRegistration)
}

func (r *Registry) Unregister(ctx context.Context, instanceId, serviceName string) error {
	middleware.LogGlobal.Log.Info("unregistering service", "service", serviceName, "instanceId", instanceId)
	return r.client.Agent().ServiceDeregister(instanceId)
}

func (r *Registry) Discover(ctx context.Context, instanceId string) ([]string, error) {
	services, _, err := r.client.Health().Service(instanceId, "", true, nil)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, service := range services {
		addrs = append(addrs, fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port))
	}
	return addrs, nil
}

func (r *Registry) HealthCheck(instanceId, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceId, "online", consul.HealthPassing)
}
