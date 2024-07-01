package app

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	accessApi "github.com/passsquale/auth/internal/api/access"
	authApi "github.com/passsquale/auth/internal/api/auth"
	userApi "github.com/passsquale/auth/internal/api/user"
	"github.com/passsquale/auth/internal/client/db"
	"github.com/passsquale/auth/internal/client/db/pg"
	"github.com/passsquale/auth/internal/client/db/transaction"
	"github.com/passsquale/auth/internal/closer"
	"github.com/passsquale/auth/internal/config"
	"github.com/passsquale/auth/internal/repository"
	userRepository "github.com/passsquale/auth/internal/repository/user"
	"github.com/passsquale/auth/internal/service"
	accessService "github.com/passsquale/auth/internal/service/access"
	authService "github.com/passsquale/auth/internal/service/auth"
	userService "github.com/passsquale/auth/internal/service/user"
	"github.com/passsquale/auth/internal/storage"
	cache "github.com/passsquale/auth/internal/storage/redis"
	"github.com/passsquale/auth/internal/utils"
	"log"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig
	redisConfig   config.RedisConfig
	jwtConfig     config.JWTConfig

	redisClient storage.Redis
	dbClient    db.Client
	txManager   db.TxManager

	userRepository repository.UserRepository

	userService   service.UserService
	authService   service.AuthService
	accessService service.AccessService

	userImpl   *userApi.UserImplementation
	accessImpl *accessApi.Implementation
	authImpl   *authApi.Implementation

	accessChecker utils.AccessChecker
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

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
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

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %v", err)
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to get jwt config: %v", err)
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) RedisClient() storage.Redis {
	if s.redisClient == nil {
		cl, err := cache.NewRedisConnection(&redis.Options{
			Addr:     s.RedisConfig().Address(),
			Password: s.RedisConfig().Password(),
			DB:       0,
		})
		if err != nil {
			log.Fatalf("failed to create redis client: %v", err)
		}

		err = cl.Ping()
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}

		closer.Add(cl.Close)

		s.redisClient = cl

		s.routesMigrate()
	}

	return s.redisClient
}

func (s *serviceProvider) routesMigrate() {
	routes := s.RedisConfig().RoutesAccesses()

	for route, roles := range routes {
		rolesJSON, err := json.Marshal(roles)
		if err != nil {
			log.Fatalf("error at json marshal")
		}
		_, err = s.RedisClient().Set(route, rolesJSON, 0).Result()
		if err != nil {
			log.Fatalf("error at migration routes to redis")
		}
	}
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
			s.UserRepository(ctx), s.TxManager(ctx), s.RedisClient())
	}
	return s.userService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(
			s.JWTConfig(),
			s.AccessChecker(ctx))
	}

	return s.accessService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.RedisClient(),
			s.UserRepository(ctx),
			s.JWTConfig(),
		)
	}

	return s.authService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *userApi.UserImplementation {
	if s.userImpl == nil {
		s.userImpl = userApi.NewUserImplementation(s.UserService(ctx))
	}
	return s.userImpl
}

func (s *serviceProvider) AuthImplementation(ctx context.Context) *authApi.Implementation {
	if s.authImpl == nil {
		s.authImpl = authApi.NewImplementation(
			s.AuthService(ctx),
		)
	}

	return s.authImpl
}

func (s *serviceProvider) AccessImplementation(ctx context.Context) *accessApi.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = accessApi.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) AccessChecker(ctx context.Context) utils.AccessChecker {
	if s.accessChecker == nil {
		s.accessChecker = utils.NewRouteAccessChecker(
			s.JWTConfig(),
			s.RedisClient(),
			s.UserRepository(ctx),
		)
	}

	return s.accessChecker
}
