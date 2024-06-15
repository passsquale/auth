package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/passsquale/auth/internal/config"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	pool *pgxpool.Pool
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
				Role:  desc.Role_USER,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	buildInsert := squirrel.Insert(usersTable).
		PlaceholderFormat(squirrel.Dollar).
		Columns(NameColumn, EmailColumn, RoleColumn, PasswordColumn).
		Values(req.GetInfo().GetName(), req.GetInfo().GetEmail(), req.GetInfo().GetRole(), req.GetInfo().GetPassword()).
		Suffix(fmt.Sprintf("RETURNING %s", IDColumn))
	query, args, err := buildInsert.ToSql()
	if err != nil {
		return nil, err
	}
	var ID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&ID)
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: ID,
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	log.Printf("Update user id: %d  name: %+v", req.GetWrap().GetId(), req.GetWrap().GetName())
	return &empty.Empty{}, nil
}
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("Delete user id: %d", req.GetId())
	return &empty.Empty{}, nil
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

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %s", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
