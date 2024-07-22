package consul

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/constants"
	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/go-resty/resty/v2"
	consul "github.com/hashicorp/consul/api"
	"golang.org/x/exp/rand"
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

func (r *Registry) DiscoveOne(ctx context.Context, serviceName string) (domain string, err error) {
	services, err := r.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(services) == 0 {
		return "", fmt.Errorf(constants.ErrServiceUnavailable, serviceName)
	}
	return services[rand.Intn(len(services))], nil
}

// @Param: serviceName - the name of the service to trigger Eg. user-service
// @Usage: Post to another service
// @Param: path - the path to trigger Eg. /v1/user
// @Param: lomToken - the token to authenticate the request
// @Param: body - the body to send to the service (struct || interface{} || []byte)
// @Param: resultParseType - the type to parse the result into
// @Return: *resty.Response - the response from the service
// @Return: error - if the service is unavailable
func (r *Registry) Post(ctx context.Context, serviceName, path string, lomToken *string, body any, resultParseType *interface{}) (*resty.Response, error) {
	domain, err := r.DiscoveOne(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	restyClient := resty.New()
	//restyClient.SetDebug(true)
	restyClient.SetRetryCount(3)
	restyClient.SetRetryWaitTime(2 * time.Second)
	if lomToken != nil {
		restyClient.SetHeader(middleware.LOMCookieAuthPrefix, *lomToken)
	}

	req := restyClient.R()
	req.SetBody(body)
	if resultParseType != nil {
		req.SetResult(resultParseType)
	}

	res, err := req.Post(fmt.Sprintf("http://%s%s", domain, path))
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		return nil, fmt.Errorf("failed to post to service %s", serviceName)
	}
	return res, nil
}

// @Param: serviceName - the name of the service to trigger Eg. user-service
// @Usage: Get to another service
// @Param: path - the path to trigger Eg. /v1/user
// @Param: lomToken - the token to authenticate the request
// @Param: body - the body to send to the service (struct || interface{} || []byte)
// @Param: resultParseType - the type to parse the result into
// @Return: *resty.Response - the response from the service
// @Return: error - if the service is unavailable
func (r *Registry) Get(ctx context.Context, serviceName, path string, lomToken *string, body any, resultParseType *interface{}) (*resty.Response, error) {
	domain, err := r.DiscoveOne(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	restyClient := resty.New()
	//restyClient.SetDebug(true)
	restyClient.SetRetryCount(3)
	restyClient.SetRetryWaitTime(2 * time.Second)
	if lomToken != nil {
		restyClient.SetHeader(middleware.LOMCookieAuthPrefix, *lomToken)
	}

	req := restyClient.R()
	req.SetBody(body)
	if resultParseType != nil {
		req.SetResult(resultParseType)
	}

	res, err := req.Get(fmt.Sprintf("http://%s%s", domain, path))
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		return nil, fmt.Errorf("failed to get to service %s", serviceName)
	}
	return res, nil
}

// @Param: serviceName - the name of the service to trigger Eg. user-service
// @Usage: Patch to another service
// @Param: path - the path to trigger Eg. /v1/user
// @Param: lomToken - the token to authenticate the request
// @Param: body - the body to send to the service (struct || interface{} || []byte)
// @Param: resultParseType - the type to parse the result into
// @Return: *resty.Response - the response from the service
// @Return: error - if the service is unavailable
func (r *Registry) Patch(ctx context.Context, serviceName, path string, lomToken *string, body any, resultParseType *interface{}) (*resty.Response, error) {
	domain, err := r.DiscoveOne(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	restyClient := resty.New()
	restyClient.SetRetryCount(3)
	restyClient.SetRetryWaitTime(2 * time.Second)
	restyClient.SetRetryMaxWaitTime(2 * time.Second)
	restyClient.SetDebug(true)
	if lomToken != nil {
		restyClient.SetHeader(middleware.LOMCookieAuthPrefix, *lomToken)
	}

	req := restyClient.R()
	req.SetBody(body)
	if resultParseType != nil {
		req.SetResult(resultParseType)
	}

	res, err := req.Patch(fmt.Sprintf("http://%s%s", domain, path))
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		if res.Error() != nil {
			return nil, res.Error().(error)
		}
		return nil, fmt.Errorf("failed to patch to service %s", serviceName)
	}
	return res, nil
}

// @Param: serviceName - the name of the service to trigger Eg. user-service
// @Usage: Put to another service
// @Param: path - the path to trigger Eg. /v1/user
// @Param: lomToken - the token to authenticate the request
// @Param: body - the body to send to the service (struct || interface{} || []byte)
// @Param: resultParseType - the type to parse the result into
// @Return: *resty.Response - the response from the service
// @Return: error - if the service is unavailable
func (r *Registry) Put(ctx context.Context, serviceName, path string, lomToken *string, body any, resultParseType *interface{}) (*resty.Response, error) {
	domain, err := r.DiscoveOne(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	restyClient := resty.New()
	//restyClient.SetDebug(true)
	restyClient.SetRetryCount(3)
	restyClient.SetRetryWaitTime(2 * time.Second)
	if lomToken != nil {
		restyClient.SetHeader(middleware.LOMCookieAuthPrefix, *lomToken)
	}

	req := restyClient.R()
	req.SetBody(body)
	if resultParseType != nil {
		req.SetResult(resultParseType)
	}

	res, err := req.Put(fmt.Sprintf("http://%s%s", domain, path))
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		return nil, fmt.Errorf("failed to put to service %s", serviceName)
	}
	return res, nil
}

// @Param: serviceName - the name of the service to trigger Eg. user-service
// @Usage: Delete to another service
// @Param: path - the path to trigger Eg. /v1/user
// @Param: lomToken - the token to authenticate the request
// @Param: body - the body to send to the service (struct || interface{} || []byte)
// @Param: resultParseType - the type to parse the result into
// @Return: *resty.Response - the response from the service
// @Return: error - if the service is unavailable
func (r *Registry) Delete(ctx context.Context, serviceName, path string, lomToken *string, body any, resultParseType *interface{}) (*resty.Response, error) {
	domain, err := r.DiscoveOne(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	restyClient := resty.New()
	//restyClient.SetDebug(true)
	restyClient.SetRetryCount(3)
	restyClient.SetRetryWaitTime(2 * time.Second)
	if lomToken != nil {
		restyClient.SetHeader(middleware.LOMCookieAuthPrefix, *lomToken)
	}

	req := restyClient.R()
	req.SetBody(body)
	if resultParseType != nil {
		req.SetResult(resultParseType)
	}

	res, err := req.Delete(fmt.Sprintf("http://%s%s", domain, path))
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		return nil, fmt.Errorf("failed to delete to service %s", serviceName)
	}
	return res, nil
}

func (r *Registry) HealthCheck(instanceId, serviceName string) error {
	return r.Client.Agent().UpdateTTL(instanceId, "online", consul.HealthPassing)
}
