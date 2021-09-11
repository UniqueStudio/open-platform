// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package uni_sms

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SMSServiceClient is the client API for SMSService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SMSServiceClient interface {
	PushSMS(ctx context.Context, in *PushSMSRequest, opts ...grpc.CallOption) (*PushSMSResponse, error)
	AddSMSSignature(ctx context.Context, in *AddSMSSignatureRequest, opts ...grpc.CallOption) (*UniformResponse, error)
	AddSMSTemplate(ctx context.Context, in *AddSMSTemplateRequest, opts ...grpc.CallOption) (*UniformResponse, error)
}

type sMSServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSMSServiceClient(cc grpc.ClientConnInterface) SMSServiceClient {
	return &sMSServiceClient{cc}
}

func (c *sMSServiceClient) PushSMS(ctx context.Context, in *PushSMSRequest, opts ...grpc.CallOption) (*PushSMSResponse, error) {
	out := new(PushSMSResponse)
	err := c.cc.Invoke(ctx, "/sms.SMSService/PushSMS", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sMSServiceClient) AddSMSSignature(ctx context.Context, in *AddSMSSignatureRequest, opts ...grpc.CallOption) (*UniformResponse, error) {
	out := new(UniformResponse)
	err := c.cc.Invoke(ctx, "/sms.SMSService/AddSMSSignature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sMSServiceClient) AddSMSTemplate(ctx context.Context, in *AddSMSTemplateRequest, opts ...grpc.CallOption) (*UniformResponse, error) {
	out := new(UniformResponse)
	err := c.cc.Invoke(ctx, "/sms.SMSService/AddSMSTemplate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SMSServiceServer is the server API for SMSService service.
// All implementations must embed UnimplementedSMSServiceServer
// for forward compatibility
type SMSServiceServer interface {
	PushSMS(context.Context, *PushSMSRequest) (*PushSMSResponse, error)
	AddSMSSignature(context.Context, *AddSMSSignatureRequest) (*UniformResponse, error)
	AddSMSTemplate(context.Context, *AddSMSTemplateRequest) (*UniformResponse, error)
	mustEmbedUnimplementedSMSServiceServer()
}

// UnimplementedSMSServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSMSServiceServer struct {
}

func (UnimplementedSMSServiceServer) PushSMS(context.Context, *PushSMSRequest) (*PushSMSResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushSMS not implemented")
}
func (UnimplementedSMSServiceServer) AddSMSSignature(context.Context, *AddSMSSignatureRequest) (*UniformResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSMSSignature not implemented")
}
func (UnimplementedSMSServiceServer) AddSMSTemplate(context.Context, *AddSMSTemplateRequest) (*UniformResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSMSTemplate not implemented")
}
func (UnimplementedSMSServiceServer) mustEmbedUnimplementedSMSServiceServer() {}

// UnsafeSMSServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SMSServiceServer will
// result in compilation errors.
type UnsafeSMSServiceServer interface {
	mustEmbedUnimplementedSMSServiceServer()
}

func RegisterSMSServiceServer(s grpc.ServiceRegistrar, srv SMSServiceServer) {
	s.RegisterService(&SMSService_ServiceDesc, srv)
}

func _SMSService_PushSMS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushSMSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SMSServiceServer).PushSMS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sms.SMSService/PushSMS",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SMSServiceServer).PushSMS(ctx, req.(*PushSMSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SMSService_AddSMSSignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddSMSSignatureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SMSServiceServer).AddSMSSignature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sms.SMSService/AddSMSSignature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SMSServiceServer).AddSMSSignature(ctx, req.(*AddSMSSignatureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SMSService_AddSMSTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddSMSTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SMSServiceServer).AddSMSTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sms.SMSService/AddSMSTemplate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SMSServiceServer).AddSMSTemplate(ctx, req.(*AddSMSTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SMSService_ServiceDesc is the grpc.ServiceDesc for SMSService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SMSService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sms.SMSService",
	HandlerType: (*SMSServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushSMS",
			Handler:    _SMSService_PushSMS_Handler,
		},
		{
			MethodName: "AddSMSSignature",
			Handler:    _SMSService_AddSMSSignature_Handler,
		},
		{
			MethodName: "AddSMSTemplate",
			Handler:    _SMSService_AddSMSTemplate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sms.proto",
}