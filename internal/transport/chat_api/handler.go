package chat

import (
	"context"
	"log"
	"strings"

	conv "github.com/Mobo140/microservices/chat/internal/converter"
	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/Mobo140/microservices/chat/internal/service"
	"github.com/Mobo140/microservices/chat/internal/transport"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
)

var _ transport.ChatAPIHandler = (*implementation)(nil)

type implementation struct {
	desc.UnimplementedChatV1Server
	chatAPIService service.ChatService
}

func NewImplementation(chatService service.ChatService) *implementation {
	return &implementation{chatAPIService: chatService}
}

func (i *implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	info := conv.ToChatInfoFromDesc(req.Info)

	id, err := i.chatAPIService.Create(ctx, info)
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{Id: id}, nil

}

func (i *implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

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

func (i *implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatAPIService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (i *implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	messageInfo := conv.ToMessageFromDesc(req.Message)

	message := &model.Message{
		ChatID: req.ChatId,
		Info: model.MessageInfo{
			From:      messageInfo.From,
			Text:      messageInfo.Text,
			Timestamp: messageInfo.Timestamp,
		},
	}

	err := i.chatAPIService.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
