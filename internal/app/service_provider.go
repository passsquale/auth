package app

import (
	"context"
	userApi "github.com/passsquale/auth/internal/api/user"
	"github.com/passsquale/auth/internal/client/db"
	"github.com/passsquale/auth/internal/client/db/pg"
	"github.com/passsquale/auth/internal/client/db/transaction"
	"github.com/passsquale/auth/internal/closer"
	"github.com/passsquale/auth/internal/config"
	"github.com/passsquale/auth/internal/repository"
	userRepository "github.com/passsquale/auth/internal/repository/user"
	"github.com/passsquale/auth/internal/service"
	userService "github.com/passsquale/auth/internal/service/user"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository

	userService service.UserService

	userImpl *userApi.UserImplementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(cl.Close)
		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewUserRepository(s.DBClient(ctx))
	}
	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewUserService(
			s.UserRepository(ctx), s.TxManager(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *userApi.UserImplementation {
	if s.userImpl == nil {
		s.userImpl = userApi.NewUserImplementation(s.UserService(ctx))
	}
	return s.userImpl
}
