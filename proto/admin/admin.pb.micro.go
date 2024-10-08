// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/admin/admin.proto

package go_micro_service_adminUser

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// Api Endpoints for AdminUser service

func NewAdminUserEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for AdminUser service

type AdminUserService interface {
	Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error)
	Reg(ctx context.Context, in *RegRequest, opts ...client.CallOption) (*RegResponse, error)
	UserList(ctx context.Context, in *UserListReq, opts ...client.CallOption) (*UserListRep, error)
	UserEdit(ctx context.Context, in *EditRequest, opts ...client.CallOption) (*EditResponse, error)
	UserDel(ctx context.Context, in *DelRequest, opts ...client.CallOption) (*DelResponse, error)
	MenuList(ctx context.Context, in *MenuReq, opts ...client.CallOption) (*MenuListRep, error)
	MenuShowTree(ctx context.Context, in *MenuReq, opts ...client.CallOption) (*MenuTreeRep, error)
	MenuTree(ctx context.Context, in *MenuReq, opts ...client.CallOption) (*MenuTreeRep, error)
	MenuAdd(ctx context.Context, in *MenuAddReq, opts ...client.CallOption) (*MenuRep, error)
	MenuEdit(ctx context.Context, in *MenuEditReq, opts ...client.CallOption) (*MenuRep, error)
	MenuDel(ctx context.Context, in *MenuDelReq, opts ...client.CallOption) (*MenuRep, error)
	MenuListByIds(ctx context.Context, in *IdsReq, opts ...client.CallOption) (*MenuListRep, error)
	RoleList(ctx context.Context, in *RoleReq, opts ...client.CallOption) (*RoleListRep, error)
	RoleAdd(ctx context.Context, in *RoleAddReq, opts ...client.CallOption) (*RoleRep, error)
	RoleEdit(ctx context.Context, in *RoleEditReq, opts ...client.CallOption) (*RoleRep, error)
	RoleDel(ctx context.Context, in *RoleDelReq, opts ...client.CallOption) (*RoleRep, error)
	FindMenuIdsByRoleIds(ctx context.Context, in *IdsReq, opts ...client.CallOption) (*MenuIdsRep, error)
}

type adminUserService struct {
	c    client.Client
	name string
}

func NewAdminUserService(name string, c client.Client) AdminUserService {
	return &adminUserService{
		c:    c,
		name: name,
	}
}

