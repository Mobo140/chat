package chat

import (
	"context"
	"strconv"
	"strings"
	"sync"

	cl "github.com/Mobo140/chat/internal/client"
	conv "github.com/Mobo140/chat/internal/converter"
	"github.com/Mobo140/chat/internal/model"
	"github.com/Mobo140/chat/internal/service"
	"github.com/Mobo140/platform_common/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/Mobo140/chat/pkg/chat_v1"
)

type Implementation struct {
	desc.UnimplementedChatV1Server
	chatAPIService      service.ChatService
	accessServiceClient cl.AccessServiceClient

	chats  map[string]*Chat
	mxChat sync.RWMutex

	channels  map[string]chan *desc.Message
	mxChannel sync.RWMutex
}

func NewImplementation(chatService service.ChatService, accessServiceClient cl.AccessServiceClient) *Implementation {
	return &Implementation{
		chatAPIService:      chatService,
		accessServiceClient: accessServiceClient,
		chats:               make(map[string]*Chat),
		channels:            make(map[string]chan *desc.Message),
	}
}

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Create chat")
	defer span.Finish()

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

	i.channels[strconv.FormatInt(id, 10)] = make(chan *desc.Message, 100)

	logger.Info("Create chat: ", zap.Int64("id", id))

	return &desc.CreateResponse{Id: id}, nil
}

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get chat")
	defer span.Finish()

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "Delete chat")
	defer span.Finish()

	logger.Info("Deletting chat...", zap.Any("info", req.GetId()))

	err := i.chatAPIService.Delete(ctx, req.GetId())
	if err != nil {
		logger.Error("Failed to delete chat by id", zap.Int64("id", req.GetId()), zap.Error(err))

		return nil, err
	}

	logger.Info("Delete chat: ", zap.Int64("id", req.GetId()))

	return &emptypb.Empty{}, nil
}

func (i *Implementation) ConnectChat(req *desc.ConnectChatRequest, stream desc.ChatV1_ConnectChatServer) error {
	span, _ := opentracing.StartSpanFromContext(stream.Context(), "ConnectChat")
	defer span.Finish()

	logger.Info("Attempting to connect to chat...", 
		zap.String("chat_id", req.GetChatId()),
		zap.String("username", req.GetUsername()))

	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetChatId()]
	i.mxChannel.RUnlock()

	if !ok {
		logger.Error("Chat not found", zap.String("chat_id", req.GetChatId()))
		return status.Errorf(codes.NotFound, "chat not found")
	}

	i.mxChat.Lock()
	if _, okChat := i.chats[req.GetChatId()]; !okChat {
		logger.Info("Creating new chat instance", zap.String("chat_id", req.GetChatId()))
		i.chats[req.GetChatId()] = NewChat()
	}
	i.mxChat.Unlock()

	i.chats[req.GetChatId()].m.Lock()
	i.chats[req.GetChatId()].streams[req.GetUsername()] = stream
	i.chats[req.GetChatId()].m.Unlock()

	logger.Info("Successfully connected to chat", 
		zap.String("chat_id", req.GetChatId()),
		zap.String("username", req.GetUsername()))

	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				logger.Info("Chat channel closed", 
					zap.String("chat_id", req.GetChatId()),
					zap.String("username", req.GetUsername()))
				return nil
			}

			logger.Info("Received message in chat", 
				zap.String("chat_id", req.GetChatId()),
				zap.Any("message", msg))

			for username, st := range i.chats[req.GetChatId()].streams {
				if username == req.GetUsername() {
					continue
				}

				logger.Info("Attempting to send message to user", 
					zap.String("to_username", username),
					zap.String("from_username", req.GetUsername()))

				if err := st.Send(msg); err != nil {
					logger.Error("Failed to send message to stream", 
						zap.String("to_username", username),
						zap.String("from_username", req.GetUsername()),
						zap.Error(err))
					return err
				}

				logger.Info("Successfully sent message to user", 
					zap.String("to_username", username),
					zap.String("from_username", req.GetUsername()))
			}

		case <-stream.Context().Done():
			logger.Info("Disconnecting from chat", 
				zap.String("chat_id", req.GetChatId()),
				zap.String("username", req.GetUsername()))

			i.chats[req.GetChatId()].m.Lock()
			delete(i.chats[req.GetChatId()].streams, req.GetUsername())
			i.chats[req.GetChatId()].m.Unlock()

			return nil
		}
	}
}

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendMessage")
	defer span.Finish()

	type Session struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
		Username     string `json:"username"`
	}

	logger.Info("Checking the access token...")

	err := i.accessServiceClient.Check(ctx, "/user_v1.ChatV1/SendMessage")
	if err != nil {
		logger.Error("Failed to check the access token", zap.Error(err))

		return nil, err
	}

	logger.Info("Access granted")

	logger.Info("Getting chat channel...")

	i.mxChannel.RLock()
	chatChan, ok := i.channels[strconv.FormatInt(req.GetChatId(), 10)]
	i.mxChannel.RUnlock()

	if !ok {
		logger.Error("Chat not found", zap.Int64("chat id", req.GetChatId()))

		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	logger.Info("Sending message to chat...", zap.Any("chat id", req.GetChatId()), zap.Any("message", req.GetMessage()))

	messageInfo, err := conv.ToMessageFromDesc(req.Message)
	if err != nil {
		logger.Error("Failed to convert message to desc",
			zap.Any("chat id", req.GetChatId()),
			zap.Any("message", req.GetMessage()),
			zap.Error(err),
		)

		return nil, err
	}

	message := &model.SendMessage{
		ChatID: req.GetChatId(),
		Message: model.Message{
			From:      messageInfo.From,
			Text:      messageInfo.Text,
			CreatedAt: messageInfo.CreatedAt,
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

	chatChan <- req.GetMessage()

	logger.Info("Send messsage to chat: ", zap.Any("chat id", req.GetChatId()), zap.Any("message", req.GetMessage()))

	return &emptypb.Empty{}, nil
}
