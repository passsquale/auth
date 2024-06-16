package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5/pgxpool"
	userApi "github.com/passsquale/auth/internal/api/user"
	"github.com/passsquale/auth/internal/config"
	userRepo "github.com/passsquale/auth/internal/repository/user"
	userService "github.com/passsquale/auth/internal/service/user"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config")
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
	userRepo := userRepo.NewRepository(pool)
	userServ := userService.NewUserService(ctx, userRepo)
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, userApi.NewImplementation(userServ))

	log.Printf("server listening at %s", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
