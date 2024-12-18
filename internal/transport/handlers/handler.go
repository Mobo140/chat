package transport

import (
	"context"

	desc "github.com/Mobo140/chat/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatAPIHandler interface {
	Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error)
	Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error)
	Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error)
	// Update(ctx context.Context, info *model.UpdateInfo) error
	SendMessage(cfg context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error)
	ConnectChat(ctx context.Context, req *desc.ConnectChatRequest, stream desc.ChatV1_ConnectChatServer)  error
	// GetMessagesByChatID()
}
