package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/passsquale/auth/internal/api/user"
	"github.com/passsquale/auth/internal/model"
	"github.com/passsquale/auth/internal/repository"
	repositoryMocks "github.com/passsquale/auth/internal/repository/mocks"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		password  = gofakeit.Password(true, false, false, false, false, 10)
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		role = desc.Role(rand.Intn(3))

		repositoryErr = fmt.Errorf("repository error")

		req = &desc.GetRequest{
			Id: id,
		}

		repositoryRes = &model.User{
			ID: id,
			UserInfo: model.UserInfo{
				Name:  name,
				Email: email,
				Role:  model.Role(role),
			},
			Password:  password,
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}

		res = &desc.GetResponse{
			User: &desc.User{
				Id: id,
				Info: &desc.UserInfo{
					Name:  name,
					Email: email,
					Role:  role,
				},
				Password:  password,
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}
	)

	tests := []struct {
		name               string
		args               args
		want               *desc.GetResponse
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(repositoryRes, nil)
				return mock
			},
		},
		{
			name: "server error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repositoryErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepositoryMock := tt.userRepositoryMock(mc)
			api := user.NewUserImplementation(userRepositoryMock)

			resHandler, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}

}
