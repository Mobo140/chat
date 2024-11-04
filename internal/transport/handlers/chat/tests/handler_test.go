package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/Mobo140/microservices/chat/internal/service"
	serviceMocks "github.com/Mobo140/microservices/chat/internal/service/mocks"
	chatHandler "github.com/Mobo140/microservices/chat/internal/transport/handlers/chat"
	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	// "google.golang.org/protobuf/types/known/emptypb"
	// "google.golang.org/protobuf/types/known/timestamppb"
	// "google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		req *desc.CreateRequest
	}

	var (
		ctxValue = context.Background()
		mc       = minimock.NewController(t)

		id        = gofakeit.Int64()
		usernames = []string{gofakeit.Username()}

		serviceErr      = fmt.Errorf("service create error")
		conversationErr = fmt.Errorf("chatInfo is empty")

		req = &desc.CreateRequest{
			Info: &desc.ChatInfo{
				Usernames: usernames,
			},
		}

		info = &model.ChatInfo{
			Usernames: usernames,
		}

		res = &desc.CreateResponse{
			Id: id,
		}

		unknownChat = (int64)(-1)
	)

	tests := []struct {
		name            string
		args            args
		userServiceMock userServiceMockFunc
		want            *desc.CreateResponse
		err             error
	}{
		{
			name: "success",
			args: args{
				req: req,
			},
			userServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateMock.Expect(ctxValue, info).Return(id, nil)
				return mock
			},
			want: res,
			err:  nil,
		},
		{
			name: "service error case",
			args: args{
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateMock.Expect(ctxValue, info).Return(unknownChat, serviceErr)
				return mock
			},
		},
		{
			name: "conversation error case",
			args: args{
				req: &desc.CreateRequest{
					Info: nil,
				},
			},
			userServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
			want: nil,
			err:  conversationErr,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			transport := chatHandler.NewImplementation(userServiceMock)

			newID, err := transport.Create(ctxValue, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		req *desc.GetRequest
	}

	var (
		ctxValue = context.Background()
		mc       = minimock.NewController(t)

		id        = gofakeit.Int64()
		usernames = []string{gofakeit.Username()}

		serviceErr = fmt.Errorf("service error")

		req = &desc.GetRequest{
			Id: id,
		}

		chat = &model.Chat{
			ID: id,
			Info: model.ChatInfo{
				Usernames: usernames,
			},
		}

		res = &desc.GetResponse{
			Chat: &desc.Chat{
				Id: id,
				Info: &desc.ChatInfo{
					Usernames: usernames,
				},
			},
		}
	)

	tests := []struct {
		name            string
		args            args
		chatServiceMock chatServiceMockFunc
		want            *desc.GetResponse
		err             error
	}{
		{
			name: "success",
			args: args{
				req: req,
			},
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetMock.Expect(ctxValue, id).Return(chat, nil)
				return mock
			},
			want: res,
			err:  nil,
		},
		{
			name: "service error case",
			args: args{
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetMock.Expect(ctxValue, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			transport := chatHandler.NewImplementation(chatServiceMock)

			response, err := transport.Get(ctxValue, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		req *desc.DeleteRequest
	}

	var (
		ctxValue = context.Background()
		mc       = minimock.NewController(t)

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("service error")

		req = &desc.DeleteRequest{
			Id: id,
		}

		res = &emptypb.Empty{}
	)

	tests := []struct {
		name            string
		args            args
		userServiceMock userServiceMockFunc
		want            *emptypb.Empty
		err             error
	}{
		{
			name: "success",
			args: args{
				req: req,
			},
			userServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.DeleteMock.Expect(ctxValue, id).Return(nil)
				return mock
			},
			want: res,
			err:  nil,
		},
		{
			name: "service error case",
			args: args{
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.DeleteMock.Expect(ctxValue, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.userServiceMock(mc)
			transport := chatHandler.NewImplementation(chatServiceMock)

			response, err := transport.Delete(ctxValue, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		req *desc.SendMessageRequest
	}

	var (
		ctxValue = context.Background()
		mc       = minimock.NewController(t)

		id   = gofakeit.Int64()
		from = gofakeit.Name()
		text = gofakeit.Color()

		serviceErr  = fmt.Errorf("service update error")
		converseErr = fmt.Errorf("message is empty")

		req = &desc.SendMessageRequest{
			ChatId: id,
			Message: &desc.Message{
				From: from,
				Text: text,
			},
		}

		message = &model.SendMessage{
			ChatID: id,
			Message: model.Message{
				From: from,
				Text: text,
			},
		}

		res = &emptypb.Empty{}
	)

	tests := []struct {
		name            string
		args            args
		chatServiceMock chatServiceMockFunc
		want            *emptypb.Empty
		err             error
	}{
		{
			name: "success",
			args: args{
				req: req,
			},
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctxValue, message).Return(nil)
				return mock
			},
			want: res,
			err:  nil,
		},
		{
			name: "service error case",
			args: args{
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctxValue, message).Return(serviceErr)
				return mock
			},
		},
		{
			name: "conversation error case",
			args: args{
				req: &desc.SendMessageRequest{
					ChatId:  id,
					Message: nil,
				},
			},
			want: nil,
			err:  converseErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			transport := chatHandler.NewImplementation(chatServiceMock)

			response, err := transport.SendMessage(ctxValue, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
