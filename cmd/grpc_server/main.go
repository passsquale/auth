package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/passsquale/auth/internal/config"
	"github.com/passsquale/auth/internal/converter"
	"github.com/passsquale/auth/internal/repository/user"
	"github.com/passsquale/auth/internal/service"
	userSrv "github.com/passsquale/auth/internal/service/user"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.userService.Create(ctx, converter.ToUserInfoFromDesc(req))
	if err != nil {
		return nil, err
	}
	log.Printf("inserted user with id: %d", id)
	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := s.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	converteredUser := converter.ToUserFromService(userObj)
	return &desc.GetResponse{
		Id:        converteredUser.Id,
		Info:      converteredUser.Info,
		CreatedAt: converteredUser.CreatedAt,
		UpdatedAt: converteredUser.UpdatedAt,
	}, nil
}
func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("cannot load godotenv: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}
	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := user.NewRepository(pool)
	userSrv := userSrv.NewService(userRepo)
	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterUserV1Server(s, &server{userService: userSrv})
	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