func (c *adminUserService) Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	req := c.c.NewRequest(c.name, "AdminUser.Login", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) Reg(ctx context.Context, in *RegRequest, opts ...client.CallOption) (*RegResponse, error) {
	req := c.c.NewRequest(c.name, "AdminUser.Reg", in)
	out := new(RegResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) UserList(ctx context.Context, in *UserListReq, opts ...client.CallOption) (*UserListRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.UserList", in)
	out := new(UserListRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) UserEdit(ctx context.Context, in *EditRequest, opts ...client.CallOption) (*EditResponse, error) {
	req := c.c.NewRequest(c.name, "AdminUser.UserEdit", in)
	out := new(EditResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) UserDel(ctx context.Context, in *DelRequest, opts ...client.CallOption) (*DelResponse, error) {
	req := c.c.NewRequest(c.name, "AdminUser.UserDel", in)
	out := new(DelResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) MenuList(ctx context.Context, in *MenuReq, opts ...client.CallOption) (*MenuListRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.MenuList", in)
	out := new(MenuListRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) MenuShowTree(ctx context.Context, in *MenuReq, opts ...client.CallOption) (*MenuTreeRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.MenuShowTree", in)
	out := new(MenuTreeRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) MenuTree(ctx context.Context, in *MenuReq, opts ...client.CallOption) (*MenuTreeRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.MenuTree", in)
	out := new(MenuTreeRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) MenuAdd(ctx context.Context, in *MenuAddReq, opts ...client.CallOption) (*MenuRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.MenuAdd", in)
	out := new(MenuRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) MenuEdit(ctx context.Context, in *MenuEditReq, opts ...client.CallOption) (*MenuRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.MenuEdit", in)
	out := new(MenuRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) MenuDel(ctx context.Context, in *MenuDelReq, opts ...client.CallOption) (*MenuRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.MenuDel", in)
	out := new(MenuRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) MenuListByIds(ctx context.Context, in *IdsReq, opts ...client.CallOption) (*MenuListRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.MenuListByIds", in)
	out := new(MenuListRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) RoleList(ctx context.Context, in *RoleReq, opts ...client.CallOption) (*RoleListRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.RoleList", in)
	out := new(RoleListRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) RoleAdd(ctx context.Context, in *RoleAddReq, opts ...client.CallOption) (*RoleRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.RoleAdd", in)
	out := new(RoleRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) RoleEdit(ctx context.Context, in *RoleEditReq, opts ...client.CallOption) (*RoleRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.RoleEdit", in)
	out := new(RoleRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) RoleDel(ctx context.Context, in *RoleDelReq, opts ...client.CallOption) (*RoleRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.RoleDel", in)
	out := new(RoleRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminUserService) FindMenuIdsByRoleIds(ctx context.Context, in *IdsReq, opts ...client.CallOption) (*MenuIdsRep, error) {
	req := c.c.NewRequest(c.name, "AdminUser.FindMenuIdsByRoleIds", in)
	out := new(MenuIdsRep)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AdminUser service

type AdminUserHandler interface {
	Login(context.Context, *LoginRequest, *LoginResponse) error
	Reg(context.Context, *RegRequest, *RegResponse) error
	UserList(context.Context, *UserListReq, *UserListRep) error
	UserEdit(context.Context, *EditRequest, *EditResponse) error
	UserDel(context.Context, *DelRequest, *DelResponse) error
	MenuList(context.Context, *MenuReq, *MenuListRep) error
	MenuShowTree(context.Context, *MenuReq, *MenuTreeRep) error
	MenuTree(context.Context, *MenuReq, *MenuTreeRep) error
	MenuAdd(context.Context, *MenuAddReq, *MenuRep) error
	MenuEdit(context.Context, *MenuEditReq, *MenuRep) error
	MenuDel(context.Context, *MenuDelReq, *MenuRep) error
	MenuListByIds(context.Context, *IdsReq, *MenuListRep) error
	RoleList(context.Context, *RoleReq, *RoleListRep) error
	RoleAdd(context.Context, *RoleAddReq, *RoleRep) error
	RoleEdit(context.Context, *RoleEditReq, *RoleRep) error
	RoleDel(context.Context, *RoleDelReq, *RoleRep) error
	FindMenuIdsByRoleIds(context.Context, *IdsReq, *MenuIdsRep) error
}

func RegisterAdminUserHandler(s server.Server, hdlr AdminUserHandler, opts ...server.HandlerOption) error {
	type adminUser interface {
		Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error
		Reg(ctx context.Context, in *RegRequest, out *RegResponse) error
		UserList(ctx context.Context, in *UserListReq, out *UserListRep) error
		UserEdit(ctx context.Context, in *EditRequest, out *EditResponse) error
		UserDel(ctx context.Context, in *DelRequest, out *DelResponse) error
		MenuList(ctx context.Context, in *MenuReq, out *MenuListRep) error
		MenuShowTree(ctx context.Context, in *MenuReq, out *MenuTreeRep) error
		MenuTree(ctx context.Context, in *MenuReq, out *MenuTreeRep) error
		MenuAdd(ctx context.Context, in *MenuAddReq, out *MenuRep) error
		MenuEdit(ctx context.Context, in *MenuEditReq, out *MenuRep) error
		MenuDel(ctx context.Context, in *MenuDelReq, out *MenuRep) error
		MenuListByIds(ctx context.Context, in *IdsReq, out *MenuListRep) error
		RoleList(ctx context.Context, in *RoleReq, out *RoleListRep) error
		RoleAdd(ctx context.Context, in *RoleAddReq, out *RoleRep) error
		RoleEdit(ctx context.Context, in *RoleEditReq, out *RoleRep) error
		RoleDel(ctx context.Context, in *RoleDelReq, out *RoleRep) error
		FindMenuIdsByRoleIds(ctx context.Context, in *IdsReq, out *MenuIdsRep) error
	}
	type AdminUser struct {
		adminUser
	}
	h := &adminUserHandler{hdlr}
	return s.Handle(s.NewHandler(&AdminUser{h}, opts...))
}

type adminUserHandler struct {
	AdminUserHandler
}

func (h *adminUserHandler) Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error {
	return h.AdminUserHandler.Login(ctx, in, out)
}

func (h *adminUserHandler) Reg(ctx context.Context, in *RegRequest, out *RegResponse) error {
	return h.AdminUserHandler.Reg(ctx, in, out)
}

func (h *adminUserHandler) UserList(ctx context.Context, in *UserListReq, out *UserListRep) error {
	return h.AdminUserHandler.UserList(ctx, in, out)
}

func (h *adminUserHandler) UserEdit(ctx context.Context, in *EditRequest, out *EditResponse) error {
	return h.AdminUserHandler.UserEdit(ctx, in, out)
}

func (h *adminUserHandler) UserDel(ctx context.Context, in *DelRequest, out *DelResponse) error {
	return h.AdminUserHandler.UserDel(ctx, in, out)
}

func (h *adminUserHandler) MenuList(ctx context.Context, in *MenuReq, out *MenuListRep) error {
	return h.AdminUserHandler.MenuList(ctx, in, out)
}

func (h *adminUserHandler) MenuShowTree(ctx context.Context, in *MenuReq, out *MenuTreeRep) error {
	return h.AdminUserHandler.MenuShowTree(ctx, in, out)
}

func (h *adminUserHandler) MenuTree(ctx context.Context, in *MenuReq, out *MenuTreeRep) error {
	return h.AdminUserHandler.MenuTree(ctx, in, out)
}

func (h *adminUserHandler) MenuAdd(ctx context.Context, in *MenuAddReq, out *MenuRep) error {
	return h.AdminUserHandler.MenuAdd(ctx, in, out)
}

func (h *adminUserHandler) MenuEdit(ctx context.Context, in *MenuEditReq, out *MenuRep) error {
	return h.AdminUserHandler.MenuEdit(ctx, in, out)
}

func (h *adminUserHandler) MenuDel(ctx context.Context, in *MenuDelReq, out *MenuRep) error {
	return h.AdminUserHandler.MenuDel(ctx, in, out)
}

func (h *adminUserHandler) MenuListByIds(ctx context.Context, in *IdsReq, out *MenuListRep) error {
	return h.AdminUserHandler.MenuListByIds(ctx, in, out)
}

func (h *adminUserHandler) RoleList(ctx context.Context, in *RoleReq, out *RoleListRep) error {
	return h.AdminUserHandler.RoleList(ctx, in, out)
}

func (h *adminUserHandler) RoleAdd(ctx context.Context, in *RoleAddReq, out *RoleRep) error {
	return h.AdminUserHandler.RoleAdd(ctx, in, out)
}

func (h *adminUserHandler) RoleEdit(ctx context.Context, in *RoleEditReq, out *RoleRep) error {
	return h.AdminUserHandler.RoleEdit(ctx, in, out)
}

func (h *adminUserHandler) RoleDel(ctx context.Context, in *RoleDelReq, out *RoleRep) error {
	return h.AdminUserHandler.RoleDel(ctx, in, out)
}

func (h *adminUserHandler) FindMenuIdsByRoleIds(ctx context.Context, in *IdsReq, out *MenuIdsRep) error {
	return h.AdminUserHandler.FindMenuIdsByRoleIds(ctx, in, out)
}
