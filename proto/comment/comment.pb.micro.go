// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/comment/comment.proto

package go_micro_service_comment

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

// Api Endpoints for CommentSrv service

func NewCommentSrvEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for CommentSrv service

type CommentSrvService interface {
	AddComment(ctx context.Context, in *AddCommentRequest, opts ...client.CallOption) (*CommonResponse, error)
	GetComments(ctx context.Context, in *CommentsRequest, opts ...client.CallOption) (*CommentResponse, error)
	DianZan(ctx context.Context, in *DianZanRequest, opts ...client.CallOption) (*CommentResponse, error)
}

type commentSrvService struct {
	c    client.Client
	name string
}

func NewCommentSrvService(name string, c client.Client) CommentSrvService {
	return &commentSrvService{
		c:    c,
		name: name,
	}
}

func (c *commentSrvService) AddComment(ctx context.Context, in *AddCommentRequest, opts ...client.CallOption) (*CommonResponse, error) {
	req := c.c.NewRequest(c.name, "CommentSrv.AddComment", in)
	out := new(CommonResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentSrvService) GetComments(ctx context.Context, in *CommentsRequest, opts ...client.CallOption) (*CommentResponse, error) {
	req := c.c.NewRequest(c.name, "CommentSrv.GetComments", in)
	out := new(CommentResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentSrvService) DianZan(ctx context.Context, in *DianZanRequest, opts ...client.CallOption) (*CommentResponse, error) {
	req := c.c.NewRequest(c.name, "CommentSrv.DianZan", in)
	out := new(CommentResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CommentSrv service

type CommentSrvHandler interface {
	AddComment(context.Context, *AddCommentRequest, *CommonResponse) error
	GetComments(context.Context, *CommentsRequest, *CommentResponse) error
	DianZan(context.Context, *DianZanRequest, *CommentResponse) error
}

func RegisterCommentSrvHandler(s server.Server, hdlr CommentSrvHandler, opts ...server.HandlerOption) error {
	type commentSrv interface {
		AddComment(ctx context.Context, in *AddCommentRequest, out *CommonResponse) error
		GetComments(ctx context.Context, in *CommentsRequest, out *CommentResponse) error
		DianZan(ctx context.Context, in *DianZanRequest, out *CommentResponse) error
	}
	type CommentSrv struct {
		commentSrv
	}
	h := &commentSrvHandler{hdlr}
	return s.Handle(s.NewHandler(&CommentSrv{h}, opts...))
}

type commentSrvHandler struct {
	CommentSrvHandler
}

func (h *commentSrvHandler) AddComment(ctx context.Context, in *AddCommentRequest, out *CommonResponse) error {
	return h.CommentSrvHandler.AddComment(ctx, in, out)
}

func (h *commentSrvHandler) GetComments(ctx context.Context, in *CommentsRequest, out *CommentResponse) error {
	return h.CommentSrvHandler.GetComments(ctx, in, out)
}

func (h *commentSrvHandler) DianZan(ctx context.Context, in *DianZanRequest, out *CommentResponse) error {
	return h.CommentSrvHandler.DianZan(ctx, in, out)
}
