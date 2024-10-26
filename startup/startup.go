package startup

import (
	"news-master/app"
	"news-master/auth"
	"news-master/logger"
	"news-master/repository"
	"sync"
)

var loadEnvOnce sync.Once

func Init() {
	loadEnvOnce.Do(_init)
}

func _init() {
	app.LoadEnvVars()
	logger.InitLogger(app.Config.LogLevel)
	auth.LoadKeys()
	repository.Migrate()
}
