package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Mobo140/chat/internal/model"
	"github.com/Mobo140/chat/internal/service"
	serviceMocks "github.com/Mobo140/chat/internal/service/mocks"
	chatHandler "github.com/Mobo140/chat/internal/transport/handlers/chat"
	desc "github.com/Mobo140/chat/pkg/chat_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	value int64 = 4
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type setupMocks func(mockService *serviceMocks.ChatServiceMock)

	type args struct {
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		usernames = []string{gofakeit.Username(), gofakeit.Username()}

		serviceErr      = fmt.Errorf("service create error")
		conversationErr = fmt.Errorf("chatInfo is empty")

		info = &model.ChatInfo{
			Usernames: usernames,
		}

		req = &desc.CreateRequest{
			Info: &desc.ChatInfo{
				Usernames: usernames,
			},
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name         string
		args         args
		setupMocks   setupMocks
		expectedResp *desc.CreateResponse
		expectedErr  error
	}{
		{
			name: "success case",
			args: args{
				req: req,
			},
			setupMocks: func(mock *serviceMocks.ChatServiceMock) {
				mock.CreateMock.Expect(ctx, info).Return(id, nil)
			},
			expectedResp: res,
			expectedErr:  nil,
		},
		{
			name: "empty request",
			args: args{
				req: nil,
			},
			setupMocks: func(mock *serviceMocks.ChatServiceMock) {
				// Мок не должен вызываться
			},
			expectedResp: nil,
			expectedErr:  conversationErr,
		},
		{
			name: "service error",
			args: args{
				req: req,
			},
			setupMocks: func(mock *serviceMocks.ChatServiceMock) {
				mock.CreateMock.Expect(ctx, info).Return(0, serviceErr)
			},
			expectedResp: nil,
			expectedErr:  serviceErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Создаем мок сервиса
			mockService := serviceMocks.NewChatServiceMock(mc)
			// Настраиваем мок
			tt.setupMocks(mockService)
			// Создаем handler
			handler := chatHandler.NewImplementation(mockService, nil)

			// Выполняем тест
			resp, err := handler.Create(ctx, tt.args.req)

			// Проверяем результаты
			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedResp, resp)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	type setupMocks func(mockService *serviceMocks.ChatServiceMock)

	type args struct {
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		usernames = []string{gofakeit.Username(), gofakeit.Username()}

		serviceErr = fmt.Errorf("service error")

		chat = &model.Chat{
			ID: id,
			Info: model.ChatInfo{
				Usernames: usernames,
			},
		}

		req = &desc.GetRequest{
			Id: id,
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
		name         string
		args         args
		setupMocks   setupMocks
		expectedResp *desc.GetResponse
		expectedErr  error
	}{
		{
			name: "success case",
			args: args{
				req: req,
			},
			setupMocks: func(mock *serviceMocks.ChatServiceMock) {
				mock.GetMock.Expect(ctx, id).Return(chat, nil)
			},
			expectedResp: res,
			expectedErr:  nil,
		},
		{
			name: "service error",
			args: args{
				req: req,
			},
			setupMocks: func(mock *serviceMocks.ChatServiceMock) {
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
			},
			expectedResp: nil,
			expectedErr:  serviceErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := serviceMocks.NewChatServiceMock(mc)
			tt.setupMocks(mockService)
			handler := chatHandler.NewImplementation(mockService, nil)

			resp, err := handler.Get(ctx, tt.args.req)
			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedResp, resp)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

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
		chatServiceMock chatServiceMockFunc
		want            *emptypb.Empty
		err             error
	}{
		{
			name: "success case",
			args: args{
				req: req,
			},
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
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
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			handler := chatHandler.NewImplementation(chatServiceMock, nil)

			response, err := handler.Delete(ctx, tt.args.req)
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
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = value
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
			ChatID: value,
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
			name: "success case",
			args: args{
				req: req,
			},
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(nil)
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
				mock.SendMessageMock.Expect(ctx, message).Return(serviceErr)
				return mock
			},
		},
		{
			name: "empty message case",
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
			handler := chatHandler.NewImplementation(chatServiceMock, nil)

			response, err := handler.SendMessage(ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
