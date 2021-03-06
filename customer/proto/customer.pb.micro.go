// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: customer/proto/customer.proto

package proto

import (
	_ "cdcon21builddriver/globalProtos"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/struct"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for CustomerSrv service

func NewCustomerSrvEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for CustomerSrv service

type CustomerSrvService interface {
	GetCustomerById(ctx context.Context, in *SearchId, opts ...client.CallOption) (*Customer, error)
	GetCustomers(ctx context.Context, in *SearchParams, opts ...client.CallOption) (*Customers, error)
	CreateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*Response, error)
	UpdateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*Response, error)
	DeleteCustomer(ctx context.Context, in *SearchId, opts ...client.CallOption) (*Response, error)
	BeforeCreateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*ValidationErr, error)
	BeforeUpdateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*ValidationErr, error)
	BeforeDeleteCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*ValidationErr, error)
	AfterCreateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*AfterFuncErr, error)
	AfterUpdateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*AfterFuncErr, error)
	AfterDeleteCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*AfterFuncErr, error)
}

type customerSrvService struct {
	c    client.Client
	name string
}

func NewCustomerSrvService(name string, c client.Client) CustomerSrvService {
	return &customerSrvService{
		c:    c,
		name: name,
	}
}

func (c *customerSrvService) GetCustomerById(ctx context.Context, in *SearchId, opts ...client.CallOption) (*Customer, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.GetCustomerById", in)
	out := new(Customer)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) GetCustomers(ctx context.Context, in *SearchParams, opts ...client.CallOption) (*Customers, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.GetCustomers", in)
	out := new(Customers)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) CreateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.CreateCustomer", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) UpdateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.UpdateCustomer", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) DeleteCustomer(ctx context.Context, in *SearchId, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.DeleteCustomer", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) BeforeCreateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*ValidationErr, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.BeforeCreateCustomer", in)
	out := new(ValidationErr)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) BeforeUpdateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*ValidationErr, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.BeforeUpdateCustomer", in)
	out := new(ValidationErr)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) BeforeDeleteCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*ValidationErr, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.BeforeDeleteCustomer", in)
	out := new(ValidationErr)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) AfterCreateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*AfterFuncErr, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.AfterCreateCustomer", in)
	out := new(AfterFuncErr)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) AfterUpdateCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*AfterFuncErr, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.AfterUpdateCustomer", in)
	out := new(AfterFuncErr)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerSrvService) AfterDeleteCustomer(ctx context.Context, in *Customer, opts ...client.CallOption) (*AfterFuncErr, error) {
	req := c.c.NewRequest(c.name, "CustomerSrv.AfterDeleteCustomer", in)
	out := new(AfterFuncErr)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CustomerSrv service

type CustomerSrvHandler interface {
	GetCustomerById(context.Context, *SearchId, *Customer) error
	GetCustomers(context.Context, *SearchParams, *Customers) error
	CreateCustomer(context.Context, *Customer, *Response) error
	UpdateCustomer(context.Context, *Customer, *Response) error
	DeleteCustomer(context.Context, *SearchId, *Response) error
	BeforeCreateCustomer(context.Context, *Customer, *ValidationErr) error
	BeforeUpdateCustomer(context.Context, *Customer, *ValidationErr) error
	BeforeDeleteCustomer(context.Context, *Customer, *ValidationErr) error
	AfterCreateCustomer(context.Context, *Customer, *AfterFuncErr) error
	AfterUpdateCustomer(context.Context, *Customer, *AfterFuncErr) error
	AfterDeleteCustomer(context.Context, *Customer, *AfterFuncErr) error
}

func RegisterCustomerSrvHandler(s server.Server, hdlr CustomerSrvHandler, opts ...server.HandlerOption) error {
	type customerSrv interface {
		GetCustomerById(ctx context.Context, in *SearchId, out *Customer) error
		GetCustomers(ctx context.Context, in *SearchParams, out *Customers) error
		CreateCustomer(ctx context.Context, in *Customer, out *Response) error
		UpdateCustomer(ctx context.Context, in *Customer, out *Response) error
		DeleteCustomer(ctx context.Context, in *SearchId, out *Response) error
		BeforeCreateCustomer(ctx context.Context, in *Customer, out *ValidationErr) error
		BeforeUpdateCustomer(ctx context.Context, in *Customer, out *ValidationErr) error
		BeforeDeleteCustomer(ctx context.Context, in *Customer, out *ValidationErr) error
		AfterCreateCustomer(ctx context.Context, in *Customer, out *AfterFuncErr) error
		AfterUpdateCustomer(ctx context.Context, in *Customer, out *AfterFuncErr) error
		AfterDeleteCustomer(ctx context.Context, in *Customer, out *AfterFuncErr) error
	}
	type CustomerSrv struct {
		customerSrv
	}
	h := &customerSrvHandler{hdlr}
	return s.Handle(s.NewHandler(&CustomerSrv{h}, opts...))
}

type customerSrvHandler struct {
	CustomerSrvHandler
}

func (h *customerSrvHandler) GetCustomerById(ctx context.Context, in *SearchId, out *Customer) error {
	return h.CustomerSrvHandler.GetCustomerById(ctx, in, out)
}

func (h *customerSrvHandler) GetCustomers(ctx context.Context, in *SearchParams, out *Customers) error {
	return h.CustomerSrvHandler.GetCustomers(ctx, in, out)
}

func (h *customerSrvHandler) CreateCustomer(ctx context.Context, in *Customer, out *Response) error {
	return h.CustomerSrvHandler.CreateCustomer(ctx, in, out)
}

func (h *customerSrvHandler) UpdateCustomer(ctx context.Context, in *Customer, out *Response) error {
	return h.CustomerSrvHandler.UpdateCustomer(ctx, in, out)
}

func (h *customerSrvHandler) DeleteCustomer(ctx context.Context, in *SearchId, out *Response) error {
	return h.CustomerSrvHandler.DeleteCustomer(ctx, in, out)
}

func (h *customerSrvHandler) BeforeCreateCustomer(ctx context.Context, in *Customer, out *ValidationErr) error {
	return h.CustomerSrvHandler.BeforeCreateCustomer(ctx, in, out)
}

func (h *customerSrvHandler) BeforeUpdateCustomer(ctx context.Context, in *Customer, out *ValidationErr) error {
	return h.CustomerSrvHandler.BeforeUpdateCustomer(ctx, in, out)
}

func (h *customerSrvHandler) BeforeDeleteCustomer(ctx context.Context, in *Customer, out *ValidationErr) error {
	return h.CustomerSrvHandler.BeforeDeleteCustomer(ctx, in, out)
}

func (h *customerSrvHandler) AfterCreateCustomer(ctx context.Context, in *Customer, out *AfterFuncErr) error {
	return h.CustomerSrvHandler.AfterCreateCustomer(ctx, in, out)
}

func (h *customerSrvHandler) AfterUpdateCustomer(ctx context.Context, in *Customer, out *AfterFuncErr) error {
	return h.CustomerSrvHandler.AfterUpdateCustomer(ctx, in, out)
}

func (h *customerSrvHandler) AfterDeleteCustomer(ctx context.Context, in *Customer, out *AfterFuncErr) error {
	return h.CustomerSrvHandler.AfterDeleteCustomer(ctx, in, out)
}
