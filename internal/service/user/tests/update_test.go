package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/passsquale/auth/internal/api/user"
	"github.com/passsquale/auth/internal/model"
	"github.com/passsquale/auth/internal/repository"
	repositoryMocks "github.com/passsquale/auth/internal/repository/mocks"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math/rand"
	"testing"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = desc.Role(rand.Intn(3))

		repositoryErr = fmt.Errorf("repository error")

		req = &desc.UpdateRequest{
			Wrap: &desc.Updwrap{
				Id:    id,
				Name:  wrapperspb.String(name),
				Email: wrapperspb.String(email),
				Role:  role,
			},
		}

		repositoryRes = &model.UserUpdate{
			ID: id,
			UserInfo: model.UserInfo{
				Name:  name,
				Email: email,
				Role:  model.Role(role),
			},
		}

		res = &empty.Empty{}
	)

	tests := []struct {
		name               string
		args               args
		want               *empty.Empty
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
				mock.UpdateMock.Expect(ctx, repositoryRes).Return(nil)
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
				mock.UpdateMock.Expect(ctx, repositoryRes).Return(repositoryErr)
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

			resHandler, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}

}
