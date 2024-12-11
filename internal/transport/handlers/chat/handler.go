package chat

import (
	"context"
	"errors"
	"strings"

	cl "github.com/Mobo140/microservices/chat/internal/client"
	conv "github.com/Mobo140/microservices/chat/internal/converter"
	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/Mobo140/microservices/chat/internal/service"
	"github.com/Mobo140/platform_common/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
)

type Implementation struct {
	desc.UnimplementedChatV1Server
	chatAPIService      service.ChatService
	accessServiceClient cl.AccessServiceClient
}

func NewImplementation(chatService service.ChatService, accessServiceClient cl.AccessServiceClient) *Implementation {
	return &Implementation{chatAPIService: chatService, accessServiceClient: accessServiceClient}
}

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger.Info("Creating chat...", zap.Any("info", req.GetInfo()))

	info, err := conv.ToChatInfoFromDesc(req.Info)
	if err != nil {
		logger.Error("Failed to convert to chat info from desc", zap.Error(err))

		return nil, err
	}

	id, err := i.chatAPIService.Create(ctx, info)
	if err != nil {
		logger.Error("Failed to create chat", zap.Error(err))

		return nil, err
	}

	logger.Info("Create chat: ", zap.Int64("id", id))

	return &desc.CreateResponse{Id: id}, nil
}

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	logger.Info("Getting chat...", zap.Any("info", req.GetId()))

	chat, err := i.chatAPIService.Get(ctx, req.GetId())
	if err != nil {
		logger.Error("Failed to get to chat by id", zap.Int64("id", req.GetId()), zap.Error(err))

		return nil, err
	}

	logger.Info("Get chat: ", zap.Int64("id", chat.ID), zap.Any("usernames", strings.Join(chat.Info.Usernames, ", ")))

	chatDesc := conv.ToChatFromService(chat)

	return &desc.GetResponse{
		Chat: chatDesc,
	}, nil
}

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	logger.Info("Deletting chat...", zap.Any("info", req.GetId()))

	err := i.chatAPIService.Delete(ctx, req.GetId())
	if err != nil {
		logger.Error("Failed to delete chat by id", zap.Int64("id", req.GetId()), zap.Error(err))

		return nil, err
	}

	logger.Info("Delete chat: ", zap.Int64("id", req.GetId()))

	return &emptypb.Empty{}, nil
}

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendMessage")
	defer span.Finish()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := errors.New("metadata is not provided")
		logger.Error("Failed to get metadata from context: ",
			zap.Any("request", req),
			zap.Error(err),
		)

		return nil, err
	}

	for key, values := range md {
		for _, value := range values {
			ctx = metadata.AppendToOutgoingContext(ctx, key, value)
		}
	}

	logger.Info("Checking the access token...")

	err := i.accessServiceClient.Check(ctx, "/user_v1.ChatV1/SendMessage")
	if err != nil {
		logger.Error("Failed to check the access token", zap.Error(err))

		return nil, err
	}

	logger.Info("Access granted")

	logger.Info("Sending message to chat...", zap.Any("chat id", req.GetChatId()), zap.Any("message", req.GetMessage()))

	messageInfo, err := conv.ToMessageFromDesc(req.Message)
	if err != nil {
		logger.Info("Sending message to chat...", zap.Any("chat id", req.GetChatId()), zap.Any("message", req.GetMessage()))

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
		logger.Info("Failed to send message to chat...",
			zap.Any("chat id", req.GetChatId()),
			zap.Any("message", req.GetMessage()),
			zap.Error(err),
		)

		return nil, err
	}

	logger.Info("Send messsage to chat: ", zap.Any("chat id", req.GetChatId()), zap.Any("message", req.GetMessage()))

	return &emptypb.Empty{}, nil
}
