package main

import (
	"context"
	"flag"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/passsquale/auth/internal/config"
	"github.com/passsquale/auth/internal/repository"
	"github.com/passsquale/auth/internal/repository/user"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	usersTable      = "users"
	IDColumn        = "id"
	NameColumn      = "name"
	EmailColumn     = "email"
	RoleColumn      = "role"
	PasswordColumn  = "password"
	CreatedAtColumn = "created_at"
	UpdatedAtColumn = "updated_at"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config")
}

type server struct {
	desc.UnimplementedUserV1Server
	userRepository repository.UserRepository
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.userRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &desc.GetResponse{
		User: user,
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.userRepository.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: id}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	err := s.userRepository.Update(ctx, req.GetWrap())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	err := s.userRepository.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, err
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
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
	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	userRepo := user.NewRepository(pool)
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{userRepository: userRepo})

	log.Printf("server listening at %s", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
