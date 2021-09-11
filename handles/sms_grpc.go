package handles

import (
	"context"

	"github.com/UniqueStudio/open-platform/pb/uni_sms"
)

type TencentSMSClient struct {
	uni_sms.UnimplementedSMSServiceServer
}

func NewTencentSMSGrpcServer() *TencentSMSClient {
	return &TencentSMSClient{}
}

func (tsc *TencentSMSClient) PushSMS(ctx context.Context, in *uni_sms.PushSMSRequest) (*uni_sms.PushSMSResponse, error) {
	return &uni_sms.PushSMSResponse{
		SMSStatus: []*uni_sms.SMSStatus{
			{
				SerialNo:    "abababa",
				PhoneNumber: "18272008762",
			},
		},
	}, nil
}

func (tsc *TencentSMSClient) AddSMSSignature(ctx context.Context, in *uni_sms.AddSMSSignatureRequest) (*uni_sms.UniformResponse, error) {
	return &uni_sms.UniformResponse{
		Success: false,
		Message: "error",
		Id:      "",
	}, nil
}

func (tsc *TencentSMSClient) AddSMSTemplate(ctx context.Context, in *uni_sms.AddSMSTemplateRequest) (*uni_sms.UniformResponse, error) {
	return &uni_sms.UniformResponse{
		Success: false,
		Message: "error",
		Id:      "",
	}, nil
}
