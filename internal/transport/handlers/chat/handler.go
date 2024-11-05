package chat

import (
	"context"
	"log"
	"strings"

	conv "github.com/Mobo140/microservices/chat/internal/converter"
	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/Mobo140/microservices/chat/internal/service"
	transport "github.com/Mobo140/microservices/chat/internal/transport/handlers"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
)

var _ transport.ChatAPIHandler = (*Implementation)(nil)

type Implementation struct {
	desc.UnimplementedChatV1Server
	chatAPIService service.ChatService
}

func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{chatAPIService: chatService}
}

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	info, err := conv.ToChatInfoFromDesc(req.Info)
	if err != nil {
		return nil, err
	}

	id, err := i.chatAPIService.Create(ctx, info)
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{Id: id}, nil
}

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	chat, err := i.chatAPIService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, usernames: %s\n",
		chat.ID, strings.Join(chat.Info.Usernames, ", "),
	)

	chatDesc := conv.ToChatFromService(chat)

	return &desc.GetResponse{
		Chat: chatDesc,
	}, nil
}

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatAPIService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	messageInfo, err := conv.ToMessageFromDesc(req.Message)
	if err != nil {
		return nil, err
	}
	message := &model.SendMessage{
		ChatID: req.ChatId,
		Message: model.Message{
			From: messageInfo.From,
			Text: messageInfo.Text,
		},
	}

	err = i.chatAPIService.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
