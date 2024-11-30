package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	repositoryTx "github.com/Mobo140/platform_common/pkg/db"
	dbTxMocks "github.com/Mobo140/platform_common/pkg/db/mocks"
	"github.com/Mobo140/microservices/chat/internal/model"
	repositoryMocks "github.com/Mobo140/microservices/chat/internal/repository/mocks"
	chatService "github.com/Mobo140/microservices/chat/internal/service/chat"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type setupMocks func(
		chatRepo *repositoryMocks.ChatRepositoryMock,
		messageRepo *repositoryMocks.MessageRepositoryMock,
		logRepo *repositoryMocks.LogRepositoryMock,
		txManager *dbTxMocks.TxManagerMock,
	)

	type args struct {
		req *model.ChatInfo
	}

	var (
		ctxValue = context.Background()
		mc       = minimock.NewController(t)

		id        = gofakeit.Int64()
		usernames = []string{gofakeit.Username()}

		repositoryErr  = fmt.Errorf("create chatRepo error")
		logErr         = fmt.Errorf("create log error")
		transactionErr = fmt.Errorf("transaction error")

		info = &model.ChatInfo{
			Usernames: usernames,
		}
		req = &model.ChatInfo{
			Usernames: usernames,
		}

		logEntry = &model.LogEntry{
			ChatID:   id,
			Activity: fmt.Sprintf("Create chat: usernames:%s", strings.Join(info.Usernames, ", ")),
		}

		unknownChat = (int64)(-1)
	)

	tests := []struct {
		name       string
		args       args
		setupMocks setupMocks
		want       int64
		err        error
	}{
		{
			name: "success case",
			args: args{
				req: req,
			},
			want: id,
			err:  nil,
			setupMocks: func(chatRepo *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				logRepo *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				chatRepo.CreateMock.Expect(ctxValue, info).Return(id, nil)
				logRepo.CreateMock.Expect(ctxValue, logEntry).Return(nil)
				txManager.ReadCommitedMock.Set(func(ctx context.Context, f repositoryTx.Handler) error {
					return f(ctx)
				})
			},
		},
		{
			name: "transaction error",
			args: args{
				req: req,
			},
			want: unknownChat,
			err:  transactionErr,
			setupMocks: func(_ *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				_ *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				txManager.ReadCommitedMock.Set(func(_ context.Context, _ repositoryTx.Handler) error {
					return transactionErr
				})
			},
		},
		{
			name: "chatRepo error",
			args: args{
				req: req,
			},
			want: unknownChat,
			err:  repositoryErr,
			setupMocks: func(chatRepo *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				_ *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				chatRepo.CreateMock.Expect(ctxValue, info).Return(unknownChat, repositoryErr)
				txManager.ReadCommitedMock.Set(func(ctxValue context.Context, f repositoryTx.Handler) error {
					return f(ctxValue)
				})
			},
		},
		{
			name: "creating log in db error",
			args: args{
				req: req,
			},
			want: unknownChat,
			err:  logErr,
			setupMocks: func(chatRepo *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				logRepo *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				chatRepo.CreateMock.Expect(ctxValue, info).Return(id, nil)
				logRepo.CreateMock.Expect(ctxValue, logEntry).Return(logErr)
				txManager.ReadCommitedMock.Set(func(ctxValue context.Context, f repositoryTx.Handler) error {
					return f(ctxValue)
				})
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepo := repositoryMocks.NewChatRepositoryMock(mc)
			messageRepo := repositoryMocks.NewMessageRepositoryMock(mc)
			logRepo := repositoryMocks.NewLogRepositoryMock(mc)
			txManager := dbTxMocks.NewTxManagerMock(mc)

			// Настройка моков в соответствии с тестами
			tt.setupMocks(chatRepo, messageRepo, logRepo, txManager)

			service := chatService.NewService(chatRepo, messageRepo, logRepo, txManager)

			gotID, err := service.Create(ctxValue, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, gotID)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	type setupMocks func(
		chatRepo *repositoryMocks.ChatRepositoryMock,
		messageRepo *repositoryMocks.MessageRepositoryMock,
		logRepo *repositoryMocks.LogRepositoryMock,
		txManager *dbTxMocks.TxManagerMock,
	)

	type args struct {
		req int64
	}

	var (
		ctxValue = context.Background()
		mc       = minimock.NewController(t)

		id        = gofakeit.Int64()
		usernames = []string{gofakeit.Username()}

		repositoryErr  = fmt.Errorf("update userRepo error")
		logErr         = fmt.Errorf("update log error")
		transactionErr = fmt.Errorf("transaction error")

		chat = &model.Chat{
			ID: id,
			Info: model.ChatInfo{
				Usernames: usernames,
			},
		}

		logEntry = &model.LogEntry{
			ChatID:   id,
			Activity: fmt.Sprintf("Get chat: Id:%d, Usernames:%s", id, strings.Join(chat.Info.Usernames, ", ")),
		}
	)

	tests := []struct {
		name       string
		setupMocks setupMocks
		args       args
		want       *model.Chat
		err        error
	}{
		{
			name: "success case",
			want: chat,
			args: args{
				req: id,
			},
			err: nil,
			setupMocks: func(chatRepo *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				logRepo *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				chatRepo.GetMock.Expect(ctxValue, id).Return(chat, nil)
				logRepo.CreateMock.Expect(ctxValue, logEntry).Return(nil)
				txManager.ReadCommitedMock.Set(func(ctx context.Context, f repositoryTx.Handler) error {
					return f(ctx)
				})
			},
		},
		{
			name: "transaction error",
			want: nil,
			args: args{
				req: id,
			},
			err: transactionErr,
			setupMocks: func(_ *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				_ *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				txManager.ReadCommitedMock.Set(func(_ context.Context, _ repositoryTx.Handler) error {
					return transactionErr
				})
			},
		},
		{
			name: "chatRepo error",
			want: nil,
			args: args{
				req: id,
			},
			err: repositoryErr,
			setupMocks: func(userRepo *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				_ *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				userRepo.GetMock.Expect(ctxValue, id).Return(nil, repositoryErr)
				txManager.ReadCommitedMock.Set(func(ctxValue context.Context, f repositoryTx.Handler) error {
					return f(ctxValue)
				})
			},
		},
		{
			name: "creating log in db error",
			want: nil,
			args: args{
				req: id,
			},
			err: logErr,
			setupMocks: func(userRepo *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				logRepo *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				userRepo.GetMock.Expect(ctxValue, id).Return(chat, nil)
				logRepo.CreateMock.Expect(ctxValue, logEntry).Return(logErr)
				txManager.ReadCommitedMock.Set(func(ctxValue context.Context, f repositoryTx.Handler) error {
					return f(ctxValue)
				})
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepo := repositoryMocks.NewChatRepositoryMock(mc)
			messageRepo := repositoryMocks.NewMessageRepositoryMock(mc)
			logRepo := repositoryMocks.NewLogRepositoryMock(mc)
			txManager := dbTxMocks.NewTxManagerMock(mc)

			// Настройка моков в соответствии с тестами
			tt.setupMocks(chatRepo, messageRepo, logRepo, txManager)

			service := chatService.NewService(chatRepo, messageRepo, logRepo, txManager)

			gotID, err := service.Get(ctxValue, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, gotID)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	type setupMocks func(
		userRepo *repositoryMocks.ChatRepositoryMock,
		messageRepo *repositoryMocks.MessageRepositoryMock,
		logRepo *repositoryMocks.LogRepositoryMock,
		txManager *dbTxMocks.TxManagerMock,
	)

	type args struct {
		req int64
	}

	var (
		ctxValue = context.Background()
		mc       = minimock.NewController(t)

		id = gofakeit.Int64()

		repositoryErr  = fmt.Errorf("delete chatRepo error")
		logErr         = fmt.Errorf("delete log error")
		transactionErr = fmt.Errorf("transaction error")

		logEntry = &model.LogEntry{
			ChatID:   id,
			Activity: fmt.Sprintf("Delete chat: ID=%d", id),
		}
	)

	tests := []struct {
		name       string
		setupMocks setupMocks
		args       args
		err        error
	}{
		{
			name: "success case",
			err:  nil,
			args: args{
				req: id,
			},
			setupMocks: func(chatRepo *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				logRepo *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				chatRepo.DeleteMock.Expect(ctxValue, id).Return(nil)
				logRepo.CreateMock.Expect(ctxValue, logEntry).Return(nil)
				txManager.ReadCommitedMock.Set(func(ctx context.Context, f repositoryTx.Handler) error {
					return f(ctx)
				})
			},
		},
		{
			name: "transaction error",
			err:  transactionErr,
			args: args{
				req: id,
			},
			setupMocks: func(_ *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				_ *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				txManager.ReadCommitedMock.Set(func(_ context.Context, _ repositoryTx.Handler) error {
					return transactionErr
				})
			},
		},
		{
			name: "chatRepo error",
			err:  repositoryErr,
			args: args{
				req: id,
			},
			setupMocks: func(chatRepo *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				_ *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				chatRepo.DeleteMock.Expect(ctxValue, id).Return(repositoryErr)
				txManager.ReadCommitedMock.Set(func(ctxValue context.Context, f repositoryTx.Handler) error {
					return f(ctxValue)
				})
			},
		},
		{
			name: "creating log in db error",
			err:  logErr,
			args: args{
				req: id,
			},
			setupMocks: func(userRepo *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				logRepo *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				userRepo.DeleteMock.Expect(ctxValue, id).Return(nil)
				logRepo.CreateMock.Expect(ctxValue, logEntry).Return(logErr)
				txManager.ReadCommitedMock.Set(func(ctxValue context.Context, f repositoryTx.Handler) error {
					return f(ctxValue)
				})
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepo := repositoryMocks.NewChatRepositoryMock(mc)
			messageRepo := repositoryMocks.NewMessageRepositoryMock(mc)
			logRepo := repositoryMocks.NewLogRepositoryMock(mc)
			txManager := dbTxMocks.NewTxManagerMock(mc)

			// Настройка моков в соответствии с тестами
			tt.setupMocks(userRepo, messageRepo, logRepo, txManager)

			service := chatService.NewService(userRepo, messageRepo, logRepo, txManager)

			err := service.Delete(ctxValue, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type setupMocks func(
		userRepo *repositoryMocks.ChatRepositoryMock,
		messageRepo *repositoryMocks.MessageRepositoryMock,
		logRepo *repositoryMocks.LogRepositoryMock,
		txManager *dbTxMocks.TxManagerMock,
	)

	type args struct {
		req *model.SendMessage
	}

	var (
		ctxValue = context.Background()
		mc       = minimock.NewController(t)

		id   = gofakeit.Int64()
		from = gofakeit.Name()
		text = gofakeit.Color()

		repositoryErr  = fmt.Errorf("sendMessage messageRepo error")
		logErr         = fmt.Errorf("sendMessage log error")
		transactionErr = fmt.Errorf("transaction error")

		logEntry = &model.LogEntry{
			ChatID: id,
			Activity: fmt.Sprintf(
				"Send message to chat: ChatID:%d, From:%s, Text:%s",
				id,
				from,
				text,
			),
		}

		message = &model.SendMessage{
			ChatID: id,
			Message: model.Message{
				From: from,
				Text: text,
			},
		}

		req = &model.SendMessage{
			ChatID: id,
			Message: model.Message{
				From: from,
				Text: text,
			},
		}
	)

	tests := []struct {
		name       string
		setupMocks setupMocks
		args       args
		err        error
	}{
		{
			name: "success case",
			err:  nil,
			args: args{
				req: req,
			},
			setupMocks: func(_ *repositoryMocks.ChatRepositoryMock,
				messageRepo *repositoryMocks.MessageRepositoryMock,
				logRepo *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				messageRepo.SendMessageMock.Expect(ctxValue, message).Return(nil)
				logRepo.CreateMock.Expect(ctxValue, logEntry).Return(nil)
				txManager.ReadCommitedMock.Set(func(ctx context.Context, f repositoryTx.Handler) error {
					return f(ctx)
				})
			},
		},
		{
			name: "transaction error",
			err:  transactionErr,
			args: args{
				req: req,
			},
			setupMocks: func(_ *repositoryMocks.ChatRepositoryMock,
				_ *repositoryMocks.MessageRepositoryMock,
				_ *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				txManager.ReadCommitedMock.Set(func(_ context.Context, _ repositoryTx.Handler) error {
					return transactionErr
				})
			},
		},
		{
			name: "messageRepo error",
			err:  repositoryErr,
			args: args{
				req: req,
			},
			setupMocks: func(_ *repositoryMocks.ChatRepositoryMock,
				messageRepo *repositoryMocks.MessageRepositoryMock,
				_ *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				messageRepo.SendMessageMock.Expect(ctxValue, message).Return(repositoryErr)
				txManager.ReadCommitedMock.Set(func(ctxValue context.Context, f repositoryTx.Handler) error {
					return f(ctxValue)
				})
			},
		},
		{
			name: "creating log in db error",
			err:  logErr,
			args: args{
				req: req,
			},
			setupMocks: func(_ *repositoryMocks.ChatRepositoryMock,
				messageRepo *repositoryMocks.MessageRepositoryMock,
				logRepo *repositoryMocks.LogRepositoryMock,
				txManager *dbTxMocks.TxManagerMock,
			) {
				messageRepo.SendMessageMock.Expect(ctxValue, message).Return(nil)
				logRepo.CreateMock.Expect(ctxValue, logEntry).Return(logErr)
				txManager.ReadCommitedMock.Set(func(ctxValue context.Context, f repositoryTx.Handler) error {
					return f(ctxValue)
				})
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepo := repositoryMocks.NewChatRepositoryMock(mc)
			messageRepo := repositoryMocks.NewMessageRepositoryMock(mc)
			logRepo := repositoryMocks.NewLogRepositoryMock(mc)
			txManager := dbTxMocks.NewTxManagerMock(mc)

			// Настройка моков в соответствии с тестами
			tt.setupMocks(userRepo, messageRepo, logRepo, txManager)

			service := chatService.NewService(userRepo, messageRepo, logRepo, txManager)

			err := service.SendMessage(ctxValue, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
