package handlers

import (
	"time"

	"github.com/Abdulazizxoshimov/Hospital/config"
	"github.com/Abdulazizxoshimov/Hospital/internal/repo"
	redis "github.com/Abdulazizxoshimov/Hospital/internal/repo/redisdb"
	"github.com/Abdulazizxoshimov/Hospital/pkg/logger"
	tokens "github.com/Abdulazizxoshimov/Hospital/pkg/token"

	"github.com/casbin/casbin/v2"
)

type HandlerV1 struct {
	Config         config.Config
	Logger         logger.Logger
	ContextTimeout time.Duration
	redisStorage   redis.Cache
	RefreshToken   tokens.JWTHandler
	Enforcer       *casbin.Enforcer
	Service        repo.StorageI
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Config         config.Config
	Logger         logger.Logger
	ContextTimeout time.Duration
	Redis          redis.Cache
	RefreshToken   tokens.JWTHandler
	Enforcer       *casbin.Enforcer
	Service        repo.StorageI
}

// New ...
func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		ContextTimeout: c.ContextTimeout,
		redisStorage:   c.Redis,
		Enforcer:       c.Enforcer,
		RefreshToken:   c.RefreshToken,
		Service:        c.Service,
	}
}
