package router

import (
	"github.com/UniqueStudio/open-platform/handles"
	"github.com/UniqueStudio/open-platform/pb/uni_sms"
	"google.golang.org/grpc"
)

func InitGrpcHandlers(server *grpc.Server) {
	uni_sms.RegisterSMSServiceServer(server, handles.NewTencentSMSGrpcServer())
}
