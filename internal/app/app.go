package app

import (
	"fmt"

	"github.com/Abdulazizxoshimov/Hospital/api"
	"github.com/Abdulazizxoshimov/Hospital/api/server"
	"github.com/Abdulazizxoshimov/Hospital/config"
	repo "github.com/Abdulazizxoshimov/Hospital/internal/repo"
	redisrepo "github.com/Abdulazizxoshimov/Hospital/internal/repo/redisdb"
	"github.com/Abdulazizxoshimov/Hospital/pkg/logger"
	"github.com/Abdulazizxoshimov/Hospital/pkg/storage"

	"net/http"

	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"

	"github.com/casbin/casbin/v2"
)

type App struct {
	Config   config.Config
	Logger   logger.Logger
	DB       *storage.PostgresDB
	server   *http.Server
	Enforcer *casbin.Enforcer
	RedisDB  *storage.RedisDB
	StorageI repo.StorageI
}

func NewApp(cfg config.Config) (*App, error) {
	// init logger
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.App+".log")
	if err != nil {
		return nil, err
	}

	//init redis
	redisdb, err := storage.NewRedis(&cfg)
	if err != nil {
		return nil, err
	}

	//init casbin enforcer
	enforcer, err := casbin.NewEnforcer("./config/auth.conf", "./config/auth.csv")
	if err != nil {
		return nil, err
	}

	// init db
	db, err := storage.New(&cfg)
	if err != nil {
		return nil, err
	}

	storageI := repo.NewStoragePg(db)

	return &App{
		Config:   cfg,
		Logger:   logger,
		RedisDB:  redisdb,
		Enforcer: enforcer,
		StorageI: storageI,
	}, nil
}

func (a *App) Run() error {

	// initialize cache
	cache := redisrepo.NewCache(a.RedisDB)

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: a.Config.Context.Timeout,
		Cache:          cache,
		Enforcer:       a.Enforcer,
		Service:        a.StorageI,
	})

	//for Casbin init
	err := a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	roleManager := a.Enforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl)

	roleManager.AddMatchingFunc("keyMatch", util.KeyMatch)
	roleManager.AddMatchingFunc("keyMatch3", util.KeyMatch3)

	// server init
	a.server, err = server.NewServer(&a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {
	// database connection
	a.DB.Close()

	// zap logger sync
	a.Logger.Sync()
}
