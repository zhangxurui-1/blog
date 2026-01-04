package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/config"
)

var (
	Config     *config.Config
	Log        *zap.Logger
	BlackCache local_cache.Cache
	DB         *gorm.DB
	Redis      redis.Client
	ESClient   *elasticsearch.TypedClient
)
