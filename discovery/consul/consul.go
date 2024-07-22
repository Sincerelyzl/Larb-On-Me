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
	Client *consul.Client
}

func NewRegistry(addr, serviceName string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{Client: client}, nil
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

	return r.Client.Agent().ServiceRegister(agentRegistration)
}

func (r *Registry) Unregister(ctx context.Context, instanceId, serviceName string) error {
	middleware.LogGlobal.Log.Info("unregistering service", "service", serviceName, "instanceId", instanceId)
	return r.Client.Agent().ServiceDeregister(instanceId)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	services, _, err := r.Client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, service := range services {
		addrs = append(addrs, fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port))
	}
	return addrs, nil
}

func (r *Registry) Trigger(ctx context.Context, serviceName string) error {
	//@ TODO: Sincerelyzl - implement trigger
	// userServices, err := uc.registry.Discover(ctx, "user-service")
	// if err != nil {
	// 	return nil, err
	// }
	// if len(userServices) == 0 {
	// 	return nil, fmt.Errorf(constants.ErrServiceUnavailable, "user-service")
	// }
	// userService := userServices[0]
	// userServiceClient := resty.New()
	// userServiceClient.SetDebug(true)
	// userServiceClient.SetRetryCount(3)
	// userServiceClient.SetRetryWaitTime(2 * time.Second)
	// userServiceClient.SetHeader(middleware.LOMCookieAuthPrefix, lomToken)
	// body := models.UserAddChatRoomRequest{
	// 	Uuid: uuidV7StringChatRoom,
	// }
	// res, err := userServiceClient.R().SetBody(body).Patch(fmt.Sprintf("http://%s/v1/user/add.chatroom.uuid", userService))
	// if err != nil {
	// 	return nil, err
	// }
	// if !res.IsSuccess() {
	// 	return nil, fmt.Errorf("failed to update user model")
	// }
	return nil
}

func (r *Registry) HealthCheck(instanceId, serviceName string) error {
	return r.Client.Agent().UpdateTTL(instanceId, "online", consul.HealthPassing)
}
