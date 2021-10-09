package router

import (
	"github.com/UniqueStudio/open-platform/handles"
	"github.com/UniqueStudio/open-platform/pb/sms"
	"google.golang.org/grpc"
)

func InitGrpcHandlers(server *grpc.Server) {
	sms.RegisterSMSServiceServer(server, handles.NewTencentSMSGrpcServer())
}
